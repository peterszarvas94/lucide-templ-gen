//go:build integration

package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIWorkflowIntegration(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "icons")

	assertCommandContains(t, runCLI(t, "help"), "Commands:")
	assertCommandContains(t, runCLI(t, "version"), "lucide-gen version 2.0.0")

	assertCommandContains(t, runCLI(t, "list", "--output", outputDir), "Total: 0")
	assertCommandContains(t, runCLI(t, "add", "arrow-left,arrow-right", "--output", outputDir), "Added 2 icon(s). Before: 0, After: 2")

	listAfterAdd := runCLI(t, "list", "--output", outputDir)
	assertCommandContains(t, listAfterAdd, "arrow-left")
	assertCommandContains(t, listAfterAdd, "arrow-right")
	assertCommandContains(t, listAfterAdd, "Total: 2")

	assertCommandContains(t, runCLI(t, "remove", "arrow-left", "--output", outputDir), "Removed 1 icon(s). Before: 2, After: 1")

	listAfterRemove := runCLI(t, "list", "--output", outputDir)
	assertCommandContains(t, listAfterRemove, "arrow-right")
	assertCommandNotContains(t, listAfterRemove, "arrow-left")
	assertCommandContains(t, listAfterRemove, "Total: 1")

	assertCommandContains(t, runCLI(t, "sync", "--output", outputDir), "Synced icons. Before: 1, After: 1")
	assertCommandContains(t, runCLI(t, "add", "--all", "--output", outputDir), "Added all icons. Before: 1, After:")
	assertCommandContains(t, runCLI(t, "remove", "--all", "--output", outputDir), "Removed all icons. Before:")
	assertCommandContains(t, runCLI(t, "list", "--output", outputDir), "Total: 0")
}

func runCLI(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: go run . %s\nerror: %v\noutput:\n%s", strings.Join(args, " "), err, string(out))
	}
	return string(out)
}

func assertCommandContains(t *testing.T, output, needle string) {
	t.Helper()
	if !strings.Contains(output, needle) {
		t.Fatalf("expected output to contain %q\nactual output:\n%s", needle, output)
	}
}

func assertCommandNotContains(t *testing.T, output, needle string) {
	t.Helper()
	if strings.Contains(output, needle) {
		t.Fatalf("expected output not to contain %q\nactual output:\n%s", needle, output)
	}
}
