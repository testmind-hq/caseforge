// internal/methodology/required_omission_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func makeRequiredOmissionOp() *spec.Operation {
	return &spec.Operation{
		Method: "POST",
		Path:   "/users",
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"email", "name"},
						Properties: map[string]*spec.Schema{
							"email": {Type: "string", Format: "email"},
							"name":  {Type: "string"},
							"age":   {Type: "integer"},
						},
					},
					Examples: map[string]*spec.Example{
						"validUser": {
							Summary: "A valid user",
							Value: map[string]any{
								"email": "test@example.com",
								"name":  "Alice",
								"age":   30,
							},
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {}},
	}
}

func TestRequiredOmissionTechnique_Applies_WhenHasRequiredField(t *testing.T) {
	op := makeRequiredOmissionOp()
	tech := NewRequiredOmissionTechnique()
	assert.True(t, tech.Applies(op))
}

func TestRequiredOmissionTechnique_Applies_False_NoRequired(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/items",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"name": {Type: "string"},
					},
					// no Required slice
				}},
			},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}
	tech := NewRequiredOmissionTechnique()
	assert.False(t, tech.Applies(op))
}

func TestRequiredOmissionTechnique_Generate_OneCasePerRequiredField(t *testing.T) {
	op := makeRequiredOmissionOp()
	tech := NewRequiredOmissionTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)
	// email and name are required → 2 cases
	assert.Len(t, cases, 2)
}

func TestRequiredOmissionTechnique_Generate_FieldIsAbsent_NotNull(t *testing.T) {
	op := makeRequiredOmissionOp()
	tech := NewRequiredOmissionTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	for _, tc := range cases {
		body, ok := tc.Steps[0].Body.(map[string]any)
		require.True(t, ok, "body must be map[string]any")
		// Identify which field should be absent from Source.SpecPath
		// e.g. "POST /users requestBody.email" => "email" absent
		for _, req := range op.RequestBody.Content["application/json"].Schema.Required {
			if tc.Source.SpecPath == "POST /users requestBody."+req {
				_, present := body[req]
				assert.False(t, present, "required field %q should be absent (not null) in body", req)
			}
		}
	}
}

func TestRequiredOmissionTechnique_Generate_Expects4xx(t *testing.T) {
	op := makeRequiredOmissionOp()
	tech := NewRequiredOmissionTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	for _, tc := range cases {
		assertions := tc.Steps[0].Assertions
		require.Len(t, assertions, 2)
		assert.Equal(t, "gte", assertions[0].Operator)
		assert.Equal(t, 400, assertions[0].Expected)
		assert.Equal(t, "lt", assertions[1].Operator)
		assert.Equal(t, 500, assertions[1].Expected)
	}
}

func TestRequiredOmissionTechnique_Generate_UsesExampleAsSeed(t *testing.T) {
	op := makeRequiredOmissionOp()
	tech := NewRequiredOmissionTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	// The example has age=30; when we omit "email", age should still be present
	for _, tc := range cases {
		if tc.Source.SpecPath == "POST /users requestBody.email" {
			body, ok := tc.Steps[0].Body.(map[string]any)
			require.True(t, ok)
			// age comes from example seed
			_, hasAge := body["age"]
			assert.True(t, hasAge, "body should include age from example seed when email is omitted")
			_, hasEmail := body["email"]
			assert.False(t, hasEmail, "email should be absent")
		}
	}
}
