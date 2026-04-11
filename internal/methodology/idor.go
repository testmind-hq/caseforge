// internal/methodology/idor.go
package methodology

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

const (
	idorAltUUID = "00000000-0000-0000-0000-000000000001"
	idorNilUUID = "00000000-0000-0000-0000-000000000000"
	idorAltInt  = 99999
	idorZeroInt = 0
)

// IDORTechnique generates test cases for Insecure Direct Object Reference (IDOR)
// vulnerabilities. For each path or query parameter that looks like an ID, it
// substitutes alternative values and expects the server to return 403 or 404.
type IDORTechnique struct {
	gen *datagen.Generator
}

func NewIDORTechnique() *IDORTechnique {
	return &IDORTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *IDORTechnique) Name() string { return "idor" }

// Applies returns true if the operation has at least one path or query parameter
// that matches the ID heuristic.
func (t *IDORTechnique) Applies(op *spec.Operation) bool {
	for _, p := range op.Parameters {
		if p.In != "path" && p.In != "query" {
			continue
		}
		if isIDParam(p) {
			return true
		}
	}
	return false
}

// Generate produces 2 test cases per detected ID parameter: one with an
// alternative value and one with a zero/nil value.
func (t *IDORTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	var cases []schema.TestCase

	for _, p := range op.Parameters {
		if p.In != "path" && p.In != "query" {
			continue
		}
		if !isIDParam(p) {
			continue
		}

		isUUID := p.Schema != nil && (p.Schema.Format == "uuid" || p.Schema.Format == "guid")

		type mutation struct {
			label string
			value string
		}
		var mutations []mutation
		if isUUID {
			mutations = []mutation{
				{label: "alt_uuid", value: idorAltUUID},
				{label: "nil_uuid", value: idorNilUUID},
			}
		} else {
			mutations = []mutation{
				{label: "alt_id", value: fmt.Sprintf("%d", idorAltInt)},
				{label: "zero_id", value: fmt.Sprintf("%d", idorZeroInt)},
			}
		}

		for _, m := range mutations {
			specPath := fmt.Sprintf("%s %s parameters.%s", op.Method, op.Path, p.Name)
			title := fmt.Sprintf("IDOR %s=%s (%s)", p.Name, m.value, m.label)
			tc := buildTestCase(op, nil, title, specPath)
			tc.Priority = "P1"

			// Override the path to substitute the ID value
			tc.Steps[0].Path = buildIDORPath(op, p, m.value)

			// Accept any 4xx: 403 (forbidden) or 404 (resource-hiding) are both
			// valid secure responses; 2xx signals a likely IDOR vulnerability.
			tc.Steps[0].Assertions = []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 400},
				{Target: "status_code", Operator: "lt", Expected: 500},
			}
			tc.Source = schema.CaseSource{
				Technique: "idor",
				SpecPath:  specPath,
				Rationale: fmt.Sprintf("IDOR probe: substituting %s=%s to test authorization boundary", p.Name, m.value),
				Scenario:  "IDOR_PARAM",
			}
			cases = append(cases, tc)
		}
	}

	return cases, nil
}

// buildIDORPath constructs the request path with the given parameter substituted.
func buildIDORPath(op *spec.Operation, p *spec.Parameter, substituteValue string) string {
	if p.In == "path" {
		// Replace {paramName} placeholder in path
		path := op.Path
		// First replace all other path params with valid placeholders
		for _, other := range op.Parameters {
			if other.In == "path" && other.Name != p.Name {
				path = strings.ReplaceAll(path, "{"+other.Name+"}", "1")
			}
		}
		path = strings.ReplaceAll(path, "{"+p.Name+"}", substituteValue)
		return path
	}
	// For query parameters: build base path (replacing other path params), then add this param
	path := op.Path
	for _, other := range op.Parameters {
		if other.In == "path" {
			path = strings.ReplaceAll(path, "{"+other.Name+"}", "1")
		}
	}
	params := map[string]any{p.Name: substituteValue}
	return buildPathWithQuery(path, params)
}

// isIDParam returns true if the parameter looks like an ID based on its name or schema.
func isIDParam(p *spec.Parameter) bool {
	// Name exactly "id" (case-insensitive)
	if strings.EqualFold(p.Name, "id") {
		return true
	}
	// UUID or GUID format
	if p.Schema != nil && (p.Schema.Format == "uuid" || p.Schema.Format == "guid") {
		return true
	}
	// Integer type with ID-like name
	if p.Schema != nil && p.Schema.Type == "integer" && isIDLike(p.Name) {
		return true
	}
	return false
}

// isIDLike returns true if the parameter name contains "id" (case-insensitive).
// This catches: id, userId, user_id, orderId, order_id, etc.
func isIDLike(name string) bool {
	return strings.Contains(strings.ToLower(name), "id")
}
