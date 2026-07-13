# Test the Hub (WebSocket connection manager)

Run: `go test ./...`

## Test Strategy

- **hub_test.go** — Tests for WebSocket hub registration, unregistration, and broadcasting
- **inject_test.go** — Tests for HTML live reload script injection
- **watcher_test.go** — Tests for the file system watcher (uses temp directories)