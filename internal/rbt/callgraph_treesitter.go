// internal/rbt/callgraph_treesitter.go
package rbt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gotreesitter "github.com/odvcencio/gotreesitter"
)

// TreeSitterCallGraphBuilder uses the gotreesitter library to extract function
// definitions and call edges from source files.
type TreeSitterCallGraphBuilder struct{}

func NewTreeSitterCallGraphBuilder() *TreeSitterCallGraphBuilder {
	return &TreeSitterCallGraphBuilder{}
}

// ExtractFuncs returns function definitions and call edges for filePath.
// Returns empty slices (no error) if the file does not exist or the language
// is unsupported.
func (b *TreeSitterCallGraphBuilder) ExtractFuncs(filePath string) ([]CallNode, []CallEdge, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, nil, nil
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	lang, ok := langByExt[ext]
	if !ok {
		return nil, nil, nil
	}
	langFn, ok := langConstructors[lang]
	if !ok {
		return nil, nil, nil
	}
	l := langFn()

	defQuery := callGraphDefQueryForLang(lang)
	callQuery := callGraphCallQueryForLang(lang)
	if defQuery == "" || callQuery == "" {
		return nil, nil, nil
	}

	// Read and parse the file once; share src and tree across both queries.
	src, tree, err := parseFile(filePath, l)
	if err != nil {
		return nil, nil, nil
	}

	// Extract function definitions with line numbers.
	defResults, err := runCallGraphQuery(src, tree, defQuery, l)
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
	callResults, err := runCallGraphQuery(src, tree, callQuery, l)
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

// parseFile reads filePath and parses it with l, returning the source bytes
// and parsed tree. Returns an error if reading or parsing fails.
func parseFile(filePath string, l *gotreesitter.Language) ([]byte, *gotreesitter.Tree, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	tree, err := gotreesitter.NewParser(l).Parse(src)
	if err != nil {
		return nil, nil, fmt.Errorf("tree-sitter parse %s: %w", filePath, err)
	}
	return src, tree, nil
}

// runCallGraphQuery runs query against an already-parsed tree, returning
// (name, line) pairs for every capture in every match.
func runCallGraphQuery(src []byte, tree *gotreesitter.Tree, query string, l *gotreesitter.Language) ([]captureResult, error) {
	q, err := gotreesitter.NewQuery(query, l)
	if err != nil {
		return nil, nil
	}
	var results []captureResult
	for _, match := range q.Execute(tree) {
		for _, cap := range match.Captures {
			results = append(results, captureResult{
				name: cap.Text(src),
				line: int(cap.Node.StartPoint().Row),
			})
		}
	}
	return results, nil
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
