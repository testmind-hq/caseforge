// internal/rbt/mapper_test.go
package rbt

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
)

// stubParser is a test double for SourceParser.
type stubParser struct {
	via      string
	routes   []RouteMapping
	claimsIf func(path string) bool
}

func (s *stubParser) ExtractRoutes(ctx context.Context, srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	var out []RouteMapping
	for _, f := range files {
		if s.claimsIf == nil || s.claimsIf(f.Path) {
			for _, r := range s.routes {
				r.SourceFile = f.Path
				r.Via = s.via
				out = append(out, r)
			}
		}
	}
	return out, nil
}

func TestMapChain_FirstParserClaimsAll(t *testing.T) {
	p1 := &stubParser{via: "p1", routes: []RouteMapping{{Method: "GET", RoutePath: "/a"}}, claimsIf: func(string) bool { return true }}
	p2 := &stubParser{via: "p2", routes: []RouteMapping{{Method: "POST", RoutePath: "/b"}}, claimsIf: func(string) bool { return true }}

	files := []ChangedFile{{Path: "handler.go"}}
	mappings, err := MapChain([]SourceParser{p1, p2}, ".", files)
	require.NoError(t, err)
	assert.Len(t, mappings, 1)
	assert.Equal(t, "p1", mappings[0].Via)
}

func TestMapChain_FallsThrough(t *testing.T) {
	p1 := &stubParser{via: "p1", routes: nil, claimsIf: func(string) bool { return false }}
	p2 := &stubParser{via: "p2", routes: []RouteMapping{{Method: "POST", RoutePath: "/b"}}, claimsIf: func(string) bool { return true }}

	files := []ChangedFile{{Path: "handler.go"}}
	mappings, err := MapChain([]SourceParser{p1, p2}, ".", files)
	require.NoError(t, err)
	assert.Len(t, mappings, 1)
	assert.Equal(t, "p2", mappings[0].Via)
}

func TestMapChain_NoFiles_ReturnsEmpty(t *testing.T) {
	p1 := &stubParser{via: "p1", routes: []RouteMapping{{Method: "GET", RoutePath: "/a"}}}
	mappings, err := MapChain([]SourceParser{p1}, ".", nil)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

// ── MapFileParser tests ──────────────────────────────────────────────────────

func TestMapFileParser_ParsesYAML(t *testing.T) {
	dir := t.TempDir()
	mapContent := `
mappings:
  - source: internal/user/service.go
    operations:
      - POST /users
      - GET /users/{id}
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "caseforge-map.yaml"), []byte(mapContent), 0644))

	parser := NewMapFileParser(filepath.Join(dir, "caseforge-map.yaml"))
	files := []ChangedFile{{Path: "internal/user/service.go"}}
	mappings, err := parser.ExtractRoutes(context.Background(), dir, files)
	require.NoError(t, err)
	require.Len(t, mappings, 2)
	assert.Equal(t, "POST", mappings[0].Method)
	assert.Equal(t, "/users", mappings[0].RoutePath)
	assert.Equal(t, "GET", mappings[1].Method)
	assert.Equal(t, "/users/{id}", mappings[1].RoutePath)
	assert.Equal(t, "mapfile", mappings[0].Via)
}

func TestMapFileParser_NoMapFile_ReturnsEmpty(t *testing.T) {
	parser := NewMapFileParser("/nonexistent/caseforge-map.yaml")
	files := []ChangedFile{{Path: "handler.go"}}
	mappings, err := parser.ExtractRoutes(context.Background(), ".", files)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

func TestMapFileParser_UnmatchedFile_Skipped(t *testing.T) {
	dir := t.TempDir()
	mapContent := `
mappings:
  - source: internal/user/service.go
    operations:
      - POST /users
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "caseforge-map.yaml"), []byte(mapContent), 0644))
	parser := NewMapFileParser(filepath.Join(dir, "caseforge-map.yaml"))
	files := []ChangedFile{{Path: "internal/order/handler.go"}}
	mappings, err := parser.ExtractRoutes(context.Background(), dir, files)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

// ── RegexParser tests ─────────────────────────────────────────────────────────

func TestRegexParser_GoGinRoutes(t *testing.T) {
	dir := t.TempDir()
	src := `package handler

func Register(r *gin.Engine) {
    r.GET("/users", ListUsers)
    r.POST("/users", CreateUser)
    r.DELETE("/users/:id", DeleteUser)
}
`
	srcPath := filepath.Join(dir, "handler.go")
	require.NoError(t, os.WriteFile(srcPath, []byte(src), 0644))

	parser := NewRegexParser()
	files := []ChangedFile{{Path: srcPath}}
	mappings, err := parser.ExtractRoutes(context.Background(), dir, files)
	require.NoError(t, err)
	assert.Len(t, mappings, 3)
	methods := make(map[string]string)
	for _, m := range mappings {
		methods[m.Method] = m.RoutePath
	}
	assert.Equal(t, "/users", methods["GET"])
	assert.Equal(t, "/users", methods["POST"])
	assert.Equal(t, "/users/:id", methods["DELETE"])
	assert.Equal(t, "regex", mappings[0].Via)
}

func TestRegexParser_NoRoutes_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "util.go")
	require.NoError(t, os.WriteFile(srcPath, []byte("package util\nfunc Helper() {}\n"), 0644))
	parser := NewRegexParser()
	files := []ChangedFile{{Path: srcPath}}
	mappings, err := parser.ExtractRoutes(context.Background(), dir, files)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

// ── LLMParser tests ───────────────────────────────────────────────────────────

type fakeLLMProvider struct {
	response string
}

func (f *fakeLLMProvider) Complete(_ context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	return &llm.CompletionResponse{Text: f.response}, nil
}
func (f *fakeLLMProvider) IsAvailable() bool { return true }
func (f *fakeLLMProvider) Name() string       { return "fake" }

func TestLLMParser_ParsesJSONResponse(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "user.go")
	require.NoError(t, os.WriteFile(srcPath, []byte("package service\n"), 0644))

	resp := `[{"method":"POST","path":"/users","confidence":0.9}]`
	provider := &fakeLLMProvider{response: resp}
	parser := NewLLMParser(provider, nil)
	files := []ChangedFile{{Path: srcPath}}
	mappings, err := parser.ExtractRoutes(context.Background(), dir, files)
	require.NoError(t, err)
	require.Len(t, mappings, 1)
	assert.Equal(t, "POST", mappings[0].Method)
	assert.Equal(t, "/users", mappings[0].RoutePath)
	assert.InDelta(t, 0.9, mappings[0].Confidence, 0.001)
	assert.Equal(t, "llm", mappings[0].Via)
}

func TestLLMParser_UnavailableProvider_ReturnsEmpty(t *testing.T) {
	parser := NewLLMParser(nil, nil)
	files := []ChangedFile{{Path: "service/user.go"}}
	mappings, err := parser.ExtractRoutes(context.Background(), ".", files)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

// ── TreeSitterParser tests ────────────────────────────────────────────────────

func TestTreeSitterParser_NotInstalled_ReturnsEmpty(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	parser := NewTreeSitterParser()
	files := []ChangedFile{{Path: "handler.go"}}
	mappings, err := parser.ExtractRoutes(context.Background(), ".", files)
	require.NoError(t, err)
	assert.Empty(t, mappings)
}

