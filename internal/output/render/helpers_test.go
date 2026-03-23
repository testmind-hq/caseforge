package render

import (
	"os"
	"testing"
)

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile(%q): %v", path, err)
	}
	return string(data)
}
