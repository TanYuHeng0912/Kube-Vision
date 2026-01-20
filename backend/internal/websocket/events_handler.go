package websocket

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// DockerEvent represents a Docker event
type DockerEvent struct {
	Type   string                 `json:"type"`
	Action string                 `json:"action"`
	Actor  map[string]interface{} `json:"actor"`
	Time   int64                  `json:"time"`
	TimeNano int64                `json:"timeNano"`
}

// EventsHandler handles WebSocket connections for Docker events
func EventsHandler(dockerClient interface {
	Events(ctx context.Context, options events.ListOptions) (<-chan events.Message, <-chan error)
}, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Get event filters from query params
		eventTypes := c.QueryArray("type")
		eventActions := c.QueryArray("action")

		// Create filters
		eventFilters := filters.NewArgs()
		if len(eventTypes) > 0 {
			for _, t := range eventTypes {
				eventFilters.Add("type", t)
			}
		}
		if len(eventActions) > 0 {
			for _, a := range eventActions {
				eventFilters.Add("action", a)
			}
		}

		// Start listening to Docker events
		eventChan, errChan := dockerClient.Events(ctx, events.ListOptions{
			Filters: eventFilters,
		})

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

		// Main loop: send events to client
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errChan:
				if err != nil {
					logger.Error("Docker events error", zap.Error(err))
					_ = conn.SetWriteDeadline(time.Now().Add(WriteWait))
					_ = conn.WriteJSON(gin.H{"error": "Events stream error"})
					return
				}
			case event := <-eventChan:
				dockerEvent := DockerEvent{
					Type:     string(event.Type),
					Action:   string(event.Action),
					Time:     event.Time,
					TimeNano: event.TimeNano,
					Actor:    make(map[string]interface{}),
				}

				// Convert Actor to map
				if event.Actor.ID != "" {
					dockerEvent.Actor["ID"] = event.Actor.ID
					dockerEvent.Actor["Attributes"] = event.Actor.Attributes
				}

				_ = conn.SetWriteDeadline(time.Now().Add(WriteWait))
				if err := conn.WriteJSON(dockerEvent); err != nil {
					logger.Error("Failed to write event", zap.Error(err))
					return
				}
			}
		}
	}
}

