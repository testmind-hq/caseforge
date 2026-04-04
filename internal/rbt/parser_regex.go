// internal/rbt/parser_regex.go
package rbt

import (
	"bufio"
	"context"
	"os"
	"regexp"
	"strings"
)

var httpMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

var routeRegex = regexp.MustCompile(
	`(?i)\.(?:` + strings.Join(httpMethods, "|") + `)\s*\(\s*["` + "`" + `]([/][^"` + "`" + `\s]*)["` + "`" + `]` +
		`|@(?:Get|Post|Put|Delete|Patch|Head|Options)Mapping\s*\(\s*"([/][^"]*)"`,
)

var methodExtract = regexp.MustCompile(`(?i)\.(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s*\(|@(Get|Post|Put|Delete|Patch|Head|Options)Mapping`)

type RegexParser struct{}

func NewRegexParser() *RegexParser { return &RegexParser{} }

func (p *RegexParser) ExtractRoutes(ctx context.Context, srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	var mappings []RouteMapping
	for _, f := range files {
		routes, err := extractRoutesFromFile(f.Path)
		if err != nil {
			continue
		}
		for _, r := range routes {
			r.SourceFile = f.Path
			mappings = append(mappings, r)
		}
	}
	return mappings, nil
}

func extractRoutesFromFile(path string) ([]RouteMapping, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var mappings []RouteMapping
	lineNum := 0
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matches := routeRegex.FindAllStringSubmatch(line, -1)
		if len(matches) == 0 {
			continue
		}
		methodMatches := methodExtract.FindAllStringSubmatch(line, -1)
		for i, match := range matches {
			routePath := match[1]
			if routePath == "" {
				routePath = match[2]
			}
			if routePath == "" {
				continue
			}

			method := ""
			if i < len(methodMatches) {
				m := methodMatches[i]
				method = strings.ToUpper(m[1] + m[2])
			}
			if method == "" {
				continue
			}

			mappings = append(mappings, RouteMapping{
				Line:       lineNum,
				Method:     method,
				RoutePath:  routePath,
				Via:        "regex",
				Confidence: 0.7,
			})
		}
	}
	return mappings, scanner.Err()
}
