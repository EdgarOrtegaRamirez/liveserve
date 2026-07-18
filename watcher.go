package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// startWatcher watches the served directory for file changes.
// It sends a signal on reloadCh whenever a file is created, modified, or deleted.
func startWatcher(dir string, reloadCh chan<- struct{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error creating file watcher: %v (live reload disabled)", err)
		return
	}
	defer watcher.Close()

	// Watch the root directory and all subdirectories
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // skip inaccessible paths
		}
		if fi.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking directory tree: %v (live reload disabled)", err)
		return
	}

	// Debounce timer
	var debounce *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Ignore hidden files and temp files
			base := filepath.Base(event.Name)
			if len(base) > 0 && base[0] == '.' {
				continue
			}
			// Ignore vim/nano temp files
			if base[len(base)-1] == '~' {
				continue
			}

			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) != 0 {
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(100*time.Millisecond, func() {
					select {
					case reloadCh <- struct{}{}:
					default:
					}
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
