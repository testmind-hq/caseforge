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

// ── Shared helpers ────────────────────────────────────────────────────────────

// PriorityAllure maps a TestCase priority to an Allure severity label.
func PriorityAllure(p string) string {
	switch p {
	case "P0":
		return "blocker"
	case "P1":
		return "critical"
	case "P2":
		return "normal"
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
	case "P2":
		return "Medium"
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
	case "P2":
		return "Medium"
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

func uuidFromString(s string) string {
	h := md5.Sum([]byte(s))
	// RFC 4122 UUID v3: set version nibble to 0x3 and variant bits to 0x8x
	b := h // copy
	b[6] = (b[6] & 0x0f) | 0x30
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ── Stubs (replaced one-by-one in Tasks 2, 3, 4) ────────────────────────────
// Each stub is removed ATOMICALLY when its implementation file is created.

type AllureExporter struct{}

func (e *AllureExporter) Format() string                              { return "allure" }
func (e *AllureExporter) Export(_ []schema.TestCase, _ string) error { return nil }

type XrayExporter struct{}

func (e *XrayExporter) Format() string                              { return "xray" }
func (e *XrayExporter) Export(_ []schema.TestCase, _ string) error { return nil }

type TestRailExporter struct{}

func (e *TestRailExporter) Format() string                              { return "testrail" }
func (e *TestRailExporter) Export(_ []schema.TestCase, _ string) error { return nil }
