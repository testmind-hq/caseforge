// internal/rbt/callgraph_test.go
package rbt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCallGraphBuilder returns pre-canned defs and calls for specific file paths.
type mockCallGraphBuilder struct {
	data map[string]struct {
		defs  []CallNode
		calls []CallEdge
	}
}

func (m *mockCallGraphBuilder) ExtractFuncs(filePath string) ([]CallNode, []CallEdge, error) {
	if d, ok := m.data[filePath]; ok {
		return d.defs, d.calls, nil
	}
	return nil, nil, nil
}

func TestBuildCallGraph_BuildsInvertedEdges(t *testing.T) {
	builder := &mockCallGraphBuilder{
		data: map[string]struct {
			defs  []CallNode
			calls []CallEdge
		}{
			"/app/handler.go": {
				defs:  []CallNode{{File: "/app/handler.go", FuncName: "Register"}},
				calls: []CallEdge{{CallerFile: "/app/handler.go", CallerFunc: "Register", CalleeName: "CreateUser"}},
			},
			"/app/service.go": {
				defs:  []CallNode{{File: "/app/service.go", FuncName: "CreateUser"}},
				calls: []CallEdge{{CallerFile: "/app/service.go", CallerFunc: "CreateUser", CalleeName: "validate"}},
			},
			"/app/utils.go": {
				defs:  []CallNode{{File: "/app/utils.go", FuncName: "validate"}},
				calls: nil,
			},
		},
	}

	files := []ChangedFile{
		{Path: "/app/handler.go"},
		{Path: "/app/service.go"},
		{Path: "/app/utils.go"},
	}

	cg, defsByFile := BuildCallGraph(files, builder)
	_ = defsByFile
	require.NotNil(t, cg)

	callers := cg.Edges[CallNodeKey("/app/utils.go", "validate")]
	require.Len(t, callers, 1)
	assert.Equal(t, "CreateUser", callers[0].FuncName)
	assert.Equal(t, "/app/service.go", callers[0].File)

	callers2 := cg.Edges[CallNodeKey("/app/service.go", "CreateUser")]
	require.Len(t, callers2, 1)
	assert.Equal(t, "Register", callers2[0].FuncName)
}

func TestTraceToRoutes_FindsRouteAcross3Layers(t *testing.T) {
	cg := &CallGraph{
		Edges: map[string][]CallNode{
			CallNodeKey("/app/utils.go", "validate"):     {{File: "/app/service.go", FuncName: "CreateUser"}},
			CallNodeKey("/app/service.go", "CreateUser"): {{File: "/app/handler.go", FuncName: "Register"}},
		},
	}
	routeFiles := map[string][]RouteMapping{
		"/app/handler.go": {
			{SourceFile: "/app/handler.go", Method: "POST", RoutePath: "/users", Via: "treesitter"},
		},
	}
	startNodes := []CallNode{{File: "/app/utils.go", FuncName: "validate"}}

	result, covered := TraceToRoutes(cg, startNodes, routeFiles, 0, "callgraph", 0.8)
	require.Len(t, result, 1)
	assert.Equal(t, "POST", result[0].Method)
	assert.Equal(t, "/users", result[0].RoutePath)
	assert.Equal(t, "callgraph", result[0].Via)
	assert.InDelta(t, 0.8, result[0].Confidence, 0.001)
	assert.True(t, covered["/app/utils.go"], "utils.go should be covered")
}

func TestTraceToRoutes_DepthCapPreventsReachingRoute(t *testing.T) {
	cg := &CallGraph{
		Edges: map[string][]CallNode{
			CallNodeKey("/app/utils.go", "validate"):     {{File: "/app/service.go", FuncName: "CreateUser"}},
			CallNodeKey("/app/service.go", "CreateUser"): {{File: "/app/handler.go", FuncName: "Register"}},
		},
	}
	routeFiles := map[string][]RouteMapping{
		"/app/handler.go": {{Method: "POST", RoutePath: "/users", Via: "treesitter"}},
	}
	startNodes := []CallNode{{File: "/app/utils.go", FuncName: "validate"}}

	result, covered := TraceToRoutes(cg, startNodes, routeFiles, 1, "callgraph", 0.8)
	assert.Empty(t, result)
	assert.Empty(t, covered)
}

func TestTraceToRoutes_CycleDoesNotHang(t *testing.T) {
	cg := &CallGraph{
		Edges: map[string][]CallNode{
			CallNodeKey("/app/a.go", "A"): {{File: "/app/b.go", FuncName: "B"}},
			CallNodeKey("/app/b.go", "B"): {{File: "/app/c.go", FuncName: "C"}},
			CallNodeKey("/app/c.go", "C"): {{File: "/app/a.go", FuncName: "A"}},
		},
	}
	routeFiles := map[string][]RouteMapping{}
	startNodes := []CallNode{{File: "/app/a.go", FuncName: "A"}}

	result, covered := TraceToRoutes(cg, startNodes, routeFiles, 0, "callgraph", 0.8)
	assert.Empty(t, result)
	assert.Empty(t, covered)
}

func TestSubtractFiles_RemovesClaimed(t *testing.T) {
	all := []ChangedFile{{Path: "/a.go"}, {Path: "/b.go"}, {Path: "/c.go"}}
	claimed := []RouteMapping{{SourceFile: "/a.go"}, {SourceFile: "/b.go"}}
	remaining := subtractFiles(all, claimed)
	require.Len(t, remaining, 1)
	assert.Equal(t, "/c.go", remaining[0].Path)
}

func TestTraceToRoutes_LLMViaAndConfidence(t *testing.T) {
	cg := &CallGraph{
		Edges: map[string][]CallNode{
			CallNodeKey("/app/service.go", "CreateUser"): {{File: "/app/handler.go", FuncName: "Register"}},
		},
	}
	routeFiles := map[string][]RouteMapping{
		"/app/handler.go": {
			{SourceFile: "/app/handler.go", Method: "POST", RoutePath: "/users", Via: "treesitter"},
		},
	}
	startNodes := []CallNode{{File: "/app/service.go", FuncName: "CreateUser"}}

	result, covered := TraceToRoutes(cg, startNodes, routeFiles, 0, "callgraph-llm", 0.65)
	require.Len(t, result, 1)
	assert.Equal(t, "callgraph-llm", result[0].Via)
	assert.InDelta(t, 0.65, result[0].Confidence, 0.001)
	assert.True(t, covered["/app/service.go"])
}

func TestFallbackCallGraphBuilder_UsesLLMWhenPrimaryEmpty(t *testing.T) {
	// primary returns empty; fallback (LLM) returns defs
	primary := &mockCallGraphBuilder{data: map[string]struct {
		defs  []CallNode
		calls []CallEdge
	}{}}
	llmFallback := &mockCallGraphBuilder{
		data: map[string]struct {
			defs  []CallNode
			calls []CallEdge
		}{
			"/app/service.go": {
				defs:  []CallNode{{File: "/app/service.go", FuncName: "CreateUser"}},
				calls: nil,
			},
		},
	}

	fb := &fallbackCallGraphBuilder{primary: primary, fallback: llmFallback}
	assert.False(t, fb.hasUsedLLM, "should not have used LLM yet")

	defs, _, err := fb.ExtractFuncs("/app/service.go")
	require.NoError(t, err)
	assert.Len(t, defs, 1)
	assert.True(t, fb.hasUsedLLM, "should have used LLM after primary returned empty")
}

func TestSubtractChangedFiles(t *testing.T) {
	all := []ChangedFile{{Path: "/a.go"}, {Path: "/b.go"}, {Path: "/c.go"}}
	remove := []ChangedFile{{Path: "/a.go"}, {Path: "/b.go"}}
	remaining := subtractChangedFiles(all, remove)
	require.Len(t, remaining, 1)
	assert.Equal(t, "/c.go", remaining[0].Path)
}

func TestTraceToRoutes_CoveredFilesTracksOriginFile(t *testing.T) {
	// 3-layer chain: utils.go → service.go → handler.go (route)
	cg := &CallGraph{
		Edges: map[string][]CallNode{
			CallNodeKey("/app/utils.go", "validate"):     {{File: "/app/service.go", FuncName: "CreateUser"}},
			CallNodeKey("/app/service.go", "CreateUser"): {{File: "/app/handler.go", FuncName: "Register"}},
		},
	}
	routeFiles := map[string][]RouteMapping{
		"/app/handler.go": {{SourceFile: "/app/handler.go", Method: "POST", RoutePath: "/users", Via: "treesitter"}},
	}
	// Start from utils.go — 3 hops to route
	startNodes := []CallNode{{File: "/app/utils.go", FuncName: "validate"}}

	_, covered := TraceToRoutes(cg, startNodes, routeFiles, 0, "callgraph", 0.8)
	assert.True(t, covered["/app/utils.go"], "utils.go is the origin file and should be covered")
	assert.False(t, covered["/app/service.go"], "service.go was not a start node")
}
