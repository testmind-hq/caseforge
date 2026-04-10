// internal/methodology/chain.go
package methodology

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	assertpkg "github.com/testmind-hq/caseforge/internal/assert"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// ChainTechnique generates multi-step CRUD chain cases by grouping operations
// that share a resource path (e.g., POST /users + GET /users/{id}).
type ChainTechnique struct {
	gen *datagen.Generator
}

func NewChainTechnique() *ChainTechnique {
	return &ChainTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *ChainTechnique) Name() string { return "chain_crud" }

func (t *ChainTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	groups := groupByResource(s.Operations)
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var cases []schema.TestCase
	for _, resourcePath := range keys {
		group := groups[resourcePath]
		if group.create == nil || group.read == nil {
			continue // Need at least POST + GET to form a chain
		}
		tc := t.buildChainCase(resourcePath, group)
		cases = append(cases, tc)
	}
	return cases, nil
}

// chainGroup holds the operations for a single resource path.
type chainGroup struct {
	create *spec.Operation // POST /collection
	read   *spec.Operation // GET /collection/{id}
	update *spec.Operation // PUT or PATCH /collection/{id} (optional)
	delete *spec.Operation // DELETE /collection/{id} (optional)
}

// groupByResource groups operations by their collection path.
// e.g., POST /users and GET /users/{userId} both belong to resource "/users".
func groupByResource(ops []*spec.Operation) map[string]*chainGroup {
	groups := map[string]*chainGroup{}

	ensure := func(key string) *chainGroup {
		if groups[key] == nil {
			groups[key] = &chainGroup{}
		}
		return groups[key]
	}

	for _, op := range ops {
		if isItemPath(op.Path) {
			collectionPath := collectionOf(op.Path)
			g := ensure(collectionPath)
			switch op.Method {
			case "GET":
				g.read = op
			case "PUT", "PATCH":
				if g.update == nil {
					g.update = op
				}
				if op.Method == "PUT" {
					g.update = op // PUT wins over PATCH
				}
			case "DELETE":
				g.delete = op
			}
		} else {
			if op.Method == "POST" {
				ensure(op.Path).create = op
			}
		}
	}
	return groups
}

// isItemPath returns true for paths like /users/{userId} (ends with a path param).
func isItemPath(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return false
	}
	last := parts[len(parts)-1]
	return strings.HasPrefix(last, "{") && strings.HasSuffix(last, "}")
}

// collectionOf strips the last path segment: /users/{userId} → /users
func collectionOf(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx <= 0 {
		return path
	}
	return path[:idx]
}

// captureVarName extracts the variable name from a path param, e.g. {userId} → userId.
func captureVarName(itemPath string) string {
	parts := strings.Split(strings.Trim(itemPath, "/"), "/")
	last := parts[len(parts)-1]
	return strings.Trim(last, "{}")
}

// findIDField returns the name of an `id`-like field in a 2xx response body schema.
func findIDField(op *spec.Operation) string {
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
		if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
			// Prefer exact "id" match first, then fall back to suffix match.
			// Iterating maps is non-deterministic in Go, so we collect candidates.
			if _, ok := mt.Schema.Properties["id"]; ok {
				return "id"
			}
			var candidates []string
			for name := range mt.Schema.Properties {
				if strings.HasSuffix(strings.ToLower(name), "id") {
					candidates = append(candidates, name)
				}
			}
			if len(candidates) > 0 {
				sort.Strings(candidates)
				return candidates[0]
			}
		}
	}
	return "id" // sensible default
}

func (t *ChainTechnique) buildChainCase(resourcePath string, g *chainGroup) schema.TestCase {
	id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
	captureName := captureVarName(g.read.Path)
	idField := findIDField(g.create)
	captureFrom := inferCaptureFrom(g.create, idField)

	var steps []schema.Step

	// Step 1: setup — POST to create the resource
	setupBody := buildValidBody(t.gen, g.create)
	var setupBodyAny any
	if setupBody != nil {
		setupBodyAny = setupBody
	}
	setupHeaders := map[string]string{}
	if setupBodyAny != nil {
		setupHeaders["Content-Type"] = "application/json"
	}
	steps = append(steps, schema.Step{
		ID:         "step-setup",
		Title:      fmt.Sprintf("create %s", resourcePath),
		Type:       "setup",
		Method:     g.create.Method,
		Path:       g.create.Path,
		Headers:    setupHeaders,
		Body:       setupBodyAny,
		Assertions: assertpkg.BasicAssertions(g.create),
		Captures:   []schema.Capture{{Name: captureName, From: captureFrom}},
	})

	// Step 2 (optional): update — PUT/PATCH the created resource
	if g.update != nil {
		updateParamName := captureVarName(g.update.Path)
		updatePath := strings.ReplaceAll(g.update.Path,
			fmt.Sprintf("{%s}", updateParamName),
			fmt.Sprintf("{{%s}}", captureName))
		updateBody := buildValidBody(t.gen, g.update)
		var updateBodyAny any
		if updateBody != nil {
			updateBodyAny = updateBody
		}
		updateHeaders := map[string]string{}
		if updateBodyAny != nil {
			updateHeaders["Content-Type"] = "application/json"
		}
		steps = append(steps, schema.Step{
			ID:         "step-update",
			Title:      fmt.Sprintf("update %s", resourcePath),
			Type:       "update",
			Method:     g.update.Method,
			Path:       updatePath,
			Headers:    updateHeaders,
			Body:       updateBodyAny,
			Assertions: assertpkg.BasicAssertions(g.update),
			DependsOn:  []string{"step-setup"},
		})
	}

	// Step 3: test — GET the created/updated resource
	readPath := strings.ReplaceAll(g.read.Path,
		fmt.Sprintf("{%s}", captureName),
		fmt.Sprintf("{{%s}}", captureName))
	var testDeps []string
	if g.update != nil {
		testDeps = []string{"step-update"}
	} else {
		testDeps = []string{"step-setup"}
	}
	steps = append(steps, schema.Step{
		ID:         "step-test",
		Title:      fmt.Sprintf("read %s by id", resourcePath),
		Type:       "test",
		Method:     g.read.Method,
		Path:       readPath,
		DependsOn:  testDeps,
		Assertions: assertpkg.BasicAssertions(g.read),
	})

	// Step 4 (optional): teardown — DELETE
	if g.delete != nil {
		deleteParamName := captureVarName(g.delete.Path)
		deletePath := strings.ReplaceAll(g.delete.Path,
			fmt.Sprintf("{%s}", deleteParamName),
			fmt.Sprintf("{{%s}}", captureName))
		steps = append(steps, schema.Step{
			ID:         "step-teardown",
			Title:      fmt.Sprintf("delete %s", resourcePath),
			Type:       "teardown",
			Method:     g.delete.Method,
			Path:       deletePath,
			DependsOn:  []string{"step-test"},
			Assertions: assertpkg.BasicAssertions(g.delete),
		})
	}

	stepDesc := "create → read"
	if g.update != nil {
		stepDesc = "create → update → read"
	}
	if g.delete != nil {
		stepDesc += " → delete"
	}

	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       id,
		Title:    fmt.Sprintf("CRUD chain: %s", resourcePath),
		Kind:     "chain",
		Priority: "P1",
		Tags:     g.create.Tags,
		Source: schema.CaseSource{
			Technique: "chain_crud",
			SpecPath:  resourcePath,
			Rationale: fmt.Sprintf("CRUD lifecycle: %s for %s", stepDesc, resourcePath),
		},
		Steps:       steps,
		GeneratedAt: time.Now(),
	}
}
