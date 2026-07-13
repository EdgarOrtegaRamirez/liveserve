package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	dir := flag.String("dir", ".", "Directory to serve")
	noReload := flag.Bool("no-reload", false, "Disable live reload")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("Invalid directory: %v", err)
	}

	// Verify directory exists
	if fi, err := os.Stat(absDir); err != nil || !fi.IsDir() {
		log.Fatalf("Not a valid directory: %s", absDir)
	}

	// File watcher for live reload
	reloadCh := make(chan struct{}, 1)
	if !*noReload {
		go startWatcher(absDir, reloadCh)
	}

	// WebSocket hub for broadcasting reload signals
	hub := newHub()
	go hub.run()

	// Watch for file changes and broadcast reload
	if !*noReload {
		go func() {
			for range reloadCh {
				hub.broadcast()
			}
		}()
	}

	// Static file handler with optional live reload injection
	fileServer := http.FileServer(http.Dir(absDir))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !*noReload && r.URL.Path != "/__livereload" {
			// Try to serve the file first
			serveFile := filepath.Join(absDir, r.URL.Path)
			if fi, err := os.Stat(serveFile); err == nil && !fi.IsDir() {
				// Check if it's an HTML file for injection
				ext := filepath.Ext(serveFile)
				if ext == ".html" || ext == ".htm" {
					injectLiveReload(w, r, serveFile)
					return
				}
			}
		}
		fileServer.ServeHTTP(w, r)
	})

	// WebSocket endpoint
	if !*noReload {
		http.Handle("/__livereload", websocketHandler(hub))
	}
	http.Handle("/", handler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("🌐 Liveserve running at http://localhost%s\n", addr)
	fmt.Printf("   Serving directory: %s\n", absDir)
	if !*noReload {
		fmt.Println("   Live reload: enabled (modify files to trigger refresh)")
	} else {
		fmt.Println("   Live reload: disabled")
	}
	fmt.Println()
	log.Fatal(http.ListenAndServe(addr, nil))
}