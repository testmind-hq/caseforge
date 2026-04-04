package export_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/export"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

func TestNewExporter_ValidFormats(t *testing.T) {
	for _, format := range []string{"allure", "xray", "testrail"} {
		exp, err := export.New(format)
		require.NoError(t, err, "format %s", format)
		assert.Equal(t, format, exp.Format())
	}
}

func TestNewExporter_InvalidFormat(t *testing.T) {
	_, err := export.New("unknown")
	assert.ErrorContains(t, err, "unknown export format")
}

func TestPriorityLabel_AllValues(t *testing.T) {
	cases := []struct{ in, allure, xray, testrail string }{
		{"P0", "blocker", "Highest", "Critical"},
		{"P1", "critical", "High", "High"},
		{"P2", "normal", "Medium", "Medium"},
		{"P3", "minor", "Low", "Low"},
		{"", "normal", "Medium", "Medium"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.allure, export.PriorityAllure(tc.in), "allure %s", tc.in)
		assert.Equal(t, tc.xray, export.PriorityXray(tc.in), "xray %s", tc.in)
		assert.Equal(t, tc.testrail, export.PriorityTestRail(tc.in), "testrail %s", tc.in)
	}
}

func TestAssertionsSummary(t *testing.T) {
	a := []schema.Assertion{
		{Target: "status_code", Operator: "eq", Expected: 201},
		{Target: "duration_ms", Operator: "lt", Expected: 2000},
	}
	got := export.AssertionsSummary(a)
	assert.Equal(t, "status_code eq 201; duration_ms lt 2000", got)
}

func TestAssertionsSummary_Empty(t *testing.T) {
	assert.Equal(t, "", export.AssertionsSummary(nil))
}
