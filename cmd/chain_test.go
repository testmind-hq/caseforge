// cmd/chain_test.go
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const chainSpec = `
openapi: "3.0.0"
info: {title: TestAPI, version: "1"}
paths:
  /orders:
    post:
      operationId: createOrder
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                product_id: {type: string}
                quantity:   {type: integer}
      responses:
        "201":
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:     {type: string}
                  status: {type: string}
  /orders/{orderId}:
    get:
      operationId: getOrder
      parameters:
        - name: orderId
          in: path
          required: true
          schema: {type: string}
      responses:
        "200": {description: OK}
    put:
      operationId: updateOrder
      parameters:
        - name: orderId
          in: path
          required: true
          schema: {type: string}
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                status: {type: string}
      responses:
        "200": {description: OK}
    delete:
      operationId: deleteOrder
      parameters:
        - name: orderId
          in: path
          required: true
          schema: {type: string}
      responses:
        "204": {description: No Content}
`

func TestChainCommand_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "chain" {
			found = true
			break
		}
	}
	assert.True(t, found, "chain command must be registered on rootCmd")
}

func TestChainCommand_HasRequiredFlags(t *testing.T) {
	var cmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "chain" {
			cmd = c
			break
		}
	}
	require.NotNil(t, cmd)
	assert.NotNil(t, cmd.Flags().Lookup("spec"))
	assert.NotNil(t, cmd.Flags().Lookup("depth"))
	assert.NotNil(t, cmd.Flags().Lookup("output"))
}

func TestChainCommand_GeneratesChainCases(t *testing.T) {
	specFile := filepath.Join(t.TempDir(), "chain.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(chainSpec), 0644))
	outDir := t.TempDir()

	chainDepth = 2
	chainSpecPath = specFile
	chainOutput = outDir

	require.NoError(t, runChain(chainCmd, nil))

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases, "chain command must produce chain cases")
	for _, tc := range cases {
		assert.Equal(t, "chain", tc.Kind)
		assert.GreaterOrEqual(t, len(tc.Steps), 1)
	}
}

func TestChainCommand_Depth1_SingleOpCases(t *testing.T) {
	specFile := filepath.Join(t.TempDir(), "chain.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(chainSpec), 0644))
	outDir := t.TempDir()

	chainDepth = 1
	chainSpecPath = specFile
	chainOutput = outDir

	require.NoError(t, runChain(chainCmd, nil))

	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
	for _, tc := range cases {
		assert.Equal(t, "chain", tc.Kind)
		assert.Equal(t, 1, len(tc.Steps), "depth-1 cases must have exactly 1 step")
	}
}
