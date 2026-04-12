// cmd/import.go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/har"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import test cases from external sources",
}

var importHARCmd = &cobra.Command{
	Use:          "har <file>",
	Short:        "Import test cases from an HTTP Archive (HAR) file",
	Args:         cobra.ExactArgs(1),
	RunE:         runImportHAR,
	SilenceUsage: true,
}

func init() {
	importCmd.AddCommand(importHARCmd)
	rootCmd.AddCommand(importCmd)
	importHARCmd.Flags().String("output", "./cases", "Output directory for generated test cases")
	importHARCmd.Flags().String("format", "hurl", "Output format: hurl|postman|k6|markdown|csv")
}

func runImportHAR(cmd *cobra.Command, args []string) error {
	harFile := args[0]
	outputDir, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")

	data, err := os.ReadFile(harFile)
	if err != nil {
		return fmt.Errorf("reading HAR file: %w", err)
	}

	entries, err := har.Parse(data)
	if err != nil {
		return fmt.Errorf("parsing HAR file: %w", err)
	}

	// Deduplicate: keep first occurrence per METHOD+path (ignoring query)
	seen := map[string]bool{}
	var cases []schema.TestCase
	for i, e := range entries {
		path := e.Request.URL // already stripped to /path?query by parser
		key := e.Request.Method + " " + stripQuery(path)
		if seen[key] {
			continue
		}
		seen[key] = true

		tc := harEntryToTestCase(i, e)
		cases = append(cases, tc)
	}

	// Select renderer based on format flag
	var renderer render.Renderer
	switch format {
	case "markdown":
		renderer = render.NewMarkdownRenderer()
	case "csv":
		renderer = render.NewCSVRenderer()
	case "postman":
		renderer = render.NewPostmanRenderer()
	case "k6":
		renderer = render.NewK6Renderer()
	default: // "hurl" and anything unrecognised
		renderer = render.NewHurlRenderer("{{base_url}}")
	}

	return renderer.Render(cases, outputDir)
}

func harEntryToTestCase(idx int, e har.Entry) schema.TestCase {
	id := fmt.Sprintf("TC-%04d", idx+1)
	path := e.Request.URL
	title := fmt.Sprintf("[har_replay] %s %s", e.Request.Method, stripQuery(path))

	var body any
	if e.Request.Body != "" && strings.Contains(e.Request.MIMEType, "json") {
		_ = json.Unmarshal([]byte(e.Request.Body), &body)
		// if unmarshal fails, body stays nil (skip non-JSON bodies)
	}

	headers := map[string]string{}
	for _, h := range e.Request.Headers {
		headers[h.Name] = h.Value
	}

	step := schema.Step{
		ID:     "step-main",
		Title:  title,
		Type:   "test",
		Method: e.Request.Method,
		Path:   path,
		Body:   body,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: "eq", Expected: e.Response.Status},
		},
	}
	if len(headers) > 0 {
		step.Headers = headers
	}

	return schema.TestCase{
		Schema:      schema.SchemaBaseURL,
		Version:     "1",
		ID:          id,
		Title:       title,
		Kind:        "single",
		Priority:    "P2",
		Steps:       []schema.Step{step},
		Source: schema.CaseSource{
			Technique: "har_replay",
			SpecPath:  fmt.Sprintf("%s %s", e.Request.Method, stripQuery(path)),
			Rationale: "replayed from HAR traffic recording",
		},
		GeneratedAt: time.Now(),
	}
}

func stripQuery(path string) string {
	if i := strings.IndexByte(path, '?'); i >= 0 {
		return path[:i]
	}
	return path
}
