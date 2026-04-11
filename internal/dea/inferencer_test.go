// internal/dea/inferencer_test.go
package dea

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInferRule_RequiredField_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		ID:          "H-001",
		Kind:        KindRequiredField,
		Operation:   "POST /pets",
		FieldPath:   "requestBody.name",
		Description: "omitting required field 'name' must return 4xx",
		Status:      StatusConfirmed,
		Evidence:    &Evidence{ActualStatus: 400, ActualBody: `{"error":"name required"}`},
	}

	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.False(t, rule.Implicit, "required field is declared in spec — not implicit")
	assert.Equal(t, ConfidenceHigh, rule.Confidence)
	assert.Contains(t, rule.Description, "name")
	assert.Len(t, rule.Evidence, 1)
	assert.Equal(t, 400, rule.Evidence[0].ActualStatus)
}

func TestInferRule_RequiredField_Refuted_SpecMismatch(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindRequiredField,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 201},
	}

	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.True(t, rule.Implicit)
	assert.Contains(t, rule.Description, "spec declares required but server accepts omission")
}

func TestInferRule_ImplicitMax_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindStringImplicitMax,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400, ActualBody: `{"error":"too long"}`},
	}

	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit, "implicit max is not in spec")
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.Contains(t, rule.Description, "undeclared max length")
}

func TestInferRule_ImplicitMax_Refuted_NoRule(t *testing.T) {
	h := &HypothesisNode{
		Kind:     KindStringImplicitMax,
		Status:   StatusRefuted,
		Evidence: &Evidence{ActualStatus: 201},
	}

	rule := InferRule(h)
	assert.Nil(t, rule, "refuted implicit max = no undeclared constraint")
}

func TestInferRule_SpecMax_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindStringMaxLength,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}

	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit, "spec-declared max confirmed by server")
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_SpecMax_Refuted_Mismatch(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindStringMaxLength,
		FieldPath: "requestBody.name",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 201},
	}

	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.Contains(t, rule.Description, "spec declares max but server does not enforce")
}

func TestInferRule_PendingHypothesis_Panics(t *testing.T) {
	h := &HypothesisNode{Kind: KindRequiredField, Status: StatusPending}
	assert.Panics(t, func() { InferRule(h) }, "must panic on pending hypothesis")
}

func TestInferRule_EnumViolation_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindEnumViolation,
		Operation: "POST /orders",
		FieldPath: "requestBody.status",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 422},
	}

	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_OptionalField_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindOptionalField,
		Operation: "POST /pets",
		FieldPath: "requestBody.tag",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 201},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.False(t, rule.Implicit)
}

func TestInferRule_OptionalField_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindOptionalField,
		FieldPath: "requestBody.tag",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.True(t, rule.Implicit)
}

func TestInferRule_SpecMin_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindStringMinLength,
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_ImplicitMin_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindStringImplicitMin,
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
}

func TestInferRule_ImplicitMin_Refuted_NoRule(t *testing.T) {
	h := &HypothesisNode{
		Kind:     KindStringImplicitMin,
		Status:   StatusRefuted,
		Evidence: &Evidence{ActualStatus: 200},
	}
	rule := InferRule(h)
	assert.Nil(t, rule)
}

func TestInferRule_NumericMin_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindNumericMin,
		FieldPath: "requestBody.age",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_NumericMax_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindNumericMax,
		FieldPath: "requestBody.age",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
}

func TestInferRule_NullValue_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindNullValue,
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_NullValue_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindNullValue,
		FieldPath: "requestBody.name",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 200},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategoryBehavior, rule.Category)
}

func TestInferRule_NilEvidence_Panics(t *testing.T) {
	h := &HypothesisNode{Kind: KindRequiredField, Status: StatusConfirmed, Evidence: nil}
	assert.Panics(t, func() { InferRule(h) })
}

func TestInferRule_ArrayMinItems_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindArrayMinItems,
		FieldPath: "requestBody.tags",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 422},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.Contains(t, rule.Description, "tags")
}

func TestInferRule_ArrayMinItems_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindArrayMinItems,
		FieldPath: "requestBody.tags",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 200},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.Contains(t, rule.Description, "minItems")
}

func TestInferRule_ArrayMaxItems_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindArrayMaxItems,
		FieldPath: "requestBody.tags",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_ArrayMaxItems_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindArrayMaxItems,
		FieldPath: "requestBody.tags",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 201},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.Contains(t, rule.Description, "maxItems")
}

func TestInferRule_RequiredQueryParam_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindRequiredQueryParam,
		Operation: "GET /pets",
		FieldPath: "query.status",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.Contains(t, rule.Description, "status")
}

func TestInferRule_RequiredQueryParam_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindRequiredQueryParam,
		FieldPath: "query.status",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 200},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.Contains(t, rule.Description, "spec declares required")
}

func TestInferRule_FormatViolation_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindFormatViolation,
		FieldPath: "requestBody.email",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 422},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.Contains(t, rule.Description, "email")
}

func TestInferRule_FormatViolation_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindFormatViolation,
		FieldPath: "requestBody.email",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 201},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.Contains(t, rule.Description, "format")
}

func TestInferRule_TypeCoercion_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindTypeCoercion,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 422},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.False(t, rule.Implicit)
	assert.Contains(t, rule.Description, "name")
}

func TestInferRule_TypeCoercion_Refuted(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindTypeCoercion,
		FieldPath: "requestBody.age",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 200},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryBehavior, rule.Category)
	assert.True(t, rule.Implicit)
	assert.Contains(t, rule.Description, "coercion")
}

func TestInferRule_UnicodeControl_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindUnicodeControl,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.False(t, rule.Implicit)
	assert.Contains(t, rule.Description, "name")
}

func TestInferRule_MassAssignment_Confirmed(t *testing.T) {
	h := &HypothesisNode{
		Kind:      KindMassAssignment,
		Operation: "POST /users",
		FieldPath: "requestBody",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 200},
	}
	rule := InferRule(h)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryBehavior, rule.Category)
	assert.True(t, rule.Implicit)
	assert.Contains(t, rule.Description, "Mass assignment")
}
