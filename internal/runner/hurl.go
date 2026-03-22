// internal/runner/hurl.go
package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HurlRunner struct{}

func NewHurlRunner() *HurlRunner { return &HurlRunner{} }

// Run executes all .hurl files in casesDir against the given variables.
// Returns the number of passed, failed tests, and any execution error.
func (r *HurlRunner) Run(casesDir string, vars map[string]string) (passed, failed int, err error) {
	if _, err := exec.LookPath("hurl"); err != nil {
		return 0, 0, fmt.Errorf("hurl not found on PATH — run `caseforge doctor` to check dependencies")
	}

	files, err := filepath.Glob(filepath.Join(casesDir, "*.hurl"))
	if err != nil || len(files) == 0 {
		return 0, 0, fmt.Errorf("no .hurl files found in %s", casesDir)
	}

	args := []string{"--test"}
	for k, v := range vars {
		args = append(args, "--variable", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, files...)

	cmd := exec.Command("hurl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			_ = exitErr
			return 0, len(files), nil // hurl exits non-zero on test failures
		}
		return 0, 0, err
	}
	return len(files), 0, nil
}

// ParseVars converts "key=value" strings to a map.
func ParseVars(varFlags []string) map[string]string {
	result := make(map[string]string)
	for _, v := range varFlags {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
