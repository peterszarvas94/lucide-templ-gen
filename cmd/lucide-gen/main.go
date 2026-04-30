package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	lucidegen "github.com/peterszarvas94/lucide-templ-gen/v2"
)

const version = "2.0.2"

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		runAdd(args)
	case "remove":
		runRemove(args)
	case "list":
		runList(args)
	case "sync":
		runSync(args)
	case "help", "--help", "-h":
		showHelpFor(args)
	case "version", "--version", "-version":
		fmt.Printf("lucide-gen version %s\n", version)
	default:
		exitErr("unknown command: %s", command)
	}
}

func runAdd(args []string) {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		showAddHelp()
		return
	}

	outputDir, all, iconsArg := parseOperationArgs("add", args)

	current, err := lucidegen.ReadRegistryIconNames(outputDir)
	if err != nil {
		exitErr("failed to read current registry: %v", err)
	}

	currentSet := sliceToSet(current)
	before := len(currentSet)

	if all {
		allNames, err := lucidegen.ListAvailableIconNames()
		if err != nil {
			exitErr("add --all failed: %v", err)
		}
		result, err := lucidegen.GenerateFromIconNames(lucidegen.Config{OutputDir: outputDir, PackageName: "icons"}, allNames)
		if err != nil {
			exitErr("add --all failed: %v", err)
		}
		fmt.Printf("Added all icons. Before: %d, After: %d\n", before, result.IconsGenerated)
		return
	}

	toAdd := parseIconsArg(iconsArg)
	for _, name := range toAdd {
		currentSet[name] = struct{}{}
	}

	finalNames := sortedSetKeys(currentSet)
	result, err := lucidegen.GenerateFromIconNames(lucidegen.Config{OutputDir: outputDir, PackageName: "icons"}, finalNames)
	if err != nil {
		exitErr("add failed: %v", err)
	}

	addedCount := len(currentSet) - before
	fmt.Printf("Added %d icon(s). Before: %d, After: %d\n", addedCount, before, result.IconsGenerated)
}

func runRemove(args []string) {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		showRemoveHelp()
		return
	}

	outputDir, all, iconsArg := parseOperationArgs("remove", args)

	current, err := lucidegen.ReadRegistryIconNames(outputDir)
	if err != nil {
		exitErr("failed to read current registry: %v", err)
	}

	currentSet := sliceToSet(current)
	before := len(currentSet)

	if all {
		result, err := lucidegen.GenerateFromIconNames(lucidegen.Config{OutputDir: outputDir, PackageName: "icons"}, []string{})
		if err != nil {
			exitErr("remove --all failed: %v", err)
		}
		fmt.Printf("Removed all icons. Before: %d, After: %d\n", before, result.IconsGenerated)
		return
	}

	toRemove := parseIconsArg(iconsArg)
	removed := 0
	for _, name := range toRemove {
		if _, ok := currentSet[name]; ok {
			delete(currentSet, name)
			removed++
		}
	}

	finalNames := sortedSetKeys(currentSet)
	result, err := lucidegen.GenerateFromIconNames(lucidegen.Config{OutputDir: outputDir, PackageName: "icons"}, finalNames)
	if err != nil {
		exitErr("remove failed: %v", err)
	}

	fmt.Printf("Removed %d icon(s). Before: %d, After: %d\n", removed, before, result.IconsGenerated)
}

func runList(args []string) {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		showListHelp()
		return
	}

	outputDir, err := parsePathOnlyArgs("list", args)
	if err != nil {
		exitErr("%v", err)
	}

	current, err := lucidegen.ReadRegistryIconNames(outputDir)
	if err != nil {
		exitErr("failed to read current registry: %v", err)
	}

	sort.Strings(current)
	for _, name := range current {
		fmt.Println(name)
	}
	fmt.Printf("Total: %d\n", len(current))
}

func runSync(args []string) {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		showSyncHelp()
		return
	}

	outputDir, err := parsePathOnlyArgs("sync", args)
	if err != nil {
		exitErr("%v", err)
	}

	current, err := lucidegen.ReadRegistryIconNames(outputDir)
	if err != nil {
		exitErr("failed to read current registry: %v", err)
	}

	before := len(current)
	result, err := lucidegen.GenerateFromIconNames(lucidegen.Config{OutputDir: outputDir, PackageName: "icons"}, current)
	if err != nil {
		exitErr("sync failed: %v", err)
	}

	fmt.Printf("Synced icons. Before: %d, After: %d\n", before, result.IconsGenerated)
}

func parseOperationArgs(command string, args []string) (outputDir string, all bool, iconsArg string) {
	outputDir = "./icons"
	remaining := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		if args[i] == "--output" {
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				exitErr("%s: --output requires a value", command)
			}
			outputDir = args[i+1]
			i++
			continue
		}
		remaining = append(remaining, args[i])
	}

	if len(remaining) != 1 {
		exitErr("%s requires exactly one arg: <icons> or --all", command)
	}

	if remaining[0] == "--all" {
		return outputDir, true, ""
	}

	if strings.HasPrefix(remaining[0], "-") {
		exitErr("%s invalid arg: %s", command, remaining[0])
	}

	return outputDir, false, remaining[0]
}

func parsePathOnlyArgs(command string, args []string) (string, error) {
	outputDir := "./icons"
	if len(args) == 0 {
		return outputDir, nil
	}
	if len(args) == 2 && args[0] == "--output" && strings.TrimSpace(args[1]) != "" {
		return args[1], nil
	}
	return "", fmt.Errorf("%s only supports optional --output <dir>", command)
}

func parseIconsArg(value string) []string {
	icons := make(map[string]struct{})
	for _, name := range strings.Split(value, ",") {
		normalized := strings.ToLower(strings.TrimSpace(name))
		if normalized == "" {
			continue
		}
		icons[normalized] = struct{}{}
	}
	if len(icons) == 0 {
		exitErr("icons arg is empty")
	}
	return sortedSetKeys(icons)
}

func sliceToSet(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

func sortedSetKeys(set map[string]struct{}) []string {
	items := make([]string, 0, len(set))
	for k := range set {
		items = append(items, k)
	}
	sort.Strings(items)
	return items
}

func exitErr(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func showHelp() {
	fmt.Printf(`Usage:
  lucide-gen <command> [arg]
Commands:
  add <icons>       Add icon(s) to current set
  add --all         Add all Lucide icons
  remove <icons>    Remove icon(s) from current set
  remove --all      Remove all icons
  list              List current icons
  sync              Regenerate files from current icon set
  help              Show help
  version           Show version
Arg:
  <icons>           Comma-separated icon names
                    Example: "arrow-left,arrow-right"
Options:
  --output <dir>    Output folder (default: ./icons)
Examples:
  lucide-gen add "arrow-left,arrow-right"
  lucide-gen add --all
  lucide-gen remove "arrow-left"
  lucide-gen remove --all
  lucide-gen list
  lucide-gen sync
`)
}

func showHelpFor(args []string) {
	if len(args) == 0 {
		showHelp()
		return
	}

	switch args[0] {
	case "add":
		showAddHelp()
	case "remove":
		showRemoveHelp()
	case "list":
		showListHelp()
	case "sync":
		showSyncHelp()
	case "help":
		showHelp()
	case "version":
		fmt.Println("Usage: lucide-gen version")
	default:
		exitErr("unknown help topic: %s", args[0])
	}
}

func showAddHelp() {
	fmt.Printf(`Usage:
  lucide-gen add <icons> [--output <dir>]
  lucide-gen add --all [--output <dir>]

Description:
  Add icon(s) to current set and regenerate files.

Arg:
  <icons>  Comma-separated icon names
           Example: "arrow-left,arrow-right"

Option:
  --output <dir>  Output folder (default: ./icons)
`)
}

func showRemoveHelp() {
	fmt.Printf(`Usage:
  lucide-gen remove <icons> [--output <dir>]
  lucide-gen remove --all [--output <dir>]

Description:
  Remove icon(s) from current set and regenerate files.

Arg:
  <icons>  Comma-separated icon names
           Example: "arrow-left,arrow-right"

Option:
  --output <dir>  Output folder (default: ./icons)
`)
}

func showListHelp() {
	fmt.Printf(`Usage:
  lucide-gen list [--output <dir>]

Description:
  List current icons from registry.

Option:
  --output <dir>  Output folder (default: ./icons)
`)
}

func showSyncHelp() {
	fmt.Printf(`Usage:
  lucide-gen sync [--output <dir>]

Description:
  Regenerate files from current registry set.

Option:
  --output <dir>  Output folder (default: ./icons)
`)
}
