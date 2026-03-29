// internal/dea/seeder.go
package dea

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// SeedHypotheses generates the initial set of hypothesis nodes for an operation.
// Each node represents one falsifiable claim to be tested.
func SeedHypotheses(op *spec.Operation) []*HypothesisNode {
	var nodes []*HypothesisNode
	nodes = append(nodes, seedBodyHypotheses(op)...)
	nodes = append(nodes, seedQueryParamHypotheses(op)...)
	return nodes
}

func seedBodyHypotheses(op *spec.Operation) []*HypothesisNode {
	if op.RequestBody == nil {
		return nil
	}
	mt, ok := op.RequestBody.Content["application/json"]
	if !ok || mt.Schema == nil {
		return nil
	}
	s := mt.Schema
	var nodes []*HypothesisNode

	requiredSet := make(map[string]bool)
	for _, r := range s.Required {
		requiredSet[r] = true
	}

	for fieldName, fieldSchema := range s.Properties {
		prefix := fmt.Sprintf("requestBody.%s", fieldName)

		if requiredSet[fieldName] {
			nodes = append(nodes, newHypothesis(op, KindRequiredField, prefix,
				fmt.Sprintf("omitting required field '%s' must return 4xx", fieldName)))
			nodes = append(nodes, newHypothesis(op, KindNullValue, prefix,
				fmt.Sprintf("null value for required field '%s' must return 4xx", fieldName)))
		} else {
			nodes = append(nodes, newHypothesis(op, KindOptionalField, prefix,
				fmt.Sprintf("optional field '%s' can be omitted without error", fieldName)))
		}

		nodes = append(nodes, seedFieldConstraintHypotheses(op, fieldName, fieldSchema, prefix)...)
	}
	return nodes
}

func seedFieldConstraintHypotheses(op *spec.Operation, fieldName string, s *spec.Schema, prefix string) []*HypothesisNode {
	var nodes []*HypothesisNode
	switch s.Type {
	case "string":
		if s.MinLength != nil {
			nodes = append(nodes, newHypothesis(op, KindStringMinLength, prefix,
				fmt.Sprintf("string field '%s' with length < minLength(%d) must return 4xx", fieldName, *s.MinLength)))
		} else {
			nodes = append(nodes, newHypothesis(op, KindStringImplicitMin, prefix,
				fmt.Sprintf("empty string for field '%s' may be rejected (undeclared min)", fieldName)))
		}
		if s.MaxLength != nil {
			nodes = append(nodes, newHypothesis(op, KindStringMaxLength, prefix,
				fmt.Sprintf("string field '%s' with length > maxLength(%d) must return 4xx", fieldName, *s.MaxLength)))
		} else {
			nodes = append(nodes, newHypothesis(op, KindStringImplicitMax, prefix,
				fmt.Sprintf("very long string for field '%s' may be rejected (undeclared max)", fieldName)))
		}
		if len(s.Enum) > 0 {
			nodes = append(nodes, newHypothesis(op, KindEnumViolation, prefix,
				fmt.Sprintf("invalid enum value for field '%s' must return 4xx", fieldName)))
		}
	case "integer", "number":
		if s.Minimum != nil {
			nodes = append(nodes, newHypothesis(op, KindNumericMin, prefix,
				fmt.Sprintf("numeric field '%s' below minimum(%.0f) must return 4xx", fieldName, *s.Minimum)))
		}
		if s.Maximum != nil {
			nodes = append(nodes, newHypothesis(op, KindNumericMax, prefix,
				fmt.Sprintf("numeric field '%s' above maximum(%.0f) must return 4xx", fieldName, *s.Maximum)))
		}
	}
	return nodes
}

func seedQueryParamHypotheses(op *spec.Operation) []*HypothesisNode {
	var nodes []*HypothesisNode
	for _, p := range op.Parameters {
		if p.In != "query" || p.Schema == nil {
			continue
		}
		prefix := fmt.Sprintf("query.%s", p.Name)
		if p.Schema.Type == "integer" || p.Schema.Type == "number" {
			if p.Schema.Minimum != nil {
				nodes = append(nodes, newHypothesis(op, KindNumericMin, prefix,
					fmt.Sprintf("query param '%s' below minimum must return 4xx", p.Name)))
			}
			if p.Schema.Maximum != nil {
				nodes = append(nodes, newHypothesis(op, KindNumericMax, prefix,
					fmt.Sprintf("query param '%s' above maximum must return 4xx", p.Name)))
			}
		}
	}
	return nodes
}

func newHypothesis(op *spec.Operation, kind HypothesisKind, fieldPath, desc string) *HypothesisNode {
	return NewHypothesisNode(
		fmt.Sprintf("H-%s", strings.ToUpper(uuid.New().String()[:6])),
		kind,
		fmt.Sprintf("%s %s", op.Method, op.Path),
		fieldPath,
		desc,
	)
}
