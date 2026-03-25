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
	RunE:  runExport,
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
