// internal/output/render/filename.go
package render

import (
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// FilenameFor derives a human-readable base name (no extension) for a test
// case output file. Format: {path-slug}_{method}_{description-slug}_{short-id}
//
// The path and method are extracted from the title (which carries the spec
// template path, e.g. /pets/{id}) so files group naturally by API resource.
//
// Examples:
//
//	"POST /users - valid request"              → users_post_valid_request_3bb73ff9
//	"GET /pets/{id} - missing id"              → pets_id_get_missing_id_a1b2c3d4
//	"[owasp] DELETE /admin/teams/{id} - sqli"  → admin_teams_id_delete_owasp_sqli_e5f6g7h8
//	"CRUD chain: /users"                        → users_crud_chain_b40d6929
func FilenameFor(tc schema.TestCase) string {
	slug := caseSlug(tc)
	id := shortID(tc.ID)
	if slug == "" {
		return id
	}
	return slug + "_" + id
}

// caseSlug builds {path-slug}_{method}_{description-slug} by parsing the
// title, which always carries the spec template path (e.g. /pets/{id}).
// Falls back to step data then plain titleSlug when the title has no path.
func caseSlug(tc schema.TestCase) string {
	titlePath := extractTitlePath(tc.Title)

	if titlePath != "" {
		pathPart := titleSlug(titlePath)
		titleMethod := extractTitleMethod(tc.Title, titlePath)

		var methodPart, descPart string
		if titleMethod != "" {
			methodPart = strings.ToLower(titleMethod)
			descPart = stripMethodPathFromTitle(tc.Title, titleMethod, titlePath)
		} else {
			// Chain titles like "CRUD chain: /users" — no HTTP method before path.
			descPart = stripPathFromTitle(tc.Title, titlePath)
		}

		parts := make([]string, 0, 3)
		if pathPart != "" {
			parts = append(parts, pathPart)
		}
		if methodPart != "" {
			parts = append(parts, methodPart)
		}
		if descPart != "" {
			parts = append(parts, descPart)
		}
		return capSlug(strings.Join(parts, "_"))
	}

	// No path in title — fall back to step data.
	if len(tc.Steps) > 0 {
		first := tc.Steps[0]
		pathPart := titleSlug(first.Path)
		methodPart := strings.ToLower(first.Method)
		descPart := titleSlug(tc.Title)
		parts := make([]string, 0, 3)
		if pathPart != "" {
			parts = append(parts, pathPart)
		}
		if methodPart != "" {
			parts = append(parts, methodPart)
		}
		if descPart != "" {
			parts = append(parts, descPart)
		}
		return capSlug(strings.Join(parts, "_"))
	}

	return titleSlug(tc.Title)
}

// extractTitlePath returns the first token starting with "/" in the title.
// Titles almost always follow "METHOD /path - description" or "[tag] METHOD /path - description".
func extractTitlePath(title string) string {
	for _, token := range strings.Fields(title) {
		if strings.HasPrefix(token, "/") {
			return token
		}
	}
	return ""
}

// extractTitleMethod returns the HTTP method word immediately before path in
// the title, or "" when none is found (e.g. "CRUD chain: /users").
func extractTitleMethod(title, path string) string {
	lower := strings.ToLower(title)
	lpath := strings.ToLower(path)
	idx := strings.Index(lower, lpath)
	if idx == 0 {
		return ""
	}
	before := strings.TrimRight(lower[:idx], " ")
	lastSpace := strings.LastIndex(before, " ")
	var word string
	if lastSpace < 0 {
		word = before
	} else {
		word = before[lastSpace+1:]
	}
	for _, m := range []string{"get", "post", "put", "patch", "delete", "head", "options"} {
		if word == m {
			return strings.ToUpper(m)
		}
	}
	return ""
}

// stripMethodPathFromTitle removes "METHOD /path" (case-insensitive) from the
// title and slugifies the remainder. Any prefix before method+path (e.g.
// "[owasp] ") is preserved in the result.
func stripMethodPathFromTitle(title, method, path string) string {
	needle := strings.ToLower(method) + " " + strings.ToLower(path)
	lower := strings.ToLower(title)
	idx := strings.Index(lower, needle)
	if idx < 0 {
		return titleSlug(title)
	}
	prefix := strings.TrimSpace(title[:idx])
	rest := strings.TrimLeft(title[idx+len(needle):], " -")
	combined := prefix
	if rest != "" {
		if combined != "" {
			combined += " "
		}
		combined += rest
	}
	return titleSlug(combined)
}

// stripPathFromTitle removes the path portion from a chain-case title and
// slugifies the remainder.
func stripPathFromTitle(title, path string) string {
	lower := strings.ToLower(title)
	lpath := strings.ToLower(path)
	idx := strings.Index(lower, lpath)
	if idx < 0 {
		return titleSlug(title)
	}
	prefix := strings.TrimSpace(title[:idx])
	suffix := strings.TrimSpace(title[idx+len(path):])
	combined := prefix
	if suffix != "" {
		if combined != "" {
			combined += " "
		}
		combined += suffix
	}
	return titleSlug(strings.TrimSpace(combined))
}

// capSlug trims trailing underscores and caps at 80 runes.
func capSlug(slug string) string {
	slug = strings.TrimRight(slug, "_")
	runes := []rune(slug)
	if len(runes) > 80 {
		slug = strings.TrimRight(string(runes[:80]), "_")
	}
	return slug
}

// titleSlug converts a human-readable string into a lowercase, filesystem-safe
// slug. Non-alphanumeric runs collapse to a single underscore; result capped at
// 80 runes.
func titleSlug(title string) string {
	var b strings.Builder
	prevUnderscore := true // suppress leading underscores
	for _, r := range strings.ToLower(title) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			prevUnderscore = false
		default:
			if !prevUnderscore {
				b.WriteByte('_')
				prevUnderscore = true
			}
		}
	}
	slug := strings.TrimRight(b.String(), "_")
	runes := []rune(slug)
	if len(runes) > 80 {
		slug = strings.TrimRight(string(runes[:80]), "_")
	}
	return slug
}

// shortID extracts the hash portion from a TC-{hash} identifier.
// For any other format the full ID is returned unchanged.
func shortID(id string) string {
	if strings.HasPrefix(id, "TC-") {
		return id[3:]
	}
	return id
}
