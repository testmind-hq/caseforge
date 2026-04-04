// internal/methodology/engine_annotation_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSemanticAnnotation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
	}{
		{
			name: "valid JSON",
			input: `{"resource_type":"user","action_type":"create","has_state_machine":false,"state_field":"","unique_fields":[],"implicit_rules":[]}`,
		},
		{
			name: "JSON in markdown code fence",
			input: "```json\n{\"resource_type\":\"order\",\"action_type\":\"read\",\"has_state_machine\":true,\"state_field\":\"status\",\"unique_fields\":[],\"implicit_rules\":[]}\n```",
		},
		{
			name:  "json embedded in prose",
			input: "The annotation is: {\"resource_type\":\"user\",\"action_type\":\"create\",\"has_state_machine\":false,\"state_field\":\"\",\"unique_fields\":[],\"implicit_rules\":[]}",
		},
		{
			name:    "invalid text",
			input:   "this is not json at all",
			wantNil: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSemanticAnnotation(tt.input)
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
			}
		})
	}
}
