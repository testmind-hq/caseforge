package mutation

import "sort"

// ComputeOperationScores groups CaseMutationResults by operation and returns
// per-operation scores sorted weakest-first (score ascending, then alphabetically).
// Results with empty Operation are grouped under "(unknown)".
func ComputeOperationScores(results []CaseMutationResult) []OperationScore {
	type stats struct{ killed, survivors int }
	byOp := map[string]*stats{}
	for _, r := range results {
		op := r.Operation
		if op == "" {
			op = "(unknown)"
		}
		if byOp[op] == nil {
			byOp[op] = &stats{}
		}
		if r.Survived {
			byOp[op].survivors++
		} else {
			byOp[op].killed++
		}
	}
	scores := make([]OperationScore, 0, len(byOp))
	for op, s := range byOp {
		total := s.killed + s.survivors
		score := 0.0
		if total > 0 {
			score = float64(s.killed) / float64(total)
		}
		scores = append(scores, OperationScore{
			Operation:     op,
			TotalRuns:     total,
			Killed:        s.killed,
			Survivors:     s.survivors,
			MutationScore: score,
		})
	}
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].MutationScore != scores[j].MutationScore {
			return scores[i].MutationScore < scores[j].MutationScore
		}
		return scores[i].Operation < scores[j].Operation
	})
	return scores
}
