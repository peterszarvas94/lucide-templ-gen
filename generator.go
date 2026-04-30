package lucidegen

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Config holds the configuration for icon generation
type Config struct {
	OutputDir   string // Output directory path
	PackageName string // Go package name
	Prefix      string // Optional function/constant name prefix
}

// ListAvailableIconNames returns all icon names available upstream.
func ListAvailableIconNames() ([]string, error) {
	icons, err := fetchLucideIcons(false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch icons: %w", err)
	}
	names := make([]string, 0, len(icons))
	for _, icon := range icons {
		names = append(names, icon.Name)
	}
	sort.Strings(names)
	return names, nil
}

// GenerateFromIconNames generates files from an explicit icon set.
// Unlike Generate, an empty iconNames slice means "generate empty set".
func GenerateFromIconNames(config Config, iconNames []string) (*GenerationResult, error) {
	start := time.Now()

	if config.PackageName == "" {
		config.PackageName = "icons"
	}

	allIcons, err := fetchLucideIcons(false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch icons: %w", err)
	}

	unknown := findUnknownRequestedIcons(allIcons, iconNames)
	if len(unknown) > 0 {
		return nil, fmt.Errorf("unknown icons: %s", strings.Join(unknown, ","))
	}

	requestedSet := normalizeRequestedIconSet(iconNames)
	icons := make([]IconData, 0, len(requestedSet))
	for _, icon := range allIcons {
		if _, ok := requestedSet[icon.Name]; ok {
			icons = append(icons, icon)
		}
	}
	sort.Slice(icons, func(i, j int) bool {
		return icons[i].Name < icons[j].Name
	})

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	files, err := generateFiles(icons, config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate files: %w", err)
	}

	return &GenerationResult{
		IconsGenerated: len(icons),
		FilesCreated:   files,
		Duration:       time.Since(start),
	}, nil
}

// ReadRegistryIconNames reads icon names from registry.templ in outputDir.
func ReadRegistryIconNames(outputDir string) ([]string, error) {
	registryPath := filepath.Join(outputDir, "registry.templ")
	names, err := readIconNamesFromRegistryFile(registryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	return names, nil
}

// IconData represents a parsed Lucide icon
type IconData struct {
	Name         string   `json:"name"`
	FuncName     string   `json:"func_name"`
	ViewBox      string   `json:"view_box"`
	Content      string   `json:"content"`
	Tags         []string `json:"tags"`
	Contributors []string `json:"contributors"`
	Keywords     []string `json:"keywords"` // Deprecated: use Tags instead
	Deprecated   bool     `json:"deprecated"`
}

// GenerationResult contains information about the generation process
type GenerationResult struct {
	IconsGenerated int           `json:"icons_generated"`
	FilesCreated   []string      `json:"files_created"`
	Duration       time.Duration `json:"duration"`
}

// SVGElement represents the parsed SVG structure
type SVGElement struct {
	ViewBox string `xml:"viewBox,attr"`
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
	Content string `xml:",innerxml"`
}

// IconMetadata represents the JSON metadata for a Lucide icon
type IconMetadata struct {
	Schema       string   `json:"$schema"`
	Contributors []string `json:"contributors"`
	Tags         []string `json:"tags"`
}


func readIconNamesFromRegistryFile(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	matcher := regexp.MustCompile(`(?m)^\s*[A-Za-z_][A-Za-z0-9_]*\s+IconName\s+=\s+"([^"]+)"\s*$`)
	matches := matcher.FindAllStringSubmatch(string(content), -1)

	names := make([]string, 0, len(matches))
	seen := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		name := match[1]
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}

	if len(names) == 0 && strings.Contains(string(content), "IconName =") {
		return nil, fmt.Errorf("registry parse error in %s: found IconName tokens but no valid constants", path)
	}

	return names, nil
}

// fetchLucideIcons retrieves icon data from the Lucide GitHub repository via git clone
func fetchLucideIcons(verbose bool, includeMetadata bool) ([]IconData, error) {
	if verbose {
		fmt.Println("Cloning Lucide repository...")
	}

	// Create temporary directory for git clone
	tempDir, err := os.MkdirTemp("", "lucide-clone-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Clone the repository at specific release tag (shallow clone for speed)
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", "1.14.0", "https://github.com/lucide-icons/lucide.git", tempDir)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	iconsDir := filepath.Join(tempDir, "icons")

	// Read all SVG files
	files, err := os.ReadDir(iconsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read icons directory: %w", err)
	}

	var icons []IconData
	svgCount := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".svg") {
			svgCount++
		}
	}

	if verbose {
		fmt.Printf("Found %d icons to process\n", svgCount)
	}

	processed := 0
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".svg") {
			continue
		}

		processed++
		if verbose && processed%50 == 0 {
			fmt.Printf("Processing icon %d/%d...\n", processed, svgCount)
		}

		iconName := strings.TrimSuffix(file.Name(), ".svg")
		svgPath := filepath.Join(iconsDir, file.Name())

		// Parse SVG and optionally JSON metadata from local files
		iconData, err := parseLocalIcon(svgPath, iconName, iconsDir, includeMetadata)
		if err != nil {
			if verbose {
				fmt.Printf("Warning: failed to process %s: %v\n", iconName, err)
			}
			continue
		}

		icons = append(icons, *iconData)
	}

	if verbose {
		fmt.Printf("Successfully processed %d icons\n", len(icons))
	}

	return icons, nil
}

// parseLocalIcon parses an SVG file and optionally its JSON metadata from local files
func parseLocalIcon(svgPath, iconName, iconsDir string, includeMetadata bool) (*IconData, error) {
	// Read SVG file
	svgData, err := os.ReadFile(svgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SVG file: %w", err)
	}

	var svg SVGElement
	if err := xml.Unmarshal(svgData, &svg); err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Read JSON metadata if requested
	var metadata *IconMetadata
	if includeMetadata {
		jsonPath := filepath.Join(iconsDir, iconName+".json")
		if jsonData, err := os.ReadFile(jsonPath); err == nil {
			metadata = &IconMetadata{}
			if err := json.Unmarshal(jsonData, metadata); err != nil {
				// If JSON parsing fails, use empty metadata
				metadata = &IconMetadata{}
			}
		} else {
			// If JSON file doesn't exist, use empty metadata
			metadata = &IconMetadata{}
		}
	} else {
		metadata = &IconMetadata{}
	}

	return &IconData{
		Name:         iconName,
		FuncName:     toFunctionName(iconName),
		ViewBox:      svg.ViewBox,
		Content:      strings.TrimSpace(svg.Content),
		Tags:         metadata.Tags,
		Contributors: metadata.Contributors,
	}, nil
}

func normalizeRequestedIconSet(requested []string) map[string]struct{} {
	requestedSet := make(map[string]struct{}, len(requested))
	for _, name := range requested {
		normalized := strings.ToLower(strings.TrimSpace(name))
		if normalized == "" {
			continue
		}
		requestedSet[normalized] = struct{}{}
	}
	return requestedSet
}

// findUnknownRequestedIcons returns a sorted list of requested icon names not present in icons.
func findUnknownRequestedIcons(icons []IconData, requested []string) []string {
	requestedSet := normalizeRequestedIconSet(requested)
	if len(requestedSet) == 0 {
		return nil
	}

	available := make(map[string]struct{}, len(icons))
	for _, icon := range icons {
		available[icon.Name] = struct{}{}
	}

	unknown := make([]string, 0)
	for name := range requestedSet {
		if _, ok := available[name]; !ok {
			unknown = append(unknown, name)
		}
	}
	sort.Strings(unknown)
	return unknown
}

// filterIconsByRequestedNames filters icons to only include explicitly requested names.
func filterIconsByRequestedNames(icons []IconData, requested []string) []IconData {
	requestedSet := normalizeRequestedIconSet(requested)
	if len(requestedSet) == 0 {
		return icons
	}

	filtered := make([]IconData, 0, len(icons))
	for _, icon := range icons {
		if _, ok := requestedSet[icon.Name]; ok {
			filtered = append(filtered, icon)
		}
	}

	return filtered
}

// toFunctionName converts an icon name to a valid Go function name
func toFunctionName(name string) string {
	// Convert kebab-case to PascalCase
	parts := strings.Split(name, "-")
	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			// Capitalize first letter, keep rest as-is
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}

	// Ensure it starts with a letter
	funcName := result.String()
	if len(funcName) > 0 && !isLetter(rune(funcName[0])) {
		funcName = "Icon" + funcName
	}

	return funcName
}

// isLetter checks if a rune is a letter
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// toConstantName converts an icon name to a constant name
func toConstantName(name, prefix string) string {
	funcName := toFunctionName(name)
	if prefix != "" {
		return prefix + funcName
	}
	return "Icon" + funcName
}

// generateFiles creates all the template files
func generateFiles(icons []IconData, config Config) ([]string, error) {
	var createdFiles []string

	// Generate main icons file
	iconsFile := filepath.Join(config.OutputDir, "icons.templ")
	if err := generateIconsFile(icons, config, iconsFile); err != nil {
		return nil, fmt.Errorf("failed to generate icons file: %w", err)
	}
	createdFiles = append(createdFiles, iconsFile)

	// Generate registry file
	registryFile := filepath.Join(config.OutputDir, "registry.templ")
	if err := generateRegistryFile(icons, config, registryFile); err != nil {
		return nil, fmt.Errorf("failed to generate registry file: %w", err)
	}
	createdFiles = append(createdFiles, registryFile)

	return createdFiles, nil
}
