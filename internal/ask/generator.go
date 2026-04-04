// internal/ask/generator.go
package ask

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

const systemPrompt = `You are a test case generator. Given a natural language description of an API operation, return ONLY a JSON array of test case objects. No explanation, no markdown prose — just the raw JSON array.

Each object must have:
- "title": string
- "kind": "single" or "chain"
- "priority": "P0", "P1", "P2", or "P3"
- "tags": array of strings (can be empty)
- "steps": array of step objects

Each step object must have:
- "title": string
- "method": HTTP method (GET, POST, PUT, PATCH, DELETE)
- "path": URL path (e.g. "/users")
- "body": object or null
- "assertions": array of {"target":"status_code","operator":"eq","expected":<number>}

Generate 3-7 diverse test cases covering happy path and error scenarios.`

const llmCallTimeout = 30 * time.Second

// llmTestCase is the structure the LLM is instructed to return.
type llmTestCase struct {
	Title    string    `json:"title"`
	Kind     string    `json:"kind"`
	Priority string    `json:"priority"`
	Tags     []string  `json:"tags"`
	Steps    []llmStep `json:"steps"`
}

type llmStep struct {
	Title      string             `json:"title"`
	Method     string             `json:"method"`
	Path       string             `json:"path"`
	Headers    map[string]string  `json:"headers,omitempty"`
	Body       any                `json:"body,omitempty"`
	Assertions []schema.Assertion `json:"assertions,omitempty"`
}

// Generator generates test cases from a natural language description using an LLM.
type Generator struct {
	provider llm.LLMProvider
}

// NewGenerator creates a Generator backed by the given LLM provider.
func NewGenerator(provider llm.LLMProvider) *Generator {
	return &Generator{provider: provider}
}

// Generate calls the LLM with the description and returns parsed test cases.
func (g *Generator) Generate(ctx context.Context, description string) ([]schema.TestCase, error) {
	if !g.provider.IsAvailable() {
		return nil, fmt.Errorf("AI provider unavailable: ask requires a configured LLM (set ANTHROPIC_API_KEY or equivalent)")
	}

	// Ensure a deadline: if the caller didn't set one, add our default.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, llmCallTimeout)
		defer cancel()
	}

	req := &llm.CompletionRequest{
		System:    systemPrompt,
		Messages:  []llm.Message{{Role: "user", Content: description}},
		MaxTokens: 4096,
	}
	resp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return g.provider.Complete(ctx, req)
	})
	if err != nil {
		return nil, fmt.Errorf("LLM completion failed: %w", err)
	}

	text := llm.ExtractJSON(resp.Text)

	var raw []llmTestCase
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		return nil, fmt.Errorf("parsing LLM response as JSON: %w", err)
	}

	cases := make([]schema.TestCase, len(raw))
	for i, r := range raw {
		cases[i] = toTestCase(r, description)
	}
	return cases, nil
}

func toTestCase(lc llmTestCase, description string) schema.TestCase {
	kind := lc.Kind
	if kind == "" {
		kind = "single"
	}
	priority := lc.Priority
	if priority == "" {
		priority = "P1"
	}

	steps := make([]schema.Step, len(lc.Steps))
	for i, ls := range lc.Steps {
		steps[i] = schema.Step{
			ID:         fmt.Sprintf("step-%d", i+1),
			Title:      ls.Title,
			Type:       "test",
			Method:     ls.Method,
			Path:       ls.Path,
			Headers:    ls.Headers,
			Body:       ls.Body,
			Assertions: ls.Assertions,
		}
	}

	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
		Title:    lc.Title,
		Kind:     kind,
		Priority: priority,
		Tags:     lc.Tags,
		Source: schema.CaseSource{
			Technique: "ask",
			Rationale: description,
		},
		Steps:       steps,
		GeneratedAt: time.Now(),
	}
}
