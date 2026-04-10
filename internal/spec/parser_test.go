// internal/spec/parser_test.go
package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseResponseHeaders(t *testing.T) {
	yaml := []byte(`
openapi: "3.0.0"
info: {title: T, version: "1"}
paths:
  /items:
    post:
      operationId: createItem
      responses:
        "201":
          description: Created
          headers:
            Location:
              schema:
                type: string
          content: {}
`)
	ps, err := parseRawSpec(yaml, "")
	require.NoError(t, err)
	require.Len(t, ps.Operations, 1)
	resp, ok := ps.Operations[0].Responses["201"]
	require.True(t, ok)
	assert.Equal(t, "string", resp.Headers["Location"])
}
