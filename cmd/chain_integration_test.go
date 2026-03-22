//go:build integration

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

func TestEndToEndChainCase(t *testing.T) {
	outDir := t.TempDir()
	t.Cleanup(func() {
		genSpec = ""
		genOutput = "./cases"
		genNoAI = false
		genFormat = "hurl"
	})

	genSpec = "../testdata/crud_users.yaml"
	genOutput = outDir
	genNoAI = true
	genFormat = "hurl"

	err := runGen(genCmd, nil)
	require.NoError(t, err)

	// Read index.json and find chain cases
	w := writer.NewJSONSchemaWriter()
	cases, err := w.Read(filepath.Join(outDir, "index.json"))
	require.NoError(t, err)

	var chainCases []schema.TestCase
	for _, c := range cases {
		if c.Kind == "chain" {
			chainCases = append(chainCases, c)
		}
	}
	require.NotEmpty(t, chainCases, "expected at least one chain case from crud_users.yaml")

	// Find the rendered .hurl file for the chain case
	hurlFile := filepath.Join(outDir, chainCases[0].ID+".hurl")
	content, err := os.ReadFile(hurlFile)
	require.NoError(t, err)
	contentStr := string(content)

	assert.Contains(t, contentStr, "[Captures]", "chain hurl file must have [Captures] block")
	assert.Contains(t, contentStr, "userId:", "chain hurl file must capture userId")
	assert.Contains(t, contentStr, "{{userId}}", "read step must use captured variable")
	// Verify [Captures] appears before [Asserts]
	capturesIdx := strings.Index(contentStr, "[Captures]")
	assertsIdx := strings.Index(contentStr, "[Asserts]")
	if assertsIdx >= 0 && capturesIdx >= 0 {
		assert.Less(t, capturesIdx, assertsIdx, "[Captures] must appear before [Asserts]")
	}
}

func TestEndToEndPhase2CSVWithChainCases(t *testing.T) {
	outDir := t.TempDir()
	t.Cleanup(func() {
		genSpec = ""
		genOutput = "./cases"
		genNoAI = false
		genFormat = "hurl"
	})

	genSpec = "../testdata/crud_users.yaml"
	genOutput = outDir
	genNoAI = true
	genFormat = "csv"

	err := runGen(genCmd, nil)
	require.NoError(t, err)

	// Verify CSV exists and has chain case
	data, err := os.ReadFile(filepath.Join(outDir, "cases.csv"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "chain")
	assert.Contains(t, string(data), "chain_crud")
}
