// internal/rbt/mapper.go
package rbt

import "context"

// SourceParser extracts route mappings from source files.
type SourceParser interface {
	ExtractRoutes(ctx context.Context, srcDir string, files []ChangedFile) ([]RouteMapping, error)
}

// MapChain runs parsers in order, passing only unclaimed files to each parser.
func MapChain(parsers []SourceParser, srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	return MapChainContext(context.Background(), parsers, srcDir, files)
}

// MapChainContext runs parsers in order with a caller-provided context.
func MapChainContext(ctx context.Context, parsers []SourceParser, srcDir string, files []ChangedFile) ([]RouteMapping, error) {
	if len(files) == 0 {
		return nil, nil
	}

	remaining := make([]ChangedFile, len(files))
	copy(remaining, files)

	var allMappings []RouteMapping

	for _, p := range parsers {
		if len(remaining) == 0 {
			break
		}
		mappings, err := p.ExtractRoutes(ctx, srcDir, remaining)
		if err != nil {
			return allMappings, err
		}
		if len(mappings) == 0 {
			continue
		}

		claimed := make(map[string]bool)
		for _, m := range mappings {
			claimed[m.SourceFile] = true
		}

		allMappings = append(allMappings, mappings...)

		var nextRemaining []ChangedFile
		for _, f := range remaining {
			if !claimed[f.Path] {
				nextRemaining = append(nextRemaining, f)
			}
		}
		remaining = nextRemaining
	}

	return allMappings, nil
}
