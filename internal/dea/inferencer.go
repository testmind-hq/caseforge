// internal/dea/inferencer.go
package dea

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// InferRule converts a resolved hypothesis to a DiscoveredRule.
// Returns nil when no noteworthy rule is produced (e.g. refuted implicit probe).
// Panics if called on a pending hypothesis.
func InferRule(h *HypothesisNode) *DiscoveredRule {
	if h.Status == StatusPending {
		panic("InferRule called on pending hypothesis: resolve it first")
	}
	if h.Evidence == nil {
		panic("InferRule called on hypothesis with nil Evidence: run probe first")
	}

	switch h.Kind {
	case KindRequiredField:
		return inferRequiredField(h)
	case KindOptionalField:
		return inferOptionalField(h)
	case KindStringMaxLength:
		return inferSpecStringMax(h)
	case KindStringMinLength:
		return inferSpecStringMin(h)
	case KindStringImplicitMax:
		return inferImplicitStringMax(h)
	case KindStringImplicitMin:
		return inferImplicitStringMin(h)
	case KindNumericMin:
		return inferNumericMin(h)
	case KindNumericMax:
		return inferNumericMax(h)
	case KindNullValue:
		return inferNullValue(h)
	case KindEnumViolation:
		return inferEnumViolation(h)
	}
	return nil
}

func newRule(h *HypothesisNode, category RuleCategory, description string, implicit bool) *DiscoveredRule {
	return &DiscoveredRule{
		ID:          fmt.Sprintf("RULE-%s", strings.ToUpper(uuid.New().String()[:6])),
		Operation:   h.Operation,
		FieldPath:   h.FieldPath,
		Category:    category,
		Description: description,
		Implicit:    implicit,
		Confidence:  ConfidenceHigh,
		Evidence: []RuleEvidence{{
			ProbeDescription: h.Description,
			ActualStatus:     h.Evidence.ActualStatus,
			ActualBody:       h.Evidence.ActualBody,
		}},
	}
}

func inferRequiredField(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s' is required and server validates it (returns %d when omitted)",
				field, h.Evidence.ActualStatus),
			false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec declares required but server accepts omission (returned %d)",
			field, h.Evidence.ActualStatus),
		true)
}

func inferOptionalField(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s' is optional and confirmed by server", field),
			false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec marks optional but server requires it (returned %d when omitted)",
			field, h.Evidence.ActualStatus),
		true)
}

func inferSpecStringMax(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': spec-declared maxLength is enforced by server", field),
			false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec declares max but server does not enforce it (returned %d)",
			field, h.Evidence.ActualStatus),
		true)
}

func inferSpecStringMin(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': spec-declared minLength is enforced by server", field),
			false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec declares min but server does not enforce it (returned %d)",
			field, h.Evidence.ActualStatus),
		true)
}

func inferImplicitStringMax(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': undeclared max length — 256 chars causes %d (spec does not declare maxLength)",
				field, h.Evidence.ActualStatus),
			true)
	}
	return nil
}

func inferImplicitStringMin(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': undeclared min length — empty string causes %d (spec does not declare minLength)",
				field, h.Evidence.ActualStatus),
			true)
	}
	return nil
}

func inferNumericMin(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': minimum boundary enforced by server", field), false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec declares minimum but server does not enforce it", field), true)
}

func inferNumericMax(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': maximum boundary enforced by server", field), false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec declares maximum but server does not enforce it", field), true)
}

func inferNullValue(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': null value rejected by server (%d)", field, h.Evidence.ActualStatus),
			false)
	}
	return newRule(h, CategoryBehavior,
		fmt.Sprintf("Field '%s': server accepts null value — consider adding nullable: true to spec", field),
		true)
}

func inferEnumViolation(h *HypothesisNode) *DiscoveredRule {
	field := extractFieldName(h.FieldPath)
	if h.Status == StatusConfirmed {
		return newRule(h, CategoryFieldConstraint,
			fmt.Sprintf("Field '%s': enum validation enforced by server (%d on invalid value)",
				field, h.Evidence.ActualStatus),
			false)
	}
	return newRule(h, CategorySpecMismatch,
		fmt.Sprintf("Field '%s': spec declares enum but server accepts any value", field),
		true)
}

func extractFieldName(fieldPath string) string {
	if fieldPath == "" {
		return ""
	}
	parts := strings.Split(fieldPath, ".")
	return parts[len(parts)-1]
}
