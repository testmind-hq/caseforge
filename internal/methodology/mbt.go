// internal/methodology/mbt.go
package methodology

import (
	"fmt"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// classification represents one dimension of test conditions (e.g. one parameter)
// where each element is a distinct testable value (leaf in the tree).
type classification struct {
	name   string // parameter or field name
	leaves []any  // ordered list of distinct test values
}

// ClassificationTreeTechnique implements MBT Classification Tree testing.
//
// The technique builds a classification tree where each classification (internal node)
// corresponds to one discrete parameter or body field, and each leaf is a distinct
// testable value for that parameter. It then generates test cases using Each-Choice
// Coverage (ECT): every leaf appears in at least one test case while keeping other
// parameters at their default (first) value. The number of test cases equals the
// maximum leaf count across all classifications.
//
// Distinct from PairwiseTechnique (which covers all value *pairs*) and from
// EquivalenceTechnique (which focuses on missing-field cases): ClassificationTree
// ensures systematic per-parameter variation across all discrete choices.
//
// Applies when at least one parameter has 2+ discrete values (enum or boolean).
type ClassificationTreeTechnique struct {
	gen *datagen.Generator
}

func NewClassificationTreeTechnique() *ClassificationTreeTechnique {
	return &ClassificationTreeTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *ClassificationTreeTechnique) Name() string { return "classification_tree" }

// Applies returns true when at least one parameter has 2 or more discrete values.
func (t *ClassificationTreeTechnique) Applies(op *spec.Operation) bool {
	return len(buildClassifications(op)) >= 1
}

func (t *ClassificationTreeTechnique) Generate(op *spec.Operation) ([]schema.TestCase, error) {
	classes := buildClassifications(op)
	if len(classes) == 0 {
		return nil, nil
	}

	// Each-Choice Coverage (ECT): number of test cases = max leaves across classifications.
	maxLeaves := 0
	for _, c := range classes {
		if len(c.leaves) > maxLeaves {
			maxLeaves = len(c.leaves)
		}
	}

	var cases []schema.TestCase
	for row := 0; row < maxLeaves; row++ {
		// Build query params: for each classification, pick leaf[row % len(leaves)].
		queryParams := make(map[string]any, len(classes))
		leafLabels := make([]string, 0, len(classes))
		for _, c := range classes {
			val := c.leaves[row%len(c.leaves)]
			queryParams[c.name] = val
			leafLabels = append(leafLabels, fmt.Sprintf("%s=%v", c.name, val))
		}

		base := buildValidBody(t.gen, op)
		tc := buildTestCase(op, base,
			fmt.Sprintf("classification tree row %d: %v", row+1, leafLabels),
			"")
		tc.Priority = "P2"
		tc.Source = schema.CaseSource{
			Technique: "classification_tree",
			SpecPath:  fmt.Sprintf("%s %s parameters", op.Method, op.Path),
			Rationale: fmt.Sprintf("ECT row %d — each-choice coverage: %v", row+1, leafLabels),
		}
		tc.Steps[0].Path = buildPathWithQuery(op.Path, queryParams)
		cases = append(cases, tc)
	}
	return cases, nil
}

// buildClassifications extracts one classification per discrete parameter (enum or boolean).
func buildClassifications(op *spec.Operation) []classification {
	var out []classification
	for _, p := range op.Parameters {
		if p.Schema == nil {
			continue
		}
		if len(p.Schema.Enum) >= 2 {
			out = append(out, classification{name: p.Name, leaves: p.Schema.Enum})
		} else if p.Schema.Type == "boolean" {
			out = append(out, classification{name: p.Name, leaves: []any{true, false}})
		}
	}
	return out
}
