// Package cmd provides the notel CLI command tree via Cobra.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stennie/notel/internal/output"
)

// Version identifies the application build version. It defaults to "dev" and
// is intended to be overridden at build time with -ldflags.
var Version = "dev"

const (
	Copyright = "Copyright 2026 Stennie Steneker"
	License   = "License: Apache 2.0"
	GitHubURL = "GitHub: https://github.com/stennie/notel"
)

var rootCmd = &cobra.Command{
	Use:           "notel",
	Short:         "Audit telemetry opt-out settings for developer tools",
	Version:       Version,
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `notel checks whether common developer tools are installed and verifies
that their telemetry opt-out environment variables are properly set.

Run 'notel audit' to see the current status of all installed tools,
	or 'notel audit --verbose' for a detailed audit with fix hints.`,
}

// Execute is the entry point called from main.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		output.PrintHelp(cmd.OutOrStdout(), cmd, Copyright)
	})
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the application version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s version %s\n%s\n%s\n%s\n", cmd.Root().Name(), Version, Copyright, License, GitHubURL)
	},
}
