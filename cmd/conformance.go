// cmd/conformance.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/oracle"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/runner"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var conformanceCmd = &cobra.Command{
	Use:   "conformance",
	Short: "Check spec-vs-implementation conformance using LLM-mined response body constraints",
	Long: `Conformance mines response body constraints from OAS descriptions using
LLM Observation-Confirmation prompting, generates assertion tests, runs them
against a live API, and reports spec-vs-implementation mismatches.

Requires an LLM provider configured in .caseforge.yaml and a reachable target URL.`,
	RunE:         runConformance,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(conformanceCmd)
	conformanceCmd.Flags().String("spec", "", "OpenAPI spec file (required)")
	conformanceCmd.Flags().String("target", "", "API base URL to test against (required)")
	conformanceCmd.Flags().String("output", "", "Directory to write conformance report JSON (optional)")
	_ = conformanceCmd.MarkFlagRequired("spec")
	_ = conformanceCmd.MarkFlagRequired("target")
}

// ConformanceIssue represents a single spec-vs-implementation mismatch.
type ConformanceIssue struct {
	Operation  string `json:"operation"`
	Constraint string `json:"constraint"`
	Detail     string `json:"detail"`
}

// ConformanceReport is the top-level report written to conformance-report.json.
type ConformanceReport struct {
	Target      string             `json:"target"`
	SpecPath    string             `json:"spec_path"`
	TestedOps   int                `json:"tested_ops"`
	FailedOps   int                `json:"failed_ops"`
	Issues      []ConformanceIssue `json:"issues"`
	GeneratedAt string             `json:"generated_at"`
}

func runConformance(cmd *cobra.Command, _ []string) error {
	specPath, _ := cmd.Flags().GetString("spec")
	target, _ := cmd.Flags().GetString("target")
	outputDir, _ := cmd.Flags().GetString("output")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	provider := llm.NewProviderWithConfig(llm.ProviderConfig{
		APIKey:   cfg.AI.APIKey,
		Provider: cfg.AI.Provider,
		Model:    cfg.AI.Model,
		BaseURL:  cfg.AI.BaseURL,
	})
	if !provider.IsAvailable() {
		return fmt.Errorf("LLM provider not available — conformance requires an LLM (set provider in .caseforge.yaml or use env var)")
	}

	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(specPath)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Mining response constraints for %d operation(s)...\n", len(parsedSpec.Operations))

	var cases []schema.TestCase
	ctx := context.Background()
	for _, op := range parsedSpec.Operations {
		constraints, mineErr := oracle.Mine(ctx, op, provider)
		if mineErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "  warn: oracle mining failed for %s %s: %v\n", op.Method, op.Path, mineErr)
			continue
		}
		if len(constraints) == 0 {
			continue
		}
		cases = append(cases, buildConformanceCase(op, constraints))
	}

	if len(cases) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No testable constraints found. Ensure OAS descriptions are detailed.")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Generated %d conformance test case(s). Running against %s...\n", len(cases), target)

	tmpDir, err := os.MkdirTemp("", "caseforge-conformance-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write cases to index.json
	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, tmpDir, writer.WriteOptions{}); err != nil {
		return fmt.Errorf("writing test cases: %w", err)
	}

	// Render to hurl files so the runner can find them
	hurlRenderer := render.NewHurlRenderer("")
	if err := hurlRenderer.Render(cases, tmpDir); err != nil {
		return fmt.Errorf("rendering hurl files: %w", err)
	}

	r := runner.NewHurlRunner()
	result, runErr := r.Run(tmpDir, map[string]string{"BASE_URL": target})
	if runErr != nil {
		return fmt.Errorf("running conformance tests: %w", runErr)
	}

	caseByID := make(map[string]schema.TestCase)
	for _, tc := range cases {
		caseByID[tc.ID] = tc
	}
	var issues []ConformanceIssue
	failedOps := 0
	for _, cr := range result.Cases {
		if cr.Passed {
			continue
		}
		failedOps++
		if tc, ok := caseByID[cr.ID]; ok {
			issues = append(issues, ConformanceIssue{
				Operation:  tc.Source.SpecPath,
				Constraint: tc.Source.Rationale,
				Detail:     "assertion failed — response body did not match expected constraints",
			})
		}
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "\nConformance Results: %d/%d operations passed\n", result.Passed, result.Passed+result.Failed)
	if len(issues) > 0 {
		fmt.Fprintln(out, "\nSpec-vs-Implementation Mismatches:")
		for _, iss := range issues {
			fmt.Fprintf(out, "  ✗ %s\n    %s\n", iss.Operation, iss.Constraint)
		}
	} else {
		fmt.Fprintln(out, "No mismatches found.")
	}

	if outputDir != "" {
		report := ConformanceReport{
			Target:      target,
			SpecPath:    specPath,
			TestedOps:   len(cases),
			FailedOps:   failedOps,
			Issues:      issues,
			GeneratedAt: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		}
		if err := writeConformanceReport(outputDir, report); err != nil {
			return fmt.Errorf("writing report: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "Report written to: %s\n", filepath.Join(outputDir, "conformance-report.json"))
	}

	return nil
}

func buildConformanceCase(op *spec.Operation, constraints []oracle.Constraint) schema.TestCase {
	specPath := fmt.Sprintf("%s %s", op.Method, op.Path)
	assertions := []schema.Assertion{
		{Target: "status_code", Operator: schema.OperatorGte, Expected: 200},
		{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
	}
	var constraintDescs []string
	for _, c := range constraints {
		assertions = append(assertions, c.ToAssertions()...)
		constraintDescs = append(constraintDescs, c.Detail)
	}

	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
		Title:    fmt.Sprintf("[conformance] %s", specPath),
		Kind:     "single",
		Priority: "P0",
		Steps: []schema.Step{{
			ID:         "step-1",
			Title:      fmt.Sprintf("check %s response body constraints", specPath),
			Type:       "test",
			Method:     op.Method,
			Path:       op.Path,
			Assertions: assertions,
		}},
		Source: schema.CaseSource{
			Technique: "conformance",
			SpecPath:  specPath,
			Rationale: fmt.Sprintf("OC-mined constraints: %v", constraintDescs),
		},
		GeneratedAt: time.Now(),
	}
}

func writeConformanceReport(outputDir string, report ConformanceReport) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputDir, "conformance-report.json"), data, 0644)
}
