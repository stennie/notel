// Package check performs detection and telemetry-setting audits for developer tools.
package check

import (
	"os"
	"os/exec"

	"github.com/stennie/notel/internal/tools"
)

var (
	lookPath = exec.LookPath
	getEnv   = os.Getenv
)

// EnvCheckResult holds the outcome of checking a single environment variable.
type EnvCheckResult struct {
	Check   tools.EnvCheck
	Value   string // actual value found in the environment (empty if not set)
	Passing bool   // true when the value satisfies the opt-out requirement
}

// ToolResult holds the full audit result for a single tool.
type ToolResult struct {
	Tool       tools.Tool
	Installed  bool
	BinaryPath string
	Checks     []EnvCheckResult
	AllPassing bool // true only when installed AND every env check passes
}

// Run checks every tool in the default registry.
func Run() []ToolResult {
	return RunFor(tools.Registry())
}

// RunFor checks a caller-supplied slice of tools (useful for testing/filtering).
func RunFor(toolList []tools.Tool) []ToolResult {
	results := make([]ToolResult, 0, len(toolList))
	for _, t := range toolList {
		results = append(results, auditTool(t))
	}
	return results
}

// Summary holds aggregate counts derived from a slice of ToolResults.
type Summary struct {
	Total     int // total tools in the registry
	Installed int // installed tools
	Passing   int // installed tools where all env checks pass
	Failing   int // installed tools with at least one failing check
}

// Summarise computes aggregate counts from a result set.
func Summarise(results []ToolResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Installed {
			s.Installed++
			if r.AllPassing {
				s.Passing++
			} else {
				s.Failing++
			}
		}
	}
	return s
}

// ── internal ────────────────────────────────────────────────────────────────

func auditTool(t tools.Tool) ToolResult {
	result := ToolResult{Tool: t}

	if t.Binary != "" {
		path, err := lookPath(t.Binary)
		if err == nil {
			result.Installed = true
			result.BinaryPath = path
		}
	} else {
		// No binary specified — treat as always installed (e.g. env-only tools).
		result.Installed = true
	}

	allPassing := true
	for _, chk := range t.EnvChecks {
		cr := checkEnvVar(chk)
		result.Checks = append(result.Checks, cr)
		if !cr.Passing {
			allPassing = false
		}
	}
	// Only mark AllPassing when actually installed.
	result.AllPassing = result.Installed && allPassing

	return result
}

func checkEnvVar(chk tools.EnvCheck) EnvCheckResult {
	value := getEnv(chk.Name)
	passing := false

	if value != "" {
		if len(chk.ValidValues) == 0 {
			// Any non-empty value satisfies the check.
			passing = true
		} else {
			for _, valid := range chk.ValidValues {
				if value == valid {
					passing = true
					break
				}
			}
		}
	}

	return EnvCheckResult{
		Check:   chk,
		Value:   value,
		Passing: passing,
	}
}
