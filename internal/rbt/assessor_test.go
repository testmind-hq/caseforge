// internal/rbt/assessor_test.go
package rbt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func writeCaseFile(t *testing.T, dir string, tc schema.TestCase) {
	t.Helper()
	data, err := json.Marshal(tc)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, tc.ID+".json"), data, 0644))
}

func makeCase(id, specPath string) schema.TestCase {
	return schema.TestCase{
		ID:          id,
		Title:       "test " + id,
		Source:      schema.CaseSource{SpecPath: specPath},
		GeneratedAt: time.Now(),
	}
}

func TestScanCases_Empty(t *testing.T) {
	dir := t.TempDir()
	idx, err := ScanCases(dir)
	require.NoError(t, err)
	assert.Empty(t, idx)
}

func TestScanCases_ParsesSpecPath(t *testing.T) {
	dir := t.TempDir()
	writeCaseFile(t, dir, makeCase("c1", "POST /users"))
	writeCaseFile(t, dir, makeCase("c2", "GET /users/{id} requestBody.properties.name"))
	writeCaseFile(t, dir, makeCase("c3", "POST /users"))

	idx, err := ScanCases(dir)
	require.NoError(t, err)
	assert.Len(t, idx["POST /users"], 2, "two cases for POST /users")
	assert.Len(t, idx["GET /users/{id}"], 1, "one case for GET /users/{id}")
}

func TestScanCases_IgnoresNonJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("hello"), 0644))
	idx, err := ScanCases(dir)
	require.NoError(t, err)
	assert.Empty(t, idx)
}

func makeSpec(ops ...*spec.Operation) *spec.ParsedSpec {
	return &spec.ParsedSpec{Operations: ops}
}

func makeOp(method, path, id string) *spec.Operation {
	return &spec.Operation{Method: method, Path: path, OperationID: id}
}

func TestAssess_NoAffectedOps(t *testing.T) {
	sp := makeSpec(makeOp("GET", "/users", "listUsers"))
	report := Assess(sp, nil, nil, "HEAD~1", "HEAD", nil)
	assert.Equal(t, 0, report.TotalAffected)
	assert.InDelta(t, 0.0, report.RiskScore, 0.001)
	assert.Len(t, report.Operations, 1)
	assert.Equal(t, RiskNone, report.Operations[0].Risk)
}

func TestAssess_HighRisk_NoTestCases(t *testing.T) {
	sp := makeSpec(makeOp("POST", "/users", "createUser"))
	affected := []RouteMapping{{Method: "POST", RoutePath: "/users"}}
	report := Assess(sp, affected, nil, "HEAD~1", "HEAD", nil)
	assert.Equal(t, 1, report.TotalAffected)
	assert.Equal(t, 1, report.TotalUncovered)
	assert.Equal(t, RiskHigh, report.Operations[0].Risk)
	assert.InDelta(t, 1.0, report.RiskScore, 0.001)
}

func TestAssess_MediumRisk_OneTestCase(t *testing.T) {
	sp := makeSpec(makeOp("POST", "/users", "createUser"))
	affected := []RouteMapping{{Method: "POST", RoutePath: "/users"}}
	idx := map[string][]TestCaseRef{
		"POST /users": {{File: "c1.json", CaseID: "c1"}},
	}
	report := Assess(sp, affected, idx, "HEAD~1", "HEAD", nil)
	assert.Equal(t, RiskMedium, report.Operations[0].Risk)
	assert.Equal(t, 0, report.TotalUncovered)
	assert.Equal(t, 1, report.TotalCovered)
}

func TestAssess_LowRisk_TwoOrMoreTestCases(t *testing.T) {
	sp := makeSpec(makeOp("GET", "/pets", "listPets"))
	affected := []RouteMapping{{Method: "GET", RoutePath: "/pets"}}
	idx := map[string][]TestCaseRef{
		"GET /pets": {
			{File: "c1.json", CaseID: "c1"},
			{File: "c2.json", CaseID: "c2"},
		},
	}
	report := Assess(sp, affected, idx, "HEAD~1", "HEAD", nil)
	assert.Equal(t, RiskLow, report.Operations[0].Risk)
}

func TestAssess_UncertainRisk_LowConfidence(t *testing.T) {
	sp := makeSpec(makeOp("DELETE", "/users/{id}", "deleteUser"))
	affected := []RouteMapping{{Method: "DELETE", RoutePath: "/users/{id}", Confidence: 0.3}}
	report := Assess(sp, affected, nil, "HEAD~1", "HEAD", nil)
	assert.Equal(t, RiskUncertain, report.Operations[0].Risk)
}
