// internal/lint/rules.go
package lint

import "github.com/testmind-hq/caseforge/internal/spec"

type LintIssue struct {
	RuleID   string `json:"rule_id"`
	Severity string `json:"severity"` // "error"|"warning"
	Message  string `json:"message"`
	Path     string `json:"path"` // e.g. "GET /users"
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
