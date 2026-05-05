// cmd/onboard.go
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Interactive setup wizard — get started with CaseForge in minutes",
	RunE:  runOnboard,
}

func init() {
	rootCmd.AddCommand(onboardCmd)
	onboardCmd.Flags().BoolP("yes", "y", false, "Non-interactive: accept detected defaults")
}

// providerInfo holds detected provider info.
type providerInfo struct {
	name      string // "anthropic", "openai", etc.
	envKey    string // which env var to check
	available bool
	model     string // default model
}

type checkboxOption struct {
	label  string
	detail string
}

// promptCheckbox prints a numbered list and lets the user select items by number.
// Returns 0-based indices of selected items. Blank input selects nothing.
func promptCheckbox(out io.Writer, in *bufio.Reader, title string, opts []checkboxOption) []int {
	fmt.Fprintln(out, title)
	for i, o := range opts {
		fmt.Fprintf(out, "  [%d] %s  (%s)\n", i+1, o.label, o.detail)
	}
	fmt.Fprint(out, "Select [numbers e.g. 1 2], or enter to skip: ")
	line := strings.TrimSpace(readLine(in))

	seen := make(map[int]bool)
	for _, field := range strings.Fields(line) {
		var n int
		if _, err := fmt.Sscan(field, &n); err == nil && n >= 1 && n <= len(opts) {
			seen[n-1] = true
		}
	}
	result := make([]int, 0, len(seen))
	for i := range opts {
		if seen[i] {
			result = append(result, i)
		}
	}
	return result
}

// promptProviderDetails asks for model, api_key (and base_url for openai-compat)
// after the user has selected a provider.
func promptProviderDetails(out io.Writer, in *bufio.Reader, p providerInfo) (model, apiKey, baseURL string) {
	fmt.Fprintf(out, "Model [%s]: ", p.model)
	if m := strings.TrimSpace(readLine(in)); m != "" {
		model = m
	} else {
		model = p.model
	}

	// bedrock uses the AWS credential chain — no api_key prompt needed.
	// noop has no provider at all. Both are excluded from the wizard for now;
	// this guard future-proofs the function for when bedrock is added to detectProviders.
	if p.name != "noop" && p.name != "bedrock" {
		if p.available {
			fmt.Fprintf(out, "API key: (detected via %s) [Enter to keep, or paste to override]: ", p.envKey)
		} else {
			fmt.Fprintf(out, "API key for %s (leave blank to set %s later): ", p.name, p.envKey)
		}
		apiKey = strings.TrimSpace(readLine(in))
	}

	if p.name == "openai-compat" {
		for baseURL == "" {
			fmt.Fprint(out, "Base URL (e.g. https://api.deepseek.com/v1): ")
			baseURL = strings.TrimSpace(readLine(in))
		}
	}
	return
}

func detectProviders() []providerInfo {
	providers := []providerInfo{
		{"anthropic", "ANTHROPIC_API_KEY", os.Getenv("ANTHROPIC_API_KEY") != "", "claude-sonnet-4-6"},
		{"openai", "OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY") != "", "gpt-4o-mini"},
		{"gemini", "GEMINI_API_KEY", os.Getenv("GEMINI_API_KEY") != "" || os.Getenv("GOOGLE_API_KEY") != "", "gemini-2.5-flash"},
		{"openai-compat", "OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY") != "", "gpt-4o-mini"},
		{"noop", "", true, ""},
	}
	return providers
}

func runOnboard(cmd *cobra.Command, _ []string) error {
	onboardYes, _ := cmd.Flags().GetBool("yes")
	out := cmd.OutOrStdout()
	in := bufio.NewReader(cmd.InOrStdin())

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolving home directory: %w", err)
	}
	configPath := filepath.Join(home, ".caseforge.yaml")

	fmt.Fprintln(out, "Welcome to CaseForge!")
	fmt.Fprintln(out, "This wizard will set up your environment in a few steps.")
	fmt.Fprintln(out)

	// Step 1: Check existing config
	if _, err := os.Stat(configPath); err == nil {
		if onboardYes {
			fmt.Fprintf(out, "⚠  %s already exists — overwriting (--yes).\n", configPath)
		} else {
			fmt.Fprintf(out, "%s already exists. Overwrite? [y/N]: ", configPath)
			ans := readLine(in)
			if !strings.EqualFold(strings.TrimSpace(ans), "y") {
				fmt.Fprintln(out, "Keeping existing config. Done.")
				return nil
			}
		}
	}

	// Step 2: Detect providers
	providers := detectProviders()
	fmt.Fprintln(out, "Detected LLM providers:")
	for _, p := range providers[:4] { // exclude noop from detection display
		status := "✗ not available"
		if p.available {
			status = "✓ available (" + p.envKey + " is set)"
		}
		fmt.Fprintf(out, "  %-14s %s\n", p.name, status)
	}
	fmt.Fprintln(out)

	// Step 3: Choose provider
	var chosenProvider providerInfo
	if onboardYes {
		// Pick first available non-noop, fallback to noop
		chosenProvider = providers[4] // noop default
		for _, p := range providers[:4] {
			if p.available {
				chosenProvider = p
				break
			}
		}
		fmt.Fprintf(out, "Provider: %s (auto-selected)\n", chosenProvider.name)
	} else {
		fmt.Fprintln(out, "Choose your primary LLM provider:")
		for i, p := range providers {
			marker := "  "
			if p.available && p.name != "noop" {
				marker = "✓ "
			}
			fmt.Fprintf(out, "  [%d] %s%s\n", i+1, marker, p.name)
		}
		choice := promptInt(out, in, "Provider", 1, len(providers))
		chosenProvider = providers[choice-1]
	}

	// Step 3b+4: model, api_key, base_url sub-prompts
	model := chosenProvider.model
	apiKey := ""
	baseURL := ""
	if !onboardYes && chosenProvider.name != "noop" {
		model, apiKey, baseURL = promptProviderDetails(out, in, chosenProvider)
	}

	// Step 5: Choose output format
	formats := []string{"hurl", "postman", "k6", "markdown", "csv"}
	chosenFormat := "hurl"
	if onboardYes {
		fmt.Fprintf(out, "Output format: hurl (auto-selected)\n")
	} else {
		fmt.Fprintln(out, "\nChoose output format:")
		for i, f := range formats {
			fmt.Fprintf(out, "  [%d] %s\n", i+1, f)
		}
		choice := promptInt(out, in, "Format", 1, len(formats))
		chosenFormat = formats[choice-1]
	}

	// Step 6: Install MCP server (multi-select)
	if !onboardYes {
		desktopPath := claudeDesktopConfigPath(home)
		claudeCodePath := filepath.Join(home, ".claude.json")
		codexPath := filepath.Join(home, ".codex", "config.json")

		mcpOpts := []checkboxOption{
			{label: "Claude Desktop", detail: desktopPath},
			{label: "Claude Code", detail: claudeCodePath},
			{label: "Codex CLI", detail: codexPath},
		}
		selected := promptCheckbox(out, in, "\nInstall CaseForge as MCP server?", mcpOpts)
		paths := []string{desktopPath, claudeCodePath, codexPath}
		for _, idx := range selected {
			if err := installMCPToFile(paths[idx]); err != nil {
				fmt.Fprintf(out, "  ⚠  MCP install failed (%s): %v\n", mcpOpts[idx].label, err)
			} else {
				fmt.Fprintf(out, "  ✓ Registered in %s\n", paths[idx])
			}
		}
	}

	// Step 7: Install skill (multi-select)
	if !onboardYes {
		claudeSkillLink := filepath.Join(home, ".claude", "skills", "caseforge")
		universalDst := filepath.Join(home, ".agents", "skills", "caseforge", "SKILL.md")

		skillOpts := []checkboxOption{
			{label: "Claude Code / Desktop", detail: claudeSkillLink},
			{label: "Universal AI CLI (Gemini, Codex…)", detail: universalDst},
		}
		selected := promptCheckbox(out, in, "\nInstall CaseForge skill?", skillOpts)
		if len(selected) > 0 {
			src := findSkillFile()
			if src == "" {
				fmt.Fprintln(out, "  ⚠  Skill file not found (run from caseforge source directory).")
			} else {
				for _, idx := range selected {
					var installErr error
					var displayPath string
					if idx == 0 {
						installErr = installClaudeCodeSkill(home, src)
						displayPath = claudeSkillLink
					} else {
						installErr = copySkillFile(src, universalDst)
						displayPath = universalDst
					}
					if installErr != nil {
						fmt.Fprintf(out, "  ⚠  Skill install failed (%s): %v\n", skillOpts[idx].label, installErr)
					} else {
						fmt.Fprintf(out, "  ✓ Skill installed at %s\n", displayPath)
					}
				}
			}
		}
	}

	// Step 8: Write config
	if err := writeOnboardConfig(configPath, chosenProvider.name, model, apiKey, baseURL, chosenFormat); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	fmt.Fprintf(out, "\n✓ %s written.\n", configPath)
	fmt.Fprintln(out, "  Tip: place a .caseforge.yaml in any project directory to override the global config.")

	// Next steps
	printNextSteps(out, chosenFormat)
	return nil
}

func writeOnboardConfig(path, provider, model, apiKey, baseURL, format string) error {
	cfg := map[string]any{
		"ai": map[string]any{
			"provider": provider,
			"model":    model,
			"api_key":  apiKey,
			"base_url": baseURL,
		},
		"output": map[string]any{
			"default_format": format,
			"dir":            "./cases",
		},
		"lint": map[string]any{
			"fail_on": "error",
		},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	header := []byte("# .caseforge.yaml — generated by caseforge onboard\n")
	return os.WriteFile(path, append(header, data...), 0644)
}

// installMCPToFile merges the caseforge MCP entry into a JSON config file.
func installMCPToFile(path string) error {
	var cfg map[string]any

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
	}
	if cfg == nil {
		cfg = map[string]any{}
	}

	servers, _ := cfg["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}
	if _, exists := servers["caseforge"]; exists {
		return nil // already installed — idempotent
	}
	servers["caseforge"] = map[string]any{
		"command": "caseforge",
		"args":    []string{"mcp"},
	}
	cfg["mcpServers"] = servers

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

// installClaudeCodeSkill copies src to ~/.agents/skills/caseforge/SKILL.md and
// creates a relative symlink ~/.claude/skills/caseforge → ../../.agents/skills/caseforge,
// matching the convention used by other Claude Code skills.
func installClaudeCodeSkill(home, src string) error {
	agentsDst := filepath.Join(home, ".agents", "skills", "caseforge", "SKILL.md")
	if err := copySkillFile(src, agentsDst); err != nil {
		return err
	}
	claudeLink := filepath.Join(home, ".claude", "skills", "caseforge")
	if _, err := os.Lstat(claudeLink); err == nil {
		return nil // already exists — idempotent
	}
	if err := os.MkdirAll(filepath.Join(home, ".claude", "skills"), 0755); err != nil {
		return err
	}
	// Relative target: from ~/.claude/skills/, ../../.agents/skills/caseforge resolves to ~/.agents/skills/caseforge
	return os.Symlink(filepath.Join("..", "..", ".agents", "skills", "caseforge"), claudeLink)
}

// copySkillFile copies src → dst, creating parent dirs as needed.
func copySkillFile(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return nil // already installed — idempotent
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// findSkillFile locates skills/caseforge/SKILL.md relative to the executable or cwd.
func findSkillFile() string {
	// Try relative to cwd (development / source directory)
	if _, err := os.Stat("skills/caseforge/SKILL.md"); err == nil {
		abs, _ := filepath.Abs("skills/caseforge/SKILL.md")
		return abs
	}
	// Try relative to executable binary
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), "..", "skills", "caseforge", "SKILL.md")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func claudeDesktopConfigPath(home string) string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	}
	return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
}

func printNextSteps(out io.Writer, format string) {
	fmt.Fprintln(out, "\nNext steps:")
	fmt.Fprintln(out, "  1. Generate test cases:")
	fmt.Fprintln(out, "     caseforge gen --spec openapi.yaml")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "  2. Run them (format: %s):\n", format)
	fmt.Fprintf(out, "     caseforge run --cases ./cases --format %s\n", format)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  3. Need help?")
	fmt.Fprintln(out, "     caseforge --help")
}

func promptInt(out io.Writer, in *bufio.Reader, label string, min, max int) int {
	for {
		fmt.Fprintf(out, "%s [%d-%d]: ", label, min, max)
		line := strings.TrimSpace(readLine(in))
		var n int
		if _, err := fmt.Sscan(line, &n); err == nil && n >= min && n <= max {
			return n
		}
		fmt.Fprintf(out, "  Please enter a number between %d and %d.\n", min, max)
	}
}

func readLine(in *bufio.Reader) string {
	line, _ := in.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}
