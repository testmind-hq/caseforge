//go:build integration

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestEndToEndPetstore(t *testing.T) {
	outDir := t.TempDir()
	noop := &llm.NoopProvider{}

	loader := spec.NewLoader()
	ps, err := loader.Load("testdata/petstore.yaml")
	require.NoError(t, err)
	assert.NotEmpty(t, ps.Operations)

	engine := methodology.NewEngine(noop,
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
	)
	cases, err := engine.Generate(ps)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)

	// Validate traceability: every case must have full source info
	for _, tc := range cases {
		assert.NotEmpty(t, tc.ID, "case must have ID")
		assert.Equal(t, "1", tc.Version, "case must have version=1")
		assert.NotEmpty(t, tc.Source.Technique, "case must have technique")
		assert.NotEmpty(t, tc.Source.SpecPath, "case must have spec_path")
		assert.NotEmpty(t, tc.Source.Rationale, "case must have rationale")
	}

	// Write index.json
	w := writer.NewJSONSchemaWriter()
	require.NoError(t, w.Write(cases, outDir))
	assert.FileExists(t, filepath.Join(outDir, "index.json"))

	// Validate index.json is parseable
	data, err := os.ReadFile(filepath.Join(outDir, "index.json"))
	require.NoError(t, err)
	var index writer.IndexFile
	require.NoError(t, json.Unmarshal(data, &index))
	assert.Equal(t, len(cases), len(index.TestCases))

	// Render Hurl files
	renderer := render.NewHurlRenderer("{{base_url}}")
	require.NoError(t, renderer.Render(cases, outDir))
	hurlFiles, _ := filepath.Glob(filepath.Join(outDir, "*.hurl"))
	assert.NotEmpty(t, hurlFiles, "should produce at least one .hurl file")
}

func TestEndToEndComplexParams(t *testing.T) {
	_ = t.TempDir()
	noop := &llm.NoopProvider{}

	loader := spec.NewLoader()
	ps, err := loader.Load("testdata/complex_params.yaml")
	require.NoError(t, err)

	engine := methodology.NewEngine(noop,
		methodology.NewEquivalenceTechnique(),
		methodology.NewPairwiseTechnique(),
	)
	cases, err := engine.Generate(ps)
	require.NoError(t, err)
	assert.NotEmpty(t, cases)

	// Pairwise should have been applied (5 parameters with enums/booleans)
	pairwiseCases := 0
	for _, tc := range cases {
		if tc.Source.Technique == "pairwise" {
			pairwiseCases++
		}
	}
	assert.Greater(t, pairwiseCases, 0, "expected pairwise cases for complex_params spec")
}
