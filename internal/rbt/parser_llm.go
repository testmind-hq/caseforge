// internal/rbt/parser_llm.go
package rbt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/testmind-hq/caseforge/internal/llm"
)

type llmRouteResult struct {
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Confidence float64 `json:"confidence"`
}

type LLMParser struct {
	provider llm.LLMProvider
	specYAML string
}

func NewLLMParser(provider llm.LLMProvider, specYAML string) *LLMParser {
	return &LLMParser{provider: provider, specYAML: specYAML}
}

func (p *LLMParser) ExtractRoutes(srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	if p.provider == nil || !p.provider.IsAvailable() {
		return nil, nil
	}

	var mappings []RouteMapping
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, f := range files {
		content, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		routes, err := p.inferRoutes(ctx, f.Path, string(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: LLM route inference failed for %s: %v\n", f.Path, err)
			continue
		}
		for _, r := range routes {
			mappings = append(mappings, RouteMapping{
				SourceFile: f.Path,
				Method:     strings.ToUpper(r.Method),
				RoutePath:  r.Path,
				Via:        "llm",
				Confidence: r.Confidence,
			})
		}
	}
	return mappings, nil
}

func (p *LLMParser) inferRoutes(ctx context.Context, filePath, content string) ([]llmRouteResult, error) {
	prompt := fmt.Sprintf(`Given this source file and OpenAPI spec, which API operations does this file likely implement or call?

File: %s
`+"```"+`
%s
`+"```"+`

OpenAPI spec (excerpt):
%s

Return JSON array only, no explanation:
[{"method":"POST","path":"/users","confidence":0.9}]

Rules:
- method: uppercase HTTP method
- path: exact path from spec (with {param} placeholders)
- confidence: 0.0-1.0
- Return [] if no operations are clearly affected`, filePath, content, p.specYAML)

	req := &llm.CompletionRequest{
		Messages:  []llm.Message{{Role: "user", Content: prompt}},
		MaxTokens: 512,
	}
	resp, err := llm.Retry(ctx, 3, func() (*llm.CompletionResponse, error) {
		return p.provider.Complete(ctx, req)
	})
	if err != nil {
		return nil, err
	}

	text := llm.ExtractJSON(resp.Text)
	var results []llmRouteResult
	if err := json.Unmarshal([]byte(text), &results); err != nil {
		return nil, fmt.Errorf("parsing LLM route response for %s: %w", filePath, err)
	}
	return results, nil
}
