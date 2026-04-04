// internal/rbt/generator_test.go
package rbt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Uses makeSpec/makeOp helpers from assessor_test.go (same package).

func TestHighRiskOperations_ReturnsOnlyHighRisk(t *testing.T) {
	sp := makeSpec(
		makeOp("GET", "/pets", "listPets"),
		makeOp("POST", "/pets", "createPet"),
		makeOp("DELETE", "/pets/{id}", "deletePet"),
	)
	report := RiskReport{Operations: []OperationCoverage{
		{Method: "GET", Path: "/pets", Risk: RiskHigh},
		{Method: "POST", Path: "/pets", Risk: RiskLow},
		{Method: "DELETE", Path: "/pets/{id}", Risk: RiskNone},
	}}
	ops := HighRiskOperations(report, sp)
	require.Len(t, ops, 1)
	assert.Equal(t, "GET", ops[0].Method)
	assert.Equal(t, "/pets", ops[0].Path)
}

func TestHighRiskOperations_NoneHighRisk_ReturnsNil(t *testing.T) {
	sp := makeSpec(makeOp("GET", "/pets", "listPets"))
	report := RiskReport{Operations: []OperationCoverage{
		{Method: "GET", Path: "/pets", Risk: RiskLow},
	}}
	assert.Nil(t, HighRiskOperations(report, sp))
}

func TestHighRiskOperations_MultipleHighRisk(t *testing.T) {
	sp := makeSpec(
		makeOp("GET", "/pets", "listPets"),
		makeOp("POST", "/pets", "createPet"),
		makeOp("DELETE", "/pets/{id}", "deletePet"),
	)
	report := RiskReport{Operations: []OperationCoverage{
		{Method: "GET", Path: "/pets", Risk: RiskHigh},
		{Method: "POST", Path: "/pets", Risk: RiskHigh},
		{Method: "DELETE", Path: "/pets/{id}", Risk: RiskMedium},
	}}
	ops := HighRiskOperations(report, sp)
	require.Len(t, ops, 2)
}

func TestHighRiskOperations_PreservesSpecOrder(t *testing.T) {
	sp := makeSpec(
		makeOp("GET", "/a", "opA"),
		makeOp("POST", "/b", "opB"),
		makeOp("DELETE", "/c", "opC"),
	)
	// Report lists /c before /a, but result should follow spec order.
	report := RiskReport{Operations: []OperationCoverage{
		{Method: "DELETE", Path: "/c", Risk: RiskHigh},
		{Method: "GET", Path: "/a", Risk: RiskHigh},
		{Method: "POST", Path: "/b", Risk: RiskLow},
	}}
	ops := HighRiskOperations(report, sp)
	require.Len(t, ops, 2)
	assert.Equal(t, "/a", ops[0].Path, "spec order: /a before /c")
	assert.Equal(t, "/c", ops[1].Path)
}

func TestHighRiskOperations_CaseInsensitiveMethod(t *testing.T) {
	// Spec uses lowercase "get", report uses uppercase "GET".
	sp := makeSpec(makeOp("get", "/pets", "listPets"))
	report := RiskReport{Operations: []OperationCoverage{
		{Method: "GET", Path: "/pets", Risk: RiskHigh},
	}}
	ops := HighRiskOperations(report, sp)
	require.Len(t, ops, 1)
}

func TestHighRiskOperations_EmptyReport_ReturnsNil(t *testing.T) {
	sp := makeSpec(makeOp("GET", "/pets", "listPets"))
	assert.Nil(t, HighRiskOperations(RiskReport{}, sp))
}
