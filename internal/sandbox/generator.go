// internal/sandbox/generator.go
package sandbox

import (
	"fmt"
	"strings"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// ResponseGenerator tries strategies in order and returns the first successful result.
type ResponseGenerator struct {
	strategies []ResponseStrategy
}

// newResponseGenerator builds the strategy chain for the given format option.
// "auto" = ExampleStrategy → SchemaStrategy → FakerStrategy
// "schema" = SchemaStrategy only
// "faker"  = FakerStrategy only
func newResponseGenerator(format string) *ResponseGenerator {
	var strategies []ResponseStrategy
	switch format {
	case "schema":
		strategies = []ResponseStrategy{&schemaStrategy{}}
	case "faker":
		strategies = []ResponseStrategy{newFakerStrategy()}
	default: // "auto"
		strategies = []ResponseStrategy{
			&exampleStrategy{},
			&schemaStrategy{},
			newFakerStrategy(),
		}
	}
	return &ResponseGenerator{strategies: strategies}
}

// Generate produces a (statusCode, body) for the given operation.
// For write operations (POST/PUT/PATCH), the body is always a map[string]any
// with an "id" field guaranteed to be present.
func (g *ResponseGenerator) Generate(op *spec.Operation, pathParams map[string]string) (int, map[string]any) {
	statusCode := successStatusCode(op)
	var body any
	for _, s := range g.strategies {
		if result, ok := s.Generate(op, fmt.Sprintf("%d", statusCode)); ok {
			body = result
			break
		}
	}
	// Normalize to map[string]any
	m := toMap(body)
	// Guarantee an id field is present for write operations so StateStore can key on it
	if isWriteOp(op.Method) {
		if _, ok := m["id"]; !ok {
			m["id"] = gofakeit.UUID()
		}
	}
	return statusCode, m
}

func successStatusCode(op *spec.Operation) int {
	for _, code := range []string{"200", "201", "202"} {
		if _, ok := op.Responses[code]; ok {
			n := 0
			fmt.Sscanf(code, "%d", &n)
			return n
		}
	}
	return 200
}

func isWriteOp(method string) bool {
	m := strings.ToUpper(method)
	return m == "POST" || m == "PUT" || m == "PATCH"
}

func toMap(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

// extractID looks for common ID field names in a response body, trying path-param names too.
func extractID(body map[string]any, opPath string) string {
	candidates := []string{"id", "uuid", "ID"}
	for _, seg := range strings.Split(opPath, "/") {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			candidates = append(candidates, seg[1:len(seg)-1])
		}
	}
	for _, field := range candidates {
		if val, ok := body[field]; ok {
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}
