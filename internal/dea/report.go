package dea

import (
	"time"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

type RuleCategory string

const (
	CategoryFieldConstraint RuleCategory = "field_constraint"
	CategorySpecMismatch    RuleCategory = "spec_mismatch"
	CategoryBehavior        RuleCategory = "behavior"
)

type Confidence string

const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
	ConfidenceLow    Confidence = "low"
)

// DiscoveredRule is an inferred API constraint or behavior.
type DiscoveredRule struct {
	ID        string       `json:"id"`
	Operation string       `json:"operation"`  // e.g. "POST /pets"
	FieldPath string       `json:"field_path"` // e.g. "requestBody.name"
	Category  RuleCategory `json:"category"`
	// Description is a human-readable explanation of the rule.
	Description string `json:"description"`
	// Implicit: true = server enforces constraint not declared in spec.
	// Implicit: false = spec declares it and server confirms it.
	Implicit   bool           `json:"implicit"`
	Confidence Confidence     `json:"confidence"`
	Evidence   []RuleEvidence `json:"evidence"`
}

// RuleEvidence records one probe-response pair that supports the rule.
type RuleEvidence struct {
	ProbeDescription string `json:"probe_description"`
	ActualStatus     int    `json:"actual_status"`
	ActualBody       string `json:"actual_body,omitempty"`
	SpecDeclared     string `json:"spec_declared,omitempty"`
}

// ExplorationReport is the full output of one DEA run.
type ExplorationReport struct {
	SpecPath    string           `json:"spec_path"`
	TargetURL   string           `json:"target_url"`
	ExploredAt  time.Time        `json:"explored_at"`
	TotalProbes int              `json:"total_probes"`
	Rules       []DiscoveredRule `json:"rules"`
	// TestCases contains schema.TestCase entries for rules with Implicit=true.
	TestCases []schema.TestCase `json:"test_cases,omitempty"`
}
