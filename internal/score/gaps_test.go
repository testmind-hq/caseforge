package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestComputeGaps_AllCovered(t *testing.T) {
	cases := []schema.TestCase{
		{
			Source: schema.CaseSource{SpecPath: "GET /users"},
			Steps: []schema.Step{{Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 200},
				{Target: "status_code", Operator: "lt", Expected: 300},
			}}},
		},
		{
			Source: schema.CaseSource{SpecPath: "GET /users"},
			Steps: []schema.Step{{Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 400},
			}}},
		},
	}
	gaps := ComputeGaps(cases)
	assert.Empty(t, gaps)
}

func TestComputeGaps_Missing4xx(t *testing.T) {
	cases := []schema.TestCase{
		{
			Source: schema.CaseSource{SpecPath: "POST /users"},
			Steps: []schema.Step{{Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "lt", Expected: 300},
			}}},
		},
	}
	gaps := ComputeGaps(cases)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "POST", gaps[0].Method)
	assert.Equal(t, "/users", gaps[0].Path)
	assert.True(t, gaps[0].Has2xx)
	assert.False(t, gaps[0].Has4xx)
}

func TestComputeGaps_Missing2xx(t *testing.T) {
	cases := []schema.TestCase{
		{
			Source: schema.CaseSource{SpecPath: "POST /items"},
			Steps: []schema.Step{{Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "gte", Expected: 400},
			}}},
		},
	}
	gaps := ComputeGaps(cases)
	assert.Len(t, gaps, 1)
	assert.Equal(t, "POST", gaps[0].Method)
	assert.Equal(t, "/items", gaps[0].Path)
	assert.False(t, gaps[0].Has2xx)
	assert.True(t, gaps[0].Has4xx)
}

func TestComputeGaps_BothMissing(t *testing.T) {
	cases := []schema.TestCase{
		{
			Source: schema.CaseSource{SpecPath: "DELETE /items/{id}"},
			Steps: []schema.Step{{Assertions: []schema.Assertion{
				{Target: "body.name", Operator: "eq", Expected: "test"},
			}}},
		},
	}
	gaps := ComputeGaps(cases)
	assert.Len(t, gaps, 1)
	assert.False(t, gaps[0].Has2xx)
	assert.False(t, gaps[0].Has4xx)
}
