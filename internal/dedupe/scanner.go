// internal/dedupe/scanner.go
package dedupe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// ScanCases reads every *.json file in casesDir and returns a LoadedCase for each
// successfully parsed schema.TestCase. Unparseable files are silently skipped.
// Returns a non-nil error if casesDir does not exist.
func ScanCases(casesDir string) ([]LoadedCase, error) {
	entries, err := os.ReadDir(casesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cases directory not found: %s", casesDir)
		}
		return nil, fmt.Errorf("read cases dir: %w", err)
	}

	var loaded []LoadedCase
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		fp := filepath.Join(casesDir, e.Name())
		data, readErr := os.ReadFile(fp)
		if readErr != nil {
			continue
		}
		var tc schema.TestCase
		if unmarshalErr := json.Unmarshal(data, &tc); unmarshalErr != nil {
			continue
		}
		loaded = append(loaded, LoadedCase{FilePath: fp, TC: snapshotFrom(tc)})
	}
	return loaded, nil
}

// snapshotFrom extracts duplicate-detection fields from a schema.TestCase.
func snapshotFrom(tc schema.TestCase) TestCaseSnapshot {
	snap := TestCaseSnapshot{}

	// Use the first "test"-typed step; fall back to the very first step.
	for _, step := range tc.Steps {
		if step.Type == "test" || snap.Method == "" {
			snap.Method = strings.ToUpper(step.Method)
			snap.Path = step.Path
			for _, a := range step.Assertions {
				if a.Target == "status_code" {
					if v, ok := toInt(a.Expected); ok {
						snap.ExpectedStatus = v
					}
					break
				}
			}
			if step.Body != nil {
				snap.BodyJSON = normalizeBodyJSON(step.Body)
			}
			if step.Type == "test" {
				break
			}
		}
	}

	// Collect unique assertion targets across ALL steps.
	seen := map[string]struct{}{}
	for _, step := range tc.Steps {
		for _, a := range step.Assertions {
			if _, exists := seen[a.Target]; !exists {
				seen[a.Target] = struct{}{}
				snap.AssertionTargets = append(snap.AssertionTargets, a.Target)
			}
		}
	}
	return snap
}

// normalizeBodyJSON marshals v to a compact JSON string for equality comparison.
func normalizeBodyJSON(v any) string {
	raw, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	var m any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(m); err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}

// toInt converts JSON-decoded numeric values to int.
func toInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case int:
		return x, true
	case int64:
		return int(x), true
	}
	return 0, false
}
