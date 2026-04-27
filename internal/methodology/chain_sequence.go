package methodology

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/score"
	"github.com/testmind-hq/caseforge/internal/spec"
)

const similarityThreshold = 0.4

// ChainSequenceTechnique generates multi-step chain cases for non-CRUD
// producer-consumer relationships via field-name Jaccard similarity.
type ChainSequenceTechnique struct {
	gen *datagen.Generator
}

func NewChainSequenceTechnique() *ChainSequenceTechnique {
	return &ChainSequenceTechnique{gen: datagen.NewGenerator(nil)}
}

func (t *ChainSequenceTechnique) Name() string { return "chain_sequence" }

func (t *ChainSequenceTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	edges := BuildSimilarityEdges(s.Operations, similarityThreshold)
	var cases []schema.TestCase
	for _, edge := range edges {
		cases = append(cases, t.buildSequenceCase(edge))
	}
	return cases, nil
}

func (t *ChainSequenceTechnique) buildSequenceCase(edge SimilarityEdge) schema.TestCase {
	id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
	captureName := edge.PathParam

	producerBody := buildValidBody(t.gen, edge.Producer)
	var producerBodyAny any
	if producerBody != nil {
		producerBodyAny = producerBody
	}
	producerHeaders := map[string]string{}
	if producerBodyAny != nil {
		producerHeaders["Content-Type"] = "application/json"
	}
	setupStep := schema.Step{
		ID:    "step-setup",
		Title: fmt.Sprintf("create via %s %s", edge.Producer.Method, edge.Producer.Path),
		Type:  "setup",
		Method: edge.Producer.Method,
		Path:   edge.Producer.Path,
		Headers: producerHeaders,
		Body:    producerBodyAny,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
		},
		Captures: []schema.Capture{{Name: captureName, From: edge.CaptureFrom}},
	}

	consumerPath := strings.ReplaceAll(edge.Consumer.Path,
		fmt.Sprintf("{%s}", captureName),
		fmt.Sprintf("{{%s}}", captureName))
	consumerBody := buildValidBody(t.gen, edge.Consumer)
	var consumerBodyAny any
	if consumerBody != nil {
		consumerBodyAny = consumerBody
	}
	consumerHeaders := map[string]string{}
	if consumerBodyAny != nil {
		consumerHeaders["Content-Type"] = "application/json"
	}
	testStep := schema.Step{
		ID:    "step-test",
		Title: fmt.Sprintf("use via %s %s", edge.Consumer.Method, edge.Consumer.Path),
		Type:  "test",
		Method: edge.Consumer.Method,
		Path:   consumerPath,
		Headers: consumerHeaders,
		Body:    consumerBodyAny,
		Assertions: []schema.Assertion{
			{Target: "status_code", Operator: schema.OperatorLt, Expected: 300},
		},
		DependsOn: []string{"step-setup"},
	}

	resourceDesc := fmt.Sprintf("%s → %s %s", edge.Producer.Path, edge.Consumer.Method, edge.Consumer.Path)
	return schema.TestCase{
		Schema:   schema.SchemaBaseURL,
		Version:  "1",
		ID:       id,
		Title:    fmt.Sprintf("sequence chain: %s", resourceDesc),
		Kind:     "chain",
		Priority: "P1",
		Source: schema.CaseSource{
			Technique: "chain_sequence",
			SpecPath:  edge.Producer.Path,
			Rationale: fmt.Sprintf("field-similarity chain (score %.2f): %s → %s param %s",
				edge.Score, edge.Producer.Path, edge.Consumer.Path, captureName),
			Scenario: string(score.ScenarioChainSequence),
		},
		Steps:       []schema.Step{setupStep, testStep},
		GeneratedAt: time.Now(),
	}
}
