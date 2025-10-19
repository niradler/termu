package agent

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/niradler/termu/internal/config"
	"github.com/niradler/termu/internal/tools"
)

type Provider interface {
	GenerateResponse(ctx context.Context, messages []*ai.Message) (string, error)
}

type Agent struct {
	provider Provider
	genkit   *genkit.Genkit
	tools    []ai.Tool
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
	var g *genkit.Genkit
	var allTools []ai.Tool

	switch cfg.Model.Provider {
	case "ollama":
		ollamaProvider, err := newOllamaProvider(ctx, cfg)
		if err != nil {
			return nil, err
		}
		provider = ollamaProvider
		g = ollamaProvider.genkit

		fsTools := tools.DefineFilesystemTools(g, cfg.Workdir)
		shellTool := tools.DefineShellTool(g, cfg.Workdir)
		clipboardTools := tools.DefineClipboardTools(g)
		allTools = append(fsTools, shellTool)
		allTools = append(allTools, clipboardTools...)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Model.Provider)
	}

	return &Agent{
		provider: provider,
		genkit:   g,
		tools:    allTools,
	}, nil
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

func (a *Agent) Generate(ctx context.Context, userInput string, history []ai.Message) (*Response, error) {
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

	toolRefs := make([]ai.ToolRef, len(a.tools))
	for i, t := range a.tools {
		toolRefs[i] = t
	}

	resp, err := genkit.Generate(ctx, a.genkit,
		ai.WithModel(a.getModel()),
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

func (a *Agent) getModel() ai.Model {
	if ollamaProvider, ok := a.provider.(*OllamaProvider); ok {
		return ollamaProvider.model
	}
	return nil
}
