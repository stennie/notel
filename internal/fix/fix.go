// Package fix builds shell-specific env-var commands to disable telemetry.
package fix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/tools"
)

// Shell identifies a supported shell syntax for fix commands.
type Shell string

const (
	Bash       Shell = "bash"
	Zsh        Shell = "zsh"
	Fish       Shell = "fish"
	PowerShell Shell = "powershell"
)

var getEnv = os.Getenv

// Suggestion represents one env var that should be set to disable telemetry.
type Suggestion struct {
	Category    string
	ToolName    string
	Installed   bool
	Current     string
	Check       tools.EnvCheck
	TargetValue string
}

// ParseShell normalizes a shell name and validates that it is supported.
func ParseShell(name string) (Shell, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "bash":
		return Bash, nil
	case "zsh":
		return Zsh, nil
	case "fish":
		return Fish, nil
	case "powershell", "pwsh", "pwsh.exe", "powershell.exe":
		return PowerShell, nil
	default:
		return "", fmt.Errorf("unsupported shell %q (supported: bash, fish, powershell, zsh)", name)
	}
}

// DetectShell infers the current shell from the environment.
func DetectShell() (Shell, error) {
	if shellPath := strings.TrimSpace(getEnv("SHELL")); shellPath != "" {
		return ParseShell(filepath.Base(shellPath))
	}
	if strings.TrimSpace(getEnv("PSModulePath")) != "" {
		return PowerShell, nil
	}
	if comspec := strings.ToLower(strings.TrimSpace(getEnv("ComSpec"))); strings.Contains(comspec, "powershell") || strings.Contains(comspec, "pwsh") {
		return PowerShell, nil
	}
	return "", fmt.Errorf("could not detect current shell (supported: bash, fish, powershell, zsh)")
}

// Suggestions returns env vars that should be included in a reusable shell env block.
func Suggestions(results []check.ToolResult, showAll bool) []Suggestion {
	out := make([]Suggestion, 0)
	for _, result := range results {
		if !showAll && !result.Installed {
			continue
		}
		for _, envCheck := range result.Checks {
			out = append(out, Suggestion{
				Category:    result.Tool.Category,
				ToolName:    result.Tool.Name,
				Installed:   result.Installed,
				Current:     envCheck.Value,
				Check:       envCheck.Check,
				TargetValue: preferredValue(envCheck.Check),
			})
		}
	}
	return out
}

// Comment renders a shell-safe comment heading for a category block.
func Comment(category string) string {
	return fmt.Sprintf("# %s Telemetry", category)
}

// Command renders one shell-specific command for the suggested env var.
func Command(shell Shell, suggestion Suggestion) string {
	value := quoteValue(shell, suggestion.TargetValue)
	switch shell {
	case Bash, Zsh:
		return fmt.Sprintf("export %s=%s", suggestion.Check.Name, value)
	case Fish:
		return fmt.Sprintf("set -gx %s %s", suggestion.Check.Name, value)
	case PowerShell:
		return fmt.Sprintf("$Env:%s = %s", suggestion.Check.Name, value)
	default:
		return ""
	}
}

func preferredValue(check tools.EnvCheck) string {
	if len(check.ValidValues) > 0 {
		return check.ValidValues[0]
	}
	return "1"
}

func quoteValue(shell Shell, value string) string {
	escaped := strings.ReplaceAll(value, "'", "\\'")
	switch shell {
	case PowerShell:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''"))
	default:
		return fmt.Sprintf("'%s'", escaped)
	}
}
