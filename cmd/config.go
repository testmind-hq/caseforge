// cmd/config.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/config"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CaseForge configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the effective configuration",
	RunE:  runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
}

func runConfigShow(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Mask API key: show only first 6 chars for keys longer than 8 chars
	apiKey := cfg.AI.APIKey
	if len(apiKey) > 8 {
		apiKey = apiKey[:6] + "..."
	}

	display := map[string]any{
		"ai": map[string]any{
			"provider": cfg.AI.Provider,
			"model":    cfg.AI.Model,
			"api_key":  apiKey,
			"base_url": cfg.AI.BaseURL,
		},
		"output": map[string]any{
			"default_format": cfg.Output.DefaultFormat,
			"dir":            cfg.Output.Dir,
		},
		"lint": map[string]any{
			"fail_on": cfg.Lint.FailOn,
		},
	}

	data, err := yaml.Marshal(display)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	fmt.Fprint(cmd.OutOrStdout(), string(data))
	return nil
}
