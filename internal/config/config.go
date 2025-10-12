package config

import (
	"os"
	"path/filepath"

	"github.com/yourusername/olloco/internal/tools"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Model    ModelConfig    `yaml:"model"`
	Security SecurityConfig `yaml:"security"`
	Tools    ToolsConfig    `yaml:"tools"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type ModelConfig struct {
	Provider string `yaml:"provider"`
	Name     string `yaml:"name"`
	Server   string `yaml:"server"`
	Timeout  int    `yaml:"timeout"`
}

type SecurityConfig struct {
	AllowedCommands   []string `yaml:"allowed_commands"`
	RestrictedFolders []string `yaml:"restricted_folders"`
	AllowedFolders    []string `yaml:"allowed_folders"`
	HighRiskCommands  []string `yaml:"high_risk_commands"`
	BlockedPatterns   []string `yaml:"blocked_patterns"`
	SandboxMode       bool     `yaml:"sandbox_mode"`
	AlwaysApprove     bool     `yaml:"always_approve"`
	YoloMode          bool     `yaml:"yolo_mode"`
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
			SandboxMode:   false,
			AlwaysApprove: false,
			YoloMode:      false,
		},
		Tools: ToolsConfig{
			AutoInstall:  false,
			PreferModern: true,
		},
		Logging: LoggingConfig{
			Level: "info",
			File:  "~/.olloco/logs/olloco.log",
		},
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
		".olloco.yaml",
		".olloco.yml",
	}

	home, _ := os.UserHomeDir()
	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".olloco.yaml"),
			filepath.Join(home, ".config", "olloco", "config.yaml"),
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
