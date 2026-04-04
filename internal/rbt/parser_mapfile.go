// internal/rbt/parser_mapfile.go
package rbt

import (
	"context"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type mapFileConfig struct {
	Mappings []mapFileEntry `yaml:"mappings"`
}

type mapFileEntry struct {
	Source     string   `yaml:"source"`
	Operations []string `yaml:"operations"`
}

type MapFileParser struct {
	mapFilePath string
	config      *mapFileConfig
}

func NewMapFileParser(mapFilePath string) *MapFileParser {
	return &MapFileParser{mapFilePath: mapFilePath}
}

func (p *MapFileParser) load() (*mapFileConfig, error) {
	if p.config != nil {
		return p.config, nil
	}
	data, err := os.ReadFile(p.mapFilePath)
	if err != nil {
		return nil, nil
	}
	var cfg mapFileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	p.config = &cfg
	return p.config, nil
}

func (p *MapFileParser) ExtractRoutes(ctx context.Context, srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	cfg, err := p.load()
	if err != nil || cfg == nil {
		return nil, err
	}

	sourceOps := make(map[string][]string)
	for _, entry := range cfg.Mappings {
		sourceOps[entry.Source] = entry.Operations
	}

	var mappings []RouteMapping
	for _, f := range files {
		ops, ok := sourceOps[f.Path]
		if !ok {
			continue
		}
		for _, op := range ops {
			method, path, ok := parseOpString(op)
			if !ok {
				continue
			}
			mappings = append(mappings, RouteMapping{
				SourceFile: f.Path,
				Method:     method,
				RoutePath:  path,
				Via:        "mapfile",
				Confidence: 1.0,
			})
		}
	}
	return mappings, nil
}

func parseOpString(op string) (method, path string, ok bool) {
	parts := strings.Fields(op)
	if len(parts) < 2 {
		return "", "", false
	}
	return strings.ToUpper(parts[0]), parts[1], true
}
