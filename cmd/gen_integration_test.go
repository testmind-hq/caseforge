// cmd/gen_integration_test.go
//go:build integration

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenCommandWithPetstoreSpec(t *testing.T) {
	outDir := t.TempDir()
	genSpec = "../testdata/petstore.yaml"
	genOutput = outDir
	genNoAI = true
	genFormat = "hurl"
	t.Cleanup(func() {
		genSpec = ""
		genOutput = "./cases"
		genNoAI = false
		genFormat = "hurl"
	})

	err := runGen(genCmd, nil)
	require.NoError(t, err)

	// index.json must exist
	indexPath := filepath.Join(outDir, "index.json")
	assert.FileExists(t, indexPath)

	// At least one .hurl file
	hurls, _ := filepath.Glob(filepath.Join(outDir, "*.hurl"))
	assert.NotEmpty(t, hurls)
}
