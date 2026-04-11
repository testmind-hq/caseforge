// cmd/gen_flags_test.go
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// --- splitTrimmed ---

func TestSplitTrimmed_BasicComma(t *testing.T) {
	got := splitTrimmed("equivalence_partitioning,boundary_value")
	assert.Equal(t, []string{"equivalence_partitioning", "boundary_value"}, got)
}

func TestSplitTrimmed_TrimsSpaces(t *testing.T) {
	got := splitTrimmed(" a , b , c ")
	assert.Equal(t, []string{"a", "b", "c"}, got)
}

func TestSplitTrimmed_SingleItem(t *testing.T) {
	assert.Equal(t, []string{"listPets"}, splitTrimmed("listPets"))
}

func TestSplitTrimmed_Empty(t *testing.T) {
	assert.Empty(t, splitTrimmed(""))
}

// --- filterTechniques ---

func TestFilterTechniques_EmptyFilterReturnsAll(t *testing.T) {
	ops := []methodology.Technique{
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
	}
	specs := []methodology.SpecTechnique{methodology.NewChainTechnique()}

	gotOps, gotSpecs := filterTechniques(ops, specs, "")
	assert.Len(t, gotOps, 2)
	assert.Len(t, gotSpecs, 1)
}

func TestFilterTechniques_ByName(t *testing.T) {
	ops := []methodology.Technique{
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
	}
	specs := []methodology.SpecTechnique{
		methodology.NewChainTechnique(),
		methodology.NewSecuritySpecTechnique(),
	}

	gotOps, gotSpecs := filterTechniques(ops, specs, "equivalence_partitioning,chain_crud")
	require.Len(t, gotOps, 1)
	assert.Equal(t, "equivalence_partitioning", gotOps[0].Name())
	require.Len(t, gotSpecs, 1)
	assert.Equal(t, "chain_crud", gotSpecs[0].Name())
}

func TestFilterTechniques_UnknownNameReturnsEmpty(t *testing.T) {
	ops := []methodology.Technique{methodology.NewEquivalenceTechnique()}
	specs := []methodology.SpecTechnique{methodology.NewChainTechnique()}

	gotOps, gotSpecs := filterTechniques(ops, specs, "nonexistent_technique")
	assert.Empty(t, gotOps)
	assert.Empty(t, gotSpecs)
}

// --- filterByPriority ---

func makeCase(id, priority string) schema.TestCase {
	return schema.TestCase{ID: id, Priority: priority}
}

func TestFilterByPriority_P0KeepsOnlyP0(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("a", "P0"),
		makeCase("b", "P1"),
		makeCase("c", "P2"),
	}
	got := filterByPriority(cases, "P0")
	require.Len(t, got, 1)
	assert.Equal(t, "a", got[0].ID)
}

func TestFilterByPriority_P2KeepsP0P1P2(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("a", "P0"),
		makeCase("b", "P1"),
		makeCase("c", "P2"),
		makeCase("d", "P3"),
	}
	got := filterByPriority(cases, "P2")
	require.Len(t, got, 3)
	ids := []string{got[0].ID, got[1].ID, got[2].ID}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, ids)
}

func TestFilterByPriority_P3KeepsAll(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("a", "P0"), makeCase("b", "P1"),
		makeCase("c", "P2"), makeCase("d", "P3"),
	}
	got := filterByPriority(cases, "P3")
	assert.Len(t, got, 4)
}

func TestFilterByPriority_UnknownPriorityExcluded(t *testing.T) {
	cases := []schema.TestCase{
		makeCase("a", "P0"),
		makeCase("b", ""),      // empty priority
		makeCase("c", "HIGH"),  // unrecognised
	}
	got := filterByPriority(cases, "P3")
	require.Len(t, got, 1)
	assert.Equal(t, "a", got[0].ID)
}

func TestGenCommand_HasFilterFlags(t *testing.T) {
	assert.NotNil(t, genCmd.Flags().Lookup("include-path"), "--include-path flag required")
	assert.NotNil(t, genCmd.Flags().Lookup("exclude-path"), "--exclude-path flag required")
	assert.NotNil(t, genCmd.Flags().Lookup("include-tag"),  "--include-tag flag required")
	assert.NotNil(t, genCmd.Flags().Lookup("exclude-tag"),  "--exclude-tag flag required")
}
