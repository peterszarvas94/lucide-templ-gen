package lucidegen

import (
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
