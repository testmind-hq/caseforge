// internal/dedupe/merger.go
package dedupe

import (
	"fmt"
	"io"
	"os"
)

// MergeOptions controls how Merge behaves.
type MergeOptions struct {
	DryRun bool
	Out    io.Writer
}

// Merge deletes every non-kept file in each DuplicateGroup.
// Returns the count of files deleted (or would-be-deleted in dry-run).
func Merge(groups []DuplicateGroup, opts MergeOptions) (int, error) {
	count := 0
	for _, g := range groups {
		for _, cs := range g.Cases {
			if cs.Keep {
				continue
			}
			if opts.DryRun {
				fmt.Fprintf(opts.Out, "[dry-run] would delete: %s\n", cs.FilePath)
				count++
				continue
			}
			if err := os.Remove(cs.FilePath); err != nil && !os.IsNotExist(err) {
				return count, fmt.Errorf("delete %s: %w", cs.FilePath, err)
			}
			fmt.Fprintf(opts.Out, "deleted: %s\n", cs.FilePath)
			count++
		}
	}
	return count, nil
}
