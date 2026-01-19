# KubeVision

A high-performance, zero-dependency lightweight Docker monitoring solution that provides real-time container metrics and logs visualization through WebSocket streaming.

## 📚 Tech Stack

### Backend
- **Go 1.24+** with Gin framework
- **Docker SDK v28+** for Docker API integration
- **WebSocket** (gorilla/websocket) for real-time streaming
- **Zap** for structured logging
- **Viper** for configuration management
- **Embedded Frontend** using `go:embed` for single binary deployment

### Frontend
- **React 18** with Functional Components
- **TypeScript** (Strict Mode)
- **Vite** as build tool
- **Tailwind CSS** for styling
- **Zustand** for state management
- **ECharts** with incremental updates for metrics visualization
- **xterm.js** with FitAddon and SearchAddon for terminal emulation
- **WebSocket** client for real-time updates

## 🚀 Installation Guide

Choose the installation method that best fits your needs:

| Method | Difficulty | Speed | Best For |
|--------|-----------|-------|----------|
| Method 1: Docker Compose | ⭐ Easy | ⚡ Fast | Beginners, Production |
| Method 2: Manual Docker Run | ⭐⭐ Medium | ⚡⚡ Medium | Full control, Custom config |
| Method 3: Manual Build | ⭐⭐⭐ Hard | ⚡⚡⚡ Slow | Developers, Debugging |

### Method 1: Docker Compose (Recommended ⭐)

**Best for:** Beginners, production deployments, managing complex configurations

#### Prerequisites
- Docker Desktop (Windows/Mac) or Docker Engine (Linux)
- Docker Compose v2.0+ (included with Docker Desktop)

#### Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd kubevision
   ```

2. **Build and run**
   ```bash
   docker compose -f docker/docker-compose.yml up --build
   ```
   
   Or run in detached mode:
   ```bash
   docker compose -f docker/docker-compose.yml up -d --build
   ```

3. **Check status**
   ```bash
   docker compose -f docker/docker-compose.yml ps
   ```

4. **Access the application**
   - Frontend: http://localhost:8080
   - API: http://localhost:8080/api/health
   - Containers: http://localhost:8080/api/containers

**Note for Windows users:** Use `docker compose` (with space) instead of `docker-compose` (with hyphen).

### Method 2: Manual Docker Run Commands

**Best for:** Maximum control, custom configurations, understanding how containers work

#### Prerequisites
- Docker Desktop or Docker Engine

#### Steps

```bash
# 1. Build the image
docker build -f docker/Dockerfile -t kubevision:latest .

# 2. Create network (optional)
docker network create kubevision-network

# 3. Run container
docker run -d \
  --name kubevision \
  --network kubevision-network \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e LOG_LEVEL=info \
  -e AUTH_ENABLED=false \
  -e CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000 \
  --restart unless-stopped \
  kubevision:latest
```

**Windows Docker Desktop:**
```powershell
docker run -d `
  --name kubevision `
  -p 8080:8080 `
  -v \\.\pipe\docker_engine:\\.\pipe\docker_engine `
  -e PORT=8080 `
  -e LOG_LEVEL=info `
  --restart unless-stopped `
  kubevision:latest
```

### Method 3: Manual Build & Run (For Developers)

**Best for:** Developers who want to modify code, debug, or understand the application structure

#### Prerequisites
- Go 1.24+
- Node.js 18+ and npm
- Docker Desktop or Docker Engine (for container monitoring)

#### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables (optional)**
   ```bash
   # Create .env file (optional - defaults work for most cases)
   export PORT=8080
   export HOST=0.0.0.0
   export LOG_LEVEL=info
   export AUTH_ENABLED=false
   export CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
   ```
   
   **Note:** `DOCKER_HOST` is optional. The Docker client will auto-detect:
   - Windows: `npipe:////./pipe/docker_engine`
   - Linux/Mac: `unix:///var/run/docker.sock`

4. **Run server**
   ```bash
   go run ./cmd/server
   ```
   
   Or build and run:
   ```bash
   go build ./cmd/server
   ./server        # Linux/Mac
   .\server.exe    # Windows
   ```

   Backend starts on http://localhost:8080

#### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   # Or if peer dependency conflicts occur:
   npm install --legacy-peer-deps
   ```

3. **Configure API URL (optional - defaults to http://localhost:8080)**
   ```bash
   echo "VITE_API_URL=http://localhost:8080" > .env
   ```

4. **Start development server**
   ```bash
   npm run dev
   ```

   Frontend starts on http://localhost:5173

## 🧪 Testing

### Backend Tests

Run all Go tests:
```bash
cd backend
go test ./... -v
```

Run tests for a specific package:
```bash
go test ./internal/docker -v
go test ./internal/api -v
go test ./internal/middleware -v
```

**Test Coverage:**
- `internal/docker/metrics_test.go` - Metrics calculation and CPU/memory stats
- `internal/api/containers_test.go` - Container handler and name extraction
- `internal/middleware/cors_test.go` - CORS middleware functionality

### Frontend Tests

Run all frontend tests:
```bash
cd frontend
npm test
```

Run tests in watch mode:
```bash
npm test -- --watch
```

Run tests with UI:
```bash
npm run test:ui
```

Run tests with coverage:
```bash
npm run test:coverage
```

**Test Coverage:**
- `src/services/api.test.ts` - API service methods
- `src/hooks/useWebSocket.test.ts` - WebSocket hook behavior
- `src/stores/containerStore.test.ts` - Zustand store operations

## ✨ Key Features

### Real-time Monitoring
- **WebSocket-based streaming** for container statistics
- **CPU and Memory metrics** with real-time updates
- **Network and I/O statistics** tracking
- **Optimized for 20+ containers** simultaneously

### Visualization
- **Beautiful charts** with ECharts for CPU and memory metrics
- **Incremental updates** for performance optimization
- **Container status indicators** with color coding
- **Responsive design** with Tailwind CSS

### Log Streaming
- **Real-time container logs** with xterm.js terminal emulation
- **Log search functionality** with SearchAddon
- **Persistent logs** - logs persist when scrolling (no data loss)
- **Auto-scroll toggle** for following new logs
- **Clear and search** controls

### Container Control
- **Start, stop, restart** containers
- **Pause and unpause** operations
- **Token-based authentication** for container operations (optional)
- **Protected API endpoints**

### Performance
- Browser CPU usage: <15% (20 containers)
- Browser memory: <200MB
- Backend memory: <512MB (20 containers)
- WebSocket latency: <100ms
- API response time: <200ms (p95)

## 📡 API Endpoints

### REST API

- `GET /api/health` - Health check
- `GET /api/containers` - List all containers
- `GET /api/containers/:id` - Get container details
- `POST /api/containers/:id/start` - Start container (requires auth)
- `POST /api/containers/:id/stop` - Stop container (requires auth)
- `POST /api/containers/:id/restart` - Restart container (requires auth)
- `POST /api/containers/:id/pause` - Pause container (requires auth)
- `POST /api/containers/:id/unpause` - Unpause container (requires auth)

### WebSocket

- `WS /ws/stats/:id` - Real-time container statistics
- `WS /ws/logs/:id` - Real-time container logs
  - Query params: `follow=true`, `tail=100`, `since=timestamp`

## 🔒 Security Features

- **Non-root Docker containers**: Application runs as kubevision user (when not accessing Docker socket)
- **Secure base images**: Uses Alpine Linux for minimal attack surface
- **Environment variables**: All secrets externalized
- **Read-only Docker socket**: Socket mounted as read-only (`:ro`)
- **Internal networking**: Optional Docker network isolation
- **Health checks**: Built-in health monitoring endpoints
- **CORS configuration**: Configurable allowed origins
- **Token-based authentication**: Optional JWT-like token auth for container operations
- **Security options**: `no-new-privileges:true` in Docker Compose

## 🛠️ Configuration

### Backend Environment Variables

Environment variables can be set via `.env` file or environment:

```env
# Optional: Docker connection (auto-detected if not set)
# Windows: npipe:////./pipe/docker_engine
# Linux/Mac: unix:///var/run/docker.sock
DOCKER_HOST=

PORT=8080
HOST=0.0.0.0
LOG_LEVEL=info
AUTH_ENABLED=false
AUTH_TOKEN=your-secret-token-here
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

### Frontend Environment Variables

Create `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080
```

