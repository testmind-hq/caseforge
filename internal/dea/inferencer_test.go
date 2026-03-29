// internal/dea/inferencer_test.go
package dea

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestInferRule_RequiredField_Confirmed(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{
		ID:          "H-001",
		Kind:        KindRequiredField,
		Operation:   "POST /pets",
		FieldPath:   "requestBody.name",
		Description: "omitting required field 'name' must return 4xx",
		Status:      StatusConfirmed,
		Evidence:    &Evidence{ActualStatus: 400, ActualBody: `{"error":"name required"}`},
	}

	rule := InferRule(h, op)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.False(t, rule.Implicit, "required field is declared in spec — not implicit")
	assert.Equal(t, ConfidenceHigh, rule.Confidence)
	assert.Contains(t, rule.Description, "name")
	assert.Len(t, rule.Evidence, 1)
	assert.Equal(t, 400, rule.Evidence[0].ActualStatus)
}

func TestInferRule_RequiredField_Refuted_SpecMismatch(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{
		Kind:      KindRequiredField,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 201},
	}

	rule := InferRule(h, op)
	require.NotNil(t, rule)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.True(t, rule.Implicit)
	assert.Contains(t, rule.Description, "spec declares required but server accepts omission")
}

func TestInferRule_ImplicitMax_Confirmed(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{
		Kind:      KindStringImplicitMax,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400, ActualBody: `{"error":"too long"}`},
	}

	rule := InferRule(h, op)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit, "implicit max is not in spec")
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
	assert.Contains(t, rule.Description, "undeclared max length")
}

func TestInferRule_ImplicitMax_Refuted_NoRule(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{
		Kind:     KindStringImplicitMax,
		Status:   StatusRefuted,
		Evidence: &Evidence{ActualStatus: 201},
	}

	rule := InferRule(h, op)
	assert.Nil(t, rule, "refuted implicit max = no undeclared constraint")
}

func TestInferRule_SpecMax_Confirmed(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{
		Kind:      KindStringMaxLength,
		Operation: "POST /pets",
		FieldPath: "requestBody.name",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 400},
	}

	rule := InferRule(h, op)
	require.NotNil(t, rule)
	assert.False(t, rule.Implicit, "spec-declared max confirmed by server")
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}

func TestInferRule_SpecMax_Refuted_Mismatch(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{
		Kind:      KindStringMaxLength,
		FieldPath: "requestBody.name",
		Status:    StatusRefuted,
		Evidence:  &Evidence{ActualStatus: 201},
	}

	rule := InferRule(h, op)
	require.NotNil(t, rule)
	assert.True(t, rule.Implicit)
	assert.Equal(t, CategorySpecMismatch, rule.Category)
	assert.Contains(t, rule.Description, "spec declares max but server does not enforce")
}

func TestInferRule_PendingHypothesis_Panics(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/pets"}
	h := &HypothesisNode{Kind: KindRequiredField, Status: StatusPending}
	assert.Panics(t, func() { InferRule(h, op) }, "must panic on pending hypothesis")
}

func TestInferRule_EnumViolation_Confirmed(t *testing.T) {
	op := &spec.Operation{Method: "POST", Path: "/orders"}
	h := &HypothesisNode{
		Kind:      KindEnumViolation,
		Operation: "POST /orders",
		FieldPath: "requestBody.status",
		Status:    StatusConfirmed,
		Evidence:  &Evidence{ActualStatus: 422},
	}

	rule := InferRule(h, op)
	require.NotNil(t, rule)
	assert.Equal(t, CategoryFieldConstraint, rule.Category)
}
