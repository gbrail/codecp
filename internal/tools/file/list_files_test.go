package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestListFiles(t *testing.T) {
	ctx := context.Background()
	td := filepath.Join(findRoot(), "testdata")

	t.Run("default", func(t *testing.T) {
		t.Logf("Listing in %s", td)
		dflt, err := doListFiles(ctx, &ListFilesArgs{
			Directory: td,
		})
		if err != nil {
			t.Fatalf("Error listing files: %v", err)
		}
		t.Logf("Result: %v", dflt.Files)
	})
}

func findRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "" // Reached filesystem root
		}
		dir = parent
	}
}
