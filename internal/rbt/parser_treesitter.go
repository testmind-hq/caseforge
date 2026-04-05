// internal/rbt/parser_treesitter.go
package rbt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gotreesitter "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
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

// langConstructors maps lang name to its gotreesitter Language constructor.
// Shared by parser_treesitter.go and callgraph_treesitter.go.
var langConstructors = map[string]func() *gotreesitter.Language{
	"go":         grammars.GoLanguage,
	"python":     grammars.PythonLanguage,
	"typescript": grammars.TypescriptLanguage,
	"tsx":        grammars.TsxLanguage,
	"javascript": grammars.JavascriptLanguage,
	"java":       grammars.JavaLanguage,
	"ruby":       grammars.RubyLanguage,
	"rust":       grammars.RustLanguage,
}

type treeSitterResult struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Line   int    `json:"line"`
}

type TreeSitterParser struct{}

func NewTreeSitterParser() *TreeSitterParser { return &TreeSitterParser{} }

func (p *TreeSitterParser) ExtractRoutes(_ context.Context, _ string, files []ChangedFile) ([]RouteMapping, error) {
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
	langFn, ok := langConstructors[lang]
	if !ok {
		return nil, nil
	}

	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	l := langFn()
	parser := gotreesitter.NewParser(l)
	tree, err := parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter parse: %w", err)
	}

	q, err := gotreesitter.NewQuery(query, l)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter query: %w", err)
	}

	var results []treeSitterResult
	for _, match := range q.Execute(tree) {
		var cur treeSitterResult
		hasMethod, hasPath := false, false
		for _, cap := range match.Captures {
			text := cap.Text(src)
			switch cap.Name {
			case "method":
				cur.Method = text
				hasMethod = true
			case "path":
				// interpreted_string_literal and Python string nodes include
				// surrounding quotes; string_fragment (TS/JS) does not.
				cur.Path = stripQuotes(text)
				cur.Line = int(cap.Node.StartPoint().Row)
				hasPath = true
			}
		}
		if hasMethod && hasPath {
			results = append(results, cur)
		}
	}
	return results, nil
}

// stripQuotes removes a single layer of surrounding double or single quotes.
func stripQuotes(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
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
