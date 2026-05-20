package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/fix"
	"github.com/stennie/notel/internal/tools"
)

func TestResolveFixShell(t *testing.T) {
	originalDetectShell := detectShell
	originalParseShell := parseShell
	t.Cleanup(func() {
		detectShell = originalDetectShell
		parseShell = originalParseShell
	})

	detectShell = func() (fix.Shell, error) { return fix.Zsh, nil }
	parseShell = func(name string) (fix.Shell, error) {
		if name != "powershell" {
			t.Fatalf("parseShell called with %q, want powershell", name)
		}
		return fix.PowerShell, nil
	}

	got, err := resolveFixShell("auto")
	if err != nil {
		t.Fatalf("resolveFixShell(auto) error = %v", err)
	}
	if got != fix.Zsh {
		t.Fatalf("resolveFixShell(auto) = %q, want %q", got, fix.Zsh)
	}

	got, err = resolveFixShell("powershell")
	if err != nil {
		t.Fatalf("resolveFixShell(powershell) error = %v", err)
	}
	if got != fix.PowerShell {
		t.Fatalf("resolveFixShell(powershell) = %q, want %q", got, fix.PowerShell)
	}
}

func TestAuditFixUsesDetectedShellAndShowAll(t *testing.T) {
	originalShowAll := auditShowAll
	originalVerbose := auditVerbose
	originalFix := auditFix
	originalRunChecks := runChecks
	originalDetectShell := detectShell
	originalBuildFixSuggestions := buildFixSuggestions
	originalPrintFixCommands := printFixCommands

	t.Cleanup(func() {
		auditShowAll = originalShowAll
		auditVerbose = originalVerbose
		auditFix = originalFix
		runChecks = originalRunChecks
		detectShell = originalDetectShell
		buildFixSuggestions = originalBuildFixSuggestions
		printFixCommands = originalPrintFixCommands
		_ = auditCmd.Flags().Set("all", "false")
		_ = auditCmd.Flags().Set("verbose", "false")
		_ = auditCmd.Flags().Set("fix", "")
	})

	runChecks = func() []check.ToolResult {
		return []check.ToolResult{{Installed: true}}
	}
	detectShell = func() (fix.Shell, error) { return fix.Fish, nil }

	var capturedShowAll bool
	buildFixSuggestions = func(_ []check.ToolResult, showAll bool) []fix.Suggestion {
		capturedShowAll = showAll
		return []fix.Suggestion{{ToolName: "Homebrew"}}
	}

	var capturedShell fix.Shell
	printFixCommands = func(_ io.Writer, _ io.Writer, shell fix.Shell, suggestions []fix.Suggestion) {
		capturedShell = shell
		if len(suggestions) != 1 {
			t.Fatalf("printFixCommands suggestions len = %d, want 1", len(suggestions))
		}
	}

	if err := auditCmd.Flags().Set("all", "true"); err != nil {
		t.Fatalf("setting all flag: %v", err)
	}
	if err := auditCmd.Flags().Set("fix", "auto"); err != nil {
		t.Fatalf("setting fix flag: %v", err)
	}

	if err := auditCmd.RunE(auditCmd, nil); err != nil {
		t.Fatalf("auditCmd.RunE() error = %v", err)
	}
	if !capturedShowAll {
		t.Fatal("buildFixSuggestions did not receive showAll=true")
	}
	if capturedShell != fix.Fish {
		t.Fatalf("printFixCommands shell = %q, want %q", capturedShell, fix.Fish)
	}
}

func TestAuditFixWithExplicitShellBypassesDetection(t *testing.T) {
	originalFix := auditFix
	originalRunChecks := runChecks
	originalDetectShell := detectShell
	originalParseShell := parseShell
	originalBuildFixSuggestions := buildFixSuggestions
	originalPrintFixCommands := printFixCommands

	t.Cleanup(func() {
		auditFix = originalFix
		runChecks = originalRunChecks
		detectShell = originalDetectShell
		parseShell = originalParseShell
		buildFixSuggestions = originalBuildFixSuggestions
		printFixCommands = originalPrintFixCommands
		_ = auditCmd.Flags().Set("fix", "")
	})

	runChecks = func() []check.ToolResult { return nil }
	detectShell = func() (fix.Shell, error) {
		t.Fatal("detectShell should not be called for explicit fix shell")
		return "", nil
	}
	parseShell = func(name string) (fix.Shell, error) {
		if name != "bash" {
			t.Fatalf("parseShell called with %q, want bash", name)
		}
		return fix.Bash, nil
	}
	buildFixSuggestions = func(_ []check.ToolResult, _ bool) []fix.Suggestion { return nil }

	called := false
	printFixCommands = func(_ io.Writer, _ io.Writer, shell fix.Shell, _ []fix.Suggestion) {
		called = true
		if shell != fix.Bash {
			t.Fatalf("printFixCommands shell = %q, want %q", shell, fix.Bash)
		}
	}

	if err := auditCmd.Flags().Set("fix", "bash"); err != nil {
		t.Fatalf("setting fix flag: %v", err)
	}
	if err := auditCmd.RunE(auditCmd, nil); err != nil {
		t.Fatalf("auditCmd.RunE() error = %v", err)
	}
	if !called {
		t.Fatal("printFixCommands was not called")
	}
}

func TestAuditFixFlagHasFShorthand(t *testing.T) {
	flag := auditCmd.Flags().Lookup("fix")
	if flag == nil {
		t.Fatal("fix flag is not registered")
	}
	if flag.Shorthand != "f" {
		t.Fatalf("fix flag shorthand = %q, want %q", flag.Shorthand, "f")
	}
}

func TestAuditFixStdoutContainsOnlyCommands(t *testing.T) {
	originalFix := auditFix
	originalRunChecks := runChecks
	originalParseShell := parseShell
	originalBuildFixSuggestions := buildFixSuggestions
	originalPrintFixCommands := printFixCommands

	t.Cleanup(func() {
		auditFix = originalFix
		runChecks = originalRunChecks
		parseShell = originalParseShell
		buildFixSuggestions = originalBuildFixSuggestions
		printFixCommands = originalPrintFixCommands
		auditCmd.SetOut(nil)
		auditCmd.SetErr(nil)
		_ = auditCmd.Flags().Set("fix", "")
	})

	runChecks = func() []check.ToolResult { return nil }
	parseShell = func(name string) (fix.Shell, error) { return fix.Bash, nil }
	buildFixSuggestions = func(_ []check.ToolResult, _ bool) []fix.Suggestion {
		return []fix.Suggestion{
			{Category: "Package Manager", Check: tools.EnvCheck{Name: "TEST_VAR"}, TargetValue: "1"},
			{Category: "Package Manager", Check: tools.EnvCheck{Name: "ALREADY_SET"}, Current: "1", TargetValue: "1"},
			{Category: "Framework", Check: tools.EnvCheck{Name: "FRAMEWORK_VAR"}, TargetValue: "0"},
		}
	}

	original := printFixCommands
	printFixCommands = original

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	auditCmd.SetOut(&stdout)
	auditCmd.SetErr(&stderr)

	if err := auditCmd.Flags().Set("fix", "bash"); err != nil {
		t.Fatalf("setting fix flag: %v", err)
	}
	if err := auditCmd.RunE(auditCmd, nil); err != nil {
		t.Fatalf("auditCmd.RunE() error = %v", err)
	}

	if got := strings.TrimSpace(stdout.String()); got != "# Package Manager Telemetry\nexport TEST_VAR='1'\nexport ALREADY_SET='1'\n\n# Framework Telemetry\nexport FRAMEWORK_VAR='0'" {
		t.Fatalf("stdout = %q, want full env block", got)
	}
	if !strings.Contains(stderr.String(), "Fix Commands") {
		t.Fatalf("stderr missing contextual output:\n%s", stderr.String())
	}
}
