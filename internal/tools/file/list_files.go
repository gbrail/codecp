package file

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const defaultDepth = 1

var separatorStr = string(filepath.Separator)

type ListFilesArgs struct {
	Directory string `json:"directory"`
	Depth     int    `json:"depth"`
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
a list of file names or an error. The "depth" parameter specifies how deep
to recurse. This parameter defaults to 1.
`,
	}, listFiles)
}

func listFiles(tc tool.Context, args *ListFilesArgs) (*ListFilesResult, error) {
	return doListFiles(tc, args)
}

func doListFiles(tc context.Context, args *ListFilesArgs) (*ListFilesResult, error) {
	slog.DebugContext(tc, "list files", "directory", args.Directory, "depth", args.Depth)

	depth := args.Depth
	if depth == 0 {
		depth = defaultDepth
	}

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
		nameLen := len(strings.Split(path, separatorStr))
		if nameLen > depth {
			return fs.SkipDir
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return &ListFilesResult{
			Error: fmt.Sprintf("Error reading %q: %v", args.Directory, err),
		}, nil
	}
	slog.DebugContext(tc, "list files", "files", files)
	return &ListFilesResult{Files: files}, nil
}
