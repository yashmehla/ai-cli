package llm

import (
	"context"

	"google.golang.org/genai"
)

type Gemini struct {
	client *genai.Client
}

func New(apiKey string) (*Gemini, error) {

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
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

	ctx := context.Background()

	resp, err := g.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}
