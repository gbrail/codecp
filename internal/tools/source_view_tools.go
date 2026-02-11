package tools

import (
	"github.com/gbrail/codecp/internal/tools/file"
	"github.com/gbrail/codecp/internal/tools/git"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
)

var sourceTools = []func() (tool.Tool, error){
	file.ListFilesTool,
	file.ReadFileTool,
	git.StatusTool,
	git.LogTool,
	git.DiffTool,
}

type SourceViewTools struct {
}

func (t *SourceViewTools) Name() string {
	return "SourceViewingTools"
}

func (t *SourceViewTools) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	var tools []tool.Tool
	for _, t := range sourceTools {
		newTool, err := t()
		if err != nil {
			return nil, err
		}
		tools = append(tools, newTool)
	}
	return tools, nil
}
