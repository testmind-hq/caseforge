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

	// Check ANTHROPIC_API_KEY
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		color.Green("  ✓ ANTHROPIC_API_KEY set")
	} else {
		color.Yellow("  ⚠ ANTHROPIC_API_KEY not set — AI features disabled (use --no-ai or set the key)")
	}

	// Check OPENAI_API_KEY
	if os.Getenv("OPENAI_API_KEY") != "" {
		color.Green("  ✓ OPENAI_API_KEY set")
	} else {
		color.Yellow("  ⚠ OPENAI_API_KEY not set — openai/openai-compat provider unavailable")
	}

	// Check GEMINI_API_KEY / GOOGLE_API_KEY
	if os.Getenv("GEMINI_API_KEY") != "" || os.Getenv("GOOGLE_API_KEY") != "" {
		color.Green("  ✓ Gemini API key set (GEMINI_API_KEY or GOOGLE_API_KEY)")
	} else {
		color.Yellow("  ⚠ GEMINI_API_KEY not set — gemini provider unavailable")
	}

	if !ok {
		return fmt.Errorf("environment check failed")
	}
	return nil
}
