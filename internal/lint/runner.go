// internal/lint/runner.go
package lint

import "github.com/testmind-hq/caseforge/internal/spec"

// RunAll executes all registered rules against the spec, skipping any rule
// whose ID is present in skip.
func RunAll(ps *spec.ParsedSpec, skip map[string]bool) []LintIssue {
	var issues []LintIssue
	for _, rule := range registeredRules {
		if skip[rule.ID()] {
			continue
		}
		issues = append(issues, rule.Check(ps)...)
	}
	return issues
}
