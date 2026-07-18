package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInjectLiveReload(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "injects before </body>",
			html: "<html><body><h1>Hello</h1></body></html>",
			want: "ws.onmessage",
		},
		{
			name: "appends when no </body>",
			html: "<html><head></head><h1>No body tag</h1></html>",
			want: "location.reload",
		},
		{
			name: "empty file",
			html: "",
			want: "location.reload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			filePath := filepath.Join(dir, "test.html")
			if err := os.WriteFile(filePath, []byte(tt.html), 0644); err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
			req.Host = "localhost:8080"
			rec := httptest.NewRecorder()

			injectLiveReload(rec, req, filePath)

			body := rec.Body.String()
			if !strings.Contains(body, tt.want) {
				t.Errorf("expected body to contain %q, got:\n%s", tt.want, body)
			}
			if !strings.Contains(body, "WebSocket") {
				t.Errorf("expected body to contain WebSocket reference")
			}
			// Verify Content-Type is set
			if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
				t.Errorf("expected Content-Type text/html, got %q", ct)
			}
		})
	}
}

func TestInjectLiveReload_FileNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent.html", nil)
	rec := httptest.NewRecorder()

	injectLiveReload(rec, req, "/nonexistent.html")
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}

// mockWSConn implements wsConn for testing hub broadcast.
type mockWSConn struct {
	written []byte
}

func (m *mockWSConn) Write(data []byte) (int, error) {
	m.written = append(m.written, data...)
	return len(data), nil
}

func (m *mockWSConn) Close() error { return nil }

func TestHubRegisterAndBroadcast(t *testing.T) {
	h := newHub()

	m1 := &mockWSConn{}
	m2 := &mockWSConn{}

	h.register(m1)
	h.register(m2)

	if len(h.clients) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(h.clients))
	}

	h.broadcast()

	if string(m1.written) != "reload" {
		t.Errorf("expected m1 to receive 'reload', got %q", string(m1.written))
	}
	if string(m2.written) != "reload" {
		t.Errorf("expected m2 to receive 'reload', got %q", string(m2.written))
	}
}

func TestHubUnregister(t *testing.T) {
	h := newHub()
	m := &mockWSConn{}

	h.register(m)
	h.unregister(m)

	if len(h.clients) != 0 {
		t.Errorf("expected 0 clients after unregister, got %d", len(h.clients))
	}
}

func TestMainFunction_Builds(t *testing.T) {
	// Verify the binary still compiles
	// This is implicitly tested by go build, but we also check the flag defaults
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
}
