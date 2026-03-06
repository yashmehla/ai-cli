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

var allowedCommands = []string{
	"ls",
	"pwd",
	"whoami",
	"date",
	"cat",
}

var blockedCommands = []string{
	"rm",
	"rmdir",
	"unlink",
	"mv",
	"chmod",
	"chown",
	"sudo",
	"dd",
	"mkfs",
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func blocked(cmd string) bool {
	for _, b := range blockedCommands {
		if strings.Contains(cmd, b) {
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

	if blocked(input) {
		return "", errors.New("command blocked for safety")
	}

	if !contains(allowedCommands, parts[0]) {
		return "", errors.New("command not allowed")
	}

	out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()

	return string(out), err
}
