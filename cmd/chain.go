// cmd/chain.go
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	assertpkg "github.com/testmind-hq/caseforge/internal/assert"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/methodology"
	"github.com/testmind-hq/caseforge/internal/output/render"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/output/writer"
	"github.com/testmind-hq/caseforge/internal/spec"
)

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Generate multi-step chain cases via BFS over the dependency graph",
	Long: `Builds a dependency graph from the spec (producer-consumer relationships between operations)
and generates chain test cases of length 1..depth using BFS.

Depth 1: one single-step test per operation.
Depth 2: two-step chains (create → read/update/delete).
Depth 3: three-step chains (create → update → read, etc.).`,
	RunE: runChain,
}

var (
	chainSpecPath    string
	chainOutput      string
	chainDepth       int
	chainFormat      string
	chainDataPool    string // path to DataPool JSON (from explore --export-pool)
	chainSeedPostman string // path to Postman Collection v2.1 JSON
)

func init() {
	rootCmd.AddCommand(chainCmd)
	chainCmd.Flags().StringVar(&chainSpecPath, "spec", "", "OpenAPI spec file or URL (required)")
	chainCmd.Flags().StringVar(&chainOutput, "output", "./chains", "Output directory")
	chainCmd.Flags().IntVar(&chainDepth, "depth", 2, "Maximum chain depth (1..4)")
	chainCmd.Flags().StringVar(&chainFormat, "format", "hurl", "Output format: hurl|markdown|csv|postman|k6")
	chainCmd.Flags().StringVar(&chainDataPool, "data-pool", "", "JSON data pool file (from explore --export-pool)")
	chainCmd.Flags().StringVar(&chainSeedPostman, "seed-postman", "", "Postman Collection v2.1 JSON file; extracts body field values as seed data")
	_ = chainCmd.MarkFlagRequired("spec")
}

func runChain(cmd *cobra.Command, args []string) error {
	if chainDepth < 1 || chainDepth > 4 {
		return fmt.Errorf("invalid --depth %d: must be between 1 and 4", chainDepth)
	}

	loader := spec.NewLoader()
	parsedSpec, err := loader.Load(chainSpecPath)
	if err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}

	depGraph := methodology.BuildDepGraph(parsedSpec.Operations)
	gen := datagen.NewGenerator(nil)
	if chainDataPool != "" {
		pool, err := datagen.LoadDataPool(chainDataPool)
		if err != nil {
			return fmt.Errorf("loading data pool: %w", err)
		}
		gen.Pool = pool
	}
	if chainSeedPostman != "" {
		pmPool, err := datagen.ParsePostmanCollection(chainSeedPostman)
		if err != nil {
			return fmt.Errorf("loading postman collection: %w", err)
		}
		if gen.Pool == nil {
			gen.Pool = pmPool
		} else {
			gen.Pool.Merge(pmPool)
		}
	}
	cases := bfsChainCases(parsedSpec.Operations, depGraph, chainDepth, gen)

	if err := os.MkdirAll(chainOutput, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	w := writer.NewJSONSchemaWriter()
	if err := w.Write(cases, chainOutput, writer.WriteOptions{CaseforgeVersion: Version}); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	var renderer render.Renderer
	switch chainFormat {
	case "markdown":
		renderer = render.NewMarkdownRenderer()
	case "csv":
		renderer = render.NewCSVRenderer()
	case "postman":
		renderer = render.NewPostmanRenderer()
	case "k6":
		renderer = render.NewK6Renderer()
	default:
		renderer = render.NewHurlRenderer("")
	}
	if err := renderer.Render(cases, chainOutput); err != nil {
		return fmt.Errorf("rendering output: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Generated %d chain cases (depth 1..%d) → %s\n",
		len(cases), chainDepth, chainOutput)
	return nil
}

// bfsChainCases generates all sequences of length 1..maxDepth using the dep graph.
// Each sequence becomes a chain TestCase with steps connected by captures.
func bfsChainCases(ops []*spec.Operation, g *methodology.DepGraph, maxDepth int, gen *datagen.Generator) []schema.TestCase {
	var cases []schema.TestCase

	// Depth 1: one case per operation (single-step)
	if maxDepth >= 1 {
		for _, op := range ops {
			cases = append(cases, singleOpCase(op, gen))
		}
	}

	if maxDepth < 2 {
		return cases
	}

	// Depth 2+: extend from each dep edge
	type sequence struct {
		steps       []schema.Step
		lastOpPath  string // path of the last operation (for further extension)
		captureName string
	}

	// Seed sequences from each edge (depth 2: creator → consumer)
	var seqs []sequence
	for _, edge := range g.Edges {
		steps, captureName := buildEdgeSteps(edge, gen)
		seqs = append(seqs, sequence{
			steps:       steps,
			lastOpPath:  edge.Consumer.Path,
			captureName: captureName,
		})
		if len(steps) >= 2 {
			// Append DELETE teardown if consumer is not already DELETE
			if edge.Consumer.Method != "DELETE" {
				if td := findTeardownEdge(g, edge.Creator.Path); td != nil {
					lastID := steps[len(steps)-1].ID
					// Derive the path-param name from the DELETE path itself to avoid a
					// mismatch when td.PathParam is Link-derived and may differ from the
					// capture variable name used in prior steps.
					tdParamName := lastPathParam(td.Consumer.Path)
					tdPath := strings.ReplaceAll(td.Consumer.Path,
						fmt.Sprintf("{%s}", tdParamName),
						fmt.Sprintf("{{%s}}", captureName))
					tdStep := schema.Step{
						ID:    fmt.Sprintf("step-%d", len(steps)+1),
						Title: fmt.Sprintf("teardown: %s %s", td.Consumer.Method, tdPath),
						Type:  "teardown",
						Method: td.Consumer.Method,
						Path:   tdPath,
						Assertions: []schema.Assertion{
							{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
						},
						DependsOn: []string{lastID},
					}
					steps = append(steps, tdStep)
				}
			}
			cases = append(cases, chainCase(steps, edge.Creator.Path, "chain_bfs", 2))
		}
	}

	// Depth 3+: extend existing sequences
	for depth := 3; depth <= maxDepth; depth++ {
		var nextSeqs []sequence
		for _, seq := range seqs {
			for _, edge := range g.Edges {
				if edge.Creator.Path == seq.lastOpPath ||
					strings.HasPrefix(seq.lastOpPath, edge.Creator.Path+"/") {
					prevID := seq.steps[len(seq.steps)-1].ID
				consumerStep := buildConsumerStep(edge, seq.captureName, gen, len(seq.steps), prevID)
					newSteps := append(append([]schema.Step{}, seq.steps...), consumerStep)
					newSeq := sequence{
						steps:       newSteps,
						lastOpPath:  edge.Consumer.Path,
						captureName: seq.captureName,
					}
					nextSeqs = append(nextSeqs, newSeq)
					cases = append(cases, chainCase(newSteps, edge.Creator.Path, "chain_bfs", depth))
				}
			}
		}
		seqs = nextSeqs
		if len(seqs) == 0 {
			break
		}
	}

	return cases
}

func singleOpCase(op *spec.Operation, gen *datagen.Generator) schema.TestCase {
	body := buildValidBodyForOp(op, gen)
	var bodyAny any
	if body != nil {
		bodyAny = body
	}
	headers := map[string]string{}
	if bodyAny != nil {
		headers["Content-Type"] = "application/json"
	}
	step := schema.Step{
		ID:         "step-main",
		Title:      fmt.Sprintf("%s %s", op.Method, op.Path),
		Type:       "test",
		Method:     op.Method,
		Path:       op.Path,
		Headers:    headers,
		Body:       bodyAny,
		Assertions: assertpkg.BasicAssertions(op),
	}
	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
		Title:    fmt.Sprintf("BFS depth-1: %s %s", op.Method, op.Path),
		Kind:     "chain",
		Priority: "P2",
		Tags:     op.Tags,
		Source: schema.CaseSource{
			Technique: "chain_bfs",
			SpecPath:  fmt.Sprintf("%s %s", op.Method, op.Path),
			Rationale: "BFS depth-1 single-operation case",
		},
		Steps:       []schema.Step{step},
		GeneratedAt: time.Now(),
	}
}

func buildEdgeSteps(edge methodology.DepEdge, gen *datagen.Generator) ([]schema.Step, string) {
	captureName := edge.PathParam
	setupBody := buildValidBodyForOp(edge.Creator, gen)
	var setupBodyAny any
	if setupBody != nil {
		setupBodyAny = setupBody
	}
	setupHeaders := map[string]string{}
	if setupBodyAny != nil {
		setupHeaders["Content-Type"] = "application/json"
	}
	setupStep := schema.Step{
		ID:      "step-setup",
		Title:   fmt.Sprintf("create via %s %s", edge.Creator.Method, edge.Creator.Path),
		Type:    "setup",
		Method:  edge.Creator.Method,
		Path:    edge.Creator.Path,
		Headers: setupHeaders,
		Body:    setupBodyAny,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
		},
		Captures: []schema.Capture{{Name: captureName, From: edge.CaptureFrom}},
	}

	consumerStep := buildConsumerStep(edge, captureName, gen, 1, "step-setup")
	return []schema.Step{setupStep, consumerStep}, captureName
}

// buildConsumerStep creates a step for the consumer side of a dep edge.
// stepIdx is the 1-based position in the full step sequence (1 = first consumer).
// prevStepID is the actual ID of the preceding step in the sequence.
func buildConsumerStep(edge methodology.DepEdge, captureName string, gen *datagen.Generator, stepIdx int, prevStepID string) schema.Step {
	paramName := edge.PathParam
	path := strings.ReplaceAll(edge.Consumer.Path,
		fmt.Sprintf("{%s}", paramName),
		fmt.Sprintf("{{%s}}", captureName))

	body := buildValidBodyForOp(edge.Consumer, gen)
	var bodyAny any
	if body != nil {
		bodyAny = body
	}
	headers := map[string]string{}
	if bodyAny != nil {
		headers["Content-Type"] = "application/json"
	}

	id := fmt.Sprintf("step-%d", stepIdx+1)
	if stepIdx == 1 {
		id = "step-test"
	}

	return schema.Step{
		ID:         id,
		Title:      fmt.Sprintf("%s %s", edge.Consumer.Method, path),
		Type:       "test",
		Method:     edge.Consumer.Method,
		Path:       path,
		Headers:    headers,
		Body:       bodyAny,
		Assertions: assertpkg.BasicAssertions(edge.Consumer),
		DependsOn:  []string{prevStepID},
	}
}

func chainCase(steps []schema.Step, resourcePath, technique string, depth int) schema.TestCase {
	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
		Title:    fmt.Sprintf("BFS depth-%d chain: %s", depth, resourcePath),
		Kind:     "chain",
		Priority: "P1",
		Source: schema.CaseSource{
			Technique: technique,
			SpecPath:  resourcePath,
			Rationale: fmt.Sprintf("BFS-generated %d-step chain starting from %s", depth, resourcePath),
		},
		Steps:       steps,
		GeneratedAt: time.Now(),
	}
}

// lastPathParam extracts the path-parameter name from the last segment of a path,
// e.g. "/items/{itemId}" → "itemId". Returns "" if the last segment is not a parameter.
func lastPathParam(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	last := parts[len(parts)-1]
	if strings.HasPrefix(last, "{") && strings.HasSuffix(last, "}") {
		return strings.Trim(last, "{}")
	}
	return ""
}

// findTeardownEdge returns the first DELETE edge from the given creator path, or nil.
// Used to automatically append cleanup steps to non-DELETE BFS chains.
func findTeardownEdge(g *methodology.DepGraph, creatorPath string) *methodology.DepEdge {
	for i := range g.Edges {
		e := &g.Edges[i]
		if e.Creator.Path == creatorPath && e.Consumer.Method == "DELETE" {
			return e
		}
	}
	return nil
}

// buildValidBodyForOp generates a valid request body for an operation.
func buildValidBodyForOp(op *spec.Operation, gen *datagen.Generator) map[string]any {
	if op.RequestBody == nil {
		return nil
	}
	mt, ok := op.RequestBody.Content["application/json"]
	if !ok || mt.Schema == nil {
		return nil
	}
	body := map[string]any{}
	for fieldName, fieldSchema := range mt.Schema.Properties {
		body[fieldName] = gen.Generate(fieldSchema, fieldName)
	}
	return body
}
