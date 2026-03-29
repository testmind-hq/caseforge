package dea

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExplorationReportJSON(t *testing.T) {
	r := &ExplorationReport{
		SpecPath:  "petstore.yaml",
		TargetURL: "http://localhost:8080",
		Rules: []DiscoveredRule{
			{
				ID:          "RULE-001",
				Operation:   "POST /pets",
				Category:    CategoryFieldConstraint,
				Description: "Field 'name' is required by server",
				Implicit:    false,
				Confidence:  ConfidenceHigh,
			},
		},
	}
	data, err := json.Marshal(r)
	require.NoError(t, err)
	var decoded ExplorationReport
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, "RULE-001", decoded.Rules[0].ID)
}

func TestDiscoveredRuleImplicitFlag(t *testing.T) {
	implicit := DiscoveredRule{Implicit: true, Confidence: ConfidenceHigh}
	explicit := DiscoveredRule{Implicit: false, Confidence: ConfidenceMedium}
	assert.True(t, implicit.Implicit)
	assert.False(t, explicit.Implicit)
}
