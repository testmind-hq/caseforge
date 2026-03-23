// internal/runner/k6.go
package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var _ Runner = (*K6Runner)(nil)

// K6Runner executes k6 test scripts via the k6 binary.
type K6Runner struct {
	k6Bin string // path to k6 binary; defaults to "k6" (from PATH)
}

func NewK6Runner() *K6Runner {
	return &K6Runner{k6Bin: "k6"}
}

func (r *K6Runner) Run(casesDir string, vars map[string]string) (RunResult, error) {
	k6Bin := r.k6Bin
	if _, err := exec.LookPath(k6Bin); err != nil {
		return RunResult{}, fmt.Errorf("k6 binary not found in PATH: %w", err)
	}

	scriptPath := filepath.Join(casesDir, "k6_tests.js")
	if _, err := os.Stat(scriptPath); err != nil {
		return RunResult{}, fmt.Errorf("k6_tests.js not found in %s: %w", casesDir, err)
	}

	tmpDir, err := os.MkdirTemp("", "caseforge-k6-*")
	if err != nil {
		return RunResult{}, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	summaryPath := filepath.Join(tmpDir, "summary.json")

	args := []string{"run", "--summary-export=" + summaryPath}
	for k, v := range vars {
		args = append(args, "--env", k+"="+v)
	}
	args = append(args, scriptPath)

	cmd := exec.Command(k6Bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// k6 exits non-zero when checks fail — that's expected, not an infrastructure error.
	// We intentionally ignore the exit code here and detect failures via the summary file.
	_ = cmd.Run()

	// If the summary file was not written, k6 crashed before finishing (infrastructure failure).
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		return RunResult{}, fmt.Errorf("k6 did not produce a summary (crashed or bad script): %w", err)
	}
	return parseK6Summary(data)
}

type k6Summary struct {
	Metrics struct {
		Checks *struct {
			Passes int `json:"passes"`
			Fails  int `json:"fails"`
		} `json:"checks"`
	} `json:"metrics"`
}

func parseK6Summary(data []byte) (RunResult, error) {
	var s k6Summary
	if err := json.Unmarshal(data, &s); err != nil {
		return RunResult{}, fmt.Errorf("parsing k6 summary: %w", err)
	}
	if s.Metrics.Checks == nil {
		return RunResult{}, nil
	}
	return RunResult{
		Passed: s.Metrics.Checks.Passes,
		Failed: s.Metrics.Checks.Fails,
	}, nil
}
