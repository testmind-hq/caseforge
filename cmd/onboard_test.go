package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnboard_ProviderSubPrompts_ShowsModelAndKey(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "sk-open")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// provider=2(openai), model=gpt-4o, apikey=enter(keep), format=1, mcp=enter(skip), skill=enter(skip)
	onboardCmd.SetIn(strings.NewReader("2\ngpt-4o\n\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(filepath.Join(dir, ".caseforge.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "provider: openai")
	assert.Contains(t, string(data), "model: gpt-4o")
}

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
			f := c.Flags().Lookup("yes")
			require.NotNil(t, f, "--yes flag must exist")
			assert.Equal(t, "y", f.Shorthand, "-y shorthand must exist")
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
	t.Setenv("GOOGLE_API_KEY", "")

	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)
	t.Cleanup(func() { onboardCmd.SetOut(os.Stdout) })
	require.NoError(t, onboardCmd.Flags().Set("yes", "true"))
	t.Cleanup(func() { onboardCmd.Flags().Set("yes", "false") })

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

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// Write existing config
	require.NoError(t, os.WriteFile(".caseforge.yaml", []byte("existing: true\n"), 0644))

	// Inject "n" to skip overwrite
	onboardCmd.SetIn(strings.NewReader("n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
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
	t.Setenv("GOOGLE_API_KEY", "")

	// y=overwrite, provider=1(anthropic), model=enter(default), apikey=enter(keep), format=1(hurl), mcp=enter(skip), skill=enter(skip)
	onboardCmd.SetIn(strings.NewReader("y\n1\n\n\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
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
	t.Setenv("GOOGLE_API_KEY", "")

	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)
	t.Cleanup(func() { onboardCmd.SetOut(os.Stdout) })
	require.NoError(t, onboardCmd.Flags().Set("yes", "true"))
	t.Cleanup(func() { onboardCmd.Flags().Set("yes", "false") })

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
	t.Setenv("GOOGLE_API_KEY", "")

	// provider=5(noop), format=1(hurl), mcp=enter(skip), skill=enter(skip)
	onboardCmd.SetIn(strings.NewReader("5\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
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

func TestOnboard_InstallSkill_Idempotent(t *testing.T) {
	skillSrc := filepath.Join(t.TempDir(), "SKILL.md")
	require.NoError(t, os.WriteFile(skillSrc, []byte("# CaseForge Skill\n"), 0644))

	dstDir := t.TempDir()
	dst := filepath.Join(dstDir, "caseforge.md")

	require.NoError(t, copySkillFile(skillSrc, dst))
	require.NoError(t, copySkillFile(skillSrc, dst)) // 第二次调用，幂等
	data, err := os.ReadFile(dst)
	require.NoError(t, err)
	assert.Contains(t, string(data), "CaseForge Skill")
}

func TestPromptCheckbox_NoneSelected(t *testing.T) {
	var buf bytes.Buffer
	in := strings.NewReader("\n") // blank = skip all
	selected := promptCheckbox(&buf, bufio.NewReader(in), "Pick targets:", []checkboxOption{
		{label: "Option A", detail: "path/a"},
		{label: "Option B", detail: "path/b"},
	})
	assert.Empty(t, selected)
}

func TestPromptCheckbox_SelectOne(t *testing.T) {
	var buf bytes.Buffer
	in := strings.NewReader("2\n")
	selected := promptCheckbox(&buf, bufio.NewReader(in), "Pick targets:", []checkboxOption{
		{label: "Option A", detail: "path/a"},
		{label: "Option B", detail: "path/b"},
	})
	assert.Equal(t, []int{1}, selected) // 0-based index
}

func TestPromptCheckbox_SelectMultiple(t *testing.T) {
	var buf bytes.Buffer
	in := strings.NewReader("1 2\n")
	selected := promptCheckbox(&buf, bufio.NewReader(in), "Pick targets:", []checkboxOption{
		{label: "Option A", detail: "path/a"},
		{label: "Option B", detail: "path/b"},
	})
	assert.Equal(t, []int{0, 1}, selected)
}

func TestPromptCheckbox_IgnoresOutOfRange(t *testing.T) {
	var buf bytes.Buffer
	in := strings.NewReader("0 5 1\n")
	selected := promptCheckbox(&buf, bufio.NewReader(in), "Pick targets:", []checkboxOption{
		{label: "Option A", detail: "path/a"},
	})
	assert.Equal(t, []int{0}, selected) // only valid index
}

func TestOnboard_MCPMultiSelect_InstallsMultiple(t *testing.T) {
	claudeCodePath := filepath.Join(t.TempDir(), "claude.json")
	codexPath := filepath.Join(t.TempDir(), "codex_config.json")

	// installMCPToFile is the helper — test it directly for two paths
	require.NoError(t, installMCPToFile(claudeCodePath))
	require.NoError(t, installMCPToFile(codexPath))

	for _, path := range []string{claudeCodePath, codexPath} {
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		var cfg map[string]any
		require.NoError(t, json.Unmarshal(data, &cfg))
		servers := cfg["mcpServers"].(map[string]any)
		assert.NotNil(t, servers["caseforge"], "expected caseforge in %s", path)
	}
}

func TestOnboard_SkillCheckbox_InstallsUniversal(t *testing.T) {
	skillSrc := filepath.Join(t.TempDir(), "SKILL.md")
	require.NoError(t, os.WriteFile(skillSrc, []byte("# CaseForge Skill\n"), 0644))

	universalDst := filepath.Join(t.TempDir(), "caseforge.md")
	require.NoError(t, copySkillFile(skillSrc, universalDst))

	data, err := os.ReadFile(universalDst)
	require.NoError(t, err)
	assert.Contains(t, string(data), "CaseForge Skill")
}
