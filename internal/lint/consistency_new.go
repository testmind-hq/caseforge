// internal/lint/consistency_new.go
package lint

import (
	"fmt"
	"sort"
	"strings"

	"github.com/testmind-hq/caseforge/internal/security"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
	register(&ruleL016{})
	register(&ruleL017{})
	register(&ruleL018{})
}

// L016: duplicate operationId (error, P0)
type ruleL016 struct{}

func (r *ruleL016) ID() string       { return "L016" }
func (r *ruleL016) Severity() string { return "error" }
func (r *ruleL016) Check(ps *spec.ParsedSpec) []LintIssue {
	seen := map[string]string{} // operationId → "METHOD /path"
	var issues []LintIssue
	for _, op := range ps.Operations {
		if op.OperationID == "" {
			continue
		}
		first, exists := seen[op.OperationID]
		if exists {
			issues = append(issues, LintIssue{
				RuleID:   "L016",
				Severity: "error",
				Message:  fmt.Sprintf("duplicate operationId %q (also used by %s)", op.OperationID, first),
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		} else {
			seen[op.OperationID] = fmt.Sprintf("%s %s", op.Method, op.Path)
		}
	}
	return issues
}

// L017: inconsistent path versioning (warning, P1)
type ruleL017 struct{}

func (r *ruleL017) ID() string       { return "L017" }
func (r *ruleL017) Severity() string { return "warning" }
func (r *ruleL017) Check(ps *spec.ParsedSpec) []LintIssue {
	v1Paths, v2Paths := security.FindVersionedPaths(ps.Operations)
	if len(v1Paths) == 0 {
		return nil
	}
	return []LintIssue{{
		RuleID:   "L017",
		Severity: "warning",
		Message:  fmt.Sprintf("mixed path versions: v1 in [%s] and v2 in [%s]", strings.Join(v1Paths, ", "), strings.Join(v2Paths, ", ")),
		Path:     "spec",
	}}
}

// L018: inconsistent response Content-Type (warning, P1)
type ruleL018 struct{}

func (r *ruleL018) ID() string       { return "L018" }
func (r *ruleL018) Severity() string { return "warning" }
func (r *ruleL018) Check(ps *spec.ParsedSpec) []LintIssue {
	contentTypes := map[string]bool{}
	for _, op := range ps.Operations {
		for _, resp := range op.Responses {
			for ct := range resp.Content {
				contentTypes[ct] = true
			}
		}
	}
	if len(contentTypes) <= 1 {
		return nil
	}
	types := make([]string, 0, len(contentTypes))
	for ct := range contentTypes {
		types = append(types, ct)
	}
	sort.Strings(types)
	return []LintIssue{{
		RuleID:   "L018",
		Severity: "warning",
		Message:  fmt.Sprintf("inconsistent response Content-Type: %s", strings.Join(types, ", ")),
		Path:     "spec",
	}}
}
