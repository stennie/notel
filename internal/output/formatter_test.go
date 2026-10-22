package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/tools"
)

func TestPrintReportResultsIncludesDocumentationURL(t *testing.T) {
	output := captureStdout(t, func() {
		PrintReportResults([]check.ToolResult{
			{
				Tool: tools.Tool{
					Name:             "Homebrew",
					Category:         "Package Manager",
					DocumentationURL: "https://docs.example.com/homebrew",
					EnvChecks: []tools.EnvCheck{{
						Name:        "HOMEBREW_NO_ANALYTICS",
						ValidValues: []string{"1"},
						Description: "Disables analytics",
					}},
				},
				Installed:  true,
				BinaryPath: "/usr/local/bin/brew",
				Checks: []check.EnvCheckResult{{
					Check: tools.EnvCheck{
						Name:        "HOMEBREW_NO_ANALYTICS",
						ValidValues: []string{"1"},
						Description: "Disables analytics",
					},
					Value:   "1",
					Passing: true,
				}},
				AllPassing: true,
			},
			{
				Tool: tools.Tool{
					Name:             "Netlify CLI",
					Category:         "Deployment",
					DocumentationURL: "https://docs.example.com/netlify",
				},
				Installed: false,
			},
		}, true)
	})

	for _, want := range []string{
		"Docs:",
		"https://docs.example.com/homebrew",
		"https://docs.example.com/netlify",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("PrintReportResults output missing %q:\n%s", want, output)
		}
	}
}

func TestPrintReportResultsFiltersUninstalledToolsUnlessShowAll(t *testing.T) {
	results := []check.ToolResult{
		{
			Tool: tools.Tool{
				Name:             "Installed Tool",
				Category:         "Runtime",
				DocumentationURL: "https://docs.example.com/installed",
			},
			Installed:  true,
			BinaryPath: "/usr/local/bin/tool",
		},
		{
			Tool: tools.Tool{
				Name:             "Missing Tool",
				Category:         "Runtime",
				DocumentationURL: "https://docs.example.com/missing",
			},
			Installed: false,
		},
	}

	installedOnly := captureStdout(t, func() {
		PrintReportResults(results, false)
	})
	if strings.Contains(installedOnly, "Missing Tool") {
		t.Fatalf("PrintReportResults(..., false) included uninstalled tool:\n%s", installedOnly)
	}
	if !strings.Contains(installedOnly, "Installed Tool") {
		t.Fatalf("PrintReportResults(..., false) omitted installed tool:\n%s", installedOnly)
	}

	allTools := captureStdout(t, func() {
		PrintReportResults(results, true)
	})
	if !strings.Contains(allTools, "Missing Tool") {
		t.Fatalf("PrintReportResults(..., true) omitted uninstalled tool:\n%s", allTools)
	}
}

func TestPrintReportResultsWarnsOnUnexpectedValue(t *testing.T) {
	output := captureStdout(t, func() {
		PrintReportResults([]check.ToolResult{
			{
				Tool: tools.Tool{
					Name:             "Homebrew",
					Category:         "Package Manager",
					DocumentationURL: "https://docs.example.com/homebrew",
				},
				Installed:  true,
				BinaryPath: "/usr/local/bin/brew",
				Checks: []check.EnvCheckResult{{
					Check: tools.EnvCheck{
						Name:        "HOMEBREW_NO_ANALYTICS",
						ValidValues: []string{"1"},
						Description: "Disables Homebrew analytics reporting",
					},
					Value:   "0",
					Passing: false,
				}},
			},
		}, true)
	})

	for _, want := range []string{
		"unexpected value; assuming telemetry enabled",
		"HOMEBREW_NO_ANALYTICS=0",
		"→ set HOMEBREW_NO_ANALYTICS=1",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("PrintReportResults output missing %q:\n%s", want, output)
		}
	}
}

func TestPrintCheckResultsShowsDetectedEnvForUninstalledToolWhenIncluded(t *testing.T) {
	output := captureStdout(t, func() {
		PrintCheckResults([]check.ToolResult{
			{
				Tool: tools.Tool{
					Name:     "Yarn",
					Category: "Package Manager",
				},
				Installed: false,
				Checks: []check.EnvCheckResult{
					{
						Check:   tools.EnvCheck{Name: "YARN_ENABLE_TELEMETRY"},
						Value:   "0",
						Passing: true,
					},
				},
			},
		}, true)
	})

	for _, want := range []string{
		"Yarn",
		"not installed",
		"YARN_ENABLE_TELEMETRY=0",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("PrintCheckResults output missing %q:\n%s", want, output)
		}
	}
}

func TestPrintReportResultsShowsDetectedEnvForUninstalledToolWhenIncluded(t *testing.T) {
	output := captureStdout(t, func() {
		PrintReportResults([]check.ToolResult{
			{
				Tool: tools.Tool{
					Name:             "Yarn",
					Category:         "Package Manager",
					DocumentationURL: "https://docs.example.com/yarn",
				},
				Installed: false,
				Checks: []check.EnvCheckResult{
					{
						Check: tools.EnvCheck{
							Name:        "YARN_ENABLE_TELEMETRY",
							ValidValues: []string{"0"},
							Description: "Disables Yarn telemetry",
						},
						Value:   "0",
						Passing: true,
					},
				},
			},
		}, true)
	})

	for _, want := range []string{
		"Yarn",
		"not installed",
		"YARN_ENABLE_TELEMETRY=0",
		"Disables Yarn telemetry",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("PrintReportResults output missing %q:\n%s", want, output)
		}
	}
}

func TestWriteReportResultsShowsDefaultOffAndUnverifiedCommand(t *testing.T) {
	var buf bytes.Buffer
	WriteReportResults(&buf, []check.ToolResult{
		{
			Tool: tools.Tool{
				Name:             "Google Cloud SDK",
				Category:         "Cloud",
				DocumentationURL: "https://docs.example.com/gcloud",
				Default:          tools.DefaultDisabled,
				Commands:         []string{"gcloud config set disable_usage_reporting true"},
			},
			Installed:  true,
			BinaryPath: "/usr/local/bin/gcloud",
			AllPassing: true,
			Checks: []check.EnvCheckResult{{
				Check:   tools.EnvCheck{Name: "CLOUDSDK_CORE_DISABLE_USAGE_REPORTING"},
				Passing: true,
				Reason:  check.ReasonDefaultOff,
			}},
		},
	}, true)

	output := buf.String()
	for _, want := range []string{
		"not set; telemetry is off unless you opt in",
		"→ also (not verified):",
		"gcloud config set disable_usage_reporting true",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("WriteReportResults output missing %q:\n%s", want, output)
		}
	}
}

func TestWriteCheckResultsDimsUnusedAnyOfAlternative(t *testing.T) {
	var buf bytes.Buffer
	WriteCheckResults(&buf, []check.ToolResult{
		{
			Tool: tools.Tool{
				Name:     "GitHub CLI",
				Category: "Cloud",
				OptOut:   tools.AnyOf,
			},
			Installed:  true,
			AllPassing: true,
			Checks: []check.EnvCheckResult{
				{Check: tools.EnvCheck{Name: "GH_TELEMETRY"}, Passing: false, Reason: check.ReasonUnset},
				{Check: tools.EnvCheck{Name: "DO_NOT_TRACK", ValidValues: []string{"1"}}, Value: "1", Passing: true, Reason: check.ReasonMatched},
			},
		},
	}, true)

	output := buf.String()
	if !strings.Contains(output, "alternative opt-out, not required") {
		t.Fatalf("compact AnyOf row missing unused-alternative hint:\n%s", output)
	}
	if !strings.Contains(output, "DO_NOT_TRACK=1") {
		t.Fatalf("compact AnyOf row missing the matching variable:\n%s", output)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}

	os.Stdout = w
	defer func() {
		os.Stdout = originalStdout
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error = %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("r.Close() error = %v", err)
	}

	return buf.String()
}
