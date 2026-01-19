# Changelog

All notable changes to KubeVision will be documented in this file.

## [0.1.0] - 2024-01-XX

### Added
- Initial release of KubeVision
- Real-time container statistics via WebSocket
- Container metrics visualization with ECharts
- Container logs streaming with xterm.js
- Container control operations (start, stop, restart, pause, unpause)
- RESTful API for container management
- WebSocket endpoints for stats and logs
- Token-based authentication for container operations
- CORS support
- Docker Compose deployment
- Multi-stage Dockerfile for production builds
- CI/CD pipeline with GitHub Actions
- Comprehensive documentation

### Features
- **Backend (Go)**
  - Docker client with singleton pattern
  - CPU and memory metrics calculation
  - WebSocket streaming with backpressure control
  - Graceful shutdown handling
  - Structured logging with zap
  - Health check endpoint

- **Frontend (React + TypeScript)**
  - Real-time metrics charts
  - Container list with status indicators
  - Log terminal with search functionality
  - Intersection Observer for performance optimization
  - WebSocket auto-reconnection with exponential backoff
  - Zustand state management
  - TailwindCSS styling

### Security
- Token-based authentication
- Read-only Docker socket mount
- Non-root user in Docker container
- CORS configuration

### Performance
- Optimized for 20+ containers
- Virtual scrolling support
- Data retention limits (60 data points)
- Efficient WebSocket message handling

