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

func TestBuildFilterSet_ExcludePath(t *testing.T) {
	f := buildFilterSet("", "^/admin", "", "")
	assert.Equal(t, spec.FilterSet{ExcludePaths: []string{"^/admin"}}, f)
}

func TestBuildFilterSet_ExcludeTag(t *testing.T) {
	f := buildFilterSet("", "", "", "deprecated")
	assert.Equal(t, []string{"deprecated"}, f.ExcludeTags)
}

func TestChainCommand_HasFilterFlags(t *testing.T) {
	assert.NotNil(t, chainCmd.Flags().Lookup("include-path"), "--include-path flag required")
	assert.NotNil(t, chainCmd.Flags().Lookup("exclude-path"), "--exclude-path flag required")
	assert.NotNil(t, chainCmd.Flags().Lookup("include-tag"), "--include-tag flag required")
	assert.NotNil(t, chainCmd.Flags().Lookup("exclude-tag"), "--exclude-tag flag required")
}

func TestExploreCommand_HasFilterFlags(t *testing.T) {
	assert.NotNil(t, exploreCmd.Flags().Lookup("include-path"), "--include-path flag required")
	assert.NotNil(t, exploreCmd.Flags().Lookup("exclude-path"), "--exclude-path flag required")
	assert.NotNil(t, exploreCmd.Flags().Lookup("include-tag"), "--include-tag flag required")
	assert.NotNil(t, exploreCmd.Flags().Lookup("exclude-tag"), "--exclude-tag flag required")
}
