package mutation_test

import (
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func TestComputeOperationScores_Basic(t *testing.T) {
	results := []mutation.CaseMutationResult{
		{CaseID: "TC-0001", Operation: "GET /pets", Operator: "field_drop", Survived: false},
		{CaseID: "TC-0001", Operation: "GET /pets", Operator: "status_swap_2xx", Survived: false},
		{CaseID: "TC-0002", Operation: "POST /pets", Operator: "field_drop", Survived: true},
		{CaseID: "TC-0002", Operation: "POST /pets", Operator: "status_swap_2xx", Survived: false},
	}
	scores := mutation.ComputeOperationScores(results)

	if len(scores) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(scores))
	}
	// weakest first: POST /pets = 0.5, GET /pets = 1.0
	if scores[0].Operation != "POST /pets" {
		t.Errorf("expected weakest first, got %s", scores[0].Operation)
	}
	if scores[0].Killed != 1 || scores[0].Survivors != 1 || scores[0].TotalRuns != 2 {
		t.Errorf("POST /pets stats wrong: %+v", scores[0])
	}
	if scores[1].Operation != "GET /pets" {
		t.Errorf("expected GET /pets second, got %s", scores[1].Operation)
	}
	if scores[1].MutationScore != 1.0 {
		t.Errorf("GET /pets should be 100%%, got %f", scores[1].MutationScore)
	}
}

func TestComputeOperationScores_Empty(t *testing.T) {
	scores := mutation.ComputeOperationScores(nil)
	if len(scores) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(scores))
	}
}

func TestComputeOperationScores_UnknownFallback(t *testing.T) {
	results := []mutation.CaseMutationResult{
		{CaseID: "TC-0001", Operation: "", Operator: "field_drop", Survived: true},
	}
	scores := mutation.ComputeOperationScores(results)
	if len(scores) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(scores))
	}
	if scores[0].Operation != "(unknown)" {
		t.Errorf("expected '(unknown)', got %q", scores[0].Operation)
	}
}

func TestComputeOperationScores_TieBreak(t *testing.T) {
	// Equal scores must be sorted alphabetically by operation name
	results := []mutation.CaseMutationResult{
		{Operation: "POST /z", Operator: "field_drop", Survived: true},
		{Operation: "POST /a", Operator: "field_drop", Survived: true},
	}
	scores := mutation.ComputeOperationScores(results)
	if scores[0].Operation != "POST /a" {
		t.Errorf("expected alphabetical tiebreak, got %s first", scores[0].Operation)
	}
}
