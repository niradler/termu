package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

type ReadFileInput struct {
	Path string `json:"path" jsonschema:"description=Path to the file to read (relative to working directory)"`
}

type WriteFileInput struct {
	Path    string `json:"path" jsonschema:"description=Path to the file to write (relative to working directory)"`
	Content string `json:"content" jsonschema:"description=Complete content to write to the file"`
}

type SearchReplaceInput struct {
	Path       string `json:"path" jsonschema:"description=Path to the file to modify"`
	OldText    string `json:"old_text" jsonschema:"description=Exact text to search for and replace"`
	NewText    string `json:"new_text" jsonschema:"description=Text to replace with"`
	ReplaceAll bool   `json:"replace_all,omitempty" jsonschema:"description=Replace all occurrences (default: false, only first occurrence)"`
}

type ListDirectoryInput struct {
	Path      string `json:"path" jsonschema:"description=Directory path to list (relative to working directory)"`
	Recursive bool   `json:"recursive,omitempty" jsonschema:"description=List recursively (default: false)"`
}

func DefineFilesystemTools(g *genkit.Genkit, workdir string) []ai.Tool {
	readFileTool := genkit.DefineTool(g, "read_file",
		"Reads the complete contents of a file from the filesystem",
		func(ctx *ai.ToolContext, input ReadFileInput) (string, error) {
			fullPath := filepath.Join(workdir, input.Path)

			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("failed to read file %s: %w", input.Path, err)
			}

			return string(content), nil
		},
	)

	writeFileTool := genkit.DefineTool(g, "write_file",
		"Writes content to a file, creating or overwriting it",
		func(ctx *ai.ToolContext, input WriteFileInput) (string, error) {
			fullPath := filepath.Join(workdir, input.Path)

			dir := filepath.Dir(fullPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", fmt.Errorf("failed to create directory: %w", err)
			}

			if err := os.WriteFile(fullPath, []byte(input.Content), 0644); err != nil {
				return "", fmt.Errorf("failed to write file %s: %w", input.Path, err)
			}

			return fmt.Sprintf("Successfully wrote %d bytes to %s", len(input.Content), input.Path), nil
		},
	)

	searchReplaceTool := genkit.DefineTool(g, "search_replace",
		"Performs exact string search and replace in a file",
		func(ctx *ai.ToolContext, input SearchReplaceInput) (string, error) {
			fullPath := filepath.Join(workdir, input.Path)

			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("failed to read file %s: %w", input.Path, err)
			}

			originalContent := string(content)
			var newContent string
			var count int

			if input.ReplaceAll {
				count = strings.Count(originalContent, input.OldText)
				newContent = strings.ReplaceAll(originalContent, input.OldText, input.NewText)
			} else {
				if !strings.Contains(originalContent, input.OldText) {
					return "", fmt.Errorf("text not found: %s", input.OldText)
				}
				newContent = strings.Replace(originalContent, input.OldText, input.NewText, 1)
				count = 1
			}

			if err := os.WriteFile(fullPath, []byte(newContent), 0644); err != nil {
				return "", fmt.Errorf("failed to write file %s: %w", input.Path, err)
			}

			return fmt.Sprintf("Replaced %d occurrence(s) in %s", count, input.Path), nil
		},
	)

	listDirectoryTool := genkit.DefineTool(g, "list_directory",
		"Lists files and directories in a path",
		func(ctx *ai.ToolContext, input ListDirectoryInput) (string, error) {
			fullPath := filepath.Join(workdir, input.Path)

			if input.Recursive {
				var files []string
				err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					relPath, _ := filepath.Rel(workdir, path)
					if info.IsDir() {
						files = append(files, relPath+"/")
					} else {
						files = append(files, relPath)
					}
					return nil
				})
				if err != nil {
					return "", fmt.Errorf("failed to walk directory: %w", err)
				}
				return strings.Join(files, "\n"), nil
			}

			entries, err := os.ReadDir(fullPath)
			if err != nil {
				return "", fmt.Errorf("failed to read directory %s: %w", input.Path, err)
			}

			var result []string
			for _, entry := range entries {
				name := entry.Name()
				if entry.IsDir() {
					name += "/"
				}
				result = append(result, name)
			}

			return strings.Join(result, "\n"), nil
		},
	)

	return []ai.Tool{readFileTool, writeFileTool, searchReplaceTool, listDirectoryTool}
}
