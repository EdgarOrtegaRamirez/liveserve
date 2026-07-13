# Liveserve

A minimal, zero-dependency static file server with live reload, written in Go.

```bash
liveserve                  # Serve current directory on :8080
liveserve --port 3000      # Custom port
liveserve --dir ./public   # Custom directory
liveserve --no-reload      # Disable live reload
```

## Features

- **Static file serving** — serves any directory over HTTP with automatic `index.html` support
- **Live reload** — injects a WebSocket script into HTML pages; when files change, connected browsers auto-refresh
- **File watching** — monitors the entire directory tree (recursive) with debounced change detection
- **Zero config** — just run `liveserve` in any directory

## Install

```bash
# From source
go install github.com/EdgarOrtegaRamirez/liveserve@latest

# Or build from repo
git clone https://github.com/EdgarOrtegaRamirez/liveserve.git
cd liveserve
go build -o liveserve .
sudo mv liveserve /usr/local/bin/
```

## Usage

```
Usage: liveserve [--port N] [--dir PATH] [--no-reload]

Options:
  --port int       HTTP server port (default 8080)
  --dir string     Directory to serve (default ".")
  --no-reload      Disable live reload (default false)
  --help           Show this help message
```

### Examples

```bash
# Serve current directory
cd my-project
liveserve

# Custom port
liveserve --port 3000

# Serve a specific directory
liveserve --dir ./build

# Production mode (no live reload)
liveserve --no-reload
```

## How Live Reload Works

1. `liveserve` starts an HTTP server and a WebSocket server at `/__livereload`
2. When serving `.html` files, it automatically injects a small JavaScript snippet before `</body>`
3. The JavaScript opens a WebSocket connection back to the server
4. When a file change is detected (via `fsnotify`), the server broadcasts a `"reload"` message to all connected clients
5. Clients call `location.reload()` upon receiving the message

Changes are debounced (100ms) to avoid multiple reloads during rapid saves.

![CI](https://github.com/EdgarOrtegaRamirez/liveserve/actions/workflows/ci.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)