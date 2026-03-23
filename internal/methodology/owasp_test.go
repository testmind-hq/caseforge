package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func securitySpec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				Method: "GET", Path: "/users/{userId}",
				Security:   []string{"bearerAuth"},
				Parameters: []*spec.Parameter{{Name: "userId", In: "path", Required: true, Schema: &spec.Schema{Type: "integer"}}},
				Responses:  map[string]*spec.Response{"200": {Description: "OK"}},
			},
			{
				Method: "PATCH", Path: "/users/{userId}",
				Security: []string{"bearerAuth"},
				RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "object",
						Properties: map[string]*spec.Schema{"name": {Type: "string"}}}},
				}},
				Responses: map[string]*spec.Response{"200": {Description: "OK"}},
			},
			{
				Method: "GET", Path: "/items",
				Parameters: []*spec.Parameter{
					{Name: "limit", In: "query", Schema: &spec.Schema{Type: "integer"}},
					{Name: "q", In: "query", Schema: &spec.Schema{Type: "string"}},
				},
				Responses: map[string]*spec.Response{"200": {Description: "OK"}},
			},
		},
	}
}

func TestSecurityTechniqueAPI1_BOLA(t *testing.T) {
	st := NewSecurityTechnique()
	op := securitySpec().Operations[0] // GET /users/{userId}
	require.True(t, st.Applies(op))
	cases, err := st.Generate(op)
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api1-bola" {
				found = true
				assert.Equal(t, "P0", tc.Priority)
				assert.Contains(t, tc.Steps[0].Path, "{{other_resource_id}}")
				assert.Equal(t, 403, tc.Steps[0].Assertions[0].Expected)
			}
		}
	}
	assert.True(t, found, "should generate API1 BOLA case")
}

func TestSecurityTechniqueAPI2_Auth(t *testing.T) {
	st := NewSecurityTechnique()
	op := securitySpec().Operations[0] // GET /users/{userId} has Security
	cases, err := st.Generate(op)
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api2-broken-auth" {
				found = true
				_, hasAuth := tc.Steps[0].Headers["Authorization"]
				assert.False(t, hasAuth, "API2 case must not have Authorization header")
				assert.Equal(t, 401, tc.Steps[0].Assertions[0].Expected)
			}
		}
	}
	assert.True(t, found)
}

func TestSecurityTechniqueAPI3_BOPLA(t *testing.T) {
	st := NewSecurityTechnique()
	op := securitySpec().Operations[1] // PATCH /users/{userId}
	require.True(t, st.Applies(op))
	cases, err := st.Generate(op)
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api3-bopla" {
				found = true
				body, ok := tc.Steps[0].Body.(map[string]any)
				require.True(t, ok)
				assert.Equal(t, true, body["is_admin"])
				assert.Equal(t, "admin", body["role"])
			}
		}
	}
	assert.True(t, found)
}

func TestSecurityTechniqueAPI4_Pagination(t *testing.T) {
	st := NewSecurityTechnique()
	op := securitySpec().Operations[2] // GET /items has limit param
	require.True(t, st.Applies(op))
	cases, err := st.Generate(op)
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api4-resource-consumption" {
				found = true
				assert.Contains(t, tc.Steps[0].Path, "limit=99999")
			}
		}
	}
	assert.True(t, found)
}

func TestSecurityTechniqueAPI7_Injection(t *testing.T) {
	st := NewSecurityTechnique()
	op := securitySpec().Operations[2] // GET /items has string param "q"
	cases, err := st.Generate(op)
	require.NoError(t, err)
	var count int
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api7-injection" {
				count++
			}
		}
	}
	assert.Equal(t, 3, count, "should generate 3 injection sub-cases (XSS, SQLi, path traversal)")
}

func TestSecurityTechniqueAPI10_SSRF(t *testing.T) {
	st := NewSecurityTechnique()
	op := &spec.Operation{
		Method: "POST", Path: "/webhooks",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{Type: "object",
				Properties: map[string]*spec.Schema{"url": {Type: "string"}}}},
		}},
		Responses: map[string]*spec.Response{"200": {Description: "OK"}},
	}
	require.True(t, st.Applies(op))
	cases, err := st.Generate(op)
	require.NoError(t, err)
	var found bool
	for _, tc := range cases {
		for _, tag := range tc.Tags {
			if tag == "api10-ssrf" {
				found = true
			}
		}
	}
	assert.True(t, found)
}

func TestSecurityTechniqueNotApplies(t *testing.T) {
	st := NewSecurityTechnique()
	op := &spec.Operation{
		Method: "GET", Path: "/health",
		Responses: map[string]*spec.Response{"200": {}},
	}
	// No ID param, no security, no body, no pagination, no string params
	assert.False(t, st.Applies(op))
}
