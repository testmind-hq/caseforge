// internal/mutation/types.go
package mutation

import (
	"net/http"
	"time"
)

// Operator applies a single mutation to an HTTP response body.
// Header-mutating operators set headers directly on resp inside Apply.
type Operator interface {
	Name() string
	Apply(resp *http.Response, body []byte) ([]byte, error)
}

// CaseMutationResult records whether one operator survived against one test case.
type CaseMutationResult struct {
	CaseID    string `json:"case_id"`
	Title     string `json:"title"`
	Operation string `json:"operation"` // from source.spec_path; "(unknown)" if absent
	Operator  string `json:"operator"`
	Survived  bool   `json:"survived"` // true = mutation was not caught by assertions
}

// SurvivorCluster groups all operators that survived a single test case.
type SurvivorCluster struct {
	CaseID    string   `json:"case_id"`
	Title     string   `json:"title"`
	Operators []string `json:"operators"`
	RiskScore float64  `json:"risk_score"` // len(Operators) / total operators in run
}

// OperationScore holds the per-operation mutation score breakdown.
type OperationScore struct {
	Operation     string  `json:"operation"`      // "METHOD /path"
	TotalRuns     int     `json:"total_runs"`
	Killed        int     `json:"killed"`
	Survivors     int     `json:"survivors"`
	MutationScore float64 `json:"mutation_score"` // killed / total_runs
}

// SuggestedAssertion is one LLM-suggested assertion to add to an existing test case.
type SuggestedAssertion struct {
	Target   string `json:"target"`
	Operator string `json:"operator"`
	Expected any    `json:"expected,omitempty"`
}

// FeedbackItem is one LLM-generated diagnosis for a SurvivorCluster.
type FeedbackItem struct {
	CaseID              string               `json:"case_id"`
	Title               string               `json:"title"`
	RiskScore           float64              `json:"risk_score"`
	Diagnosis           string               `json:"diagnosis"`
	SuggestedAssertions []SuggestedAssertion `json:"suggested_assertions"`
}

// MutationRun is the top-level report for a complete mutation run.
type MutationRun struct {
	Target        string               `json:"target"`
	CasesDir      string               `json:"cases_dir"`
	Operators     []string             `json:"operators"`
	TotalCases    int                  `json:"total_cases"`
	TotalRuns     int                  `json:"total_runs"`    // cases × operators
	Killed        int                  `json:"killed"`
	Survivors     int                  `json:"survivors"`
	MutationScore float64              `json:"mutation_score"` // killed / total_runs
	Results         []CaseMutationResult `json:"results"`
	Clusters        []SurvivorCluster    `json:"clusters,omitempty"` // Phase 2
	OperationScores []OperationScore    `json:"operation_scores,omitempty"`
	Feedback        []FeedbackItem      `json:"feedback,omitempty"` // Phase 2
	GeneratedAt     string              `json:"generated_at"`
}

// RunOptions configures a mutation run.
type RunOptions struct {
	Target              string
	CasesDir            string
	OutputDir           string
	Operators           []Operator
	Concurrency         int // cases per operator (default 4)
	OperatorConcurrency int // operators in parallel (default 2)
}

func (o *RunOptions) setDefaults() {
	if o.Concurrency <= 0 {
		o.Concurrency = 4
	}
	if o.OperatorConcurrency <= 0 {
		o.OperatorConcurrency = 2
	}
}

// persistDir returns the directory for run JSON files.
func persistDir() string {
	return ".caseforge/mutation/runs"
}

// timestamp returns a UTC timestamp string suitable for filenames.
func timestamp() string {
	return time.Now().UTC().Format("2006-01-02T15-04-05Z")
}
