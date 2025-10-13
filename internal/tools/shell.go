package tools

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type ExecuteCommandInput struct {
	Command string `json:"command" jsonschema:"description=Shell command to execute in the working directory"`
}

func DefineShellTool(g *genkit.Genkit, workdir string) ai.Tool {
	return genkit.DefineTool(g, "execute_command",
		`Executes shell commands for exploration, searching, and information gathering.

Use this tool for:
- Searching code: rg "pattern" --type go, rg -i "search term"
- Finding files: fd "*.py", fd -e go -e mod
- File previews: bat file.go, bat -n file.py
- Directory listings: eza -l, eza --tree --level 2
- Git operations: git status, git log --oneline -10, git diff
- Disk usage: dua, dua aggregate
- JSON/YAML: jq '.' data.json, yq '.key' config.yaml

The command runs in the current working directory context. Output includes both stdout and stderr.`,
		func(ctx *ai.ToolContext, input ExecuteCommandInput) (string, error) {
			shell, shellArg := getShell()
			cmd := exec.CommandContext(context.Background(), shell, shellArg, input.Command)
			cmd.Dir = workdir

			output, err := cmd.CombinedOutput()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					return fmt.Sprintf("Command failed with exit code %d:\n%s", exitErr.ExitCode(), string(output)), nil
				}
				return "", fmt.Errorf("failed to execute command: %w", err)
			}

			return string(output), nil
		},
	)
}

func getShell() (string, string) {
	if runtime.GOOS == "windows" {
		return "powershell.exe", "-Command"
	}
	return "/bin/sh", "-c"
}
