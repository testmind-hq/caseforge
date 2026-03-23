// cmd/mcp.go
package cmd

import (
	"context"
	"io"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	caseforgemcp "github.com/testmind-hq/caseforge/internal/mcp"
)

var mcpCmd = &cobra.Command{
	Use:          "mcp",
	Short:        "Start CaseForge as an MCP server (stdio transport)",
	SilenceUsage: true,
	Long: `Starts CaseForge as a Model Context Protocol server over stdio.

Add to Claude Desktop's MCP config:
  {
    "mcpServers": {
      "caseforge": {
        "command": "caseforge",
        "args": ["mcp"]
      }
    }
  }

Claude Desktop will then be able to call generate_test_cases,
using its own model for LLM inference (no API key required).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := caseforgemcp.NewServer()
		err := s.Run(context.Background(), &mcpsdk.StdioTransport{})
		if err == nil || err == io.EOF {
			return nil
		}
		// "server is closing: EOF" — normal shutdown when the client closes stdin.
		if strings.Contains(err.Error(), "EOF") {
			return nil
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
