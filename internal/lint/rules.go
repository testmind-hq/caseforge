// internal/lint/rules.go
package lint

import "github.com/testmind-hq/caseforge/internal/spec"

type LintIssue struct {
	RuleID   string
	Severity string // "error"|"warning"
	Message  string
	Path     string // e.g. "GET /users"
}

type LintRule interface {
	ID() string
	Severity() string
	Check(ps *spec.ParsedSpec) []LintIssue
}

var registeredRules []LintRule

func register(r LintRule) {
	registeredRules = append(registeredRules, r)
}
