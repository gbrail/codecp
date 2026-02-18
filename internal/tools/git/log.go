package git

import (
	"log/slog"
	"os/exec"
	"strconv"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const defaultMaxEntries = 100

type LogArgs struct {
	MaxEntries int `json:"maxEntries,omitempty"`
}

type LogResult struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
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
	slog.DebugContext(tc, "git log", "maxEntries", maxEntries)

	cmd := exec.CommandContext(tc, "git", "log", "-n", strconv.Itoa(maxEntries))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &LogResult{
			Error: err.Error(),
		}, nil
	}
	slog.DebugContext(tc, "git log", "result", string(output))
	return &LogResult{
		Result: string(output),
	}, nil
}
