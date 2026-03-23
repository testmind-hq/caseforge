// internal/llm/openai.go
package llm

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIProvider struct {
	client  *openai.Client
	model   string
	baseURL string
}

type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

func NewOpenAIProvider(cfg OpenAIConfig) *OpenAIProvider {
	opts := []option.RequestOption{option.WithAPIKey(cfg.APIKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}
	client := openai.NewClient(opts...)
	return &OpenAIProvider{client: &client, model: cfg.Model, baseURL: cfg.BaseURL}
}

func (p *OpenAIProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	msgs := []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(req.System)}
	for _, m := range req.Messages {
		if m.Role == "user" {
			msgs = append(msgs, openai.UserMessage(m.Content))
		} else {
			msgs = append(msgs, openai.AssistantMessage(m.Content))
		}
	}
	resp, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(p.model),
		Messages:  msgs,
		MaxTokens: openai.Int(int64(req.MaxTokens)),
	})
	if err != nil {
		return nil, fmt.Errorf("openai complete: %w", err)
	}
	text := ""
	if len(resp.Choices) > 0 {
		text = resp.Choices[0].Message.Content
	}
	return &CompletionResponse{Text: text, Tokens: int(resp.Usage.TotalTokens)}, nil
}

func (p *OpenAIProvider) IsAvailable() bool { return true }
func (p *OpenAIProvider) Name() string {
	if p.baseURL != "" {
		return "openai-compat:" + p.model
	}
	return "openai:" + p.model
}
