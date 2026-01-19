# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install --legacy-peer-deps
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Copy frontend build (after backend to ensure proper structure)
COPY --from=frontend-builder /app/frontend/dist ./web/dist

# Build the application (without embed tag, using filesystem)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o kubevision \
    ./cmd/server

# Stage 3: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary from builder
COPY --from=backend-builder /app/kubevision .

# Copy frontend files for runtime (if not embedded)
COPY --from=backend-builder /app/web/dist ./web/dist

# Create non-root user
RUN addgroup -g 1000 kubevision && \
    adduser -D -u 1000 -G kubevision kubevision && \
    chown -R kubevision:kubevision /root

USER kubevision

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

CMD ["./kubevision"]

