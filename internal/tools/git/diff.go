package git

import (
	"os/exec"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type DiffArgs struct {
	Commit string `json:"commit,omitempty"` // If empty, diffs against the index (staging area).
	Path   string `json:"path,omitempty"`   // Optional file path to diff.
}

type DiffResult struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func DiffTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "GitDiff",
		Description: `
Return the git diff for changes in the working tree, optionally against a specific
commit or reference. If no commits are specified, it returns the diff between the working
tree and the index (staging area). filePath can be used to limit the diff to a specific file.`,
	}, gitDiff)
}

func gitDiff(tc tool.Context, inputArgs *DiffArgs) (*DiffResult, error) {
	var args []string
	args = append(args, "diff")
	if inputArgs.Commit != "" && inputArgs.Path != "" {
		args = append(args, inputArgs.Commit, inputArgs.Path)
	} else if inputArgs.Commit != "" {
		args = append(args, inputArgs.Commit)
	} else if inputArgs.Path != "" {
		args = append(args, inputArgs.Path)
	}

	cmd := exec.CommandContext(tc, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &DiffResult{
			Error: err.Error(),
		}, nil
	}
	return &DiffResult{
		Result: string(output),
	}, nil
}
