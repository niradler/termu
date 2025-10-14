package agent

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/niradler/termu/internal/config"
	"github.com/niradler/termu/internal/tools"
	"github.com/openai/openai-go/option"
)

type Agent struct {
	genkit *genkit.Genkit
	model  ai.Model
	tools  []ai.Tool
}

type Response struct {
	Text    string
	Command string
}

func New(ctx context.Context, cfg *config.Config) (*Agent, error) {
	var g *genkit.Genkit
	var model ai.Model

	switch cfg.Model.Provider {
	case "ollama":
		plugin := &ollama.Ollama{
			ServerAddress: cfg.Model.Server,
			Timeout:       cfg.Model.Timeout,
		}
		g = genkit.Init(ctx, genkit.WithPlugins(plugin))
		model = plugin.DefineModel(g,
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

	case "openai":
		opts := []option.RequestOption{option.WithAPIKey(cfg.Model.APIKey)}
		if cfg.Model.BaseURL != "" {
			opts = append(opts, option.WithBaseURL(cfg.Model.BaseURL))
		}

		plugin := &compat_oai.OpenAICompatible{
			Opts:     opts,
			Provider: "openai",
			APIKey:   cfg.Model.APIKey,
			BaseURL:  cfg.Model.BaseURL,
		}
		g = genkit.Init(ctx, genkit.WithPlugins(plugin))
		model = plugin.DefineModel("openai", cfg.Model.Name, ai.ModelOptions{
			Label: "OpenAI Compatible - " + cfg.Model.Name,
			Supports: &ai.ModelSupports{
				Multiturn:  true,
				SystemRole: true,
				Tools:      true,
				Media:      false,
			},
		})

	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Model.Provider)
	}

	// Initialize tools once
	fsTools := tools.DefineFilesystemTools(g, cfg.Workdir)
	shellTool := tools.DefineShellTool(g, cfg.Workdir)
	allTools := append(fsTools, shellTool)

	return &Agent{
		genkit: g,
		model:  model,
		tools:  allTools,
	}, nil
}

func (a *Agent) Generate(ctx context.Context, userInput string, history []ai.Message) (*Response, error) {
	messages := []*ai.Message{
		{
			Role:    ai.RoleSystem,
			Content: []*ai.Part{ai.NewTextPart(GetSystemPrompt())},
		},
	}

	for i := range history {
		messages = append(messages, &history[i])
	}

	messages = append(messages, &ai.Message{
		Role:    ai.RoleUser,
		Content: []*ai.Part{ai.NewTextPart(userInput)},
	})

	toolRefs := make([]ai.ToolRef, len(a.tools))
	for i, t := range a.tools {
		toolRefs[i] = t
	}

	resp, err := genkit.Generate(ctx, a.genkit,
		ai.WithModel(a.model),
		ai.WithMessages(messages...),
		ai.WithTools(toolRefs...),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate with tools: %w", err)
	}

	return &Response{
		Text:    resp.Text(),
		Command: "",
	}, nil
}
