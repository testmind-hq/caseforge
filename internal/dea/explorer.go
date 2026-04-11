// internal/dea/explorer.go
package dea

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/testmind-hq/caseforge/internal/datagen"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// Explorer runs the hypothesis-test-infer loop against a live API.
type Explorer struct {
	TargetURL           string
	MaxProbes           int
	DryRun              bool // if true, seed hypotheses but skip HTTP execution
	PrioritizeUncovered bool // if true, use two-pass scheduling: breadth first, then depth on non-2xx ops
	MaxFailures         int  // stop exploration after this many rules discovered (0 = unlimited)

	gen       *datagen.Generator
	pool      *datagen.DataPool
	seenRules map[ruleKey]bool
}

// ruleKey is the deduplication key for discovered rules within one exploration run.
// Rules with the same (operation, category, fieldPath) are considered duplicates.
type ruleKey struct {
	operation string
	category  RuleCategory
	fieldPath string
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
		seenRules: make(map[ruleKey]bool),
	}
}

// appendRule appends rule to report.Rules only if no rule with the same
// (operation, category, fieldPath) key has been added in this run.
// Also enforces the MaxFailures cap when non-zero.
// Returns true if the rule was appended, false if it was a duplicate or capped.
func (e *Explorer) appendRule(report *ExplorationReport, rule *DiscoveredRule) bool {
	k := ruleKey{operation: rule.Operation, category: rule.Category, fieldPath: rule.FieldPath}
	if e.seenRules[k] {
		return false
	}
	e.seenRules[k] = true
	report.Rules = append(report.Rules, *rule)
	return true
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

	if e.PrioritizeUncovered && !e.DryRun {
		return e.exploreWithPriority(ctx, s, report)
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
				if e.appendRule(report, rule) && e.MaxFailures > 0 && len(report.Rules) >= e.MaxFailures {
					report.TotalProbes = probesRun
					return report, nil
				}
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
				if rule := validateProbeResponse(op, probe, ev); rule != nil {
					e.appendRule(report, rule)
				}
			}

			rule := InferRule(h)
			if rule != nil {
				if e.appendRule(report, rule) && e.MaxFailures > 0 && len(report.Rules) >= e.MaxFailures {
					report.TotalProbes = probesRun
					return report, nil
				}
			}
		}
	}

	report.TotalProbes = probesRun
	return report, nil
}

// exploreWithPriority implements a two-pass probe scheduling strategy inspired by
// EvoMaster's focused-search mode:
//
//	Pass 1: run the first hypothesis for every operation (breadth scan, 1 probe/op).
//	Pass 2: allocate remaining budget to operations that didn't return 2xx in pass 1,
//	         running their remaining hypotheses in full before moving to covered ops.
func (e *Explorer) exploreWithPriority(ctx context.Context, s *spec.ParsedSpec, report *ExplorationReport) (*ExplorationReport, error) {
	type opState struct {
		op         *spec.Operation
		hypotheses []*HypothesisNode
		got2xx     bool
	}

	states := make([]opState, 0, len(s.Operations))
	for _, op := range s.Operations {
		hyps := SeedHypotheses(op)
		if len(hyps) > 0 {
			states = append(states, opState{op: op, hypotheses: hyps})
		}
	}

	probesRun := 0

	// Pass 1: one probe per operation for breadth coverage.
	for i := range states {
		if ctx.Err() != nil || probesRun >= e.MaxProbes {
			break
		}
		h := states[i].hypotheses[0]
		probe := DesignProbe(h, states[i].op, e.gen)
		ev, err := RunProbe(ctx, e.TargetURL, probe)
		if err != nil {
			continue
		}
		probesRun++
		states[i].got2xx = ev.ActualStatus >= 200 && ev.ActualStatus < 300
		expectedIs4xx := probe.ExpectedStatus >= 400
		confirmed := (expectedIs4xx && ev.ActualStatus >= 400) ||
			(!expectedIs4xx && states[i].got2xx)
		h.Resolve(ev, confirmed)
		if rule := InferRule(h); rule != nil {
			if e.appendRule(report, rule) && e.MaxFailures > 0 && len(report.Rules) >= e.MaxFailures {
				report.TotalProbes = probesRun
				return report, nil
			}
		}
		if states[i].got2xx && ev.ActualBody != "" {
			extractBodyToPool(e.pool, ev.ActualBody)
			if rule := validateProbeResponse(states[i].op, probe, ev); rule != nil {
				e.appendRule(report, rule)
			}
		}
	}

	// Pass 2: prioritize ops that didn't get 2xx — run their remaining hypotheses first.
	// Sort: failing ops (got2xx=false) before covered ops.
	sort.SliceStable(states, func(i, j int) bool {
		return !states[i].got2xx && states[j].got2xx
	})

	for _, st := range states {
		for _, h := range st.hypotheses[1:] {
			if ctx.Err() != nil || probesRun >= e.MaxProbes {
				break
			}
			probe := DesignProbe(h, st.op, e.gen)
			ev, err := RunProbe(ctx, e.TargetURL, probe)
			if err != nil {
				continue
			}
			probesRun++
			is2xx := ev.ActualStatus >= 200 && ev.ActualStatus < 300
			expectedIs4xx := probe.ExpectedStatus >= 400
			confirmed := (expectedIs4xx && ev.ActualStatus >= 400) || (!expectedIs4xx && is2xx)
			h.Resolve(ev, confirmed)
			if rule := InferRule(h); rule != nil {
				if e.appendRule(report, rule) && e.MaxFailures > 0 && len(report.Rules) >= e.MaxFailures {
					report.TotalProbes = probesRun
					return report, nil
				}
			}
			if is2xx && ev.ActualBody != "" {
				extractBodyToPool(e.pool, ev.ActualBody)
				if rule := validateProbeResponse(st.op, probe, ev); rule != nil {
					e.appendRule(report, rule)
				}
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
