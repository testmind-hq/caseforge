// cmd/opfilter_test.go
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestBuildFilterSet_Empty(t *testing.T) {
	f := buildFilterSet("", "", "", "")
	assert.True(t, f.IsEmpty())
}

func TestBuildFilterSet_IncludePath(t *testing.T) {
	f := buildFilterSet("^/users", "", "", "")
	assert.Equal(t, spec.FilterSet{IncludePaths: []string{"^/users"}}, f)
}

func TestBuildFilterSet_IncludeTag(t *testing.T) {
	f := buildFilterSet("", "", "users,admin", "")
	assert.Equal(t, []string{"users", "admin"}, f.IncludeTags)
}

func TestBuildFilterSet_ExcludeTag(t *testing.T) {
	f := buildFilterSet("", "", "", "deprecated")
	assert.Equal(t, []string{"deprecated"}, f.ExcludeTags)
}
