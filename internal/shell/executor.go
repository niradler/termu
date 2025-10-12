package shell

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

type Executor struct {
	workdir string
	sandbox bool
}

type ExecutionResult struct {
	Command  string
	Output   string
	Error    string
	ExitCode int
	Executed bool
}

func New(workdir string, sandbox bool) *Executor {
	return &Executor{
		workdir: workdir,
		sandbox: sandbox,
	}
}

func (e *Executor) Execute(ctx context.Context, command string) (*ExecutionResult, error) {
	result := &ExecutionResult{
		Command:  command,
		Executed: !e.sandbox,
	}

	if e.sandbox {
		result.Output = "[SANDBOX MODE] Command would be executed: " + command
		return result, nil
	}

	shell, shellArg := getShell()
	cmd := exec.CommandContext(ctx, shell, shellArg, command)
	cmd.Dir = e.workdir

	output, err := cmd.CombinedOutput()
	result.Output = string(output)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Error = err.Error()
		} else {
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return result, nil
}

func getShell() (string, string) {
	if runtime.GOOS == "windows" {
		return "powershell.exe", "-Command"
	}
	return "/bin/sh", "-c"
}
