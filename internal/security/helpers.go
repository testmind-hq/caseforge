// internal/security/helpers.go
package security

import (
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// sensitiveKeywords are substrings that flag a field as sensitive (case-insensitive).
var sensitiveKeywords = []string{
	"password", "secret", "token", "private_key", "access_token", "refresh_token",
}

// HasIDPathParam returns true if the operation path contains a path parameter
// whose name ends with "id" (case-insensitive), e.g. {id}, {userId}, {orderId}.
func HasIDPathParam(op *spec.Operation) bool {
	for _, seg := range strings.Split(op.Path, "/") {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			name := strings.ToLower(seg[1 : len(seg)-1])
			if name == "id" || strings.HasSuffix(name, "id") {
				return true
			}
		}
	}
	return false
}

// FindSensitiveFields returns field names in schema.Properties that contain
// a sensitive keyword (case-insensitive substring match).
func FindSensitiveFields(s *spec.Schema) []string {
	if s == nil {
		return nil
	}
	var found []string
	for name := range s.Properties {
		lower := strings.ToLower(name)
		for _, kw := range sensitiveKeywords {
			if strings.Contains(lower, kw) {
				found = append(found, name)
				break
			}
		}
	}
	return found
}

// IsAuthRequired returns true if the operation declares at least one security scheme.
func IsAuthRequired(op *spec.Operation) bool {
	return len(op.Security) > 0
}

// FindVersionedPaths returns two slices: paths containing "/v1/" (or starting with "/v1/")
// and paths containing "/v2/" (or similar next-version prefix).
// Returns nil slices when no versioned pair exists.
func FindVersionedPaths(ops []*spec.Operation) (v1Paths, v2Paths []string) {
	seen := map[string][]string{} // version prefix → paths
	for _, op := range ops {
		p := op.Path
		for _, ver := range []string{"/v1/", "/v2/", "/v3/"} {
			if strings.Contains(p, ver) || strings.HasPrefix(p, ver[1:]) {
				seen[ver] = append(seen[ver], p)
				break
			}
		}
	}
	if len(seen["/v1/"]) > 0 && len(seen["/v2/"]) > 0 {
		return seen["/v1/"], seen["/v2/"]
	}
	return nil, nil
}
