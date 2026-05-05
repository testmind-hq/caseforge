// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AI       AIConfig        `mapstructure:"ai"`
	Output   OutputConfig    `mapstructure:"output"`
	Lint     LintConfig      `mapstructure:"lint"`
	Webhooks []WebhookConfig `mapstructure:"webhooks"`
}

// WebhookConfig defines a single webhook endpoint to receive CaseForge events.
type WebhookConfig struct {
	URL         string   `mapstructure:"url"`
	Events      []string `mapstructure:"events"`          // "on_generate", "on_run_complete"
	Secret      string   `mapstructure:"secret"`          // optional HMAC-SHA256 signing key
	TimeoutSecs int      `mapstructure:"timeout_seconds"` // default 10
	MaxRetries  int      `mapstructure:"max_retries"`     // default 3
}

type AIConfig struct {
	Provider string `mapstructure:"provider"` // "anthropic"|"openai"|"openai-compat"|"gemini"|"bedrock"|"noop"
	Model    string `mapstructure:"model"`
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"` // openai-compat only (DeepSeek, Qwen, Azure, etc.)
	Region   string `mapstructure:"region"`   // bedrock only; falls back to AWS_REGION, then AWS_DEFAULT_REGION
}

type OutputConfig struct {
	DefaultFormat string `mapstructure:"default_format"` // "hurl"|"markdown"|"csv"
	Dir           string `mapstructure:"dir"`
}

type LintConfig struct {
	FailOn    string   `mapstructure:"fail_on"`    // "warning"|"error"
	SkipRules []string `mapstructure:"skip_rules"` // rules to skip
}

func Load() (*Config, error) {
	viper.SetDefault("ai.provider", "noop")
	viper.SetDefault("ai.model", "claude-sonnet-4-6")
	viper.SetDefault("output.default_format", "hurl")
	viper.SetDefault("output.dir", "./cases")
	viper.SetDefault("lint.fail_on", "error")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	// API key override from environment
	if cfg.AI.APIKey == "" {
		switch cfg.AI.Provider {
		case "openai", "openai-compat":
			cfg.AI.APIKey = os.Getenv("OPENAI_API_KEY")
		case "gemini":
			cfg.AI.APIKey = firstNonEmpty(os.Getenv("GEMINI_API_KEY"), os.Getenv("GOOGLE_API_KEY"))
		default: // "anthropic" and anything unrecognised
			cfg.AI.APIKey = os.Getenv("ANTHROPIC_API_KEY")
		}
	}
	// Region override from environment
	if cfg.AI.Provider == "bedrock" {
		if cfg.AI.Region == "" {
			cfg.AI.Region = firstNonEmpty(os.Getenv("AWS_REGION"), os.Getenv("AWS_DEFAULT_REGION"))
		}
	}
	if err := cfg.AI.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

var knownProviders = map[string]bool{
	"anthropic": true, "openai": true, "gemini": true,
	"openai-compat": true, "bedrock": true, "noop": true,
}

// Validate checks that the AI config fields are self-consistent.
// It is called automatically by Load() and returns a user-friendly error.
func (c *AIConfig) Validate() error {
	if !knownProviders[c.Provider] {
		return fmt.Errorf("config: unknown ai.provider %q — valid values: anthropic, openai, gemini, openai-compat, bedrock, noop", c.Provider)
	}
	if strings.HasPrefix(c.APIKey, "http://") || strings.HasPrefix(c.APIKey, "https://") {
		return fmt.Errorf("config: ai.api_key looks like a URL (%q) — did you swap api_key and base_url?", c.APIKey)
	}
	if c.BaseURL != "" && !strings.HasPrefix(c.BaseURL, "http://") && !strings.HasPrefix(c.BaseURL, "https://") {
		return fmt.Errorf("config: ai.base_url %q has no HTTP scheme — it must start with https:// (or http://)", c.BaseURL)
	}
	if c.Provider == "openai-compat" && c.BaseURL == "" {
		return fmt.Errorf("config: provider \"openai-compat\" requires ai.base_url to be set")
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
