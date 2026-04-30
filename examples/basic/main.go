package main

import (
	"log"

	lucidegen "github.com/peterszarvas94/lucide-templ-gen/v2"
)

func main() {
	// Basic usage - generate all icons
	all, err := lucidegen.ListAvailableIconNames()
	if err != nil {
		log.Fatalf("Failed to load icon list: %v", err)
	}

	config := lucidegen.Config{OutputDir: "./icons", PackageName: "icons"}
	result, err := lucidegen.GenerateFromIconNames(config, all)
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	log.Printf("Generated %d icons in %v", result.IconsGenerated, result.Duration)
}
