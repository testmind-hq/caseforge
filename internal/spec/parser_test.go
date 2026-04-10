// internal/spec/parser_test.go
package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsedSpec_LinksPopulated(t *testing.T) {
	const specYAML = `
openapi: "3.0.0"
info:
  title: Links Test
  version: "1"
paths:
  /users:
    post:
      operationId: createUser
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
      responses:
        "201":
          description: created
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
          links:
            GetUserById:
              operationId: getUser
              parameters:
                userId: "$response.body#/id"
  /users/{userId}:
    get:
      operationId: getUser
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: ok
`
	ps, err := parseRawSpec([]byte(specYAML), "")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var createOp *Operation
	for _, op := range ps.Operations {
		if op.OperationID == "createUser" {
			createOp = op
			break
		}
	}
	if createOp == nil {
		t.Fatal("createUser operation not found")
	}
	if len(createOp.Links) == 0 {
		t.Fatal("expected links on createUser, got none")
	}
	link := createOp.Links[0]
	if link.Name != "GetUserById" {
		t.Errorf("link name = %q, want GetUserById", link.Name)
	}
	if link.OperationID != "getUser" {
		t.Errorf("link operationId = %q, want getUser", link.OperationID)
	}
	if link.ResponseCode != "201" {
		t.Errorf("link responseCode = %q, want 201", link.ResponseCode)
	}
	if expr, ok := link.Parameters["userId"]; !ok || expr != "$response.body#/id" {
		t.Errorf("link parameter userId = %q, want $response.body#/id", expr)
	}
}

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
