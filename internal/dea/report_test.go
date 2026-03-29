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

func TestDiscoveredRuleJSONTags(t *testing.T) {
	rule := DiscoveredRule{
		ID:         "RULE-001",
		Implicit:   true,
		Confidence: ConfidenceHigh,
		Category:   CategorySpecMismatch,
	}
	data, err := json.Marshal(rule)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"implicit":true`)
	assert.Contains(t, string(data), `"confidence":"high"`)
	assert.Contains(t, string(data), `"category":"spec_mismatch"`)
}
