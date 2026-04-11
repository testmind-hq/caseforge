package dea


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
	KindArrayMinItems     HypothesisKind = "array_min_items"
	KindArrayMaxItems     HypothesisKind = "array_max_items"
	KindRequiredQueryParam HypothesisKind = "required_query_param"
	KindFormatViolation   HypothesisKind = "format_violation"
	KindTypeCoercion      HypothesisKind = "type_coercion"
	KindUnicodeControl    HypothesisKind = "unicode_control"
	KindMassAssignment    HypothesisKind = "mass_assignment"
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
	ID          string           `json:"id"`
	Kind        HypothesisKind   `json:"kind"`
	Description string           `json:"description"`
	Operation   string           `json:"operation"`
	FieldPath   string           `json:"field_path"`
	Probe       Probe            `json:"probe"`
	Status      HypothesisStatus `json:"status"`
	Evidence    *Evidence        `json:"evidence,omitempty"`
	Children    []*HypothesisNode `json:"children,omitempty"`
}

// Probe is a concrete HTTP request to be executed against the target API.
type Probe struct {
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	Headers        map[string]string `json:"headers,omitempty"`
	Body           any               `json:"body,omitempty"`
	QueryParams    map[string]string `json:"query_params,omitempty"`
	ExpectedStatus int               `json:"expected_status"`
}

// Evidence is the observed HTTP response from running a probe.
type Evidence struct {
	ActualStatus  int               `json:"actual_status"`
	ActualBody    string            `json:"actual_body,omitempty"`
	ActualHeaders map[string]string `json:"actual_headers,omitempty"`
	DurationMs    int64             `json:"duration_ms"`
}

// NewHypothesisNode creates a HypothesisNode with StatusPending as the initial status.
func NewHypothesisNode(id string, kind HypothesisKind, operation, fieldPath, description string) *HypothesisNode {
	return &HypothesisNode{
		ID:          id,
		Kind:        kind,
		Operation:   operation,
		FieldPath:   fieldPath,
		Description: description,
		Status:      StatusPending,
	}
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
