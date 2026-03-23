// internal/mcp/server.go
package mcp

import (
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewServer creates the CaseForge MCP server.
// Actual per-call sampling is done in makeGenerateHandler via req.Session.
func NewServer() *mcpsdk.Server {
	s := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "caseforge", Version: "1.0.0"}, nil)

	s.AddTool(generateTestCasesTool(), makeGenerateHandler)

	return s
}
