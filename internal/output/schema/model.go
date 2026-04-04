// internal/output/schema/model.go
package schema

import "time"

const SchemaBaseURL = "https://caseforge.dev/schema/v1/testcase.json"

// Assertion operator constants.
const (
	// OperatorEq checks exact equality.
	OperatorEq = "eq"
	// OperatorNe checks inequality.
	OperatorNe = "ne"
	// OperatorLt checks less-than (numeric).
	OperatorLt = "lt"
	// OperatorGt checks greater-than (numeric).
	OperatorGt = "gt"
	// OperatorContains checks substring / element containment.
	OperatorContains = "contains"
	// OperatorMatches checks regex match.
	OperatorMatches = "matches"
	// OperatorExists checks that the target field is present in the response.
	// Expected is not evaluated; use nil or omit it.
	OperatorExists = "exists"
	// OperatorIsISO8601 checks that the target field value is a valid ISO 8601
	// date-time string (e.g. "2024-01-02T15:04:05Z"). Expected is not evaluated.
	OperatorIsISO8601 = "is_iso8601"
	// OperatorIsUUID checks that the target field value is a valid UUID v4 string.
	// Expected is not evaluated.
	OperatorIsUUID = "is_uuid"
)

type TestCase struct {
	Schema      string            `json:"$schema"`
	Version     string            `json:"version"`
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Kind        string            `json:"kind"`     // "single"|"chain"
	Priority    string            `json:"priority"` // "P0"|"P1"|"P2"|"P3"
	Tags        []string          `json:"tags"`
	Source      CaseSource        `json:"source"`
	Steps       []Step            `json:"steps"`
	Labels      map[string]string `json:"labels,omitempty"`
	GeneratedAt time.Time         `json:"generated_at"`
}

type Step struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Type is "test" for single-step cases; "setup"/"test"/"teardown" for chain cases.
	Type       string            `json:"type"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Headers    map[string]string `json:"headers,omitempty"`
	// Body is any for flexible serialization. On JSON unmarshal, becomes map[string]interface{}.
	Body       any         `json:"body,omitempty"`
	Assertions []Assertion `json:"assertions"`
	Captures   []Capture   `json:"captures,omitempty"`   // populated in chain setup steps
	DependsOn  []string    `json:"depends_on,omitempty"` // populated in chain test/teardown steps
}

// Assertion is a single check applied to a step's response.
// Valid operators: eq, ne, lt, gt, contains, matches, exists, is_iso8601, is_uuid.
// For exists, is_iso8601, and is_uuid the Expected field is not evaluated by runners.
type Assertion struct {
	Target   string `json:"target"`   // "status_code"|"jsonpath $.<field>"|"header <Name>"|"duration_ms"|"body.<field>"
	Operator string `json:"operator"` // see Operator* constants
	Expected any    `json:"expected,omitempty"`
}

type CaseSource struct {
	Technique string `json:"technique"` // e.g. "equivalence_partitioning"
	SpecPath  string `json:"spec_path"` // e.g. "/users POST requestBody.properties.email"
	Rationale string `json:"rationale"`
}

// Capture records a value extracted from a step response for use in subsequent steps.
type Capture struct {
	Name   string `json:"name"`
	From   string `json:"from"`
	Filter string `json:"filter,omitempty"`
}
