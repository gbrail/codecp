package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/gbrail/codecp/internal"
	"github.com/gbrail/codecp/internal/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	genai "google.golang.org/genai"
)

const (
	appName        = "codecp"
	gcpLocation    = "global"
	defaultModel   = "gemini-3-flash-preview"
	defaultSession = "default"
	defaultUser    = "default"
)

func main() {
	model := flag.String("f", defaultModel, "LLM model to use")
	debug := flag.Bool("d", false, "Print debug output")
	help := flag.Bool("h", false, "Print this message")
	flag.Parse()
	if !flag.Parsed() || *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var lvl slog.Level
	if *debug {
		lvl = slog.LevelDebug
	} else {
		lvl = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
	slog.Debug("debug is on")

	ctx := context.Background()

	config := &genai.ClientConfig{
		Project:  internal.GetGCPProject(),
		Location: gcpLocation,
		Backend:  genai.BackendVertexAI,
	}

	m, err := gemini.NewModel(ctx, *model, config)
	if err != nil {
		log.Fatalf("failed to create model: %v", err)
	}

	llmAgent, err := llmagent.New(llmagent.Config{
		Name:                 "codecp",
		Description:          "A helpful assistant",
		Model:                m,
		Instruction:          "You are a helpful assistant.",
		Toolsets:             []tool.Toolset{&tools.SourceViewTools{}},
		BeforeModelCallbacks: []llmagent.BeforeModelCallback{internal.BeforeModel},
		AfterModelCallbacks:  []llmagent.AfterModelCallback{internal.AfterModel},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{internal.BeforeTool},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{internal.AfterTool},
	})
	if err != nil {
		log.Fatalf("failed to create agent: %v", err)
	}

	sessions := session.InMemoryService()

	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          llmAgent,
		SessionService: sessions,
	})
	if err != nil {
		log.Fatalf("failed to create runner: %v", err)
	}

	_, err = sessions.Create(ctx, &session.CreateRequest{
		AppName:   appName,
		UserID:    defaultUser,
		SessionID: defaultSession,
	})
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("codecp> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "/exit" {
			break
		}

		userContent := &genai.Content{
			Parts: []*genai.Part{{Text: input}},
			Role:  "user",
		}

		events := r.Run(ctx, defaultUser, defaultSession, userContent, agent.RunConfig{})

		for event, err := range events {
			if err != nil {
				fmt.Printf("\nError: %v\n", err)
				break
			}
			if event.Content != nil {
				for _, part := range event.Content.Parts {
					if part.Text != "" {
						fmt.Print(part.Text)
					}
				}
			}
			if event.IsFinalResponse() {
				fmt.Println()
			}
		}
	}
}
