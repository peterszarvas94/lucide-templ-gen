package main

import (
	"log"

	lucidegen "github.com/peterszarvas94/lucide-templ-gen"
)

func main() {
	// Basic usage - generate all icons
	config := lucidegen.Config{
		OutputDir:   "./icons",
		PackageName: "icons",
		Verbose:     true,
	}

	result, err := lucidegen.Generate(config)
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	log.Printf("Generated %d icons in %v", result.IconsGenerated, result.Duration)
	log.Printf("Categories: %v", result.Categories)
}
