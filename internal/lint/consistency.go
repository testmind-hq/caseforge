// internal/lint/consistency.go
package lint

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

func init() {
	register(&ruleL007{})
	register(&ruleL008{})
	register(&ruleL009{})
	register(&ruleL010{})
}

var verbSegments = []string{"get", "create", "update", "delete", "list", "fetch", "add", "remove"}

// L007: verb in path segment
type ruleL007 struct{}

func (r *ruleL007) ID() string       { return "L007" }
func (r *ruleL007) Severity() string { return "warning" }
func (r *ruleL007) Check(ps *spec.ParsedSpec) []LintIssue {
	var issues []LintIssue
	for _, op := range ps.Operations {
		found := ""
		for _, seg := range strings.Split(strings.Trim(op.Path, "/"), "/") {
			if strings.HasPrefix(seg, "{") {
				continue // skip path parameter placeholders like {listId}
			}
			lower := strings.ToLower(seg)
			for _, verb := range verbSegments {
				if strings.HasPrefix(lower, verb) {
					found = verb
					break
				}
			}
			if found != "" {
				break
			}
		}
		if found != "" {
			issues = append(issues, LintIssue{
				RuleID:   "L007",
				Severity: "warning",
				Message:  fmt.Sprintf("verb %q found in path segment", found),
				Path:     fmt.Sprintf("%s %s", op.Method, op.Path),
			})
		}
	}
	return issues
}

// L008: inconsistent naming style (camelCase vs snake_case)
type ruleL008 struct{}

func (r *ruleL008) ID() string       { return "L008" }
func (r *ruleL008) Severity() string { return "warning" }
func (r *ruleL008) Check(ps *spec.ParsedSpec) []LintIssue {
	var allNames []string
	for _, op := range ps.Operations {
		for _, p := range op.Parameters {
			allNames = append(allNames, p.Name)
		}
		if op.RequestBody != nil {
			if mt, ok := op.RequestBody.Content["application/json"]; ok && mt.Schema != nil {
				for name := range mt.Schema.Properties {
					allNames = append(allNames, name)
				}
			}
		}
	}
	hasCamel, hasSnake := false, false
	for _, name := range allNames {
		if strings.Contains(name, "_") {
			hasSnake = true
		} else if name != strings.ToLower(name) {
			hasCamel = true
		}
	}
	if hasCamel && hasSnake {
		return []LintIssue{{
			RuleID:   "L008",
			Severity: "warning",
			Message:  "mixed naming styles: camelCase and snake_case both present",
			Path:     "spec",
		}}
	}
	return nil
}

// L009: inconsistent pagination style
type ruleL009 struct{}

func (r *ruleL009) ID() string       { return "L009" }
func (r *ruleL009) Severity() string { return "warning" }
func (r *ruleL009) Check(ps *spec.ParsedSpec) []LintIssue {
	pageStyle, offsetStyle := false, false
	pageStyleOps, offsetStyleOps := []string{}, []string{}
	for _, op := range ps.Operations {
		for _, p := range op.Parameters {
			if p.In != "query" {
				continue
			}
			n := strings.ToLower(p.Name)
			if n == "page" || n == "size" {
				pageStyle = true
				pageStyleOps = append(pageStyleOps, fmt.Sprintf("%s %s", op.Method, op.Path))
			}
			if n == "offset" || n == "limit" {
				offsetStyle = true
				offsetStyleOps = append(offsetStyleOps, fmt.Sprintf("%s %s", op.Method, op.Path))
			}
		}
	}
	if pageStyle && offsetStyle {
		return []LintIssue{{
			RuleID:   "L009",
			Severity: "warning",
			Message:  fmt.Sprintf("mixed pagination styles: page/size in [%s] and offset/limit in [%s]", strings.Join(pageStyleOps, ", "), strings.Join(offsetStyleOps, ", ")),
			Path:     "spec",
		}}
	}
	return nil
}

// L010: inconsistent error response schema
type ruleL010 struct{}

func (r *ruleL010) ID() string       { return "L010" }
func (r *ruleL010) Severity() string { return "warning" }
func (r *ruleL010) Check(ps *spec.ParsedSpec) []LintIssue {
	seen := map[string]bool{} // sorted field key → true
	for _, op := range ps.Operations {
		for code, resp := range op.Responses {
			n, err := strconv.Atoi(code)
			if err != nil || n < 400 {
				continue
			}
			if mt, ok := resp.Content["application/json"]; ok && mt.Schema != nil {
				fields := make([]string, 0, len(mt.Schema.Properties))
				for name := range mt.Schema.Properties {
					fields = append(fields, name)
				}
				sort.Strings(fields)
				seen[strings.Join(fields, ",")] = true
			}
		}
	}
	if len(seen) >= 2 {
		var structures []string
		for k := range seen {
			structures = append(structures, "{"+k+"}")
		}
		sort.Strings(structures)
		return []LintIssue{{
			RuleID:   "L010",
			Severity: "warning",
			Message:  fmt.Sprintf("inconsistent error response schemas: %s", strings.Join(structures, " vs ")),
			Path:     "spec",
		}}
	}
	return nil
}
