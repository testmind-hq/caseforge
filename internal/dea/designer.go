// internal/dea/designer.go
package dea

import (
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// DesignProbe builds a concrete HTTP probe from a hypothesis node.
// It starts from a valid baseline body and mutates the target field
// according to the hypothesis kind.
func DesignProbe(h *HypothesisNode, op *spec.Operation, gen *datagen.Generator) Probe {
	headers := map[string]string{}
	var bodyAny any
	var queryParams map[string]string

	bodyFieldName := extractBodyFieldName(h.FieldPath)
	queryFieldName := extractQueryParamName(h.FieldPath)

	if hasJSONBody(op) {
		body := buildBaseBody(op, gen)
		headers["Content-Type"] = "application/json"
		if h.Kind == KindMassAssignment {
			body["admin"] = true
			body["role"] = "__probe__"
			body["isAdmin"] = true
		} else if bodyFieldName != "" {
			mutateField(body, bodyFieldName, h.Kind, op)
		}
		bodyAny = body
	}

	if queryFieldName != "" {
		queryParams = buildBaseQueryParams(op, gen)
		mutateQueryParam(queryParams, queryFieldName, h.Kind, op)
	}

	return Probe{
		Method:         op.Method,
		Path:           op.Path,
		Headers:        headers,
		Body:           bodyAny,
		QueryParams:    queryParams,
		ExpectedStatus: expectedStatusFor(h.Kind, op),
	}
}

func mutateField(body map[string]any, fieldName string, kind HypothesisKind, op *spec.Operation) {
	s := fieldSchema(op, fieldName)

	switch kind {
	case KindRequiredField, KindOptionalField:
		delete(body, fieldName)

	case KindNullValue:
		body[fieldName] = nil

	case KindStringMinLength:
		body[fieldName] = ""

	case KindStringMaxLength:
		if s != nil && s.MaxLength != nil {
			body[fieldName] = strings.Repeat("a", int(*s.MaxLength)+1)
		} else {
			body[fieldName] = strings.Repeat("a", 257)
		}

	case KindStringImplicitMin:
		body[fieldName] = ""

	case KindStringImplicitMax:
		body[fieldName] = strings.Repeat("a", 256)

	case KindNumericMin:
		if s != nil && s.Minimum != nil {
			body[fieldName] = *s.Minimum - 1
		} else {
			body[fieldName] = float64(-1)
		}

	case KindNumericMax:
		if s != nil && s.Maximum != nil {
			body[fieldName] = *s.Maximum + 1
		} else {
			body[fieldName] = float64(99999)
		}

	case KindEnumViolation:
		body[fieldName] = "__INVALID__ENUM_VALUE"

	case KindFormatViolation:
		body[fieldName] = "not-a-valid-format-value"

	case KindArrayMinItems:
		body[fieldName] = []any{}

	case KindArrayMaxItems:
		if s != nil && s.MaxItems != nil {
			arr := make([]any, int(*s.MaxItems)+1)
			for i := range arr {
				arr[i] = "item"
			}
			body[fieldName] = arr
		} else {
			body[fieldName] = make([]any, 101)
		}

	case KindTypeCoercion:
		if s != nil {
			switch s.Type {
			case "string":
				body[fieldName] = float64(123)
			case "integer", "number":
				body[fieldName] = "not_a_number"
			case "boolean":
				body[fieldName] = "not_a_boolean"
			case "array":
				body[fieldName] = "not_an_array"
			default:
				body[fieldName] = float64(999)
			}
		} else {
			body[fieldName] = float64(999)
		}

	case KindUnicodeControl:
		body[fieldName] = "hello\u0000world"
	}
}

func expectedStatusFor(kind HypothesisKind, op *spec.Operation) int {
	switch kind {
	case KindOptionalField:
		for code := range op.Responses {
			var n int
			if _, err := fmt.Sscanf(code, "%d", &n); err == nil && n >= 200 && n < 300 {
				return n
			}
		}
		return 200
	case KindTypeCoercion, KindUnicodeControl:
		return 400
	case KindMassAssignment:
		return 200
	}
	return 400
}

func hasJSONBody(op *spec.Operation) bool {
	if op.RequestBody == nil {
		return false
	}
	_, ok := op.RequestBody.Content["application/json"]
	return ok
}

func buildBaseBody(op *spec.Operation, gen *datagen.Generator) map[string]any {
	mt := op.RequestBody.Content["application/json"]
	if mt == nil || mt.Schema == nil {
		return map[string]any{}
	}
	body := map[string]any{}
	for name, fs := range mt.Schema.Properties {
		body[name] = gen.Generate(fs, name)
	}
	return body
}

func extractBodyFieldName(fieldPath string) string {
	const prefix = "requestBody."
	if strings.HasPrefix(fieldPath, prefix) {
		return strings.TrimPrefix(fieldPath, prefix)
	}
	return ""
}

func fieldSchema(op *spec.Operation, fieldName string) *spec.Schema {
	if op.RequestBody == nil {
		return nil
	}
	mt, ok := op.RequestBody.Content["application/json"]
	if !ok || mt.Schema == nil {
		return nil
	}
	return mt.Schema.Properties[fieldName]
}

func extractQueryParamName(fieldPath string) string {
	const prefix = "query."
	if strings.HasPrefix(fieldPath, prefix) {
		return strings.TrimPrefix(fieldPath, prefix)
	}
	return ""
}

func buildBaseQueryParams(op *spec.Operation, gen *datagen.Generator) map[string]string {
	params := map[string]string{}
	for _, p := range op.Parameters {
		if p.In != "query" || p.Schema == nil {
			continue
		}
		val := gen.Generate(p.Schema, p.Name)
		params[p.Name] = fmt.Sprintf("%v", val)
	}
	return params
}

func queryParamSchema(op *spec.Operation, paramName string) *spec.Schema {
	for _, p := range op.Parameters {
		if p.Name == paramName && p.In == "query" {
			return p.Schema
		}
	}
	return nil
}

func mutateQueryParam(params map[string]string, paramName string, kind HypothesisKind, op *spec.Operation) {
	s := queryParamSchema(op, paramName)
	switch kind {
	case KindRequiredQueryParam:
		delete(params, paramName)

	case KindNumericMin:
		if s != nil && s.Minimum != nil {
			params[paramName] = fmt.Sprintf("%g", *s.Minimum-1)
		} else {
			params[paramName] = "-1"
		}
	case KindNumericMax:
		if s != nil && s.Maximum != nil {
			params[paramName] = fmt.Sprintf("%g", *s.Maximum+1)
		} else {
			params[paramName] = "99999"
		}
	}
}
