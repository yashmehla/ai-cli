package agent

import (
	"ai-cli/internal/llm"
	"ai-cli/internal/tools"
	"encoding/json"
	"regexp"
	"strings"
)

type Agent struct {
	model *llm.Gemini
	tools *tools.Registry
}

const systemPrompt = `
You are an AI CLI assistant.

You can request tools using JSON.

Available tools:

shell(command)
Runs safe shell commands.

Example tool call:

{
"tool": "shell",
"input": "ls"
}

STRICT SAFETY RULES:
- NEVER delete files or directories
- NEVER run commands like rm, rmdir, unlink, mv, sudo, chmod, chown
- NEVER modify system files
- ONLY run read-only commands

Allowed command types:
- list files
- inspect directories
- show system info

If a request involves deletion or modification, refuse the request.

Return JSON only when calling a tool.
Do NOT wrap JSON in markdown.
`
func NewAgent(apiKey string) *Agent {

	model, err := llm.New(apiKey)
	if err != nil {
		panic(err)
	}

	registry := tools.NewRegistry()
	registry.Register(tools.ShellTool{})

	return &Agent{
		model: model,
		tools: registry,
	}
}

func extractJSON(resp string) string {

	// remove markdown code blocks
	re := regexp.MustCompile("(?s)```.*?```")
	match := re.FindString(resp)

	if match != "" {
		match = strings.TrimPrefix(match, "```json")
		match = strings.TrimPrefix(match, "```")
		match = strings.TrimSuffix(match, "```")
		return strings.TrimSpace(match)
	}

	return resp
}

func (a *Agent) Handle(input string) string {

	input = strings.TrimSpace(input)

	prompt := systemPrompt + "\nUser: " + input

	resp, err := a.model.Generate(prompt)
	if err != nil {
		return err.Error()
	}

	jsonText := extractJSON(resp)

	if strings.Contains(jsonText, `"tool"`) {

		var toolCall struct {
			Tool  string `json:"tool"`
			Input string `json:"input"`
		}

		err := json.Unmarshal([]byte(jsonText), &toolCall)
		if err != nil {
			return resp
		}

		tool := a.tools.Get(toolCall.Tool)
		if tool == nil {
			return "Unknown tool requested"
		}

		result, err := tool.Run(toolCall.Input)
		if err != nil {
			return err.Error()
		}

		finalPrompt := "Tool result:\n" + result + "\nExplain this to the user."

		finalResp, err := a.model.Generate(finalPrompt)
		if err != nil {
			return err.Error()
		}

		return finalResp
	}

	return resp
}
