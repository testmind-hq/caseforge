// cmd/mutate_test.go
package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutateCmdRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "mutate" {
			found = true
			break
		}
	}
	assert.True(t, found, "mutate command must be registered on rootCmd")
}

func TestMutateCmdRequiresFlags(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"mutate"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("mutate without --cases and --target must return error")
	}
	rootCmd.SetArgs(nil)
}
