// internal/llm/factory.go
package llm

import "os"

// NewProvider returns the appropriate LLMProvider based on config.
func NewProvider(apiKey, providerName, model string) LLMProvider {
	return NewProviderWithConfig(apiKey, providerName, model, "")
}

// NewProviderWithConfig is like NewProvider but also accepts baseURL for
// openai-compat providers (DeepSeek, Qwen, Moonshot, Azure, etc.).
func NewProviderWithConfig(apiKey, providerName, model, baseURL string) LLMProvider {
	switch providerName {
	case "anthropic":
		key := firstNonEmpty(apiKey, os.Getenv("ANTHROPIC_API_KEY"))
		if key == "" {
			return &NoopProvider{}
		}
		return &AnthropicProvider{
			client: newAnthropicClient(key),
			model:  firstNonEmpty(model, "claude-sonnet-4-6"),
		}
	case "openai", "openai-compat":
		// Both openai and openai-compat use OPENAI_API_KEY. For compat providers
		// (DeepSeek, Qwen, Moonshot, Azure) set OPENAI_API_KEY to the provider's key.
		key := firstNonEmpty(apiKey, os.Getenv("OPENAI_API_KEY"))
		if key == "" {
			return &NoopProvider{}
		}
		return NewOpenAIProvider(OpenAIConfig{
			APIKey:  key,
			Model:   firstNonEmpty(model, "gpt-4o-mini"),
			BaseURL: baseURL,
		})
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
