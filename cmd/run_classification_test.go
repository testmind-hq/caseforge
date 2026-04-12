package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/runner"
)

func TestClassifyFailure_SecurityRegression(t *testing.T) {
	tc := schema.TestCase{
		Source: schema.CaseSource{Technique: "owasp_api_top10"},
		Steps: []schema.Step{{
			Assertions: []schema.Assertion{{Target: "status_code", Operator: "eq", Expected: 403}},
		}},
	}
	assert.Equal(t, runner.CategorySecurityRegression, classifyFailure(tc))
}

func TestClassifyFailure_MissingValidation(t *testing.T) {
	tc := schema.TestCase{
		Source: schema.CaseSource{Technique: "mutation"},
		Steps: []schema.Step{{
			Assertions: []schema.Assertion{{Target: "status_code", Operator: "gte", Expected: 400}},
		}},
	}
	assert.Equal(t, runner.CategoryMissingValidation, classifyFailure(tc))
}

func TestClassifyFailure_AuthFailure(t *testing.T) {
	tc := schema.TestCase{
		Source: schema.CaseSource{Technique: "auth_chain"},
	}
	assert.Equal(t, runner.CategoryAuthFailure, classifyFailure(tc))
}

func TestClassifyFailure_ServerError(t *testing.T) {
	tc := schema.TestCase{
		Source: schema.CaseSource{Technique: "equivalence_partitioning"},
		Steps: []schema.Step{{
			Assertions: []schema.Assertion{{Target: "status_code", Operator: "gte", Expected: 200}},
		}},
	}
	assert.Equal(t, runner.CategoryServerError, classifyFailure(tc))
}
