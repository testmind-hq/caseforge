// internal/diff/suggest.go
package diff

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// AffectedCase describes a test case that may be broken by a spec change.
type AffectedCase struct {
	ID     string
	Title  string
	Reason string
}

// Suggest returns the test cases likely affected by Breaking or PotentiallyBreaking changes.
// cases should be loaded via writer.NewJSONSchemaWriter().Read(indexPath) before calling.
func Suggest(result DiffResult, cases []schema.TestCase) []AffectedCase {
	// Build a set of breaking change paths for fast lookup
	breakingPaths := map[string]Change{}
	for _, c := range result.Changes {
		if c.Kind == Breaking || c.Kind == PotentiallyBreaking {
			breakingPaths[c.Path] = c
		}
	}
	if len(breakingPaths) == 0 {
		return nil
	}

	var affected []AffectedCase
	for _, tc := range cases {
		reason := ""
		// Match via Source.SpecPath ("METHOD /path" → extract path)
		specPath := extractPath(tc.Source.SpecPath)
		if change, ok := breakingPaths[specPath]; ok {
			reason = fmt.Sprintf("%s change on %s: %s", change.Kind, change.Path, change.Description)
		}
		// Match via Steps[].Path (normalize {{x}} → {x}, then match structurally)
		if reason == "" {
			for _, step := range tc.Steps {
				normalizedPath := normalizeTemplatePath(step.Path)
				for changePath, change := range breakingPaths {
					if pathsMatch(normalizedPath, changePath) {
						reason = fmt.Sprintf("step path %s matches %s change: %s", step.Path, change.Kind, change.Description)
						break
					}
				}
				if reason != "" {
					break
				}
			}
		}
		if reason != "" {
			affected = append(affected, AffectedCase{ID: tc.ID, Title: tc.Title, Reason: reason})
		}
	}
	return affected
}

// extractPath extracts the path from a SpecPath string like "GET /users/{id}" → "/users/{id}".
func extractPath(specPath string) string {
	parts := strings.SplitN(specPath, " ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return specPath
}

// pathsMatch returns true if two paths have the same structure, treating any
// {param} segment as a wildcard so {id} and {userId} both match.
func pathsMatch(a, b string) bool {
	if a == b {
		return true
	}
	aParts := strings.Split(strings.Trim(a, "/"), "/")
	bParts := strings.Split(strings.Trim(b, "/"), "/")
	if len(aParts) != len(bParts) {
		return false
	}
	for i := range aParts {
		aIsParam := strings.HasPrefix(aParts[i], "{") && strings.HasSuffix(aParts[i], "}")
		bIsParam := strings.HasPrefix(bParts[i], "{") && strings.HasSuffix(bParts[i], "}")
		if aIsParam && bIsParam {
			// Both are params — match regardless of name
			continue
		}
		if aParts[i] != bParts[i] {
			return false
		}
	}
	return true
}

// normalizeTemplatePath converts Hurl/Postman {{varName}} to OpenAPI {varName}.
func normalizeTemplatePath(path string) string {
	// Replace {{x}} with {x}
	result := strings.Builder{}
	i := 0
	for i < len(path) {
		if i+1 < len(path) && path[i] == '{' && path[i+1] == '{' {
			// find closing }}
			end := strings.Index(path[i+2:], "}}")
			if end >= 0 {
				varName := path[i+2 : i+2+end]
				result.WriteString("{")
				result.WriteString(varName)
				result.WriteString("}")
				i = i + 2 + end + 2
				continue
			}
		}
		result.WriteByte(path[i])
		i++
	}
	return result.String()
}
