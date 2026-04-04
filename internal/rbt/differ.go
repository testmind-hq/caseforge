// internal/rbt/differ.go
package rbt

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Diff returns the list of files changed between base and head in the git
// repo at repoDir. Returns empty (no error) if repoDir is not a git repo.
func Diff(repoDir, base, head string) ([]ChangedFile, error) {
	// Check if it's a git repo
	check := exec.Command("git", "rev-parse", "--git-dir")
	check.Dir = repoDir
	if err := check.Run(); err != nil {
		return nil, nil // not a git repo
	}

	// Get changed file list (Added, Modified, Deleted)
	listCmd := exec.Command("git", "diff", "--name-status", "--diff-filter=AMD",
		base+".."+head)
	listCmd.Dir = repoDir
	out, err := listCmd.Output()
	if err != nil {
		// base ref may not exist (e.g. shallow clone, first commit); return empty
		return nil, nil
	}

	var files []ChangedFile
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		status, path := parts[0], parts[1]
		cf := ChangedFile{Path: path}
		switch {
		case strings.HasPrefix(status, "A"):
			cf.IsNew = true
		case strings.HasPrefix(status, "D"):
			cf.IsDeleted = true
		}
		if !cf.IsDeleted {
			cf.ChangedLines = changedLines(repoDir, base, head, path)
		}
		files = append(files, cf)
	}
	return files, scanner.Err()
}

// changedLines returns the line numbers added/changed in path between base and head.
func changedLines(repoDir, base, head, path string) []int {
	cmd := exec.Command("git", "diff", "-U0", base+".."+head, "--", path)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parseHunkLines(string(out))
}

// parseHunkLines parses @@ -a,b +c,d @@ hunk headers and returns new-file line numbers.
func parseHunkLines(diff string) []int {
	var lines []int
	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "@@") {
			continue
		}
		// Format: @@ -old +new[,count] @@
		// Extract the +N[,M] part
		parts := strings.Fields(line)
		for _, p := range parts {
			if !strings.HasPrefix(p, "+") || p == "+++ " {
				continue
			}
			p = strings.TrimPrefix(p, "+")
			startStr, countStr, hasComma := strings.Cut(p, ",")
			start, err := strconv.Atoi(startStr)
			if err != nil {
				continue
			}
			count := 1
			if hasComma {
				count, _ = strconv.Atoi(countStr)
			}
			for i := 0; i < count; i++ {
				lines = append(lines, start+i)
			}
			break
		}
	}
	// deduplicate
	seen := make(map[int]bool)
	var deduped []int
	for _, l := range lines {
		if !seen[l] {
			seen[l] = true
			deduped = append(deduped, l)
		}
	}
	return deduped
}

// AbsChangedFiles resolves ChangedFile.Path to absolute paths given repoDir.
func AbsChangedFiles(files []ChangedFile, repoDir string) []ChangedFile {
	out := make([]ChangedFile, len(files))
	for i, f := range files {
		out[i] = f
		out[i].Path = filepath.Join(repoDir, f.Path)
	}
	return out
}
