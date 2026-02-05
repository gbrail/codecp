package internal

import (
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

func BeforeModel(ac agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
	fmt.Printf("Calling %s...", req.Model)
	return nil, nil
}

func AfterModel(ac agent.CallbackContext, resp *model.LLMResponse, err error) (*model.LLMResponse, error) {
	fmt.Printf("model done.\n")
	return nil, nil
}

func BeforeTool(tc tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
	fmt.Printf("Using %s...", tool.Name())
	return nil, nil
}

func AfterTool(tc tool.Context, tool tool.Tool, result map[string]any, r map[string]any, err error) (map[string]any, error) {
	fmt.Printf("tool done.\n")
	return nil, nil
}
