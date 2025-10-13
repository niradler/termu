package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/niradler/termu/internal/config"
	"github.com/niradler/termu/internal/tools"
	"github.com/niradler/termu/internal/tui"
	"github.com/spf13/cobra"
)

var (
	configFile  string
	sandboxMode bool
)

var rootCmd = &cobra.Command{
	Use:     "termu [command]",
	Short:   "termu - Your Terminal Sidekick",
	Long:    `termu is a terminal chat agent that helps you accomplish tasks using shell commands.`,
	Version: "0.2.0",
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive chat session",
	Long:  `Start a chat session with termu in your current directory`,
	RunE:  runChat,
}

var runCmd = &cobra.Command{
	Use:   "run [prompt]",
	Short: "Execute a quick command",
	Long:  `Generate and execute a command from natural language`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCommand,
}

var installToolsCmd = &cobra.Command{
	Use:   "install-tools",
	Short: "Install cross-platform CLI tools",
	Long:  `Install termu's favorite cross-platform CLI tools`,
	RunE:  installTools,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".termu.yaml", "config file")
	rootCmd.PersistentFlags().BoolVar(&sandboxMode, "sandbox", false, "sandbox mode (dry-run)")

	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installToolsCmd)
}

func runChat(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if sandboxMode || cfg.Security.SandboxMode {
		sandboxMode = true
	}

	ctx := context.Background()
	model, err := tui.NewModel(ctx, cfg, sandboxMode)
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start chat: %w", err)
	}

	return nil
}

func runCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if sandboxMode || cfg.Security.SandboxMode {
		sandboxMode = true
	}

	fmt.Printf("Processing: %s\n", args[0])
	fmt.Println("Quick run mode coming soon. Use 'termu chat' for now.")

	return nil
}

func installTools(cmd *cobra.Command, args []string) error {
	installer := tools.NewInstaller()
	return installer.InstallAll()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
