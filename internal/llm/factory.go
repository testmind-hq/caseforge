// internal/llm/factory.go
package llm

import (
	"os"
)

// ProviderConfig holds all parameters for constructing an LLMProvider.
type ProviderConfig struct {
	APIKey   string
	Provider string
	Model    string
	BaseURL  string // openai-compat only
	Region   string // bedrock only
}

// NewProvider is a convenience wrapper for the three common fields.
func NewProvider(apiKey, providerName, model string) LLMProvider {
	return NewProviderWithConfig(ProviderConfig{
		APIKey:   apiKey,
		Provider: providerName,
		Model:    model,
	})
}

// NewProviderWithConfig constructs an LLMProvider from a ProviderConfig.
func NewProviderWithConfig(cfg ProviderConfig) LLMProvider {
	switch cfg.Provider {
	case "anthropic":
		key := firstNonEmpty(cfg.APIKey, os.Getenv("ANTHROPIC_API_KEY"))
		if key == "" {
			return &NoopProvider{}
		}
		return &AnthropicProvider{
			client: newAnthropicClient(key),
			model:  firstNonEmpty(cfg.Model, "claude-sonnet-4-6"),
		}
	case "openai", "openai-compat":
		key := firstNonEmpty(cfg.APIKey, os.Getenv("OPENAI_API_KEY"))
		if key == "" {
			return &NoopProvider{}
		}
		return NewOpenAIProvider(OpenAIConfig{
			APIKey:  key,
			Model:   firstNonEmpty(cfg.Model, "gpt-4o-mini"),
			BaseURL: cfg.BaseURL,
		})
	case "gemini":
		key := firstNonEmpty(cfg.APIKey, os.Getenv("GEMINI_API_KEY"), os.Getenv("GOOGLE_API_KEY"))
		if key == "" {
			return &NoopProvider{}
		}
		p, err := newGeminiProvider(key, firstNonEmpty(cfg.Model, "gemini-2.5-flash"))
		if err != nil {
			return &NoopProvider{}
		}
		return p
	}
	return &NoopProvider{}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
