package agent

import (
	"ai-cli/internal/llm"
)

type Agent struct {
	model *llm.Gemini
}

func NewAgent(apiKey string) *Agent {

	model, err := llm.New(apiKey)

	if err != nil {
		panic(err)
	}

	return &Agent{
		model: model,
	}
}

func (a *Agent) Handle(input string) string {

	resp, err := a.model.Generate(input)

	if err != nil {
		return err.Error()
	}

	return resp
}
