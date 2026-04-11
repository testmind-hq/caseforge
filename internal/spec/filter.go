// internal/spec/filter.go
package spec

import (
	"fmt"
	"regexp"
)

// FilterSet holds include/exclude predicates for operation filtering.
// An empty FilterSet (all fields nil/empty) passes all operations unchanged.
// Filtering is applied as AND across dimensions: an operation must satisfy
// ALL active filters to be included.
//
// Path filters use Go regexp partial-match semantics (like grep). Anchor
// with ^ or $ for full-path matching:
//
//	--include-path "^/users" matches /users and /users/123 but not /admin/users
//	--include-path "/users"  matches /users, /users/123, and /admin/users
type FilterSet struct {
	IncludePaths []string // regex patterns; op passes if path matches ≥1 (if set)
	ExcludePaths []string // regex patterns; op excluded if path matches any
	IncludeTags  []string // exact tag names; op passes if tags intersect (if set)
	ExcludeTags  []string // exact tag names; op excluded if tags intersect any
}

// IsEmpty reports whether no filters are set.
func (f *FilterSet) IsEmpty() bool {
	return len(f.IncludePaths) == 0 && len(f.ExcludePaths) == 0 &&
		len(f.IncludeTags) == 0 && len(f.ExcludeTags) == 0
}

// Validate compiles all regex patterns and returns the first error encountered.
func (f *FilterSet) Validate() error {
	for _, p := range f.IncludePaths {
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("invalid --include-path regex %q: %w", p, err)
		}
	}
	for _, p := range f.ExcludePaths {
		if _, err := regexp.Compile(p); err != nil {
			return fmt.Errorf("invalid --exclude-path regex %q: %w", p, err)
		}
	}
	return nil
}

// Apply returns the subset of ops that pass all active filters.
// Returns ops unchanged when the FilterSet is empty.
func (f *FilterSet) Apply(ops []*Operation) []*Operation {
	if f.IsEmpty() {
		return ops
	}
	out := make([]*Operation, 0, len(ops))
	for _, op := range ops {
		if f.passes(op) {
			out = append(out, op)
		}
	}
	return out
}

func (f *FilterSet) passes(op *Operation) bool {
	// IncludePath: op.Path must match at least one pattern (if any set)
	if len(f.IncludePaths) > 0 {
		matched := false
		for _, p := range f.IncludePaths {
			if ok, _ := regexp.MatchString(p, op.Path); ok {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	// ExcludePath: op.Path must not match any pattern
	for _, p := range f.ExcludePaths {
		if ok, _ := regexp.MatchString(p, op.Path); ok {
			return false
		}
	}
	// IncludeTag: op.Tags must intersect with IncludeTags (if any set)
	if len(f.IncludeTags) > 0 {
		matched := false
	outer:
		for _, includeTag := range f.IncludeTags {
			for _, opTag := range op.Tags {
				if opTag == includeTag {
					matched = true
					break outer
				}
			}
		}
		if !matched {
			return false
		}
	}
	// ExcludeTag: op.Tags must not intersect with ExcludeTags
	for _, excludeTag := range f.ExcludeTags {
		for _, opTag := range op.Tags {
			if opTag == excludeTag {
				return false
			}
		}
	}
	return true
}
