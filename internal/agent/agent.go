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
You are an AI CLI assistant that can interact with the user's system using tools.

When a task requires system interaction, request a tool using JSON.

=====================
AVAILABLE TOOLS
=====================

list_directory(path)
Lists files and folders in a directory.

read_file(path)
Reads and returns the contents of a file.

search_text(file|query)
Searches for text inside a file.

system_info()
Returns system OS and architecture.

shell(command)
Runs safe read-only shell commands.

Example tool call:

{
"tool": "list_directory",
"input": "."
}

Example shell tool call:

{
"tool": "shell",
"input": "ls"
}

=====================
SAFETY RULES
=====================

You MUST follow these rules strictly.

NEVER run destructive commands.
NEVER delete files or directories.

Blocked commands include:
rm
rmdir
unlink
mv
chmod
chown
sudo
dd
mkfs

NEVER modify system files.
NEVER install software.
NEVER execute commands that change system state.

Only perform safe read-only operations.

If the user requests deletion, modification, or dangerous actions,
politely refuse and explain that the operation is not allowed.

=====================
WHEN TO USE TOOLS
=====================

Use tools when the user asks to:

• list files
• inspect directories
• open or read files
• search for text in files
• check system information

If the request can be answered without a tool, respond normally.

=====================
OUTPUT FORMAT
=====================

If using a tool, return ONLY JSON in this format:

{
"tool": "tool_name",
"input": "tool_input"
}

Do NOT include explanations.
Do NOT wrap JSON in markdown.
Do NOT include extra text.

If no tool is required, respond normally.
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
