// internal/export/exporter.go
package export

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// Exporter converts a slice of TestCase into a third-party platform format.
type Exporter interface {
	// Format returns the format name: "allure", "xray", or "testrail".
	Format() string
	// Export writes output files into outDir (created if missing).
	Export(cases []schema.TestCase, outDir string) error
}

// New returns the Exporter for the given format name.
func New(format string) (Exporter, error) {
	switch format {
	case "allure":
		return &AllureExporter{}, nil
	case "xray":
		return &XrayExporter{}, nil
	case "testrail":
		return &TestRailExporter{}, nil
	}
	return nil, fmt.Errorf("unknown export format %q — supported: allure, xray, testrail", format)
}

// ── Shared helpers ─────────────────────────────────────────────────────────────

// PriorityAllure maps a TestCase priority to an Allure severity label.
func PriorityAllure(p string) string {
	switch p {
	case "P0":
		return "blocker"
	case "P1":
		return "critical"
	case "P3":
		return "minor"
	default:
		return "normal"
	}
}

// PriorityXray maps a TestCase priority to a Jira Xray priority string.
func PriorityXray(p string) string {
	switch p {
	case "P0":
		return "Highest"
	case "P1":
		return "High"
	case "P3":
		return "Low"
	default:
		return "Medium"
	}
}

// PriorityTestRail maps a TestCase priority to a TestRail priority string.
func PriorityTestRail(p string) string {
	switch p {
	case "P0":
		return "Critical"
	case "P1":
		return "High"
	case "P3":
		return "Low"
	default:
		return "Medium"
	}
}

// AssertionsSummary formats a slice of Assertion into a human-readable string.
// Example: "status_code eq 201; duration_ms lt 2000"
func AssertionsSummary(assertions []schema.Assertion) string {
	if len(assertions) == 0 {
		return ""
	}
	var b strings.Builder
	for i, a := range assertions {
		if i > 0 {
			b.WriteString("; ")
		}
		fmt.Fprintf(&b, "%s %s %v", a.Target, a.Operator, a.Expected)
	}
	return b.String()
}

// toUUID derives a deterministic UUID-shaped string from an arbitrary string.
// Uses crypto/md5 so the same ID always produces the same UUID.
func toUUID(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x-%x-%x-%x-%x", h[0:4], h[4:6], h[6:8], h[8:10], h[10:16])
}
