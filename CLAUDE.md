# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A web-based serial port terminal server. Users connect to local serial ports (e.g., `/dev/ttyUSB0`, COM ports) through a browser-based xterm.js terminal. The Go backend bridges WebSocket connections to physical serial ports, and a Vue 3 SPA serves as the frontend.

## Build & Run Commands

### Go backend
```sh
# Build (frontend dist/ must exist first — it's committed to git)
go build -o remote-serial-port-server

# Build with version info (as used in CI/Dockerfile)
go build -trimpath \
  -ldflags="-X 'github.com/iceking2nd/remote-serial-port-server/global.Version=<VERSION>' \
            -X 'github.com/iceking2nd/remote-serial-port-server/global.BuildTime=<DATE>' \
            -X github.com/iceking2nd/remote-serial-port-server/global.GitCommit=<SHA>" \
  -o remote-serial-port-server

# Run (default: 127.0.0.1 on random port)
./remote-serial-port-server -l 0.0.0.0 -p 8192

# With debug logging
./remote-serial-port-server --log-level 5 --log-file /tmp/rsps.log
```

### Frontend (static/)
```sh
cd static
npm install
npm run dev       # Dev server with hot reload (proxies to localhost:8192)
npm run build     # Production build → static/dist/
npm run lint      # ESLint + Prettier check
```

### Docker
```sh
docker build --build-arg VERSION=x.y.z \
  --build-arg BUILD_DATE="$(date +'%Y-%m-%d %H:%M:%S')" \
  --build-arg SHA="$(git rev-parse HEAD)" \
  -t remote-serial-port-server .
```

There are no Go tests (`*_test.go`) in this project.

## Architecture

### Request Flow
```
Browser (xterm.js) → WebSocket /api/port/open?port=...&baudrate=... → Gorilla WebSocket → go.bug.st/serial → Physical serial port
```

### Backend Structure (Gin HTTP framework)
- **`cmd/root.go`** — Cobra CLI entry point; parses flags, initializes logrus logger, generates UUID API key (`global.APIKey`), creates Gin server, registers routes, handles graceful shutdown on SIGINT/SIGTERM
- **`app/routers/setup.go`** — Registers three route groups: `/debug/*` (pprof), `/api/*` (core API), `/` (static SPA)
- **`app/controllers/PortController/open.go`** — Core WebSocket↔serial bridge: upgrades HTTP to WebSocket, opens serial port with query params, spawns goroutine for serial→WebSocket, main loop for WebSocket→serial
- **`app/controllers/PortController/list.go`** — Lists available serial ports via `serial.GetPortsList()`
- **`app/controllers/SystemController/key.go`** — Returns the server-generated API key (no auth required)
- **`app/middlewares/key.go`** — Checks `X-API-Key` header or `?key=` query param against `global.APIKey`
- **`global/vars.go`** — Package-level variables: `Version`, `BuildTime`, `GitCommit`, `Log`, `APIKey`

### Frontend Structure (Vue 3 + Vite in `static/`)
- **`static/src/App.vue`** — Single-page terminal UI with Element Plus controls for serial port configuration, xterm.js terminal, and i18n (en-US, zh-CN)
- **`static/static.go`** — Go `embed` directives; `RootFS` and `AssetsFS` embed `dist/` into the Go binary

### API Key Model
The API key is a UUID generated at startup and exposed publicly via `/api/system/key`. It acts as a session token rather than a true auth mechanism — anyone who can reach the server can retrieve it. It changes on every restart.

### Static Asset Embedding
`static/dist/` is committed to git and embedded into the Go binary at build time via `//go:embed`. This allows building the Go binary without the Node.js toolchain if `dist/` already exists.

### CLI Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--listen-address` | `-l` | `127.0.0.1` | Bind address |
| `--listen-port` | `-p` | `0` | Port (0 = random) |
| `--log-file` | | `""` | Log file path (empty = stdout) |
| `--log-level` | | `3` | 0–6 (3=Warn, 5=Debug) |

## Key Dependencies
- **Gin** (`github.com/gin-gonic/gin`) — HTTP framework
- **Gorilla WebSocket** (`github.com/gorilla/websocket`) — WebSocket upgrade and I/O
- **go.bug.st/serial** — Cross-platform serial port access
- **Cobra** (`github.com/spf13/cobra`) — CLI framework
- **Logrus** (`github.com/sirupsen/logrus`) — Structured logging
- **xterm.js** + **Element Plus** — Frontend terminal and UI components
