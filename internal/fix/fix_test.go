package fix

import (
	"testing"

	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/tools"
)

func TestParseShell(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Shell
		wantErr bool
	}{
		{name: "bash", input: "bash", want: Bash},
		{name: "zsh", input: "zsh", want: Zsh},
		{name: "fish", input: "fish", want: Fish},
		{name: "pwsh", input: "pwsh", want: PowerShell},
		{name: "powershell exe", input: "powershell.exe", want: PowerShell},
		{name: "invalid", input: "cmd", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseShell(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("ParseShell() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseShell() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseShell() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetectShell(t *testing.T) {
	originalGetEnv := getEnv
	t.Cleanup(func() { getEnv = originalGetEnv })

	tests := []struct {
		name    string
		env     map[string]string
		want    Shell
		wantErr bool
	}{
		{name: "shell path zsh", env: map[string]string{"SHELL": "/bin/zsh"}, want: Zsh},
		{name: "powershell via module path", env: map[string]string{"PSModulePath": "x"}, want: PowerShell},
		{name: "undetected", env: map[string]string{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getEnv = func(key string) string { return tt.env[key] }
			got, err := DetectShell()
			if tt.wantErr {
				if err == nil {
					t.Fatal("DetectShell() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("DetectShell() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("DetectShell() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSuggestions(t *testing.T) {
	results := []check.ToolResult{
		{
			Tool:      tools.Tool{Name: "Installed"},
			Installed: true,
			Checks: []check.EnvCheckResult{
				{Check: tools.EnvCheck{Name: "GOOD", ValidValues: []string{"1"}}, Value: "1", Passing: true},
				{Check: tools.EnvCheck{Name: "BAD", ValidValues: []string{"1"}}, Value: "0", Passing: false},
			},
		},
		{
			Tool:      tools.Tool{Name: "Missing"},
			Installed: false,
			Checks: []check.EnvCheckResult{
				{Check: tools.EnvCheck{Name: "MISSING", ValidValues: []string{"0"}}, Value: "", Passing: false},
			},
		},
	}

	got := Suggestions(results, false)
	if len(got) != 2 {
		t.Fatalf("Suggestions(..., false) len = %d, want 2", len(got))
	}
	if got[0].Category != "" {
		t.Fatalf("Suggestions(..., false) first category = %q, want empty category when tool category not set", got[0].Category)
	}
	if got[0].Check.Name != "GOOD" || got[1].Check.Name != "BAD" {
		t.Fatalf("Suggestions(..., false) = %#v, want both installed checks in order", got)
	}

	got = Suggestions(results, true)
	if len(got) != 3 {
		t.Fatalf("Suggestions(..., true) len = %d, want 3", len(got))
	}
}

func TestCommand(t *testing.T) {
	suggestion := Suggestion{
		Check:       tools.EnvCheck{Name: "HOMEBREW_NO_ANALYTICS"},
		TargetValue: "1",
	}

	tests := []struct {
		shell Shell
		want  string
	}{
		{shell: Bash, want: "export HOMEBREW_NO_ANALYTICS='1'"},
		{shell: Zsh, want: "export HOMEBREW_NO_ANALYTICS='1'"},
		{shell: Fish, want: "set -gx HOMEBREW_NO_ANALYTICS '1'"},
		{shell: PowerShell, want: "$Env:HOMEBREW_NO_ANALYTICS = '1'"},
	}

	for _, tt := range tests {
		t.Run(string(tt.shell), func(t *testing.T) {
			if got := Command(tt.shell, suggestion); got != tt.want {
				t.Fatalf("Command() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestComment(t *testing.T) {
	if got := Comment("Package Manager"); got != "# Package Manager Telemetry" {
		t.Fatalf("Comment() = %q, want %q", got, "# Package Manager Telemetry")
	}
}
