// internal/config/config.go
package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AI     AIConfig     `mapstructure:"ai"`
	Output OutputConfig `mapstructure:"output"`
	Lint   LintConfig   `mapstructure:"lint"`
}

type AIConfig struct {
	Provider string `mapstructure:"provider"` // "anthropic"|"openai"|"openai-compat"|"gemini"|"noop"
	Model    string `mapstructure:"model"`
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"` // openai-compat only (DeepSeek, Qwen, Azure, etc.)
}

type OutputConfig struct {
	DefaultFormat string `mapstructure:"default_format"` // "hurl"|"markdown"|"csv"
	Dir           string `mapstructure:"dir"`
}

type LintConfig struct {
	FailOn string `mapstructure:"fail_on"` // "warning"|"error"
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
	return &cfg, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
