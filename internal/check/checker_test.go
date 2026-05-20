package check

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stennie/notel/internal/tools"
)

func TestCheckEnvVar(t *testing.T) {
	t.Setenv("TEST_TELEMETRY", "1")
	t.Setenv("TEST_ANY_NON_EMPTY", "enabled")

	tests := []struct {
		name  string
		check tools.EnvCheck
		want  EnvCheckResult
	}{
		{
			name: "matches accepted value",
			check: tools.EnvCheck{
				Name:        "TEST_TELEMETRY",
				ValidValues: []string{"1", "true"},
			},
			want: EnvCheckResult{
				Check: tools.EnvCheck{
					Name:        "TEST_TELEMETRY",
					ValidValues: []string{"1", "true"},
				},
				Value:   "1",
				Passing: true,
			},
		},
		{
			name: "accepts any non-empty value when valid values omitted",
			check: tools.EnvCheck{
				Name: "TEST_ANY_NON_EMPTY",
			},
			want: EnvCheckResult{
				Check: tools.EnvCheck{
					Name: "TEST_ANY_NON_EMPTY",
				},
				Value:   "enabled",
				Passing: true,
			},
		},
		{
			name: "fails when env var missing",
			check: tools.EnvCheck{
				Name:        "TEST_MISSING",
				ValidValues: []string{"1"},
			},
			want: EnvCheckResult{
				Check: tools.EnvCheck{
					Name:        "TEST_MISSING",
					ValidValues: []string{"1"},
				},
				Value:   "",
				Passing: false,
			},
		},
		{
			name: "fails when env var has unexpected value",
			check: tools.EnvCheck{
				Name:        "TEST_TELEMETRY",
				ValidValues: []string{"0"},
			},
			want: EnvCheckResult{
				Check: tools.EnvCheck{
					Name:        "TEST_TELEMETRY",
					ValidValues: []string{"0"},
				},
				Value:   "1",
				Passing: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkEnvVar(tt.check)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("checkEnvVar() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAuditToolInstalledAndPassing(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() {
		lookPath = originalLookPath
	})

	lookPath = func(file string) (string, error) {
		if file != "brew" {
			t.Fatalf("lookPath called with %q, want brew", file)
		}
		return "/usr/local/bin/brew", nil
	}

	t.Setenv("HOMEBREW_NO_ANALYTICS", "1")

	got := auditTool(tools.Tool{
		Name:   "Homebrew",
		Binary: "brew",
		EnvChecks: []tools.EnvCheck{{
			Name:        "HOMEBREW_NO_ANALYTICS",
			ValidValues: []string{"1"},
		}},
	})

	if !got.Installed {
		t.Fatal("auditTool() did not mark installed tool as installed")
	}
	if got.BinaryPath != "/usr/local/bin/brew" {
		t.Fatalf("auditTool() binary path = %q, want %q", got.BinaryPath, "/usr/local/bin/brew")
	}
	if !got.AllPassing {
		t.Fatal("auditTool() did not mark passing installed tool as passing")
	}
	if len(got.Checks) != 1 || !got.Checks[0].Passing {
		t.Fatalf("auditTool() checks = %#v, want one passing check", got.Checks)
	}
}

func TestAuditToolNotInstalledNeverAllPassing(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() {
		lookPath = originalLookPath
	})

	lookPath = func(file string) (string, error) {
		return "", errors.New("not found")
	}

	t.Setenv("HOMEBREW_NO_ANALYTICS", "1")

	got := auditTool(tools.Tool{
		Name:   "Homebrew",
		Binary: "brew",
		EnvChecks: []tools.EnvCheck{{
			Name:        "HOMEBREW_NO_ANALYTICS",
			ValidValues: []string{"1"},
		}},
	})

	if got.Installed {
		t.Fatal("auditTool() marked missing binary as installed")
	}
	if got.AllPassing {
		t.Fatal("auditTool() marked missing binary as all passing")
	}
	if len(got.Checks) != 1 || !got.Checks[0].Passing {
		t.Fatalf("auditTool() checks = %#v, want env check result preserved", got.Checks)
	}
}

func TestRunForAndSummarise(t *testing.T) {
	originalLookPath := lookPath
	t.Cleanup(func() {
		lookPath = originalLookPath
	})

	lookPath = func(file string) (string, error) {
		switch file {
		case "brew":
			return "/usr/local/bin/brew", nil
		case "node":
			return "", errors.New("not found")
		default:
			t.Fatalf("unexpected binary lookup: %q", file)
			return "", nil
		}
	}

	t.Setenv("HOMEBREW_NO_ANALYTICS", "1")
	t.Setenv("NODE_NO_TELEMETRY", "1")

	results := RunFor([]tools.Tool{
		{
			Name:   "Homebrew",
			Binary: "brew",
			EnvChecks: []tools.EnvCheck{{
				Name:        "HOMEBREW_NO_ANALYTICS",
				ValidValues: []string{"1"},
			}},
		},
		{
			Name:   "Node.js",
			Binary: "node",
			EnvChecks: []tools.EnvCheck{{
				Name:        "NODE_NO_TELEMETRY",
				ValidValues: []string{"1"},
			}},
		},
		{
			Name: "Env Only Tool",
			EnvChecks: []tools.EnvCheck{{
				Name: "EMPTY_ACCEPTS_NOTHING",
			}},
		},
	})

	if len(results) != 3 {
		t.Fatalf("RunFor() len = %d, want 3", len(results))
	}
	if results[0].Tool.Name != "Homebrew" || results[1].Tool.Name != "Node.js" || results[2].Tool.Name != "Env Only Tool" {
		t.Fatalf("RunFor() did not preserve tool ordering: %#v", results)
	}

	summary := Summarise(results)
	want := Summary{
		Total:     3,
		Installed: 2,
		Passing:   1,
		Failing:   1,
	}
	if summary != want {
		t.Fatalf("Summarise() = %#v, want %#v", summary, want)
	}
}
