# Remote Serial Port Server

A web-based serial port terminal server. Connect to local serial ports (e.g., `/dev/ttyUSB0`, COM ports) through a browser-based terminal powered by xterm.js. The Go backend bridges WebSocket connections to physical serial ports, allowing multiple clients to share a single serial port session.

## Features

- Browser-based serial port terminal — no client software needed
- Multiple WebSocket clients can share a single serial port (pub/sub model)
- Configurable serial port parameters: baud rate, data bits, parity, stop bits, RTS/DTR
- Multi-language UI (English, 简体中文, 繁體中文, Deutsch, Français, Русский, Polski)
- Single binary deployment — frontend assets are embedded via `go:embed`
- Docker support with multi-platform images

## Quick Start

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/iceking2nd/remote-serial-port-server/releases).

### Docker

```sh
docker run -d --name rsps \
  --device=/dev/ttyUSB0 \
  -p 8192:8192 \
  iceking2nd/remote-serial-port-server:latest
```

### Build from Source

Prerequisites: Go 1.25+, Node.js (for frontend build)

```sh
# Build frontend first
cd static && npm install && npm run build && cd ..

# Build Go binary
go build -o remote-serial-port-server

# Run
./remote-serial-port-server -l 0.0.0.0 -p 8192
```

Then open `http://localhost:8192` in your browser.

### Build with Docker

```sh
docker build \
  --build-arg VERSION=$(git describe --tags --always) \
  --build-arg BUILD_DATE="$(date +'%Y-%m-%d %H:%M:%S')" \
  --build-arg SHA="$(git rev-parse HEAD)" \
  -t remote-serial-port-server .
```

## Usage

```
Usage:
  remote-serial-port-server [flags]

Flags:
  -l, --listen-address string   Bind address (default "127.0.0.1")
  -p, --listen-port int         Listen port, 0 = random (default 0)
      --log-file string         Log file path (empty = stdout)
      --log-level int           Log level 0-6 (default 3, 3=Warn, 5=Debug)
```

## How It Works

```
Browser (xterm.js)
    │
    ▼ WebSocket /api/port/open?port=...&baudrate=...
Go Backend (Gin + Gorilla WebSocket)
    │
    ▼ go.bug.st/serial
Physical Serial Port
```

When multiple clients connect to the same serial port, they share a single `PortSession`. Data read from the port is broadcast to all subscribers; when the last subscriber disconnects, the serial port is closed automatically.

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/system/key` | None | Get the server-generated API key |
| GET | `/api/port/` | API Key | List available serial ports |
| GET | `/api/port/open` | API Key | Open WebSocket↔serial bridge |

The API key is a UUID generated at startup and available via `/api/system/key`. It acts as a session token and changes on every server restart.

## Development

### Frontend Dev Server

```sh
cd static
npm install
npm run dev    # Hot-reload dev server, proxies API to localhost:8192
```

### Run Tests

```sh
go test ./...
go test -v ./app/services/            # Single package
go test -v -run TestPortSessionBroadcast ./app/services/  # Single test
```

### Frontend Lint

```sh
cd static && npm run lint
```

## License

[MIT](LICENSE) © 2025 Daniel Wu
