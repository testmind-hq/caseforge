// internal/rbt/differ_test.go
package rbt

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupGitRepo creates a temp git repo with an initial commit.
func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test.com")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "git command failed: %s\n%s", args, out)
	}
	run("git", "init")
	run("git", "checkout", "-b", "main")
	// initial commit with a file
	require.NoError(t, os.WriteFile(filepath.Join(dir, "foo.go"), []byte("package main\n"), 0644))
	run("git", "add", ".")
	run("git", "commit", "-m", "init")
	return dir
}

func TestDiff_NoGitRepo_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir() // not a git repo
	files, err := Diff(dir, "HEAD~1", "HEAD")
	require.NoError(t, err)
	assert.Empty(t, files)
}

func TestDiff_NewFile(t *testing.T) {
	dir := setupGitRepo(t)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "bar.go"), []byte("package main\n"), 0644))
	shCmd := exec.Command("sh", "-c", "git add . && git commit -m 'add bar'")
	shCmd.Dir = dir
	shCmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test.com")
	require.NoError(t, shCmd.Run())

	files, err := Diff(dir, "HEAD~1", "HEAD")
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "bar.go", files[0].Path)
	assert.True(t, files[0].IsNew)
	assert.False(t, files[0].IsDeleted)
}

func TestDiff_ModifiedFile(t *testing.T) {
	dir := setupGitRepo(t)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "foo.go"), []byte("package main\nfunc A(){}\n"), 0644))
	shCmd := exec.Command("sh", "-c", "git add . && git commit -m 'modify foo'")
	shCmd.Dir = dir
	shCmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test.com")
	require.NoError(t, shCmd.Run())

	files, err := Diff(dir, "HEAD~1", "HEAD")
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "foo.go", files[0].Path)
	assert.False(t, files[0].IsNew)
	assert.False(t, files[0].IsDeleted)
	assert.NotEmpty(t, files[0].ChangedLines)
}

func TestDiff_DeletedFile(t *testing.T) {
	dir := setupGitRepo(t)
	require.NoError(t, os.Remove(filepath.Join(dir, "foo.go")))
	shCmd := exec.Command("sh", "-c", "git add -A && git commit -m 'delete foo'")
	shCmd.Dir = dir
	shCmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test.com")
	require.NoError(t, shCmd.Run())

	files, err := Diff(dir, "HEAD~1", "HEAD")
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, "foo.go", files[0].Path)
	assert.True(t, files[0].IsDeleted)
}
