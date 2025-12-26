package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/KOU050223/specs-viewer/internal/parser"
	"github.com/KOU050223/specs-viewer/internal/watcher"
)

type Server struct {
	port      int
	specPaths []string
	templates embed.FS
	watcher   *watcher.FileWatcher
	upgrader  websocket.Upgrader
}

func New(port int, specPaths []string, templates embed.FS, fw *watcher.FileWatcher) *Server {
	return &Server{
		port:      port,
		specPaths: specPaths,
		templates: templates,
		watcher:   fw,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/tree", s.handleTree)
	http.HandleFunc("/api/file", s.handleFile)
	http.HandleFunc("/ws", s.handleWebSocket)

	addr := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(s.templates, "web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Specs Viewer",
		"Paths": s.specPaths,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleTree(w http.ResponseWriter, r *http.Request) {
	trees, err := parser.ParseMultipleDirectories(s.specPaths)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trees)
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	// Security check: ensure the file is within one of the spec directories
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	allowed := false
	for _, specPath := range s.specPaths {
		absSpecPath, err := filepath.Abs(specPath)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absSpecPath) {
			allowed = true
			break
		}
	}

	if !allowed {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	file, err := parser.GetFileContent(absPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(file)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Subscribe to file changes
	changes := s.watcher.Subscribe()
	defer s.watcher.Unsubscribe(changes)

	// Send initial connection message
	if err := conn.WriteJSON(map[string]string{
		"type":    "connected",
		"message": "WebSocket connected",
	}); err != nil {
		log.Printf("Write error: %v", err)
		return
	}

	// Listen for file changes and send updates
	for {
		select {
		case changedPath, ok := <-changes:
			if !ok {
				return
			}

			// Re-parse the changed file
			file, err := parser.GetFileContent(changedPath)
			if err != nil {
				log.Printf("Error parsing changed file: %v", err)
				continue
			}

			// Send update to client
			if err := conn.WriteJSON(map[string]interface{}{
				"type": "file_changed",
				"file": file,
			}); err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		}
	}
}
