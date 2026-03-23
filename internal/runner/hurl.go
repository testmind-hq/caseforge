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

type hurlReport struct {
	Entries []struct {
		Filename string `json:"filename"`
		Success  bool   `json:"success"`
	} `json:"entries"`
}

// Run executes all .hurl files in casesDir and returns a RunResult.
// Infrastructure failures (binary not found, no files) are returned as error.
// Test failures are expressed through result.Failed > 0, not as an error.
func (r *HurlRunner) Run(casesDir string, vars map[string]string) (RunResult, error) {
	if _, err := exec.LookPath("hurl"); err != nil {
		return RunResult{}, fmt.Errorf("hurl not found on PATH — run `caseforge doctor` to check dependencies")
	}

	files, err := filepath.Glob(filepath.Join(casesDir, "*.hurl"))
	if err != nil || len(files) == 0 {
		return RunResult{}, fmt.Errorf("no .hurl files found in %s", casesDir)
	}

	reportDir, err := os.MkdirTemp("", "caseforge-report-*")
	if err != nil {
		return RunResult{}, fmt.Errorf("creating report dir: %w", err)
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

	reportPath := filepath.Join(reportDir, "report.json")
	data, readErr := os.ReadFile(reportPath)
	if readErr == nil {
		result := buildRunResult(data)
		if result.Passed+result.Failed == 0 && runErr != nil {
			return RunResult{}, fmt.Errorf("hurl infrastructure error: %w", runErr)
		}
		return result, nil
	}

	// Fallback: no report file — use exit code
	if runErr != nil {
		return RunResult{Failed: len(files)}, nil
	}
	return RunResult{Passed: len(files)}, nil
}

func buildRunResult(data []byte) RunResult {
	var report hurlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return RunResult{}
	}
	result := RunResult{}
	for _, e := range report.Entries {
		id := strings.TrimSuffix(filepath.Base(e.Filename), ".hurl")
		result.Cases = append(result.Cases, CaseResult{ID: id, Title: id, Passed: e.Success})
		if e.Success {
			result.Passed++
		} else {
			result.Failed++
		}
	}
	return result
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
