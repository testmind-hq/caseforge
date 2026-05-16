// internal/mutation/engine.go
package mutation

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/testmind-hq/caseforge/internal/runner"
)

type indexFile struct {
	TestCases []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Source struct {
			SpecPath string `json:"spec_path"`
		} `json:"source"`
	} `json:"test_cases"`
}

// Run executes the full mutation loop and returns the aggregate MutationRun.
func Run(opts RunOptions) (MutationRun, error) {
	opts.setDefaults()

	cases, err := loadCases(opts.CasesDir)
	if err != nil {
		return MutationRun{}, fmt.Errorf("loading index.json: %w", err)
	}
	if len(cases) == 0 {
		return MutationRun{}, fmt.Errorf("no test cases found in %s/index.json", opts.CasesDir)
	}

	operators := opts.Operators
	if len(operators) == 0 {
		operators = Registry()
	}

	var (
		mu      sync.Mutex
		results []CaseMutationResult
	)

	opSem := make(chan struct{}, opts.OperatorConcurrency)
	g, _ := errgroup.WithContext(context.Background())

	for _, op := range operators {
		g.Go(func() error {
			opSem <- struct{}{}
			defer func() { <-opSem }()

			proxy, err := NewProxy(opts.Target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warn: operator %s: proxy start failed: %v\n", op.Name(), err)
				return nil // skip this operator; don't abort others
			}
			defer proxy.Close()
			proxy.SetActive(op)

			caseSem := make(chan struct{}, opts.Concurrency)
			cg, _ := errgroup.WithContext(context.Background())
			for _, tc := range cases {
				cg.Go(func() error {
					caseSem <- struct{}{}
					defer func() { <-caseSem }()
					survived := runOnce(proxy.Addr(), opts.CasesDir, tc.ID)
					mu.Lock()
					results = append(results, CaseMutationResult{
						CaseID:    tc.ID,
						Title:     tc.Title,
						Operation: tc.Operation,
						Operator:  op.Name(),
						Survived:  survived,
					})
					mu.Unlock()
					return nil
				})
			}
			_ = cg.Wait()
			return nil
		})
	}
	// goroutines always return nil (proxy failures are logged and skipped)
	_ = g.Wait()

	killed, survivors := 0, 0
	for _, r := range results {
		if r.Survived {
			survivors++
		} else {
			killed++
		}
	}

	opNames := make([]string, len(operators))
	for i, op := range operators {
		opNames[i] = op.Name()
	}

	total := killed + survivors
	score := 0.0
	if total > 0 {
		score = float64(killed) / float64(total)
	}

	return MutationRun{
		Target:          opts.Target,
		CasesDir:        opts.CasesDir,
		Operators:       opNames,
		TotalCases:      len(cases),
		TotalRuns:       total,
		Killed:          killed,
		Survivors:       survivors,
		MutationScore:   score,
		Results:         results,
		OperationScores: ComputeOperationScores(results),
		GeneratedAt:     time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}

// runOnce runs a single .hurl file against the proxy.
// Returns true if the mutation survived (hurl passed = mutation not caught).
func runOnce(proxyAddr, casesDir, caseID string) bool {
	hurlFile := filepath.Join(casesDir, caseID+".hurl")
	if _, err := os.Stat(hurlFile); err != nil {
		return false // missing file → count as killed
	}

	tmpDir, err := os.MkdirTemp("", "caseforge-mutate-*")
	if err != nil {
		return false
	}
	defer os.RemoveAll(tmpDir)

	src, err := os.ReadFile(hurlFile)
	if err != nil {
		return false
	}
	if err := os.WriteFile(filepath.Join(tmpDir, caseID+".hurl"), src, 0644); err != nil {
		return false
	}

	r := runner.NewHurlRunner()
	result, err := r.Run(tmpDir, map[string]string{"base_url": "http://" + proxyAddr})
	if err != nil {
		return false
	}

	// Mutation survived if hurl passed (assertions didn't catch the mutation)
	return result.Passed > 0 && result.Failed == 0
}

type caseRef struct {
	ID        string
	Title     string
	Operation string // source.spec_path; "(unknown)" if absent
}

func loadCases(casesDir string) ([]caseRef, error) {
	data, err := os.ReadFile(filepath.Join(casesDir, "index.json"))
	if err != nil {
		return nil, err
	}
	var idx indexFile
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	refs := make([]caseRef, len(idx.TestCases))
	for i, tc := range idx.TestCases {
		op := tc.Source.SpecPath
		if op == "" {
			op = "(unknown)"
		}
		refs[i] = caseRef{ID: tc.ID, Title: tc.Title, Operation: op}
	}
	return refs, nil
}
