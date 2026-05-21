# notel

**notel** audits telemetry opt-out settings across common developer tools.

Telemetry collection is a privacy concern for many developers, especially where collection
of data may not be obvious without consulting the documentation. Historically every tool has tended to choose unique environment variables or configuration settings which makes it difficult to opt-out by default.

This tool helps identify and opt-out of telemetry for popular developer tools. It checks
whether tools are installed and whether the environment variables that disable telemetry
collection are correctly set.

### Why disable telemetry?

- Disabling telemetry exercises your rights to privacy and control over your data and development environments.
- Poorly considered telemetry collection can potentially impact performance or resources.
- Telemetry collection may leak unexpected details (internal project names, dependencies, or build configurations) from private build or CI environments.
- Vendors aren't always up front about telemetry collection practices and details.

### Why enable telemetry?

- Telemetry signals can provide insight into tool adoption, usage patterns, and crash information that will help vendors improve the product experience.
- Aggregated telemetry signals can help product users understand usage of packages, features, and versions.
- Well-implemented telemetry can have negligible performance impact and can be configured to exclude CI environments.

### DO_NOT_TRACK

There is a proposed [`DO_NOT_TRACK`](https://donottrack.sh/) standard which aims to unambiguously express a user's intent  to opt out of:

 - Ad tracking
 - Usage reporting (anonymous or otherwise)
 - Telemetry
 - Crash reporting
 - Non-essential-to-functionality requests

However, this hasn't been widely adopted yet.

## Requirements

- Go 1.26.3+
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [Just](https://github.com/casey/just) — Task runner

## Build

- Local build: `just build` -> `bin/notel`
- Cross-platform release artifacts: `just release` -> `dist/`
- Full release preflight plus artifacts: `just release-check`
- Release packaging is shared by local and CI builds via `scripts/release.sh`

`just release` builds archives for:

- `darwin/amd64`
- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`
- `windows/amd64`
- `windows/arm64`

Release artifacts follow the convention `notel_<version>_<os>_<arch>.tar.gz` or `.zip`, and `dist/SHA256SUMS` is generated for the packaged archives.

## Usage

```
notel [command]
```

### Commands

| Command | Description |
|---------|-------------|
| `notel audit` | Show telemetry status for all **installed** tools |
| `notel audit --all` | Show status for all tools (including those not installed) |
| `notel audit --verbose` | Detailed audit with fix hints for all **installed** tools |
| `notel audit --verbose --all` | Detailed audit for all tools (including those not installed) |
| `notel audit --fix` / `notel audit -f` | Print shell commands to disable telemetry for the current detected shell |
| `notel audit --fix=zsh` / `notel audit -f=zsh` | Print shell commands for a specific shell (`bash`, `fish`, `powershell`, `zsh`) |
| `notel list` / `notel ls` | List all supported tools and their opt-out variables |
| `notel version` / `notel --version` | Show the application version |
| `notel help` | Show help |

When `--fix` is combined with `--all` and/or `--verbose`, the audit report is still shown on `stderr` and the shell commands remain redirect-safe on `stdout`.

### Example output

```
  notel  —  DevTools Telemetry Auditor
  ────────────────────────────────────────────────────────────

  Package Manager
  ✓  Homebrew              HOMEBREW_NO_ANALYTICS=1
  ✗  Yarn                  YARN_ENABLE_TELEMETRY (not set)

  Runtime
  ✓  Node.js               NODE_NO_TELEMETRY=1

  ...

  ✗  2/3  installed tools have telemetry disabled  (1 need attention)
```

## Supported Tools by Category

### All The Things

| Tool | Environment Variable | Opt-out Value | Data Collection |
|------|---------------------|---------------|-----------------|
| [Do Not Track](https://donottrack.sh) | `DO_NOT_TRACK` | `1` | Signals a global preference to disable telemetry, analytics, crash reporting, and non-essential tracking. |

### Cloud & Deployment

| Tool | Environment Variable | Opt-out Value | Data Collection |
|------|---------------------|---------------|-----------------|
| [Algolia CLI](https://www.algolia.com/doc/tools/cli/telemetry) | `ALGOLIA_CLI_TELEMETRY` | `0` | Command usage, operating system details, CLI version, and local profile metadata. |
| [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest#cli-configuration-values-and-environment-variables) | `AZURE_CORE_COLLECT_TELEMETRY` | `0` | Command usage, performance metrics, and error-rate telemetry.<p><p>`az config set core.collect_telemetry=no`|
| [CockroachDB](https://www.cockroachlabs.com/docs/stable/telemetry) | `COCKROACH_SKIP_ENABLING`<br/>`_DIAGNOSTIC_REPORTING` | `true` | Cluster diagnostics, telemetry, and crash reports sent to Cockroach Labs.<p><p>NOTE: The environment variable only has effect if be set before initialising the first node of the cluster. After a cluster is running update the  [`diagnostics.reporting.enabled`](https://www.cockroachlabs.com/docs/stable/diagnostics-reporting#at-cluster-initialization) setting. Telemetry is required during the 30 day self-service Enterprise Trial (see: [Licensing FAQs](https://www.cockroachlabs.com/docs/stable/licensing-faqs)).|
| [GitHub CLI](https://docs.github.com/en/github-cli/github-cli/github-cli-telemetry) | `GH_TELEMETRY` / `DO_NOT_TRACK` | `false` / `1` | Command usage, feature usage, and diagnostic telemetry from the GitHub CLI. <p><p>`gh config get telemetry`<p><p>`gh config set telemetry disabled`<p><p>`gh config set telemetry enabled`|
| [Google Cloud SDK](https://cloud.google.com/sdk/docs/usage-statistics) | `CLOUDSDK_CORE_DISABLE`<br/>`_USAGE_REPORTING` | `true` | Anonymized command execution metrics, timing, and error status.<p><p> NOTE: Unless you opt-in during Google Cloud CLI installation, gcloud CLI software does not collect usage statistics.<p><p>`gcloud config set disable_usage_reporting` `true|false`|
| [Netlify CLI](https://docs.netlify.com/api-and-cli-guides/cli-guides/get-started-with-cli/#usage-data-collection) | `NETLIFY_TELEMETRY_DISABLED` | `1` | Anonymous CLI usage and diagnostic telemetry. <p><p>`netlify --telemetry-disable`<p><p>`netlify --telemetry-enable`
| [Turborepo](https://turborepo.com/docs/telemetry) | `TURBO_TELEMETRY_DISABLED` | `1` | Anonymous command usage, host information, and repo or task metrics. <p><p>`turbo telemetry disable`<p><p>`turbo telemetry enable`<p><p>`turbo telemetry status`|


### Framework

| Tool | Environment Variable | Opt-out Value | Data Collection |
|------|---------------------|---------------|-----------------|
| [Angular CLI](https://angular.dev/cli/analytics) | `NG_CLI_ANALYTICS` | `false` or `ci` | Command usage, selected flags, workspace shape, and local version metadata.<p><p>`ng analytics [disable|enable|info|prompt]` |
| [Astro](https://docs.astro.build/en/reference/cli-reference/#astro-telemetry) | `ASTRO_TELEMETRY_DISABLED` | `1` | Anonymous command usage, integration usage, and project metadata. <p><p>`astro telemetry disable`<p><p>`astro telemetry enable`|
| [Expo CLI](https://docs.expo.dev/more/expo-cli/#telemetry) | `EXPO_NO_TELEMETRY` | `1` | Anonymous CLI usage and diagnostics for Expo development workflows. |
| [Gatsby](https://www.gatsbyjs.com/docs/telemetry/) | `GATSBY_TELEMETRY_DISABLED` | `1` | Anonymous command usage, plugin usage, and machine characteristics. <p><p>`gatsby telemetry --disable`|
| [Next.js](https://nextjs.org/telemetry) | `NEXT_TELEMETRY_DISABLED` | `1` | Anonymous command usage, session timing, and project or machine characteristics. <p><p>`npx @nuxt/telemetry [status|enable|disable]`<br/>`[-g,--global] [dir]`|
| [Nuxt](https://github.com/nuxt/telemetry#nuxt-telemetry-module) | `NUXT_TELEMETRY_DISABLED` | `1` | Anonymous command usage, module usage, and environment characteristics. |
| [Storybook](https://storybook.js.org/docs/configure/telemetry) | `STORYBOOK_DISABLE_TELEMETRY` | `true` | Anonymous framework usage, addon data, and environment metadata.<p><p>`npm run storybook -- --disable-telemetry` |

### Package Manager

| Tool | Environment Variable | Opt-out Value | Data Collection |
|------|---------------------|---------------|-----------------|
| [Homebrew](https://docs.brew.sh/Analytics) | `HOMEBREW_NO_ANALYTICS` | `1` | Anonymous install events, command usage, and build error metadata.<p><p> `brew analytics state`<p><p>`brew analytics off`<p><p>`brew analytics on`|
| [Poetry](https://python-poetry.org/docs/configuration/#using-environment-variables) | `POETRY_TELEMETRY_ENABLED` | `0` | Anonymous usage statistics and telemetry about Poetry command execution. |
| [Yarn (v2+)](https://yarnpkg.com/advanced/telemetry) | `YARN_ENABLE_TELEMETRY` | `0` | Anonymous command usage and version information. |

### Runtimes & SDKs

| Tool | Environment Variable | Opt-out Value | Data Collection |
|------|---------------------|---------------|-----------------|
| [Bun](https://bun.sh/docs/runtime/bunfig#telemetry) | `DO_NOT_TRACK` | `1` | Anonymous telemetry and crash-reporting signals from Bun tooling. |
| [.NET SDK](https://learn.microsoft.com/en-us/dotnet/core/tools/telemetry) | `DOTNET_CLI_TELEMETRY_OPTOUT` | `1` | CLI command usage, SDK versions, and exception metadata. |
| [Flutter](https://docs.flutter.dev/reference/crash-reporting) | `FLUTTER_CLI_CRASH_REPORTING` | `false` | CLI crash reports, tool usage signals, and local environment metadata.<p><p><p><p>Note: ⁠FLUTTER_CLI_CRASH_REPORTING does not disable the general analytics reporting setting which must be set via the CLI: <p><p>`flutter config --no-analytics`|
| [Node.js](https://nodejs.org/api/cli.html#node_no_telemetry1) | `DISABLE_TELEMETRY` / `NODE_NO_TELEMETRY` | `1` | Anonymous runtime and CLI telemetry in Node features that support it. |

### Security & Observability

| Tool | Environment Variable | Opt-out Value | Data Collection |
|------|---------------------|---------------|-----------------|
| [Semgrep](https://semgrep.dev/docs/metrics) | `SEMGREP_SEND_METRICS` | `off` | Usage metrics about Semgrep runs, registry usage, and login-related activity. <p><p>`semgrep --metrics auto|on|off ...`|
| [Sentry CLI](https://cli.sentry.dev/configuration/) | `SENTRY_CLI_NO_TELEMETRY` | `1` | Anonymous CLI usage and diagnostic telemetry. |

## Adding a new tool

Add an entry to `internal/tools/registry.go`:

```go
{
    Name:        "My Tool",
    Description: "Description of the tool",
    DocumentationURL: "https://example.com/docs/telemetry",
    DataCollection:   "Anonymous command usage and diagnostic metadata.",
    Binary:           "mytool",        // binary name checked via PATH
    Category:         "Framework",
    EnvChecks: []EnvCheck{
        {
            Name:        "MYTOOL_TELEMETRY_DISABLED",
            ValidValues: []string{"1", "true"},
            Description: "Disables My Tool telemetry",
        },
    },
},
```

Keep registry entries sorted alphabetically by `Category` and then by `Name`.

## Project structure

```
notel/
├── main.go
├── cmd/
│   ├── root.go      # Cobra root command
│   ├── audit.go     # notel audit
│   └── list.go      # notel list
└── internal/
    ├── tools/
    │   └── registry.go   # Tool definitions (add new tools here)
    ├── check/
    │   └── checker.go    # Detection & env-var audit logic
    └── output/
        └── formatter.go  # Lipgloss-styled output
```

## License
Licensed under the Apache License, Version 2.0. See `LICENSE`.

## Feedback or suggestions?

Bug reports, feature requests, and questions can be posted in the [Issues](https://github.com/stennie/notel/issues) section on GitHub.
