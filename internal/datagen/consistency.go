// internal/datagen/consistency.go
package datagen

import (
	"sort"
	"strings"
	"time"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// ApplyCrossFieldConstraints enforces two kinds of consistency in a generated body:
//  1. Temporal ordering: datetime fields are assigned monotonically increasing values
//     ordered by semantic role (created < updated < deleted/expires).
//  2. Range ordering: numeric fields with min_/max_ prefix or _min/_max suffix are
//     swapped when min > max.
//
// The body is modified in-place and returned.
func ApplyCrossFieldConstraints(body map[string]any, s *spec.Schema) map[string]any {
	if len(body) == 0 || s == nil {
		return body
	}
	enforceTemporalOrder(body, s)
	enforceRangeOrder(body)
	return body
}

// --- temporal ordering ---

// temporalRoles maps name fragments to a rank: lower rank = earlier in time.
var temporalRoles = []struct {
	fragments []string
	rank      int
}{
	{[]string{"creat", "insert", "add", "start", "begin", "open", "from"}, 0},
	{[]string{"updat", "modif", "chang", "edit", "touch"}, 1},
	{[]string{"end", "finish", "complet", "clos", "to", "until"}, 2},
	{[]string{"delet", "remov", "expir", "purge", "archiv"}, 3},
}

// temporalRankOf returns the rank for a field, or -1 if not recognized as temporal.
func temporalRankOf(fieldName string, fieldSchema *spec.Schema) int {
	if fieldSchema == nil {
		return -1
	}
	// Must be a date/date-time schema OR have a clearly date-like suffix in its name.
	// Only underscore-prefixed suffixes are accepted to avoid false positives on
	// words like "runtime", "uptime", "candidate", "mandate".
	isDateType := fieldSchema.Format == "date-time" || fieldSchema.Format == "date"
	lower := strings.ToLower(fieldName)
	hasDateSuffix := strings.HasSuffix(lower, "_at") ||
		strings.HasSuffix(lower, "_date") ||
		strings.HasSuffix(lower, "_time")
	if !isDateType && !hasDateSuffix {
		return -1
	}
	for _, role := range temporalRoles {
		for _, frag := range role.fragments {
			if strings.Contains(lower, frag) {
				return role.rank
			}
		}
	}
	// Recognized as a temporal field but no known role — assign middle rank.
	return 1
}

func enforceTemporalOrder(body map[string]any, s *spec.Schema) {
	type tf struct {
		key    string
		rank   int
		isDate bool // true → "date" format, false → "date-time"
	}
	var fields []tf

	for fieldName, fieldSchema := range s.Properties {
		if _, exists := body[fieldName]; !exists {
			continue
		}
		rank := temporalRankOf(fieldName, fieldSchema)
		if rank < 0 {
			continue
		}
		fields = append(fields, tf{fieldName, rank, fieldSchema.Format == "date"})
	}
	if len(fields) < 2 {
		return // nothing to enforce with a single temporal field
	}

	// Sort by rank, then by key name for stability.
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].rank != fields[j].rank {
			return fields[i].rank < fields[j].rank
		}
		return fields[i].key < fields[j].key
	})

	// Assign monotonically increasing times spaced 24 h apart.
	// Only overwrite string values — skip if the generator produced a non-string
	// (e.g., schema has Format:"date-time" but Type:"integer", unusual but safe to skip).
	base := time.Now().UTC().Truncate(24 * time.Hour).Add(-time.Duration(len(fields)) * 24 * time.Hour)
	for i, f := range fields {
		if _, ok := body[f.key].(string); !ok {
			continue
		}
		t := base.Add(time.Duration(i+1) * 24 * time.Hour)
		if f.isDate {
			body[f.key] = t.Format("2006-01-02")
		} else {
			body[f.key] = t.Format("2006-01-02T15:04:05Z")
		}
	}
}

// --- range ordering ---

// enforceRangeOrder swaps min/max numeric field pairs when min > max.
// Recognises: min_X / max_X prefix pairs and X_min / X_max suffix pairs.
func enforceRangeOrder(body map[string]any) {
	// Build lowercase → original key map for fast lookup.
	lowerToKey := make(map[string]string, len(body))
	for k := range body {
		lowerToKey[strings.ToLower(k)] = k
	}

	for lowerName, origKey := range lowerToKey {
		var peerLower string

		switch {
		case strings.HasPrefix(lowerName, "min_"):
			peerLower = "max_" + lowerName[4:]
		case strings.HasSuffix(lowerName, "_min"):
			peerLower = lowerName[:len(lowerName)-4] + "_max"
		default:
			continue // not a min field; max fields are handled via their min counterpart
		}

		peerKey, ok := lowerToKey[peerLower]
		if !ok {
			continue
		}

		minVal := toFloat64(body[origKey])
		maxVal := toFloat64(body[peerKey])
		if minVal == nil || maxVal == nil {
			continue
		}
		if *minVal > *maxVal {
			body[origKey], body[peerKey] = body[peerKey], body[origKey]
		}
	}
}

// toFloat64 converts common numeric types to *float64. Returns nil for non-numerics.
func toFloat64(v any) *float64 {
	switch n := v.(type) {
	case float64:
		return &n
	case float32:
		f := float64(n)
		return &f
	case int:
		f := float64(n)
		return &f
	case int64:
		f := float64(n)
		return &f
	case int32:
		f := float64(n)
		return &f
	}
	return nil
}
