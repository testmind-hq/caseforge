// internal/mcp/tools.go
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
)

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

func makeGenerateHandler() mcpsdk.ToolHandler {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		// Parse arguments from raw JSON
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

		// Create session-scoped LLM provider that samples via MCP sampling protocol
		provider := llm.NewMCPSamplingProvider(func(ctx context.Context, cr *llm.CompletionRequest) (*llm.CompletionResponse, error) {
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
			text := ""
			if tc, ok := result.Content.(*mcpsdk.TextContent); ok {
				text = tc.Text
			}
			return &llm.CompletionResponse{Text: text}, nil
		})

		// Load spec
		loader := spec.NewLoader()
		parsedSpec, err := loader.Load(args.Spec)
		if err != nil {
			r := &mcpsdk.CallToolResult{}
			r.SetError(fmt.Errorf("loading spec: %w", err))
			return r, nil
		}

		// Build engine with all techniques (mirrors cmd/gen.go)
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

		// Create output directory
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			r := &mcpsdk.CallToolResult{}
			r.SetError(fmt.Errorf("creating output dir: %w", err))
			return r, nil
		}

		// Write index.json
		w := writer.NewJSONSchemaWriter()
		if err := w.Write(cases, outDir); err != nil {
			r := &mcpsdk.CallToolResult{}
			r.SetError(fmt.Errorf("writing index: %w", err))
			return r, nil
		}

		// Render to target format
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
		default:
			renderer = render.NewHurlRenderer("")
		}
		if err := renderer.Render(cases, outDir); err != nil {
			r := &mcpsdk.CallToolResult{}
			r.SetError(fmt.Errorf("rendering: %w", err))
			return r, nil
		}

		summary := fmt.Sprintf("Generated %d test cases → %s (%s format)", len(cases), outDir, format)
		return &mcpsdk.CallToolResult{Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: summary}}}, nil
	}
}
