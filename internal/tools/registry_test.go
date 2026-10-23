package tools

import (
	"strings"
	"testing"
)

func TestRegistryMetadataComplete(t *testing.T) {
	registry := Registry()
	if len(registry) == 0 {
		t.Fatal("Registry() returned no tools")
	}

	for _, tool := range registry {
		if tool.Name == "" {
			t.Fatal("registry contains tool with empty name")
		}
		if tool.Description == "" {
			t.Fatalf("%s is missing a description", tool.Name)
		}
		if tool.DocumentationURL == "" {
			t.Fatalf("%s is missing a documentation URL", tool.Name)
		}
		if tool.DataCollection == "" {
			t.Fatalf("%s is missing a data collection description", tool.Name)
		}
		if tool.Category == "" {
			t.Fatalf("%s is missing a category", tool.Name)
		}
		if len(tool.EnvChecks) == 0 {
			t.Fatalf("%s has no env checks", tool.Name)
		}

		for _, check := range tool.EnvChecks {
			if check.Name == "" {
				t.Fatalf("%s has an env check with empty name", tool.Name)
			}
			if check.Description == "" {
				t.Fatalf("%s env check %s is missing a description", tool.Name, check.Name)
			}
		}
	}
}

func TestRegistryToolNamesUnique(t *testing.T) {
	seen := map[string]struct{}{}

	for _, tool := range Registry() {
		if _, exists := seen[tool.Name]; exists {
			t.Fatalf("duplicate tool name in registry: %s", tool.Name)
		}
		seen[tool.Name] = struct{}{}
	}
}

func TestRegistryOptOutPolicies(t *testing.T) {
	byName := map[string]Tool{}
	for _, tool := range Registry() {
		byName[tool.Name] = tool
	}

	anyOf := []string{"Hugging Face CLI", "Railway CLI"}
	for _, name := range anyOf {
		tool, ok := byName[name]
		if !ok {
			t.Fatalf("registry missing %s", name)
		}
		if tool.OptOut != AnyOf {
			t.Fatalf("%s OptOut = %v, want AnyOf", name, tool.OptOut)
		}
		if len(tool.EnvChecks) < 2 {
			t.Fatalf("%s AnyOf policy needs at least two env checks", name)
		}
	}

	githubCLI := byName["GitHub CLI"]
	if githubCLI.OptOut != FirstSet {
		t.Fatalf("GitHub CLI OptOut = %v, want FirstSet", githubCLI.OptOut)
	}
	if len(githubCLI.EnvChecks) != 2 || githubCLI.EnvChecks[0].Name != "GH_TELEMETRY" || !githubCLI.EnvChecks[0].TrimSpace {
		t.Fatalf("GitHub CLI checks = %#v, want GH_TELEMETRY then DO_NOT_TRACK with trim", githubCLI.EnvChecks)
	}
	if githubCLI.EnvChecks[0].ValidValues[0] != "false" {
		t.Fatalf("GH_TELEMETRY preferred value = %q, want false", githubCLI.EnvChecks[0].ValidValues[0])
	}

	gcloud := byName["Google Cloud SDK"]
	if gcloud.Default != DefaultDisabled {
		t.Fatalf("Google Cloud SDK Default = %v, want DefaultDisabled", gcloud.Default)
	}

	flutter := byName["Flutter"]
	if len(flutter.Commands) == 0 {
		t.Fatal("Flutter is missing the analytics CLI hint")
	}

	// AnyOf plus DefaultDisabled would treat a single unset variable as a
	// pass and hide the other alternatives. The catalog does not use that
	// combination. This test is the reminder to decide the sentence in the
	// Tool comment before adding one.
	for _, tool := range Registry() {
		if tool.OptOut == AnyOf && tool.Default == DefaultDisabled {
			t.Fatalf("%s combines AnyOf with DefaultDisabled", tool.Name)
		}
		if tool.OptOut == FirstSet && tool.Default == DefaultDisabled {
			t.Fatalf("%s combines FirstSet with DefaultDisabled", tool.Name)
		}
	}

	for _, tool := range Registry() {
		for _, check := range tool.EnvChecks {
			if !check.IgnoreCase {
				continue
			}
			for _, value := range check.ValidValues {
				switch value {
				case "True", "TRUE", "False", "FALSE", "Off", "OFF":
					t.Fatalf("%s check %s repeats capitalization %q; IgnoreCase already covers it", tool.Name, check.Name, value)
				}
			}
		}
	}
}

func TestCaseInsensitiveCopiesValues(t *testing.T) {
	values := []string{"1", "true"}
	check := CaseInsensitive("DO_NOT_TRACK", "convention", values...)
	values[0] = "mutated"
	if check.ValidValues[0] != "1" || !check.IgnoreCase {
		t.Fatalf("CaseInsensitive() = %#v, want an independent case-insensitive check", check)
	}
}

func TestRegistrySortedByCategoryThenName(t *testing.T) {
	registry := Registry()
	for i := 1; i < len(registry); i++ {
		prev := registry[i-1]
		curr := registry[i]

		prevCategory := strings.ToLower(prev.Category)
		currCategory := strings.ToLower(curr.Category)
		if prevCategory > currCategory {
			t.Fatalf("registry category order incorrect at %q before %q", prev.Category, curr.Category)
		}
		if prevCategory == currCategory && strings.ToLower(prev.Name) > strings.ToLower(curr.Name) {
			t.Fatalf("registry name order incorrect within %q: %q before %q", curr.Category, prev.Name, curr.Name)
		}
	}
}
