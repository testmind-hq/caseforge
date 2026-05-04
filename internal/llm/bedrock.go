// internal/llm/bedrock.go
package llm

import (
	"context"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

type BedrockProvider struct {
	client *anthropic.Client
	model  string
}

func newBedrockProvider(ctx context.Context, model, region string) (*BedrockProvider, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("bedrock: load AWS config: %w", err)
	}
	c := anthropic.NewClient(bedrock.WithConfig(cfg))
	return &BedrockProvider{client: &c, model: model}, nil
}

func (p *BedrockProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
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
		return nil, fmt.Errorf("bedrock complete: %w", err)
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

func (p *BedrockProvider) IsAvailable() bool { return true }
func (p *BedrockProvider) Name() string      { return "bedrock:" + p.model }
