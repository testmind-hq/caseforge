// internal/methodology/semantic_annotation_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// makeNullableOp builds a POST operation with a nullable field in its JSON body.
func makeNullableOp() *spec.Operation {
	return &spec.Operation{
		Method: "POST", Path: "/users",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"email":    {Type: "string"},
					"nickname": {Type: "string", Nullable: true},
				},
				Required: []string{"email"},
			}},
		}},
		Responses: map[string]*spec.Response{"201": {Description: "created"}},
	}
}

// makeReadOnlyOp builds a POST operation with a readOnly field in its JSON body.
func makeReadOnlyOp() *spec.Operation {
	return &spec.Operation{
		Method: "POST", Path: "/users",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"email": {Type: "string"},
					"id":    {Type: "integer", ReadOnly: true},
				},
				Required: []string{"email"},
			}},
		}},
		Responses: map[string]*spec.Response{"201": {Description: "created"}},
	}
}

// makeWriteOnlyGETOp builds a GET operation whose 200 response schema has a writeOnly field.
func makeWriteOnlyGETOp() *spec.Operation {
	return &spec.Operation{
		Method: "GET", Path: "/users/1",
		Responses: map[string]*spec.Response{
			"200": {
				Description: "ok",
				Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"id":       {Type: "integer"},
							"email":    {Type: "string"},
							"password": {Type: "string", WriteOnly: true},
						},
					}},
				},
			},
		},
	}
}

// makeNoAnnotationsOp builds a GET operation with no semantic annotations.
func makeNoAnnotationsOp() *spec.Operation {
	return &spec.Operation{
		Method: "GET", Path: "/items",
		Responses: map[string]*spec.Response{
			"200": {
				Description: "ok",
				Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"id":   {Type: "integer"},
							"name": {Type: "string"},
						},
					}},
				},
			},
		},
	}
}

func TestSemanticAnnotationTechnique_Applies_Nullable(t *testing.T) {
	op := makeNullableOp()
	assert.True(t, NewSemanticAnnotationTechnique().Applies(op))
}

func TestSemanticAnnotationTechnique_Applies_ReadOnly(t *testing.T) {
	op := makeReadOnlyOp()
	assert.True(t, NewSemanticAnnotationTechnique().Applies(op))
}

func TestSemanticAnnotationTechnique_Applies_WriteOnly_GET(t *testing.T) {
	op := makeWriteOnlyGETOp()
	assert.True(t, NewSemanticAnnotationTechnique().Applies(op))
}

func TestSemanticAnnotationTechnique_Applies_False_NoAnnotations(t *testing.T) {
	op := makeNoAnnotationsOp()
	assert.False(t, NewSemanticAnnotationTechnique().Applies(op))
}

func TestSemanticAnnotationTechnique_Generate_NullableCase_Expects2xx(t *testing.T) {
	op := makeNullableOp()
	cases, err := NewSemanticAnnotationTechnique().Generate(op)
	require.NoError(t, err)

	var nullableCase *spec.Operation
	_ = nullableCase

	found := false
	for _, tc := range cases {
		if tc.Source.Scenario != "NULLABLE_ACCEPTANCE" {
			continue
		}
		require.Len(t, tc.Steps, 1)
		// Body must contain nil for the nullable field
		body, ok := tc.Steps[0].Body.(map[string]any)
		require.True(t, ok)
		val, exists := body["nickname"]
		assert.True(t, exists, "nullable field 'nickname' should be in body")
		assert.Nil(t, val, "nullable field should be set to nil")

		// Assertions: status_code gte 200 AND lt 300
		operators := make(map[string]bool)
		for _, a := range tc.Steps[0].Assertions {
			if a.Target == "status_code" {
				operators[a.Operator] = true
			}
		}
		assert.True(t, operators["gte"], "should have gte assertion on status_code")
		assert.True(t, operators["lt"], "should have lt assertion on status_code")
		found = true
	}
	assert.True(t, found, "expected at least one NULLABLE_ACCEPTANCE case")
}

func TestSemanticAnnotationTechnique_Generate_ReadOnlyCase_Expects4xx(t *testing.T) {
	op := makeReadOnlyOp()
	cases, err := NewSemanticAnnotationTechnique().Generate(op)
	require.NoError(t, err)

	found := false
	for _, tc := range cases {
		if tc.Source.Scenario != "READ_ONLY_WRITE" {
			continue
		}
		require.Len(t, tc.Steps, 1)

		// Body must include the readOnly field
		body, ok := tc.Steps[0].Body.(map[string]any)
		require.True(t, ok)
		_, hasID := body["id"]
		assert.True(t, hasID, "readOnly field 'id' should be in body")

		// Assertions: status_code gte 400 AND lt 500
		operators := make(map[string]bool)
		for _, a := range tc.Steps[0].Assertions {
			if a.Target == "status_code" {
				operators[a.Operator] = true
			}
		}
		assert.True(t, operators["gte"], "should have gte assertion on status_code")
		assert.True(t, operators["lt"], "should have lt assertion on status_code")

		assert.Equal(t, "P2", tc.Priority)
		found = true
	}
	assert.True(t, found, "expected at least one READ_ONLY_WRITE case")
}

func TestSemanticAnnotationTechnique_Generate_WriteOnlyCase_FieldAbsent(t *testing.T) {
	op := makeWriteOnlyGETOp()
	cases, err := NewSemanticAnnotationTechnique().Generate(op)
	require.NoError(t, err)

	found := false
	for _, tc := range cases {
		if tc.Source.Scenario != "WRITE_ONLY_READ" {
			continue
		}
		require.Len(t, tc.Steps, 1)

		// Check jsonpath assertion for the writeOnly field
		hasJSONPathAssertion := false
		for _, a := range tc.Steps[0].Assertions {
			if a.Target == "jsonpath $.password" {
				hasJSONPathAssertion = true
				assert.Equal(t, schema.OperatorNotExists, a.Operator)
			}
		}
		assert.True(t, hasJSONPathAssertion, "should have jsonpath assertion for 'password' field")
		assert.Equal(t, "P2", tc.Priority)
		found = true
	}
	assert.True(t, found, "expected at least one WRITE_ONLY_READ case")
}

func TestSemanticAnnotationTechnique_Generate_Scenario_Populated(t *testing.T) {
	// Build an op that has all three annotation types
	op := &spec.Operation{
		Method: "POST", Path: "/items",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"name":     {Type: "string"},
					"id":       {Type: "integer", ReadOnly: true},
					"nickname": {Type: "string", Nullable: true},
				},
				Required: []string{"name"},
			}},
		}},
		Responses: map[string]*spec.Response{"201": {Description: "created"}},
	}
	cases, err := NewSemanticAnnotationTechnique().Generate(op)
	require.NoError(t, err)

	scenarios := make(map[string]bool)
	for _, tc := range cases {
		assert.NotEmpty(t, tc.Source.Scenario, "Source.Scenario should be set for case %s", tc.ID)
		scenarios[tc.Source.Scenario] = true
	}

	assert.True(t, scenarios["NULLABLE_ACCEPTANCE"], "NULLABLE_ACCEPTANCE scenario should be present")
	assert.True(t, scenarios["READ_ONLY_WRITE"], "READ_ONLY_WRITE scenario should be present")
}
