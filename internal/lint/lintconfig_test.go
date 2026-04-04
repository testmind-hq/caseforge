// internal/lint/lintconfig_test.go
package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadLintFileConfig_FileExists(t *testing.T) {
	dir := t.TempDir()
	content := "skip_rules:\n  - L001\n  - L003\nfail_on: warning\n"
	err := os.WriteFile(filepath.Join(dir, ".caseforgelint.yaml"), []byte(content), 0644)
	assert.NoError(t, err)

	cfg, err := LoadLintFileConfig(dir)
	assert.NoError(t, err)
	assert.Equal(t, []string{"L001", "L003"}, cfg.SkipRules)
	assert.Equal(t, "warning", cfg.FailOn)
}

func TestLoadLintFileConfig_FileMissing(t *testing.T) {
	dir := t.TempDir()
	cfg, err := LoadLintFileConfig(dir)
	assert.NoError(t, err)
	assert.Empty(t, cfg.SkipRules)
	assert.Empty(t, cfg.FailOn)
}

func TestLoadLintFileConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, ".caseforgelint.yaml"), []byte(""), 0644)
	assert.NoError(t, err)

	cfg, err := LoadLintFileConfig(dir)
	assert.NoError(t, err)
	assert.Empty(t, cfg.SkipRules)
}
