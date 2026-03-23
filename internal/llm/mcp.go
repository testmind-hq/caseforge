// internal/llm/mcp.go
package llm

import "context"

// SamplerFunc is a function that fulfills an LLM completion request
// by forwarding it to an external system (e.g. MCP host).
type SamplerFunc func(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

// MCPSamplingProvider implements LLMProvider by delegating to a SamplerFunc.
// The sampler is typically a closure that calls back to a connected MCP client
// via the MCP sampling/createMessage protocol.
type MCPSamplingProvider struct {
	sampler SamplerFunc
}

// NewMCPSamplingProvider creates an LLMProvider backed by the given sampler function.
func NewMCPSamplingProvider(sampler SamplerFunc) *MCPSamplingProvider {
	return &MCPSamplingProvider{sampler: sampler}
}

func (p *MCPSamplingProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	return p.sampler(ctx, req)
}

func (p *MCPSamplingProvider) IsAvailable() bool { return true }
func (p *MCPSamplingProvider) Name() string      { return "mcp-sampling" }
