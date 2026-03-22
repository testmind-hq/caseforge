// cmd/diff_test.go
//go:build integration

package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffCommand_BasicText(t *testing.T) {
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "text"
	diffCases = ""
	t.Cleanup(func() {
		diffOld = ""
		diffNew = ""
		diffFormat = "text"
		diffCases = ""
	})
	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)
	// breaking changes → runDiff returns errBreakingChanges; that's expected
	err := runDiff(diffCmd, nil)
	require.ErrorIs(t, err, errBreakingChanges)
	output := buf.String()
	assert.Contains(t, output, "BREAKING")
	assert.Contains(t, output, "/pets/{petId}")
}

func TestDiffCommand_JSONFormat(t *testing.T) {
	diffOld = "../testdata/petstore_v1.yaml"
	diffNew = "../testdata/petstore_v2.yaml"
	diffFormat = "json"
	diffCases = ""
	t.Cleanup(func() {
		diffOld = ""
		diffNew = ""
		diffFormat = "text"
		diffCases = ""
	})
	buf := &bytes.Buffer{}
	diffCmd.SetOut(buf)
	diffCmd.SetErr(buf)
	err := runDiff(diffCmd, nil)
	require.ErrorIs(t, err, errBreakingChanges)
	output := buf.String()
	assert.Contains(t, output, `"kind"`)
	assert.Contains(t, output, "BREAKING")
}
