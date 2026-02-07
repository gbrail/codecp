package git

import (
	"github.com/go-git/go-git/v6"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type StatusArgs struct{}

type StatusResult struct {
	Modified  []string `json:"modified,omitempty"`
	Untracked []string `json:"untracked,omitempty"`
	Error     string   `json:"error,omitempty"`
}

func StatusTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "GitStatus",
		Description: `
Return a list of files in the current working directory which have been modified
and are under the control of git, and a list of files that are not under the
control of git.`,
	}, gitStatus)
}

func gitStatus(tc tool.Context, args *StatusArgs) (*StatusResult, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return &StatusResult{
			Error: err.Error(),
		}, nil
	}

	tree, err := repo.Worktree()
	if err != nil {
		return &StatusResult{
			Error: err.Error(),
		}, nil
	}

	status, err := tree.Status()
	if err != nil {
		return &StatusResult{
			Error: err.Error(),
		}, nil
	}
	var modified []string
	var untracked []string
	for fileName, s := range status {
		if s.Worktree == git.Untracked {
			untracked = append(untracked, fileName)
		} else if s.Worktree != git.Unmodified || s.Staging != git.Unmodified {
			modified = append(modified, fileName)
		}
	}
	return &StatusResult{
		Modified:  modified,
		Untracked: untracked,
	}, nil
}
