// internal/methodology/variable_irrelevance.go
package methodology

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// VariableIrrelevanceTechnique detects parameter dependency groups and generates
// test cases demonstrating that dependent parameters are correctly ignored when
// their controlling parameter is disabled.
//
// Detection heuristic: a boolean (or 2-value enum) parameter P is a "controller"
// if at least one other parameter's name starts with P+"_" (e.g., "sort" controls
// "sort_field" and "sort_order"). This mirrors Tcases' variable irrelevance / NA
// marking: when sort=false, sort_field and sort_order are not applicable (NA).
type VariableIrrelevanceTechnique struct{}

func NewVariableIrrelevanceTechnique() *VariableIrrelevanceTechnique {
	return &VariableIrrelevanceTechnique{}
}

func (t *VariableIrrelevanceTechnique) Name() string { return "variable_irrelevance" }

func (t *VariableIrrelevanceTechnique) Applies(op *spec.Operation) bool {
	return len(detectParamGroups(op)) > 0
}

func (t *VariableIrrelevanceTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	groups := detectParamGroups(op)
	var cases []schema.TestCase

	for _, g := range groups {
		// Build path with controller=false, all dependents OMITTED
		queryParams := map[string]any{g.controller: false}
		// Include other independent params with placeholder values
		for _, p := range op.Parameters {
			if p.Name == g.controller {
				continue
			}
			isDependent := false
			for _, dep := range g.controlled {
				if p.Name == dep {
					isDependent = true
					break
				}
			}
			if !isDependent && p.Schema != nil {
				queryParams[p.Name] = defaultParamValue(p.Schema)
			}
		}

		path := buildPathWithQuery(resolvePathParams(op), queryParams)
		rationale := fmt.Sprintf(
			"variable irrelevance: %s=false makes %s not applicable (NA); omitting them is correct",
			g.controller, strings.Join(g.controlled, ", "),
		)

		tc := buildTestCase(op, nil,
			fmt.Sprintf("N/A when %s=false: %s omitted", g.controller, strings.Join(g.controlled, ", ")),
			fmt.Sprintf("%s %s parameters.%s", op.Method, op.Path, g.controller))
		tc.Priority = "P2"
		tc.Steps[0].Path = path
		tc.Steps[0].Body = nil
		tc.Source = schema.CaseSource{
			Technique: "variable_irrelevance",
			SpecPath:  fmt.Sprintf("%s %s parameters.%s", op.Method, op.Path, g.controller),
			Rationale: rationale,
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

// paramGroup describes a controlling parameter and its dependent parameters.
type paramGroup struct {
	controller string
	controlled []string
}

// detectParamGroups scans operation parameters for dependency patterns.
// A boolean or 2-value enum param P is a controller if ≥1 other param name
// starts with P+"_".
func detectParamGroups(op *spec.Operation) []paramGroup {
	// Build set of all param names
	names := make(map[string]bool, len(op.Parameters))
	for _, p := range op.Parameters {
		names[p.Name] = true
	}

	var groups []paramGroup
	for _, p := range op.Parameters {
		if p.Schema == nil {
			continue
		}
		// Only boolean or 2-value enum params can be controllers
		isController := p.Schema.Type == "boolean" ||
			(len(p.Schema.Enum) == 2)
		if !isController {
			continue
		}
		prefix := p.Name + "_"
		var controlled []string
		for name := range names {
			if name != p.Name && strings.HasPrefix(name, prefix) {
				controlled = append(controlled, name)
			}
		}
		if len(controlled) > 0 {
			groups = append(groups, paramGroup{
				controller: p.Name,
				controlled: controlled,
			})
		}
	}
	return groups
}

// defaultParamValue returns a simple default value for the parameter's schema.
func defaultParamValue(s *spec.Schema) any {
	if s == nil {
		return "value"
	}
	if len(s.Enum) > 0 {
		return s.Enum[0]
	}
	switch s.Type {
	case "boolean":
		return true
	case "integer", "number":
		return 1
	default:
		return "value"
	}
}

// resolvePathParams replaces {param} in path with "1".
func resolvePathParams(op *spec.Operation) string {
	path := op.Path
	for _, p := range op.Parameters {
		if p.In == "path" {
			path = strings.ReplaceAll(path, "{"+p.Name+"}", "1")
		}
	}
	return path
}
