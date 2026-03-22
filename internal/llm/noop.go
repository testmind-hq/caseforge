// internal/llm/noop.go
package llm

import "context"

type NoopProvider struct{}

func (p *NoopProvider) Complete(_ context.Context, _ *CompletionRequest) (*CompletionResponse, error) {
	return &CompletionResponse{Text: "", Tokens: 0}, nil
}
func (p *NoopProvider) IsAvailable() bool { return false }
func (p *NoopProvider) Name() string      { return "noop" }
