package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KOU050223/specs-viewer/internal/server"
	"github.com/KOU050223/specs-viewer/internal/watcher"
)

//go:embed web/templates/*
var templates embed.FS

func main() {
	port := flag.Int("port", 4829, "Port to run the server on")
	flag.Parse()

	args := flag.Args()

	var specPaths []string
	if len(args) < 1 {
		// 引数なしの場合、標準的なspec駆動開発のディレクトリを探す
		candidates := []string{"specs", ".specify"}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				absPath, err := filepath.Abs(candidate)
				if err != nil {
					log.Fatalf("Failed to get absolute path for %s: %v", candidate, err)
				}
				specPaths = append(specPaths, absPath)
				fmt.Printf("📂 Auto-detected spec directory: %s\n", candidate)
			}
		}

		if len(specPaths) == 0 {
			fmt.Fprintln(os.Stderr, "Error: No spec directory found")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Usage: specs-viewer [path-to-spec-directory...]")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "If no path is provided, specs-viewer will look for:")
			fmt.Fprintln(os.Stderr, "  - ./specs")
			fmt.Fprintln(os.Stderr, "  - ./.specify")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Example: specs-viewer ./my-specs")
			fmt.Fprintln(os.Stderr, "Example (multiple): specs-viewer ./specs ./.specify")
			os.Exit(1)
		}
	} else {
		// 引数で指定された全てのパスを追加
		for _, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				log.Fatalf("Failed to get absolute path for %s: %v", arg, err)
			}

			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				log.Fatalf("Path does not exist: %s", absPath)
			}

			specPaths = append(specPaths, absPath)
		}
	}

	// Create file watcher for all paths
	fw, err := watcher.NewMulti(specPaths)
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	defer fw.Close()

	// Start server
	srv := server.New(*port, specPaths, templates, fw)

	fmt.Printf("🚀 Specs Viewer running at http://localhost:%d\n", *port)
	fmt.Printf("📂 Watching %d director%s:\n", len(specPaths), map[bool]string{true: "y", false: "ies"}[len(specPaths) == 1])
	for _, path := range specPaths {
		fmt.Printf("   - %s\n", path)
	}
	fmt.Println("Press Ctrl+C to stop")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
