// cmd/export_test.go
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportCmd_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "export" {
			found = true
			break
		}
	}
	assert.True(t, found, "export command must be registered on rootCmd")
}

func TestExportCmd_RequiredFlags(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "export" {
			assert.NotNil(t, c.Flags().Lookup("cases"), "--cases flag must exist")
			assert.NotNil(t, c.Flags().Lookup("format"), "--format flag must exist")
			assert.NotNil(t, c.Flags().Lookup("output"), "--output flag must exist")
			return
		}
	}
	t.Fatal("export command not found")
}
