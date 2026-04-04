// internal/methodology/engine.go
package methodology

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

type Engine struct {
	techniques     []Technique
	specTechniques []SpecTechnique
	llm            llm.LLMProvider
	sink           event.Sink
}

func NewEngine(provider llm.LLMProvider, techniques ...Technique) *Engine {
	return &Engine{techniques: techniques, llm: provider}
}

// AddSpecTechnique registers a spec-level technique with the engine.
func (e *Engine) AddSpecTechnique(t SpecTechnique) {
	e.specTechniques = append(e.specTechniques, t)
}

// SetSink registers an event sink for progress events.
func (e *Engine) SetSink(s event.Sink) {
	e.sink = s
}

func (e *Engine) emit(ev event.Event) {
	if e.sink != nil {
		e.sink.Emit(ev)
	}
}

// Generate annotates all operations with LLM semantic info, then
// dispatches each operation to applicable techniques.
func (e *Engine) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	// Step 1: Enrich all operations with LLM semantic annotations
	e.annotateOperations(s.Operations)

	// Step 2: For each operation, apply all applicable techniques
	var allCases []schema.TestCase
	for _, op := range s.Operations {
		for _, tech := range e.techniques {
			if !tech.Applies(op) {
				continue
			}
			cases, err := tech.Generate(op)
			if err != nil {
				return nil, fmt.Errorf("technique %s on %s %s: %w",
					tech.Name(), op.Method, op.Path, err)
			}
			for range cases {
				e.emit(event.Event{Type: event.EventCaseGenerated})
			}
			allCases = append(allCases, cases...)
		}
		e.emit(event.Event{Type: event.EventOperationDone, Payload: op.Path})
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
