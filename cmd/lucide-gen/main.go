package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	lucidegen "github.com/peterszarvas94/lucide-templ-gen"
)

const version = "1.3.3"

func main() {
	var (
		outputDir      = flag.String("output", ".", "Output directory")
		packageName    = flag.String("package", "icons", "Package name")
		prefix         = flag.String("prefix", "", "Function name prefix")
		categories     = flag.String("categories", "", "Comma-separated categories to include (empty = all)")
		icons          = flag.String("icons", "", "Comma-separated icon names to include (empty = all)")
		removeIcons    = flag.String("remove", "", "Comma-separated icon names to remove from final output")
		skipRegistry   = flag.Bool("skip-registry", false, "Skip generating registry.templ")
		skipCategories = flag.Bool("skip-categories", false, "Skip generating categories.go")
		includeSearch  = flag.Bool("search", false, "Include search functionality (fetches metadata)")
		mergeExisting  = flag.Bool("merge", true, "Merge with already generated icons in output directory")
		dryRun         = flag.Bool("dry-run", false, "Show what would be generated without creating files")
		verbose        = flag.Bool("verbose", false, "Enable verbose output")
		showVersion    = flag.Bool("version", false, "Show version information")
		help           = flag.Bool("help", false, "Show help information")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *showVersion {
		fmt.Printf("lucide-gen version %s\n", version)
		return
	}

	// Parse categories
	var categoryList []string
	if *categories != "" {
		categoryList = strings.Split(*categories, ",")
		for i, cat := range categoryList {
			categoryList[i] = strings.TrimSpace(cat)
		}
	}

	// Create config
	requestedIconSet := parseIconsFlag(*icons)
	requestedIcons := make([]string, 0, len(requestedIconSet))
	for name := range requestedIconSet {
		requestedIcons = append(requestedIcons, name)
	}

	removedIconSet := parseIconsFlag(*removeIcons)
	removedIcons := make([]string, 0, len(removedIconSet))
	for name := range removedIconSet {
		removedIcons = append(removedIcons, name)
	}

	config := lucidegen.Config{
		OutputDir:      *outputDir,
		PackageName:    *packageName,
		Prefix:         *prefix,
		Categories:     categoryList,
		RequestedIcons: requestedIcons,
		RemovedIcons:   removedIcons,
		SkipRegistry:   *skipRegistry,
		SkipCategories: *skipCategories,
		IncludeSearch:  *includeSearch,
		MergeExisting:  *mergeExisting,
		DryRun:         *dryRun,
		Verbose:        *verbose,
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Generate icons
	result, err := lucidegen.Generate(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Generation failed: %v\n", err)
		os.Exit(1)
	}

	// Show results
	if !config.DryRun {
		fmt.Printf("✅ Successfully generated %d icons in %v\n", result.IconsGenerated, result.Duration)
		if config.Verbose {
			fmt.Printf("📁 Files created:\n")
			for _, file := range result.FilesCreated {
				fmt.Printf("   %s\n", file)
			}
			fmt.Printf("📂 Categories: %s\n", strings.Join(result.Categories, ", "))
		}
	}
}

func validateConfig(config lucidegen.Config) error {
	if config.OutputDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	if config.PackageName == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	// Validate package name
	if !isValidPackageName(config.PackageName) {
		return fmt.Errorf("invalid package name: %s", config.PackageName)
	}

	return nil
}

func isValidPackageName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// Must start with letter or underscore
	if !isLetter(rune(name[0])) && name[0] != '_' {
		return false
	}

	// Rest can be letters, digits, or underscores
	for _, r := range name[1:] {
		if !isLetter(r) && !isDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func parseIconsFlag(value string) map[string]struct{} {
	icons := make(map[string]struct{})
	if strings.TrimSpace(value) == "" {
		return icons
	}

	for _, name := range strings.Split(value, ",") {
		normalized := strings.ToLower(strings.TrimSpace(name))
		if normalized == "" {
			continue
		}
		icons[normalized] = struct{}{}
	}

	return icons
}

func showHelp() {
	fmt.Printf(`lucide-gen - Generate type-safe Templ components for Lucide Icons

Usage:
  lucide-gen [options]

Options:
  -output string      Output directory (default ".")
  -package string     Package name (default "icons")
  -prefix string      Function name prefix (default "")
  -categories string  Comma-separated categories to include (default: all)
  -icons string       Comma-separated icon names to include (default: all)
  -remove string      Comma-separated icon names to remove from final output
  -skip-registry      Skip generating registry.templ
  -skip-categories    Skip generating categories.go
  -search            Include search functionality (fetches metadata)
  -merge             Merge with already generated icons in output directory (default true)
  -dry-run           Show what would be generated without creating files
  -verbose           Enable verbose output
  -version           Show version information
  -help              Show this help message

Examples:
  # Generate all icons in current directory
  lucide-gen

  # Generate specific categories
  lucide-gen -categories "navigation,actions,media"

  # Generate only selected icons
  lucide-gen -icons "a-arrow-down,search,x"

  # Generate with custom package and prefix
  lucide-gen -output ./icons -package icons -prefix Lucide

  # Dry run to preview
  lucide-gen -dry-run -verbose

Categories:
  navigation      - home, menu, chevron-*, arrow-*, etc.
  actions         - plus, minus, edit, trash, save, etc.
  media          - play, pause, stop, volume, etc.
  communication  - mail, phone, message, bell, etc.
  files          - file, folder, download, upload, etc.
  ui             - eye, lock, search, check, etc.
  data           - database, server, cloud, chart, etc.
  devices        - smartphone, laptop, monitor, etc.
  social         - heart, star, share, thumbs-up, etc.
  weather        - sun, moon, cloud, rain, etc.
  transportation - car, plane, bike, etc.
  business       - briefcase, building, wallet, etc.

For more information, visit: https://github.com/peterszarvas94/lucide-templ-gen
`)
}
