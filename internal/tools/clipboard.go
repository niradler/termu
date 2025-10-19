package tools

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type WriteClipboardInput struct {
	Content string `json:"content" jsonschema:"description=Text content to copy to clipboard"`
}

func DefineClipboardTools(g *genkit.Genkit) []ai.Tool {
	readClipboardTool := genkit.DefineTool(g, "read_clipboard",
		"Reads the current text content from the system clipboard",
		func(ctx *ai.ToolContext, input struct{}) (string, error) {
			content, err := clipboard.ReadAll()
			if err != nil {
				return "", fmt.Errorf("failed to read clipboard: %w", err)
			}

			if content == "" {
				return "Clipboard is empty", nil
			}

			return content, nil
		},
	)

	writeClipboardTool := genkit.DefineTool(g, "write_clipboard",
		"Writes text content to the system clipboard",
		func(ctx *ai.ToolContext, input WriteClipboardInput) (string, error) {
			err := clipboard.WriteAll(input.Content)
			if err != nil {
				return "", fmt.Errorf("failed to write clipboard: %w", err)
			}

			return fmt.Sprintf("Successfully copied %d characters to clipboard", len(input.Content)), nil
		},
	)

	return []ai.Tool{readClipboardTool, writeClipboardTool}
}
