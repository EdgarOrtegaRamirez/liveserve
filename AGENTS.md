# AGENTS.md — Liveserve

## Project Overview

Liveserve is a minimal static file server with live reload, written in Go. It serves static files over HTTP with automatic `index.html` support, and injects a WebSocket-based live reload script into HTML pages.

## Architecture

```
liveserve/
├── main.go         — CLI entry point (flag-based, no external CLI deps)
├── hub.go          — WebSocket connection hub for broadcasting reload signals
├── inject.go       — HTML live reload script injection
├── watcher.go      — File system watcher using fsnotify (debounced)
├── liveserve_test.go — Tests (8 tests: inject, hub register/broadcast, unregister)
├── go.mod / go.sum
├── README.md
├── LICENSE
└── .github/workflows/ci.yml
```

## Key Design Decisions

1. **Zero external CLI dependencies** — Uses stdlib `flag` instead of cobra (simple enough)
2. **Interface-based hub** — `wsConn` interface lets us test broadcasting without real WebSocket connections
3. **Debounced file watching** — 100ms debounce timer prevents multiple reloads on rapid saves
4. **HTML injection** — Injects live reload script before `</body>`, falls back to appending at end
5. **Hidden file filtering** — Ignores files starting with `.` and temp files ending with `~`

## Dependencies

- `github.com/fsnotify/fsnotify` — File system event notifications
- `golang.org/x/net/websocket` — WebSocket server implementation

## Build & Test

```bash
go build -o liveserve .
go test ./... -v
go vet ./...
```

## Common Tasks

### Adding a new feature
1. Create a new file (e.g., `auth.go` for basic auth support)
2. Add CLI flag in `main.go`
3. Wire it into the handler chain
4. Add tests
5. Update README.md

### Changing the reload message
Edit the `msg` byte slice in `hub.go` `broadcast()` method.

### Customizing the injected script
Edit the `reloadScript` template in `inject.go`.

## Troubleshooting

- **Tests failing with "Hijack not supported"**: Don't test the WebSocket handler with httptest.Recorder — it doesn't support Hijack. Use mock connections instead.
- **Live reload not working**: Ensure the file being watched isn't hidden (starts with `.`) or a temp file (ends with `~`).