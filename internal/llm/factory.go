// internal/llm/factory.go
package llm

import "os"

// NewProvider returns the appropriate LLMProvider based on config.
// Falls back to NoopProvider if no API key is available.
func NewProvider(apiKey, providerName, model string) LLMProvider {
	// Resolve API key: explicit arg > env var
	key := apiKey
	if key == "" {
		key = os.Getenv("ANTHROPIC_API_KEY")
	}

	switch providerName {
	case "anthropic":
		if key == "" {
			return &NoopProvider{}
		}
		return &AnthropicProvider{
			client: newAnthropicClient(key),
			model:  firstNonEmpty(model, "claude-sonnet-4-6"),
		}
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
