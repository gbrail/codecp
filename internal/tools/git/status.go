package git

import (
	"os/exec"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type StatusArgs struct{}

type StatusResult struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
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

func gitStatus(tc tool.Context, _ *StatusArgs) (*StatusResult, error) {
	cmd := exec.CommandContext(tc, "git", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &StatusResult{
			Error: err.Error(),
		}, nil
	}
	return &StatusResult{
		Result: string(output),
	}, nil
}
