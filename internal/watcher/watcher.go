package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher   *fsnotify.Watcher
	listeners []chan string
	mu        sync.RWMutex
}

func New(path string) (*FileWatcher, error) {
	return NewMulti([]string{path})
}

func NewMulti(paths []string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher:   watcher,
		listeners: make([]chan string, 0),
	}

	for _, path := range paths {
		if err := fw.addRecursive(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	go fw.watch()

	return fw, nil
}

func (fw *FileWatcher) addRecursive(path string) error {
	err := filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil {
			return nil
		}

		// Skip hidden directories
		if strings.HasPrefix(filepath.Base(walkPath), ".") && walkPath != path {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Add directory to watcher
		if info.IsDir() {
			return fw.watcher.Add(walkPath)
		}

		return nil
	})

	return err
}

func (fw *FileWatcher) watch() {
	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Only notify on write or create events for .md files
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				if strings.HasSuffix(event.Name, ".md") {
					log.Printf("File changed: %s", event.Name)
					fw.notify(event.Name)
				}
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (fw *FileWatcher) notify(path string) {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	for _, listener := range fw.listeners {
		select {
		case listener <- path:
		default:
			// Skip if channel is full
		}
	}
}

func (fw *FileWatcher) Subscribe() chan string {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	ch := make(chan string, 10)
	fw.listeners = append(fw.listeners, ch)
	return ch
}

func (fw *FileWatcher) Unsubscribe(ch chan string) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	for i, listener := range fw.listeners {
		if listener == ch {
			fw.listeners = append(fw.listeners[:i], fw.listeners[i+1:]...)
			close(ch)
			break
		}
	}
}

func (fw *FileWatcher) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	for _, listener := range fw.listeners {
		close(listener)
	}
	fw.listeners = nil

	return fw.watcher.Close()
}
