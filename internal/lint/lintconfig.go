// internal/lint/lintconfig.go
package lint

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LintFileConfig holds settings from .caseforgelint.yaml.
type LintFileConfig struct {
	SkipRules []string `yaml:"skip_rules"`
	FailOn    string   `yaml:"fail_on"`
}

// LoadLintFileConfig reads .caseforgelint.yaml from dir.
// Returns an empty config (no error) if the file does not exist.
func LoadLintFileConfig(dir string) (LintFileConfig, error) {
	path := filepath.Join(dir, ".caseforgelint.yaml")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return LintFileConfig{}, nil
	}
	if err != nil {
		return LintFileConfig{}, fmt.Errorf("reading .caseforgelint.yaml: %w", err)
	}
	var cfg LintFileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return LintFileConfig{}, fmt.Errorf("parsing .caseforgelint.yaml: %w", err)
	}
	return cfg, nil
}
