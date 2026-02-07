package git

import (
	"io"

	"github.com/go-git/go-git/v6"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const defaultMaxEntries = 100

type LogArgs struct {
	MaxEntries int `json:"maxEntries,omitempty"`
}

type LogResult struct {
	Commits []string `json:"commits,omitempty"`
	Error   string   `json:"error,omitempty"`
}

func LogTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "GitLog",
		Description: `
Return a log of git commits in the current branch. If not specified
it will return the last 100 entries.`,
	}, gitLog)
}

func gitLog(tc tool.Context, args *LogArgs) (*LogResult, error) {
	maxEntries := defaultMaxEntries
	if args.MaxEntries != 0 {
		maxEntries = args.MaxEntries
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		return &LogResult{
			Error: err.Error(),
		}, nil
	}

	var commits []string
	l, err := repo.Log(&git.LogOptions{})
	i := 0
	for {
		commit, err := l.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return &LogResult{
				Error: err.Error(),
			}, nil
		}
		if commit == nil || i == maxEntries {
			break
		}
		commits = append(commits, commit.String())
		i++
	}
	return &LogResult{
		Commits: commits,
	}, nil
}
