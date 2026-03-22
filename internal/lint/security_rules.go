// internal/lint/security_rules.go
package lint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/testmind-hq/caseforge/internal/security"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
	register(&ruleL011{})
	register(&ruleL012{})
}

var excludedPathSubstrings = []string{"/public", "/health", "/login", "/register"}

// L011: non-GET operation missing security scheme declaration
type ruleL011 struct{}

func (r *ruleL011) ID() string       { return "L011" }
func (r *ruleL011) Severity() string { return "error" }
func (r *ruleL011) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.Method == "GET" {
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
		if len(op.Security) == 0 {
			issues = append(issues, LintIssue{
				RuleID:   "L011",
				Severity: "error",
				Message:  "non-GET operation has no security scheme declared",
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L012: sensitive field exposed in 2xx response schema
type ruleL012 struct{}

func (r *ruleL012) ID() string       { return "L012" }
func (r *ruleL012) Severity() string { return "error" }
func (r *ruleL012) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		for code, resp := range op.Responses {
			n, err := strconv.Atoi(code)
			if err != nil || n < 200 || n >= 300 {
				continue
			}
			if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
				for _, fieldName := range security.FindSensitiveFields(mt.Schema) {
					issues = append(issues, LintIssue{
						RuleID:   "L012",
						Severity: "error",
						Message:  fmt.Sprintf("sensitive field %q exposed in %s response", fieldName, code),
						Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
					})
				}
			}
		}
	}
	return issues
}
