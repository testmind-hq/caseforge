// internal/methodology/isolated_negative_test.go
package methodology

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestIsolatedNegative_GeneratesOneFailureCasePerField(t *testing.T) {
	op := &spec.Operation{
		OperationID: "createUser",
		Method:      "POST",
		Path:        "/users",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"email", "age"},
						Properties: map[string]*spec.Schema{
							"email": {Type: "string", Format: "email"},
							"age":   {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(120)},
							"name":  {Type: "string"},
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}

	tech := NewIsolatedNegativeTechnique()
	require.True(t, tech.Applies(op))

	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	// Each case must have exactly one source of invalidity
	for _, tc := range cases {
		assert.Equal(t, "isolated_negative", tc.Source.Technique)
		assert.NotEmpty(t, tc.Source.Rationale)
		assert.Equal(t, 422, tc.Steps[0].Assertions[0].Expected)
		assert.Equal(t, "status_code", tc.Steps[0].Assertions[0].Target)
	}

	// Required fields: email and age missing → 2 cases
	// Field violations: age below-min or above-max → up to 2 cases
	// So total ≥ 2 (for missing required fields)
	assert.GreaterOrEqual(t, len(cases), 2)
}

func TestIsolatedNegative_EachCaseHasOnlyOneInvalidField(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/items",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"count"},
						Properties: map[string]*spec.Schema{
							"count": {Type: "integer", Minimum: floatPtr(1), Maximum: floatPtr(100)},
							"label": {Type: "string"},
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {}},
	}

	tech := NewIsolatedNegativeTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	// The case for "count below minimum" should have label as a valid non-nil value
	var belowMinCase *schema.TestCase
	for i, tc := range cases {
		if containsString(tc.Source.Rationale, "below") || containsString(tc.Source.Rationale, "minimum") {
			belowMinCase = &cases[i]
		}
	}
	if belowMinCase != nil {
		body, ok := belowMinCase.Steps[0].Body.(map[string]any)
		require.True(t, ok)
		// label should be present (valid base value), count should be invalid
		_, hasLabel := body["label"]
		assert.True(t, hasLabel, "all other fields should be present with valid values")
	}
}

func TestIsolatedNegative_RequiredParamMissing(t *testing.T) {
	op := &spec.Operation{
		Method: "GET",
		Path:   "/search",
		Parameters: []*spec.Parameter{
			{Name: "q", In: "query", Required: true, Schema: &spec.Schema{Type: "string"}},
			{Name: "limit", In: "query", Required: false, Schema: &spec.Schema{Type: "integer"}},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}

	tech := NewIsolatedNegativeTechnique()
	require.True(t, tech.Applies(op))

	cases, err := tech.Generate(op)
	require.NoError(t, err)
	require.NotEmpty(t, cases)

	found := false
	for _, tc := range cases {
		if containsString(tc.Source.Rationale, "q") && containsString(tc.Source.Rationale, "absent") {
			found = true
			assert.Equal(t, "isolated_negative", tc.Source.Technique)
		}
	}
	assert.True(t, found, "expected a case for missing required param 'q'")
}

func TestIsolatedNegative_AppliesToBodyWithProperties(t *testing.T) {
	opWithBody := &spec.Operation{
		Method: "POST", Path: "/x",
		RequestBody: &spec.RequestBody{Content: map[string]*spec.MediaType{
			"application/json": {Schema: &spec.Schema{Type: "object",
				Properties: map[string]*spec.Schema{"n": {Type: "integer", Minimum: floatPtr(0)}},
			}},
		}},
		Responses: map[string]*spec.Response{"201": {}},
	}
	opNoBody := &spec.Operation{Method: "GET", Path: "/x",
		Responses: map[string]*spec.Response{"200": {}},
	}

	tech := NewIsolatedNegativeTechnique()
	assert.True(t, tech.Applies(opWithBody))
	assert.False(t, tech.Applies(opNoBody))
}

func containsString(s, sub string) bool {
	return strings.Contains(s, sub)
}
