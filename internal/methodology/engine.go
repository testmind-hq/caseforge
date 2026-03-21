// internal/methodology/engine.go
package methodology

import (
	"github.com/testmind-hq/caseforge/internal/llm"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type Technique interface {
	Name() string
	Applies(op *spec.Operation) bool
	Generate(op *spec.Operation) ([]schema.TestCase, error)
}

type Engine struct {
	techniques []Technique
	llm        llm.LLMProvider
}

func NewEngine(llm llm.LLMProvider, techniques ...Technique) *Engine {
	return &Engine{techniques: techniques, llm: llm}
}

// Generate annotates all operations with LLM semantic info, then
// dispatches each operation to applicable techniques.
func (e *Engine) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	// TODO: implement in Task 14 (Week 3)
	return nil, nil
}
