# notel

**notel** audits telemetry opt-out settings across common developer tools.

Telemetry collection is a privacy concern for many developers, especially where collection
of data may not be obvious without consulting the documentation.

This tool helps identify and opt-out of telemetry for popular developer tools. It checks
whether tools are installed and whether the environment variables that disable telemetry
collection are correctly set.

Licensed under the Apache License, Version 2.0. See `LICENSE`.

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

## Supported Tools

| Tool | Environment Variable | Opt-out Value |
|------|---------------------|---------------|
| Homebrew | `HOMEBREW_NO_ANALYTICS` | `1` |
| Yarn (v2+) | `YARN_ENABLE_TELEMETRY` | `0` |
| Node.js | `DISABLE_TELEMETRY` / `NODE_NO_TELEMETRY` | `1` |
| .NET SDK | `DOTNET_CLI_TELEMETRY_OPTOUT` | `1` |
| Flutter | `FLUTTER_CLI_CRASH_REPORTING` | `false` |
| Next.js | `NEXT_TELEMETRY_DISABLED` | `1` |
| Gatsby | `GATSBY_TELEMETRY_DISABLED` | `1` |
| Nuxt | `NUXT_TELEMETRY_DISABLED` | `1` |
| Angular CLI | `NG_CLI_ANALYTICS` | `false` or `ci` |
| Astro | `ASTRO_TELEMETRY_DISABLED` | `1` |
| Turborepo | `TURBO_TELEMETRY_DISABLED` | `1` |
| Terraform | `CHECKPOINT_DISABLE` | `1` |
| Google Cloud SDK | `CLOUDSDK_CORE_DISABLE_USAGE_REPORTING` | `true` |
| Azure CLI | `AZURE_CORE_COLLECT_TELEMETRY` | `0` |
| Netlify CLI | `NETLIFY_TELEMETRY_DISABLED` | `1` |
| Sentry CLI | `SENTRY_CLI_NO_TELEMETRY` | `1` |

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
