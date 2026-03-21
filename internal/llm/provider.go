// internal/llm/provider.go
package llm

import "context"

type LLMProvider interface {
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	IsAvailable() bool
	Name() string
}

type CompletionRequest struct {
	System    string
	Messages  []Message
	MaxTokens int
}

type Message struct {
	Role    string // "user"|"assistant"
	Content string
}

type CompletionResponse struct {
	Text   string
	Tokens int
}
