// internal/mutation/types_test.go
package mutation_test

import (
	"net/http"
	"testing"

	"github.com/testmind-hq/caseforge/internal/mutation"
)

func TestOperatorInterface(t *testing.T) {
	op := &noopOperator{}
	if op.Name() == "" {
		t.Fatal("Name() must not be empty")
	}
	body := []byte(`{"id":1}`)
	out, err := op.Apply(&http.Response{}, body)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(body) {
		t.Fatalf("noop should return body unchanged, got %s", out)
	}
}

func TestMutationRunScore(t *testing.T) {
	run := mutation.MutationRun{
		TotalRuns: 10,
		Killed:    7,
		Survivors: 3,
	}
	run.MutationScore = float64(run.Killed) / float64(run.TotalRuns)
	if run.MutationScore < 0.69 || run.MutationScore > 0.71 {
		t.Fatalf("expected ~0.70, got %f", run.MutationScore)
	}
}

// noopOperator satisfies the Operator interface for tests.
type noopOperator struct{}

func (n *noopOperator) Name() string                                         { return "noop" }
func (n *noopOperator) Apply(_ *http.Response, body []byte) ([]byte, error) { return body, nil }
