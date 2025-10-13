package tui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/firebase/genkit/go/ai"
	"github.com/niradler/termu/internal/agent"
	"github.com/niradler/termu/internal/config"
	"github.com/niradler/termu/internal/security"
	"github.com/niradler/termu/internal/shell"
)

type SessionState int

const (
	StateInput SessionState = iota
	StateThinking
	StateApproval
	StateExecuting
)

type Message struct {
	Role    string
	Content string
}

type Model struct {
	ctx         context.Context
	state       SessionState
	textarea    textarea.Model
	viewport    viewport.Model
	messages    []Message
	aiHistory   []ai.Message
	currentCmd  string
	width       int
	height      int
	sandboxMode bool
	mdRenderer  *glamour.TermRenderer
	agent       *agent.Agent
	validator   *security.Validator
	executor    *shell.Executor
	workdir     string
}

type ExecutionCompleteMsg struct {
	Output       string
	Error        error
	NewMessages  []Message
	NewAIHistory []ai.Message
}

type ApprovalRequestMsg struct {
	Command string
}

func NewModel(ctx context.Context, cfg *config.Config, sandboxMode bool) (Model, error) {
	ta := textarea.New()
	ta.Placeholder = "Describe what you want to do..."
	ta.Focus()
	ta.CharLimit = 500
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingRight(2)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(78),
	)

	ag, err := agent.New(ctx, cfg)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create agent: %w", err)
	}

	workdir, _ := os.Getwd()

	return Model{
		ctx:         ctx,
		state:       StateInput,
		textarea:    ta,
		viewport:    vp,
		messages:    []Message{},
		aiHistory:   []ai.Message{},
		width:       80,
		height:      24,
		sandboxMode: sandboxMode,
		mdRenderer:  renderer,
		agent:       ag,
		validator:   security.New(cfg),
		executor:    shell.New(workdir, sandboxMode),
		workdir:     workdir,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyCtrlD:
			return m, tea.Quit

		case tea.KeyEnter:
			if m.state == StateInput && m.textarea.Value() != "" {
				userInput := m.textarea.Value()
				m.messages = append(m.messages, Message{
					Role:    "user",
					Content: userInput,
				})
				m.textarea.Reset()
				m.state = StateThinking
				m.updateViewport()
				return m, m.processUserInput(userInput)
			} else if m.state == StateApproval {
				m.validator.ApproveCommand(m.currentCmd)
				m.state = StateExecuting
				m.updateViewport()
				return m, m.executeCommand()
			}

		case tea.KeyEsc:
			if m.state == StateApproval {
				m.messages = append(m.messages, Message{
					Role:    "system",
					Content: "❌ Command rejected by user",
				})
				m.state = StateInput
				m.currentCmd = ""
				m.updateViewport()
			}

		case tea.KeyCtrlL:
			m.messages = []Message{}
			m.updateViewport()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10
		m.textarea.SetWidth(msg.Width - 4)
		m.updateViewport()

	case ApprovalRequestMsg:
		m.currentCmd = msg.Command
		m.state = StateApproval
		m.updateViewport()

	case ExecutionCompleteMsg:
		m.messages = append(m.messages, msg.NewMessages...)
		m.aiHistory = append(m.aiHistory, msg.NewAIHistory...)

		if msg.Error != nil {
			m.messages = append(m.messages, Message{
				Role:    "error",
				Content: fmt.Sprintf("Execution failed: %v", msg.Error),
			})
		} else if msg.Output != "" {
			m.messages = append(m.messages, Message{
				Role:    "output",
				Content: msg.Output,
			})
		}
		m.state = StateInput
		m.currentCmd = ""
		m.updateViewport()
	}

	if m.state == StateInput {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	header := m.renderHeader()
	b.WriteString(header)
	b.WriteString("\n\n")

	b.WriteString(m.viewport.View())
	b.WriteString("\n\n")

	if m.state == StateInput {
		b.WriteString(PromptStyle.Render("→ You:"))
		b.WriteString("\n")
		b.WriteString(m.textarea.View())
	} else if m.state == StateThinking {
		b.WriteString(InfoStyle.Render("🤔 AI is thinking..."))
	} else if m.state == StateApproval {
		b.WriteString(m.renderApproval())
	} else if m.state == StateExecuting {
		b.WriteString(InfoStyle.Render("⚡ Executing command..."))
	}

	b.WriteString("\n\n")
	b.WriteString(m.renderFooter())

	return BaseStyle.Render(b.String())
}

func (m Model) renderHeader() string {
	title := TitleStyle.Render("🤖 termu - Your Terminal Sidekick")

	var mode string
	if m.sandboxMode {
		mode = SandboxStyle.Render(" SANDBOX ")
	} else {
		mode = StatusBarStyle.Render(" SESSION ")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, title, " ", mode)
}

func (m Model) renderFooter() string {
	help := HelpStyle.Render(
		"Enter: Send/Approve • Esc: Reject • Ctrl+L: Clear • Ctrl+D: Exit",
	)
	return help
}

func (m Model) renderApproval() string {
	var b strings.Builder

	b.WriteString(ApprovalStyle.Render("⚠️  Command Approval Required"))
	b.WriteString("\n\n")
	b.WriteString(CommandStyle.Render(m.currentCmd))
	b.WriteString("\n\n")
	b.WriteString(SuccessStyle.Render("Press Enter to approve"))
	b.WriteString(" • ")
	b.WriteString(ErrorStyle.Render("Press Esc to reject"))

	return b.String()
}

func (m *Model) updateViewport() {
	var b strings.Builder

	for _, msg := range m.messages {
		switch msg.Role {
		case "user":
			b.WriteString(PromptStyle.Render("You: "))
			b.WriteString(msg.Content)
			b.WriteString("\n\n")

		case "assistant":
			b.WriteString(InfoStyle.Render("🤖 termu: "))
			b.WriteString("\n")
			if rendered, err := m.mdRenderer.Render(msg.Content); err == nil {
				b.WriteString(rendered)
			} else {
				b.WriteString(msg.Content)
			}
			b.WriteString("\n")

		case "command":
			b.WriteString(WarningStyle.Render("📝 Generated Command:"))
			b.WriteString("\n")
			b.WriteString(CommandStyle.Render(msg.Content))
			b.WriteString("\n")

		case "output":
			b.WriteString(SuccessStyle.Render("✅ Output:"))
			b.WriteString("\n")
			b.WriteString(OutputStyle.Render(msg.Content))
			b.WriteString("\n")

		case "error":
			b.WriteString(ErrorStyle.Render("❌ Error:"))
			b.WriteString("\n")
			b.WriteString(OutputStyle.Render(msg.Content))
			b.WriteString("\n")

		case "system":
			b.WriteString(HelpStyle.Render("ℹ️  " + msg.Content))
			b.WriteString("\n")
		}
	}

	m.viewport.SetContent(b.String())
	m.viewport.GotoBottom()
}

func (m Model) processUserInput(input string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.agent.GenerateCommand(m.ctx, input, m.aiHistory)
		if err != nil {
			return ExecutionCompleteMsg{
				Error: fmt.Errorf("failed to generate command: %w", err),
			}
		}

		newMessages := []Message{
			{Role: "assistant", Content: resp.Text},
		}

		newAIHistory := []ai.Message{
			{Role: ai.RoleUser, Content: []*ai.Part{ai.NewTextPart(input)}},
			{Role: ai.RoleModel, Content: []*ai.Part{ai.NewTextPart(resp.Text)}},
		}

		if resp.Command == "" {
			return ExecutionCompleteMsg{
				Output:       "No command was generated. Try rephrasing your request.",
				NewMessages:  newMessages,
				NewAIHistory: newAIHistory,
			}
		}

		newMessages = append(newMessages, Message{
			Role:    "command",
			Content: resp.Command,
		})

		validation := m.validator.Validate(resp.Command, m.workdir)
		if !validation.Allowed {
			return ExecutionCompleteMsg{
				Error:        fmt.Errorf("command blocked: %s", validation.Reason),
				NewMessages:  newMessages,
				NewAIHistory: newAIHistory,
			}
		}

		if validation.NeedsApproval {
			return ApprovalRequestMsg{Command: resp.Command}
		}

		result, err := m.executor.Execute(m.ctx, resp.Command)
		if err != nil {
			return ExecutionCompleteMsg{
				Error:        fmt.Errorf("execution failed: %w", err),
				NewMessages:  newMessages,
				NewAIHistory: newAIHistory,
			}
		}

		return ExecutionCompleteMsg{
			Output:       result.Output,
			NewMessages:  newMessages,
			NewAIHistory: newAIHistory,
		}
	}
}

func (m Model) executeCommand() tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(m.ctx, m.currentCmd)

		if err != nil {
			return ExecutionCompleteMsg{
				Error: fmt.Errorf("execution failed: %w", err),
			}
		}

		if result.ExitCode != 0 {
			return ExecutionCompleteMsg{
				Output: result.Output,
				Error:  fmt.Errorf("command exited with code %d: %s", result.ExitCode, result.Error),
			}
		}

		return ExecutionCompleteMsg{
			Output: result.Output,
		}
	}
}
