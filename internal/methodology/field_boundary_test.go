// internal/methodology/field_boundary_test.go
package methodology

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/score"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func makeFieldBoundaryOp() *spec.Operation {
	minVal := float64(1)
	maxVal := float64(100)
	minLen := int64(5)
	maxLen := int64(10)
	return &spec.Operation{
		Method: "POST",
		Path:   "/orders",
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type:     "object",
					Required: []string{"quantity"},
					Properties: map[string]*spec.Schema{
						"quantity": {Type: "integer", Minimum: &minVal, Maximum: &maxVal},
						"address": {
							Type: "object",
							Properties: map[string]*spec.Schema{
								"zip": {Type: "string", MinLength: &minLen, MaxLength: &maxLen},
							},
						},
					},
				}},
			},
		},
		Responses: map[string]*spec.Response{"201": {}},
	}
}

func TestFieldBoundaryTechnique_Applies_WhenHasBoundaryConstraint(t *testing.T) {
	op := makeFieldBoundaryOp()
	tech := NewFieldBoundaryTechnique()
	assert.True(t, tech.Applies(op))
}

func TestFieldBoundaryTechnique_Applies_False_NoBoundary(t *testing.T) {
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
				}},
			},
		},
		Responses: map[string]*spec.Response{"200": {}},
	}
	tech := NewFieldBoundaryTechnique()
	assert.False(t, tech.Applies(op))
}

func TestFieldBoundaryTechnique_Generate_4CasesPerConstrainedField(t *testing.T) {
	// quantity has both Minimum and Maximum → 4 cases (valid_min, invalid_below_min, valid_max, invalid_above_max)
	op := makeFieldBoundaryOp()
	tech := NewFieldBoundaryTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	// Count cases for quantity
	var qCases int
	for _, tc := range cases {
		if strings.Contains(tc.Source.SpecPath, "quantity") {
			qCases++
		}
	}
	assert.Equal(t, 4, qCases, "expected 4 cases for quantity (min+max boundaries)")
}

func TestFieldBoundaryTechnique_Generate_ValidBoundaryExpects2xx(t *testing.T) {
	op := makeFieldBoundaryOp()
	tech := NewFieldBoundaryTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	for _, tc := range cases {
		if tc.Source.Scenario != score.ScenarioFieldBoundaryValid {
			continue
		}
		assertions := tc.Steps[0].Assertions
		require.Len(t, assertions, 2)
		assert.Equal(t, "gte", assertions[0].Operator)
		assert.Equal(t, 200, assertions[0].Expected)
		assert.Equal(t, "lt", assertions[1].Operator)
		assert.Equal(t, 300, assertions[1].Expected)
	}
}

func TestFieldBoundaryTechnique_Generate_InvalidBoundaryExpects4xx(t *testing.T) {
	op := makeFieldBoundaryOp()
	tech := NewFieldBoundaryTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	for _, tc := range cases {
		if tc.Source.Scenario != score.ScenarioFieldBoundaryInvalid {
			continue
		}
		assertions := tc.Steps[0].Assertions
		require.Len(t, assertions, 2)
		assert.Equal(t, "gte", assertions[0].Operator)
		assert.Equal(t, 400, assertions[0].Expected)
		assert.Equal(t, "lt", assertions[1].Operator)
		assert.Equal(t, 500, assertions[1].Expected)
	}
}

func TestFieldBoundaryTechnique_Generate_NestedField_DotPath(t *testing.T) {
	op := makeFieldBoundaryOp()
	tech := NewFieldBoundaryTechnique()
	cases, err := tech.Generate(op)
	require.NoError(t, err)

	var foundNested bool
	for _, tc := range cases {
		if strings.Contains(tc.Source.SpecPath, "address.zip") {
			foundNested = true
			break
		}
	}
	assert.True(t, foundNested, "expected at least one case with dot-path 'address.zip'")
}
