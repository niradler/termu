package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/olloco/internal/config"
)

type Validator struct {
	config *config.Config
}

type ValidationResult struct {
	Allowed       bool
	Reason        string
	RiskLevel     RiskLevel
	NeedsApproval bool
}

type RiskLevel int

const (
	RiskLow RiskLevel = iota
	RiskMedium
	RiskHigh
	RiskCritical
)

func New(cfg *config.Config) *Validator {
	return &Validator{config: cfg}
}

func (v *Validator) Validate(command string, workdir string) *ValidationResult {
	command = strings.TrimSpace(command)

	if v.config.Security.YoloMode {
		return &ValidationResult{
			Allowed:       true,
			RiskLevel:     RiskLow,
			NeedsApproval: false,
		}
	}

	if blocked := v.checkBlockedPatterns(command); blocked != nil {
		return blocked
	}

	if restricted := v.checkRestrictedFolders(command, workdir); restricted != nil {
		return restricted
	}

	if allowed := v.checkAllowedFolders(command, workdir); allowed != nil {
		return allowed
	}

	baseCmd := extractBaseCommand(command)

	if !v.isCommandAllowed(baseCmd) {
		return &ValidationResult{
			Allowed:   false,
			Reason:    fmt.Sprintf("Command '%s' is not in the allowed list", baseCmd),
			RiskLevel: RiskMedium,
		}
	}

	riskLevel := v.assessRisk(command)
	needsApproval := v.needsApproval(command, riskLevel)

	return &ValidationResult{
		Allowed:       true,
		RiskLevel:     riskLevel,
		NeedsApproval: needsApproval,
	}
}

func (v *Validator) checkBlockedPatterns(command string) *ValidationResult {
	for _, pattern := range v.config.Security.BlockedPatterns {
		if strings.Contains(command, pattern) {
			return &ValidationResult{
				Allowed:   false,
				Reason:    fmt.Sprintf("Command contains blocked pattern: %s", pattern),
				RiskLevel: RiskCritical,
			}
		}
	}
	return nil
}

func (v *Validator) checkRestrictedFolders(command, workdir string) *ValidationResult {
	for _, restricted := range v.config.Security.RestrictedFolders {
		expanded := expandPath(restricted)
		if strings.Contains(command, expanded) || strings.HasPrefix(workdir, expanded) {
			return &ValidationResult{
				Allowed:   false,
				Reason:    fmt.Sprintf("Access to restricted folder: %s", restricted),
				RiskLevel: RiskCritical,
			}
		}
	}
	return nil
}

func (v *Validator) checkAllowedFolders(command, workdir string) *ValidationResult {
	if len(v.config.Security.AllowedFolders) == 0 {
		return nil
	}

	for _, allowed := range v.config.Security.AllowedFolders {
		expanded := expandPath(allowed)
		absWorkdir, _ := filepath.Abs(workdir)
		if strings.HasPrefix(absWorkdir, expanded) {
			return nil
		}
	}

	return &ValidationResult{
		Allowed:   false,
		Reason:    "Command would execute outside allowed folders",
		RiskLevel: RiskHigh,
	}
}

func (v *Validator) isCommandAllowed(cmd string) bool {
	if len(v.config.Security.AllowedCommands) == 0 {
		return true
	}

	for _, allowed := range v.config.Security.AllowedCommands {
		if cmd == allowed {
			return true
		}
	}
	return false
}

func (v *Validator) assessRisk(command string) RiskLevel {
	baseCmd := extractBaseCommand(command)

	for _, highRisk := range v.config.Security.HighRiskCommands {
		if strings.HasPrefix(command, highRisk) {
			return RiskHigh
		}
	}

	destructivePatterns := []string{
		"rm ", "del ", "remove", "delete",
		"mv ", "move", "chmod", "chown",
		">", ">>",
	}

	for _, pattern := range destructivePatterns {
		if strings.Contains(command, pattern) {
			return RiskMedium
		}
	}

	if baseCmd == "curl" && strings.Contains(command, "-X POST") {
		return RiskMedium
	}

	return RiskLow
}

func (v *Validator) needsApproval(command string, risk RiskLevel) bool {
	if v.config.Security.AlwaysApprove {
		return true
	}

	if risk >= RiskMedium {
		return true
	}

	baseCmd := extractBaseCommand(command)
	for _, highRisk := range v.config.Security.HighRiskCommands {
		if baseCmd == highRisk {
			return true
		}
	}

	return false
}

func extractBaseCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	return filepath.Base(parts[0])
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	abs, _ := filepath.Abs(path)
	return abs
}
