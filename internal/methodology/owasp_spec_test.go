package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func apiSpecWithRBAC() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/users/me", Responses: map[string]*spec.Response{"200": {}}},
			{Method: "GET", Path: "/admin/users", Responses: map[string]*spec.Response{"200": {}}},
			{Method: "DELETE", Path: "/users/{id}", Responses: map[string]*spec.Response{"204": {}}},
		},
	}
}

func apiSpecWithVersions() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/v1/users", Responses: map[string]*spec.Response{"200": {}}},
			{Method: "GET", Path: "/v2/users", Responses: map[string]*spec.Response{"200": {}}},
			{Method: "POST", Path: "/v2/orders", Responses: map[string]*spec.Response{"201": {}}},
		},
	}
}

func TestSecuritySpecAPI5_FunctionLevel(t *testing.T) {
	sst := NewSecuritySpecTechnique()
	cases, err := sst.Generate(apiSpecWithRBAC())
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api5-function-level-auth" {
				found = true
				assert.Equal(t, "P0", tc.Priority)
				assert.Contains(t, tc.Steps[0].Headers["Authorization"], "{{user_token}}")
				assert.Equal(t, 403, tc.Steps[0].Assertions[0].Expected)
			}
		}
	}
	assert.True(t, found, "should generate API5 cases when low+high privilege paths exist")
}

func TestSecuritySpecAPI8_CORS(t *testing.T) {
	sst := NewSecuritySpecTechnique()
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
			{Method: "POST", Path: "/users", Responses: map[string]*spec.Response{"201": {}}},
			{Method: "GET", Path: "/orders", Responses: map[string]*spec.Response{"200": {}}},
		},
	}
	cases, err := sst.Generate(ps)
	require.NoError(t, err)
	var corsCases []string
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api8-cors" {
				corsCases = append(corsCases, tc.Steps[0].Path)
			}
		}
	}
	// /users appears twice (GET+POST) but should produce only 1 CORS case
	assert.Equal(t, 2, len(corsCases), "one CORS case per unique path")
	assert.NotEqual(t, corsCases[0], corsCases[1])
}

func TestSecuritySpecAPI9_AssetManagement(t *testing.T) {
	sst := NewSecuritySpecTechnique()
	cases, err := sst.Generate(apiSpecWithVersions())
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api9-asset-management" {
				found = true
				assert.Equal(t, 404, tc.Steps[0].Assertions[0].Expected)
				assert.Contains(t, tc.Steps[0].Path, "/v1/")
			}
		}
	}
	assert.True(t, found)
}

func TestSecuritySpecNoAPI5_WithoutLowPriv(t *testing.T) {
	sst := NewSecuritySpecTechnique()
	ps := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		},
	}
	cases, err := sst.Generate(ps)
	require.NoError(t, err)
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			assert.NotEqual(t, "api5-function-level-auth", tag)
		}
	}
}
