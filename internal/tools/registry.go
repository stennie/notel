package tools

// Registry returns the developer tools notel can audit.
//
// To add a tool, append one Tool in category order, then name order.
// The zero values of Default and OptOut mean "telemetry is on until every
// variable matches". Set those fields only when that sentence is false.
// CaseInsensitive is for vendors that accept any capitalization. An exact
// EnvCheck literal is for vendors that document one spelling.
//
// Commands are printed as hints. notel does not run them and does not treat
// them as proof that telemetry is off.
func Registry() []Tool {
	return []Tool{
		// ── All The Things ───────────────────────────────────────────────────────────
		{
			Name:             "Do Not Track",
			Description:      "Cross-tool environment variable convention for disabling telemetry and tracking",
			DocumentationURL: "https://donottrack.sh",
			DataCollection:   "Signals a global preference to disable telemetry, analytics, crash reporting, and non-essential tracking.",
			// No binary. This entry is the convention itself, so it is always in scope.
			Category: "All The Things",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"DO_NOT_TRACK",
					"Signals a global do-not-track preference to tools that honor it",
					"1", "true",
				),
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
				CaseInsensitive(
					"AZURE_CORE_COLLECT_TELEMETRY",
					"Disables Azure CLI telemetry collection",
					"0", "false",
				),
			},
			Commands: []string{"az config set core.collect_telemetry=no"},
		},
		{
			Name:             "Cloudflare Wrangler CLI",
			Description:      "Cloudflare CLI for building and deploying Workers and Pages applications",
			DocumentationURL: "https://github.com/cloudflare/workers-sdk/blob/main/packages/wrangler/telemetry.md#wrangler-cli-telemetry",
			DataCollection:   "Command names, anonymized arguments and flags, anonymized project metadata, API timing, CLI and runtime details, session duration, CI or interactive usage, AI agent usage, sanitized error details, and general machine information.",
			Binary:           "wrangler",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"WRANGLER_SEND_METRICS",
					"Disables Cloudflare Wrangler CLI telemetry",
					"false",
				),
			},
			Commands: []string{"wrangler telemetry disable"},
		},
		{
			Name:             "CockroachDB",
			Description:      "Distributed SQL database and CockroachDB command-line tools",
			DocumentationURL: "https://www.cockroachlabs.com/docs/stable/telemetry",
			DataCollection:   "Cluster diagnostics, telemetry, and crash reports sent to Cockroach Labs.",
			Binary:           "cockroach",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"COCKROACH_SKIP_ENABLING_DIAGNOSTIC_REPORTING",
					"Disables CockroachDB diagnostic reporting when set before the first node of a cluster is initialised",
					"true", "1",
				),
			},
		},
		{
			Name:             "Docker Agent",
			Description:      "Docker CLI agent for running AI agents from Docker",
			DocumentationURL: "https://docs.docker.com/ai/docker-agent/community/telemetry/",
			DataCollection:   "Command names, agent and model types, tool-call outcomes, token counts, estimated costs, and session metadata.",
			Binary:           "docker",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"TELEMETRY_ENABLED",
					"Disables Docker Agent telemetry",
					"false",
				),
			},
		},
		{
			Name:             "GitHub CLI",
			Description:      "Official command-line tool for GitHub",
			DocumentationURL: "https://docs.github.com/en/github-cli/github-cli/github-cli-telemetry",
			DataCollection:   "Command usage, feature usage, and diagnostic telemetry from the GitHub CLI.",
			Binary:           "gh",
			Category:         "Cloud & Deployment",
			// Either variable is a complete opt-out. Requiring both reported a
			// failure on machines that had already set DO_NOT_TRACK=1.
			OptOut: AnyOf,
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"GH_TELEMETRY",
					"Disables GitHub CLI telemetry",
					"false",
				),
				CaseInsensitive(
					"DO_NOT_TRACK",
					"Disables GitHub CLI telemetry via the global do-not-track convention",
					"1", "true",
				),
			},
			Commands: []string{"gh config set telemetry disabled"},
		},
		{
			Name:             "Google Cloud SDK",
			Description:      "Google Cloud command-line tools",
			DocumentationURL: "https://cloud.google.com/sdk/docs/usage-statistics",
			DataCollection:   "Anonymized command execution metrics, timing, and error status.",
			Binary:           "gcloud",
			Category:         "Cloud & Deployment",
			// Usage statistics are collected only after an explicit opt-in at
			// install time. An unset variable is already the disabled state.
			// A value that is set but is not a recognized opt-out still fails.
			Default: DefaultDisabled,
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"CLOUDSDK_CORE_DISABLE_USAGE_REPORTING",
					"Disables Google Cloud SDK usage reporting",
					"true", "1",
				),
			},
			Commands: []string{"gcloud config set disable_usage_reporting true"},
		},
		{
			Name:             "Hasura CLI",
			Description:      "Hasura command-line interface for managing GraphQL Engine projects",
			DocumentationURL: "https://hasura.io/docs/2.0/policies/telemetry/#cli",
			DataCollection:   "Command name, timestamp, error status, server version and UUID, operating system platform and architecture, and CLI version.",
			Binary:           "hasura",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"HASURA_GRAPHQL_ENABLE_TELEMETRY",
					"Disables Hasura CLI telemetry",
					"false",
				),
			},
		},
		{
			Name:             "Hugging Face CLI",
			Description:      "Hugging Face command-line interface for interacting with the Hub",
			DocumentationURL: "https://huggingface.co/docs/huggingface_hub/package_reference/environment_variables#hfhubdisabletelemetry",
			DataCollection:   "Usage data collected by Hugging Face libraries to monitor usage, debug issues, and help prioritize features.",
			Binary:           "hf",
			Category:         "Cloud & Deployment",
			// The Hub client honors either the global convention or its own variable.
			OptOut: AnyOf,
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"DO_NOT_TRACK",
					"Disables Hugging Face telemetry via the global do-not-track convention",
					"1", "true",
				),
				CaseInsensitive(
					"HF_HUB_DISABLE_TELEMETRY",
					"Disables Hugging Face telemetry",
					"1", "true",
				),
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
				CaseInsensitive(
					"NETLIFY_TELEMETRY_DISABLED",
					"Disables Netlify CLI telemetry",
					"1", "true",
				),
			},
			Commands: []string{"netlify --telemetry-disable"},
		},
		{
			Name:             "Railway CLI",
			Description:      "Railway command-line interface for deploying and managing projects",
			DocumentationURL: "https://docs.railway.com/cli/telemetry",
			DataCollection:   "Command and subcommand names, duration, success, truncated error message, OS, architecture, CLI version, and whether the command ran in CI.",
			Binary:           "railway",
			Category:         "Cloud & Deployment",
			OptOut:           AnyOf,
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"DO_NOT_TRACK",
					"Disables Railway CLI telemetry via the global do-not-track convention",
					"1", "true",
				),
				CaseInsensitive(
					"RAILWAY_NO_TELEMETRY",
					"Disables Railway CLI telemetry",
					"1", "true",
				),
			},
		},
		{
			Name:             "Redocly CLI",
			Description:      "Redocly command-line interface for API governance and documentation workflows",
			DocumentationURL: "https://redocly.com/docs/realm/reunite/project/telemetry",
			DataCollection:   "Anonymous usage information about Redocly project feature usage, with additional filtering to avoid sending secrets and other sensitive data.",
			Binary:           "redocly",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"REDOCLY_TELEMETRY",
					"Disables Redocly CLI telemetry",
					"off",
				),
			},
		},
		{
			Name:             "Stripe CLI",
			Description:      "Stripe command-line interface for developer workflows and API interactions",
			DocumentationURL: "https://docs.stripe.com/cli/telemetry",
			DataCollection:   "Command usage patterns, error rates and types, performance metrics, operating system details, CLI version information, and detected AI coding agent names.",
			Binary:           "stripe",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"STRIPE_CLI_TELEMETRY_OPTOUT",
					"Disables Stripe CLI telemetry",
					"1", "true",
				),
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
				CaseInsensitive(
					"TURBO_TELEMETRY_DISABLED",
					"Disables Turborepo telemetry",
					"1", "true",
				),
			},
			Commands: []string{"turbo telemetry disable"},
		},
		{
			Name:             "Vercel CLI",
			Description:      "Frontend cloud platform CLI for deployment and project management",
			DocumentationURL: "https://vercel.com/docs/cli/about-telemetry",
			DataCollection:   "Command usage, arguments, CLI version, and general machine information.",
			Binary:           "vercel",
			Category:         "Cloud & Deployment",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"VERCEL_TELEMETRY_DISABLED",
					"Disables Vercel CLI telemetry",
					"1", "true",
				),
			},
			Commands: []string{"vercel telemetry disable"},
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
				CaseInsensitive(
					"NG_CLI_ANALYTICS",
					"Disables Angular CLI analytics (set to false or ci)",
					"false", "ci", "0",
				),
			},
			Commands: []string{"ng analytics disable"},
		},
		{
			Name:             "Astro",
			Description:      "The web framework for content-driven websites",
			DocumentationURL: "https://docs.astro.build/en/reference/cli-reference/#astro-telemetry",
			DataCollection:   "Anonymous command usage, integration usage, and project metadata.",
			Binary:           "astro",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"ASTRO_TELEMETRY_DISABLED",
					"Disables Astro telemetry",
					"1", "true",
				),
			},
			Commands: []string{"astro telemetry disable"},
		},
		{
			Name:             "Expo CLI",
			Description:      "Expo command-line tools for React Native and universal apps",
			DocumentationURL: "https://docs.expo.dev/more/expo-cli/#telemetry",
			DataCollection:   "Anonymous CLI usage and diagnostics for Expo development workflows.",
			Binary:           "expo",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"EXPO_NO_TELEMETRY",
					"Disables Expo CLI telemetry",
					"1", "true",
				),
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
				CaseInsensitive(
					"GATSBY_TELEMETRY_DISABLED",
					"Disables Gatsby telemetry",
					"1", "true",
				),
			},
			Commands: []string{"gatsby telemetry --disable"},
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
					// Next.js documents the exact value "1". Other spellings
					// are not accepted, so this check stays case-sensitive
					// and does not list "true".
					Name:        "NEXT_TELEMETRY_DISABLED",
					ValidValues: []string{"1"},
					Description: "Disables Next.js telemetry collection",
				},
			},
			Commands: []string{"next telemetry disable"},
		},
		{
			Name:             "Nuxt",
			Description:      "The intuitive Vue framework",
			DocumentationURL: "https://github.com/nuxt/telemetry#nuxt-telemetry-module",
			DataCollection:   "Anonymous command usage, module usage, and environment characteristics.",
			Binary:           "nuxt",
			Category:         "Framework",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"NUXT_TELEMETRY_DISABLED",
					"Disables Nuxt telemetry collection",
					"1", "true",
				),
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
				CaseInsensitive(
					"STORYBOOK_DISABLE_TELEMETRY",
					"Disables Storybook telemetry",
					"true",
				),
			},
			Commands: []string{"npm run storybook -- --disable-telemetry"},
		},

		// ── Package Manager ─────────────────────────────────────────────────────────
		{
			Name:             "CocoaPods",
			Description:      "Dependency manager for Swift and Objective-C projects",
			DocumentationURL: "https://blog.cocoapods.org/Stats/",
			DataCollection:   "Anonymous pod install statistics, including hashed project target identifiers, target product types, pod names, and CocoaPods version information.",
			Binary:           "pod",
			Category:         "Package Manager",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"COCOAPODS_DISABLE_STATS",
					"Disables CocoaPods stats collection",
					"true", "1",
				),
			},
		},
		{
			Name:             "Homebrew",
			Description:      "The missing package manager for macOS",
			DocumentationURL: "https://docs.brew.sh/Analytics",
			DataCollection:   "Anonymous install events, command usage, and build error metadata.",
			Binary:           "brew",
			Category:         "Package Manager",
			EnvChecks: []EnvCheck{
				{
					// Homebrew documents the exact value "1".
					Name:        "HOMEBREW_NO_ANALYTICS",
					ValidValues: []string{"1"},
					Description: "Disables Homebrew analytics reporting",
				},
			},
			Commands: []string{"brew analytics off"},
		},
		{
			Name:             "Poetry",
			Description:      "Python dependency management and packaging tool",
			DocumentationURL: "https://python-poetry.org/docs/configuration/#using-environment-variables",
			DataCollection:   "Anonymous usage statistics and telemetry about Poetry command execution.",
			Binary:           "poetry",
			Category:         "Package Manager",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"POETRY_TELEMETRY_ENABLED",
					"Disables Poetry telemetry collection",
					"0", "false",
				),
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
				CaseInsensitive(
					"YARN_ENABLE_TELEMETRY",
					"Disables Yarn telemetry (Yarn 2+)",
					"0", "false",
				),
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
				CaseInsensitive(
					"DOTNET_CLI_TELEMETRY_OPTOUT",
					"Disables .NET CLI telemetry collection",
					"1", "true",
				),
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
				CaseInsensitive(
					"DO_NOT_TRACK",
					"Disables Bun telemetry via the global do-not-track convention",
					"1", "true",
				),
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
				CaseInsensitive(
					"FLUTTER_CLI_CRASH_REPORTING",
					"Disables Flutter CLI crash reporting. This does not disable general analytics.",
					"false", "0",
				),
			},
			// The analytics setting has no environment variable. The command is
			// a hint only. A passing crash-reporting check does not mean
			// `flutter config --no-analytics` has been run.
			Commands: []string{"flutter config --no-analytics"},
		},
		{
			Name:             "Node.js",
			Description:      "JavaScript runtime built on Chrome's V8 engine",
			DocumentationURL: "https://nodejs.org/api/cli.html#node_no_telemetry1",
			DataCollection:   "Anonymous runtime and CLI telemetry in Node features that support it.",
			Binary:           "node",
			Category:         "Runtimes & SDKs",
			// Both checks are required. NODE_NO_TELEMETRY is Node's own switch.
			// DISABLE_TELEMETRY is a separate convention used by libraries that
			// ship with Node programs. They are not alternatives.
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"DISABLE_TELEMETRY",
					"Disables telemetry for tools that respect this variable",
					"1", "true",
				),
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
					// Semgrep documents the exact value "off".
					Name:        "SEMGREP_SEND_METRICS",
					ValidValues: []string{"off"},
					Description: "Disables Semgrep metrics reporting",
				},
			},
			Commands: []string{"semgrep --metrics=off"},
		},
		{
			Name:             "Sentry CLI",
			Description:      "Sentry error tracking command-line tool",
			DocumentationURL: "https://cli.sentry.dev/configuration/",
			DataCollection:   "Anonymous CLI usage and diagnostic telemetry.",
			Binary:           "sentry-cli",
			Category:         "Security & Observability",
			EnvChecks: []EnvCheck{
				CaseInsensitive(
					"SENTRY_CLI_NO_TELEMETRY",
					"Disables Sentry CLI telemetry",
					"1", "true",
				),
			},
		},
	}
}
