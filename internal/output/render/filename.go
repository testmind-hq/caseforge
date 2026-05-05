// internal/output/render/filename.go
package render

import (
	"strings"
	"unicode"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// FilenameFor derives a human-readable base name (no extension) for a test
// case output file. Format: {title-slug}_{short-id}
//
// Examples:
//
//	"POST /users - valid request"   → post_users_valid_request_3bb73ff9
//	"GET /pets/{id} - missing id"   → get_pets_id_missing_id_a1b2c3d4
//	"[owasp] DELETE /admin - sqli"  → owasp_delete_admin_sqli_e5f6g7h8
func FilenameFor(tc schema.TestCase) string {
	slug := titleSlug(tc.Title)
	id := shortID(tc.ID)
	if slug == "" {
		return id
	}
	return slug + "_" + id
}

// titleSlug converts a human-readable title into a lowercase, filesystem-safe
// slug. Punctuation and path characters become underscores; consecutive
// underscores are collapsed; result is capped at 80 chars.
func titleSlug(title string) string {
	var b strings.Builder
	prevUnderscore := true // suppress leading underscores
	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			prevUnderscore = false
		default:
			// Collapse any non-alphanumeric run into a single underscore.
			if !prevUnderscore {
				b.WriteByte('_')
				prevUnderscore = true
			}
		}
	}
	slug := strings.TrimRight(b.String(), "_")
	if len(slug) > 80 {
		slug = strings.TrimRight(slug[:80], "_")
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
