// internal/rbt/generator.go
package rbt

import (
	"fmt"
	"io"
)

// GenerateForHighRisk generates test cases for all HIGH-risk operations.
// This is a stub for V1 — full methodology pipeline integration in a future iteration.
func GenerateForHighRisk(w io.Writer, report RiskReport, casesDir string) error {
	var highRisk []OperationCoverage
	for _, op := range report.Operations {
		if op.Risk == RiskHigh {
			highRisk = append(highRisk, op)
		}
	}
	if len(highRisk) == 0 {
		fmt.Fprintln(w, "No HIGH-risk operations to generate tests for.")
		return nil
	}
	for _, op := range highRisk {
		fmt.Fprintf(w, "⚠  --generate not yet implemented for %s %s\n", op.Method, op.Path)
	}
	return nil
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
