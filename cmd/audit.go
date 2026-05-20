package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/fix"
	"github.com/stennie/notel/internal/output"
)

var (
	auditShowAll bool
	auditVerbose bool
	auditFix     string

	runChecks            = check.Run
	writeTitle           = output.WriteTitle
	writeCheckResults    = output.WriteCheckResults
	writeDetailedResults = output.WriteReportResults
	writeAuditSummary    = output.WriteSummary
	printFixCommands     = output.PrintFixCommands
	detectShell          = fix.DetectShell
	parseShell           = fix.ParseShell
	buildFixSuggestions  = fix.Suggestions
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit telemetry settings for installed tools",
	Args:  cobra.NoArgs,
	Long: `Checks all installed developer tools and reports whether their
telemetry opt-out environment variables are correctly set.

By default only installed tools are shown. Use --all to include tools
that are not detected in PATH. Use --verbose for detailed per-tool output.
Use --fix[=<shell>] to print shell commands that disable telemetry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		results := runChecks()
		stdout := cmd.OutOrStdout()
		stderr := cmd.ErrOrStderr()
		if auditFix != "" {
			shell, err := resolveFixShell(auditFix)
			if err != nil {
				return err
			}
			writeTitle(stderr)
			if auditShowAll || auditVerbose {
				writeAuditView(stderr, results)
			}
			printFixCommands(stdout, stderr, shell, buildFixSuggestions(results, auditShowAll))
			return nil
		}

		writeTitle(stdout)
		writeAuditView(stdout, results)
		return nil
	},
}

func init() {
	auditCmd.Flags().BoolVarP(&auditShowAll, "all", "a", false, "Show all tools, including those not installed")
	auditCmd.Flags().BoolVarP(&auditVerbose, "verbose", "v", false, "Show detailed per-tool output")
	auditCmd.Flags().StringVarP(&auditFix, "fix", "f", "", "Print shell commands to disable telemetry; optionally pass bash, fish, powershell, or zsh")
	auditCmd.Flags().Lookup("fix").NoOptDefVal = "auto"
}

func resolveFixShell(value string) (fix.Shell, error) {
	if value == "" || value == "auto" {
		shell, err := detectShell()
		if err != nil {
			return "", fmt.Errorf("%w; use --fix=<shell> to choose one explicitly", err)
		}
		return shell, nil
	}
	return parseShell(value)
}

func writeAuditView(w io.Writer, results []check.ToolResult) {
	if auditVerbose {
		writeDetailedResults(w, results, auditShowAll)
	} else {
		writeCheckResults(w, results, auditShowAll)
	}
	writeAuditSummary(w, results)
}
