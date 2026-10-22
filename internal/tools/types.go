// Package tools defines the registry of developer tools that notel can audit.
package tools

// DefaultState is the telemetry state when none of a tool's opt-out
// variables are present in the environment.
//
// Most developer tools collect telemetry until the user opts out, so the
// zero value is DefaultEnabled and an unset variable fails the audit.
// A few tools do the opposite. Google Cloud SDK does not upload usage
// statistics unless the user opted in during installation, so an unset
// CLOUDSDK_CORE_DISABLE_USAGE_REPORTING is already the private state and
// must not be reported as a failure.
type DefaultState int

const (
	// DefaultEnabled means an unset opt-out variable fails the check.
	// The tool is treated as still collecting telemetry.
	DefaultEnabled DefaultState = iota

	// DefaultDisabled means an unset opt-out variable passes the check.
	// A value that is set but not a recognized opt-out still fails, because
	// notel will not guess that an unknown value means "disabled".
	DefaultDisabled
)

// Requirement is how a tool's EnvChecks combine into one pass/fail result.
//
// The historical rule, and the zero value, is AllOf: every variable has to
// match. That is correct when each variable controls a distinct switch.
// It is the wrong rule when the variables are alternatives. GitHub CLI
// honors either GH_TELEMETRY=false or DO_NOT_TRACK=1. Forcing both made a
// machine that had already opted out via DO_NOT_TRACK look like it was
// still sending GitHub CLI telemetry.
type Requirement int

const (
	// AllOf requires every EnvCheck to pass.
	AllOf Requirement = iota

	// AnyOf requires one EnvCheck to pass. The other checks are still
	// evaluated and shown, so the report can say which alternative matched
	// and which ones were left unset.
	AnyOf
)

// EnvCheck defines an environment variable requirement for a tool.
type EnvCheck struct {
	// Name is the environment variable, for example HOMEBREW_NO_ANALYTICS.
	Name string

	// ValidValues are the accepted opt-out spellings.
	// An empty list keeps the original rule: any non-empty value counts.
	// The first entry is the spelling --fix suggests. Put the preferred
	// form first ("1" before "true", "false" before "0").
	ValidValues []string

	// IgnoreCase compares the current value with strings.EqualFold.
	// Set it when the vendor accepts any capitalization. The registry then
	// lists "true" once instead of "true", "True", and "TRUE".
	// Leave it false when the vendor documents one exact spelling, such as
	// Homebrew's HOMEBREW_NO_ANALYTICS=1 or Next.js's NEXT_TELEMETRY_DISABLED=1.
	IgnoreCase bool

	// Description is the human-readable line shown in verbose and list output.
	Description string
}

// Tool represents a developer tool that can be audited for telemetry settings.
type Tool struct {
	Name             string
	Description      string
	DocumentationURL string
	DataCollection   string

	// Binary is the executable looked up on PATH.
	// An empty Binary means the tool has no executable of its own and is
	// always in scope. Do Not Track uses that: it is a convention, not a CLI.
	Binary string

	Category string

	// Default is the state when every opt-out variable is unset.
	// Leave it as DefaultEnabled unless the vendor documents that collection
	// is opt-in.
	Default DefaultState

	// OptOut combines EnvChecks. Leave it as AllOf unless any one of the
	// variables is sufficient on its own.
	OptOut Requirement

	EnvChecks []EnvCheck

	// Commands are extra ways to disable telemetry that are not environment
	// variables: `brew analytics off`, `gh config set telemetry disabled`,
	// `flutter config --no-analytics`.
	//
	// notel prints these hints. It does not run them, and it does not treat
	// them as evidence that telemetry is off. A tool can pass its env checks
	// while one of these commands is still required. Flutter is the case
	// that motivated the field: FLUTTER_CLI_CRASH_REPORTING does not disable
	// the general analytics setting.
	Commands []string
}

// CaseInsensitive builds an EnvCheck that accepts any capitalization of values.
//
// List the canonical spellings only. "true" matches "True" and "TRUE"
// because the audit compares with strings.EqualFold. The suggested fix
// value is values[0], so the preferred spelling has to come first.
//
// The returned ValidValues slice is a copy. Callers can reuse or append to
// the input slice later without changing the check.
func CaseInsensitive(name, description string, values ...string) EnvCheck {
	copied := make([]string, len(values))
	copy(copied, values)
	return EnvCheck{
		Name:        name,
		Description: description,
		ValidValues: copied,
		IgnoreCase:  true,
	}
}
