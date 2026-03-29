// internal/output/schema/model.go
package schema

import "time"

const SchemaBaseURL = "https://caseforge.dev/schema/v1/testcase.json"

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
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	// Type is "test" for single-step cases; "setup"/"test"/"teardown" for chain cases.
	Type       string            `json:"type"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Headers    map[string]string `json:"headers,omitempty"`
	// Body is any for flexible serialization. On JSON unmarshal, becomes map[string]interface{}.
	Body       any               `json:"body,omitempty"`
	Assertions []Assertion       `json:"assertions"`
	Captures   []Capture         `json:"captures,omitempty"`   // populated in chain setup steps
	DependsOn  []string          `json:"depends_on,omitempty"` // populated in chain test/teardown steps
}

type Assertion struct {
	Target   string `json:"target"`   // "status_code"|"jsonpath $.<field>"|"header <HeaderName>"|"duration_ms"
	Operator string `json:"operator"` // "eq"|"ne"|"lt"|"gt"|"contains"|"matches"
	Expected any    `json:"expected"`
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
