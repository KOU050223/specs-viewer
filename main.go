package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yourusername/specs-viewer/internal/server"
	"github.com/yourusername/specs-viewer/internal/watcher"
)

//go:embed web/templates/*
var templates embed.FS

func main() {
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: specs-viewer <path-to-spec-directory>")
		fmt.Fprintln(os.Stderr, "Example: specs-viewer ./specs")
		os.Exit(1)
	}

	specPath := args[0]
	absPath, err := filepath.Abs(specPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %s", absPath)
	}

	// Create file watcher
	fw, err := watcher.New(absPath)
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	defer fw.Close()

	// Start server
	srv := server.New(*port, absPath, templates, fw)

	fmt.Printf("🚀 Specs Viewer running at http://localhost:%d\n", *port)
	fmt.Printf("📂 Watching directory: %s\n", absPath)
	fmt.Println("Press Ctrl+C to stop")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
