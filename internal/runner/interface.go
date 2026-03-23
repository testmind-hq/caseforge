// internal/runner/interface.go
package runner

// RunResult holds the outcome of a test run.
type RunResult struct {
	Passed int
	Failed int
	Cases  []CaseResult // per-case results (empty if runner doesn't support it)
}

// CaseResult holds the pass/fail result of a single test case.
type CaseResult struct {
	ID     string // test case ID (e.g. "TC-0001")
	Title  string // test case title or filename
	Passed bool
}

// Runner executes test cases from an output directory.
type Runner interface {
	Run(casesDir string, vars map[string]string) (RunResult, error)
}
