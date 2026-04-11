// internal/methodology/depgraph.go
package methodology

import (
	"fmt"
	"sort"
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// DepEdge represents an inferred producer-consumer relationship between two operations.
type DepEdge struct {
	Creator     *spec.Operation // POST /collection (creates the resource)
	Consumer    *spec.Operation // GET|PUT|PATCH|DELETE /collection/{param}
	PathParam   string          // OpenAPI path param name, e.g. "userId"
	IDField     string          // field name in Creator response body, e.g. "id"
	CaptureFrom string          // resolved capture expr: "jsonpath $.id" or "header Location"
}

// DepGraph is the full set of inferred producer-consumer edges for a spec.
type DepGraph struct {
	Edges []DepEdge
}

// BuildDepGraph scans ops and infers producer-consumer edges:
// for each item-path operation (GET/PUT/PATCH/DELETE /x/{param}),
// it looks for a POST on the collection path /x and records the edge.
func BuildDepGraph(ops []*spec.Operation) *DepGraph {
	creators := make(map[string]*spec.Operation)
	for _, op := range ops {
		if op.Method == "POST" && !isItemPath(op.Path) {
			creators[op.Path] = op
		}
	}

	seen := make(map[string]bool)
	var edges []DepEdge
	for _, op := range ops {
		if !isItemPath(op.Path) {
			continue
		}
		switch op.Method {
		case "GET", "PUT", "PATCH", "DELETE":
		default:
			continue
		}
		colPath := collectionOf(op.Path)
		creator, ok := creators[colPath]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%s|%s|%s", creator.Path, op.Method, op.Path)
		if seen[key] {
			continue
		}
		seen[key] = true

		paramName := captureVarName(op.Path)
		// Determine idField: prefer nested path (e.g. "data.id") when present
		idField := findIDField(creator)
		if nested := findNestedIDPath(creator); nested != "" {
			idField = nested
		}
		captureFrom := inferCaptureFrom(creator, findIDField(creator))
		edges = append(edges, DepEdge{
			Creator:     creator,
			Consumer:    op,
			PathParam:   paramName,
			IDField:     idField,
			CaptureFrom: captureFrom,
		})
	}

	// Pass 2: incorporate explicitly declared OpenAPI Links.
	// Links provide authoritative producer→consumer relationships with exact parameter bindings,
	// supplementing the path-heuristic approach above.
	opsByID := make(map[string]*spec.Operation)
	for _, op := range ops {
		if op.OperationID != "" {
			opsByID[op.OperationID] = op
		}
	}
	for _, op := range ops {
		for _, link := range op.Links {
			consumer, ok := opsByID[link.OperationID]
			if !ok {
				continue
			}
			for paramName, paramExpr := range link.Parameters {
				captureFrom := parseLinkExpression(paramExpr)
				if captureFrom == "" {
					continue
				}
				// idField: strip the "jsonpath $." prefix produced by parseLinkExpression
				// so both values share a single source of truth.
				idField := strings.TrimPrefix(captureFrom, "jsonpath $.")

				key := fmt.Sprintf("%s|%s|%s", op.Path, consumer.Method, consumer.Path)
				if seen[key] {
					continue
				}
				seen[key] = true
				edges = append(edges, DepEdge{
					Creator:     op,
					Consumer:    consumer,
					PathParam:   paramName,
					IDField:     idField,
					CaptureFrom: captureFrom,
				})
			}
		}
	}

	// Sort for determinism
	sort.Slice(edges, func(i, j int) bool {
		ki := edges[i].Creator.Path + edges[i].Consumer.Method + edges[i].Consumer.Path
		kj := edges[j].Creator.Path + edges[j].Consumer.Method + edges[j].Consumer.Path
		return ki < kj
	})
	return &DepGraph{Edges: edges}
}

// inferCaptureFrom decides the capture expression for a creator operation.
// Prefers Location header (REST 201 convention) when documented; falls back to jsonpath.
func inferCaptureFrom(creator *spec.Operation, idField string) string {
	if resp, ok := creator.Responses["201"]; ok {
		if _, hasLoc := resp.Headers["Location"]; hasLoc {
			return "header Location"
		}
	}
	if nested := findNestedIDPath(creator); nested != "" {
		return fmt.Sprintf("jsonpath $.%s", nested)
	}
	return fmt.Sprintf("jsonpath $.%s", idField)
}

// findNestedIDPath returns a dotted path like "data.id" when the response body
// wraps the resource in a common envelope field (data, result, payload, response).
func findNestedIDPath(op *spec.Operation) string {
	wrappers := []string{"data", "result", "payload", "response"}

	// Sort codes for deterministic iteration
	codes := make([]string, 0, len(op.Responses))
	for code := range op.Responses {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	for _, code := range codes {
		resp := op.Responses[code]
		n := 0
		fmt.Sscanf(code, "%d", &n)
		if n < 200 || n >= 300 {
			continue
		}
		mt, ok := resp.Content["application/json"]
		if !ok || mt == nil || mt.Schema == nil {
			continue
		}
		for _, w := range wrappers {
			ws, ok := mt.Schema.Properties[w]
			if !ok || ws == nil || ws.Type != "object" {
				continue
			}
			if _, hasID := ws.Properties["id"]; hasID {
				return w + ".id"
			}
			var candidates []string
			for name := range ws.Properties {
				if strings.HasSuffix(strings.ToLower(name), "id") {
					candidates = append(candidates, name)
				}
			}
			if len(candidates) > 0 {
				sort.Strings(candidates)
				return w + "." + candidates[0]
			}
		}
	}
	return ""
}

// parseLinkExpression converts an OpenAPI link runtime expression to a CaseForge capture expression.
// "$response.body#/id"       → "jsonpath $.id"
// "$response.body#/data/id"  → "jsonpath $.data.id"
// Anything that is not a $response.body expression returns "" (caller should skip).
func parseLinkExpression(expr string) string {
	const prefix = "$response.body#/"
	if !strings.HasPrefix(expr, prefix) {
		return ""
	}
	path := strings.TrimPrefix(expr, prefix)
	dotPath := strings.ReplaceAll(path, "/", ".")
	return fmt.Sprintf("jsonpath $.%s", dotPath)
}
