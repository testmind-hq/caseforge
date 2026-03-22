// internal/llm/anthropic.go
package llm

import (
	"context"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	client *anthropic.Client
	model  string
}

func (p *AnthropicProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	msgs := make([]anthropic.MessageParam, len(req.Messages))
	for i, m := range req.Messages {
		if m.Role == "user" {
			msgs[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content))
		} else {
			msgs[i] = anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content))
		}
	}

	system := []anthropic.TextBlockParam{}
	if req.System != "" {
		system = []anthropic.TextBlockParam{{Text: req.System}}
	}

	resp, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: int64(req.MaxTokens),
		System:    system,
		Messages:  msgs,
	})
	if err != nil {
		return nil, fmt.Errorf("anthropic complete: %w", err)
	}
	text := ""
	if len(resp.Content) > 0 {
		text = resp.Content[0].Text
	}
	return &CompletionResponse{
		Text:   text,
		Tokens: int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
	}, nil
}

func (p *AnthropicProvider) IsAvailable() bool { return true }
func (p *AnthropicProvider) Name() string      { return "anthropic:" + p.model }

func newAnthropicClient(apiKey string) *anthropic.Client {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &c
}
