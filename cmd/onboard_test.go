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
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "sk-open")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// provider=2(openai), apikey=enter(keep), model=gpt-4o, format=1, mcp=enter(skip), skill=enter(skip)
	onboardCmd.SetIn(strings.NewReader("2\n\ngpt-4o\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(filepath.Join(home, ".caseforge.yaml"))
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
	home := t.TempDir()
	t.Setenv("HOME", home)

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

	data, err := os.ReadFile(filepath.Join(home, ".caseforge.yaml"))
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "provider: anthropic")
	assert.Contains(t, content, "default_format: hurl")
}

func TestOnboard_SkipsExistingConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	require.NoError(t, os.WriteFile(filepath.Join(home, ".caseforge.yaml"), []byte("existing: true\n"), 0644))

	onboardCmd.SetIn(strings.NewReader("n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, _ := os.ReadFile(filepath.Join(home, ".caseforge.yaml"))
	assert.Contains(t, string(data), "existing: true")
}

func TestOnboard_OverwritesOnConfirm(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	require.NoError(t, os.WriteFile(filepath.Join(home, ".caseforge.yaml"), []byte("old: true\n"), 0644))

	t.Setenv("ANTHROPIC_API_KEY", "sk-test")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// y=overwrite, provider=1(anthropic), model=enter(default), apikey=enter(keep),
	// format=1(hurl), mcp=enter(skip), skill=enter(skip)
	onboardCmd.SetIn(strings.NewReader("y\n1\n\n\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(filepath.Join(home, ".caseforge.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "provider: anthropic")
}

func TestOnboard_PrintsNextSteps(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

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
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// provider=6(noop), format=1(hurl), mcp=enter(skip), skill=enter(skip)
	onboardCmd.SetIn(strings.NewReader("6\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(filepath.Join(home, ".caseforge.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), "provider: noop")
}

func TestOnboard_BedrockProvider_PromptRegion(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIATEST")
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("AWS_DEFAULT_REGION", "ap-northeast-1")

	// provider=5(bedrock), region=enter(use default), model=enter(default), format=1, mcp=enter, skill=enter
	onboardCmd.SetIn(strings.NewReader("5\n\n\n1\n\n\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	data, err := os.ReadFile(filepath.Join(home, ".caseforge.yaml"))
	require.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "provider: bedrock")
	assert.Contains(t, content, "region: ap-northeast-1")
	assert.NotContains(t, content, "api_key")
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

func TestOnboard_InstallClaudeCodeSkill_CreatesSymlink(t *testing.T) {
	home := t.TempDir()
	skillSrc := filepath.Join(t.TempDir(), "SKILL.md")
	require.NoError(t, os.WriteFile(skillSrc, []byte("# CaseForge Skill\n"), 0644))

	require.NoError(t, installClaudeCodeSkill(home, skillSrc))

	// Real file in ~/.agents/
	agentsDst := filepath.Join(home, ".agents", "skills", "caseforge", "SKILL.md")
	data, err := os.ReadFile(agentsDst)
	require.NoError(t, err)
	assert.Contains(t, string(data), "CaseForge Skill")

	// Symlink at ~/.claude/skills/caseforge
	claudeLink := filepath.Join(home, ".claude", "skills", "caseforge")
	info, err := os.Lstat(claudeLink)
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink)
	target, err := os.Readlink(claudeLink)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("..", "..", ".agents", "skills", "caseforge"), target)
}

func TestOnboard_InstallClaudeCodeSkill_Idempotent(t *testing.T) {
	home := t.TempDir()
	skillSrc := filepath.Join(t.TempDir(), "SKILL.md")
	require.NoError(t, os.WriteFile(skillSrc, []byte("# CaseForge Skill\n"), 0644))

	require.NoError(t, installClaudeCodeSkill(home, skillSrc))
	require.NoError(t, installClaudeCodeSkill(home, skillSrc)) // second call must not error

	claudeLink := filepath.Join(home, ".claude", "skills", "caseforge")
	_, err := os.Lstat(claudeLink)
	require.NoError(t, err)
	target, err := os.Readlink(claudeLink)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("..", "..", ".agents", "skills", "caseforge"), target)
}

func TestOnboard_InstallClaudeCodeSkill_RejectsExistingNonSymlink(t *testing.T) {
	home := t.TempDir()
	skillSrc := filepath.Join(t.TempDir(), "SKILL.md")
	require.NoError(t, os.WriteFile(skillSrc, []byte("# CaseForge Skill\n"), 0644))

	claudeLink := filepath.Join(home, ".claude", "skills", "caseforge")
	require.NoError(t, os.MkdirAll(filepath.Dir(claudeLink), 0755))
	require.NoError(t, os.WriteFile(claudeLink, []byte("stale content"), 0644))

	err := installClaudeCodeSkill(home, skillSrc)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a symlink")
}

func TestOnboard_SkillCheckbox_InstallsClaudeCode(t *testing.T) {
	srcDir := t.TempDir()
	skillDir := filepath.Join(srcDir, "skills", "caseforge")
	require.NoError(t, os.MkdirAll(skillDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# CaseForge Skill\n"), 0644))

	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(srcDir))
	t.Cleanup(func() { os.Chdir(orig) })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ANTHROPIC_API_KEY", "sk-test")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// provider=1(anthropic), model=enter, apikey=enter, format=1, mcp=enter(skip), skill=1(Claude Code)
	onboardCmd.SetIn(strings.NewReader("1\n\n\n1\n\n1\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	agentsDst := filepath.Join(home, ".agents", "skills", "caseforge", "SKILL.md")
	data, err := os.ReadFile(agentsDst)
	require.NoError(t, err)
	assert.Contains(t, string(data), "CaseForge Skill")

	claudeLink := filepath.Join(home, ".claude", "skills", "caseforge")
	info, err := os.Lstat(claudeLink)
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink)
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
	// Set up a temp src dir with a fake SKILL.md so findSkillFile() returns a path
	srcDir := t.TempDir()
	skillDir := filepath.Join(srcDir, "skills", "caseforge")
	require.NoError(t, os.MkdirAll(skillDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# CaseForge Skill\n"), 0644))

	orig, _ := os.Getwd()
	require.NoError(t, os.Chdir(srcDir))
	t.Cleanup(func() { os.Chdir(orig) })

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ANTHROPIC_API_KEY", "sk-test")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	// provider=1(anthropic), model=enter, apikey=enter, format=1, mcp=enter(skip), skill=2(universal)
	onboardCmd.SetIn(strings.NewReader("1\n\n\n1\n\n2\n"))
	t.Cleanup(func() { onboardCmd.SetIn(os.Stdin); onboardCmd.SetOut(os.Stdout) })
	var buf bytes.Buffer
	onboardCmd.SetOut(&buf)

	require.NoError(t, runOnboard(onboardCmd, nil))

	universalDst := filepath.Join(home, ".agents", "skills", "caseforge", "SKILL.md")
	data, err := os.ReadFile(universalDst)
	require.NoError(t, err)
	assert.Contains(t, string(data), "CaseForge Skill")
}
