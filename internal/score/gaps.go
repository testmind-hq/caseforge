package score

import (
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// Gap represents an operation that is missing happy-path (2xx) or
// invalid-input (4xx) test coverage.
type Gap struct {
	Method string
	Path   string
	Has2xx bool
	Has4xx bool
}

// ComputeGaps returns operations that lack 2xx or 4xx coverage.
// Operations with both 2xx and 4xx coverage are excluded.
func ComputeGaps(cases []schema.TestCase) []Gap {
	type coverage struct{ has2xx, has4xx bool }
	opCov := make(map[string]*coverage)

	for _, c := range cases {
		key := canonicalGapPath(c.Source.SpecPath)
		if key == "" {
			continue
		}
		if opCov[key] == nil {
			opCov[key] = &coverage{}
		}
		for _, step := range c.Steps {
			for _, a := range step.Assertions {
				if a.Target != "status_code" {
					continue
				}
				if gapIs2xx(a) {
					opCov[key].has2xx = true
				}
				if gapIs4xx(a) {
					opCov[key].has4xx = true
				}
			}
		}
	}

	var gaps []Gap
	for key, cov := range opCov {
		if cov.has2xx && cov.has4xx {
			continue
		}
		parts := strings.SplitN(key, " ", 2)
		if len(parts) != 2 {
			continue
		}
		gaps = append(gaps, Gap{
			Method: parts[0],
			Path:   parts[1],
			Has2xx: cov.has2xx,
			Has4xx: cov.has4xx,
		})
	}
	return gaps
}

// canonicalGapPath returns "METHOD /path" from a SpecPath string.
func canonicalGapPath(specPath string) string {
	parts := strings.Fields(specPath)
	if len(parts) < 2 {
		return ""
	}
	return parts[0] + " " + parts[1]
}

// gapIs2xx is a broader 2xx detector used for gap analysis.
// It recognises gte/lt/eq assertions that indicate a happy-path expectation.
func gapIs2xx(a schema.Assertion) bool {
	n, ok := toAssertInt(a.Expected)
	if !ok {
		return false
	}
	switch a.Operator {
	case schema.OperatorGte:
		return n >= 200 && n < 300
	case schema.OperatorLt:
		return n <= 300 && n > 0
	case schema.OperatorEq:
		return n >= 200 && n < 300
	}
	return false
}

// gapIs4xx is a broader 4xx detector used for gap analysis.
func gapIs4xx(a schema.Assertion) bool {
	n, ok := toAssertInt(a.Expected)
	if !ok {
		return false
	}
	switch a.Operator {
	case schema.OperatorGte:
		return n >= 400
	case schema.OperatorEq:
		return n >= 400 && n < 500
	}
	return false
}
