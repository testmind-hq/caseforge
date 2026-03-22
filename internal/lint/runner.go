// internal/lint/runner.go
package lint

import "github.com/testmind-hq/caseforge/internal/spec"

// RunAll executes all registered rules against the spec.
func RunAll(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, rule := range registeredRules {
		issues = append(issues, rule.Check(ps)...)
	}
	return issues
}
