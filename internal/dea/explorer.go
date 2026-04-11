// internal/dea/explorer.go
package dea

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// Explorer runs the hypothesis-test-infer loop against a live API.
type Explorer struct {
	TargetURL string
	MaxProbes int
	DryRun    bool // if true, seed hypotheses but skip HTTP execution

	gen  *datagen.Generator
	pool *datagen.DataPool
}

// NewExplorer creates an Explorer with sensible defaults.
func NewExplorer(targetURL string, maxProbes int) *Explorer {
	if maxProbes <= 0 {
		maxProbes = 50
	}
	pool := datagen.NewDataPool()
	gen := datagen.NewGenerator(nil)
	gen.Pool = pool // observed values feed back into probe generation within the same run
	return &Explorer{
		TargetURL: targetURL,
		MaxProbes: maxProbes,
		gen:       gen,
		pool:      pool,
	}
}

// DataPool returns the pool of field values observed from 2xx responses.
func (e *Explorer) DataPool() *datagen.DataPool { return e.pool }

// hypothesisCategory returns the expected RuleCategory for a hypothesis kind in dry-run mode.
func hypothesisCategory(kind HypothesisKind) RuleCategory {
	switch kind {
	case KindOptionalField:
		return CategoryBehavior
	case KindNullValue:
		return CategoryBehavior
	default:
		return CategoryFieldConstraint
	}
}

// Explore seeds hypotheses for each operation, runs probes, and returns a report.
func (e *Explorer) Explore(ctx context.Context, s *spec.ParsedSpec) (*ExplorationReport, error) {
	report := &ExplorationReport{
		TargetURL:  e.TargetURL,
		ExploredAt: time.Now(),
	}

	probesRun := 0

	for _, op := range s.Operations {
		// Check context cancellation between operations
		if ctx.Err() != nil {
			return report, ctx.Err()
		}
		if probesRun >= e.MaxProbes {
			break
		}

		hypotheses := SeedHypotheses(op)

		for _, h := range hypotheses {
			// Check context cancellation between probes
			if ctx.Err() != nil {
				return report, ctx.Err()
			}
			if probesRun >= e.MaxProbes {
				break
			}

			probe := DesignProbe(h, op, e.gen)
			h.Probe = probe

			if e.DryRun {
				// In dry-run mode, record planned hypotheses as low-confidence pending rules.
				// Category is derived from hypothesis kind (not a blanket CategoryBehavior).
				rule := &DiscoveredRule{
					ID:          fmt.Sprintf("PLAN-%s", h.ID),
					Operation:   h.Operation,
					FieldPath:   h.FieldPath,
					Category:    hypothesisCategory(h.Kind),
					Description: fmt.Sprintf("[DRY RUN] Planned probe: %s", h.Description),
					Confidence:  ConfidenceLow,
				}
				report.Rules = append(report.Rules, *rule)
				probesRun++ // count planned probes toward the cap in DryRun
				continue
			}

			ev, err := RunProbe(ctx, e.TargetURL, probe)
			if err != nil {
				// Non-fatal: skip this probe and continue
				continue
			}
			probesRun++

			// Confirm or refute: check whether actual status falls in the expected range.
			is4xx := ev.ActualStatus >= 400 && ev.ActualStatus < 500
			is2xx := ev.ActualStatus >= 200 && ev.ActualStatus < 300
			expectedIs4xx := probe.ExpectedStatus >= 400 && probe.ExpectedStatus < 500
			expectedIs2xx := probe.ExpectedStatus >= 200 && probe.ExpectedStatus < 300

			confirmed := (expectedIs4xx && is4xx) || (expectedIs2xx && is2xx)
			h.Resolve(ev, confirmed)

			// Extract field values from 2xx responses into the data pool
			if is2xx && ev.ActualBody != "" {
				extractBodyToPool(e.pool, ev.ActualBody)
			}

			rule := InferRule(h)
			if rule != nil {
				report.Rules = append(report.Rules, *rule)
			}
		}
	}

	report.TotalProbes = probesRun
	return report, nil
}

// extractBodyToPool does a shallow JSON parse of a response body and adds
// all scalar field values (string, number, bool) to the pool.
func extractBodyToPool(pool *datagen.DataPool, body string) {
	var m map[string]any
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return // non-JSON or array body — skip
	}
	for k, v := range m {
		switch v.(type) {
		case string, float64, bool:
			pool.Add(k, v)
		}
	}
}
