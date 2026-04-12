// internal/runner/interface.go
package runner

// FailureCategory describes why a test case failed.
type FailureCategory = string

const (
	CategoryServerError        FailureCategory = "server_error"
	CategoryMissingValidation  FailureCategory = "missing_validation"
	CategoryAuthFailure        FailureCategory = "auth_failure"
	CategorySecurityRegression FailureCategory = "security_regression"
)

// RunResult holds the outcome of a test run.
type RunResult struct {
	Passed int
	Failed int
	Cases  []CaseResult // per-case results (empty if runner doesn't support it)
}

// CaseResult holds the pass/fail result of a single test case.
type CaseResult struct {
	ID       string          // test case ID (e.g. "TC-0001")
	Title    string          // test case title or filename
	Passed   bool
	Category FailureCategory // set for failed cases; empty for passed cases
}

// Runner executes test cases from an output directory.
type Runner interface {
	Run(casesDir string, vars map[string]string) (RunResult, error)
}
