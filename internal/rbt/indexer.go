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

	specpkg "github.com/testmind-hq/caseforge/internal/spec"
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
	               // maxDepth=N means "traverse at most N hops from the seeded functions".
	               // Both V2 (runCallGraphPhase) and V3 (runGoCallGraphPhase) check terminal
	               // before depth cap, so route files at exactly depth=N are recorded.
	Algo      string // Go call graph algorithm: "rta" (default) | "vta"
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
	// runEmbedPhase never returns an error (all failures are swallowed internally).
	embedMappings, _ := idx.runEmbedPhase(unclaimed)
	mappings = append(mappings, embedMappings...)

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

// runEmbedPhase embeds unclaimed source files and spec operations, then uses
// cosine similarity (topKAboveThreshold) to produce RouteMapping entries. Falls back
// to the regex parser when the embedder is noop (no API key) or when no
// mappings could be derived from the index.
func (idx *Indexer) runEmbedPhase(files []ChangedFile) ([]RouteMapping, error) {
	// Noop embedder: no API key configured — skip embedding, fall back to regex.
	if _, isNoop := idx.Embedder.(*NoopEmbedder); isNoop {
		regexMappings, _ := NewRegexParser().ExtractRoutes(context.Background(), idx.SrcDir, files)
		return regexMappings, nil
	}

	// Load or create LocalIndex.
	localIdx, err := idx.Store.Load()
	if err != nil || localIdx == nil {
		localIdx = &LocalIndex{}
	}

	// Phase 1: embed file chunks (incremental — skip unchanged files).
	for _, f := range files {
		data, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		hash := fmt.Sprintf("%x", sha256.Sum256(data))
		if !isChunkStale(localIdx, f.Path, hash) {
			continue // already cached
		}
		emb, err := idx.Embedder.Embed(string(data))
		if err != nil {
			continue
		}
		// Replace existing entry for this file or append a new one.
		// Fn is set to the base filename (e.g. "service.go") as a V1
		// approximation; function-level chunking via tree-sitter is planned.
		replaced := false
		for i, c := range localIdx.Chunks {
			if c.File == f.Path {
				localIdx.Chunks[i] = IndexChunk{File: f.Path, Fn: filepath.Base(f.Path), Hash: hash, Embedding: emb}
				replaced = true
				break
			}
		}
		if !replaced {
			localIdx.Chunks = append(localIdx.Chunks, IndexChunk{
				File: f.Path, Fn: filepath.Base(f.Path), Hash: hash, Embedding: emb,
			})
		}
	}

	// Phase 2: embed spec operations (incremental — skip already-cached ops).
	if idx.SpecPath != "" {
		if parsed, loadErr := specpkg.NewLoader().Load(idx.SpecPath); loadErr == nil {
			for _, op := range parsed.Operations {
				opKey := strings.ToUpper(op.Method) + " " + op.Path
				if !isSpecOpStale(localIdx, opKey) {
					continue // already cached
				}
				text := fmt.Sprintf("%s %s — %s", strings.ToUpper(op.Method), op.Path, op.Summary)
				emb, embErr := idx.Embedder.Embed(text)
				if embErr != nil {
					continue
				}
				replaced := false
				for i, s := range localIdx.SpecOps {
					if s.Operation == opKey {
						localIdx.SpecOps[i] = IndexSpecOp{Operation: opKey, Description: op.Summary, Embedding: emb}
						replaced = true
						break
					}
				}
				if !replaced {
					localIdx.SpecOps = append(localIdx.SpecOps, IndexSpecOp{
						Operation: opKey, Description: op.Summary, Embedding: emb,
					})
				}
			}
		}
	}

	// Persist updated index for future incremental runs.
	_ = idx.Store.Save(localIdx)

	// Phase 3: cosine-similarity matching — for each spec op, find the top-k
	// most similar source file chunks (above a minimum similarity threshold)
	// and emit a RouteMapping per match.
	//
	// Only chunks from the files passed to this call are considered. The
	// persisted index may contain chunks from prior runs (already claimed by
	// tree-sitter or call-graph phases); restricting to the current `files`
	// slice avoids emitting spurious "embed" mappings for those files.
	const (
		topK         = 3
		minSimilarity = 0.3 // discard near-orthogonal matches
	)
	filesSet := make(map[string]bool, len(files))
	for _, f := range files {
		filesSet[f.Path] = true
	}
	var fileBoundChunks []IndexChunk
	for _, c := range localIdx.Chunks {
		if filesSet[c.File] {
			fileBoundChunks = append(fileBoundChunks, c)
		}
	}

	seen := make(map[string]bool) // deduplicate (sourceFile, operation) pairs
	var mappings []RouteMapping
	for _, op := range localIdx.SpecOps {
		if len(op.Embedding) == 0 {
			continue
		}
		parts := strings.SplitN(op.Operation, " ", 2)
		if len(parts) != 2 {
			continue
		}
		method, routePath := parts[0], parts[1]
		for _, chunk := range topKAboveThreshold(op.Embedding, fileBoundChunks, topK, minSimilarity) {
			dedupKey := chunk.File + "|" + op.Operation
			if seen[dedupKey] {
				continue
			}
			seen[dedupKey] = true
			mappings = append(mappings, RouteMapping{
				SourceFile: chunk.File,
				Method:     method,
				RoutePath:  routePath,
				Via:        "embed",
				Confidence: 0.65,
			})
		}
	}

	if len(mappings) == 0 {
		// No spec ops loaded or no chunks available — fall back to regex.
		regexMappings, _ := NewRegexParser().ExtractRoutes(context.Background(), idx.SrcDir, files)
		return regexMappings, nil
	}

	return mappings, nil
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

// isSpecOpStale returns true if opKey is not yet cached in idx or has an empty
// embedding. Unlike isChunkStale (which is hash-based), spec ops are considered
// fresh as long as an embedding exists — there is no content hash to version by.
// If an operation's summary changes, delete index.json to force re-embedding.
func isSpecOpStale(idx *LocalIndex, opKey string) bool {
	if idx == nil {
		return true
	}
	for _, s := range idx.SpecOps {
		if s.Operation == opKey {
			return len(s.Embedding) == 0
		}
	}
	return true
}
