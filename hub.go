package main

import (
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

// wsConn is the minimal interface needed for WebSocket write+close operations.
// Both *websocket.Conn and test mocks satisfy this interface.
type wsConn interface {
	io.Writer
	Close() error
}

// hub manages WebSocket connections for live reload broadcasting.
type hub struct {
	clients map[wsConn]bool
	mu      sync.RWMutex
}

func newHub() *hub {
	return &hub{
		clients: make(map[wsConn]bool),
	}
}

func (h *hub) register(ws wsConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[ws] = true
}

func (h *hub) unregister(ws wsConn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, ws)
	ws.Close()
}

func (h *hub) broadcast() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	msg := []byte("reload")
	for ws := range h.clients {
		if _, err := ws.Write(msg); err != nil {
			log.Printf("WebSocket write error: %v", err)
			go h.unregister(ws)
		}
	}
}

func (h *hub) run() {
	// hub.run is a no-op placeholder for future heartbeat/cleanup logic.
	// Currently, cleanup happens on write error.
}

// websocketHandler returns an HTTP handler that upgrades connections
// to WebSocket and registers them with the hub.
func websocketHandler(h *hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		websocket.Server{
			Handler: func(ws *websocket.Conn) {
				h.register(ws)
				// Keep connection alive — read in a loop
				buf := make([]byte, 256)
				for {
					_, err := ws.Read(buf)
					if err != nil {
						h.unregister(ws)
						return
					}
				}
			},
		}.ServeHTTP(w, r)
	}
}
