# termu - Your Terminal Sidekick

termu is a terminal chat agent that helps you accomplish tasks using shell commands. Have conversations with termu, and it will assist you through intelligent command execution - all while keeping you in control.

## Overview

termu is your terminal sidekick - a conversational AI agent that understands what you want to accomplish and helps you get there. Powered by local LLMs (via Ollama) or OpenAI-compatible servers (like LiteLLM), termu can directly read and edit files, execute shell commands for exploration, and leverage modern cross-platform CLI tools - all while keeping you in control with session-based command approval.

## Features

### 💬 Conversational AI Agent with Direct File Access

- Natural language interactions - just tell termu what you need
- Direct filesystem access through structured tool calling
- Can read, write, and surgically edit files without shell commands
- Multi-turn conversations with context awareness
- Session-based command approval (approve once per session)
- Runs in your current working directory context
- Powered by Ollama (default: Qwen3 model) or OpenAI-compatible servers (LiteLLM, etc.)

### 💬 Chat-First Experience

- Rich interactive chat mode powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Markdown-formatted responses with [Glow](https://github.com/charmbracelet/glow)
- Command preview before execution with syntax highlighting
- Session-based conversation history
- Split-pane interface: chat on left, execution on right
- Approve/reject commands with simple keyboard shortcuts

### 🔒 Safe Command Execution

- **Session-Based Approval**: Approve/reject commands before execution; approvals persist for the entire session
- **Current Directory Context**: All commands run in the folder where termu was launched
- **Command Whitelist**: Control which commands termu can use
- **Destructive Action Guard**: Prevents dangerous operations
- **Sandbox Mode**: Test behavior without actually executing commands

### ⚡ Smart Tool Selection

termu favors modern, cross-platform CLI tools to accomplish tasks efficiently. You don't need to know how to use them - termu handles that:

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

### 🛠️ Structured Tool Calling

termu uses [Genkit's tool calling feature](https://genkit.dev/docs/tool-calling/?lang=go) for intelligent filesystem operations and command execution. Instead of relying solely on shell commands, termu has direct access to structured tools for file manipulation and code editing.

**Available Tools:**

| Tool              | Purpose                           | When to Use                                                    |
| ----------------- | --------------------------------- | -------------------------------------------------------------- |
| `read_file`       | Read complete file contents       | Understanding code, checking implementations, reviewing files  |
| `write_file`      | Create or completely overwrite    | Creating new files, major rewrites (overwrites entire content) |
| `search_replace`  | Exact string search and replace   | Targeted edits, renaming, bug fixes, surgical code changes     |
| `list_directory`  | List files and directories        | Exploring project structure, finding files (with recursion)    |
| `execute_command` | Run shell commands in working dir | Searching (fd/rg), git operations, previews (bat/eza)          |
| `read_clipboard`  | Read system clipboard content     | Accessing copied text, working with clipboard data             |
| `write_clipboard` | Write content to system clipboard | Copying results, sharing data between applications             |

**Benefits of Structured Tool Calling:**

- **Direct file access**: Read and write files without shell command parsing
- **Surgical edits**: Make precise changes with `search_replace` without rewriting entire files
- **Reliability**: Structured JSON interfaces eliminate command parsing ambiguity
- **Multi-step workflows**: Read → analyze → edit → verify in one conversation
- **Hybrid approach**: Direct file operations + shell commands for exploration
- **Context awareness**: All operations run in your current working directory

## Installation

### Prerequisites

1. **Go 1.23+**
2. **Ollama** - [Download](https://ollama.com/download)
3. **Qwen3 Model** (or your preferred model):
   ```bash
   ollama pull qwen3
   ```

### Install termu

```bash
go install github.com/niradler/termu/cmd/termu@latest
```

Or build from source:

```bash
git clone https://github.com/niradler/termu.git
cd termu
go build -o termu cmd/termu/main.go
```

### Install Favorite Cross-Platform Tools

termu works with standard commands but is supercharged with these modern tools. The built-in installer can help you set them up:

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

Create `.termu.yaml` in your home directory or project root:

```yaml
# AI Model Configuration
model:
  provider: ollama # Options: "ollama" (default) or "openai" (for OpenAI-compatible servers)
  name: qwen3 # Model name
  server: http://127.0.0.1:11434 # Ollama server address (for ollama provider)
  timeout: 60 # Request timeout in seconds

  # For OpenAI-compatible servers (e.g., LiteLLM, custom OpenAI endpoints)
  # api_key: "sk-1234"           # Your API key (required for openai provider)
  # base_url: "http://localhost:4000/v1"  # Custom endpoint URL (required for openai provider)

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

# Logging
logging:
  level: info
  file: ~/.termu/logs/termu.log
```

### Using OpenAI-Compatible Servers (LiteLLM, etc.)

termu supports OpenAI-compatible servers like [LiteLLM](https://docs.litellm.ai/), allowing you to use various LLM providers through a unified API:

```yaml
model:
  provider: "openai" # Use OpenAI-compatible provider
  name: "gpt-4" # Model name (depends on your server config)
  api_key: "sk-1234" # API key (required)
  base_url: "http://localhost:4000/v1" # Custom endpoint URL
  timeout: 60
```

**Example: Running with LiteLLM**

1. Start LiteLLM server:

   ```bash
   litellm --model gpt-4 --api_base http://localhost:4000
   ```

2. Configure `.termu.yaml` with the settings above

3. Start termu:
   ```bash
   termu chat
   ```

See `.termu.openai.example.yaml` for a complete example configuration.

## Usage

### Start a Chat Session

```bash
termu chat
```

This starts an interactive session where you can have a conversation with termu. Tell it what you want to accomplish, and it will suggest commands to help.

**Chat Features:**

- 💬 Natural conversation - just describe what you want
- 🎨 Beautiful split-pane interface
- 📝 Markdown-formatted responses
- ✅ Approve/reject commands before execution
- 🔄 Session-based approval - approve once, reuse for similar commands
- 📂 Commands run in your current working directory
- 📜 Full conversation history within the session

**Keyboard Shortcuts:**

- `Enter` - Send message / Approve command
- `Esc` - Reject command
- `Ctrl+C` - Cancel operation
- `Ctrl+D` - Exit session
- `↑/↓` - Navigate history
- `Ctrl+L` - Clear screen

### Quick Command

```bash
termu "find all Python files modified in the last week"
```

### Custom Config

```bash
termu --config ./custom-config.yaml chat
```

## How It Works

1. **Start a Session**: Launch `termu chat` in any directory
2. **Have a Conversation**: Tell termu what you want to accomplish in natural language
3. **AI Understanding**: Local LLM (via Ollama) or OpenAI-compatible server understands your request
4. **Smart Tool Selection**: termu chooses the best approach - direct file operations or shell commands
5. **Preview & Approval**: See exactly what termu will do (read/edit files or run commands)
6. **Your Control**: Approve or reject; approvals are remembered for the session
7. **Execution in Context**: Operations run in your current working directory
8. **Continue the Conversation**: Discuss results, refine, or move to the next task

## Sessions & Command Approval

termu uses a session-based approval system:

- **First Time**: When termu wants to run a command, you approve or reject it
- **Within Session**: Once approved, similar commands in the same session don't require re-approval
- **Session End**: When you exit termu, approval history is cleared
- **New Session**: Fresh start with new approval requirements

This gives you control while keeping the conversation flowing naturally.

## Security Best Practices

1. **Review Before Approving**: Always check the command preview
2. **Session Awareness**: Remember that approvals persist throughout your session

## Example Conversations

Here are some examples of what you can ask termu:

### File Management & Code Editing

- "Read the README.md and update the installation section"
- "Find all TODO comments in TypeScript files and list them"
- "Replace all occurrences of 'oldFunction' with 'newFunction' in src/utils.go"
- "Create a new config.yaml file with database settings"

### Data Analysis & Processing

- "Find duplicate files by name in ./downloads"
- "Search for the most common error in log files using ripgrep"
- "Show me the largest files taking up disk space"
- "Convert all YAML configs to JSON"

### Code Exploration

- "List all Go files in the project recursively"
- "Show me what changed in the last 5 git commits"
- "Preview the main.go file with syntax highlighting"
- "Find files that haven't been modified in 6 months"

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

Add the command to your `.termu.yaml`:

```yaml
security:
  allowed_commands:
    - your_command
```

## Development

### Project Structure

```
termu/
├── cmd/termu/           # CLI entry point
├── internal/
│   ├── agent/           # Genkit agent implementation
│   ├── security/        # Security validators and approvals
│   ├── shell/           # Shell execution in cwd context
│   ├── tools/           # Tool management and installer
│   ├── tui/             # Interactive chat interface
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
go build -o termu cmd/termu/main.go
```

## Roadmap

- [ ] Enhanced conversation memory across sessions
- [ ] Support for additional models providers
- [ ] MCP Servers support
- [ ] Multi-session management

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
- Amazing Rust/Go CLI tools that make termu powerful
