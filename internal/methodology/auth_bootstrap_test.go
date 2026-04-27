package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func makeBootstrapSpec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		GlobalSecurity: []string{"bearerAuth"},
		SecuritySchemes: []string{"bearerAuth"},
		Operations: []*spec.Operation{
			{
				Method: "POST", Path: "/auth/login",
				Responses: map[string]*spec.Response{
					"200": {Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"token": {Type: "string"},
							},
						}},
					}},
				},
			},
			{
				Method: "GET", Path: "/users",
				Security: []string{"bearerAuth"},
				Responses: map[string]*spec.Response{"200": {Description: "ok"}},
			},
		},
	}
}

func TestBootstrapAuth_ConvertsSingleCaseToChain(t *testing.T) {
	s := makeBootstrapSpec()
	input := []schema.TestCase{
		{
			ID:   "TC-0001",
			Kind: "single",
			Source: schema.CaseSource{Technique: "equivalence_partitioning", SpecPath: "GET /users"},
			Steps: []schema.Step{
				{ID: "step-1", Method: "GET", Path: "/users", Type: "test"},
			},
		},
	}
	result := BootstrapAuth(input, s)
	require.Len(t, result, 1)
	assert.Equal(t, "chain", result[0].Kind)
	require.Len(t, result[0].Steps, 2)
	assert.Equal(t, "step-auth", result[0].Steps[0].ID)
	assert.Equal(t, "step-1", result[0].Steps[1].ID)
	assert.Equal(t, []string{"step-auth"}, result[0].Steps[1].DependsOn)
}

func TestBootstrapAuth_NoAuthOp_ReturnsUnchanged(t *testing.T) {
	s := &spec.ParsedSpec{
		GlobalSecurity: []string{"bearerAuth"},
		Operations:     []*spec.Operation{},
	}
	input := []schema.TestCase{{ID: "TC-0001", Kind: "single"}}
	result := BootstrapAuth(input, s)
	require.Len(t, result, 1)
	assert.Equal(t, "single", result[0].Kind)
}

func TestBootstrapAuth_SkipsChainCases(t *testing.T) {
	s := makeBootstrapSpec()
	input := []schema.TestCase{
		{
			ID:   "TC-chain",
			Kind: "chain",
			Steps: []schema.Step{
				{ID: "step-auth", Type: "setup"},
				{ID: "step-test", Type: "test"},
			},
		},
	}
	result := BootstrapAuth(input, s)
	require.Len(t, result, 1)
	assert.Equal(t, "chain", result[0].Kind)
	assert.Len(t, result[0].Steps, 2) // unchanged
}

func TestBootstrapAuth_UnsecuredOp_Unchanged(t *testing.T) {
	s := &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{Method: "POST", Path: "/auth/login",
				Responses: map[string]*spec.Response{
					"200": {Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{"token": {Type: "string"}},
						}},
					}},
				},
			},
			{Method: "GET", Path: "/public", Responses: map[string]*spec.Response{"200": {Description: "ok"}}},
		},
	}
	input := []schema.TestCase{
		{
			ID: "TC-pub", Kind: "single",
			Source: schema.CaseSource{SpecPath: "GET /public"},
			Steps: []schema.Step{{ID: "step-1", Method: "GET", Path: "/public", Type: "test"}},
		},
	}
	result := BootstrapAuth(input, s)
	require.Len(t, result, 1)
	assert.Equal(t, "single", result[0].Kind)
}
