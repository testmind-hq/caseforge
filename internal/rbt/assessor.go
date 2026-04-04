// internal/rbt/assessor.go
package rbt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// ScanCases scans casesDir for *.json test case files and builds
// an index from "METHOD /path" → []TestCaseRef.
func ScanCases(casesDir string) (map[string][]TestCaseRef, error) {
	entries, err := os.ReadDir(casesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	idx := make(map[string][]TestCaseRef)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(casesDir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var tc schema.TestCase
		if err := json.Unmarshal(data, &tc); err != nil {
			continue
		}
		key := specPathKey(tc.Source.SpecPath)
		if key == "" {
			continue
		}
		idx[key] = append(idx[key], TestCaseRef{
			File:   path,
			CaseID: tc.ID,
			Title:  tc.Title,
		})
	}
	return idx, nil
}

// specPathKey extracts "METHOD /path" from a SpecPath string.
// SpecPath format: "METHOD /path [optional_field_info]"
func specPathKey(specPath string) string {
	parts := strings.Fields(specPath)
	if len(parts) < 2 {
		return ""
	}
	return strings.ToUpper(parts[0]) + " " + parts[1]
}

// Assess builds a RiskReport by cross-referencing the spec operations against
// the set of affected route mappings and the test case index.
// changedFiles is used for report metadata only.
func Assess(
	sp *spec.ParsedSpec,
	affected []RouteMapping,
	caseIndex map[string][]TestCaseRef,
	base, head string,
	changedFiles []ChangedFile,
) RiskReport {
	// Build a quick lookup: "METHOD /path" → []RouteMapping
	affectedByKey := make(map[string][]RouteMapping)
	for _, rm := range affected {
		key := rm.Method + " " + rm.RoutePath
		affectedByKey[key] = append(affectedByKey[key], rm)
	}

	var ops []OperationCoverage
	totalAffected, totalCovered, totalUncovered := 0, 0, 0

	for _, op := range sp.Operations {
		key := op.Method + " " + op.Path
		rms := affectedByKey[key]
		isAffected := len(rms) > 0
		cases := caseIndex[key]

		risk := computeRisk(isAffected, len(cases), lowestConfidence(rms))

		oc := OperationCoverage{
			OperationID: op.OperationID,
			Method:      op.Method,
			Path:        op.Path,
			Affected:    isAffected,
			SourceRefs:  rms,
			TestCases:   cases,
			Risk:        risk,
		}
		ops = append(ops, oc)

		if isAffected {
			totalAffected++
			if len(cases) == 0 && risk != RiskUncertain {
				totalUncovered++
			} else if len(cases) > 0 {
				totalCovered++
			}
		}
	}

	riskScore := 0.0
	if totalAffected > 0 {
		riskScore = float64(totalUncovered) / float64(totalAffected)
	}

	return RiskReport{
		DiffBase:       base,
		DiffHead:       head,
		ChangedFiles:   changedFiles,
		Operations:     ops,
		TotalAffected:  totalAffected,
		TotalCovered:   totalCovered,
		TotalUncovered: totalUncovered,
		RiskScore:      riskScore,
		GeneratedAt:    time.Now(),
	}
}

func computeRisk(isAffected bool, testCount int, minConfidence float64) RiskLevel {
	if !isAffected {
		return RiskNone
	}
	if minConfidence > 0 && minConfidence < 0.5 {
		return RiskUncertain
	}
	switch {
	case testCount == 0:
		return RiskHigh
	case testCount == 1:
		return RiskMedium
	default:
		return RiskLow
	}
}

func lowestConfidence(rms []RouteMapping) float64 {
	min := 1.0
	for _, rm := range rms {
		if rm.Confidence > 0 && rm.Confidence < min {
			min = rm.Confidence
		}
	}
	if min == 1.0 && len(rms) > 0 {
		// No confidence set → treat as fully confident
		return 1.0
	}
	return min
}
