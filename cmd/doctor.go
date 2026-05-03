// cmd/doctor.go
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment dependencies",
	RunE:  runDoctor,
}

func init() { rootCmd.AddCommand(doctorCmd) }

func runDoctor(cmd *cobra.Command, args []string) error {
	ok := true

	// Check hurl
	if _, err := exec.LookPath("hurl"); err != nil {
		color.Red("  ✗ hurl not found — install from https://hurl.dev/docs/installation.html")
		ok = false
	} else {
		color.Green("  ✓ hurl found")
	}

	// Check k6
	if _, err := exec.LookPath("k6"); err != nil {
		color.Yellow("  ⚠ k6 not found — k6 runner disabled (install from https://k6.io/docs/get-started/installation/)")
	} else {
		color.Green("  ✓ k6 found")
	}

	// Check AI provider keys — at least one is required for AI features
	hasAnthropic := os.Getenv("ANTHROPIC_API_KEY") != ""
	hasOpenAI := os.Getenv("OPENAI_API_KEY") != ""
	hasGemini := os.Getenv("GEMINI_API_KEY") != "" || os.Getenv("GOOGLE_API_KEY") != ""

	if hasAnthropic {
		color.Green("  ✓ ANTHROPIC_API_KEY set")
	} else {
		color.Yellow("  ⚠ ANTHROPIC_API_KEY not set — anthropic provider unavailable")
	}
	if hasOpenAI {
		color.Green("  ✓ OPENAI_API_KEY set")
	} else {
		color.Yellow("  ⚠ OPENAI_API_KEY not set — openai/openai-compat provider unavailable")
	}
	if hasGemini {
		color.Green("  ✓ Gemini API key set (GEMINI_API_KEY or GOOGLE_API_KEY)")
	} else {
		color.Yellow("  ⚠ neither GEMINI_API_KEY nor GOOGLE_API_KEY is set — gemini provider unavailable")
	}
	if !hasAnthropic && !hasOpenAI && !hasGemini {
		color.Yellow("  ⚠ no AI provider key set — AI features disabled (use --no-ai or set at least one key)")
	}

	if !ok {
		return fmt.Errorf("environment check failed")
	}
	return nil
}
