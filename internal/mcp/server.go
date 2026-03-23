// internal/mcp/server.go
package mcp

import (
	"context"
	"fmt"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/testmind-hq/caseforge/internal/llm"
)

// NewServer creates the CaseForge MCP server and a sentinel LLMProvider.
// The sentinel's IsAvailable() returns true and Name() returns "mcp-sampling".
// Actual per-call sampling is done in makeGenerateHandler via req.Session.
func NewServer() (*mcpsdk.Server, llm.LLMProvider) {
	s := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "caseforge", Version: "1.0.0"}, nil)

	sentinel := llm.NewMCPSamplingProvider(func(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
		return nil, fmt.Errorf("mcp sampling: no active session")
	})

	s.AddTool(generateTestCasesTool(), makeGenerateHandler())

	return s, sentinel
}
