// cmd/opfilter.go
package cmd

import "github.com/testmind-hq/caseforge/internal/spec"

// buildFilterSet constructs a spec.FilterSet from flag string values.
// includePath and excludePath are single regex strings.
// includeTag and excludeTag are comma-separated tag names.
func buildFilterSet(includePath, excludePath, includeTag, excludeTag string) spec.FilterSet {
	f := spec.FilterSet{}
	if includePath != "" {
		f.IncludePaths = []string{includePath}
	}
	if excludePath != "" {
		f.ExcludePaths = []string{excludePath}
	}
	if includeTag != "" {
		f.IncludeTags = splitTrimmed(includeTag)
	}
	if excludeTag != "" {
		f.ExcludeTags = splitTrimmed(excludeTag)
	}
	return f
}
