package tui

import "github.com/charmbracelet/lipgloss"

var (
	BaseStyle = lipgloss.NewStyle().
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(0, 1)

	CommandStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a1a")).
			Foreground(lipgloss.Color("#00ff00")).
			Padding(1).
			MarginTop(1).
			MarginBottom(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4"))

	OutputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0f0f0f")).
			Foreground(lipgloss.Color("#d0d0d0")).
			Padding(1).
			MarginTop(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#404040"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF"))

	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#404040"))

	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	ApprovalStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FFD700")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(1).
			MarginTop(1).
			MarginBottom(1).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#FFD700"))

	SandboxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FF8C00")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1)

	YoloStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FF0000")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Blink(true).
			Padding(0, 1)
)
