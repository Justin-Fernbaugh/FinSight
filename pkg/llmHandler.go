package llm

import (
	"context"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const (
	systemPrompt = `You are a expert at crafting financial summarizations from transactions. 
	Given a list of transactions create a neatly organized summary of the transactions by category with emojis and a total amount for each category.
	Include the type of purchase and the vendor name in the summary which is from the last 7 days.`
)

var (
	temperature float32 = 0.9
	tokenLimit int32 = 1000
)

func Handler(projectID string, location string, modelName string) (*genai.GenerativeModel, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, err
	}
	gemini := client.GenerativeModel(modelName)
	gemini.Temperature = &temperature
	gemini.MaxOutputTokens = &tokenLimit
	return gemini, nil
}

func GenerateResponse(ctx context.Context, gemini *genai.GenerativeModel, prompt string) (string, error) {
	resp, err := gemini.GenerateContent(ctx, genai.Text(systemPrompt+prompt))
	if err != nil {
		return "", err
	}
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned in Gemini response")
	}
	content := resp.Candidates[0].Content
	if content == nil {
		return "", fmt.Errorf("content is nil in Gemini response")
	}
	parts := content.Parts
	if len(parts) == 0 {
		return "", fmt.Errorf("no parts found in Gemini response content")
	}

	var responseText string
	for _, part := range parts {
		if text, ok := part.(genai.Text); ok {
			responseText += string(text)
		}
	}

	return responseText, nil
}

