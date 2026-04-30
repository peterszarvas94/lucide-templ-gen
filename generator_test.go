package lucidegen

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFilterIconsByRequestedNamesExactFiltering(t *testing.T) {
	icons := []IconData{{Name: "search"}, {Name: "x"}, {Name: "plus"}, {Name: "menu"}}

	filtered := filterIconsByRequestedNames(icons, []string{"search", "x", "plus"})
	unknown := findUnknownRequestedIcons(icons, []string{"search", "x", "plus"})
	if len(unknown) != 0 {
		t.Fatalf("expected no unknown icons, got %v", unknown)
	}

	got := make([]string, 0, len(filtered))
	for _, icon := range filtered {
		got = append(got, icon.Name)
	}

	want := []string{"search", "x", "plus"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected filtered icons: got %v, want %v", got, want)
	}
}

func TestFilterIconsByRequestedNamesUnknown(t *testing.T) {
	icons := []IconData{{Name: "search"}, {Name: "x"}}

	unknown := findUnknownRequestedIcons(icons, []string{"search", "does-not-exist"})
	want := []string{"does-not-exist"}
	if !reflect.DeepEqual(unknown, want) {
		t.Fatalf("unexpected unknown icons: got %v, want %v", unknown, want)
	}
}

func TestGenerateFilesCreatesIconsAndRegistry(t *testing.T) {
	tempDir := t.TempDir()
	icons := []IconData{{
		Name:     "search",
		FuncName: "Search",
		ViewBox:  "0 0 24 24",
		Content:  "<circle cx=\"11\" cy=\"11\" r=\"8\"></circle>",
	}}
	config := Config{OutputDir: tempDir, PackageName: "icons"}

	files, err := generateFiles(icons, config)
	if err != nil {
		t.Fatalf("generateFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected exactly two generated files, got %d (%v)", len(files), files)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "icons.templ")); err != nil {
		t.Fatalf("expected icons.templ to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tempDir, "registry.templ")); err != nil {
		t.Fatalf("expected registry.templ to exist: %v", err)
	}
}

func TestReadRegistryIconNamesStrictParseError(t *testing.T) {
	tempDir := t.TempDir()
	bad := "package icons\nconst (\nIconBroken IconName =\n)\n"
	if err := os.WriteFile(filepath.Join(tempDir, "registry.templ"), []byte(bad), 0644); err != nil {
		t.Fatalf("failed writing registry.templ: %v", err)
	}

	_, err := ReadRegistryIconNames(tempDir)
	if err == nil {
		t.Fatalf("expected parse error for malformed registry")
	}
}
