// internal/methodology/engine.go
package methodology

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/testmind-hq/caseforge/internal/event"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type Technique interface {
	Name() string
	Applies(op *spec.Operation) bool
	Generate(op *spec.Operation) ([]schema.TestCase, error)
}

// SpecTechnique generates test cases that require cross-operation context.
// Use this for patterns like chain cases that span multiple operations.
type SpecTechnique interface {
	Name() string
	Generate(s *spec.ParsedSpec) ([]schema.TestCase, error)
}

// Seedable is an optional interface techniques can implement to receive a
// deterministic seed for reproducible combination ordering.
type Seedable interface {
	SetSeed(seed int64)
}

type Engine struct {
	techniques     []Technique
	specTechniques []SpecTechnique
	llm            llm.LLMProvider
	sink           event.Sink
	concurrency    int   // 0 or 1 = serial; >1 = parallel worker pool
	seed           int64 // 0 = random
}

func NewEngine(provider llm.LLMProvider, techniques ...Technique) *Engine {
	return &Engine{techniques: techniques, llm: provider}
}

// SetConcurrency sets the number of operations processed concurrently.
// Values ≤1 (including 0) run serially. Values >1 enable a worker pool.
func (e *Engine) SetConcurrency(n int) {
	e.concurrency = n
}

// AddSpecTechnique registers a spec-level technique with the engine.
func (e *Engine) AddSpecTechnique(t SpecTechnique) {
	e.specTechniques = append(e.specTechniques, t)
}

// SetSink registers an event sink for progress events.
func (e *Engine) SetSink(s event.Sink) {
	e.sink = s
}

// SetSeed seeds the global data generator and all Seedable techniques for
// deterministic output. Call before Generate.
func (e *Engine) SetSeed(seed int64) {
	e.seed = seed
	gofakeit.Seed(seed)
	for _, t := range e.techniques {
		if s, ok := t.(Seedable); ok {
			s.SetSeed(seed)
		}
	}
}

func (e *Engine) emit(ev event.Event) {
	if e.sink != nil {
		e.sink.Emit(ev)
	}
}

// Generate annotates all operations with LLM semantic info, then
// dispatches each operation to applicable techniques.
// When concurrency > 1, operations are processed in parallel via a worker pool.
func (e *Engine) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	// Step 1: Enrich all operations with LLM semantic annotations
	e.annotateOperations(s.Operations)

	// Step 2: For each operation, apply all applicable techniques
	allCases, err := e.generateOperations(s.Operations)
	if err != nil {
		return nil, err
	}

	// Step 3: Apply spec-level techniques (cross-operation, e.g. chain cases)
	for _, tech := range e.specTechniques {
		cases, err := tech.Generate(s)
		if err != nil {
			return nil, fmt.Errorf("spec technique %s: %w", tech.Name(), err)
		}
		allCases = append(allCases, cases...)
	}
	return allCases, nil
}

// generateOperations runs all per-operation techniques, using a worker pool
// when concurrency > 1. Results are returned in the same order as ops.
func (e *Engine) generateOperations(ops []*spec.Operation) ([]schema.TestCase, error) {
	n := e.concurrency
	if n <= 1 {
		// Serial path — unchanged behaviour.
		var allCases []schema.TestCase
		for _, op := range ops {
			cases, err := e.generateOne(op)
			if err != nil {
				return nil, err
			}
			allCases = append(allCases, cases...)
		}
		return allCases, nil
	}

	// Parallel path: each goroutine writes to its own slot (no mutex needed on
	// results slice). event.Bus.Emit is thread-safe so e.emit() is safe here.
	type opResult struct {
		cases []schema.TestCase
		err   error
	}
	results := make([]opResult, len(ops))
	sem := make(chan struct{}, n)
	var wg sync.WaitGroup

	for i, op := range ops {
		wg.Add(1)
		go func(idx int, op *spec.Operation) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			cases, err := e.generateOne(op)
			results[idx] = opResult{cases: cases, err: err}
		}(i, op)
	}
	wg.Wait()

	var allCases []schema.TestCase
	for _, r := range results {
		if r.err != nil {
			return nil, r.err
		}
		allCases = append(allCases, r.cases...)
	}
	return allCases, nil
}

// generateOne runs all per-operation techniques for a single operation.
func (e *Engine) generateOne(op *spec.Operation) ([]schema.TestCase, error) {
	var cases []schema.TestCase
	for _, tech := range e.techniques {
		if !tech.Applies(op) {
			continue
		}
		c, err := tech.Generate(op)
		if err != nil {
			return nil, fmt.Errorf("technique %s on %s %s: %w",
				tech.Name(), op.Method, op.Path, err)
		}
		for range c {
			e.emit(event.Event{Type: event.EventCaseGenerated})
		}
		cases = append(cases, c...)
	}
	e.emit(event.Event{Type: event.EventOperationDone, Payload: event.OperationDonePayload{
		OperationID: op.OperationID,
		Method:      op.Method,
		Path:        op.Path,
		CaseCount:   len(cases),
	}})
	return cases, nil
}

func (e *Engine) annotateOperations(ops []*spec.Operation) {
	if !e.llm.IsAvailable() {
		return // NoopProvider: skip annotation, SemanticInfo stays nil
	}
	for _, op := range ops {
		annotation, err := e.annotateOperation(op)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: LLM annotation failed for %s %s: %v\n", op.Method, op.Path, err)
			continue
		}
		op.SemanticInfo = annotation
	}
}

func (e *Engine) annotateOperation(op *spec.Operation) (*spec.SemanticAnnotation, error) {
	prompt := fmt.Sprintf(
		"Analyze this API operation and return JSON:\n"+
			"Operation: %s %s\nSummary: %s\nDescription: %s\n"+
			"Return: {resource_type, action_type, has_state_machine, state_field, unique_fields, implicit_rules}",
		op.Method, op.Path, op.Summary, op.Description,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &llm.CompletionRequest{
		System:    "You are an API testing expert. Analyze operations and return structured JSON.",
		Messages:  []llm.Message{{Role: "user", Content: prompt}},
		MaxTokens: 512,
	}
	resp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return e.llm.Complete(ctx, req)
	})
	if err != nil {
		return nil, err
	}
	if resp.Text == "" {
		return nil, nil
	}
	return parseSemanticAnnotation(resp.Text), nil
}

func parseSemanticAnnotation(text string) *spec.SemanticAnnotation {
	extracted := llm.ExtractJSON(text)
	var raw struct {
		ResourceType    string   `json:"resource_type"`
		ActionType      string   `json:"action_type"`
		HasStateMachine bool     `json:"has_state_machine"`
		StateField      string   `json:"state_field"`
		UniqueFields    []string `json:"unique_fields"`
		ImplicitRules   []string `json:"implicit_rules"`
	}
	if err := json.Unmarshal([]byte(extracted), &raw); err != nil {
		return nil
	}
	return &spec.SemanticAnnotation{
		ResourceType:    raw.ResourceType,
		ActionType:      raw.ActionType,
		HasStateMachine: raw.HasStateMachine,
		StateField:      raw.StateField,
		UniqueFields:    raw.UniqueFields,
		ImplicitRules:   raw.ImplicitRules,
	}
}
