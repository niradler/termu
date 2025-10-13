package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/niradler/termu/internal/config"
)

type Provider interface {
	GenerateResponse(ctx context.Context, messages []*ai.Message) (string, error)
}

type Agent struct {
	provider Provider
}

type OllamaProvider struct {
	genkit *genkit.Genkit
	model  ai.Model
}

type Response struct {
	Text    string
	Command string
}

func New(ctx context.Context, cfg *config.Config) (*Agent, error) {
	var provider Provider
	var err error

	switch cfg.Model.Provider {
	case "ollama":
		provider, err = newOllamaProvider(ctx, cfg)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Model.Provider)
	}

	if err != nil {
		return nil, err
	}

	return &Agent{provider: provider}, nil
}

func newOllamaProvider(ctx context.Context, cfg *config.Config) (*OllamaProvider, error) {
	plugin := &ollama.Ollama{
		ServerAddress: cfg.Model.Server,
		Timeout:       cfg.Model.Timeout,
	}

	g := genkit.Init(ctx, genkit.WithPlugins(plugin))

	model := plugin.DefineModel(g,
		ollama.ModelDefinition{
			Name: cfg.Model.Name,
			Type: "chat",
		},
		&ai.ModelOptions{
			Label: "Ollama - " + cfg.Model.Name,
			Supports: &ai.ModelSupports{
				Multiturn:  true,
				SystemRole: true,
				Tools:      true,
				Media:      false,
			},
		},
	)

	return &OllamaProvider{
		genkit: g,
		model:  model,
	}, nil
}

func (p *OllamaProvider) GenerateResponse(ctx context.Context, messages []*ai.Message) (string, error) {
	resp, err := genkit.Generate(ctx, p.genkit,
		ai.WithModel(p.model),
		ai.WithMessages(messages...),
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return resp.Text(), nil
}

func (a *Agent) GenerateCommand(ctx context.Context, userInput string, history []ai.Message) (*Response, error) {
	systemPrompt := GetSystemPrompt()

	messages := []*ai.Message{
		{
			Role:    ai.RoleSystem,
			Content: []*ai.Part{ai.NewTextPart(systemPrompt)},
		},
	}

	for i := range history {
		messages = append(messages, &history[i])
	}

	messages = append(messages, &ai.Message{
		Role:    ai.RoleUser,
		Content: []*ai.Part{ai.NewTextPart(userInput)},
	})

	text, err := a.provider.GenerateResponse(ctx, messages)
	if err != nil {
		return nil, err
	}

	command := extractCommand(text)

	return &Response{
		Text:    text,
		Command: command,
	}, nil
}

func extractCommand(text string) string {
	lines := strings.Split(text, "\n")

	knownCommands := map[string]bool{
		"sd": true, "fd": true, "rg": true, "bat": true, "xsv": true,
		"jaq": true, "yq": true, "dua": true, "eza": true,
		"ls": true, "cat": true, "grep": true, "find": true,
		"echo": true, "pwd": true, "cd": true, "git": true,
		"curl": true, "wget": true, "rm": true, "mv": true,
		"cp": true, "mkdir": true, "touch": true, "chmod": true,
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "```") {
			continue
		}
		if line != "" && !strings.HasPrefix(line, "#") &&
			!strings.Contains(line, "I'll") && !strings.Contains(line, "would") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				cmd := strings.TrimPrefix(parts[0], "./")
				cmd = strings.TrimPrefix(cmd, ".\\")

				baseName := cmd
				if idx := strings.LastIndexAny(cmd, "/\\"); idx != -1 {
					baseName = cmd[idx+1:]
				}

				if knownCommands[baseName] {
					return line
				}
			}
		}
	}
	return ""
}
