// internal/dea/explorer.go
package dea

import (
	"context"
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

	gen *datagen.Generator
}

// NewExplorer creates an Explorer with sensible defaults.
func NewExplorer(targetURL string, maxProbes int) *Explorer {
	if maxProbes <= 0 {
		maxProbes = 50
	}
	return &Explorer{
		TargetURL: targetURL,
		MaxProbes: maxProbes,
		gen:       datagen.NewGenerator(nil),
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
		if probesRun >= e.MaxProbes {
			break
		}

		hypotheses := SeedHypotheses(op)

		for _, h := range hypotheses {
			if probesRun >= e.MaxProbes {
				break
			}

			probe := DesignProbe(h, op, e.gen)
			h.Probe = probe

			if e.DryRun {
				// In dry-run mode, record planned hypotheses as low-confidence pending rules
				rule := &DiscoveredRule{
					ID:          fmt.Sprintf("PLAN-%s", h.ID),
					Operation:   h.Operation,
					FieldPath:   h.FieldPath,
					Category:    CategoryBehavior,
					Description: fmt.Sprintf("[DRY RUN] Planned probe: %s", h.Description),
					Confidence:  ConfidenceLow,
				}
				report.Rules = append(report.Rules, *rule)
				continue
			}

			ev, err := RunProbe(ctx, e.TargetURL, probe)
			if err != nil {
				// Non-fatal: skip this probe and continue
				continue
			}
			probesRun++

			// Confirm or refute based on whether actual status matches expected
			is4xx := ev.ActualStatus >= 400 && ev.ActualStatus < 500
			is2xx := ev.ActualStatus >= 200 && ev.ActualStatus < 300
			expectedIs4xx := probe.ExpectedStatus >= 400

			confirmed := (expectedIs4xx && is4xx) || (!expectedIs4xx && is2xx)
			h.Resolve(ev, confirmed)

			rule := InferRule(h)
			if rule != nil {
				report.Rules = append(report.Rules, *rule)
			}
		}
	}

	report.TotalProbes = probesRun
	return report, nil
}
