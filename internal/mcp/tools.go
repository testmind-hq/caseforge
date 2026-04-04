// internal/mcp/tools.go
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/testmind-hq/caseforge/internal/ask"
	"github.com/testmind-hq/caseforge/internal/lint"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// --- generate_test_cases ---

func generateTestCasesTool() *mcpsdk.Tool {
	return &mcpsdk.Tool{
		Name:        "generate_test_cases",
		Description: "Generate API test cases from an OpenAPI spec file or URL and write them to an output directory.",
		InputSchema: json.RawMessage(`{
            "type": "object",
            "properties": {
                "spec": {
                    "type": "string",
                    "description": "Path or URL to the OpenAPI spec (e.g. ./api.yaml or https://example.com/openapi.json)"
                },
                "output": {
                    "type": "string",
                    "description": "Output directory for generated test files (default: ./cases)"
                },
                "format": {
                    "type": "string",
                    "description": "Output format: hurl|markdown|csv|postman|k6 (default: hurl)",
                    "enum": ["hurl", "markdown", "csv", "postman", "k6"]
                }
            },
            "required": ["spec"]
        }`),
	}
}

func makeGenerateHandler(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	var args struct {
		Spec   string `json:"spec"`
		Output string `json:"output"`
		Format string `json:"format"`
	}
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("invalid arguments: %w", err))
		return r, nil
	}

	outDir := args.Output
	if outDir == "" {
		outDir = "./cases"
	}
	format := args.Format
	if format == "" {
		format = "hurl"
	}

	provider := makeSamplingProvider(ctx, req)

	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(args.Spec)
	if err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("loading spec: %w", err))
		return r, nil
	}

	engine := methodology.NewEngine(provider,
		methodology.NewEquivalenceTechnique(),
		methodology.NewBoundaryTechnique(),
		methodology.NewDecisionTechnique(),
		methodology.NewStateTechnique(),
		methodology.NewIdempotentTechnique(),
		methodology.NewPairwiseTechnique(),
		methodology.NewSecurityTechnique(),
	)
	engine.AddSpecTechnique(methodology.NewChainTechnique())
	engine.AddSpecTechnique(methodology.NewSecuritySpecTechnique())

	cases, err := engine.Generate(parsedSpec)
	if err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("generating cases: %w", err))
		return r, nil
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("creating output dir: %w", err))
		return r, nil
	}

	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, outDir, writer.WriteOptions{}); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("writing index: %w", err))
		return r, nil
	}

	renderer := rendererFor(format)
	if err := renderer.Render(cases, outDir); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("rendering: %w", err))
		return r, nil
	}

	summary := fmt.Sprintf("Generated %d test cases → %s (%s format)", len(cases), outDir, format)
	return &mcpsdk.CallToolResult{Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: summary}}}, nil
}

// --- lint_spec ---

func lintSpecTool() *mcpsdk.Tool {
	return &mcpsdk.Tool{
		Name:        "lint_spec",
		Description: "Lint an OpenAPI spec and return quality issues with a score (0–100). Use fail_on to filter by severity.",
		InputSchema: json.RawMessage(`{
            "type": "object",
            "properties": {
                "spec": {
                    "type": "string",
                    "description": "Path or URL to the OpenAPI spec"
                },
                "fail_on": {
                    "type": "string",
                    "description": "Minimum severity to include in results: error|warning (default: warning)",
                    "enum": ["error", "warning"]
                }
            },
            "required": ["spec"]
        }`),
	}
}

func makeLintHandler(_ context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	var args struct {
		Spec   string `json:"spec"`
		FailOn string `json:"fail_on"`
	}
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("invalid arguments: %w", err))
		return r, nil
	}
	if args.FailOn == "" {
		args.FailOn = "warning"
	}

	loader := spec.NewLoader()
	ps, err := loader.Load(args.Spec)
	if err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("loading spec: %w", err))
		return r, nil
	}

	issues := lint.RunAll(ps, nil)

	// Filter by fail_on severity: "error" shows only errors; "warning" (default) shows all.
	var filtered []lint.LintIssue
	for _, iss := range issues {
		if args.FailOn == "error" && iss.Severity != "error" {
			continue
		}
		filtered = append(filtered, iss)
	}

	// Build text report
	var sb strings.Builder
	errors, warnings := 0, 0
	for _, iss := range filtered {
		if iss.Severity == "error" {
			errors++
		} else {
			warnings++
		}
		sb.WriteString(fmt.Sprintf("[%s] %s: %s", strings.ToUpper(iss.Severity), iss.RuleID, iss.Message))
		if iss.Path != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", iss.Path))
		}
		sb.WriteString("\n")
	}

	score := lint.Score(issues)
	summary := fmt.Sprintf("Score: %d/100 — %d error(s), %d warning(s)\n\n%s", score, errors, warnings, sb.String())
	if len(filtered) == 0 {
		summary = fmt.Sprintf("Score: %d/100 — no issues found.", score)
	}
	return &mcpsdk.CallToolResult{Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: summary}}}, nil
}

// --- ask_test_cases ---

func askTestCasesTool() *mcpsdk.Tool {
	return &mcpsdk.Tool{
		Name:        "ask_test_cases",
		Description: "Generate API test cases from a natural language description using the host LLM via MCP sampling.",
		InputSchema: json.RawMessage(`{
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "description": "Natural language description of the API operation and test scenarios (e.g. 'POST /users - test valid email, missing fields, duplicate email')"
                },
                "output": {
                    "type": "string",
                    "description": "Output directory for generated test files (default: ./cases)"
                },
                "format": {
                    "type": "string",
                    "description": "Output format: hurl|markdown|csv|postman|k6 (default: hurl)",
                    "enum": ["hurl", "markdown", "csv", "postman", "k6"]
                }
            },
            "required": ["description"]
        }`),
	}
}

func makeAskHandler(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	var args struct {
		Description string `json:"description"`
		Output      string `json:"output"`
		Format      string `json:"format"`
	}
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("invalid arguments: %w", err))
		return r, nil
	}
	if args.Description == "" {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("description is required"))
		return r, nil
	}

	outDir := args.Output
	if outDir == "" {
		outDir = "./cases"
	}
	format := args.Format
	if format == "" {
		format = "hurl"
	}

	provider := makeSamplingProvider(ctx, req)

	gen := ask.NewGenerator(provider)
	cases, err := gen.Generate(ctx, args.Description)
	if err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("generating cases: %w", err))
		return r, nil
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("creating output dir: %w", err))
		return r, nil
	}

	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, outDir, writer.WriteOptions{}); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("writing index: %w", err))
		return r, nil
	}

	renderer := rendererFor(format)
	if err := renderer.Render(cases, outDir); err != nil {
		r := &mcpsdk.CallToolResult{}
		r.SetError(fmt.Errorf("rendering: %w", err))
		return r, nil
	}

	summary := fmt.Sprintf("Generated %d test cases from description → %s (%s format)", len(cases), outDir, format)
	return &mcpsdk.CallToolResult{Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: summary}}}, nil
}

// --- shared helpers ---

// makeSamplingProvider creates an MCP sampling provider that delegates LLM calls
// to the host (Claude Desktop / Claude Code) via the MCP sampling protocol.
func makeSamplingProvider(_ context.Context, req *mcpsdk.CallToolRequest) llm.LLMProvider {
	return llm.NewMCPSamplingProvider(func(ctx context.Context, cr *llm.CompletionRequest) (*llm.CompletionResponse, error) {
		if req.Session == nil {
			return nil, fmt.Errorf("mcp sampling: no active session (sampling capability not negotiated)")
		}
		msgs := make([]*mcpsdk.SamplingMessage, 0, len(cr.Messages))
		for _, m := range cr.Messages {
			role := mcpsdk.Role("user")
			if m.Role == "assistant" {
				role = mcpsdk.Role("assistant")
			}
			msgs = append(msgs, &mcpsdk.SamplingMessage{
				Role:    role,
				Content: &mcpsdk.TextContent{Text: m.Content},
			})
		}
		result, err := req.Session.CreateMessage(ctx, &mcpsdk.CreateMessageParams{
			Messages:     msgs,
			SystemPrompt: cr.System,
			MaxTokens:    int64(cr.MaxTokens),
		})
		if err != nil {
			return nil, fmt.Errorf("mcp sampling: %w", err)
		}
		tc, ok := result.Content.(*mcpsdk.TextContent)
		if !ok {
			return nil, fmt.Errorf("mcp sampling: unexpected content type %T from host", result.Content)
		}
		return &llm.CompletionResponse{Text: tc.Text}, nil
	})
}

// rendererFor returns the renderer for the given format name.
func rendererFor(format string) render.Renderer {
	switch format {
	case "markdown":
		return render.NewMarkdownRenderer()
	case "csv":
		return render.NewCSVRenderer()
	case "postman":
		return render.NewPostmanRenderer()
	case "k6":
		return render.NewK6Renderer()
	default:
		return render.NewHurlRenderer("")
	}
}
