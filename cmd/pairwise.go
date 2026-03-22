// cmd/pairwise.go
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/testmind-hq/caseforge/internal/methodology"
)

var pairwiseCmd = &cobra.Command{
	Use:   "pairwise",
	Short: "Compute pairwise combinations for given parameters",
	Long:  `Compute a 2-way covering array (IPOG algorithm) for the given parameters.`,
	RunE:  runPairwise,
}

var pairwiseParams string

func init() {
	rootCmd.AddCommand(pairwiseCmd)
	pairwiseCmd.Flags().StringVar(&pairwiseParams, "params",
		"", `Parameters in "name:v1,v2 name2:v3,v4" format (required)`)
	_ = pairwiseCmd.MarkFlagRequired("params")
}

func runPairwise(cmd *cobra.Command, args []string) error {
	params, err := parsePairwiseParams(pairwiseParams)
	if err != nil {
		return fmt.Errorf("parsing params: %w", err)
	}
	rows := methodology.IPOG(params)
	out, _ := json.MarshalIndent(rows, "", "  ")
	fmt.Println(string(out))
	fmt.Fprintf(cmd.ErrOrStderr(), "✓ %d combinations (full factorial: %d)\n",
		len(rows), fullFactorial(params))
	return nil
}

func parsePairwiseParams(input string) ([]methodology.PairwiseParam, error) {
	var params []methodology.PairwiseParam
	for _, part := range strings.Fields(input) {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("expected name:v1,v2 format, got %q", part)
		}
		vals := strings.Split(kv[1], ",")
		anyVals := make([]any, len(vals))
		for i, v := range vals {
			anyVals[i] = v
		}
		params = append(params, methodology.PairwiseParam{Name: kv[0], Values: anyVals})
	}
	if len(params) < 2 {
		return nil, fmt.Errorf("need at least 2 parameters")
	}
	return params, nil
}

func fullFactorial(params []methodology.PairwiseParam) int {
	n := 1
	for _, p := range params {
		n *= len(p.Values)
	}
	return n
}
