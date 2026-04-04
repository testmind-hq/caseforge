// internal/rbt/callgraph.go
package rbt

// CallGraphBuilder extracts function definitions and call edges from a single file.
type CallGraphBuilder interface {
	ExtractFuncs(filePath string) (defs []CallNode, calls []CallEdge, err error)
}

// BuildCallGraph builds an inverted call graph over all source files.
// Returns the graph and a per-file definition index so callers avoid a second
// ExtractFuncs pass when collecting BFS start nodes.
func BuildCallGraph(files []ChangedFile, builder CallGraphBuilder) (*CallGraph, map[string][]CallNode) {
	cg := &CallGraph{Edges: make(map[string][]CallNode)}
	defsByFile := make(map[string][]CallNode)

	type fileData struct {
		defs  []CallNode
		calls []CallEdge
	}
	// Single pass: collect defs and calls from each file once.
	fileResults := make([]fileData, 0, len(files))
	for _, f := range files {
		defs, calls, err := builder.ExtractFuncs(f.Path)
		if err != nil {
			fileResults = append(fileResults, fileData{})
			continue
		}
		if len(defs) > 0 {
			defsByFile[f.Path] = defs
		}
		fileResults = append(fileResults, fileData{defs, calls})
	}

	// Build definition index: funcName → all nodes with that name.
	allDefs := make(map[string][]CallNode)
	for _, fd := range fileResults {
		for _, d := range fd.defs {
			allDefs[d.FuncName] = append(allDefs[d.FuncName], d)
		}
	}

	// Build inverted edges: callee key → []callers.
	for _, fd := range fileResults {
		for _, edge := range fd.calls {
			callerNode := CallNode{File: edge.CallerFile, FuncName: edge.CallerFunc}
			for _, calleeNode := range allDefs[edge.CalleeName] {
				key := CallNodeKey(calleeNode.File, calleeNode.FuncName)
				cg.Edges[key] = append(cg.Edges[key], callerNode)
			}
		}
	}
	return cg, defsByFile
}

// TraceToRoutes performs BFS from startNodes upward through the call graph.
// It stops when it reaches a file in routeFileMappings (a route-registering file)
// or when the traversal depth exceeds maxDepth (0 = unlimited / dynamic).
// via and confidence are set on every returned RouteMapping.
// The second return value is the set of start-node file paths that reached a route
// (used by the caller to subtract those files from the "unclaimed" pool).
func TraceToRoutes(
	cg *CallGraph,
	startNodes []CallNode,
	routeFileMappings map[string][]RouteMapping,
	maxDepth int,
	via string,
	confidence float64,
) ([]RouteMapping, map[string]bool) {
	type item struct {
		node       CallNode
		depth      int
		originFile string // file of the start node that initiated this BFS path
	}

	visited := make(map[string]bool)
	queue := make([]item, 0, len(startNodes))
	for _, n := range startNodes {
		key := CallNodeKey(n.File, n.FuncName)
		if !visited[key] {
			visited[key] = true
			queue = append(queue, item{n, 0, n.File})
		}
	}

	seen := make(map[string]bool) // dedup output by "METHOD /path"
	coveredFiles := make(map[string]bool)
	var result []RouteMapping

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Terminal: this file registers routes.
		if rms, ok := routeFileMappings[cur.node.File]; ok {
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
			continue // don't traverse further from a route-registering file
		}

		// Depth cap.
		if maxDepth > 0 && cur.depth >= maxDepth {
			continue
		}

		// Traverse callers.
		key := CallNodeKey(cur.node.File, cur.node.FuncName)
		for _, caller := range cg.Edges[key] {
			callerKey := CallNodeKey(caller.File, caller.FuncName)
			if !visited[callerKey] {
				visited[callerKey] = true
				queue = append(queue, item{caller, cur.depth + 1, cur.originFile})
			}
		}
	}
	return result, coveredFiles
}

// subtractFiles returns files from `all` whose path does not appear as SourceFile
// in any of the `claimed` RouteMapping entries (used for Phase 1 tree-sitter subtraction).
func subtractFiles(all []ChangedFile, claimed []RouteMapping) []ChangedFile {
	claimedSet := make(map[string]bool, len(claimed))
	for _, m := range claimed {
		claimedSet[m.SourceFile] = true
	}
	var remaining []ChangedFile
	for _, f := range all {
		if !claimedSet[f.Path] {
			remaining = append(remaining, f)
		}
	}
	return remaining
}

// subtractChangedFiles removes files in `remove` from `all` by path.
// Used for Phase 2 call-graph subtraction where the covered files are tracked
// separately from RouteMapping.SourceFile.
func subtractChangedFiles(all, remove []ChangedFile) []ChangedFile {
	removeSet := make(map[string]bool, len(remove))
	for _, f := range remove {
		removeSet[f.Path] = true
	}
	var remaining []ChangedFile
	for _, f := range all {
		if !removeSet[f.Path] {
			remaining = append(remaining, f)
		}
	}
	return remaining
}

// fallbackCallGraphBuilder tries primary first; if it returns nothing or errors,
// uses fallback. Records whether the fallback was used (for Via/Confidence labeling).
type fallbackCallGraphBuilder struct {
	primary    CallGraphBuilder
	fallback   CallGraphBuilder
	hasUsedLLM bool
}

func (b *fallbackCallGraphBuilder) ExtractFuncs(filePath string) ([]CallNode, []CallEdge, error) {
	defs, calls, err := b.primary.ExtractFuncs(filePath)
	if err == nil && (len(defs) > 0 || len(calls) > 0) {
		return defs, calls, nil
	}
	rdefs, rcalls, rerr := b.fallback.ExtractFuncs(filePath)
	if rerr == nil && (len(rdefs) > 0 || len(rcalls) > 0) {
		b.hasUsedLLM = true
	}
	return rdefs, rcalls, rerr
}
