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
	seen := make(map[string]bool)         // dedup output by "METHOD /path"
	coveredFiles := make(map[string]bool) // origin files that reached a route
	var result []RouteMapping

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Depth cap — enforced before terminal check so maxDepth=N means
		// "callers reachable in at most N−1 hops from the seeded functions".
		if maxDepth > 0 && cur.depth >= maxDepth {
			continue
		}

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
