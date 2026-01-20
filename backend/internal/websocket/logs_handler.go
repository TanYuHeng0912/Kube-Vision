package websocket

import (
	"context"
	"encoding/binary"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/kubevision/kubevision/internal/utils"
)


// LogsHandler handles WebSocket connections for container logs
func LogsHandler(dockerClient interface {
	ContainerLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error)
}, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		containerID := c.Param("id")
		if containerID == "" {
			c.JSON(400, gin.H{"error": "Container ID is required"})
			return
		}

		// Validate container ID
		if !utils.ValidateContainerID(containerID) {
			c.JSON(400, gin.H{"error": "Invalid container ID format"})
			return
		}

		// Get query parameters
		tail := c.DefaultQuery("tail", "100")
		follow := c.DefaultQuery("follow", "true") == "true"
		since := c.DefaultQuery("since", "")

		// Upgrade connection to WebSocket
		upgrader := GetUpgrader()
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("Failed to upgrade connection", zap.Error(err))
			return
		}
		defer conn.Close()

		// Set connection parameters
		conn.SetReadDeadline(time.Now().Add(PongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(PongWait))
			return nil
		})

		// Create context for this connection
		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()

		// Configure log options
		logOptions := container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     follow,
			Tail:       tail,
			Timestamps: true,
		}

		if since != "" {
			logOptions.Since = since
		}

		// Get container logs
		logsReader, err := dockerClient.ContainerLogs(ctx, containerID, logOptions)
		if err != nil {
			logger.Error("Failed to get container logs",
				zap.String("container_id", containerID),
				zap.Error(err))
			conn.WriteJSON(gin.H{"error": "Failed to get logs"})
			return
		}
		defer logsReader.Close()

		// Buffer for reading logs
		buffer := make([]byte, 8192)

		// Read and send logs
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := logsReader.Read(buffer)
				if err != nil {
					if err == io.EOF {
						// End of stream
						return
					}
					logger.Error("Failed to read logs",
						zap.String("container_id", containerID),
						zap.Error(err))
					return
				}

				if n == 0 {
					continue
				}

				// Process Docker log header and send data
				data := buffer[:n]
				processedData := stripDockerHeader(data)

				if len(processedData) > 0 {
					conn.SetWriteDeadline(time.Now().Add(WriteWait))
					if err := conn.WriteMessage(websocket.TextMessage, processedData); err != nil {
						logger.Error("Failed to write logs",
							zap.String("container_id", containerID),
							zap.Error(err))
						return
					}
				}
			}
		}
	}
}

// stripDockerHeader removes the 8-byte Docker log stream header
// Format: [STREAM_TYPE(1)][RESERVED(3)][SIZE(4)]
func stripDockerHeader(data []byte) []byte {
	if len(data) < 8 {
		return data
	}

	result := make([]byte, 0, len(data))
	offset := 0

	for offset < len(data) {
		if offset+8 > len(data) {
			// Not enough data for a header, append remaining
			result = append(result, data[offset:]...)
			break
		}

		// Extract size from bytes 4-7 (big-endian)
		size := binary.BigEndian.Uint32(data[offset+4 : offset+8])

		// Check if we have enough data
		if offset+8+int(size) > len(data) {
			// Partial frame, append remaining
			result = append(result, data[offset+8:]...)
			break
		}

		// Extract payload (skip 8-byte header)
		payload := data[offset+8 : offset+8+int(size)]
		result = append(result, payload...)

		offset += 8 + int(size)
	}

	return result
}

