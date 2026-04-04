// internal/lint/completeness.go
package lint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
	// P0 (error — must implement, affect generation correctness)
	register(&ruleL004{})
	register(&ruleL005{})
	register(&ruleL006{})
	// P1 (warning — documentation quality)
	register(&ruleL001{})
	register(&ruleL002{})
	register(&ruleL003{})
	register(&ruleL013{})
	register(&ruleL014{})
	register(&ruleL015{})
}

// L001: operation missing operationId (warning, P1)
type ruleL001 struct{}

func (r *ruleL001) ID() string       { return "L001" }
func (r *ruleL001) Severity() string { return "warning" }
func (r *ruleL001) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.OperationID == "" {
			issues = append(issues, LintIssue{
				RuleID:   "L001",
				Severity: "warning",
				Message:  "operation is missing operationId",
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L002: operation missing summary or description (warning, P1)
type ruleL002 struct{}

func (r *ruleL002) ID() string       { return "L002" }
func (r *ruleL002) Severity() string { return "warning" }
func (r *ruleL002) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.Summary == "" && op.Description == "" {
			issues = append(issues, LintIssue{
				RuleID:   "L002",
				Severity: "warning",
				Message:  "operation has no summary or description",
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L003: request body field missing description (warning, P1)
type ruleL003 struct{}

func (r *ruleL003) ID() string       { return "L003" }
func (r *ruleL003) Severity() string { return "warning" }
func (r *ruleL003) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.RequestBody == nil {
			continue
		}
		for ct, mt := range op.RequestBody.Content {
			if mt.Schema == nil {
				continue
			}
			for fieldName, prop := range mt.Schema.Properties {
				_ = ct
				if prop.Description == "" {
					issues = append(issues, LintIssue{
						RuleID:   "L003",
						Severity: "warning",
						Message:  fmt.Sprintf("request body field %q has no description", fieldName),
						Path:     fmt.Sprintf("%s %s requestBody.properties.%s", op.Method, op.Path, fieldName),
					})
				}
			}
		}
	}
	return issues
}

// L004: no 2xx response defined (error, P0)
type ruleL004 struct{}

func (r *ruleL004) ID() string       { return "L004" }
func (r *ruleL004) Severity() string { return "error" }
func (r *ruleL004) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		has2xx := false
		for code := range op.Responses {
			n, err := strconv.Atoi(code)
			if err == nil && n >= 200 && n < 300 {
				has2xx = true
				break
			}
		}
		if !has2xx {
			issues = append(issues, LintIssue{
				RuleID:   "L004",
				Severity: "error",
				Message:  "no 2xx response defined",
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L005: undefined $ref (error, P0)
// Note: kin-openapi resolves $refs during loading; any unresolved ref causes a parse error.
// L005 catches refs that parsed but whose target schema is nil (partially resolved).
type ruleL005 struct{}

func (r *ruleL005) ID() string       { return "L005" }
func (r *ruleL005) Severity() string { return "error" }
func (r *ruleL005) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.RequestBody != nil {
			for _, mt := range op.RequestBody.Content {
				if mt != nil && mt.Schema != nil && mt.Schema.Ref != "" {
					if _, ok := ps.Schemas[schemaRefName(mt.Schema.Ref)]; !ok {
						issues = append(issues, LintIssue{
							RuleID:   "L005",
							Severity: "error",
							Message:  fmt.Sprintf("undefined $ref: %s", mt.Schema.Ref),
							Path:     fmt.Sprintf("%s %s requestBody", op.Method, op.Path),
						})
					}
				}
			}
		}
	}
	return issues
}

// schemaRefName extracts the component name from a $ref path like "#/components/schemas/Foo".
func schemaRefName(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

// L006: path parameter not declared in parameters (error, P0)
type ruleL006 struct{}

func (r *ruleL006) ID() string       { return "L006" }
func (r *ruleL006) Severity() string { return "error" }
func (r *ruleL006) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		declared := make(map[string]bool)
		for _, p := range op.Parameters {
			if p.In == "path" {
				declared[p.Name] = true
			}
		}
		// Extract {param} placeholders from path
		for _, seg := range strings.Split(op.Path, "/") {
			if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
				name := seg[1 : len(seg)-1]
				if !declared[name] {
					issues = append(issues, LintIssue{
						RuleID:   "L006",
						Severity: "error",
						Message:  fmt.Sprintf("path parameter {%s} not declared in parameters", name),
						Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
					})
				}
			}
		}
	}
	return issues
}

// L013: parameter missing type (warning, P1)
type ruleL013 struct{}

func (r *ruleL013) ID() string       { return "L013" }
func (r *ruleL013) Severity() string { return "warning" }
func (r *ruleL013) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		for _, p := range op.Parameters {
			if p.Schema == nil || p.Schema.Type == "" {
				issues = append(issues, LintIssue{
					RuleID:   "L013",
					Severity: "warning",
					Message:  fmt.Sprintf("parameter %q has no type declaration", p.Name),
					Path:     fmt.Sprintf("%s %s parameter.%s", op.Method, op.Path, p.Name),
				})
			}
		}
	}
	return issues
}

// L014: no 4xx error response defined (warning, P1)
type ruleL014 struct{}

func (r *ruleL014) ID() string       { return "L014" }
func (r *ruleL014) Severity() string { return "warning" }
func (r *ruleL014) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		has4xx := false
		for code := range op.Responses {
			n, err := strconv.Atoi(code)
			if err == nil && n >= 400 && n < 500 {
				has4xx = true
				break
			}
		}
		if !has4xx {
			issues = append(issues, LintIssue{
				RuleID:   "L014",
				Severity: "warning",
				Message:  "no 4xx error response defined",
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L015: 2xx response schema properties missing example (warning, P1)
type ruleL015 struct{}

func (r *ruleL015) ID() string       { return "L015" }
func (r *ruleL015) Severity() string { return "warning" }
func (r *ruleL015) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		for code, resp := range op.Responses {
			n, err := strconv.Atoi(code)
			if err != nil || n < 200 || n >= 300 {
				continue
			}
			if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
				for fieldName, prop := range mt.Schema.Properties {
					if prop.Example == nil {
						issues = append(issues, LintIssue{
							RuleID:   "L015",
							Severity: "warning",
							Message:  fmt.Sprintf("response field %q has no example", fieldName),
							Path:     fmt.Sprintf("%s %s %s response.properties.%s", op.Method, op.Path, code, fieldName),
						})
					}
				}
			}
		}
	}
	return issues
}
