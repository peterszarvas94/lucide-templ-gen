package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	lucidegen "github.com/peterszarvas94/lucide-templ-gen"
)

func main() {
	// Generate icons programmatically
	all, err := lucidegen.ListAvailableIconNames()
	if err != nil {
		log.Fatalf("Failed to load icon list: %v", err)
	}

	config := lucidegen.Config{OutputDir: "./components", PackageName: "components"}
	result, err := lucidegen.GenerateFromIconNames(config, all)
	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}

	fmt.Printf("✅ Generated %d icons in %v\n", result.IconsGenerated, result.Duration)
	fmt.Printf("📁 Files created: %v\n", result.FilesCreated)

	// Start a demo web server
	fmt.Println("🌐 Starting demo server on http://localhost:8080")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/icons", iconsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	component := HomePage()
	component.Render(context.Background(), w)
}

func iconsHandler(w http.ResponseWriter, r *http.Request) {
	component := IconsPage()
	component.Render(context.Background(), w)
}
