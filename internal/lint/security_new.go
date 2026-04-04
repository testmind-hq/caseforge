// internal/lint/security_new.go
package lint

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/security"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
	register(&ruleL019{})
	register(&ruleL020{})
	register(&ruleL021{})
}

// L019: GET operation missing security scheme (warning, P1)
type ruleL019 struct{}

func (r *ruleL019) ID() string       { return "L019" }
func (r *ruleL019) Severity() string { return "warning" }
func (r *ruleL019) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.Method != "GET" {
			continue
		}
		excluded := false
		for _, sub := range excludedPathSubstrings {
			if strings.Contains(op.Path, sub) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		if len(op.Security) == 0 && len(ps.GlobalSecurity) == 0 {
			issues = append(issues, LintIssue{
				RuleID:   "L019",
				Severity: "warning",
				Message:  "GET operation has no security scheme declared",
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L020: sensitive field exposed as query parameter (error, P0)
type ruleL020 struct{}

func (r *ruleL020) ID() string       { return "L020" }
func (r *ruleL020) Severity() string { return "error" }
func (r *ruleL020) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		for _, p := range op.Parameters {
			if p.In != "query" {
				continue
			}
			if security.IsSensitiveName(p.Name) {
				issues = append(issues, LintIssue{
					RuleID:   "L020",
					Severity: "error",
					Message:  fmt.Sprintf("sensitive field %q exposed as query parameter", p.Name),
					Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
				})
			}
		}
	}
	return issues
}

// L021: no security schemes defined in spec (warning, P1)
type ruleL021 struct{}

func (r *ruleL021) ID() string       { return "L021" }
func (r *ruleL021) Severity() string { return "warning" }
func (r *ruleL021) Check(ps *spec.ParsedSpec) []LintIssue {
	if len(ps.SecuritySchemes) > 0 {
		return nil
	}
	return []LintIssue{{
		RuleID:   "L021",
		Severity: "warning",
		Message:  "no security schemes defined in spec",
		Path:     "spec",
	}}
}
