package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnboardCommand_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "onboard" {
			found = true
			break
		}
	}
	assert.True(t, found, "onboard command must be registered on rootCmd")
}

func TestOnboardCommand_HasYesFlag(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "onboard" {
			assert.NotNil(t, c.Flags().Lookup("yes"))
			assert.NotNil(t, c.Flags().Lookup("y"))
			return
		}
	}
	t.Fatal("onboard command not found")
}

func TestOnboard_NonInteractive_WritesConfig(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	t.Setenv("ANTHROPIC_API_KEY", "sk-test-key")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")

	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)
	onboardYes = true
	t.Cleanup(func() { onboardYes = false })

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(filepath.Join(dir, ".caseforge.yaml"))
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "provider: anthropic")
	assert.Contains(t, content, "default_format: hurl")
}

func TestOnboard_SkipsExistingConfig(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	// Write existing config
	require.NoError(t, os.WriteFile(".caseforge.yaml", []byte("existing: true\n"), 0644))

	// Inject "n" to skip overwrite
	onboardCmd.SetIn(strings.NewReader("n\n"))
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	// Original file must be untouched
	data, _ := os.ReadFile(".caseforge.yaml")
	assert.Contains(t, string(data), "existing: true")
}

func TestOnboard_OverwritesOnConfirm(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	require.NoError(t, os.WriteFile(".caseforge.yaml", []byte("old: true\n"), 0644))

	t.Setenv("ANTHROPIC_API_KEY", "sk-test")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")

	// y=overwrite, then provider=1(anthropic), format=1(hurl), mcp=3(skip), skill=n
	onboardCmd.SetIn(strings.NewReader("y\n1\n1\n3\nn\n"))
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(".caseforge.yaml")
	require.NoError(t, err)
	assert.Contains(t, string(data), "provider: anthropic")
}

func TestOnboard_PrintsNextSteps(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	t.Setenv("ANTHROPIC_API_KEY", "sk-test")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")

	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)
	onboardYes = true
	t.Cleanup(func() { onboardYes = false })

	require.NoError(t, runOnboard(onboardCmd, nil))

	out := buf.String()
	assert.Contains(t, out, "Next steps")
	assert.Contains(t, out, "caseforge gen")
}

func TestOnboard_NoopProvider_SkipsAPIKeyPrompt(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")

	// provider=5(noop), format=1(hurl), mcp=3(skip), skill=n
	onboardCmd.SetIn(strings.NewReader("5\n1\n3\nn\n"))
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(".caseforge.yaml")
	require.NoError(t, err)
	assert.Contains(t, string(data), "provider: noop")
}

func TestOnboard_InstallMCP_WritesJSON(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "claude_desktop_config.json")
	require.NoError(t, installMCPToFile(configPath))

	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var cfg map[string]any
	require.NoError(t, json.Unmarshal(data, &cfg))

	servers, ok := cfg["mcpServers"].(map[string]any)
	require.True(t, ok)
	assert.NotNil(t, servers["caseforge"])
}

func TestOnboard_InstallMCP_Idempotent(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "claude_desktop_config.json")

	require.NoError(t, installMCPToFile(configPath))
	require.NoError(t, installMCPToFile(configPath)) // second call — must not error

	data, _ := os.ReadFile(configPath)
	var cfg map[string]any
	json.Unmarshal(data, &cfg)
	servers := cfg["mcpServers"].(map[string]any)
	assert.NotNil(t, servers["caseforge"])
}

func TestOnboard_InstallSkill_CopiesFile(t *testing.T) {
	skillSrc := filepath.Join(t.TempDir(), "SKILL.md")
	require.NoError(t, os.WriteFile(skillSrc, []byte("# CaseForge Skill\n"), 0644))

	dstDir := t.TempDir()
	dst := filepath.Join(dstDir, "caseforge.md")

	require.NoError(t, copySkillFile(skillSrc, dst))

	data, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Contains(t, string(data), "CaseForge Skill")
}
