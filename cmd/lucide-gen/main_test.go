package main

import (
	"reflect"
	"testing"
)

func TestParseIconsArg(t *testing.T) {
	got := parseIconsArg(" arrow-left,arrow-right,ARROW-LEFT ")
	want := []string{"arrow-left", "arrow-right"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected icons: got %v, want %v", got, want)
	}
}

func TestParsePathOnlyArgsDefault(t *testing.T) {
	got, err := parsePathOnlyArgs("list", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "./icons" {
		t.Fatalf("unexpected output dir: got %q", got)
	}
}
