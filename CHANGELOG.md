# Changelog

All notable changes to KubeVision will be documented in this file.

## [Unreleased] - Production Enhancements

### Added
- **Security**
  - WebSocket origin validation against configured CORS allowed origins
  - Rate limiting middleware (token bucket algorithm, default: 100 req/min)
  - Input validation for container IDs to prevent path traversal attacks
  - Environment-based API/WebSocket URL configuration

- **Observability**
  - Prometheus metrics endpoint (`/metrics`)
  - Correlation ID middleware for request tracing
  - Structured error responses with request IDs and timestamps
  - Docker Compose health check configuration

- **Frontend**
  - React ErrorBoundary component for graceful error handling
  - Automatic retry logic with exponential backoff for API calls
  - Dynamic API/WebSocket URL configuration based on environment
  - Improved error messages and user feedback

- **Infrastructure**
  - Health check endpoint monitoring in Docker Compose
  - Configuration defaults for rate limiting
  - Shared WebSocket constants to reduce code duplication

### Changed
- WebSocket handlers now validate container IDs
- All API endpoints use structured error responses
- Frontend uses environment variables for API configuration
- Error handling improved across the application

### Security
- Fixed WebSocket origin validation vulnerability
- Added rate limiting to prevent DoS attacks
- Added input validation to prevent injection attacks
- Improved CORS configuration handling

### Performance
- Optimized stats update frequency using streaming API (1 second intervals)
- Improved error handling reduces unnecessary retries

## [Previous] - Initial Release

### Features
- Real-time container monitoring via WebSocket
- CPU and Memory metrics visualization
- Container log streaming
- Container control operations (start, stop, restart, pause, unpause)
- Responsive UI with Tailwind CSS
- Single binary deployment
