// cmd/import_test.go
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestImportCommandIsRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "import" {
			return
		}
	}
	t.Error("import command not registered in rootCmd")
}

func TestImportHARCommandIsRegistered(t *testing.T) {
	var imp *cobra.Command
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "import" {
			imp = sub
			break
		}
	}
	if imp == nil {
		t.Fatal("import command not found")
	}
	for _, sub := range imp.Commands() {
		if sub.Use == "har <file>" {
			return
		}
	}
	t.Error("import har subcommand not registered")
}

func TestImportHARCommand_WritesTestCases(t *testing.T) {
	outDir := t.TempDir()
	rootCmd.SetArgs([]string{"import", "har", "testdata/sample.har", "--output", outDir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("reading output dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected output files in output directory, got none")
	}
}

func TestImportHARCommand_DeduplicatesEntries(t *testing.T) {
	outDir := t.TempDir()
	rootCmd.SetArgs([]string{"import", "har", "testdata/sample.har", "--output", outDir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// sample.har has 3 entries but POST /users appears twice — expect 2 output files
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("reading output dir: %v", err)
	}

	// Count non-directory files
	var fileCount int
	for _, e := range entries {
		if !e.IsDir() {
			fileCount++
		}
	}

	if fileCount != 2 {
		t.Errorf("expected 2 test case files (deduplication of POST /users), got %d", fileCount)
		for _, e := range entries {
			t.Logf("  file: %s", filepath.Join(outDir, e.Name()))
		}
	}
}
