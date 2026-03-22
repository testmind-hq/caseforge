// internal/diff/diff.go
package diff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/testmind-hq/caseforge/internal/spec"
)

// ChangeKind classifies the impact of a spec change on API consumers.
type ChangeKind string

const (
	Breaking            ChangeKind = "BREAKING"
	PotentiallyBreaking ChangeKind = "POTENTIALLY_BREAKING"
	NonBreaking         ChangeKind = "NON_BREAKING"
)

// Change describes a single detected difference between two spec versions.
type Change struct {
	Kind        ChangeKind `json:"kind"`
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	Location    string     `json:"location,omitempty"` // "requestBody", "response.200", "param.limit", etc.
	Description string     `json:"description"`
}

// DiffResult holds all detected changes.
type DiffResult struct {
	Changes []Change
}

// HasBreaking returns true if any change is Breaking or PotentiallyBreaking.
func (r DiffResult) HasBreaking() bool {
	for _, c := range r.Changes {
		if c.Kind == Breaking || c.Kind == PotentiallyBreaking {
			return true
		}
	}
	return false
}

// Diff compares oldSpec against newSpec and returns all detected changes.
func Diff(oldSpec, newSpec *spec.ParsedSpec) DiffResult {
	var changes []Change

	oldOps := opMap(oldSpec)
	newOps := opMap(newSpec)

	// Detect removed endpoints, param changes, response field changes
	for key, oldOp := range oldOps {
		newOp, exists := newOps[key]
		if !exists {
			changes = append(changes, Change{
				Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
				Description: "endpoint removed",
			})
			continue
		}
		changes = append(changes, diffParams(oldOp, newOp)...)
		changes = append(changes, diffResponseFields(oldOp, newOp)...)
		changes = append(changes, diffRequestBodyFields(oldOp, newOp)...)
	}

	// Detect path renames (BREAKING) and new endpoints (NON_BREAKING)
	uniqueOldPaths := sortedUniqueNewPaths(oldSpec, newSpec)
	uniqueNewPaths := uniqueNewPathsSet(oldSpec, newSpec)
	renamedNewPaths := map[string]bool{}

	for _, oldPath := range uniqueOldPaths {
		best, diffCount := findBestRename(oldPath, uniqueNewPaths)
		if best != "" && diffCount > 0 {
			changes = append(changes, Change{
				Kind: Breaking, Path: oldPath,
				Description: fmt.Sprintf("path renamed: %s → %s", oldPath, best),
			})
			renamedNewPaths[best] = true
		}
	}

	// New endpoints (not renames)
	for _, newPath := range sortedSlice(uniqueNewPaths) {
		if renamedNewPaths[newPath] {
			continue
		}
		// find all ops on this new path
		for _, newOp := range newSpec.Operations {
			if newOp.Path == newPath {
				changes = append(changes, Change{
					Kind: NonBreaking, Method: newOp.Method, Path: newPath,
					Description: "new endpoint",
				})
			}
		}
	}

	return DiffResult{Changes: changes}
}

// opMap builds a "METHOD /path" → *Operation map.
func opMap(ps *spec.ParsedSpec) map[string]*spec.Operation {
	m := map[string]*spec.Operation{}
	for _, op := range ps.Operations {
		m[op.Method+" "+op.Path] = op
	}
	return m
}

func diffParams(oldOp, newOp *spec.Operation) []Change {
	var changes []Change
	oldParams := paramMap(oldOp.Parameters)
	newParams := paramMap(newOp.Parameters)

	for name, oldP := range oldParams {
		newP, exists := newParams[name]
		if !exists {
			// Parameter removed — could break clients sending it; treat as POTENTIALLY_BREAKING
			changes = append(changes, Change{
				Kind: PotentiallyBreaking, Method: oldOp.Method, Path: oldOp.Path,
				Location:    "param." + name,
				Description: fmt.Sprintf("parameter %q removed", name),
			})
			continue
		}
		// Required flag: optional → required is BREAKING
		if !oldP.Required && newP.Required {
			changes = append(changes, Change{
				Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
				Location:    "param." + name,
				Description: fmt.Sprintf("parameter %q changed from optional to required", name),
			})
		}
		// Type change is BREAKING
		if oldP.Schema != nil && newP.Schema != nil && oldP.Schema.Type != newP.Schema.Type {
			changes = append(changes, Change{
				Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
				Location:    "param." + name,
				Description: fmt.Sprintf("parameter %q type changed: %s → %s", name, oldP.Schema.Type, newP.Schema.Type),
			})
		}
	}

	// New optional params are NON_BREAKING; new required params are POTENTIALLY_BREAKING
	for name, newP := range newParams {
		if _, exists := oldParams[name]; !exists {
			kind := NonBreaking
			desc := fmt.Sprintf("new optional parameter %q added", name)
			if newP.Required {
				kind = PotentiallyBreaking
				desc = fmt.Sprintf("new required parameter %q added", name)
			}
			changes = append(changes, Change{
				Kind: kind, Method: oldOp.Method, Path: oldOp.Path,
				Location: "param." + name, Description: desc,
			})
		}
	}
	return changes
}

func diffResponseFields(oldOp, newOp *spec.Operation) []Change {
	var changes []Change
	for code, oldResp := range oldOp.Responses {
		newResp, exists := newOp.Responses[code]
		if !exists {
			// Entire response code removed — treat as BREAKING
			_ = oldResp
			changes = append(changes, Change{
				Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
				Location:    fmt.Sprintf("response.%s", code),
				Description: fmt.Sprintf("response code %s removed", code),
			})
			continue
		}
		oldSchema := responseJSONSchema(oldResp)
		newSchema := responseJSONSchema(newResp)
		if oldSchema == nil || newSchema == nil {
			continue
		}
		for fieldName, oldField := range oldSchema.Properties {
			newField, exists := newSchema.Properties[fieldName]
			if !exists {
				changes = append(changes, Change{
					Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
					Location:    fmt.Sprintf("response.%s", code),
					Description: fmt.Sprintf("response field %q removed", fieldName),
				})
				continue
			}
			if oldField.Type != newField.Type && newField.Type != "" {
				// Type widening (integer→number) is POTENTIALLY_BREAKING; other changes are BREAKING
				kind := Breaking
				if oldField.Type == "integer" && newField.Type == "number" {
					kind = PotentiallyBreaking
				}
				changes = append(changes, Change{
					Kind: kind, Method: oldOp.Method, Path: oldOp.Path,
					Location:    fmt.Sprintf("response.%s", code),
					Description: fmt.Sprintf("response field %q type changed: %s → %s", fieldName, oldField.Type, newField.Type),
				})
			}
		}
		// New response fields are NON_BREAKING
		for fieldName := range newSchema.Properties {
			if _, exists := oldSchema.Properties[fieldName]; !exists {
				changes = append(changes, Change{
					Kind: NonBreaking, Method: oldOp.Method, Path: oldOp.Path,
					Location:    fmt.Sprintf("response.%s", code),
					Description: fmt.Sprintf("new response field %q added", fieldName),
				})
			}
		}
	}
	return changes
}

func diffRequestBodyFields(oldOp, newOp *spec.Operation) []Change {
	var changes []Change
	oldSchema := requestJSONSchema(oldOp)
	newSchema := requestJSONSchema(newOp)
	if oldSchema == nil || newSchema == nil {
		return nil
	}
	// New required fields in request body are POTENTIALLY_BREAKING
	oldRequired := stringSet(oldSchema.Required)
	for _, req := range newSchema.Required {
		if !oldRequired[req] {
			changes = append(changes, Change{
				Kind: PotentiallyBreaking, Method: oldOp.Method, Path: oldOp.Path,
				Location:    "requestBody",
				Description: fmt.Sprintf("new required request body field %q added", req),
			})
		}
	}
	// Field type changes in request body are BREAKING
	for fieldName, oldField := range oldSchema.Properties {
		if newField, exists := newSchema.Properties[fieldName]; exists {
			if oldField.Type != newField.Type && newField.Type != "" {
				changes = append(changes, Change{
					Kind: Breaking, Method: oldOp.Method, Path: oldOp.Path,
					Location:    "requestBody",
					Description: fmt.Sprintf("request body field %q type changed: %s → %s", fieldName, oldField.Type, newField.Type),
				})
			}
		}
	}
	return changes
}

// --- path rename helpers ---

func sortedUniqueNewPaths(oldSpec, newSpec *spec.ParsedSpec) []string {
	oldPaths := pathSet(oldSpec)
	newPaths := pathSet(newSpec)
	var result []string
	for p := range oldPaths {
		if !newPaths[p] {
			result = append(result, p)
		}
	}
	sort.Strings(result)
	return result
}

func uniqueNewPathsSet(oldSpec, newSpec *spec.ParsedSpec) map[string]bool {
	oldPaths := pathSet(oldSpec)
	newPaths := pathSet(newSpec)
	result := map[string]bool{}
	for p := range newPaths {
		if !oldPaths[p] {
			result[p] = true
		}
	}
	return result
}

func findBestRename(oldPath string, candidates map[string]bool) (string, int) {
	oldSegs := splitPath(oldPath)
	best := ""
	bestDiff := -1
	for cand := range candidates {
		newSegs := splitPath(cand)
		if len(oldSegs) != len(newSegs) {
			continue
		}
		paramSame := true
		diffCount := 0
		for i := range oldSegs {
			oldIsParam := isParamSeg(oldSegs[i])
			newIsParam := isParamSeg(newSegs[i])
			if oldIsParam != newIsParam {
				paramSame = false
				break
			}
			if !oldIsParam && !newIsParam && oldSegs[i] != newSegs[i] {
				diffCount++
			}
		}
		// Require exactly 1 differing segment to avoid false-positive renames
		// between semantically unrelated paths (e.g. /users removed + /orders added).
		if !paramSame || diffCount == 0 || diffCount > 1 {
			continue
		}
		if bestDiff < 0 || diffCount < bestDiff || (diffCount == bestDiff && cand < best) {
			bestDiff = diffCount
			best = cand
		}
	}
	return best, bestDiff
}

// --- small helpers ---

func paramMap(params []*spec.Parameter) map[string]*spec.Parameter {
	m := map[string]*spec.Parameter{}
	for _, p := range params {
		m[p.Name] = p
	}
	return m
}

func pathSet(ps *spec.ParsedSpec) map[string]bool {
	m := map[string]bool{}
	for _, op := range ps.Operations {
		m[op.Path] = true
	}
	return m
}

func responseJSONSchema(resp *spec.Response) *spec.Schema {
	if resp == nil {
		return nil
	}
	if mt, ok := resp.Content["application/json"]; ok {
		return mt.Schema
	}
	return nil
}

func requestJSONSchema(op *spec.Operation) *spec.Schema {
	if op.RequestBody == nil {
		return nil
	}
	if mt, ok := op.RequestBody.Content["application/json"]; ok {
		return mt.Schema
	}
	return nil
}

func stringSet(ss []string) map[string]bool {
	m := map[string]bool{}
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func splitPath(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

func isParamSeg(seg string) bool {
	return strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}")
}

func sortedSlice(m map[string]bool) []string {
	var result []string
	for k := range m {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}
