// internal/spec/filter_test.go
package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterSet_IsEmpty(t *testing.T) {
	assert.True(t, (&FilterSet{}).IsEmpty())
	assert.False(t, (&FilterSet{IncludePaths: []string{"/users"}}).IsEmpty())
	assert.False(t, (&FilterSet{ExcludePaths: []string{"/admin"}}).IsEmpty())
	assert.False(t, (&FilterSet{IncludeTags: []string{"auth"}}).IsEmpty())
	assert.False(t, (&FilterSet{ExcludeTags: []string{"deprecated"}}).IsEmpty())
}

func TestFilterSet_Validate_InvalidRegex(t *testing.T) {
	err := (&FilterSet{IncludePaths: []string{"[invalid"}}).Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--include-path")

	err = (&FilterSet{ExcludePaths: []string{"[bad"}}).Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--exclude-path")
}

func TestFilterSet_Validate_ValidRegex(t *testing.T) {
	f := &FilterSet{
		IncludePaths: []string{"^/users/.*"},
		ExcludePaths: []string{"^/admin"},
	}
	assert.NoError(t, f.Validate())
}

func TestFilterSet_Apply_EmptyReturnsAll(t *testing.T) {
	ops := []*Operation{{Path: "/a"}, {Path: "/b"}}
	result := (&FilterSet{}).Apply(ops)
	assert.Equal(t, ops, result)
}

func TestFilterSet_Apply_IncludePath(t *testing.T) {
	ops := []*Operation{
		{Path: "/users", Method: "GET"},
		{Path: "/users/1", Method: "GET"},
		{Path: "/orders", Method: "GET"},
	}
	result := (&FilterSet{IncludePaths: []string{"/users"}}).Apply(ops)
	assert.Len(t, result, 2) // /users and /users/1 both contain "/users"
}

func TestFilterSet_Apply_IncludePath_AnchoredRegex(t *testing.T) {
	ops := []*Operation{
		{Path: "/users", Method: "GET"},
		{Path: "/admin/users", Method: "GET"},
	}
	result := (&FilterSet{IncludePaths: []string{"^/users"}}).Apply(ops)
	assert.Len(t, result, 1)
	assert.Equal(t, "/users", result[0].Path)
}

func TestFilterSet_Apply_ExcludePath(t *testing.T) {
	ops := []*Operation{
		{Path: "/users", Method: "GET"},
		{Path: "/admin/users", Method: "GET"},
	}
	result := (&FilterSet{ExcludePaths: []string{"^/admin"}}).Apply(ops)
	assert.Len(t, result, 1)
	assert.Equal(t, "/users", result[0].Path)
}

func TestFilterSet_Apply_IncludeTag(t *testing.T) {
	ops := []*Operation{
		{Path: "/a", Tags: []string{"users", "admin"}},
		{Path: "/b", Tags: []string{"orders"}},
		{Path: "/c"},
	}
	result := (&FilterSet{IncludeTags: []string{"users"}}).Apply(ops)
	assert.Len(t, result, 1)
	assert.Equal(t, "/a", result[0].Path)
}

func TestFilterSet_Apply_ExcludeTag(t *testing.T) {
	ops := []*Operation{
		{Path: "/a", Tags: []string{"users"}},
		{Path: "/b", Tags: []string{"deprecated"}},
		{Path: "/c", Tags: []string{"users", "deprecated"}},
	}
	result := (&FilterSet{ExcludeTags: []string{"deprecated"}}).Apply(ops)
	assert.Len(t, result, 1)
	assert.Equal(t, "/a", result[0].Path)
}

func TestFilterSet_Apply_IncludeTagNoMatch_ReturnsNone(t *testing.T) {
	ops := []*Operation{{Path: "/a", Tags: []string{"users"}}}
	result := (&FilterSet{IncludeTags: []string{"admin"}}).Apply(ops)
	assert.Empty(t, result)
}

func TestFilterSet_Apply_CombinedIncludeExclude(t *testing.T) {
	ops := []*Operation{
		{Path: "/users", Tags: []string{"users"}},
		{Path: "/users/admin", Tags: []string{"users", "admin"}},
		{Path: "/orders", Tags: []string{"orders"}},
	}
	f := &FilterSet{
		IncludeTags:  []string{"users"},
		ExcludePaths: []string{"/admin"},
	}
	result := f.Apply(ops)
	assert.Len(t, result, 1)
	assert.Equal(t, "/users", result[0].Path)
}
