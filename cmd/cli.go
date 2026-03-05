package cmd

import (
	"ai-cli/config"
	"ai-cli/internal/agent"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartCLI() {

	cfg := config.Load()

	assistant := agent.NewAgent(cfg.GeminiAPIKey)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("AI CLI Assistant")
	fmt.Println("Type 'exit' to quit")

	for {
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		resp := assistant.Handle(input)

		fmt.Println(resp)
	}
}
