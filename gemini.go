package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func ResponseToText(resp *genai.GenerateContentResponse) genai.Text {
	var content genai.Text

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content += part.(genai.Text)
			}
		}
	}

	return content
}

func GetGeminiResponse(userInput string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	var request string = fmt.Sprintf(`Write a roadmap about %s. Strictly follow the following format:
Roadmap name: name
Goal #Number: goal
Deadline: dd.mm - dd.mm
Mini-goals: mini-goal`, userInput)
	resp, err := model.GenerateContent(ctx, genai.Text(request))
	if err != nil {
		log.Fatal(err)
	}

	content := ResponseToText(resp)
	var text string = string(content)
	return text
}
