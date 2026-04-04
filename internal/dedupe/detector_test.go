// internal/dedupe/detector_test.go
package dedupe

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func writeCase(t *testing.T, dir string, tc schema.TestCase) string {
	t.Helper()
	data, err := json.Marshal(tc)
	require.NoError(t, err)
	fp := filepath.Join(dir, tc.ID+".json")
	require.NoError(t, os.WriteFile(fp, data, 0644))
	return fp
}

func makeStep(method, path string, status int, bodyKey string, extras ...string) schema.Step {
	var body any
	if bodyKey != "" {
		body = map[string]any{bodyKey: "val"}
	}
	assertions := []schema.Assertion{
		{Target: "status_code", Operator: "eq", Expected: float64(status)},
	}
	for _, tgt := range extras {
		assertions = append(assertions, schema.Assertion{Target: tgt, Operator: "eq", Expected: "x"})
	}
	return schema.Step{
		ID: "s1", Type: "test",
		Method: method, Path: path,
		Body: body, Assertions: assertions,
	}
}

func makeTC(id string, steps ...schema.Step) schema.TestCase {
	return schema.TestCase{ID: id, GeneratedAt: time.Now(), Steps: steps}
}

func TestFindDuplicates_NilInput(t *testing.T) {
	groups, err := FindDuplicates(nil, 0.9)
	require.NoError(t, err)
	assert.Empty(t, groups)
}

func TestFindDuplicates_ExactDuplicate_ReturnsOneGroup(t *testing.T) {
	dir := t.TempDir()
	step := makeStep("POST", "/users", 201, "name", "jsonpath $.id")
	pathA := writeCase(t, dir, makeTC("case-a", step))
	pathB := writeCase(t, dir, makeTC("case-b", step))

	cases, err := ScanCases(dir)
	require.NoError(t, err)

	groups, err := FindDuplicates(cases, 0.9)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	assert.Equal(t, MatchExact, groups[0].Kind)
	assert.InDelta(t, 1.0, groups[0].Similarity, 0.001)

	paths := []string{groups[0].Cases[0].FilePath, groups[0].Cases[1].FilePath}
	assert.Contains(t, paths, pathA)
	assert.Contains(t, paths, pathB)
}

func TestFindDuplicates_StructuralDuplicate_AboveThreshold(t *testing.T) {
	dir := t.TempDir()
	stepA := makeStep("GET", "/users", 200, "", "jsonpath $.id", "jsonpath $.name", "jsonpath $.email")
	stepB := makeStep("GET", "/users", 200, "", "jsonpath $.id", "jsonpath $.name")

	writeCase(t, dir, makeTC("case-a", stepA))
	writeCase(t, dir, makeTC("case-b", stepB))

	cases, err := ScanCases(dir)
	require.NoError(t, err)

	groups, err := FindDuplicates(cases, 0.5)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	assert.Equal(t, MatchStructural, groups[0].Kind)
	assert.GreaterOrEqual(t, groups[0].Similarity, 0.5)
}

func TestFindDuplicates_StructuralDuplicate_BelowThreshold_NoGroup(t *testing.T) {
	dir := t.TempDir()
	stepA := makeStep("GET", "/users", 200, "", "jsonpath $.id", "jsonpath $.name", "jsonpath $.email")
	stepB := makeStep("GET", "/users", 200, "", "jsonpath $.id", "jsonpath $.name")

	writeCase(t, dir, makeTC("case-a", stepA))
	writeCase(t, dir, makeTC("case-b", stepB))

	cases, err := ScanCases(dir)
	require.NoError(t, err)

	groups, err := FindDuplicates(cases, 0.95)
	require.NoError(t, err)
	assert.Empty(t, groups)
}

func TestFindDuplicates_UniqueFiles_NoGroups(t *testing.T) {
	dir := t.TempDir()
	writeCase(t, dir, makeTC("a", makeStep("GET", "/users", 200, "")))
	writeCase(t, dir, makeTC("b", makeStep("POST", "/users", 201, "name")))

	cases, err := ScanCases(dir)
	require.NoError(t, err)

	groups, err := FindDuplicates(cases, 0.9)
	require.NoError(t, err)
	assert.Empty(t, groups)
}

func TestFindDuplicates_Keep_MostAssertions(t *testing.T) {
	dir := t.TempDir()
	stepRich := makeStep("POST", "/users", 201, "name", "jsonpath $.id", "jsonpath $.name", "header X-Req-ID")
	stepPoor := makeStep("POST", "/users", 201, "name", "jsonpath $.id")

	writeCase(t, dir, makeTC("aaa-case", stepRich))
	writeCase(t, dir, makeTC("bbb-case", stepPoor))

	cases, err := ScanCases(dir)
	require.NoError(t, err)

	groups, err := FindDuplicates(cases, 0.9)
	require.NoError(t, err)
	require.Len(t, groups, 1)

	var kept *CaseScore
	for i := range groups[0].Cases {
		if groups[0].Cases[i].Keep {
			kept = &groups[0].Cases[i]
		}
	}
	require.NotNil(t, kept)
	assert.Equal(t, 4, kept.AssertionCount)
}

func TestFindDuplicates_TieBreak_LexicographicFilename(t *testing.T) {
	dir := t.TempDir()
	step := makeStep("GET", "/pets", 200, "", "jsonpath $.id")
	writeCase(t, dir, makeTC("aaa-case", step))
	writeCase(t, dir, makeTC("zzz-case", step))

	cases, err := ScanCases(dir)
	require.NoError(t, err)

	groups, err := FindDuplicates(cases, 0.9)
	require.NoError(t, err)
	require.Len(t, groups, 1)

	var kept *CaseScore
	for i := range groups[0].Cases {
		if groups[0].Cases[i].Keep {
			kept = &groups[0].Cases[i]
		}
	}
	require.NotNil(t, kept)
	assert.Contains(t, kept.FilePath, "aaa-case")
}

func TestJaccardSimilarity(t *testing.T) {
	assert.InDelta(t, 0.5, jaccardSimilarity(
		[]string{"a", "b", "c"},
		[]string{"b", "c", "d"},
	), 0.001)
	assert.InDelta(t, 1.0, jaccardSimilarity([]string{}, []string{}), 0.001)
	assert.InDelta(t, 1.0, jaccardSimilarity([]string{"x"}, []string{"x"}), 0.001)
	assert.InDelta(t, 0.0, jaccardSimilarity([]string{"a"}, []string{"b"}), 0.001)
}
