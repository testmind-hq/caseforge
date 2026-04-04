// internal/mcp/server_test.go
package mcp

import (
	"context"
	"encoding/json"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// callReq is a helper that builds a CallToolRequest with the given JSON arguments.
func callReq(t *testing.T, argsJSON string) *mcpsdk.CallToolRequest {
	t.Helper()
	return &mcpsdk.CallToolRequest{
		Params: &mcpsdk.CallToolParamsRaw{
			Arguments: json.RawMessage(argsJSON),
		},
	}
}

func TestNewServerReturnsNonNil(t *testing.T) {
	s := NewServer()
	require.NotNil(t, s)
}

func TestServerHasGenerateTestCasesTool(t *testing.T) {
	assert.NotNil(t, generateTestCasesTool())
}

func TestServerHasLintSpecTool(t *testing.T) {
	tool := lintSpecTool()
	require.NotNil(t, tool)
	assert.Equal(t, "lint_spec", tool.Name)
	assert.NotEmpty(t, tool.Description)

	rawSchema, _ := tool.InputSchema.(json.RawMessage)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(rawSchema, &schema))
	required, _ := schema["required"].([]any)
	assert.Contains(t, required, "spec")
}

func TestServerHasAskTestCasesTool(t *testing.T) {
	tool := askTestCasesTool()
	require.NotNil(t, tool)
	assert.Equal(t, "ask_test_cases", tool.Name)
	assert.NotEmpty(t, tool.Description)

	rawSchema, _ := tool.InputSchema.(json.RawMessage)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(rawSchema, &schema))
	required, _ := schema["required"].([]any)
	assert.Contains(t, required, "description")
}

func TestLintHandler_InvalidJSON_ReturnsToolError(t *testing.T) {
	result, err := makeLintHandler(context.Background(), callReq(t, `{bad json`))
	require.NoError(t, err, "handler must not return a Go error")
	require.NotNil(t, result)
	assert.True(t, result.IsError, "invalid JSON args must produce a tool error")
}

func TestLintHandler_MissingSpec_ReturnsToolError(t *testing.T) {
	result, err := makeLintHandler(context.Background(), callReq(t, `{"spec":"/nonexistent/api.yaml"}`))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsError, "non-existent spec must produce a tool error")
}

func TestAskHandler_InvalidJSON_ReturnsToolError(t *testing.T) {
	result, err := makeAskHandler(context.Background(), callReq(t, `{bad json`))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestAskHandler_EmptyDescription_ReturnsToolError(t *testing.T) {
	result, err := makeAskHandler(context.Background(), callReq(t, `{"description":""}`))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsError, "empty description must produce a tool error")
}

func TestAskHandler_NoSession_ReturnsToolError(t *testing.T) {
	// No session → MCP sampling unavailable → ask.Generate fails with provider error.
	result, err := makeAskHandler(context.Background(), callReq(t, `{"description":"POST /users create user"}`))
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsError, "missing session must produce a tool error")
}
