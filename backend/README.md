# KubeVision Backend

High-performance Docker monitoring backend built with Go.

## Features

- Real-time container statistics via WebSocket
- Container logs streaming
- RESTful API for container management
- Zero-dependency architecture
- High-performance streaming with backpressure control

## Prerequisites

- Go 1.21 or higher
- Docker daemon running (with access to Docker socket)

## Configuration

Copy `.env.example` to `.env` and configure:

```env
DOCKER_HOST=unix:///var/run/docker.sock
PORT=8080
LOG_LEVEL=info
AUTH_ENABLED=false
AUTH_TOKEN=your-secret-token
```

## Running

```bash
# Install dependencies
go mod download

# Run server
go run cmd/server/main.go
```

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/containers` - List all containers
- `GET /api/containers/:id` - Get container details
- `WS /ws/stats/:id` - WebSocket for container stats
- `WS /ws/logs/:id` - WebSocket for container logs

## Development

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...
```

