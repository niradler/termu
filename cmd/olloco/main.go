package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yourusername/olloco/internal/config"
	"github.com/yourusername/olloco/internal/tools"
	"github.com/yourusername/olloco/internal/tui"
)

var (
	configFile  string
	sandboxMode bool
	yoloMode    bool
)

var rootCmd = &cobra.Command{
	Use:   "olloco [command]",
	Short: "Olloco - Secure AI Shell Assistant",
	Long: `Olloco is an intelligent CLI agent that combines the power of local LLMs
with secure shell execution capabilities.`,
	Version: "0.1.0",
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive chat mode",
	Long:  `Start an interactive chat session with the AI assistant`,
	RunE:  runChat,
}

var runCmd = &cobra.Command{
	Use:   "run [prompt]",
	Short: "Execute a one-shot command",
	Long:  `Generate and execute a command from a natural language prompt`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCommand,
}

var installToolsCmd = &cobra.Command{
	Use:   "install-tools",
	Short: "Install modern CLI tools",
	Long:  `Install all recommended modern CLI tools (fd, ripgrep, bat, sd, xsv, jaq, yq, dua, eza)`,
	RunE:  installTools,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", ".olloco.yaml", "config file (default is .olloco.yaml)")
	rootCmd.PersistentFlags().BoolVar(&sandboxMode, "sandbox", false, "run in sandbox mode (dry-run)")
	rootCmd.PersistentFlags().BoolVar(&yoloMode, "yolo", false, "⚠️  YOLO mode: skip ALL security checks (DANGEROUS!)")

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

	if yoloMode || cfg.Security.YoloMode {
		cfg.Security.YoloMode = true
		fmt.Println("⚠️  ⚠️  ⚠️  WARNING: YOLO MODE ENABLED ⚠️  ⚠️  ⚠️")
		fmt.Println("All security checks are DISABLED!")
		fmt.Println("The AI can execute ANY command without validation.")
		fmt.Println("Use at your own risk!")
		fmt.Println()
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
		return fmt.Errorf("failed to start TUI: %w", err)
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

	_ = context.Background()

	fmt.Printf("Processing: %s\n", args[0])

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
