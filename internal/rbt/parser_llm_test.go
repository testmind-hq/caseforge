// internal/rbt/parser_llm_test.go
package rbt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

// llmStub is a controllable LLMProvider for rbt tests.
type llmStub struct {
	text string
	err  error
}

func (s *llmStub) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &llm.CompletionResponse{Text: s.text}, nil
}
func (s *llmStub) IsAvailable() bool { return true }
func (s *llmStub) Name() string      { return "stub" }

// countingStub records how many times Complete is called (R-3 cache test).
type countingStub struct {
	text  string
	mu    sync.Mutex
	count int
}

func (s *countingStub) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	s.mu.Lock()
	s.count++
	s.mu.Unlock()
	return &llm.CompletionResponse{Text: s.text}, nil
}
func (s *countingStub) IsAvailable() bool { return true }
func (s *countingStub) Name() string      { return "counting-stub" }

// promptCaptureStub records the prompt text from each Complete call (R-2 test).
type promptCaptureStub struct {
	text      string
	mu        sync.Mutex
	prompts   []string
}

func (s *promptCaptureStub) Complete(_ context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	s.mu.Lock()
	if len(req.Messages) > 0 {
		s.prompts = append(s.prompts, req.Messages[0].Content)
	}
	s.mu.Unlock()
	return &llm.CompletionResponse{Text: s.text}, nil
}
func (s *promptCaptureStub) IsAvailable() bool { return true }
func (s *promptCaptureStub) Name() string      { return "capture-stub" }

// --- existing inference tests (updated: nil ops instead of "") ---

func TestLLMParser_InferRoutes_BareJSON(t *testing.T) {
	p := NewLLMParser(&llmStub{
		text: `[{"method":"GET","path":"/users","confidence":0.9}]`,
	}, nil)
	routes, err := p.inferRoutes(context.Background(), "handler.go", "func listUsers() {}")
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/users", routes[0].Path)
	assert.InDelta(t, 0.9, routes[0].Confidence, 0.001)
}

func TestLLMParser_InferRoutes_FencedJSON(t *testing.T) {
	p := NewLLMParser(&llmStub{
		text: "```json\n[{\"method\":\"POST\",\"path\":\"/orders\",\"confidence\":0.8}]\n```",
	}, nil)
	routes, err := p.inferRoutes(context.Background(), "order.go", "func createOrder() {}")
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "POST", routes[0].Method)
}

func TestLLMParser_InferRoutes_TextPreamble(t *testing.T) {
	p := NewLLMParser(&llmStub{
		text: "Here are the routes:\n[{\"method\":\"DELETE\",\"path\":\"/items/{id}\",\"confidence\":0.7}]",
	}, nil)
	routes, err := p.inferRoutes(context.Background(), "item.go", "func deleteItem() {}")
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, "DELETE", routes[0].Method)
}

func TestLLMParser_InferRoutes_EmptyArray(t *testing.T) {
	p := NewLLMParser(&llmStub{text: "[]"}, nil)
	routes, err := p.inferRoutes(context.Background(), "util.go", "func helper() {}")
	require.NoError(t, err)
	assert.Empty(t, routes)
}

func TestLLMParser_InferRoutes_LLMError_ReturnsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel so retry doesn't wait
	p := NewLLMParser(&llmStub{err: errors.New("network error")}, nil)
	_, err := p.inferRoutes(ctx, "handler.go", "")
	require.Error(t, err)
}

func TestLLMParser_InferRoutes_BadJSON_ReturnsError(t *testing.T) {
	p := NewLLMParser(&llmStub{text: "not json"}, nil)
	_, err := p.inferRoutes(context.Background(), "handler.go", "")
	require.Error(t, err)
}

// --- R-3: content hash cache ---

func TestLLMParser_Cache_SameContentSkipsLLM(t *testing.T) {
	stub := &countingStub{text: `[{"method":"GET","path":"/pets","confidence":0.9}]`}
	p := NewLLMParser(stub, nil)
	content := "func listPets() {}"

	r1, err := p.inferRoutes(context.Background(), "a.go", content)
	require.NoError(t, err)
	r2, err := p.inferRoutes(context.Background(), "b.go", content) // same content, different file
	require.NoError(t, err)

	assert.Equal(t, 1, stub.count, "LLM should be called once for identical content")
	assert.Equal(t, r1, r2, "cached result should equal original")
}

func TestLLMParser_Cache_DifferentContentCallsLLM(t *testing.T) {
	stub := &countingStub{text: `[{"method":"GET","path":"/pets","confidence":0.8}]`}
	p := NewLLMParser(stub, nil)

	_, err := p.inferRoutes(context.Background(), "a.go", "func listPets() {}")
	require.NoError(t, err)
	_, err = p.inferRoutes(context.Background(), "b.go", "func createPet() {}")
	require.NoError(t, err)

	assert.Equal(t, 2, stub.count, "LLM should be called once per unique content")
}

// --- R-2: structured spec candidates in prompt ---

func TestLLMParser_PromptContainsSpecCandidates(t *testing.T) {
	stub := &promptCaptureStub{text: "[]"}
	ops := []*specpkg.Operation{
		{Method: "get", Path: "/pets", OperationID: "listPets", Summary: "List all pets"},
		{Method: "post", Path: "/pets", OperationID: "createPet"},
	}
	p := NewLLMParser(stub, ops)
	_, err := p.inferRoutes(context.Background(), "handler.go", "func list() {}")
	require.NoError(t, err)

	require.Len(t, stub.prompts, 1)
	prompt := stub.prompts[0]
	assert.Contains(t, prompt, "/pets")
	assert.Contains(t, prompt, "listPets")
	assert.Contains(t, prompt, "createPet")
	assert.Contains(t, prompt, "List all pets")
}

func TestLLMParser_PromptWithNilOps_ContainsEmptyList(t *testing.T) {
	stub := &promptCaptureStub{text: "[]"}
	p := NewLLMParser(stub, nil)
	_, err := p.inferRoutes(context.Background(), "handler.go", "func list() {}")
	require.NoError(t, err)

	require.Len(t, stub.prompts, 1)
	assert.Contains(t, stub.prompts[0], "[]") // empty candidate list injected
}

func TestLLMParser_BuildSpecCandidates_UppercasesMethod(t *testing.T) {
	ops := []*specpkg.Operation{
		{Method: "get", Path: "/users"},
		{Method: "POST", Path: "/users"},
	}
	result := buildSpecCandidates(ops)
	assert.Contains(t, result, `"GET"`)
	assert.Contains(t, result, `"POST"`)
}

func TestLLMParser_BuildSpecCandidates_NilOps(t *testing.T) {
	assert.Equal(t, "[]", buildSpecCandidates(nil))
}

// --- R-1: concurrent file processing ---

func TestLLMParser_ExtractRoutes_ProcessesAllFiles(t *testing.T) {
	dir := t.TempDir()
	const fileCount = 6 // > llmConcurrency to exercise the semaphore
	files := make([]ChangedFile, fileCount)
	for i := range files {
		path := filepath.Join(dir, fmt.Sprintf("handler%d.go", i))
		// Each file has unique content so cache deduplication does not suppress calls.
		require.NoError(t, os.WriteFile(path, []byte(fmt.Sprintf("func h%d() {}", i)), 0644))
		files[i] = ChangedFile{Path: path}
	}

	stub := &countingStub{text: `[{"method":"GET","path":"/pets","confidence":0.8}]`}
	p := NewLLMParser(stub, nil)
	mappings, err := p.ExtractRoutes(context.Background(), dir, files)
	require.NoError(t, err)
	assert.Len(t, mappings, fileCount)
	// Verify the LLM was called once per unique file (not short-circuited by cache).
	assert.Equal(t, fileCount, stub.count, "LLM should be called once per unique file")
}

func TestLLMParser_ExtractRoutes_NilProvider_ReturnsNil(t *testing.T) {
	p := NewLLMParser(nil, nil)
	mappings, err := p.ExtractRoutes(context.Background(), ".", []ChangedFile{{Path: "x.go"}})
	require.NoError(t, err)
	assert.Nil(t, mappings)
}

func TestLLMParser_ExtractRoutes_SkipsMissingFiles(t *testing.T) {
	stub := &llmStub{text: "[]"}
	p := NewLLMParser(stub, nil)
	// File does not exist — should be skipped, not error
	files := []ChangedFile{{Path: "/nonexistent/path.go"}}
	mappings, err := p.ExtractRoutes(context.Background(), ".", files)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}
