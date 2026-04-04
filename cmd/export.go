// cmd/export.go
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/export"
	"github.com/testmind-hq/caseforge/internal/output/writer"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export index.json to a third-party platform format (allure, xray, testrail)",
	Long: `Export reads index.json from the cases directory and converts test cases
into the format required by the chosen platform:

  allure    — one {uuid}-result.json per test case (Allure Report)
  xray      — single xray-import.json (Jira Xray Cloud)
  testrail  — single testrail-import.csv (TestRail CSV import)

Examples:
  caseforge export --cases ./cases --format allure
  caseforge export --cases ./cases --format xray --output ./exports
  caseforge export --cases ./cases --format testrail --output ./exports`,
	RunE: runExport,
}

var (
	exportCases  string
	exportFormat string
	exportOutput string
)

func init() {
	exportCmd.Flags().StringVar(&exportCases, "cases", "", "Directory containing index.json (required)")
	exportCmd.Flags().StringVar(&exportFormat, "format", "", "Export format: allure|xray|testrail (required)")
	exportCmd.Flags().StringVar(&exportOutput, "output", "./export", "Output directory")
	_ = exportCmd.MarkFlagRequired("cases")
	_ = exportCmd.MarkFlagRequired("format")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, _ []string) error {
	indexPath := filepath.Join(exportCases, "index.json")

	w := &writer.JSONSchemaWriter{}
	cases, err := w.Read(indexPath)
	if err != nil {
		return fmt.Errorf("reading index.json: %w", err)
	}

	exp, err := export.New(exportFormat)
	if err != nil {
		return err
	}

	outDir := filepath.Join(exportOutput, exportFormat)
	if err := exp.Export(cases, outDir); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Exported %d test cases to %s\n", len(cases), outDir)
	return nil
}
