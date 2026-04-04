// internal/rbt/callgraph_treesitter.go
package rbt

import (
	"os"
	"path/filepath"
	"strings"
)

// TreeSitterCallGraphBuilder uses the tree-sitter CLI to extract function
// definitions and call edges from source files.
type TreeSitterCallGraphBuilder struct{}

func NewTreeSitterCallGraphBuilder() *TreeSitterCallGraphBuilder {
	return &TreeSitterCallGraphBuilder{}
}

// ExtractFuncs returns function definitions and call edges for filePath.
// Returns empty slices (no error) if tree-sitter is not installed, the file
// does not exist, or the language is unsupported.
func (b *TreeSitterCallGraphBuilder) ExtractFuncs(filePath string) ([]CallNode, []CallEdge, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, nil, nil
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	lang, ok := langByExt[ext]
	if !ok {
		return nil, nil, nil
	}
	ts := NewTreeSitterParser()
	if !ts.IsAvailable() {
		return nil, nil, nil
	}

	defQuery := callGraphDefQueryForLang(lang)
	callQuery := callGraphCallQueryForLang(lang)
	if defQuery == "" || callQuery == "" {
		return nil, nil, nil
	}

	// Extract function definitions with line numbers.
	defResults, err := runCallGraphQuery(filePath, defQuery)
	if err != nil {
		return nil, nil, nil
	}

	var funcRanges []funcRange
	var defs []CallNode
	for _, r := range defResults {
		funcRanges = append(funcRanges, funcRange{r.name, r.line})
		defs = append(defs, CallNode{File: filePath, FuncName: r.name, Line: r.line})
	}

	// Extract call sites with line numbers.
	callResults, err := runCallGraphQuery(filePath, callQuery)
	if err != nil {
		return defs, nil, nil
	}

	// Attribute each call site to the nearest preceding function definition.
	var calls []CallEdge
	for _, cl := range callResults {
		callerFunc := nearestFunc(funcRanges, cl.line)
		if callerFunc == "" {
			continue
		}
		calls = append(calls, CallEdge{
			CallerFile: filePath,
			CallerFunc: callerFunc,
			CalleeName: cl.name,
		})
	}
	return defs, calls, nil
}

type captureResult struct {
	name string
	line int
}

// runCallGraphQuery runs a tree-sitter query and extracts (name, line) pairs.
func runCallGraphQuery(filePath, query string) ([]captureResult, error) {
	queryFile, err := os.CreateTemp("", "caseforge-cg-*.scm")
	if err != nil {
		return nil, err
	}
	defer os.Remove(queryFile.Name())
	if _, err := queryFile.WriteString(query); err != nil {
		return nil, err
	}
	queryFile.Close()

	out, err := runTreeSitterCmd(filePath, queryFile.Name())
	if err != nil {
		return nil, nil // tree-sitter error treated as no results
	}
	return parseCallGraphOutput(out), nil
}

// parseCallGraphOutput parses tree-sitter query output for any capture name.
// Expects lines like: `  @def: "FuncName" [row, col] - [row, col]`
func parseCallGraphOutput(output string) []captureResult {
	var results []captureResult
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		// Any capture annotation line starts with "@"
		if !strings.HasPrefix(line, "@") {
			continue
		}
		name := extractQuotedValue(line)
		lineNum := extractLineNum(line)
		if name != "" {
			results = append(results, captureResult{name, lineNum})
		}
	}
	return results
}

// nearestFunc returns the name of the function whose startLine is the largest
// value <= callLine. Returns "" if none found.
func nearestFunc(funcs []funcRange, callLine int) string {
	best := ""
	bestLine := -1
	for _, f := range funcs {
		if f.startLine <= callLine && f.startLine > bestLine {
			best = f.name
			bestLine = f.startLine
		}
	}
	return best
}

// callGraphDefQueryForLang returns the tree-sitter query for function definitions.
func callGraphDefQueryForLang(lang string) string {
	switch lang {
	case "go":
		return `[(function_declaration name: (identifier) @def)
 (method_declaration name: (field_identifier) @def)]`
	case "python":
		return `(function_definition name: (identifier) @def)`
	case "typescript", "javascript", "tsx":
		return `[(function_declaration name: (identifier) @def)
 (method_definition name: (property_identifier) @def)]`
	case "java":
		return `(method_declaration name: (identifier) @def)`
	case "rust":
		return `(function_item name: (identifier) @def)`
	default:
		return ""
	}
}

// callGraphCallQueryForLang returns the tree-sitter query for function calls.
func callGraphCallQueryForLang(lang string) string {
	switch lang {
	case "go":
		return `(call_expression function: [
 (identifier) @callee
 (selector_expression field: (field_identifier) @callee)])`
	case "python":
		return `(call function: [
 (identifier) @callee
 (attribute attribute: (identifier) @callee)])`
	case "typescript", "javascript", "tsx":
		return `(call_expression function: [
 (identifier) @callee
 (member_expression property: (property_identifier) @callee)])`
	case "java":
		return `(method_invocation name: (identifier) @callee)`
	case "rust":
		return `(call_expression function: [
 (identifier) @callee
 (field_expression field: (field_identifier) @callee)])`
	default:
		return ""
	}
}

type funcRange struct {
	name      string
	startLine int
}
