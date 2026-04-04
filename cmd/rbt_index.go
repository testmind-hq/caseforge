// cmd/rbt_index.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/rbt"
)

var rbtIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Auto-generate caseforge-map.yaml by analyzing source code",
	Long: `Index analyzes your project's source code and OpenAPI spec to automatically
generate caseforge-map.yaml, eliminating the need for manual maintenance.

Strategies:
  llm    — Analyze source with regex (full LLM inference is planned for v2)
  embed  — Use embeddings + cosine similarity (large projects, needs OPENAI_API_KEY)
  hybrid — Embed narrows candidates, LLM confirms (recommended for accuracy)

Examples:
  caseforge rbt index --spec openapi.yaml
  caseforge rbt index --spec openapi.yaml --strategy hybrid
  caseforge rbt index --spec openapi.yaml --out custom-map.yaml --overwrite`,
	RunE: runRBTIndex,
}

func init() {
	rbtCmd.AddCommand(rbtIndexCmd)
	rbtIndexCmd.Flags().String("spec", "", "OpenAPI spec file (required)")
	rbtIndexCmd.Flags().String("src", "./", "Source code root to analyze")
	rbtIndexCmd.Flags().String("out", "caseforge-map.yaml", "Output map file path")
	rbtIndexCmd.Flags().String("strategy", "llm", "Indexing strategy: llm|embed|hybrid")
	rbtIndexCmd.Flags().Bool("overwrite", false, "Overwrite existing map file")
	rbtIndexCmd.Flags().Int("depth", 0, "Call graph traversal depth (0 = dynamic, stop at route node)")
	rbtIndexCmd.Flags().String("algo", "rta", "Go call graph algorithm: rta or pta (default rta)")
	_ = rbtIndexCmd.MarkFlagRequired("spec")
}

func runRBTIndex(cmd *cobra.Command, _ []string) error {
	specPath, _ := cmd.Flags().GetString("spec")
	srcDir, _ := cmd.Flags().GetString("src")
	outPath, _ := cmd.Flags().GetString("out")
	strategy, _ := cmd.Flags().GetString("strategy")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	depth, _ := cmd.Flags().GetInt("depth")
	algo, _ := cmd.Flags().GetString("algo")

	out := cmd.OutOrStdout()

	if strategy == "llm" {
		fmt.Fprintln(out, "Note: LLM-based inference is planned for v2; using regex extraction.")
	}

	fmt.Fprintf(out, "Indexing %s (strategy: %s)...\n", srcDir, strategy)

	indexer := &rbt.Indexer{
		SrcDir:    srcDir,
		SpecPath:  specPath,
		OutPath:   outPath,
		Overwrite: overwrite,
		Store:     rbt.NewIndexStore(".caseforge-index"),
		Embedder:  rbt.NewOpenAIEmbedder(),
		Depth:     depth,
		Algo:      algo,
	}

	var err error
	switch strategy {
	case "embed", "hybrid":
		err = indexer.RunHybrid(nil)
	default: // "llm" and fallback
		err = indexer.RunRegex()
	}

	if err != nil {
		return fmt.Errorf("index: %w", err)
	}

	fmt.Fprintf(out, "✓ Map file written to %s\n", outPath)
	return nil
}
