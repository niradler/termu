# Olloco - Secure Shell AI Agent

A production-ready CLI AI agent powered by Genkit (Go) and local Ollama models, designed to execute shell commands safely with built-in security features and optimized tooling.

## Overview

Olloco is an intelligent CLI agent that combines the power of local LLMs (via Ollama) with secure shell execution capabilities. It provides a safe environment for AI-assisted command-line operations through whitelisting, approval systems, and restricted folder access.

## Features

### 🤖 AI-Powered Command Execution

- Integration with Ollama (default: Qwen3 model)
- Context-aware command suggestions
- Multi-turn conversations with tool calling
- Smart command composition using optimized tools

### 🎨 Beautiful Terminal UI

- Rich interactive mode powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Markdown rendering with [Glow](https://github.com/charmbracelet/glow)
- Syntax-highlighted command preview before execution
- Real-time streaming output with proper formatting
- Split-pane view: AI reasoning + command execution
- Conversation history with scrollback
- Keyboard shortcuts for efficient workflow

### 🔒 Security First

- **Command Whitelist**: Only approved commands can execute
- **Folder Restrictions**: Limit access to specific directories
- **Approval System**: Prompt for confirmation on risky operations
- **Destructive Action Detection**: Prevents dangerous commands (rm -rf, format, etc.)
- **Sandbox Mode**: Test agent behavior without executing commands

### ⚡ Optimized Tooling

Olloco auto-detects and encourages the use of modern, fast command-line tools:

| Tool                                             | Purpose           | Why It's Better                          |
| ------------------------------------------------ | ----------------- | ---------------------------------------- |
| [sd](https://github.com/chmln/sd)                | Search & replace  | Simpler than sed, more intuitive         |
| [fd](https://github.com/sharkdp/fd)              | File finding      | 3x faster than find, user-friendly       |
| [ripgrep](https://github.com/BurntSushi/ripgrep) | Content search    | 5x faster than grep, respects .gitignore |
| [bat](https://github.com/sharkdp/bat)            | File preview      | Syntax highlighting, Git integration     |
| [xsv](https://github.com/BurntSushi/xsv)         | CSV toolkit       | Fast CSV processing and analysis         |
| [jaq](https://github.com/01mf02/jaq)             | JSON processing   | Faster jq alternative                    |
| [yq](https://github.com/mikefarah/yq)            | YAML/JSON         | Universal config file manipulation       |
| [dua](https://github.com/Byron/dua-cli)          | Disk analysis     | Interactive disk usage visualization     |
| [eza](https://eza.rocks/)                        | Directory listing | Modern ls with colors and Git status     |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         User Input                          │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                      Cobra CLI                              │
│                    (Command Parser)                         │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                    Genkit Agent                             │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Ollama (Qwen3 Model)                              │    │
│  │  - Context understanding                           │    │
│  │  - Tool selection                                  │    │
│  │  - Command generation                              │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  Security Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐  │
│  │  Whitelist   │  │   Folder     │  │   Destructive   │  │
│  │  Validator   │  │  Restriction │  │   Action Guard  │  │
│  └──────────────┘  └──────────────┘  └─────────────────┘  │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  Approval System                            │
│           (Interactive confirmation for                     │
│            high-risk operations)                            │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  Shell Executor                             │
│        (Cross-platform command execution)                   │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                      Output/Result                          │
└─────────────────────────────────────────────────────────────┘
```

## Installation

### Prerequisites

1. **Go 1.23+**
2. **Ollama** - [Download](https://ollama.com/download)
3. **Qwen3 Model** (or your preferred model):
   ```bash
   ollama pull qwen3
   ```

### Install Olloco

```bash
go install github.com/yourusername/olloco/cmd/olloco@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/olloco.git
cd olloco
go build -o olloco cmd/olloco/main.go
```

### Install Recommended Tools

Olloco works with standard commands but is supercharged with these tools:

**macOS (Homebrew):**

```bash
brew install sd fd ripgrep bat xsv yq dua-cli eza
brew install jaq
```

**Linux (Cargo):**

```bash
cargo install sd fd-find ripgrep bat xsv jaq dua-cli eza
```

**Windows (Scoop):**

```bash
scoop install sd fd ripgrep bat yq dua eza
cargo install xsv jaq
```

## Configuration

Create `.olloco.yaml` in your home directory or project root:

```yaml
# AI Model Configuration
model:
  provider: ollama
  name: qwen3
  server: http://127.0.0.1:11434
  timeout: 60

# Security Configuration
security:
  # Command whitelist (empty = allow all from safe_commands list)
  allowed_commands:
    - sd
    - fd
    - rg
    - bat
    - xsv
    - jaq
    - yq
    - dua
    - eza
    - ls
    - cat
    - grep
    - find
    - echo
    - pwd
    - cd
    - git
    - curl
    - wget

  # Restricted folders (agent cannot access these)
  restricted_folders:
    - /System
    - /Windows
    - /etc/passwd
    - /etc/shadow
    - ~/.ssh
    - ~/.aws

  # Allowed folders (when set, restricts to these only)
  allowed_folders:
    - ./
    - ~/Documents
    - ~/Projects

  # Commands that always require approval
  high_risk_commands:
    - rm
    - mv
    - chmod
    - chown
    - sudo
    - curl -X POST
    - wget -O

  # Destructive patterns (always blocked)
  blocked_patterns:
    - "rm -rf /"
    - "rm -rf *"
    - ":(){ :|:& };:" # fork bomb
    - "mkfs"
    - "dd if="
    - "> /dev/sda"

  # Enable sandbox mode (dry-run, no actual execution)
  sandbox_mode: false

  # Require approval for all commands
  always_approve: false

# Tool Configuration
tools:
  # Auto-install missing tools
  auto_install: false

  # Prefer optimized tools over traditional ones
  prefer_modern: true

# Logging
logging:
  level: info
  file: ~/.olloco/logs/olloco.log
```

## Usage

### Basic Command

```bash
olloco "find all Python files modified in the last week"
```

### Interactive Mode

```bash
olloco chat
```

**Interactive Mode Features:**

- 🎨 Beautiful split-pane interface
- 📝 Markdown-formatted AI responses
- 👁️ Real-time command preview
- ✅ Approve/reject commands with keyboard shortcuts
- 📜 Scrollable conversation history
- 🎯 Syntax highlighting for commands and output

**Keyboard Shortcuts:**

- `Enter` - Send message / Approve command
- `Ctrl+C` - Cancel current operation
- `Ctrl+D` - Exit interactive mode
- `↑/↓` - Navigate history
- `Tab` - Autocomplete
- `Esc` - Reject command
- `Ctrl+L` - Clear screen

### Sandbox Mode (Safe Testing)

```bash
olloco --sandbox "remove all .log files older than 30 days"
# Shows what would be executed without running it
```

### ⚠️ YOLO Mode (DANGEROUS!)

```bash
olloco --yolo chat
# ⚠️  Disables ALL security checks - use at your own risk!
# See YOLO_MODE.md for full details
```

### Custom Config

```bash
olloco --config ./custom-config.yaml "analyze disk usage"
```

## How It Works

1. **User Input**: You describe what you want to accomplish
2. **AI Processing**: Qwen3 (via Ollama) understands the request
3. **Tool Selection**: Agent chooses optimal tools (prefers modern tools)
4. **Command Generation**: Creates the appropriate command
5. **Security Check**: Validates against whitelist and restrictions
6. **Approval** (if needed): Prompts for confirmation
7. **Execution**: Runs the command safely
8. **Result**: Returns output and can iterate if needed

## Tool Documentation for AI

The agent is pre-trained with comprehensive documentation for all supported tools:

### Search & Replace (sd)

```bash
# Replace text in files
sd 'old_text' 'new_text' file.txt

# Replace in multiple files
fd -e js | xargs sd 'var ' 'const '

# Preview changes
sd -p 'old' 'new' file.txt
```

### File Finding (fd)

```bash
# Find by name
fd pattern

# Find by extension
fd -e py

# Find and execute
fd -e log -x rm {}

# Find modified in last 7 days
fd --changed-within 7d
```

### Content Search (ripgrep)

```bash
# Basic search
rg "pattern"

# Search specific file types
rg -t py "import"

# Show context
rg -C 3 "error"

# Case insensitive
rg -i "pattern"
```

### File Preview (bat)

```bash
# Show file with syntax highlighting
bat file.py

# Show with line numbers
bat -n file.py

# Show specific range
bat -r 10:20 file.py
```

### CSV Toolkit (xsv)

```bash
# Show CSV statistics
xsv stats data.csv

# Select columns
xsv select name,age data.csv

# Search CSV
xsv search -s name "John" data.csv

# Sort CSV
xsv sort -s age data.csv
```

### JSON Processing (jaq)

```bash
# Parse JSON
jaq '.' data.json

# Extract field
jaq '.items[] | .name' data.json

# Filter
jaq '.items[] | select(.price > 100)' data.json
```

### YAML/JSON (yq)

```bash
# Read YAML
yq '.key' config.yaml

# Convert YAML to JSON
yq -o json config.yaml

# Update value
yq '.key = "new_value"' -i config.yaml
```

### Disk Usage (dua)

```bash
# Interactive mode
dua i

# Show top directories
dua aggregate

# Show directory size
dua /path/to/dir
```

### Directory Listing (eza)

```bash
# List with details
eza -l

# Show tree
eza --tree

# Show with Git status
eza --git -l

# Sort by size
eza -l --sort size
```

## Security Best Practices

1. **Start Restrictive**: Begin with a narrow whitelist and expand as needed
2. **Use Allowed Folders**: Restrict operations to specific directories
3. **Enable Approvals**: Set `always_approve: true` for sensitive environments
4. **Test in Sandbox**: Use `--sandbox` mode to validate agent behavior
5. **Review Logs**: Regularly check logs for unexpected behavior
6. **Keep Tools Updated**: Update CLI tools for security patches

## Examples

### Safe File Management

```bash
olloco "find duplicate files by name in ./downloads"
olloco "show me the largest files taking up space"
olloco "organize photos by date into folders"
```

### Code Analysis

```bash
olloco "find all TODO comments in TypeScript files"
olloco "count lines of code by language"
olloco "show files that haven't been modified in 6 months"
```

### Data Processing

```bash
olloco "convert all YAML configs to JSON"
olloco "find the most common error in log files"
olloco "merge these three CSV files and remove duplicates"
```

### System Maintenance

```bash
olloco "analyze which directories use the most disk space"
olloco "find and preview log files from today"
olloco "search for files containing API keys or passwords"
```

## Troubleshooting

### Ollama Connection Issues

```bash
# Check if Ollama is running
ollama list

# Verify model is installed
ollama pull qwen3

# Test connection
curl http://127.0.0.1:11434/api/tags
```

### Command Not Whitelisted

Add the command to your `.olloco.yaml`:

```yaml
security:
  allowed_commands:
    - your_command
```

### Tool Not Found

Install missing tools or disable `prefer_modern`:

```yaml
tools:
  prefer_modern: false
```

## Development

### Project Structure

```
olloco/
├── cmd/olloco/          # CLI entry point
├── internal/
│   ├── agent/           # Genkit agent implementation
│   ├── security/        # Security validators and approvals
│   ├── shell/           # Shell execution
│   ├── tools/           # Tool management and docs
│   └── config/          # Configuration handling
├── go.mod
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o olloco cmd/olloco/main.go
```

## Roadmap

- [ ] Support for additional models (Llama 3, Mistral, etc.)
- [ ] Plugin system for custom tools
- [ ] Web UI for monitoring agent activity
- [ ] Multi-agent collaboration
- [ ] Cloud-hosted Ollama support
- [ ] Windows-specific command validation
- [ ] Tool installation automation
- [ ] Learning mode (agent improves from corrections)

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Credits

Built with:

- [Genkit](https://genkit.dev/) - AI workflow framework
- [Ollama](https://ollama.com/) - Local LLM runtime
- [Cobra](https://cobra.dev/) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Glow](https://github.com/charmbracelet/glow) - Markdown rendering
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- Amazing Rust/Go CLI tools that make this agent powerful

## Security Disclosure

Found a security issue? Please email security@yourdomain.com instead of creating a public issue.
