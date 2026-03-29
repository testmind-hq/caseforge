// internal/dea/hypothesis.go
package dea

import "time"

type HypothesisKind string

const (
	KindRequiredField     HypothesisKind = "required_field"
	KindOptionalField     HypothesisKind = "optional_field"
	KindStringMinLength   HypothesisKind = "string_min_length"
	KindStringMaxLength   HypothesisKind = "string_max_length"
	KindStringImplicitMax HypothesisKind = "string_implicit_max"
	KindStringImplicitMin HypothesisKind = "string_implicit_min"
	KindNumericMin        HypothesisKind = "numeric_min"
	KindNumericMax        HypothesisKind = "numeric_max"
	KindNullValue         HypothesisKind = "null_value"
	KindEnumViolation     HypothesisKind = "enum_violation"
)

type HypothesisStatus string

const (
	StatusPending   HypothesisStatus = "pending"
	StatusConfirmed HypothesisStatus = "confirmed"
	StatusRefuted   HypothesisStatus = "refuted"
)

// HypothesisNode is one node in the hypothesis tree.
// Each node is a falsifiable claim about an API operation.
type HypothesisNode struct {
	ID          string
	Kind        HypothesisKind
	Description string

	// What operation + field this targets
	Operation string // e.g. "POST /pets"
	FieldPath string // e.g. "requestBody.name"

	// The probe to run to test this hypothesis
	Probe Probe

	// Resolved after probe execution
	Status   HypothesisStatus
	Evidence *Evidence

	// Child hypotheses (generated after this one resolves)
	Children []*HypothesisNode
}

// Probe is a concrete HTTP request to be executed against the target API.
type Probe struct {
	Method      string
	Path        string
	Headers     map[string]string
	Body        any
	QueryParams map[string]string

	// Expected HTTP status if this hypothesis is TRUE
	ExpectedStatus int
}

// Evidence is the observed HTTP response from running a probe.
type Evidence struct {
	ActualStatus  int
	ActualBody    string
	ActualHeaders map[string]string
	Duration      time.Duration
}

// Resolve marks the hypothesis as confirmed or refuted and stores the evidence.
func (h *HypothesisNode) Resolve(ev *Evidence, confirmed bool) {
	h.Evidence = ev
	if confirmed {
		h.Status = StatusConfirmed
	} else {
		h.Status = StatusRefuted
	}
}
