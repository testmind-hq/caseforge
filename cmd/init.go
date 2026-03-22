// cmd/init.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize CaseForge in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		const defaultConfig = `# .caseforge.yaml
ai:
  provider: noop        # change to "anthropic" and set ANTHROPIC_API_KEY to enable AI
  model: claude-sonnet-4-6
  api_key: ""

output:
  default_format: hurl
  dir: ./cases

lint:
  fail_on: error
`
		if _, err := os.Stat(".caseforge.yaml"); err == nil {
			return fmt.Errorf(".caseforge.yaml already exists")
		}
		return os.WriteFile(".caseforge.yaml", []byte(defaultConfig), 0644)
	},
}

func init() { rootCmd.AddCommand(initCmd) }
