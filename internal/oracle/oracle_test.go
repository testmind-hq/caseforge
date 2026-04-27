package oracle

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type callCountProvider struct {
	responses []string
	count     int
}

func (c *callCountProvider) Complete(_ context.Context, _ *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	idx := c.count
	c.count++
	if idx >= len(c.responses) {
		return &llm.CompletionResponse{Text: "[]"}, nil
	}
	return &llm.CompletionResponse{Text: c.responses[idx]}, nil
}
func (c *callCountProvider) IsAvailable() bool { return true }
func (c *callCountProvider) Name() string      { return "stub" }

func TestMine_ReturnsConfirmedConstraints(t *testing.T) {
	provider := &callCountProvider{
		responses: []string{
			`[{"type":"format","field":"email","detail":"must be valid email address"},{"type":"range","field":"age","detail":"must be >= 0"}]`,
			`[{"type":"format","field":"email","detail":"must be valid email address","confirmed":true},{"type":"range","field":"age","detail":"must be >= 0","confirmed":false}]`,
		},
	}
	op := &spec.Operation{
		Method: "POST", Path: "/users",
		Responses: map[string]*spec.Response{
			"201": {Description: "created", Content: map[string]*spec.MediaType{
				"application/json": {Schema: &spec.Schema{
					Type: "object",
					Properties: map[string]*spec.Schema{
						"email": {Type: "string", Format: "email"},
						"age":   {Type: "integer"},
					},
				}},
			}},
		},
	}
	constraints, err := Mine(context.Background(), op, provider)
	require.NoError(t, err)
	require.Len(t, constraints, 1)
	assert.Equal(t, "format", constraints[0].Type)
	assert.Equal(t, "email", constraints[0].Field)
}

func TestMine_NoopProvider_ReturnsEmpty(t *testing.T) {
	op := &spec.Operation{Method: "GET", Path: "/items"}
	constraints, err := Mine(context.Background(), op, &llm.NoopProvider{})
	require.NoError(t, err)
	assert.Empty(t, constraints)
}

func TestConstraintToAssertion_Exists(t *testing.T) {
	c := Constraint{Type: "exists", Field: "createdAt", Detail: "field must be present"}
	assertions := c.ToAssertions()
	require.Len(t, assertions, 1)
	assert.Equal(t, "jsonpath $.createdAt", assertions[0].Target)
	assert.Equal(t, schema.OperatorExists, assertions[0].Operator)
}

func TestConstraintToAssertion_Format_UUID(t *testing.T) {
	c := Constraint{Type: "format", Field: "id", Detail: "must be UUID"}
	assertions := c.ToAssertions()
	require.Len(t, assertions, 1)
	assert.Equal(t, "jsonpath $.id", assertions[0].Target)
}

func TestInjectIntoCase_Injects2xxCase(t *testing.T) {
	tc := schema.TestCase{
		Steps: []schema.Step{{
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 200},
				{Target: "status_code", Operator: "lt", Expected: 300},
			},
		}},
	}
	constraints := []Constraint{{Type: "exists", Field: "id", Detail: "must exist"}}
	result := InjectIntoCase(tc, constraints)
	assert.Greater(t, len(result.Steps[0].Assertions), 2)
}

func TestInjectIntoCase_Skips4xxCase(t *testing.T) {
	tc := schema.TestCase{
		Steps: []schema.Step{{
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 400},
			},
		}},
	}
	constraints := []Constraint{{Type: "exists", Field: "id", Detail: "must exist"}}
	result := InjectIntoCase(tc, constraints)
	assert.Len(t, result.Steps[0].Assertions, 1) // unchanged
}
