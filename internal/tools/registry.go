// Package tools defines the registry of developer tools that notel can audit.
package tools

// EnvCheck defines an environment variable requirement for a tool.
type EnvCheck struct {
	Name        string   // environment variable name
	ValidValues []string // accepted values that mean "opt-out" (empty = any non-empty value)
	Description string   // human-readable description of what this controls
}

// Tool represents a developer tool that can be audited for telemetry settings.
type Tool struct {
	Name             string
	Description      string
	DocumentationURL string
	DataCollection   string
	Binary           string // binary name to detect via PATH (empty = always consider installed)
	Category         string
	EnvChecks        []EnvCheck
}

// Registry returns the full list of developer tools notel can audit.
// To add a new tool, append a new Tool entry here.
func Registry() []Tool {
	return []Tool{
		// ── All The Things ───────────────────────────────────────────────────────────
		{
			Name:             "Do Not Track",
			Description:      "Cross-tool environment variable convention for disabling telemetry and tracking",
			DocumentationURL: "https://donottrack.sh",
			DataCollection:   "Signals a global preference to disable telemetry, analytics, crash reporting, and non-essential tracking.",
			Binary:           "",
			Category:         "All The Things",
			EnvChecks: []EnvCheck{
				{
					Name:        "DO_NOT_TRACK",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Signals a global do-not-track preference to tools that honor it",
				},
			},
		},

		// ── Cloud & Deployment ──────────────────────────────────────────────────────
		{
			Name:             "Algolia CLI",
			Description:      "Algolia command-line interface for search platform management",
			DocumentationURL: "https://www.algolia.com/doc/tools/cli/telemetry",
			DataCollection:   "Command usage, operating system details, CLI version, and local profile metadata.",
			Binary:           "algolia",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "ALGOLIA_CLI_TELEMETRY",
					ValidValues: []string{"0"},
					Description: "Disables Algolia CLI telemetry collection",
				},
			},
		},
		{
			Name:             "Azure CLI",
			Description:      "Microsoft Azure command-line interface",
			DocumentationURL: "https://learn.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest#cli-configuration-values-and-environment-variables",
			DataCollection:   "Command usage, performance metrics, and error-rate telemetry.",
			Binary:           "az",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "AZURE_CORE_COLLECT_TELEMETRY",
					ValidValues: []string{"0", "false", "False", "FALSE"},
					Description: "Disables Azure CLI telemetry collection",
				},
			},
		},
		{
			Name:             "CockroachDB",
			Description:      "Distributed SQL database and CockroachDB command-line tools",
			DocumentationURL: "https://www.cockroachlabs.com/docs/stable/telemetry",
			DataCollection:   "Cluster diagnostics, telemetry, and crash reports sent to Cockroach Labs.",
			Binary:           "cockroach",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "COCKROACH_SKIP_ENABLING_DIAGNOSTIC_REPORTING",
					ValidValues: []string{"true", "True", "TRUE", "1"},
					Description: "Disables CockroachDB diagnostic reporting for new clusters",
				},
			},
		},
		{
			Name:             "Google Cloud SDK",
			Description:      "Google Cloud command-line tools",
			DocumentationURL: "https://cloud.google.com/sdk/docs/usage-statistics",
			DataCollection:   "Anonymized command execution metrics, timing, and error status.",
			Binary:           "gcloud",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "CLOUDSDK_CORE_DISABLE_USAGE_REPORTING",
					ValidValues: []string{"true", "True", "TRUE", "1"},
					Description: "Disables Google Cloud SDK usage reporting",
				},
			},
		},

		{
			Name:             "Netlify CLI",
			Description:      "Netlify platform CLI",
			DocumentationURL: "https://docs.netlify.com/api-and-cli-guides/cli-guides/get-started-with-cli/#usage-data-collection",
			DataCollection:   "Anonymous CLI usage and diagnostic telemetry.",
			Binary:           "netlify",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "NETLIFY_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Netlify CLI telemetry",
				},
			},
		},
		{
			Name:             "Turborepo",
			Description:      "High-performance build system for JavaScript/TypeScript",
			DocumentationURL: "https://turborepo.com/docs/telemetry",
			DataCollection:   "Anonymous command usage, host information, and repo or task metrics.",
			Binary:           "turbo",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "TURBO_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Turborepo telemetry",
				},
			},
		},

		// ── Framework ───────────────────────────────────────────────────────────────
		{
			Name:             "Angular CLI",
			Description:      "Angular framework CLI",
			DocumentationURL: "https://angular.dev/cli/analytics",
			DataCollection:   "Command usage, selected flags, workspace shape, and local version metadata.",
			Binary:           "ng",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "NG_CLI_ANALYTICS",
					ValidValues: []string{"false", "ci", "0", "False", "FALSE"},
					Description: "Disables Angular CLI analytics (set to false or ci)",
				},
			},
		},
		{
			Name:             "Astro",
			Description:      "The web framework for content-driven websites",
			DocumentationURL: "https://docs.astro.build/en/reference/cli-reference/#astro-telemetry",
			DataCollection:   "Anonymous command usage, integration usage, and project metadata.",
			Binary:           "astro",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "ASTRO_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Astro telemetry",
				},
			},
		},
		{
			Name:             "Expo CLI",
			Description:      "Expo command-line tools for React Native and universal apps",
			DocumentationURL: "https://docs.expo.dev/more/expo-cli/#telemetry",
			DataCollection:   "Anonymous CLI usage and diagnostics for Expo development workflows.",
			Binary:           "expo",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "EXPO_NO_TELEMETRY",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Expo CLI telemetry",
				},
			},
		},
		{
			Name:             "Gatsby",
			Description:      "React-based static site generator",
			DocumentationURL: "https://www.gatsbyjs.com/docs/telemetry/",
			DataCollection:   "Anonymous command usage, plugin usage, and machine characteristics.",
			Binary:           "gatsby",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "GATSBY_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Gatsby telemetry",
				},
			},
		},
		{
			Name:             "Next.js",
			Description:      "The React framework for production",
			DocumentationURL: "https://nextjs.org/telemetry",
			DataCollection:   "Anonymous command usage, session timing, and project or machine characteristics.",
			Binary:           "next",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "NEXT_TELEMETRY_DISABLED",
					ValidValues: []string{"1"},
					Description: "Disables Next.js telemetry collection",
				},
			},
		},
		{
			Name:             "Nuxt",
			Description:      "The intuitive Vue framework",
			DocumentationURL: "https://github.com/nuxt/telemetry#nuxt-telemetry-module",
			DataCollection:   "Anonymous command usage, module usage, and environment characteristics.",
			Binary:           "nuxt",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "NUXT_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Nuxt telemetry collection",
				},
			},
		},
		{
			Name:             "Storybook",
			Description:      "Frontend workshop for building UI components and pages in isolation",
			DocumentationURL: "https://storybook.js.org/docs/configure/telemetry",
			DataCollection:   "Anonymous framework usage, addon data, and environment metadata.",
			Binary:           "storybook",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				{
					Name:        "STORYBOOK_DISABLE_TELEMETRY",
					ValidValues: []string{"true", "True", "TRUE"},
					Description: "Disables Storybook telemetry",
				},
			},
		},

		// ── Package Manager ─────────────────────────────────────────────────────────
		{
			Name:             "Homebrew",
			Description:      "The missing package manager for macOS",
			DocumentationURL: "https://docs.brew.sh/Analytics",
			DataCollection:   "Anonymous install events, command usage, and build error metadata.",
			Binary:           "brew",
			Category:         "Package Manager",
			EnvChecks: []EnvCheck{
				{
					Name:        "HOMEBREW_NO_ANALYTICS",
					ValidValues: []string{"1"},
					Description: "Disables Homebrew analytics reporting",
				},
			},
		},
		{
			Name:             "Poetry",
			Description:      "Python dependency management and packaging tool",
			DocumentationURL: "https://python-poetry.org/docs/configuration/#using-environment-variables",
			DataCollection:   "Anonymous usage statistics and telemetry about Poetry command execution.",
			Binary:           "poetry",
			Category:         "Package Manager",
			EnvChecks: []EnvCheck{
				{
					Name:        "POETRY_TELEMETRY_ENABLED",
					ValidValues: []string{"0", "false", "False", "FALSE"},
					Description: "Disables Poetry telemetry collection",
				},
			},
		},
		{
			Name:             "Yarn",
			Description:      "Fast, reliable JavaScript package manager",
			DocumentationURL: "https://yarnpkg.com/advanced/telemetry",
			DataCollection:   "Anonymous command usage and version information.",
			Binary:           "yarn",
			Category:         "Package Manager",
			EnvChecks: []EnvCheck{
				{
					Name:        "YARN_ENABLE_TELEMETRY",
					ValidValues: []string{"0", "false", "False", "FALSE"},
					Description: "Disables Yarn telemetry (Yarn 2+)",
				},
			},
		},

		// ── Runtimes & SDKs ─────────────────────────────────────────────────────────
		{
			Name:             ".NET SDK",
			Description:      "Microsoft .NET SDK and CLI",
			DocumentationURL: "https://learn.microsoft.com/en-us/dotnet/core/tools/telemetry",
			DataCollection:   "CLI command usage, SDK versions, and exception metadata.",
			Binary:           "dotnet",
			Category:         "Runtimes & SDKs",
			EnvChecks: []EnvCheck{
				{
					Name:        "DOTNET_CLI_TELEMETRY_OPTOUT",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables .NET CLI telemetry collection",
				},
			},
		},
		{
			Name:             "Bun",
			Description:      "Fast all-in-one JavaScript runtime and toolkit",
			DocumentationURL: "https://bun.sh/docs/runtime/bunfig#telemetry",
			DataCollection:   "Anonymous telemetry and crash-reporting signals from Bun tooling.",
			Binary:           "bun",
			Category:         "Runtimes & SDKs",
			EnvChecks: []EnvCheck{
				{
					Name:        "DO_NOT_TRACK",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Bun telemetry via the global do-not-track convention",
				},
			},
		},
		{
			Name:             "Flutter",
			Description:      "Google's UI toolkit for cross-platform apps",
			DocumentationURL: "https://docs.flutter.dev/reference/crash-reporting",
			DataCollection:   "CLI crash reports, tool usage signals, and local environment metadata.",
			Binary:           "flutter",
			Category:         "Runtimes & SDKs",
			EnvChecks: []EnvCheck{
				{
					Name:        "FLUTTER_CLI_CRASH_REPORTING",
					ValidValues: []string{"false", "False", "FALSE", "0"},
					Description: "Disables Flutter CLI crash reporting",
				},
			},
		},
		{
			Name:             "Node.js",
			Description:      "JavaScript runtime built on Chrome's V8 engine",
			DocumentationURL: "https://nodejs.org/api/cli.html#node_no_telemetry1",
			DataCollection:   "Anonymous runtime and CLI telemetry in Node features that support it.",
			Binary:           "node",
			Category:         "Runtimes & SDKs",
			EnvChecks: []EnvCheck{
				{
					Name:        "DISABLE_TELEMETRY",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables telemetry for tools that respect this variable",
				},
				{
					Name:        "NODE_NO_TELEMETRY",
					ValidValues: []string{"1"},
					Description: "Disables Node.js built-in telemetry (Node 21+)",
				},
			},
		},

		// ── Security & Observability ───────────────────────────────────────────────
		{
			Name:             "Semgrep",
			Description:      "Static analysis and code security scanning CLI",
			DocumentationURL: "https://semgrep.dev/docs/metrics",
			DataCollection:   "Usage metrics about Semgrep runs, registry usage, and login-related activity.",
			Binary:           "semgrep",
			Category:         "Security & Observability",
			EnvChecks: []EnvCheck{
				{
					Name:        "SEMGREP_SEND_METRICS",
					ValidValues: []string{"off"},
					Description: "Disables Semgrep metrics reporting",
				},
			},
		},
		{
			Name:             "Sentry CLI",
			Description:      "Sentry error tracking command-line tool",
			DocumentationURL: "https://cli.sentry.dev/configuration/",
			DataCollection:   "Anonymous CLI usage and diagnostic telemetry.",
			Binary:           "sentry-cli",
			Category:         "Security & Observability",
			EnvChecks: []EnvCheck{
				{
					Name:        "SENTRY_CLI_NO_TELEMETRY",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Sentry CLI telemetry",
				},
			},
		},
	}
}
