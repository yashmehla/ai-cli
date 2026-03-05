package llm

import (
	"context"

	"google.golang.org/genai"
)

type Gemini struct {
	client *genai.Client
}

func New(apiKey string) (*Gemini, error) {

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})

	if err != nil {
		return nil, err
	}

	return &Gemini{
		client: client,
	}, nil
}

func (g *Gemini) Generate(prompt string) (string, error) {

	resp, err := g.client.Models.GenerateContent(
		context.Background(),
		"gemini-1.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
