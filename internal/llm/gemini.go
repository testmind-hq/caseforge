// internal/llm/gemini.go
package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// GeminiProvider implements LLMProvider using Google's Gemini API (genai v1.51.0+).
type GeminiProvider struct {
	client *genai.Client
	model  string
}

// newGeminiProvider creates a GeminiProvider. The client is initialised but
// does not make any network requests at construction time.
func newGeminiProvider(apiKey, model string) (*GeminiProvider, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini client: %w", err)
	}
	return &GeminiProvider{client: client, model: model}, nil
}

func (p *GeminiProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	cfg := &genai.GenerateContentConfig{}
	if req.System != "" {
		cfg.SystemInstruction = genai.NewContentFromText(req.System, "")
	}
	if req.MaxTokens > 0 {
		cfg.MaxOutputTokens = int32(req.MaxTokens)
	}

	// All messages except the last become chat history
	var history []*genai.Content
	if len(req.Messages) > 1 {
		for _, m := range req.Messages[:len(req.Messages)-1] {
			var role genai.Role = genai.RoleUser
			if m.Role == "assistant" {
				role = genai.RoleModel
			}
			history = append(history, genai.NewContentFromText(m.Content, role))
		}
	}

	chat, err := p.client.Chats.Create(ctx, p.model, cfg, history)
	if err != nil {
		return nil, fmt.Errorf("gemini create chat: %w", err)
	}

	last := req.Messages[len(req.Messages)-1]
	resp, err := chat.SendMessage(ctx, genai.Part{Text: last.Content})
	if err != nil {
		return nil, fmt.Errorf("gemini complete: %w", err)
	}

	tokens := 0
	if resp.UsageMetadata != nil {
		tokens = int(resp.UsageMetadata.TotalTokenCount)
	}
	return &CompletionResponse{Text: resp.Text(), Tokens: tokens}, nil
}

func (p *GeminiProvider) IsAvailable() bool { return true }
func (p *GeminiProvider) Name() string      { return "gemini:" + p.model }
