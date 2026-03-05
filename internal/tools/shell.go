package tools

import (
	"errors"
	"os/exec"
	"strings"
)

type ShellTool struct{}

func (s ShellTool) Name() string {
	return "shell"
}

func (s ShellTool) Description() string {
	return "run safe shell commands"
}

var allowed = []string{
	"ls",
	"pwd",
	"whoami",
	"date",
}

func allowedCommand(cmd string) bool {

	for _, c := range allowed {
		if cmd == c {
			return true
		}
	}

	return false
}

func (s ShellTool) Run(input string) (string, error) {

	parts := strings.Fields(input)

	if len(parts) == 0 {
		return "", errors.New("empty command")
	}

	if !allowedCommand(parts[0]) {
		return "", errors.New("command not allowed")
	}

	out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()

	return string(out), err
}
