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
	StateIterating
)

type Message struct {
	Role    string
	Content string
}

type Model struct {
	ctx            context.Context
	state          SessionState
	textarea       textarea.Model
	viewport       viewport.Model
	messages       []Message
	aiHistory      []ai.Message
	currentCmd     string
	currentInput   string
	iterationCount int
	maxIterations  int
	retryCount     int
	maxRetries     int
	width          int
	height         int
	sandboxMode    bool
	mdRenderer     *glamour.TermRenderer
	agent          *agent.Agent
	validator      *security.Validator
	executor       *shell.Executor
	workdir        string
}

type AgentResponseMsg struct {
	Response *agent.Response
	Error    error
}

type CommandExecutedMsg struct {
	Command string
	Output  string
	Error   error
}

type ApprovalRequestMsg struct {
	Command string
}

type IterationCompleteMsg struct {
	FinalText string
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
		ctx:            ctx,
		state:          StateInput,
		textarea:       ta,
		viewport:       vp,
		messages:       []Message{},
		aiHistory:      []ai.Message{},
		iterationCount: 0,
		maxIterations:  cfg.Security.MaxToolIterations,
		retryCount:     0,
		maxRetries:     2,
		width:          80,
		height:         24,
		sandboxMode:    sandboxMode,
		mdRenderer:     renderer,
		agent:          ag,
		validator:      security.New(cfg),
		executor:       shell.New(workdir, sandboxMode),
		workdir:        workdir,
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
				m.currentInput = userInput
				m.iterationCount = 0
				m.retryCount = 0
				m.state = StateThinking
				m.updateViewport()
				return m, m.callAgent()
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

	case AgentResponseMsg:
		if msg.Error != nil {
			m.messages = append(m.messages, Message{
				Role:    "error",
				Content: fmt.Sprintf("Agent error: %v", msg.Error),
			})
			m.state = StateInput
			m.updateViewport()
			return m, nil
		}

		m.messages = append(m.messages, Message{
			Role:    "assistant",
			Content: msg.Response.Text,
		})

		m.aiHistory = append(m.aiHistory, ai.Message{
			Role:    ai.RoleModel,
			Content: []*ai.Part{ai.NewTextPart(msg.Response.Text)},
		})

		if msg.Response.Command == "" {
			m.state = StateInput
			m.updateViewport()
			return m, nil
		}

		m.messages = append(m.messages, Message{
			Role:    "command",
			Content: msg.Response.Command,
		})

		validation := m.validator.Validate(msg.Response.Command, m.workdir)
		if !validation.Allowed {
			m.messages = append(m.messages, Message{
				Role:    "error",
				Content: fmt.Sprintf("Command blocked: %s", validation.Reason),
			})
			m.state = StateInput
			m.updateViewport()
			return m, nil
		}

		m.currentCmd = msg.Response.Command
		m.retryCount = 0

		if validation.NeedsApproval {
			m.state = StateApproval
			m.updateViewport()
			return m, nil
		}

		m.state = StateExecuting
		m.updateViewport()
		return m, m.executeCommand()

	case CommandExecutedMsg:
		if msg.Error != nil {
			m.messages = append(m.messages, Message{
				Role:    "error",
				Content: fmt.Sprintf("❌ Error: %v", msg.Error),
			})

			if m.retryCount < m.maxRetries {
				m.retryCount++
				m.messages = append(m.messages, Message{
					Role:    "system",
					Content: fmt.Sprintf("🔄 Retry %d/%d - letting termu try to fix it...", m.retryCount, m.maxRetries),
				})

				errorFeedback := fmt.Sprintf("Command failed: %s\n\nError:\n%s\n\nPlease try a different approach.", msg.Command, msg.Output)
				m.aiHistory = append(m.aiHistory, ai.Message{
					Role:    ai.RoleUser,
					Content: []*ai.Part{ai.NewTextPart(errorFeedback)},
				})

				m.state = StateIterating
				m.updateViewport()
				return m, m.callAgent()
			}

			m.messages = append(m.messages, Message{
				Role:    "system",
				Content: fmt.Sprintf("❌ Failed after %d retries. Task incomplete.", m.maxRetries),
			})
			m.state = StateInput
			m.updateViewport()
			return m, nil
		}

		m.iterationCount++
		m.retryCount = 0

		m.messages = append(m.messages, Message{
			Role:    "output",
			Content: msg.Output,
		})

		toolOutput := fmt.Sprintf("Command executed successfully: %s\n\nOutput:\n%s", msg.Command, msg.Output)
		m.aiHistory = append(m.aiHistory, ai.Message{
			Role:    ai.RoleUser,
			Content: []*ai.Part{ai.NewTextPart(toolOutput)},
		})

		if m.iterationCount >= m.maxIterations {
			m.messages = append(m.messages, Message{
				Role:    "system",
				Content: fmt.Sprintf("⚠️ Max iterations (%d) reached", m.maxIterations),
			})
			m.state = StateInput
			m.updateViewport()
			return m, nil
		}

		m.state = StateIterating
		m.updateViewport()
		return m, m.callAgent()

	case ApprovalRequestMsg:
		m.currentCmd = msg.Command
		m.state = StateApproval
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
		status := "🤔 termu is thinking..."
		if m.iterationCount > 0 {
			status = fmt.Sprintf("🤔 termu is thinking... [%d/%d]", m.iterationCount+1, m.maxIterations)
		}
		b.WriteString(InfoStyle.Render(status))
	} else if m.state == StateIterating {
		b.WriteString(InfoStyle.Render(
			fmt.Sprintf("🔄 termu analyzing results... [%d/%d]", m.iterationCount+1, m.maxIterations)))
	} else if m.state == StateApproval {
		b.WriteString(m.renderApproval())
	} else if m.state == StateExecuting {
		status := "⚡ Executing command..."
		if m.iterationCount > 0 {
			status = fmt.Sprintf("⚡ Executing... [%d/%d]", m.iterationCount, m.maxIterations)
		}
		b.WriteString(InfoStyle.Render(status))
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

	var status string
	if m.state == StateIterating || m.state == StateExecuting {
		status = HelpStyle.Render(fmt.Sprintf(" [%d/%d]", m.iterationCount, m.maxIterations))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, title, " ", mode, status)
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

func (m Model) callAgent() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.agent.Generate(m.ctx, m.currentInput, m.aiHistory)
		return AgentResponseMsg{
			Response: resp,
			Error:    err,
		}
	}
}

func (m Model) executeCommand() tea.Cmd {
	return func() tea.Msg {
		result, err := m.executor.Execute(m.ctx, m.currentCmd)

		var cmdErr error
		if err != nil {
			cmdErr = fmt.Errorf("execution failed: %w", err)
		} else if result.ExitCode != 0 {
			cmdErr = fmt.Errorf("command exited with code %d", result.ExitCode)
		}

		output := result.Output
		if result.Error != "" && cmdErr != nil {
			output = result.Output + "\nError: " + result.Error
		}

		return CommandExecutedMsg{
			Command: m.currentCmd,
			Output:  output,
			Error:   cmdErr,
		}
	}
}
