// internal/suite/suite.go
package suite

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// ValidationError describes a single problem found when validating a TestSuite.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validate checks that a TestSuite is well-formed:
//   - ID and Title are non-empty
//   - Kind is "sequential" or "chain"
//   - Each SuiteCase has a non-empty CaseID
//   - No duplicate CaseIDs within the suite
//   - All depends_on references exist within the suite's own case list
//   - The dependency graph is acyclic (Kahn's topological sort)
//   - If knownCases is non-nil, all CaseIDs must exist in that slice
func Validate(s *schema.TestSuite, knownCases []schema.TestCase) []ValidationError {
	var errs []ValidationError

	if s.ID == "" {
		errs = append(errs, ValidationError{Field: "id", Message: "must not be empty"})
	}
	if s.Title == "" {
		errs = append(errs, ValidationError{Field: "title", Message: "must not be empty"})
	}
	if s.Kind != "sequential" && s.Kind != "chain" {
		errs = append(errs, ValidationError{
			Field:   "kind",
			Message: fmt.Sprintf("must be 'sequential' or 'chain', got %q", s.Kind),
		})
	}

	// Build set of case_ids in this suite for reference checking.
	suiteIDs := make(map[string]bool, len(s.Cases))
	for _, sc := range s.Cases {
		if sc.CaseID == "" {
			errs = append(errs, ValidationError{Field: "cases[].case_id", Message: "must not be empty"})
			continue
		}
		if suiteIDs[sc.CaseID] {
			errs = append(errs, ValidationError{
				Field:   "cases[].case_id",
				Message: fmt.Sprintf("duplicate case_id %q", sc.CaseID),
			})
		}
		suiteIDs[sc.CaseID] = true
	}

	// Build set of all known case IDs from index.json (if provided).
	knownIDs := make(map[string]bool, len(knownCases))
	for _, tc := range knownCases {
		knownIDs[tc.ID] = true
	}

	for _, sc := range s.Cases {
		// Check that the referenced case exists in index.json.
		if len(knownCases) > 0 && !knownIDs[sc.CaseID] {
			errs = append(errs, ValidationError{
				Field:   "cases[].case_id",
				Message: fmt.Sprintf("%q not found in index.json", sc.CaseID),
			})
		}
		// Check that all depends_on references exist in the suite.
		for _, dep := range sc.DependsOn {
			if !suiteIDs[dep] {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("cases[%s].depends_on", sc.CaseID),
					Message: fmt.Sprintf("unknown case_id %q", dep),
				})
			}
		}
	}

	// Cycle detection via Kahn's algorithm.
	if cycleErr := detectCycle(s); cycleErr != nil {
		errs = append(errs, ValidationError{Field: "cases.depends_on", Message: cycleErr.Error()})
	}

	return errs
}

// TopologicalOrder returns the SuiteCase slice sorted so that each case appears
// after all of its dependencies. Returns an error if the graph has cycles.
func TopologicalOrder(s *schema.TestSuite) ([]schema.SuiteCase, error) {
	// In-degree map.
	inDegree := make(map[string]int, len(s.Cases))
	// Adjacency list: dep → list of cases that depend on dep.
	adj := make(map[string][]string, len(s.Cases))
	caseMap := make(map[string]schema.SuiteCase, len(s.Cases))

	for _, sc := range s.Cases {
		inDegree[sc.CaseID] = inDegree[sc.CaseID] // ensure key exists
		caseMap[sc.CaseID] = sc
		for _, dep := range sc.DependsOn {
			adj[dep] = append(adj[dep], sc.CaseID)
			inDegree[sc.CaseID]++
		}
	}

	// Collect all zero-in-degree nodes in original insertion order.
	var queue []string
	for _, sc := range s.Cases {
		if inDegree[sc.CaseID] == 0 {
			queue = append(queue, sc.CaseID)
		}
	}

	var sorted []schema.SuiteCase
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		sorted = append(sorted, caseMap[cur])
		for _, next := range adj[cur] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(sorted) != len(s.Cases) {
		return nil, fmt.Errorf("dependency cycle detected in suite %q", s.ID)
	}
	return sorted, nil
}

// detectCycle returns a non-nil error if the suite's dependency graph contains a cycle.
func detectCycle(s *schema.TestSuite) error {
	_, err := TopologicalOrder(s)
	return err
}

// LoadSuiteFile reads and parses a suite.json file.
func LoadSuiteFile(path string) (*schema.TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading suite file: %w", err)
	}
	var s schema.TestSuite
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing suite file: %w", err)
	}
	return &s, nil
}

// WriteSuiteFile serializes a TestSuite and writes it to path.
func WriteSuiteFile(s *schema.TestSuite, path string) error {
	s.Schema = schema.SuiteSchemaURL
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling suite: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing suite file: %w", err)
	}
	return nil
}
