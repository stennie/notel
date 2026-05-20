package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stennie/notel/internal/check"
)

func TestRootCommandSilencesErrorsAndUsage(t *testing.T) {
	if !rootCmd.SilenceErrors {
		t.Fatal("rootCmd.SilenceErrors = false, want true")
	}
	if !rootCmd.SilenceUsage {
		t.Fatal("rootCmd.SilenceUsage = false, want true")
	}
}

func TestLeafCommandsRejectUnexpectedArgs(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{name: "audit", cmd: auditCmd},
		{name: "list", cmd: listCmd},
		{name: "version", cmd: versionCmd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cmd.Args(tt.cmd, []string{"unexpected"}); err == nil {
				t.Fatal("Args accepted unexpected positional argument")
			}
		})
	}
}

func TestListCommandHasLsAlias(t *testing.T) {
	found := false
	for _, alias := range listCmd.Aliases {
		if alias == "ls" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("listCmd.Aliases = %#v, want to include %q", listCmd.Aliases, "ls")
	}
}

func TestRootHelpUsesStyledHeader(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	if err := rootCmd.Help(); err != nil {
		t.Fatalf("rootCmd.Help() error = %v", err)
	}

	output := buf.String()
	for _, want := range []string{
		"notel  —  DevTools Telemetry Auditor",
		"Usage",
		"Available Commands",
		"audit",
		"Copyright 2026 Stennie Steneker",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("styled help output missing %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "Author: Stennie Steneker") {
		t.Fatalf("styled help output still contains author line:\n%s", output)
	}
}

func TestVersionCommandIncludesMetadata(t *testing.T) {
	var buf bytes.Buffer
	versionCmd.SetOut(&buf)
	versionCmd.SetErr(&buf)

	versionCmd.Run(versionCmd, nil)

	output := buf.String()
	for _, want := range []string{
		"notel version",
		"Copyright 2026 Stennie Steneker",
		"License: Apache 2.0",
		"GitHub: https://github.com/stennie/notel",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("version output missing %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "Author: Stennie Steneker") {
		t.Fatalf("version output still contains author line:\n%s", output)
	}
}

func TestAuditAllFlagWorksWithAndWithoutVerbose(t *testing.T) {
	originalShowAll := auditShowAll
	originalVerbose := auditVerbose
	originalRunChecks := runChecks
	originalWriteTitle := writeTitle
	originalWriteCheckResults := writeCheckResults
	originalWriteDetailedResults := writeDetailedResults
	originalWriteAuditSummary := writeAuditSummary

	t.Cleanup(func() {
		auditShowAll = originalShowAll
		auditVerbose = originalVerbose
		runChecks = originalRunChecks
		writeTitle = originalWriteTitle
		writeCheckResults = originalWriteCheckResults
		writeDetailedResults = originalWriteDetailedResults
		writeAuditSummary = originalWriteAuditSummary
		_ = auditCmd.Flags().Set("all", "false")
		_ = auditCmd.Flags().Set("verbose", "false")
	})

	runChecks = func() []check.ToolResult {
		return []check.ToolResult{{Installed: true}}
	}
	writeTitle = func(io.Writer) {}
	writeAuditSummary = func(io.Writer, []check.ToolResult) {}

	tests := []struct {
		name             string
		args             map[string]string
		wantCompactAll   *bool
		wantDetailedAll  *bool
		wantCompactCalls int
		wantDetailCalls  int
	}{
		{
			name:             "all without verbose",
			args:             map[string]string{"all": "true", "verbose": "false"},
			wantCompactAll:   boolPtr(true),
			wantDetailedAll:  nil,
			wantCompactCalls: 1,
			wantDetailCalls:  0,
		},
		{
			name:             "all with verbose",
			args:             map[string]string{"all": "true", "verbose": "true"},
			wantCompactAll:   nil,
			wantDetailedAll:  boolPtr(true),
			wantCompactCalls: 0,
			wantDetailCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditShowAll = false
			auditVerbose = false

			for name, value := range tt.args {
				if err := auditCmd.Flags().Set(name, value); err != nil {
					t.Fatalf("setting %s flag: %v", name, err)
				}
			}

			compactCalls := 0
			detailCalls := 0

			writeCheckResults = func(_ io.Writer, _ []check.ToolResult, showAll bool) {
				compactCalls++
				if tt.wantCompactAll == nil {
					t.Fatalf("unexpected compact render call with showAll=%v", showAll)
				}
				if showAll != *tt.wantCompactAll {
					t.Fatalf("compact render showAll=%v, want %v", showAll, *tt.wantCompactAll)
				}
			}
			writeDetailedResults = func(_ io.Writer, _ []check.ToolResult, showAll bool) {
				detailCalls++
				if tt.wantDetailedAll == nil {
					t.Fatalf("unexpected detailed render call with showAll=%v", showAll)
				}
				if showAll != *tt.wantDetailedAll {
					t.Fatalf("detailed render showAll=%v, want %v", showAll, *tt.wantDetailedAll)
				}
			}

			if err := auditCmd.RunE(auditCmd, nil); err != nil {
				t.Fatalf("auditCmd.RunE() error = %v", err)
			}
			if compactCalls != tt.wantCompactCalls {
				t.Fatalf("compact render calls = %d, want %d", compactCalls, tt.wantCompactCalls)
			}
			if detailCalls != tt.wantDetailCalls {
				t.Fatalf("detailed render calls = %d, want %d", detailCalls, tt.wantDetailCalls)
			}
		})
	}
}

func boolPtr(v bool) *bool {
	return &v
}
