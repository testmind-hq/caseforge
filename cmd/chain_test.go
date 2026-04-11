// cmd/chain_test.go
package cmd

import (
	"encoding/json"
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

func TestChainCommand_DataPool_Loaded(t *testing.T) {
	tmp := t.TempDir()
	t.Cleanup(func() { chainDataPool = "" })

	poolFile := filepath.Join(tmp, "pool.json")
	require.NoError(t, os.WriteFile(poolFile, []byte(`{"name": ["seeded-name"]}`), 0644))

	const specYAML = `
openapi: "3.0.0"
info: {title: T, version: "1"}
paths:
  /items:
    post:
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name: {type: string}
      responses:
        "201":
          description: created
          content:
            application/json:
              schema:
                type: object
                properties:
                  id: {type: integer}
  /items/{itemId}:
    get:
      parameters:
        - {name: itemId, in: path, required: true, schema: {type: integer}}
      responses:
        "200": {description: ok}
`
	specFile := filepath.Join(tmp, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(specYAML), 0644))
	outDir := filepath.Join(tmp, "chains")

	rootCmd.SetArgs([]string{
		"chain", "--spec", specFile, "--depth", "2",
		"--output", outDir, "--data-pool", poolFile,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("chain: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "index.json")); err != nil {
		t.Errorf("index.json not written: %v", err)
	}
}

func TestChainCommand_AddsTeardownForNonDeleteChains(t *testing.T) {
	const specYAML = `
openapi: "3.0.0"
info: {title: T, version: "1"}
paths:
  /items:
    post:
      operationId: createItem
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name: {type: string}
      responses:
        "201":
          description: created
          content:
            application/json:
              schema:
                type: object
                properties:
                  id: {type: integer}
  /items/{itemId}:
    get:
      operationId: getItem
      parameters:
        - {name: itemId, in: path, required: true, schema: {type: integer}}
      responses:
        "200": {description: ok}
    delete:
      operationId: deleteItem
      parameters:
        - {name: itemId, in: path, required: true, schema: {type: integer}}
      responses:
        "204": {description: deleted}
`
	tmp := t.TempDir()
	specFile := filepath.Join(tmp, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte(specYAML), 0644))

	outDir := filepath.Join(tmp, "chains")
	cmd := rootCmd
	cmd.SetArgs([]string{"chain", "--spec", specFile, "--depth", "2", "--output", outDir})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("chain command failed: %v", err)
	}

	indexPath := filepath.Join(outDir, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index.json: %v", err)
	}
	var index struct {
		TestCases []map[string]any `json:"test_cases"`
	}
	if err := json.Unmarshal(data, &index); err != nil {
		t.Fatalf("unmarshal cases: %v", err)
	}

	hasTeardown := false
	for _, tc := range index.TestCases {
		steps, _ := tc["steps"].([]any)
		for _, s := range steps {
			step := s.(map[string]any)
			if step["type"] == "teardown" {
				hasTeardown = true
			}
		}
	}
	if !hasTeardown {
		t.Error("expected at least one chain case with a teardown step, got none")
	}
}
