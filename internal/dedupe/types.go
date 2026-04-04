// internal/dedupe/types.go
package dedupe

import "time"

// LoadedCase is a parsed TestCase file with its filesystem path.
type LoadedCase struct {
	FilePath string
	TC       TestCaseSnapshot
}

// TestCaseSnapshot holds only the fields needed for duplicate detection.
type TestCaseSnapshot struct {
	Method           string
	Path             string
	ExpectedStatus   int
	BodyJSON         string   // JSON-normalized body of the first "test" step; "" if absent
	AssertionTargets []string // deduplicated Assertion.Target values across all steps
}

// MatchKind distinguishes exact vs structural duplicate groups.
type MatchKind string

const (
	MatchExact      MatchKind = "exact"
	MatchStructural MatchKind = "structural"
)

// CaseScore holds the score and keep decision for one file within a duplicate group.
type CaseScore struct {
	FilePath       string `json:"file_path"`
	AssertionCount int    `json:"assertion_count"`
	Keep           bool   `json:"keep"`
}

// DuplicateGroup represents two or more files that are duplicates of each other.
type DuplicateGroup struct {
	Kind       MatchKind   `json:"kind"`
	Similarity float64     `json:"similarity"`
	Cases      []CaseScore `json:"cases"`
}

// DedupeReport is the top-level result produced by a dedupe run.
type DedupeReport struct {
	CasesDir         string           `json:"cases_dir"`
	TotalScanned     int              `json:"total_scanned"`
	ExactGroups      int              `json:"exact_groups"`
	StructuralGroups int              `json:"structural_groups"`
	Groups           []DuplicateGroup `json:"groups"`
	GeneratedAt      time.Time        `json:"generated_at"`
}
