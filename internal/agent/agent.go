package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/yourusername/olloco/internal/config"
	"github.com/yourusername/olloco/internal/tools"
)

type Agent struct {
	genkit *genkit.Genkit
	model  ai.Model
}

type Response struct {
	Text    string
	Command string
}

func New(ctx context.Context, cfg *config.Config) (*Agent, error) {
	plugin := &ollama.Ollama{
		ServerAddress: cfg.Model.Server,
		Timeout:       cfg.Model.Timeout,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(plugin),
	)

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
				Tools:      false,
				Media:      false,
			},
		},
	)

	return &Agent{
		genkit: g,
		model:  model,
	}, nil
}

func (a *Agent) GenerateCommand(ctx context.Context, userInput string, history []ai.Message) (*Response, error) {
	systemPrompt := tools.GetSystemPrompt()

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

	resp, err := genkit.Generate(ctx, a.genkit,
		ai.WithModel(a.model),
		ai.WithMessages(messages...),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	text := resp.Text()
	command := extractCommand(text)

	return &Response{
		Text:    text,
		Command: command,
	}, nil
}

func extractCommand(text string) string {
	lines := strings.Split(text, "\n")
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

				if tools.IsKnownCommand(baseName) {
					return line
				}
			}
		}
	}
	return ""
}
