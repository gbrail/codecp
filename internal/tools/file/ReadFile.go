package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type ReadFileArgs struct {
	Path string `json:"path"`
}

type ReadFileResult struct {
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}

func ReadFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "ReadFile",
		Description: `
Read the contents of a file within the current directory.
`,
	}, readFile)
}

func readFile(tc tool.Context, args *ReadFileArgs) (*ReadFileResult, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	root, err := os.OpenRoot(wd)
	if err != nil {
		return nil, fmt.Errorf("failed to open root: %v", err)
	}
	defer root.Close()

	relPath := args.Path
	if relPath == "" {
		return &ReadFileResult{Error: "Path is required"}, nil
	}

	// Ensure the path is within the working directory
	if filepath.IsAbs(relPath) {
		rel, err := filepath.Rel(wd, relPath)
		if err != nil || strings.HasPrefix(rel, "..") {
			return &ReadFileResult{Error: "Access denied: path is outside current directory"}, nil
		}
		relPath = rel
	}
	relPath = filepath.Clean(relPath)
	if strings.HasPrefix(relPath, "..") || relPath == ".." {
		return &ReadFileResult{Error: "Access denied: path escaping detected"}, nil
	}

	// Use root.ReadFile to ensure we don't escape the root
	data, err := root.ReadFile(relPath)
	if err != nil {
		return &ReadFileResult{
			Error: fmt.Sprintf("Error reading file %q: %v", args.Path, err),
		}, nil
	}

	return &ReadFileResult{Content: string(data)}, nil
}
