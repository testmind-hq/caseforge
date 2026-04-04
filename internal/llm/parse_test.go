// internal/llm/parse_test.go
package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bare array",
			input: `[{"a":1}]`,
			want:  `[{"a":1}]`,
		},
		{
			name:  "bare object",
			input: `{"key":"val"}`,
			want:  `{"key":"val"}`,
		},
		{
			name:  "json fenced array",
			input: "```json\n[{\"a\":1}]\n```",
			want:  `[{"a":1}]`,
		},
		{
			name:  "plain fenced object",
			input: "```\n{\"key\":\"val\"}\n```",
			want:  `{"key":"val"}`,
		},
		{
			name:  "text preamble before array",
			input: "Here are the routes:\n[{\"method\":\"GET\"}]",
			want:  `[{"method":"GET"}]`,
		},
		{
			name:  "text preamble before object",
			input: "Result:\n{\"resource_type\":\"user\"}",
			want:  `{"resource_type":"user"}`,
		},
		{
			name:  "fence then preamble then json",
			input: "```\nHere it is: [{\"x\":1}]\n```",
			want:  `[{"x":1}]`,
		},
		{
			name:  "nested array",
			input: `[{"items":[1,2,3],"ok":true}]`,
			want:  `[{"items":[1,2,3],"ok":true}]`,
		},
		{
			name:  "whitespace around json",
			input: "  \n  [1,2,3]  \n  ",
			want:  "[1,2,3]",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "no json chars",
			input: "just plain text",
			want:  "just plain text",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractJSON(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
