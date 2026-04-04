// internal/rbt/callgraph_go.go
package rbt

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// GoCallGraphBuilder performs type-aware call graph analysis for Go modules
// using golang.org/x/tools/go/callgraph. It handles interface dispatch that
// V2 name-matching cannot resolve.
type GoCallGraphBuilder struct {
	SrcDir string
	Algo   string // "rta" (default) | "vta"
}

// BuildAndTrace loads the Go module in SrcDir, builds an SSA call graph via
// RTA or VTA, then BFS-traces upward from functions in unclaimed .go files to
// route-registering files. Returns the found route mappings, the subset of
// unclaimed files that were resolved, and any hard error (caller falls back to V2).
// Returns (nil, nil, nil) silently when: no go.mod, no main package, or no .go
// files in unclaimed.
//
// maxDepth controls how many BFS levels are explored from the seed functions:
//   - maxDepth=0 means unlimited (no depth cap).
//   - maxDepth=N means the BFS visits nodes at depths 0..N. A node at depth N
//     is recorded if it is a terminal (route file), but its callers are not
//     enqueued. Concretely, maxDepth=N means "traverse at most N hops from the
//     seeded functions", matching V2's TraceToRoutes semantics.
//
// Ordering: the terminal (route-file) check fires BEFORE the depth cap so that
// a route file reached exactly at depth=N is recorded when maxDepth=N.
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

	// Find main packages for RTA/VTA roots.
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
	// Collect main() roots for RTA (needed by both rta and vta paths).
	var roots []*ssa.Function
	for _, pkg := range mainPkgs {
		if fn := pkg.Func("main"); fn != nil {
			roots = append(roots, fn)
		}
	}
	if len(roots) == 0 {
		return nil, nil, nil
	}
	switch algo {
	case "vta":
		// VTA (Variable Type Analysis) is the modern replacement for the
		// deprecated go/pointer PTA. Run RTA first to get an initial call
		// graph, then refine it with VTA for more precise interface dispatch.
		rtaGraph := rta.Analyze(roots, true).CallGraph
		allFuncs := make(map[*ssa.Function]bool)
		for fn := range rtaGraph.Nodes {
			if fn != nil {
				allFuncs[fn] = true
			}
		}
		cg = vta.CallGraph(allFuncs, rtaGraph)
	default: // "rta"
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
	if algo == "vta" {
		via, confidence = "go-callgraph-vta", 0.92
	}

	// BFS upward.
	seen := make(map[string]bool)         // dedup output by "METHOD /path"
	coveredFiles := make(map[string]bool) // origin files that reached a route
	var result []RouteMapping

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Terminal: current file is a route-registering file.
		// Checked BEFORE the depth cap so that a route file reached at exactly
		// depth=maxDepth is still recorded (matches V2's TraceToRoutes ordering).
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

		// Depth cap — stop enqueuing callers once we reach the limit.
		// Applied AFTER the terminal check so route files at exactly depth=maxDepth
		// are still recorded above before we stop traversal.
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
