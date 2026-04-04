// internal/rbt/callgraph_llm.go
package rbt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/testmind-hq/caseforge/internal/llm"
)

type llmCallGraphResponse struct {
	Definitions []string       `json:"definitions"`
	Calls       []llmCallEntry `json:"calls"`
}

type llmCallEntry struct {
	Caller  string   `json:"caller"`
	Callees []string `json:"callees"`
}

// LLMCallGraphBuilder uses an LLM provider to extract function defs and calls.
// Used as a fallback when tree-sitter cannot analyse the file.
type LLMCallGraphBuilder struct {
	provider llm.LLMProvider
}

// NewLLMCallGraphBuilder creates an LLMCallGraphBuilder from an LLMParser.
// Returns a builder with nil provider (no-op) if p is nil.
func NewLLMCallGraphBuilder(p *LLMParser) *LLMCallGraphBuilder {
	if p == nil {
		return &LLMCallGraphBuilder{}
	}
	return &LLMCallGraphBuilder{provider: p.provider}
}

// ExtractFuncs asks the LLM to list function definitions and call relationships.
// Returns empty slices (no error) when the provider is unavailable or parsing fails.
func (b *LLMCallGraphBuilder) ExtractFuncs(filePath string) ([]CallNode, []CallEdge, error) {
	if b.provider == nil || !b.provider.IsAvailable() {
		return nil, nil, nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil
	}

	prompt := fmt.Sprintf(`Given this source file, list all function/method definitions and what they call.
Return JSON only, no explanation:
{"definitions":["FuncA","FuncB"],"calls":[{"caller":"FuncA","callees":["FuncB","helper"]}]}

Rules:
- Use short function/method names only (no package prefix)
- Omit standard library calls (fmt, log, os, http, etc.)
- Return {"definitions":[],"calls":[]} if file has no relevant functions

File: %s
`+"```"+`
%s
`+"```", filePath, string(content))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return b.provider.Complete(ctx, &llm.CompletionRequest{
			Messages:  []llm.Message{{Role: "user", Content: prompt}},
			MaxTokens: 512,
		})
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: LLM call graph inference failed for %s: %v\n", filePath, err)
		return nil, nil, nil
	}

	text := llm.ExtractJSON(resp.Text)

	var parsed llmCallGraphResponse
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		fmt.Fprintf(os.Stderr, "warn: LLM call graph response parse failed for %s: %v\n", filePath, err)
		return nil, nil, nil
	}

	var defs []CallNode
	for _, name := range parsed.Definitions {
		if name != "" {
			defs = append(defs, CallNode{File: filePath, FuncName: name})
		}
	}

	var calls []CallEdge
	for _, entry := range parsed.Calls {
		for _, callee := range entry.Callees {
			if callee != "" {
				calls = append(calls, CallEdge{
					CallerFile: filePath,
					CallerFunc: entry.Caller,
					CalleeName: callee,
				})
			}
		}
	}
	return defs, calls, nil
}
