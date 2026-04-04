// internal/rbt/types.go
package rbt

import "time"

type ChangedFile struct {
	Path         string
	ChangedLines []int
	IsNew        bool
	IsDeleted    bool
}

type RouteMapping struct {
	SourceFile string
	Line       int
	Method     string  // "GET", "POST", ...
	RoutePath  string  // "/users/{id}"
	Via        string  // "mapfile"|"treesitter"|"regex"|"llm"
	Confidence float64 // 0.0–1.0; <0.5 → "uncertain"
}

type RiskLevel string

const (
	RiskNone      RiskLevel = "none"
	RiskLow       RiskLevel = "low"
	RiskMedium    RiskLevel = "medium"
	RiskHigh      RiskLevel = "high"
	RiskUncertain RiskLevel = "uncertain"
)

type TestCaseRef struct {
	File   string // e.g. "cases/POST_users_201.json"
	CaseID string
	Title  string
}

type OperationCoverage struct {
	OperationID string
	Method      string
	Path        string
	Affected    bool // touched by git diff
	SourceRefs  []RouteMapping
	TestCases   []TestCaseRef
	Risk        RiskLevel
}

type RiskReport struct {
	DiffBase       string
	DiffHead       string
	ChangedFiles   []ChangedFile
	Operations     []OperationCoverage
	TotalAffected  int
	TotalCovered   int
	TotalUncovered int
	RiskScore      float64   // uncovered/affected, 0.0–1.0
	GeneratedAt    time.Time
}
