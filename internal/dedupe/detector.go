// internal/dedupe/detector.go
package dedupe

import (
	"fmt"
	"sort"
	"strings"
)

// sortedJoin returns a deterministic string from a sorted slice of strings.
func sortedJoin(ss []string) string {
	cp := make([]string, len(ss))
	copy(cp, ss)
	sort.Strings(cp)
	return strings.Join(cp, ",")
}

// FindDuplicates detects exact and structural duplicate groups.
// threshold is the minimum Jaccard similarity on assertion targets (0.0–1.0).
func FindDuplicates(cases []LoadedCase, threshold float64) ([]DuplicateGroup, error) {
	if len(cases) == 0 {
		return nil, nil
	}

	var groups []DuplicateGroup
	inGroup := map[string]bool{}

	// Pass 1: exact duplicates — method/path/status/body/assertions all identical.
	exactBuckets := map[string][]LoadedCase{}
	for _, lc := range cases {
		k := exactKey(lc.TC)
		exactBuckets[k] = append(exactBuckets[k], lc)
	}

	for _, k := range sortedKeys(exactBuckets) {
		bucket := exactBuckets[k]
		if len(bucket) < 2 {
			continue
		}
		groups = append(groups, DuplicateGroup{
			Kind:       MatchExact,
			Similarity: 1.0,
			Cases:      scoreAndRank(bucket),
		})
		for _, lc := range bucket {
			inGroup[lc.FilePath] = true
		}
	}

	// Pass 2a: structural duplicates — same method/path/status/body (non-empty) but different assertions.
	bodyBuckets := map[string][]LoadedCase{}
	for _, lc := range cases {
		if inGroup[lc.FilePath] {
			continue
		}
		k := bodyKey(lc.TC)
		if k == "" {
			// empty body — handled in pass 2b via Jaccard
			continue
		}
		bodyBuckets[k] = append(bodyBuckets[k], lc)
	}

	for _, k := range sortedKeys(bodyBuckets) {
		bucket := bodyBuckets[k]
		if len(bucket) < 2 {
			continue
		}
		groups = append(groups, DuplicateGroup{
			Kind:       MatchStructural,
			Similarity: 1.0,
			Cases:      scoreAndRank(bucket),
		})
		for _, lc := range bucket {
			inGroup[lc.FilePath] = true
		}
	}

	// Pass 2b: structural duplicates — same method/path/status, high Jaccard on assertion targets.
	jaccardBuckets := map[string][]LoadedCase{}
	for _, lc := range cases {
		if inGroup[lc.FilePath] {
			continue
		}
		k := structKey(lc.TC)
		jaccardBuckets[k] = append(jaccardBuckets[k], lc)
	}

	for _, k := range sortedKeys(jaccardBuckets) {
		bucket := jaccardBuckets[k]
		if len(bucket) < 2 {
			continue
		}
		maxSim := 0.0
		for i := 0; i < len(bucket); i++ {
			for j := i + 1; j < len(bucket); j++ {
				sim := jaccardSimilarity(
					bucket[i].TC.AssertionTargets,
					bucket[j].TC.AssertionTargets,
				)
				if sim > maxSim {
					maxSim = sim
				}
			}
		}
		if maxSim < threshold {
			continue
		}
		groups = append(groups, DuplicateGroup{
			Kind:       MatchStructural,
			Similarity: maxSim,
			Cases:      scoreAndRank(bucket),
		})
	}

	return groups, nil
}

func exactKey(snap TestCaseSnapshot) string {
	return fmt.Sprintf("%s|%s|%d|%s|%s", snap.Method, snap.Path, snap.ExpectedStatus, snap.BodyJSON, sortedJoin(snap.AssertionTargets))
}

// bodyKey groups cases by method/path/status/body (for non-empty bodies).
func bodyKey(snap TestCaseSnapshot) string {
	if snap.BodyJSON == "" {
		return ""
	}
	return fmt.Sprintf("%s|%s|%d|%s", snap.Method, snap.Path, snap.ExpectedStatus, snap.BodyJSON)
}

func structKey(snap TestCaseSnapshot) string {
	return fmt.Sprintf("%s|%s|%d", snap.Method, snap.Path, snap.ExpectedStatus)
}

// jaccardSimilarity computes |A ∩ B| / |A ∪ B|. Two empty sets → 1.0.
func jaccardSimilarity(a, b []string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	setA := toStringSet(a)
	setB := toStringSet(b)
	intersection := 0
	for k := range setA {
		if setB[k] {
			intersection++
		}
	}
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 1.0
	}
	return float64(intersection) / float64(union)
}

func toStringSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// scoreAndRank assigns CaseScore entries; marks winner (most assertions, then lex filepath) as Keep.
func scoreAndRank(bucket []LoadedCase) []CaseScore {
	scores := make([]CaseScore, len(bucket))
	for i, lc := range bucket {
		scores[i] = CaseScore{
			FilePath:       lc.FilePath,
			AssertionCount: len(lc.TC.AssertionTargets),
		}
	}
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].AssertionCount != scores[j].AssertionCount {
			return scores[i].AssertionCount > scores[j].AssertionCount
		}
		return strings.Compare(scores[i].FilePath, scores[j].FilePath) < 0
	})
	scores[0].Keep = true
	return scores
}

func sortedKeys(m map[string][]LoadedCase) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
