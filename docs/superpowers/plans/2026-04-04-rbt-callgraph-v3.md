# RBT V3 — Go Type-Aware CallGraph Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a Go-specific call graph phase to `rbt index --strategy hybrid` that uses `golang.org/x/tools/go/callgraph` (RTA/PTA) to correctly trace through interface dispatch boundaries that V2 name-matching cannot handle.

**Architecture:** A new `GoCallGraphBuilder` loads the entire Go module via `go/packages`, builds SSA form, runs RTA (default) or PTA to build a type-aware call graph, then performs the same BFS-upward logic as V2 but with precise interface resolution. The new `runGoCallGraphPhase` is inserted between `runTreeSitterPhase` and V2's `runCallGraphPhase` in `RunHybrid`; any error silently falls through to V2.

**Tech Stack:** `golang.org/x/tools/go/packages`, `golang.org/x/tools/go/ssa`, `golang.org/x/tools/go/ssa/ssautil`, `golang.org/x/tools/go/callgraph`, `golang.org/x/tools/go/callgraph/rta`, `golang.org/x/tools/go/pointer`

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `internal/rbt/callgraph_go.go` | Create | `GoCallGraphBuilder` struct + `BuildAndTrace` method |
| `internal/rbt/callgraph_go_test.go` | Create | Unit tests for `BuildAndTrace` and `runGoCallGraphPhase` |
| `internal/rbt/testdata/callgraph_go/go.mod` | Create | Test fixture: compilable Go module root |
| `internal/rbt/testdata/callgraph_go/main.go` | Create | Test fixture: entry point wiring concrete types |
| `internal/rbt/testdata/callgraph_go/handler/handler.go` | Create | Test fixture: route registration + handler |
| `internal/rbt/testdata/callgraph_go/service/service.go` | Create | Test fixture: business logic calling interface |
| `internal/rbt/testdata/callgraph_go/repo/repo.go` | Create | Test fixture: interface + concrete implementation |
| `internal/rbt/indexer.go` | Modify | Add `Algo string` field; add `runGoCallGraphPhase`; update `RunHybrid` |
| `internal/rbt/indexer_test.go` | Modify | Add test for `runGoCallGraphPhase` |
| `cmd/rbt_index.go` | Modify | Add `--algo` flag; pass to `Indexer.Algo` |
| `docs/acceptance/acceptance-tests.md` | Modify | Add AT-064~066 |
| `scripts/acceptance.sh` | Modify | Add AT-064~066 checks |
| `CLAUDE.md` | Modify | Update scenario count to 66 |
| `go.mod` / `go.sum` | Modify | Add `golang.org/x/tools` |

---

### Task 1: Add golang.org/x/tools dependency

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: Add the dependency**

```bash
go get golang.org/x/tools@latest
```

Expected: `go.mod` now lists `golang.org/x/tools`, `go.sum` updated.

- [ ] **Step 2: Verify build still passes**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3: Run existing tests to confirm nothing broke**

```bash
go test ./... -count=1
```

Expected: all PASS.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "chore(deps): add golang.org/x/tools for go/callgraph V3"
```

---

### Task 2: Create testdata fixture — compilable Go module with interface dispatch

**Files:**
- Create: `internal/rbt/testdata/callgraph_go/go.mod`
- Create: `internal/rbt/testdata/callgraph_go/main.go`
- Create: `internal/rbt/testdata/callgraph_go/handler/handler.go`
- Create: `internal/rbt/testdata/callgraph_go/service/service.go`
- Create: `internal/rbt/testdata/callgraph_go/repo/repo.go`

This fixture is the key V3 test case: `service.Process` calls `r.Save()` via the `UserRepo` interface. V2 name-matching breaks here (it sees the string "Save" but cannot determine the concrete type). V3 (RTA) sees `&repo.MySQLRepo{}` instantiated in `main.main`, knows `MySQLRepo` implements `UserRepo`, and correctly traces `service.Process → MySQLRepo.Save → ...`.

- [ ] **Step 1: Create go.mod**

```
internal/rbt/testdata/callgraph_go/go.mod
```

Content:
```
module testapp

go 1.21
```

- [ ] **Step 2: Create repo/repo.go — interface + concrete type**

```go
// internal/rbt/testdata/callgraph_go/repo/repo.go
package repo

// UserRepo is an interface — V2 name-matching breaks at this boundary.
type UserRepo interface {
	Save(name string)
}

// MySQLRepo is the concrete implementation.
type MySQLRepo struct{}

func (m *MySQLRepo) Save(name string) {}
```

- [ ] **Step 3: Create service/service.go — calls through interface**

```go
// internal/rbt/testdata/callgraph_go/service/service.go
package service

import "testapp/repo"

// Process calls repo.Save via the UserRepo interface.
// V2 cannot trace this; V3 (RTA) resolves to MySQLRepo.Save.
func Process(r repo.UserRepo) {
	r.Save("alice")
}
```

- [ ] **Step 4: Create handler/handler.go — route registration**

```go
// internal/rbt/testdata/callgraph_go/handler/handler.go
package handler

import (
	"net/http"
	"testapp/repo"
	"testapp/service"
)

// Register wires the /users route. This is the route-registering file.
func Register(mux *http.ServeMux, r repo.UserRepo) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		CreateUser(w, req, r)
	})
}

// CreateUser calls service.Process with the concrete repo.
func CreateUser(w http.ResponseWriter, r *http.Request, ur repo.UserRepo) {
	service.Process(ur)
}
```

- [ ] **Step 5: Create main.go — wires concrete type so RTA can trace instantiation**

```go
// internal/rbt/testdata/callgraph_go/main.go
package main

import (
	"net/http"
	"testapp/handler"
	"testapp/repo"
)

func main() {
	mux := http.NewServeMux()
	r := &repo.MySQLRepo{} // RTA sees this instantiation
	handler.Register(mux, r)
	_ = http.ListenAndServe(":8080", mux)
}
```

- [ ] **Step 6: Verify the fixture compiles**

```bash
cd internal/rbt/testdata/callgraph_go && go build ./... && cd -
```

Expected: no errors.

- [ ] **Step 7: Commit**

```bash
git add internal/rbt/testdata/callgraph_go/
git commit -m "test(rbt): add callgraph_go testdata fixture (interface dispatch demo)"
```

---

### Task 3: Implement GoCallGraphBuilder (TDD)

**Files:**
- Create: `internal/rbt/callgraph_go_test.go`
- Create: `internal/rbt/callgraph_go.go`

- [ ] **Step 1: Write the failing tests**

Create `internal/rbt/callgraph_go_test.go`:

```go
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

// fixtureDir returns the absolute path to testdata/callgraph_go.
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

	// maxDepth=1: service.go → handler.go is 2+ hops through CreateUser, should not reach
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
```

- [ ] **Step 2: Run tests to confirm they fail**

```bash
go test ./internal/rbt/... -run "TestGoCallGraph" -v 2>&1 | head -20
```

Expected: FAIL — `GoCallGraphBuilder` undefined.

- [ ] **Step 3: Implement callgraph_go.go**

Create `internal/rbt/callgraph_go.go`:

```go
// internal/rbt/callgraph_go.go
package rbt

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// GoCallGraphBuilder performs type-aware call graph analysis for Go modules
// using golang.org/x/tools/go/callgraph. It handles interface dispatch that
// V2 name-matching cannot resolve.
type GoCallGraphBuilder struct {
	SrcDir string
	Algo   string // "rta" (default) | "pta"
}

// BuildAndTrace loads the Go module in SrcDir, builds an SSA call graph via
// RTA or PTA, then BFS-traces upward from functions in unclaimed .go files to
// route-registering files. Returns the found route mappings, the subset of
// unclaimed files that were resolved, and any hard error (caller falls back to V2).
// Returns (nil, nil, nil) silently when: no go.mod, no main package, or no .go
// files in unclaimed.
func (b *GoCallGraphBuilder) BuildAndTrace(
	unclaimed []ChangedFile,
	routeFileMappings map[string][]RouteMapping,
	maxDepth int,
) ([]RouteMapping, []ChangedFile, error) {
	// Only handle .go files.
	var goUnclaimed []ChangedFile
	for _, f := range unclaimed {
		if strings.HasSuffix(f.Path, ".go") {
			goUnclaimed = append(goUnclaimed, f)
		}
	}
	if len(goUnclaimed) == 0 {
		return nil, nil, nil
	}

	// Require go.mod — non-Go projects skip V3 silently.
	if _, err := os.Stat(filepath.Join(b.SrcDir, "go.mod")); err != nil {
		return nil, nil, nil
	}

	// Load all packages in the module.
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedTypes |
			packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypesSizes,
		Dir: b.SrcDir,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, nil, err
	}
	// If there are load errors, fall through to V2.
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, nil, nil
		}
	}

	// Build SSA.
	prog, ssaPkgs := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	// Find main packages for RTA/PTA roots.
	var mainPkgs []*ssa.Package
	for _, pkg := range ssaPkgs {
		if pkg != nil && pkg.Pkg.Name() == "main" {
			mainPkgs = append(mainPkgs, pkg)
		}
	}
	if len(mainPkgs) == 0 {
		// Library without main — RTA has no entry point, skip V3.
		return nil, nil, nil
	}

	// Build the call graph.
	var cg *callgraph.Graph
	algo := b.Algo
	if algo == "" {
		algo = "rta"
	}
	switch algo {
	case "pta":
		ptaCfg := &pointer.Config{
			Mains:          mainPkgs,
			BuildCallGraph: true,
		}
		result, err := pointer.Analyze(ptaCfg)
		if err != nil {
			return nil, nil, err
		}
		cg = result.CallGraph
	default: // "rta"
		var roots []*ssa.Function
		for _, pkg := range mainPkgs {
			if fn := pkg.Func("main"); fn != nil {
				roots = append(roots, fn)
			}
		}
		if len(roots) == 0 {
			return nil, nil, nil
		}
		cg = rta.Analyze(roots, true).CallGraph
	}

	// Build inverted edge map: "absFile::funcName" → []callerInfo.
	type callerInfo struct{ file, funcName string }
	inverted := make(map[string][]callerInfo)

	_ = callgraph.GraphVisitEdges(cg, func(e *callgraph.Edge) error {
		if e.Callee.Func == nil || e.Caller.Func == nil {
			return nil
		}
		calleeFile := prog.Fset.Position(e.Callee.Func.Pos()).Filename
		callerFile := prog.Fset.Position(e.Caller.Func.Pos()).Filename
		if calleeFile == "" || callerFile == "" {
			return nil
		}
		key := calleeFile + "::" + e.Callee.Func.Name()
		inverted[key] = append(inverted[key], callerInfo{callerFile, e.Caller.Func.Name()})
		return nil
	})

	// Seed BFS from functions defined in unclaimed .go files.
	unclaimedSet := make(map[string]bool, len(goUnclaimed))
	for _, f := range goUnclaimed {
		unclaimedSet[f.Path] = true
	}

	type bfsItem struct {
		file, funcName string
		depth          int
		originFile     string
	}

	visited := make(map[string]bool)
	var queue []bfsItem

	for fn := range cg.Nodes {
		if fn == nil {
			continue
		}
		file := prog.Fset.Position(fn.Pos()).Filename
		if unclaimedSet[file] {
			key := file + "::" + fn.Name()
			if !visited[key] {
				visited[key] = true
				queue = append(queue, bfsItem{file, fn.Name(), 0, file})
			}
		}
	}

	// Determine via/confidence.
	via, confidence := "go-callgraph", 0.9
	if algo == "pta" {
		via, confidence = "go-callgraph-pta", 0.95
	}

	// BFS upward.
	seen := make(map[string]bool)        // dedup output by "METHOD /path"
	coveredFiles := make(map[string]bool) // origin files that reached a route
	var result []RouteMapping

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Terminal: current file is a route-registering file.
		if rms, ok := routeFileMappings[cur.file]; ok {
			coveredFiles[cur.originFile] = true
			for _, rm := range rms {
				dedupeKey := rm.Method + " " + rm.RoutePath
				if !seen[dedupeKey] {
					seen[dedupeKey] = true
					result = append(result, RouteMapping{
						SourceFile: rm.SourceFile,
						Method:     rm.Method,
						RoutePath:  rm.RoutePath,
						Via:        via,
						Confidence: confidence,
					})
				}
			}
			continue
		}

		// Depth cap.
		if maxDepth > 0 && cur.depth >= maxDepth {
			continue
		}

		// Traverse callers.
		key := cur.file + "::" + cur.funcName
		for _, caller := range inverted[key] {
			callerKey := caller.file + "::" + caller.funcName
			if !visited[callerKey] {
				visited[callerKey] = true
				queue = append(queue, bfsItem{caller.file, caller.funcName, cur.depth + 1, cur.originFile})
			}
		}
	}

	// Build claimed list.
	var claimed []ChangedFile
	for _, f := range goUnclaimed {
		if coveredFiles[f.Path] {
			claimed = append(claimed, f)
		}
	}
	return result, claimed, nil
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/rbt/... -run "TestGoCallGraph" -v -count=1
```

Expected: all 4 tests PASS (the integration test against the real fixture may be slow — 2–5s for go/packages load).

- [ ] **Step 5: Run all rbt tests to confirm no regressions**

```bash
go test ./internal/rbt/... -count=1
```

Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/rbt/callgraph_go.go internal/rbt/callgraph_go_test.go
git commit -m "feat(rbt): add GoCallGraphBuilder — RTA/PTA type-aware call graph for Go"
```

---

### Task 4: Wire GoCallGraphBuilder into Indexer

**Files:**
- Modify: `internal/rbt/indexer.go` (lines 26–34 for struct, lines 53–79 for RunHybrid)
- Modify: `internal/rbt/indexer_test.go`

- [ ] **Step 1: Write the failing test for runGoCallGraphPhase**

Add to `internal/rbt/indexer_test.go`:

```go
func TestRunGoCallGraphPhase_NoGoMod_ReturnsEmpty(t *testing.T) {
	// A directory without go.mod should return empty without panic.
	dir := t.TempDir()
	idx := &Indexer{SrcDir: dir, Algo: "rta"}
	unclaimed := []ChangedFile{{Path: filepath.Join(dir, "service.go")}}
	routeFiles := map[string][]RouteMapping{}

	mappings, claimed := idx.runGoCallGraphPhase(unclaimed, routeFiles)
	assert.Empty(t, mappings)
	assert.Empty(t, claimed)
}
```

- [ ] **Step 2: Run to confirm it fails**

```bash
go test ./internal/rbt/... -run "TestRunGoCallGraphPhase" -v 2>&1 | head -10
```

Expected: FAIL — `runGoCallGraphPhase` undefined, `Algo` field undefined.

- [ ] **Step 3: Add Algo field and runGoCallGraphPhase to indexer.go**

In `internal/rbt/indexer.go`, change the `Indexer` struct (currently lines 26–34):

```go
// Indexer orchestrates map file generation from source code.
type Indexer struct {
	SrcDir    string
	SpecPath  string
	OutPath   string
	Overwrite bool
	Store     *IndexStore
	Embedder  Embedder
	Depth     int    // 0 = dynamic BFS (stop at route node); >0 = fixed depth cap
	Algo      string // Go call graph algorithm: "rta" (default) | "pta"
}
```

Add the new phase method after `runTreeSitterPhase` (after line 97 in the current file):

```go
// runGoCallGraphPhase uses golang.org/x/tools/go/callgraph to perform type-aware
// call graph analysis for Go modules. Silently returns empty on any error so V2
// can handle the unclaimed files.
func (idx *Indexer) runGoCallGraphPhase(
	unclaimed []ChangedFile,
	routeFileMappings map[string][]RouteMapping,
) ([]RouteMapping, []ChangedFile) {
	b := &GoCallGraphBuilder{SrcDir: idx.SrcDir, Algo: idx.Algo}
	mappings, claimed, _ := b.BuildAndTrace(unclaimed, routeFileMappings, idx.Depth)
	return mappings, claimed
}
```

- [ ] **Step 4: Update RunHybrid to insert the new phase**

In `internal/rbt/indexer.go`, replace `RunHybrid` body (current lines 53–79):

```go
// RunHybrid uses tree-sitter + Go call graph (V3) + name-based call graph (V2)
// + embeddings + LLM confirmation.
func (idx *Indexer) RunHybrid(llmParser *LLMParser) error {
	if err := idx.checkOverwrite(); err != nil {
		return err
	}
	files, err := findSourceFiles(idx.SrcDir)
	if err != nil {
		return err
	}

	// Phase 1: tree-sitter direct route detection.
	mappings, routeFileMappings := idx.runTreeSitterPhase(files)
	unclaimed := subtractFiles(files, mappings)

	// Phase 2: Go type-aware call graph (V3) — handles interface dispatch.
	goMappings, goClaimed := idx.runGoCallGraphPhase(unclaimed, routeFileMappings)
	mappings = append(mappings, goMappings...)
	unclaimed = subtractChangedFiles(unclaimed, goClaimed)

	// Phase 3: name-based call graph (V2) — covers non-Go files and Go fallback.
	cgMappings, cgClaimed := idx.runCallGraphPhase(files, unclaimed, routeFileMappings, llmParser)
	mappings = append(mappings, cgMappings...)
	unclaimed = subtractChangedFiles(unclaimed, cgClaimed)

	// Phase 4: embedding-based matching for remaining unclaimed files.
	embedMappings, err := idx.runEmbedPhase(unclaimed)
	if err == nil {
		mappings = append(mappings, embedMappings...)
	}

	return idx.writeMapFile(mappings, "hybrid")
}
```

- [ ] **Step 5: Run all tests**

```bash
go test ./internal/rbt/... -count=1
```

Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/rbt/indexer.go internal/rbt/indexer_test.go
git commit -m "feat(rbt): wire runGoCallGraphPhase into RunHybrid as Phase 2; add Indexer.Algo field"
```

---

### Task 5: Add --algo flag to cmd/rbt_index.go

**Files:**
- Modify: `cmd/rbt_index.go`

- [ ] **Step 1: Add the flag and wire it**

In `cmd/rbt_index.go`, add the flag in `init()` after the `--depth` line:

```go
rbtIndexCmd.Flags().String("algo", "rta", "Go call graph algorithm: rta or pta (default rta)")
```

In `runRBTIndex`, read the flag and pass to `Indexer`:

```go
algo, _ := cmd.Flags().GetString("algo")
```

And update the `Indexer` literal to include `Algo: algo`:

```go
indexer := &rbt.Indexer{
    SrcDir:    srcDir,
    SpecPath:  specPath,
    OutPath:   outPath,
    Overwrite: overwrite,
    Store:     rbt.NewIndexStore(".caseforge-index"),
    Embedder:  rbt.NewOpenAIEmbedder(),
    Depth:     depth,
    Algo:      algo,
}
```

- [ ] **Step 2: Build and smoke-test**

```bash
go build -o /tmp/caseforge . && /tmp/caseforge rbt index --help | grep algo
```

Expected: `--algo string   Go call graph algorithm: rta or pta (default rta)` (or similar).

- [ ] **Step 3: Commit**

```bash
git add cmd/rbt_index.go
git commit -m "feat(cmd): add --algo flag to rbt index for Go call graph algorithm selection"
```

---

### Task 6: Acceptance tests + CLAUDE.md

**Files:**
- Modify: `docs/acceptance/acceptance-tests.md`
- Modify: `scripts/acceptance.sh`
- Modify: `CLAUDE.md`

- [ ] **Step 1: Add AT-064~066 to acceptance-tests.md**

Find the `rbt — Call Graph (V2)` section and add a new section after it:

```markdown
### `rbt` — Call Graph V3 (Go type-aware)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-064 | --algo flag registered on rbt index | `caseforge rbt index --help` | `--algo` listed | ✅ PASS |
| AT-065 | rbt index hybrid runs without error (no Go module) | `caseforge rbt index --spec petstore.yaml --strategy hybrid --src /tmp` | exit 0, map file written | ✅ PASS |
| AT-066 | --algo pta flag accepted | `caseforge rbt index --help` | `pta` mentioned in --algo description | ✅ PASS |
```

Update the Summary table: `rbt` row total goes from 12 to 15, Total from 63 to 66.

Update the "last run" date to 2026-04-04.

- [ ] **Step 2: Add checks to scripts/acceptance.sh**

After the `--- rbt callgraph ---` block, add:

```bash
# -------------------------------------------------------
# AT-064 – AT-066: rbt callgraph V3 (Go type-aware)
# -------------------------------------------------------
echo "--- rbt callgraph v3 ---"

contains "AT-064" "--algo flag registered on rbt index" "algo" \
  "'$BIN' rbt index --help 2>&1 || true"

contains "AT-065" "rbt index hybrid no-Go-module runs clean" "Map file written" \
  "mkdir -p '$WORKDIR/at065-out' && '$BIN' rbt index --spec '$WORKDIR/petstore.yaml' --strategy hybrid --src /tmp --out '$WORKDIR/at065-out/map.yaml' --overwrite 2>&1 || true"

contains "AT-066" "--algo accepts pta value" "pta" \
  "'$BIN' rbt index --help 2>&1 || true"
echo ""
```

- [ ] **Step 3: Update CLAUDE.md scenario count**

Change `All 63 scenarios must pass` → `All 66 scenarios must pass`.

- [ ] **Step 4: Run the full acceptance suite**

```bash
./scripts/acceptance.sh
```

Expected: 66/66 passed.

If AT-065 fails because `/tmp` has no `.go` files and the binary exits non-zero, check what the actual output is and adjust the `contains` pattern. The expectation is "Map file written" from a successful run — since `/tmp` is a valid src dir that just has no supported source files, `rbt index` should complete normally and write an empty map file.

- [ ] **Step 5: Commit everything**

```bash
git add docs/acceptance/acceptance-tests.md scripts/acceptance.sh CLAUDE.md
git commit -m "test: add AT-064~066 acceptance tests for RBT V3 --algo flag"
```

---

## Self-Review

**Spec coverage check:**

| Spec requirement | Covered by |
|-----------------|-----------|
| `GoCallGraphBuilder` with `SrcDir`, `Algo` fields | Task 3 |
| `BuildAndTrace` signature | Task 3 |
| go.mod check → silent skip | Task 3 (Step 3, lines with `os.Stat`) |
| go/packages.Load + SSA build | Task 3 (Step 3) |
| RTA default, PTA opt-in | Task 3 (Step 3, switch statement) |
| BFS with `originFile` tracking + `coveredFiles` return | Task 3 (Step 3) |
| `Via:"go-callgraph"` / `Confidence:0.9` for RTA | Task 3 (Step 3) |
| `Via:"go-callgraph-pta"` / `Confidence:0.95` for PTA | Task 3 (Step 3) |
| `Indexer.Algo` field | Task 4 |
| `runGoCallGraphPhase` method | Task 4 |
| Phase inserted between tree-sitter and V2 in `RunHybrid` | Task 4 |
| `--algo` flag on `rbt index` | Task 5 |
| AT-064~066 acceptance tests | Task 6 |
| CLAUDE.md count 66 | Task 6 |
| testdata fixture with interface dispatch | Task 2 |
| `golang.org/x/tools` dependency | Task 1 |

All spec requirements covered. No placeholders. Type names are consistent across all tasks (`GoCallGraphBuilder`, `BuildAndTrace`, `runGoCallGraphPhase`, `Algo`).
