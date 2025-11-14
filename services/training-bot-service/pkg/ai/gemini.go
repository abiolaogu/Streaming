package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIClient interface {
	Generate(prompt string) (string, error)
	GenerateResponse(messages []Message) (string, error)
}

type Message struct {
	Role    string
	Content string
}

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-pro")
	model.SetTemperature(0.7)
	model.SetTopP(0.9)
	model.SetMaxOutputTokens(2048)

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (gc *GeminiClient) Generate(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := gc.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("content generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Extract text from response
	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}

	return result, nil
}

func (gc *GeminiClient) GenerateResponse(messages []Message) (string, error) {
	ctx := context.Background()

	// Start chat session
	chat := gc.model.StartChat()
	chat.History = make([]*genai.Content, 0)

	// Add message history
	for i, msg := range messages {
		if i == 0 && msg.Role == "system" {
			// System message - configure model
			continue
		}

		role := "user"
		if msg.Role == "assistant" || msg.Role == "model" {
			role = "model"
		}

		chat.History = append(chat.History, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}

	// Get last user message
	lastMessage := messages[len(messages)-1]

	// Generate response
	resp, err := chat.SendMessage(ctx, genai.Text(lastMessage.Content))
	if err != nil {
		return "", fmt.Errorf("chat response failed: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Extract text
	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}

	return result, nil
}

func (gc *GeminiClient) Close() error {
	return gc.client.Close()
}
