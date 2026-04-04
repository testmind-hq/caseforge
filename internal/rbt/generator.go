// internal/rbt/generator.go
package rbt

import (
	"strings"

	specpkg "github.com/testmind-hq/caseforge/internal/spec"
)

// HighRiskOperations returns the spec operations that are assessed as RiskHigh
// in report, preserving the order they appear in parsedSpec.Operations.
// Returns nil when no operations have RiskHigh.
func HighRiskOperations(report RiskReport, parsedSpec *specpkg.ParsedSpec) []*specpkg.Operation {
	highRiskKeys := make(map[string]bool, len(report.Operations))
	for _, oc := range report.Operations {
		if oc.Risk == RiskHigh {
			key := strings.ToUpper(oc.Method) + " " + oc.Path
			highRiskKeys[key] = true
		}
	}
	if len(highRiskKeys) == 0 {
		return nil
	}
	var ops []*specpkg.Operation
	for _, op := range parsedSpec.Operations {
		key := strings.ToUpper(op.Method) + " " + op.Path
		if highRiskKeys[key] {
			ops = append(ops, op)
		}
	}
	return ops
}

// MatchesOperation checks whether a RouteMapping's method+path matches a spec operation.
func MatchesOperation(rm RouteMapping, method, specPath string) bool {
	if rm.Method != method {
		return false
	}
	return normalizePathParams(rm.RoutePath) == normalizePathParams(specPath)
}

// normalizePathParams converts :param style to {param} style.
func normalizePathParams(path string) string {
	out := make([]byte, 0, len(path))
	i := 0
	for i < len(path) {
		if path[i] == ':' {
			out = append(out, '{')
			i++
			for i < len(path) && path[i] != '/' {
				out = append(out, path[i])
				i++
			}
			out = append(out, '}')
		} else {
			out = append(out, path[i])
			i++
		}
	}
	return string(out)
}
