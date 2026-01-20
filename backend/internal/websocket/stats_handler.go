package websocket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/kubevision/kubevision/internal/docker"
	"github.com/kubevision/kubevision/internal/utils"
)


// StatsHandler handles WebSocket connections for container stats
func StatsHandler(dockerClient interface {
	ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponseReader, error)
}, statsCalculator *docker.StatsCalculator, logger *zap.Logger) gin.HandlerFunc {
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

		// Upgrade connection to WebSocket
		upgrader := GetUpgrader()
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("Failed to upgrade connection", zap.Error(err))
			return
		}
		defer conn.Close()

		// Set connection parameters
		_ = conn.SetReadDeadline(time.Now().Add(PongWait))
		conn.SetPongHandler(func(string) error {
			_ = conn.SetReadDeadline(time.Now().Add(PongWait))
			return nil
		})

		// Create context for this connection
		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()

		// Channel for stats with buffer to prevent blocking
		statsChan := make(chan *docker.ContainerStats, 100)

		// Goroutine to read Docker stats using streaming API
		go func() {
			defer close(statsChan)
			
			// Use streaming API for better performance
			stats, err := dockerClient.ContainerStats(ctx, containerID, true)
			if err != nil {
				logger.Error("Failed to get container stats stream",
					zap.String("container_id", containerID),
					zap.Error(err))
				return
			}
			defer stats.Body.Close()

			decoder := json.NewDecoder(stats.Body)
			var statsJSON container.StatsResponse
			lastSendTime := time.Now()
			const sendInterval = 1 * time.Second

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Decode next stats JSON from stream
					if err := decoder.Decode(&statsJSON); err != nil {
						if err.Error() == "EOF" {
							logger.Info("Stats stream ended",
								zap.String("container_id", containerID))
						} else {
							logger.Error("Failed to decode stats",
								zap.String("container_id", containerID),
								zap.Error(err))
						}
						return
					}

					// Only send if enough time has passed (throttle to 1 second)
					now := time.Now()
					if now.Sub(lastSendTime) >= sendInterval {
						// Calculate metrics
						calculatedStats, err := statsCalculator.CalculateStats(containerID, &statsJSON)
						if err != nil {
							logger.Error("Failed to calculate stats",
								zap.String("container_id", containerID),
								zap.Error(err))
							continue
						}

						// Send to channel (non-blocking)
						select {
						case statsChan <- calculatedStats:
							lastSendTime = now
						default:
							// Buffer full, skip this frame
							logger.Warn("Stats channel buffer full, skipping frame",
								zap.String("container_id", containerID))
						}
					}
				}
			}
		}()

		// Goroutine to send ping messages
		pingTicker := time.NewTicker(PingPeriod)
		defer pingTicker.Stop()

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-pingTicker.C:
					_ = conn.SetWriteDeadline(time.Now().Add(WriteWait))
					if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						return
					}
				}
			}
		}()

		// Main loop: send stats to client
		for {
			select {
			case <-ctx.Done():
				return
			case stats, ok := <-statsChan:
				_ = conn.SetWriteDeadline(time.Now().Add(WriteWait))
				if !ok {
					_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				// Send stats as JSON
				if err := conn.WriteJSON(stats); err != nil {
					logger.Error("Failed to write stats",
						zap.String("container_id", containerID),
						zap.Error(err))
					return
				}
			}
		}
	}
}

