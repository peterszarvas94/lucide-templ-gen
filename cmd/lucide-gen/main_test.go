package main

import "testing"

func TestParseIconsFlag(t *testing.T) {
	parsed := parseIconsFlag(" a-arrow-down,search ,X,,search ")

	if len(parsed) != 3 {
		t.Fatalf("expected 3 unique icons, got %d", len(parsed))
	}

	expected := []string{"a-arrow-down", "search", "x"}
	for _, name := range expected {
		if _, ok := parsed[name]; !ok {
			t.Fatalf("expected icon %q to be parsed", name)
		}
	}
}
