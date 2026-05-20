package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/fix"
	"github.com/stennie/notel/internal/tools"
)

func TestExecuteVersionCommandWritesToStdout(t *testing.T) {
	stdout, stderr, err := executeRootCommand(t, "version")
	if err != nil {
		t.Fatalf("execute version: %v", err)
	}

	for _, want := range []string{
		"notel version",
		"Copyright 2026 Stennie Steneker",
		"License: Apache 2.0",
		"GitHub: https://github.com/stennie/notel",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("version stdout missing %q:\n%s", want, stdout)
		}
	}
	if stderr != "" {
		t.Fatalf("version stderr = %q, want empty", stderr)
	}
}

func TestExecuteHelpCommandWritesStyledHelpToStdout(t *testing.T) {
	stdout, stderr, err := executeRootCommand(t, "help")
	if err != nil {
		t.Fatalf("execute help: %v", err)
	}

	for _, want := range []string{
		"notel  —  DevTools Telemetry Auditor",
		"Available Commands",
		"audit",
		"list",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("help stdout missing %q:\n%s", want, stdout)
		}
	}
	if stderr != "" {
		t.Fatalf("help stderr = %q, want empty", stderr)
	}
}

func TestExecuteAuditFixWritesScriptToStdout(t *testing.T) {
	originalRunChecks := runChecks
	originalBuildFixSuggestions := buildFixSuggestions
	originalParseShell := parseShell
	t.Cleanup(func() {
		runChecks = originalRunChecks
		buildFixSuggestions = originalBuildFixSuggestions
		parseShell = originalParseShell
		_ = auditCmd.Flags().Set("fix", "")
		_ = auditCmd.Flags().Set("all", "false")
	})

	runChecks = func() []check.ToolResult { return nil }
	buildFixSuggestions = func(_ []check.ToolResult, _ bool) []fix.Suggestion {
		return []fix.Suggestion{
			{Category: "Package Manager", Check: tools.EnvCheck{Name: "TEST_VAR"}, TargetValue: "1"},
			{Category: "Runtime", Check: tools.EnvCheck{Name: "NODE_NO_TELEMETRY"}, TargetValue: "1"},
		}
	}
	parseShell = func(name string) (fix.Shell, error) { return fix.Bash, nil }

	stdout, stderr, err := executeRootCommand(t, "audit", "--fix=bash")
	if err != nil {
		t.Fatalf("execute audit --fix: %v", err)
	}

	wantStdout := "# Package Manager Telemetry\nexport TEST_VAR='1'\n\n# Runtime Telemetry\nexport NODE_NO_TELEMETRY='1'\n"
	if stdout != wantStdout {
		t.Fatalf("audit --fix stdout = %q, want %q", stdout, wantStdout)
	}
	for _, want := range []string{"Fix Commands", "Shell: bash"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("audit --fix stderr missing %q:\n%s", want, stderr)
		}
	}
	if strings.Contains(stdout, "Fix Commands") {
		t.Fatalf("audit --fix stdout contained contextual output:\n%s", stdout)
	}
}

func TestExecuteAuditFixShorthandWritesScriptToStdout(t *testing.T) {
	originalRunChecks := runChecks
	originalBuildFixSuggestions := buildFixSuggestions
	originalParseShell := parseShell
	t.Cleanup(func() {
		runChecks = originalRunChecks
		buildFixSuggestions = originalBuildFixSuggestions
		parseShell = originalParseShell
		_ = auditCmd.Flags().Set("fix", "")
		_ = auditCmd.Flags().Set("all", "false")
	})

	runChecks = func() []check.ToolResult { return nil }
	buildFixSuggestions = func(_ []check.ToolResult, _ bool) []fix.Suggestion {
		return []fix.Suggestion{
			{Category: "Package Manager", Check: tools.EnvCheck{Name: "TEST_VAR"}, TargetValue: "1"},
		}
	}
	parseShell = func(name string) (fix.Shell, error) { return fix.Bash, nil }

	stdout, stderr, err := executeRootCommand(t, "audit", "-f=bash")
	if err != nil {
		t.Fatalf("execute audit -f: %v", err)
	}

	wantStdout := "# Package Manager Telemetry\nexport TEST_VAR='1'\n"
	if stdout != wantStdout {
		t.Fatalf("audit -f stdout = %q, want %q", stdout, wantStdout)
	}
	for _, want := range []string{"Fix Commands", "Shell: bash"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("audit -f stderr missing %q:\n%s", want, stderr)
		}
	}
}

func TestExecuteAuditFixVerboseAllAlsoRendersAuditToStderr(t *testing.T) {
	originalRunChecks := runChecks
	originalBuildFixSuggestions := buildFixSuggestions
	originalParseShell := parseShell
	t.Cleanup(func() {
		runChecks = originalRunChecks
		buildFixSuggestions = originalBuildFixSuggestions
		parseShell = originalParseShell
		_ = auditCmd.Flags().Set("fix", "")
		_ = auditCmd.Flags().Set("all", "false")
		_ = auditCmd.Flags().Set("verbose", "false")
	})

	runChecks = func() []check.ToolResult {
		return []check.ToolResult{
			{
				Tool: tools.Tool{
					Name:             "Homebrew",
					Category:         "Package Manager",
					DocumentationURL: "https://docs.example.com/homebrew",
					EnvChecks: []tools.EnvCheck{{
						Name:        "HOMEBREW_NO_ANALYTICS",
						ValidValues: []string{"1"},
						Description: "Disables Homebrew analytics",
					}},
				},
				Installed:  true,
				BinaryPath: "/usr/local/bin/brew",
				Checks: []check.EnvCheckResult{{
					Check: tools.EnvCheck{
						Name:        "HOMEBREW_NO_ANALYTICS",
						ValidValues: []string{"1"},
						Description: "Disables Homebrew analytics",
					},
					Value:   "1",
					Passing: true,
				}},
				AllPassing: true,
			},
			{
				Tool: tools.Tool{
					Name:             "Yarn",
					Category:         "Package Manager",
					DocumentationURL: "https://docs.example.com/yarn",
					EnvChecks: []tools.EnvCheck{{
						Name:        "YARN_ENABLE_TELEMETRY",
						ValidValues: []string{"0"},
						Description: "Disables Yarn telemetry",
					}},
				},
				Installed: false,
				Checks: []check.EnvCheckResult{{
					Check: tools.EnvCheck{
						Name:        "YARN_ENABLE_TELEMETRY",
						ValidValues: []string{"0"},
						Description: "Disables Yarn telemetry",
					},
					Value:   "0",
					Passing: true,
				}},
			},
		}
	}
	buildFixSuggestions = func(_ []check.ToolResult, _ bool) []fix.Suggestion {
		return []fix.Suggestion{
			{Category: "Package Manager", Check: tools.EnvCheck{Name: "HOMEBREW_NO_ANALYTICS"}, TargetValue: "1"},
			{Category: "Package Manager", Check: tools.EnvCheck{Name: "YARN_ENABLE_TELEMETRY"}, TargetValue: "0"},
		}
	}
	parseShell = func(name string) (fix.Shell, error) { return fix.Bash, nil }

	stdout, stderr, err := executeRootCommand(t, "audit", "--fix=bash", "--all", "--verbose")
	if err != nil {
		t.Fatalf("execute audit --fix --all --verbose: %v", err)
	}

	wantStdout := "# Package Manager Telemetry\nexport HOMEBREW_NO_ANALYTICS='1'\nexport YARN_ENABLE_TELEMETRY='0'\n"
	if stdout != wantStdout {
		t.Fatalf("audit --fix --all --verbose stdout = %q, want %q", stdout, wantStdout)
	}
	for _, want := range []string{
		"notel  —  DevTools Telemetry Auditor",
		"Homebrew",
		"Yarn",
		"Fix Commands",
		"Shell: bash",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("audit --fix --all --verbose stderr missing %q:\n%s", want, stderr)
		}
	}
	if strings.Contains(stdout, "Homebrew") || strings.Contains(stdout, "Fix Commands") {
		t.Fatalf("audit --fix --all --verbose stdout contained non-script output:\n%s", stdout)
	}
}

func executeRootCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)
	auditCmd.SetOut(&stdout)
	auditCmd.SetErr(&stderr)
	listCmd.SetOut(&stdout)
	listCmd.SetErr(&stderr)
	versionCmd.SetOut(&stdout)
	versionCmd.SetErr(&stderr)
	rootCmd.SetArgs(args)
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		rootCmd.SetArgs(nil)
		auditCmd.SetOut(nil)
		auditCmd.SetErr(nil)
		listCmd.SetOut(nil)
		listCmd.SetErr(nil)
		versionCmd.SetOut(nil)
		versionCmd.SetErr(nil)
	})

	err := rootCmd.Execute()
	return stdout.String(), stderr.String(), err
}
