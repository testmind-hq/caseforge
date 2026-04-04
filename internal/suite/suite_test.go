// internal/suite/suite_test.go
package suite_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/suite"
)

func newWriter() *writer.JSONSchemaWriter { return writer.NewJSONSchemaWriter() }

func writerOpts(suites []schema.TestSuite) writer.WriteOptions {
	return writer.WriteOptions{Suites: suites}
}

func validSuite() *schema.TestSuite {
	return &schema.TestSuite{
		ID:    "SUITE-001",
		Title: "E2E User Flow",
		Kind:  "chain",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-CREATE", Exports: []string{"user_id"}},
			{CaseID: "TC-READ", DependsOn: []string{"TC-CREATE"}},
			{CaseID: "TC-DELETE", DependsOn: []string{"TC-READ"}},
		},
	}
}

func knownCases() []schema.TestCase {
	return []schema.TestCase{
		{ID: "TC-CREATE"},
		{ID: "TC-READ"},
		{ID: "TC-DELETE"},
	}
}

// ── Validate ────────────────────────────────────────────────────────────────

func TestValidate_ValidSuite_NoErrors(t *testing.T) {
	errs := suite.Validate(validSuite(), knownCases())
	assert.Empty(t, errs)
}

func TestValidate_EmptyID_ReturnsError(t *testing.T) {
	s := validSuite()
	s.ID = ""
	errs := suite.Validate(s, nil)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "id")
}

func TestValidate_EmptyTitle_ReturnsError(t *testing.T) {
	s := validSuite()
	s.Title = ""
	errs := suite.Validate(s, nil)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "title")
}

func TestValidate_InvalidKind_ReturnsError(t *testing.T) {
	s := validSuite()
	s.Kind = "parallel"
	errs := suite.Validate(s, nil)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "kind")
}

func TestValidate_DuplicateCaseID_ReturnsError(t *testing.T) {
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-A"},
			{CaseID: "TC-A"}, // duplicate
		},
	}
	errs := suite.Validate(s, nil)
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Error(), "duplicate")
}

func TestValidate_UnknownDependsOn_ReturnsError(t *testing.T) {
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-A", DependsOn: []string{"TC-MISSING"}},
		},
	}
	errs := suite.Validate(s, nil)
	require.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Error(), "TC-MISSING")
}

func TestValidate_CaseNotInIndex_ReturnsError(t *testing.T) {
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{{CaseID: "TC-UNKNOWN"}},
	}
	// TC-UNKNOWN is not in knownCases
	errs := suite.Validate(s, knownCases())
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "TC-UNKNOWN")
}

func TestValidate_NilKnownCases_SkipsIndexCheck(t *testing.T) {
	// Without knownCases, we don't check whether case_ids exist in index.json.
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{{CaseID: "TC-IMAGINARY"}},
	}
	errs := suite.Validate(s, nil)
	assert.Empty(t, errs)
}

func TestValidate_CyclicDependency_ReturnsError(t *testing.T) {
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "chain",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-A", DependsOn: []string{"TC-B"}},
			{CaseID: "TC-B", DependsOn: []string{"TC-A"}}, // cycle
		},
	}
	errs := suite.Validate(s, nil)
	// Cycle is detected; the error about unknown deps may also fire
	hasCycleErr := false
	for _, e := range errs {
		if e.Field == "cases.depends_on" {
			hasCycleErr = true
		}
	}
	assert.True(t, hasCycleErr, "cycle must be detected")
}

// ── TopologicalOrder ────────────────────────────────────────────────────────

func TestTopologicalOrder_LinearChain(t *testing.T) {
	s := validSuite() // TC-CREATE → TC-READ → TC-DELETE
	sorted, err := suite.TopologicalOrder(s)
	require.NoError(t, err)
	require.Len(t, sorted, 3)
	assert.Equal(t, "TC-CREATE", sorted[0].CaseID)
	assert.Equal(t, "TC-READ", sorted[1].CaseID)
	assert.Equal(t, "TC-DELETE", sorted[2].CaseID)
}

func TestTopologicalOrder_DiamondGraph(t *testing.T) {
	// A → B, A → C, B → D, C → D
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "chain",
		Cases: []schema.SuiteCase{
			{CaseID: "A"},
			{CaseID: "B", DependsOn: []string{"A"}},
			{CaseID: "C", DependsOn: []string{"A"}},
			{CaseID: "D", DependsOn: []string{"B", "C"}},
		},
	}
	sorted, err := suite.TopologicalOrder(s)
	require.NoError(t, err)
	require.Len(t, sorted, 4)
	// A must come first, D must come last
	assert.Equal(t, "A", sorted[0].CaseID)
	assert.Equal(t, "D", sorted[3].CaseID)
}

func TestTopologicalOrder_Cycle_ReturnsError(t *testing.T) {
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "chain",
		Cases: []schema.SuiteCase{
			{CaseID: "A", DependsOn: []string{"B"}},
			{CaseID: "B", DependsOn: []string{"A"}},
		},
	}
	_, err := suite.TopologicalOrder(s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cycle")
}

func TestTopologicalOrder_NoDependencies_PreservesOrder(t *testing.T) {
	s := &schema.TestSuite{
		ID: "S", Title: "T", Kind: "sequential",
		Cases: []schema.SuiteCase{
			{CaseID: "TC-1"},
			{CaseID: "TC-2"},
			{CaseID: "TC-3"},
		},
	}
	sorted, err := suite.TopologicalOrder(s)
	require.NoError(t, err)
	require.Len(t, sorted, 3)
	assert.Equal(t, "TC-1", sorted[0].CaseID)
	assert.Equal(t, "TC-2", sorted[1].CaseID)
	assert.Equal(t, "TC-3", sorted[2].CaseID)
}

// ── File I/O ────────────────────────────────────────────────────────────────

func TestWriteAndLoadSuiteFile_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suite.json")

	s := validSuite()
	require.NoError(t, suite.WriteSuiteFile(s, path))

	loaded, err := suite.LoadSuiteFile(path)
	require.NoError(t, err)
	assert.Equal(t, s.ID, loaded.ID)
	assert.Equal(t, s.Title, loaded.Title)
	assert.Equal(t, s.Kind, loaded.Kind)
	assert.Len(t, loaded.Cases, 3)
	assert.Equal(t, schema.SuiteSchemaURL, loaded.Schema)
}

func TestLoadSuiteFile_MissingFile_ReturnsError(t *testing.T) {
	_, err := suite.LoadSuiteFile("/nonexistent/suite.json")
	require.Error(t, err)
}

func TestWriteSuiteFile_InvalidDir_ReturnsError(t *testing.T) {
	s := validSuite()
	err := suite.WriteSuiteFile(s, "/dev/null/cannot-write.json")
	require.Error(t, err)
}

// ── index.json suites field ──────────────────────────────────────────────────

func TestWriter_WriteAndReadFull_IncludesSuites(t *testing.T) {
	dir := t.TempDir()

	cases := []schema.TestCase{{ID: "TC-1", Title: "test"}}
	suites := []schema.TestSuite{*validSuite()}

	from_writer := newWriter()
	err := from_writer.Write(cases, dir, writerOpts(suites))
	require.NoError(t, err)

	_, statErr := os.Stat(filepath.Join(dir, "index.json"))
	require.NoError(t, statErr)

	full, err := from_writer.ReadFull(filepath.Join(dir, "index.json"))
	require.NoError(t, err)
	require.Len(t, full.Suites, 1)
	assert.Equal(t, "SUITE-001", full.Suites[0].ID)
}

func TestWriter_WriteNoSuites_OmitsSuitesKey(t *testing.T) {
	dir := t.TempDir()
	cases := []schema.TestCase{{ID: "TC-1", Title: "test"}}

	w := newWriter()
	require.NoError(t, w.Write(cases, dir, writerOpts(nil)))

	data, err := os.ReadFile(filepath.Join(dir, "index.json"))
	require.NoError(t, err)
	assert.NotContains(t, string(data), `"suites"`)
}
