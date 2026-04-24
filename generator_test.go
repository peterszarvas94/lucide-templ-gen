package lucidegen

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFilterIconsByRequestedNamesExactFiltering(t *testing.T) {
	icons := []IconData{
		{Name: "search", Category: "ui"},
		{Name: "x", Category: "ui"},
		{Name: "plus", Category: "actions"},
		{Name: "menu", Category: "navigation"},
	}

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
	icons := []IconData{
		{Name: "search", Category: "ui"},
		{Name: "x", Category: "ui"},
	}

	unknown := findUnknownRequestedIcons(icons, []string{"search", "does-not-exist"})
	want := []string{"does-not-exist"}
	if !reflect.DeepEqual(unknown, want) {
		t.Fatalf("unexpected unknown icons: got %v, want %v", unknown, want)
	}
}

func TestFilterIconsWithCategoriesIntersection(t *testing.T) {
	icons := []IconData{
		{Name: "search", Category: "ui"},
		{Name: "x", Category: "ui"},
		{Name: "plus", Category: "actions"},
	}

	categoryFiltered := filterIconsByCategories(icons, []string{"ui"})
	unknown := findUnknownRequestedIcons(icons, []string{"search", "plus"})
	if len(unknown) != 0 {
		t.Fatalf("expected no unknown icons, got %v", unknown)
	}
	filtered := filterIconsByRequestedNames(categoryFiltered, []string{"search", "plus"})

	if len(filtered) != 1 || filtered[0].Name != "search" {
		t.Fatalf("expected intersection to keep only search, got %v", filtered)
	}
}

func TestGenerateFilesSkipRegistryAndCategories(t *testing.T) {
	tempDir := t.TempDir()
	icons := []IconData{
		{
			Name:     "search",
			FuncName: "Search",
			ViewBox:  "0 0 24 24",
			Content:  "<circle cx=\"11\" cy=\"11\" r=\"8\"></circle>",
			Category: "ui",
		},
	}
	config := Config{
		OutputDir:      tempDir,
		PackageName:    "icons",
		SkipRegistry:   true,
		SkipCategories: true,
	}

	files, err := generateFiles(icons, config)
	if err != nil {
		t.Fatalf("generateFiles failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected exactly one generated file, got %d (%v)", len(files), files)
	}

	iconsFile := filepath.Join(tempDir, "icons.templ")
	if _, err := os.Stat(iconsFile); err != nil {
		t.Fatalf("expected icons.templ to exist: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "registry.templ")); err == nil {
		t.Fatalf("expected registry.templ to be skipped")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "categories.go")); err == nil {
		t.Fatalf("expected categories.go to be skipped")
	}
}
