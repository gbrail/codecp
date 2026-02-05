package tools

import (
	"github.com/gbrail/codecp/internal/tools/file"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
)

type SourceViewTools struct {
}

func (t *SourceViewTools) Name() string {
	return "SourceViewingTools"
}

func (t *SourceViewTools) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	var tools []tool.Tool
	newTool, err := file.ListFilesTool()
	if err != nil {
		return nil, err
	}
	tools = append(tools, newTool)
	newTool, err = file.ReadFileTool()
	if err != nil {
		return nil, err
	}
	tools = append(tools, newTool)
	return tools, nil
}
