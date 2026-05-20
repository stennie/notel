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
