// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AI     AIConfig     `mapstructure:"ai"`
	Output OutputConfig `mapstructure:"output"`
	Lint   LintConfig   `mapstructure:"lint"`
}

type AIConfig struct {
	Provider string `mapstructure:"provider"` // "anthropic"|"noop"
	Model    string `mapstructure:"model"`
	APIKey   string `mapstructure:"api_key"`
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
		cfg.AI.APIKey = viper.GetString("ANTHROPIC_API_KEY")
	}
	return &cfg, nil
}
