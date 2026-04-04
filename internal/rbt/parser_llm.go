// internal/rbt/parser_llm.go
package rbt

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/testmind-hq/caseforge/internal/llm"
	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

// llmConcurrency is the maximum number of parallel LLM calls in ExtractRoutes.
const llmConcurrency = 4

type llmRouteResult struct {
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	Confidence float64 `json:"confidence"`
}

// specCandidate is the compact route representation injected into the LLM prompt.
type specCandidate struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	OperationID string `json:"operationId,omitempty"`
	Summary     string `json:"summary,omitempty"`
}

// LLMParser infers route mappings from source files using an LLM.
//
// Improvements over the original serial implementation:
//   - R-1: Files are processed concurrently (up to llmConcurrency parallel calls).
//   - R-2: A structured spec route candidate list is injected into the prompt
//     instead of raw YAML, improving path-matching accuracy.
//   - R-3: Results are cached by SHA-256 of file content; identical content in
//     different files (or across re-runs within one session) hits the cache.
type LLMParser struct {
	provider       llm.LLMProvider
	specCandidates string   // compact JSON route list for prompt injection (R-2)
	cache          sync.Map // sha256(content) → []llmRouteResult            (R-3)
}

// NewLLMParser creates an LLMParser. ops is the parsed spec operations list
// used to build the route candidate hint injected into each prompt.
// Pass nil ops when no spec is available; the prompt will omit the candidate list.
func NewLLMParser(provider llm.LLMProvider, ops []*specpkg.Operation) *LLMParser {
	return &LLMParser{
		provider:       provider,
		specCandidates: buildSpecCandidates(ops),
	}
}

// buildSpecCandidates converts parsed spec operations into a compact JSON string
// suitable for prompt injection as a structured route candidate list.
func buildSpecCandidates(ops []*specpkg.Operation) string {
	if len(ops) == 0 {
		return "[]"
	}
	candidates := make([]specCandidate, 0, len(ops))
	for _, op := range ops {
		candidates = append(candidates, specCandidate{
			Method:      strings.ToUpper(op.Method),
			Path:        op.Path,
			OperationID: op.OperationID,
			Summary:     op.Summary,
		})
	}
	b, err := json.Marshal(candidates)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// ExtractRoutes processes files concurrently (R-1), returning route mappings
// inferred by the LLM. Files whose content has been seen before are served
// from the in-memory cache (R-3).
func (p *LLMParser) ExtractRoutes(ctx context.Context, srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	if p.provider == nil || !p.provider.IsAvailable() {
		return nil, nil
	}

	// Read all file contents upfront; empty string signals read failure.
	contents := make([]string, len(files))
	for i, f := range files {
		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		contents[i] = string(data)
	}

	// Pre-allocate per-file result slots — no mutex needed on the results slice.
	type perFileResult struct {
		routes []llmRouteResult
		err    error
	}
	results := make([]perFileResult, len(files))

	// Bounded concurrent fan-out (R-1).
	sem := make(chan struct{}, llmConcurrency)
	var wg sync.WaitGroup
	for i, f := range files {
		if contents[i] == "" {
			continue
		}
		wg.Add(1)
		i, f := i, f // capture loop vars
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			routes, err := p.inferRoutes(callCtx, f.Path, contents[i])
			results[i] = perFileResult{routes: routes, err: err}
		}()
	}
	wg.Wait()

	// Collect results in original file order.
	var mappings []RouteMapping
	for i, f := range files {
		r := results[i]
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "warn: LLM route inference failed for %s: %v\n", f.Path, r.err)
			continue
		}
		for _, route := range r.routes {
			mappings = append(mappings, RouteMapping{
				SourceFile: f.Path,
				Method:     strings.ToUpper(route.Method),
				RoutePath:  route.Path,
				Via:        "llm",
				Confidence: route.Confidence,
			})
		}
	}
	return mappings, nil
}

// contentHash returns the hex-encoded SHA-256 of the given string.
func contentHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", sum)
}

// inferRoutes asks the LLM which spec routes the given file implements or calls.
// Results are cached by content hash (R-3) to avoid redundant LLM calls.
func (p *LLMParser) inferRoutes(ctx context.Context, filePath, content string) ([]llmRouteResult, error) {
	// R-3: check cache before calling LLM.
	hash := contentHash(content)
	if cached, ok := p.cache.Load(hash); ok {
		return cached.([]llmRouteResult), nil
	}

	prompt := fmt.Sprintf(`Given this source file, which API operations from the spec does it likely implement or call?

File: %s
`+"```"+`
%s
`+"```"+`

Available operations (use paths exactly as listed):
%s

Return JSON array only, no explanation:
[{"method":"POST","path":"/users","confidence":0.9}]

Rules:
- method: uppercase HTTP method
- path: exact path from the operations list above (with {param} placeholders)
- confidence: 0.0-1.0
- Return [] if no operations are clearly affected`, filePath, content, p.specCandidates)

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
	var routeResults []llmRouteResult
	if err := json.Unmarshal([]byte(text), &routeResults); err != nil {
		return nil, fmt.Errorf("parsing LLM route response for %s: %w", filePath, err)
	}

	// Store in cache before returning.
	p.cache.Store(hash, routeResults)
	return routeResults, nil
}
