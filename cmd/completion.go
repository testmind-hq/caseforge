// cmd/completion.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate an autocompletion script for caseforge in the specified shell.

Usage:
  # Bash (add to ~/.bashrc)
  source <(caseforge completion bash)

  # Zsh (write to a $fpath directory)
  caseforge completion zsh > "${fpath[1]}/_caseforge"

  # Fish
  caseforge completion fish > ~/.config/fish/completions/caseforge.fish

  # PowerShell
  caseforge completion powershell | Out-String | Invoke-Expression`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE:      runCompletion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()
	switch args[0] {
	case "bash":
		return rootCmd.GenBashCompletion(out)
	case "zsh":
		return rootCmd.GenZshCompletion(out)
	case "fish":
		return rootCmd.GenFishCompletion(out, true)
	case "powershell":
		return rootCmd.GenPowerShellCompletionWithDesc(out)
	default:
		return fmt.Errorf("unsupported shell %q: choose one of bash, zsh, fish, powershell", args[0])
	}
}
