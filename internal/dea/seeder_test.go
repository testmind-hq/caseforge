// internal/dea/seeder_test.go
package dea

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestSeedHypotheses_RequiredField(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	require.NotEmpty(t, nodes)

	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindRequiredField, "must seed required-field hypothesis for 'name'")
}

func TestSeedHypotheses_OptionalField(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
							"tag":  {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var optionalKinds int
	for _, n := range nodes {
		if n.Kind == KindOptionalField {
			optionalKinds++
		}
	}
	assert.Equal(t, 1, optionalKinds, "exactly 1 optional field hypothesis for 'tag'")
}

func TestSeedHypotheses_StringConstraints(t *testing.T) {
	maxL := int64(100)
	minL := int64(1)
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string", MaxLength: &maxL, MinLength: &minL},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindStringMaxLength)
	assert.Contains(t, kinds, KindStringMinLength)
}

func TestSeedHypotheses_ImplicitStringConstraints(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"}, // no MinLength / MaxLength
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindStringImplicitMin)
	assert.Contains(t, kinds, KindStringImplicitMax)
}

func TestSeedHypotheses_NumericConstraints(t *testing.T) {
	min := float64(1)
	max := float64(100)
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"age": {Type: "integer", Minimum: &min, Maximum: &max},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindNumericMin)
	assert.Contains(t, kinds, KindNumericMax)
}

func TestSeedHypotheses_EnumField(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/orders",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"status": {Type: "string", Enum: []any{"active", "inactive"}},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindEnumViolation)
}

func TestSeedHypotheses_NoBodyOperation(t *testing.T) {
	op := &spec.Operation{
		Method: "GET",
		Path:   "/pets",
	}
	nodes := SeedHypotheses(op)
	for _, n := range nodes {
		assert.NotEqual(t, KindRequiredField, n.Kind, "no body hypotheses for bodyless GET")
	}
}

func TestSeedHypotheses_QueryParamConstraints(t *testing.T) {
	min := float64(1)
	max := float64(100)
	op := &spec.Operation{
		Method: "GET",
		Path:   "/pets",
		Parameters: []*spec.Parameter{
			{
				Name: "limit",
				In:   "query",
				Schema: &spec.Schema{
					Type:    "integer",
					Minimum: &min,
					Maximum: &max,
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindNumericMin)
	assert.Contains(t, kinds, KindNumericMax)
}

func TestSeedHypotheses_AllNodesHavePendingStatus(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	for _, n := range nodes {
		assert.Equal(t, StatusPending, n.Status, "all seeded nodes must start as StatusPending")
	}
}

func TestSeedHypotheses_NumericEnumField(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/orders",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"priority": {Type: "integer", Enum: []any{1, 2, 3}},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindEnumViolation, "numeric enum fields must get KindEnumViolation hypothesis")
}

func TestSeedHypotheses_AllNodesHaveOperationSet(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	for _, n := range nodes {
		assert.Equal(t, "POST /pets", n.Operation)
		assert.NotEmpty(t, n.ID)
		assert.NotEmpty(t, n.FieldPath)
	}
}

func TestSeedHypotheses_ArrayConstraints(t *testing.T) {
	minI := uint64(1)
	maxI := uint64(10)
	op := &spec.Operation{
		Method: "POST",
		Path:   "/orders",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"tags": {Type: "array", MinItems: &minI, MaxItems: &maxI},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindArrayMinItems)
	assert.Contains(t, kinds, KindArrayMaxItems)
}

func TestSeedHypotheses_ArrayNoConstraints_NoArrayKinds(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/orders",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"tags": {Type: "array"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	for _, n := range nodes {
		assert.NotEqual(t, KindArrayMinItems, n.Kind)
		assert.NotEqual(t, KindArrayMaxItems, n.Kind)
	}
}

func TestSeedHypotheses_RequiredQueryParam(t *testing.T) {
	op := &spec.Operation{
		Method: "GET",
		Path:   "/pets",
		Parameters: []*spec.Parameter{
			{
				Name:     "status",
				In:       "query",
				Required: true,
				Schema:   &spec.Schema{Type: "string"},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindRequiredQueryParam)
}

func TestSeedHypotheses_OptionalQueryParam_NoRequiredKind(t *testing.T) {
	op := &spec.Operation{
		Method: "GET",
		Path:   "/pets",
		Parameters: []*spec.Parameter{
			{
				Name:     "status",
				In:       "query",
				Required: false,
				Schema:   &spec.Schema{Type: "string"},
			},
		},
	}

	nodes := SeedHypotheses(op)
	for _, n := range nodes {
		assert.NotEqual(t, KindRequiredQueryParam, n.Kind)
	}
}

func TestSeedHypotheses_FormatViolation(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/users",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"email"},
						Properties: map[string]*spec.Schema{
							"email": {Type: "string", Format: "email"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindFormatViolation)
}

func TestSeedHypotheses_IncludesTypeCoercion(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
							"age":  {Type: "integer"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindTypeCoercion, "must seed type_coercion hypothesis for typed fields")
}

func TestSeedHypotheses_IncludesUnicodeControl(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var kinds []HypothesisKind
	for _, n := range nodes {
		kinds = append(kinds, n.Kind)
	}
	assert.Contains(t, kinds, KindUnicodeControl, "must seed unicode_control hypothesis for string fields")
}

func TestSeedHypotheses_IncludesMassAssignment(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	var massCount int
	for _, n := range nodes {
		if n.Kind == KindMassAssignment {
			massCount++
		}
	}
	require.Equal(t, 1, massCount, "must seed exactly one mass_assignment hypothesis per operation")
}

func TestSeedHypotheses_NoFormatViolation_WhenNoFormat(t *testing.T) {
	op := &spec.Operation{
		Method: "POST",
		Path:   "/pets",
		RequestBody: &spec.RequestBody{
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}

	nodes := SeedHypotheses(op)
	for _, n := range nodes {
		assert.NotEqual(t, KindFormatViolation, n.Kind)
	}
}
