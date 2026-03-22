// internal/methodology/chain.go
package methodology

import (
	"fmt"
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
	var cases []schema.TestCase
	for resourcePath, group := range groups {
		if group.create == nil || group.read == nil {
			continue // Need at least POST + GET to form a chain
		}
		tc, err := t.buildChainCase(resourcePath, group)
		if err != nil {
			return nil, err
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// chainGroup holds the operations for a single resource path.
type chainGroup struct {
	create *spec.Operation // POST /collection
	read   *spec.Operation // GET /collection/{id}
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
	for code, resp := range op.Responses {
		n := 0
		fmt.Sscanf(code, "%d", &n)
		if n < 200 || n >= 300 {
			continue
		}
		if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
			for name := range mt.Schema.Properties {
				if name == "id" || strings.HasSuffix(strings.ToLower(name), "id") {
					return name
				}
			}
		}
	}
	return "id" // sensible default
}

func (t *ChainTechnique) buildChainCase(resourcePath string, g *chainGroup) (schema.TestCase, error) {
	id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
	captureName := captureVarName(g.read.Path)
	idField := findIDField(g.create)

	var steps []schema.Step

	// Step 1: setup — POST to create the resource
	// buildValidBody is defined in internal/methodology/helpers.go (same package)
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
		Captures: []schema.Capture{
			{Name: captureName, From: fmt.Sprintf("jsonpath $.%s", idField)},
		},
	})

	// Step 2: test — GET the created resource
	readPath := strings.ReplaceAll(g.read.Path,
		fmt.Sprintf("{%s}", captureName),
		fmt.Sprintf("{{%s}}", captureName))
	steps = append(steps, schema.Step{
		ID:         "step-test",
		Title:      fmt.Sprintf("read %s by id", resourcePath),
		Type:       "test",
		Method:     g.read.Method,
		Path:       readPath,
		DependsOn:  []string{"step-setup"},
		Assertions: assertpkg.BasicAssertions(g.read),
	})

	// Step 3: teardown — DELETE (optional)
	if g.delete != nil {
		deletePath := strings.ReplaceAll(g.delete.Path,
			fmt.Sprintf("{%s}", captureName),
			fmt.Sprintf("{{%s}}", captureName))
		steps = append(steps, schema.Step{
			ID:         "step-teardown",
			Title:      fmt.Sprintf("delete %s", resourcePath),
			Type:       "teardown",
			Method:     g.delete.Method,
			Path:       deletePath,
			DependsOn:  []string{"step-setup"},
			Assertions: assertpkg.BasicAssertions(g.delete),
		})
	}

	tc := schema.TestCase{
		Schema:  schema.SchemaBaseURL,
		Version: "1",
		ID:      id,
		Title:   fmt.Sprintf("CRUD chain: %s", resourcePath),
		Kind:    "chain",
		Priority: "P1",
		Tags:    g.create.Tags,
		Source: schema.CaseSource{
			Technique: "chain_crud",
			SpecPath:  resourcePath,
			Rationale: fmt.Sprintf("CRUD lifecycle: create → read → delete for %s", resourcePath),
		},
		Steps:       steps,
		GeneratedAt: time.Now(),
	}
	return tc, nil
}
