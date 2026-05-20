package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/stennie/notel/internal/check"
	"github.com/stennie/notel/internal/output"
	"github.com/stennie/notel/internal/tools"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all supported tools and their telemetry variables",
	Args:    cobra.NoArgs,
	Long: `Lists every developer tool that notel can audit, along with the
environment variables and accepted values that disable telemetry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stdout := cmd.OutOrStdout()
		output.WriteTitle(stdout)

		subtitleStyle := lipgloss.NewStyle().Bold(true).MarginBottom(0)
		fmt.Fprintln(stdout, subtitleStyle.Render("  Supported Tools"))
		fmt.Fprintln(stdout)

		// Build results so PrintListResults can show installation status.
		results := check.RunFor(tools.Registry())
		output.WriteListResults(stdout, results)
		return nil
	},
}
