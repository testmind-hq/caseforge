// internal/runner/hurl.go
package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HurlRunner struct{}

func NewHurlRunner() *HurlRunner { return &HurlRunner{} }

// hurlReport is the top-level structure of hurl --report-json output.
type hurlReport struct {
	Entries []struct {
		Filename string `json:"filename"`
		Success  bool   `json:"success"`
	} `json:"entries"`
}

// parseHurlReport counts passed/failed from hurl --report-json output.
func parseHurlReport(data []byte) (passed, failed int) {
	var report hurlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return 0, 0
	}
	for _, e := range report.Entries {
		if e.Success {
			passed++
		} else {
			failed++
		}
	}
	return passed, failed
}

// Run executes all .hurl files in casesDir against the given variables.
// Returns the number of passed, failed tests, and any execution error.
// err == nil even when some tests fail; test failures are expressed through failed > 0.
// err is reserved for infrastructure failures (binary not found, no .hurl files, cannot create temp dir).
func (r *HurlRunner) Run(casesDir string, vars map[string]string) (passed, failed int, err error) {
	if _, err := exec.LookPath("hurl"); err != nil {
		return 0, 0, fmt.Errorf("hurl not found on PATH — run `caseforge doctor` to check dependencies")
	}

	files, err := filepath.Glob(filepath.Join(casesDir, "*.hurl"))
	if err != nil || len(files) == 0 {
		return 0, 0, fmt.Errorf("no .hurl files found in %s", casesDir)
	}

	// Use a temp dir for the JSON report
	reportDir, err := os.MkdirTemp("", "caseforge-report-*")
	if err != nil {
		return 0, 0, fmt.Errorf("creating report dir: %w", err)
	}
	defer os.RemoveAll(reportDir)

	args := []string{"--test", "--report-json", reportDir}
	for k, v := range vars {
		args = append(args, "--variable", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, files...)

	cmd := exec.Command("hurl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	runErr := cmd.Run()

	// Parse the JSON report for accurate pass/fail counts
	reportPath := filepath.Join(reportDir, "report.json")
	if data, readErr := os.ReadFile(reportPath); readErr == nil {
		passed, failed = parseHurlReport(data)
		// If the report is empty and hurl exited non-zero, this is an
		// infrastructure failure (e.g., bad flag syntax), not a test failure.
		if passed+failed == 0 && runErr != nil {
			return 0, 0, fmt.Errorf("hurl infrastructure error: %w", runErr)
		}
		return passed, failed, nil
	}

	// Fallback: no report file — use exit code
	if runErr != nil {
		return 0, len(files), nil
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
