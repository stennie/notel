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
		// ── Build ───────────────────────────────────────────────────────────────────
		{
			Name:             "Turborepo",
			Description:      "High-performance build system for JavaScript/TypeScript",
			DocumentationURL: "https://turborepo.com/docs/telemetry",
			DataCollection:   "Anonymous command usage, host information, and repo or task metrics.",
			Binary:           "turbo",
			Category:         "Build",
			EnvChecks: []EnvCheck{
				{
					Name:        "TURBO_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Turborepo telemetry",
				},
			},
		},

		// ── Cloud ───────────────────────────────────────────────────────────────────
		{
			Name:             "Azure CLI",
			Description:      "Microsoft Azure command-line interface",
			DocumentationURL: "https://learn.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest",
			DataCollection:   "Command usage, performance metrics, and error-rate telemetry.",
			Binary:           "az",
			Category:         "Cloud",
			EnvChecks: []EnvCheck{
				{
					Name:        "AZURE_CORE_COLLECT_TELEMETRY",
					ValidValues: []string{"0", "false", "False", "FALSE"},
					Description: "Disables Azure CLI telemetry collection",
				},
			},
		},
		{
			Name:             "Google Cloud SDK",
			Description:      "Google Cloud command-line tools",
			DocumentationURL: "https://cloud.google.com/sdk/docs/usage-statistics",
			DataCollection:   "Anonymized command execution metrics, timing, and error status.",
			Binary:           "gcloud",
			Category:         "Cloud",
			EnvChecks: []EnvCheck{
				{
					Name:        "CLOUDSDK_CORE_DISABLE_USAGE_REPORTING",
					ValidValues: []string{"true", "True", "TRUE", "1"},
					Description: "Disables Google Cloud SDK usage reporting",
				},
			},
		},

		// ── Deployment ──────────────────────────────────────────────────────────────
		{
			Name:             "Netlify CLI",
			Description:      "Netlify platform CLI",
			DocumentationURL: "https://docs.netlify.com/cli/get-started/#usage-data",
			DataCollection:   "Anonymous CLI usage and diagnostic telemetry.",
			Binary:           "netlify",
			Category:         "Deployment",
			EnvChecks: []EnvCheck{
				{
					Name:        "NETLIFY_TELEMETRY_DISABLED",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Netlify CLI telemetry",
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
			DocumentationURL: "https://github.com/nuxt/telemetry",
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

		// ── Infrastructure ──────────────────────────────────────────────────────────
		{
			Name:             "Terraform",
			Description:      "Infrastructure as Code tool by HashiCorp",
			DocumentationURL: "https://developer.hashicorp.com/terraform/cli/config/environment-variables#checkpoint_disable",
			DataCollection:   "Checkpoint version checks and security bulletin or alert lookups.",
			Binary:           "terraform",
			Category:         "Infrastructure",
			EnvChecks: []EnvCheck{
				{
					Name:        "CHECKPOINT_DISABLE",
					ValidValues: []string{"1"},
					Description: "Disables HashiCorp Checkpoint (telemetry and version checks)",
				},
			},
		},

		// ── Monitoring ──────────────────────────────────────────────────────────────
		{
			Name:             "Sentry CLI",
			Description:      "Sentry error tracking command-line tool",
			DocumentationURL: "https://cli.sentry.dev/configuration/",
			DataCollection:   "Anonymous CLI usage and diagnostic telemetry.",
			Binary:           "sentry-cli",
			Category:         "Monitoring",
			EnvChecks: []EnvCheck{
				{
					Name:        "SENTRY_CLI_NO_TELEMETRY",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables Sentry CLI telemetry",
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

		// ── Runtime ─────────────────────────────────────────────────────────────────
		{
			Name:             "Node.js",
			Description:      "JavaScript runtime built on Chrome's V8 engine",
			DocumentationURL: "https://nodejs.org/api/cli.html#node_no_telemetry1",
			DataCollection:   "Anonymous runtime and CLI telemetry in Node features that support it.",
			Binary:           "node",
			Category:         "Runtime",
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

		// ── SDK ─────────────────────────────────────────────────────────────────────
		{
			Name:             ".NET SDK",
			Description:      "Microsoft .NET SDK and CLI",
			DocumentationURL: "https://learn.microsoft.com/en-us/dotnet/core/tools/telemetry",
			DataCollection:   "CLI command usage, SDK versions, and exception metadata.",
			Binary:           "dotnet",
			Category:         "SDK",
			EnvChecks: []EnvCheck{
				{
					Name:        "DOTNET_CLI_TELEMETRY_OPTOUT",
					ValidValues: []string{"1", "true", "True", "TRUE"},
					Description: "Disables .NET CLI telemetry collection",
				},
			},
		},
		{
			Name:             "Flutter",
			Description:      "Google's UI toolkit for cross-platform apps",
			DocumentationURL: "https://docs.flutter.dev/reference/crash-reporting",
			DataCollection:   "CLI crash reports, tool usage signals, and local environment metadata.",
			Binary:           "flutter",
			Category:         "SDK",
			EnvChecks: []EnvCheck{
				{
					Name:        "FLUTTER_CLI_CRASH_REPORTING",
					ValidValues: []string{"false", "False", "FALSE", "0"},
					Description: "Disables Flutter CLI crash reporting",
				},
			},
		},
	}
}
