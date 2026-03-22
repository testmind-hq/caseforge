// internal/methodology/engine.go
package methodology

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type Technique interface {
	Name() string
	Applies(op *spec.Operation) bool
	Generate(op *spec.Operation) ([]schema.TestCase, error)
}

type Engine struct {
	techniques []Technique
	llm        llm.LLMProvider
}

func NewEngine(provider llm.LLMProvider, techniques ...Technique) *Engine {
	return &Engine{techniques: techniques, llm: provider}
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
			allCases = append(allCases, cases...)
		}
	}
	return allCases, nil
}

func (e *Engine) annotateOperations(ops []*spec.Operation) {
	if !e.llm.IsAvailable() {
		return // NoopProvider: skip annotation, SemanticInfo stays nil
	}
	// For each operation, call LLM to get semantic annotation
	for _, op := range ops {
		annotation, err := e.annotateOperation(op)
		if err != nil {
			// Per-operation annotation failure is non-fatal — degrade gracefully
			continue
		}
		op.SemanticInfo = annotation
	}
}

func (e *Engine) annotateOperation(op *spec.Operation) (*spec.SemanticAnnotation, error) {
	// Build prompt from template
	prompt := fmt.Sprintf(
		"Analyze this API operation and return JSON:\n"+
			"Operation: %s %s\nSummary: %s\nDescription: %s\n"+
			"Return: {resource_type, action_type, has_state_machine, state_field, unique_fields, implicit_rules}",
		op.Method, op.Path, op.Summary, op.Description,
	)
	resp, err := e.llm.Complete(context.Background(), &llm.CompletionRequest{
		System:    "You are an API testing expert. Analyze operations and return structured JSON.",
		Messages:  []llm.Message{{Role: "user", Content: prompt}},
		MaxTokens: 512,
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
	text = strings.TrimSpace(text)
	// Strip markdown code fence wrappers (```json ... ``` or ``` ... ```)
	if strings.HasPrefix(text, "```") {
		if idx := strings.Index(text, "\n"); idx != -1 {
			text = text[idx+1:]
		}
		text = strings.TrimSuffix(strings.TrimSpace(text), "```")
		text = strings.TrimSpace(text)
	}
	// Parse JSON response from LLM
	var raw struct {
		ResourceType    string   `json:"resource_type"`
		ActionType      string   `json:"action_type"`
		HasStateMachine bool     `json:"has_state_machine"`
		StateField      string   `json:"state_field"`
		UniqueFields    []string `json:"unique_fields"`
		ImplicitRules   []string `json:"implicit_rules"`
	}
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
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
