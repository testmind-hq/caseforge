// internal/rbt/callgraph_go_test.go
package rbt

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// goFixtureDir returns the absolute path to testdata/callgraph_go.
func goFixtureDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	dir := filepath.Join(filepath.Dir(thisFile), "testdata", "callgraph_go")
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("testdata/callgraph_go not found: %v", err)
	}
	return dir
}

func TestGoCallGraphBuilder_RTA_TracesInterface(t *testing.T) {
	dir := goFixtureDir(t)

	handlerFile := filepath.Join(dir, "handler", "handler.go")
	serviceFile := filepath.Join(dir, "service", "service.go")

	// routeFileMappings mocks what runTreeSitterPhase would produce for handler.go
	routeFiles := map[string][]RouteMapping{
		handlerFile: {
			{
				SourceFile: handlerFile,
				Method:     "POST",
				RoutePath:  "/users",
				Via:        "treesitter",
				Confidence: 1.0,
			},
		},
	}

	// unclaimed: service.go was changed but not directly registering routes
	unclaimed := []ChangedFile{{Path: serviceFile}}

	b := &GoCallGraphBuilder{SrcDir: dir, Algo: "rta"}
	mappings, claimed, err := b.BuildAndTrace(unclaimed, routeFiles, 0)

	require.NoError(t, err)
	require.Len(t, mappings, 1, "should find POST /users via interface dispatch")
	assert.Equal(t, "POST", mappings[0].Method)
	assert.Equal(t, "/users", mappings[0].RoutePath)
	assert.Equal(t, "go-callgraph", mappings[0].Via)
	assert.InDelta(t, 0.9, mappings[0].Confidence, 0.001)
	require.Len(t, claimed, 1)
	assert.Equal(t, serviceFile, claimed[0].Path)
}

func TestGoCallGraphBuilder_Fallback_WhenNoGoMod(t *testing.T) {
	// A directory without go.mod should return empty without error.
	dir := t.TempDir()
	b := &GoCallGraphBuilder{SrcDir: dir, Algo: "rta"}
	mappings, claimed, err := b.BuildAndTrace(
		[]ChangedFile{{Path: filepath.Join(dir, "service.go")}},
		map[string][]RouteMapping{},
		0,
	)
	require.NoError(t, err)
	assert.Empty(t, mappings)
	assert.Empty(t, claimed)
}

func TestGoCallGraphBuilder_DepthCap(t *testing.T) {
	dir := goFixtureDir(t)

	handlerFile := filepath.Join(dir, "handler", "handler.go")
	serviceFile := filepath.Join(dir, "service", "service.go")

	routeFiles := map[string][]RouteMapping{
		handlerFile: {{SourceFile: handlerFile, Method: "POST", RoutePath: "/users", Via: "treesitter"}},
	}
	unclaimed := []ChangedFile{{Path: serviceFile}}

	// maxDepth=1: service.go → handler.go is multiple hops, should not reach with depth=1
	b := &GoCallGraphBuilder{SrcDir: dir, Algo: "rta"}
	mappings, _, err := b.BuildAndTrace(unclaimed, routeFiles, 1)

	require.NoError(t, err)
	assert.Empty(t, mappings, "depth cap of 1 should prevent reaching the route handler")
}

func TestGoCallGraphBuilder_EmptyUnclaimedGoFiles(t *testing.T) {
	dir := goFixtureDir(t)
	b := &GoCallGraphBuilder{SrcDir: dir, Algo: "rta"}

	// Only non-Go files — should return immediately without loading packages
	unclaimed := []ChangedFile{{Path: "/some/file.py"}, {Path: "/other/file.ts"}}
	mappings, claimed, err := b.BuildAndTrace(unclaimed, map[string][]RouteMapping{}, 0)

	require.NoError(t, err)
	assert.Empty(t, mappings)
	assert.Empty(t, claimed)
}
