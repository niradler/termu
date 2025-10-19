package config

import (
	"os"
	"path/filepath"

	"github.com/niradler/termu/internal/tools"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Model    ModelConfig    `yaml:"model"`
	Security SecurityConfig `yaml:"security"`
	Tools    ToolsConfig    `yaml:"tools"`
	Logging  LoggingConfig  `yaml:"logging"`
	Workdir  string         `yaml:"-"`
}

type ModelConfig struct {
	Provider string `yaml:"provider"`
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Timeout  int    `yaml:"timeout"`
	APIKey   string `yaml:"api_key"`  // API key for OpenAI-compatible providers
	BaseURL  string `yaml:"base_url"` // Base URL for custom OpenAI-compatible endpoints
}

type SecurityConfig struct {
	AllowedCommands   []string `yaml:"allowed_commands"`
	RestrictedFolders []string `yaml:"restricted_folders"`
	AllowedFolders    []string `yaml:"allowed_folders"`
	HighRiskCommands  []string `yaml:"high_risk_commands"`
	BlockedPatterns   []string `yaml:"blocked_patterns"`
	SandboxMode       bool     `yaml:"sandbox_mode"`
	AlwaysApprove     bool     `yaml:"always_approve"`
	MaxToolIterations int      `yaml:"max_tool_iterations"`
}

type ToolsConfig struct {
	AutoInstall  bool `yaml:"auto_install"`
	PreferModern bool `yaml:"prefer_modern"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func DefaultConfig() *Config {
	workdir, _ := os.Getwd()
	return &Config{
		Model: ModelConfig{
			Provider: "ollama",
			Name:     "qwen3",
			Server:   "http://localhost:11434",
			Timeout:  60,
		},
		Security: SecurityConfig{
			AllowedCommands: tools.GetDefaultAllowedCommands(),
			RestrictedFolders: []string{
				"/System", "/Windows", "/etc/passwd", "/etc/shadow",
				"~/.ssh", "~/.aws",
			},
			AllowedFolders: []string{"./"},
			HighRiskCommands: []string{
				"rm", "mv", "chmod", "chown", "sudo",
			},
			BlockedPatterns: []string{
				"rm -rf /", "rm -rf *", ":(){ :|:& };:", "mkfs", "dd if=",
			},
			SandboxMode:       false,
			AlwaysApprove:     false,
			MaxToolIterations: 5,
		},
		Tools: ToolsConfig{
			AutoInstall:  false,
			PreferModern: true,
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  "~/.termu/logs/termu.log",
		},
		Workdir: workdir,
	}
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = findConfigFile()
	}

	if path == "" || !fileExists(path) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func findConfigFile() string {
	candidates := []string{
		".termu.yaml",
		".termu.yml",
	}

	home, _ := os.UserHomeDir()
	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".termu.yaml"),
			filepath.Join(home, ".config", "termu", "config.yaml"),
		)
	}

	for _, path := range candidates {
		if fileExists(path) {
			return path
		}
	}

	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
