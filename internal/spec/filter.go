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
	// IncludeTags lists exact tag names to include. An operation passes if its
	// Tags field intersects with this list. Operations that carry no tags never
	// intersect any non-empty IncludeTags list and are therefore excluded when
	// IncludeTags is set.
	IncludeTags []string
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
// Patterns are compiled once per Apply call to avoid repeated compilation.
func (f *FilterSet) Apply(ops []*Operation) []*Operation {
	if f.IsEmpty() {
		return ops
	}
	// Compile regexes once for the entire Apply call.
	inclPathRx := compileAll(f.IncludePaths)
	exclPathRx := compileAll(f.ExcludePaths)

	out := make([]*Operation, 0, len(ops))
	for _, op := range ops {
		if passes(op, inclPathRx, exclPathRx, f.IncludeTags, f.ExcludeTags) {
			out = append(out, op)
		}
	}
	return out
}

// compileAll compiles a slice of regex patterns. Patterns are assumed to have
// been validated by Validate(); errors are silently ignored here.
func compileAll(patterns []string) []*regexp.Regexp {
	rxs := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if rx, err := regexp.Compile(p); err == nil {
			rxs = append(rxs, rx)
		}
	}
	return rxs
}

func passes(op *Operation, inclPathRx, exclPathRx []*regexp.Regexp, inclTags, exclTags []string) bool {
	// IncludePath: op.Path must match at least one compiled pattern (if any set)
	if len(inclPathRx) > 0 {
		matched := false
		for _, rx := range inclPathRx {
			if rx.MatchString(op.Path) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	// ExcludePath: op.Path must not match any compiled pattern
	for _, rx := range exclPathRx {
		if rx.MatchString(op.Path) {
			return false
		}
	}
	// IncludeTag: op.Tags must intersect with IncludeTags (if any set)
	if len(inclTags) > 0 && !tagsIntersect(op.Tags, inclTags) {
		return false
	}
	// ExcludeTag: op.Tags must not intersect with ExcludeTags
	if tagsIntersect(op.Tags, exclTags) {
		return false
	}
	return true
}

// tagsIntersect reports whether a and b share at least one element.
func tagsIntersect(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x == y {
				return true
			}
		}
	}
	return false
}
