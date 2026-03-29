package dea

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func testOpForDesigner() *spec.Operation {
	maxL := int64(50)
	minL := int64(1)
	maxN := float64(100)
	minN := float64(1)
	return &spec.Operation{
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
							"age":  {Type: "integer", Minimum: &minN, Maximum: &maxN},
							"tag":  {Type: "string"},
						},
					},
				},
			},
		},
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
}

func newTestGen() *datagen.Generator {
	// datagen.NewGenerator accepts any (llm), nil is valid for tests
	return datagen.NewGenerator(nil)
}

func TestDesignProbe_RequiredField_OmitsField(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindRequiredField, FieldPath: "requestBody.name"}
	probe := DesignProbe(h, op, newTestGen())

	body, ok := probe.Body.(map[string]any)
	require.True(t, ok)
	_, hasName := body["name"]
	assert.False(t, hasName, "required-field probe must omit the target field")
	assert.Equal(t, 400, probe.ExpectedStatus)
}

func TestDesignProbe_OptionalField_OmitsField(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindOptionalField, FieldPath: "requestBody.tag"}
	probe := DesignProbe(h, op, newTestGen())

	body, ok := probe.Body.(map[string]any)
	require.True(t, ok)
	_, hasTag := body["tag"]
	assert.False(t, hasTag, "optional-field probe must omit the target field")
	assert.Equal(t, 201, probe.ExpectedStatus)
}

func TestDesignProbe_StringMaxLength_ExceedsMax(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindStringMaxLength, FieldPath: "requestBody.name"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	name, ok := body["name"].(string)
	require.True(t, ok)
	assert.Greater(t, len(name), 50)
	assert.Equal(t, 400, probe.ExpectedStatus)
}

func TestDesignProbe_StringMinLength_BelowMin(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindStringMinLength, FieldPath: "requestBody.name"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	name, ok := body["name"].(string)
	require.True(t, ok)
	assert.Equal(t, "", name)
	assert.Equal(t, 400, probe.ExpectedStatus)
}

func TestDesignProbe_ImplicitMax_Sends256Chars(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindStringImplicitMax, FieldPath: "requestBody.tag"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	tag, ok := body["tag"].(string)
	require.True(t, ok)
	assert.Equal(t, 256, len(tag))
}

func TestDesignProbe_ImplicitMin_SendsEmpty(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindStringImplicitMin, FieldPath: "requestBody.tag"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	tag, ok := body["tag"].(string)
	require.True(t, ok)
	assert.Equal(t, "", tag)
}

func TestDesignProbe_NullValue(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindNullValue, FieldPath: "requestBody.name"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	assert.Nil(t, body["name"])
	assert.Equal(t, 400, probe.ExpectedStatus)
}

func TestDesignProbe_NumericMin_BelowMin(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindNumericMin, FieldPath: "requestBody.age"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	age, ok := body["age"].(float64)
	require.True(t, ok)
	assert.Less(t, age, float64(1))
}

func TestDesignProbe_NumericMax_AboveMax(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindNumericMax, FieldPath: "requestBody.age"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	age, ok := body["age"].(float64)
	require.True(t, ok)
	assert.Greater(t, age, float64(100))
}

func TestDesignProbe_SetsContentTypeHeader(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindRequiredField, FieldPath: "requestBody.name"}
	probe := DesignProbe(h, op, newTestGen())
	assert.Equal(t, "application/json", probe.Headers["Content-Type"])
}

func TestDesignProbe_SetsOperationMethodAndPath(t *testing.T) {
	op := testOpForDesigner()
	h := &HypothesisNode{Kind: KindRequiredField, FieldPath: "requestBody.name"}
	probe := DesignProbe(h, op, newTestGen())
	assert.Equal(t, "POST", probe.Method)
	assert.Equal(t, "/pets", probe.Path)
}

func TestDesignProbe_EnumViolation(t *testing.T) {
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
		Responses: map[string]*spec.Response{"201": {Description: "Created"}},
	}
	h := &HypothesisNode{Kind: KindEnumViolation, FieldPath: "requestBody.status"}
	probe := DesignProbe(h, op, newTestGen())

	body := probe.Body.(map[string]any)
	status, ok := body["status"].(string)
	require.True(t, ok)
	assert.True(t, strings.HasPrefix(status, "__INVALID__"))
}
