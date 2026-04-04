// internal/rbt/indexer.go
package rbt

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type mapFileOutput struct {
	Mappings []mapOutputEntry `yaml:"mappings"`
}

type mapOutputEntry struct {
	Source     string   `yaml:"source"`
	Operations []string `yaml:"operations"`
}

// Indexer orchestrates map file generation from source code.
type Indexer struct {
	SrcDir    string
	SpecPath  string
	OutPath   string
	Overwrite bool
	Store     *IndexStore
	Embedder  Embedder
	Depth     int    // 0 = dynamic BFS (stop at route node); >0 = fixed depth cap.
	               // Note: V3 (runGoCallGraphPhase) checks depth before terminal, so
	               // maxDepth=N covers N-1 hops; V2 (runCallGraphPhase) checks terminal
	               // first, so maxDepth=N covers N hops.
	Algo      string // Go call graph algorithm: "rta" (default) | "pta"
}

// RunRegex uses the regex parser to extract routes and writes caseforge-map.yaml.
func (idx *Indexer) RunRegex() error {
	if err := idx.checkOverwrite(); err != nil {
		return err
	}
	files, err := findSourceFiles(idx.SrcDir)
	if err != nil {
		return err
	}
	parser := NewRegexParser()
	mappings, err := parser.ExtractRoutes(context.Background(), idx.SrcDir, files)
	if err != nil {
		return err
	}
	return idx.writeMapFile(mappings, "regex")
}

// RunHybrid uses tree-sitter + Go call graph (V3) + name-based call graph (V2)
// + embeddings + LLM confirmation.
func (idx *Indexer) RunHybrid(llmParser *LLMParser) error {
	if err := idx.checkOverwrite(); err != nil {
		return err
	}
	files, err := findSourceFiles(idx.SrcDir)
	if err != nil {
		return err
	}

	// Phase 1: tree-sitter direct route detection.
	mappings, routeFileMappings := idx.runTreeSitterPhase(files)
	unclaimed := subtractFiles(files, mappings)

	// Phase 2: Go type-aware call graph (V3) — handles interface dispatch.
	goMappings, goClaimed := idx.runGoCallGraphPhase(unclaimed, routeFileMappings)
	mappings = append(mappings, goMappings...)
	unclaimed = subtractChangedFiles(unclaimed, goClaimed)

	// Phase 3: name-based call graph (V2) — covers non-Go files and Go fallback.
	cgMappings, cgClaimed := idx.runCallGraphPhase(files, unclaimed, routeFileMappings, llmParser)
	mappings = append(mappings, cgMappings...)
	unclaimed = subtractChangedFiles(unclaimed, cgClaimed)

	// Phase 4: embedding-based matching for remaining unclaimed files.
	embedMappings, err := idx.runEmbedPhase(unclaimed)
	if err == nil {
		mappings = append(mappings, embedMappings...)
	}

	return idx.writeMapFile(mappings, "hybrid")
}

// runTreeSitterPhase extracts route mappings using the tree-sitter parser.
// Returns mappings and a map of file → []RouteMapping for route-registering files.
func (idx *Indexer) runTreeSitterPhase(files []ChangedFile) ([]RouteMapping, map[string][]RouteMapping) {
	routeFileMappings := make(map[string][]RouteMapping)
	ts := NewTreeSitterParser()
	if !ts.IsAvailable() {
		return nil, routeFileMappings
	}
	mappings, err := ts.ExtractRoutes(context.Background(), idx.SrcDir, files)
	if err != nil {
		return nil, routeFileMappings
	}
	for _, rm := range mappings {
		routeFileMappings[rm.SourceFile] = append(routeFileMappings[rm.SourceFile], rm)
	}
	return mappings, routeFileMappings
}

// runGoCallGraphPhase uses golang.org/x/tools/go/callgraph to perform type-aware
// call graph analysis for Go modules. Silently returns empty on any error so V2
// can handle the unclaimed files.
func (idx *Indexer) runGoCallGraphPhase(
	unclaimed []ChangedFile,
	routeFileMappings map[string][]RouteMapping,
) ([]RouteMapping, []ChangedFile) {
	b := &GoCallGraphBuilder{SrcDir: idx.SrcDir, Algo: idx.Algo}
	mappings, claimed, _ := b.BuildAndTrace(unclaimed, routeFileMappings, idx.Depth)
	return mappings, claimed
}

// runCallGraphPhase builds a call graph from all source files and traces
// unclaimed files upward to route-registering files.
// Returns the found route mappings and the list of unclaimed files that were resolved.
func (idx *Indexer) runCallGraphPhase(
	allFiles []ChangedFile,
	unclaimed []ChangedFile,
	routeFileMappings map[string][]RouteMapping,
	llmParser *LLMParser,
) ([]RouteMapping, []ChangedFile) {
	tsBuilder := NewTreeSitterCallGraphBuilder()
	llmBuilder := NewLLMCallGraphBuilder(llmParser)
	builder := &fallbackCallGraphBuilder{primary: tsBuilder, fallback: llmBuilder}
	return idx.runCallGraphPhaseWithBuilder(allFiles, unclaimed, routeFileMappings, builder)
}

// runCallGraphPhaseWithBuilder is the testable core of runCallGraphPhase.
// Returns the found route mappings and the list of unclaimed files that were
// resolved (i.e., whose call chains reached a route-registering file).
func (idx *Indexer) runCallGraphPhaseWithBuilder(
	allFiles []ChangedFile,
	unclaimed []ChangedFile,
	routeFileMappings map[string][]RouteMapping,
	builder CallGraphBuilder,
) ([]RouteMapping, []ChangedFile) {
	if len(unclaimed) == 0 {
		return nil, nil
	}

	// BuildCallGraph returns defsByFile so we avoid a second ExtractFuncs pass.
	cg, defsByFile := BuildCallGraph(allFiles, builder)

	// Populate RouteNodes for spec compliance and future consumers.
	for filePath := range routeFileMappings {
		cg.RouteNodes = append(cg.RouteNodes, defsByFile[filePath]...)
	}

	// Collect start nodes from the already-extracted defs (no second call needed).
	var startNodes []CallNode
	for _, f := range unclaimed {
		startNodes = append(startNodes, defsByFile[f.Path]...)
	}
	if len(startNodes) == 0 {
		return nil, nil
	}

	// Choose via/confidence based on whether the LLM fallback was activated.
	via, confidence := "callgraph", 0.8
	if fb, ok := builder.(*fallbackCallGraphBuilder); ok && fb.hasUsedLLM {
		via, confidence = "callgraph-llm", 0.65
	}

	mappings, coveredFilesMap := TraceToRoutes(cg, startNodes, routeFileMappings, idx.Depth, via, confidence)

	// Build the list of claimed unclaimed files for the caller to subtract.
	var claimed []ChangedFile
	for _, f := range unclaimed {
		if coveredFilesMap[f.Path] {
			claimed = append(claimed, f)
		}
	}
	return mappings, claimed
}

func (idx *Indexer) runEmbedPhase(files []ChangedFile) ([]RouteMapping, error) {
	// Skip embedding I/O when using the no-op embedder (no API key configured).
	if _, isNoop := idx.Embedder.(*NoopEmbedder); !isNoop {
		localIdx, err := idx.Store.Load()
		if err != nil || localIdx == nil {
			localIdx = &LocalIndex{}
		}
		for _, f := range files {
			data, err := os.ReadFile(f.Path)
			if err != nil {
				continue
			}
			hash := fmt.Sprintf("%x", sha256.Sum256(data))
			if !isChunkStale(localIdx, f.Path, hash) {
				continue
			}
			emb, err := idx.Embedder.Embed(string(data))
			if err != nil {
				continue
			}
			localIdx.Chunks = append(localIdx.Chunks, IndexChunk{
				File:      f.Path,
				Fn:        filepath.Base(f.Path),
				Hash:      hash,
				Embedding: emb,
			})
		}
		_ = idx.Store.Save(localIdx)
	}
	// V1 stub: embeddings are stored for incremental re-embedding, but cosine similarity
	// → RouteMapping conversion (TopKChunks + LLM confirmation) is not yet implemented.
	// Fall back to regex for any unclaimed files to produce a useful map file.
	regexMappings, _ := NewRegexParser().ExtractRoutes(context.Background(), idx.SrcDir, files)
	return regexMappings, nil
}

func (idx *Indexer) checkOverwrite() error {
	if _, err := os.Stat(idx.OutPath); err == nil && !idx.Overwrite {
		return fmt.Errorf("%s already exists; use --overwrite to replace it", idx.OutPath)
	}
	return nil
}

func (idx *Indexer) writeMapFile(mappings []RouteMapping, strategy string) error {
	byFile := make(map[string][]string)
	for _, m := range mappings {
		op := m.Method + " " + m.RoutePath
		byFile[m.SourceFile] = append(byFile[m.SourceFile], op)
	}

	var entries []mapOutputEntry
	for file, ops := range byFile {
		seen := make(map[string]bool)
		var deduped []string
		for _, op := range ops {
			if !seen[op] {
				seen[op] = true
				deduped = append(deduped, op)
			}
		}
		entries = append(entries, mapOutputEntry{Source: file, Operations: deduped})
	}
	for i := range entries {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].Source > entries[j].Source {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	out := mapFileOutput{Mappings: entries}
	data, err := yaml.Marshal(out)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("# caseforge-map.yaml — generated by `caseforge rbt index`\n"+
		"# Strategy: %s | Indexed: %s\n"+
		"# Review entries before committing.\n",
		strategy, time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	return os.WriteFile(idx.OutPath, append([]byte(header), data...), 0644)
}

// findSourceFiles returns all supported source files under dir.
func findSourceFiles(dir string) ([]ChangedFile, error) {
	var files []ChangedFile
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := langByExt[ext]; ok {
			files = append(files, ChangedFile{Path: path})
		}
		return nil
	})
	return files, err
}

// isChunkStale returns true if the given file's hash differs from what's in idx.
func isChunkStale(idx *LocalIndex, file, newHash string) bool {
	if idx == nil {
		return true
	}
	for _, c := range idx.Chunks {
		if c.File == file {
			return c.Hash != newHash
		}
	}
	return true
}
