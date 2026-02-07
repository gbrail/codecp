package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ListFilesArgs struct {
	Directory string `json:"directory"`
}

type ListFilesResult struct {
	Files []string `json:"result,omitempty"`
	Error string   `json:"error,omitempty"`
}

func ListFilesTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "ListFiles",
		Description: `
List the files recursively in the specified directory, returning either
a list of file names or an error.
`,
	}, listFiles)
}

func listFiles(tc tool.Context, args *ListFilesArgs) (*ListFilesResult, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	root, err := os.OpenRoot(wd)
	if err != nil {
		return nil, fmt.Errorf("failed to open root: %v", err)
	}
	defer root.Close()

	relPath := args.Directory
	if relPath == "" {
		relPath = "."
	}

	// Ensure the path is within the working directory
	if filepath.IsAbs(relPath) {
		rel, err := filepath.Rel(wd, relPath)
		if err != nil || strings.HasPrefix(rel, "..") {
			return &ListFilesResult{Error: "Access denied: path is outside current directory"}, nil
		}
		relPath = rel
	}
	relPath = filepath.Clean(relPath)
	if strings.HasPrefix(relPath, "..") {
		return &ListFilesResult{Error: "Access denied: path escaping detected"}, nil
	}

	// Open the requested directory as a sub-root to restrict access and maintain path structure
	subRoot, err := root.OpenRoot(relPath)
	if err != nil {
		return &ListFilesResult{
			Error: fmt.Sprintf("Error opening directory %q: %v", args.Directory, err),
		}, nil
	}
	defer subRoot.Close()

	var files []string
	err = fs.WalkDir(subRoot.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return &ListFilesResult{
			Error: fmt.Sprintf("Error reading %q: %v", args.Directory, err),
		}, nil
	}
	return &ListFilesResult{Files: files}, nil
}
