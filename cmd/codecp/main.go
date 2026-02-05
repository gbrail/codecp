package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gbrail/codecp/internal"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
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
	ctx := context.Background()

	// TODO check for GOOGLE_CLOUD_PROJECT
	config := &genai.ClientConfig{
		Project:  internal.GetGCPProject(),
		Location: gcpLocation,
		Backend:  genai.BackendVertexAI,
	}

	m, err := gemini.NewModel(ctx, defaultModel, config)
	if err != nil {
		log.Fatalf("failed to create model: %v", err)
	}

	llmAgent, err := llmagent.New(llmagent.Config{
		Name:        "codecp",
		Description: "A helpful assistant",
		Model:       m,
		Instruction: "You are a helpful assistant.",
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

	sessions.Create(ctx, &session.CreateRequest{
		AppName:   appName,
		UserID:    defaultUser,
		SessionID: defaultSession,
	})

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("codecp> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "/exit" {
			break
		}

		userContent := &genai.Content{
			Parts: []*genai.Part{
				&genai.Part{Text: input},
			},
			Role: "user",
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
