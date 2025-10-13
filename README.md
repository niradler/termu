# termu - Your Terminal Sidekick

termu is a terminal chat agent that helps you accomplish tasks using shell commands. Have conversations with termu, and it will assist you through intelligent command execution - all while keeping you in control.

## Overview

termu is your terminal sidekick - a conversational AI agent that understands what you want to accomplish and uses shell commands to help you get there. Powered by local LLMs (via Ollama), termu favors modern cross-platform CLI tools to be effective, while being capable of running any shell command in your current directory context. Every command requires your approval, but once approved in a session, similar commands flow smoothly without repeated confirmations.

## Features

### 💬 Conversational AI Agent

- Natural language interactions - just tell termu what you need
- Multi-turn conversations with context awareness
- Session-based command approval (approve once per session)
- Runs in your current working directory context
- Powered by Ollama (default: Qwen3 model)

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

termu uses [Genkit's tool calling feature](https://genkit.dev/docs/tool-calling/?lang=go) for intelligent filesystem operations and command execution.

**Available Tools:**

| Tool              | Purpose                     | When to Use                                  |
| ----------------- | --------------------------- | -------------------------------------------- |
| `read_file`       | Read complete file contents | Understanding code, checking implementations |
| `write_file`      | Create or overwrite files   | Creating new files, major rewrites           |
| `search_replace`  | Exact string replacements   | Targeted edits, renaming, bug fixes          |
| `list_directory`  | List files and directories  | Exploring project structure                  |
| `execute_command` | Run shell commands          | Searching (fd/rg), git operations, previews  |

**Benefits of Structured Tool Calling:**

- **Surgical edits**: Make precise changes with `search_replace` without rewriting entire files
- **Syntax awareness**: Tools understand file structure and context
- **Multi-step workflows**: Read → analyze → edit → verify in one interaction
- **Reliable execution**: Structured interfaces eliminate command parsing errors
- **Hybrid approach**: Combines direct file operations with shell commands for exploration

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
  file: ~/.termu/logs/termu.log
```

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

### Sandbox Mode (Safe Testing)

Test what termu would do without actually executing commands:

```bash
termu --sandbox chat
```

### Custom Config

```bash
termu --config ./custom-config.yaml chat
```

## How It Works

1. **Start a Session**: Launch `termu chat` in any directory
2. **Have a Conversation**: Tell termu what you want to accomplish in natural language
3. **AI Understanding**: Local LLM (via Ollama) understands your request
4. **Smart Tool Selection**: termu chooses the best tools for the job (favors modern cross-platform tools)
5. **Command Preview**: See exactly what command will run
6. **Your Approval**: Approve or reject the command; approvals are remembered for the session
7. **Execution in Context**: Command runs in your current working directory
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

**File Management**

- "find duplicate files by name in ./downloads"
- "show me the largest files taking up space"
- "organize photos by date into folders"

**Code Analysis**

- "find all TODO comments in TypeScript files"
- "count lines of code by language"
- "show files that haven't been modified in 6 months"

**Data Processing**

- "convert all YAML configs to JSON"
- "find the most common error in log files"
- "merge these CSV files and remove duplicates"

**System Maintenance**

- "analyze which directories use the most disk space"
- "find and preview log files from today"
- "search for files containing sensitive information"

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

## Security Disclosure

Found a security issue? Please email security@yourdomain.com instead of creating a public issue.
