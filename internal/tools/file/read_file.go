package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const (
	defaultReadLength = 32 * 1024
	maxReadLength     = 128 * 1024 // 128KB max
)

type ReadFileArgs struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset,omitempty"`
	Length int64  `json:"length,omitempty"`
}

type ReadFileResult struct {
	Content string `json:"content,omitempty"`
	Total   int64  `json:"total,omitempty"`
	Error   string `json:"error,omitempty"`
}

func ReadFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "ReadFile",
		Description: `
Read the contents of a file within the current directory. Supports optional offset and length.
Defaults to reading 32KB if length is not specified. Maximum length is 128KB.
`,
	}, readFile)
}

func readFile(tc tool.Context, args *ReadFileArgs) (*ReadFileResult, error) {
	if args.Offset < 0 {
		return nil, fmt.Errorf("invalid offset")
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

	f, err := root.Open(relPath)
	if err != nil {
		return &ReadFileResult{
			Error: fmt.Sprintf("Error opening file %q: %v", args.Path, err),
		}, nil
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return &ReadFileResult{Error: fmt.Sprintf("Error stating file: %v", err)}, nil
	}
	totalSize := info.Size()

	length := args.Length
	if length <= 0 {
		length = defaultReadLength
	}
	if length > maxReadLength {
		length = maxReadLength
	}

	if args.Offset >= totalSize {
		return &ReadFileResult{Content: "", Total: totalSize}, nil
	}

	if args.Offset+length > totalSize {
		length = totalSize - args.Offset
	}

	buf := make([]byte, length)
	n, err := f.ReadAt(buf, args.Offset)
	if err != nil && err != io.EOF {
		return &ReadFileResult{
			Error: fmt.Sprintf("Error reading file at offset %d: %v", args.Offset, err),
		}, nil
	}

	return &ReadFileResult{
		Content: string(buf[:n]),
		Total:   totalSize,
	}, nil
}
