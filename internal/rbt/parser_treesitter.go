// internal/rbt/parser_treesitter.go
package rbt

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var langByExt = map[string]string{
	".go":   "go",
	".py":   "python",
	".ts":   "typescript",
	".tsx":  "tsx",
	".js":   "javascript",
	".java": "java",
	".rb":   "ruby",
	".rs":   "rust",
}

type treeSitterResult struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Line   int    `json:"line"`
}

type TreeSitterParser struct{}

func NewTreeSitterParser() *TreeSitterParser { return &TreeSitterParser{} }

func (p *TreeSitterParser) IsAvailable() bool {
	_, err := exec.LookPath("tree-sitter")
	return err == nil
}

func (p *TreeSitterParser) ExtractRoutes(srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	if !p.IsAvailable() {
		return nil, nil
	}

	var mappings []RouteMapping
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Path))
		lang, ok := langByExt[ext]
		if !ok {
			continue
		}
		routes, err := runTreeSitterQuery(f.Path, lang)
		if err != nil {
			continue
		}
		for _, r := range routes {
			mappings = append(mappings, RouteMapping{
				SourceFile: f.Path,
				Line:       r.Line,
				Method:     strings.ToUpper(r.Method),
				RoutePath:  r.Path,
				Via:        "treesitter",
				Confidence: 1.0,
			})
		}
	}
	return mappings, nil
}

func runTreeSitterQuery(filePath, lang string) ([]treeSitterResult, error) {
	query := routeQueryForLang(lang)
	if query == "" {
		return nil, nil
	}

	queryFile, err := os.CreateTemp("", "caseforge-ts-*.scm")
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
		return nil, err
	}

	return parseTSSexprOutput(out), nil
}

// runTreeSitterCmd runs tree-sitter query and returns raw stdout as a string.
func runTreeSitterCmd(filePath, queryFile string) (string, error) {
	cmd := exec.Command("tree-sitter", "query", queryFile, filePath)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("tree-sitter query: %w", err)
	}
	return string(out), nil
}

func parseTSSexprOutput(output string) []treeSitterResult {
	var results []treeSitterResult
	var current treeSitterResult
	hasMethod, hasPath := false, false

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Match:") {
			if hasMethod && hasPath {
				results = append(results, current)
			}
			current = treeSitterResult{}
			hasMethod, hasPath = false, false
			continue
		}
		if strings.HasPrefix(line, "@method:") {
			current.Method = extractQuotedValue(line)
			hasMethod = true
		}
		if strings.HasPrefix(line, "@path:") {
			current.Path = extractQuotedValue(line)
			lineNum := extractLineNum(line)
			if lineNum > 0 {
				current.Line = lineNum
			}
			hasPath = true
		}
	}
	if hasMethod && hasPath {
		results = append(results, current)
	}
	return results
}

func extractQuotedValue(line string) string {
	start := strings.Index(line, `"`)
	if start == -1 {
		return ""
	}
	end := strings.Index(line[start+1:], `"`)
	if end == -1 {
		return ""
	}
	return line[start+1 : start+1+end]
}

func extractLineNum(line string) int {
	start := strings.Index(line, "[")
	if start == -1 {
		return 0
	}
	rest := line[start+1:]
	end := strings.Index(rest, ",")
	if end == -1 {
		return 0
	}
	var n int
	fmt.Sscanf(strings.TrimSpace(rest[:end]), "%d", &n)
	return n
}

func routeQueryForLang(lang string) string {
	switch lang {
	case "go":
		return `(call_expression
  function: (selector_expression
    field: (field_identifier) @method
    (#match? @method "^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)$"))
  arguments: (argument_list
    (interpreted_string_literal) @path))`
	case "python":
		return `(decorator
  (call
    function: (attribute
      attribute: (identifier) @method
      (#match? @method "^(get|post|put|delete|patch)$"))
    arguments: (argument_list
      (string) @path)))`
	case "typescript", "javascript":
		return `(call_expression
  function: (member_expression
    property: (property_identifier) @method
    (#match? @method "^(get|post|put|delete|patch)$"))
  arguments: (arguments
    (string (string_fragment) @path)))`
	default:
		return ""
	}
}

