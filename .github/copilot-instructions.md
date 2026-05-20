# Copilot instructions for `notel`

## Build and test commands

- Build the CLI binary from source: `mkdir -p bin && go build -o bin/notel .`
- Preferred release/local build with embedded version: `mkdir -p bin && go build -ldflags "-X github.com/stennie/notel/cmd.Version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" -o bin/notel .`
- Build cross-platform release artifacts into `dist/`: `just release`
- Run release preflight plus artifact builds: `just release-check`
- Shared release packaging lives in `scripts/release.sh`; keep `Justfile` and `.github/workflows/release.yml` using that script so targets and archive contents do not drift.
- Compile-check the module: `go build ./...`
- Run the full test suite: `go test ./...`
- Run a single test by name once `*_test.go` files exist: `go test ./... -run '^TestName$'`
- CI enforcement lives in `.github/workflows/ci.yml` and runs `gofmt`, `go vet`, `go test ./...`, and `go build ./...` on Go 1.26.3.

## High-level architecture

- The intended product is a single-binary Go CLI built around Cobra commands and Charm-based terminal UX. Keep the app cross-platform and dependency-light: standard library plus Cobra and Charm libraries.
- `main.go` is only the entrypoint; it delegates to `cmd.Execute()`.
- The command layer is the Cobra tree in `cmd/`. Command handlers should stay thin and hand off work to the audit and output layers.
- The telemetry source of truth is `internal/tools/registry.go`. Each `Tool` declares its display metadata, category, PATH binary name, and one or more `EnvCheck` definitions.
- The audit engine in `internal/check/checker.go` turns registry entries into `ToolResult` values by:
  - resolving installation via `exec.LookPath`
  - reading environment variables with `os.Getenv`
  - marking a tool as passing only when it is installed and every env check passes
- The rendering layer in `internal/output/formatter.go` is shared by multiple commands:
  - `audit` prints a compact installed-tool view by default, switches to the detailed per-tool block when `--verbose` is set, and can print shell-specific fix commands when `--fix` is used
  - `list` reuses registry-driven results so it can show the documented env vars alongside current installation status
- If interactive flows are added later, prefer Charm libraries (`bubbles`, `huh`) for those surfaces instead of introducing a separate TUI/form stack.

## Key conventions

- The original bootstrap scope was: Go modules, Cobra command structure, a simple `audit` command for environment variables, and Lipgloss output. Keep new work aligned with that layered structure rather than folding logic back into `main.go`.
- Keep the registry as the only place where supported tools are defined. Adding a tool means adding one `Tool` entry with its `EnvChecks`; command handlers should not duplicate tool metadata.
- A `Tool` with an empty `Binary` is treated as always installed. Preserve that behavior for env-only checks instead of adding special cases elsewhere.
- `EnvCheck.ValidValues` is authoritative:
  - if it contains values, the env var must exactly match one of them
  - if it is empty, any non-empty value counts as opted out
- Keep `internal/tools/registry.go` sorted alphabetically by `Category` and then by `Name`. That ordering also drives the grouped UI output.
- `audit` hides uninstalled tools unless `--all` is set, whether it is in compact or `--verbose` mode. `list` operates over the full registry, so changes to filtering should be made deliberately across commands.
- Styling and output rules belong in `internal/output/formatter.go`, not in Cobra handlers. Prefer extending the output package over spreading Lipgloss formatting through command files.
- The project goal is to cover common developer tools across ecosystems (package managers, runtimes, SDKs, frameworks, cloud CLIs, etc.). Prefer extending the registry in a way that keeps tool definitions data-driven and easy to append.
- Keep README, imports, and filesystem layout synchronized when moving packages so the module stays buildable.
