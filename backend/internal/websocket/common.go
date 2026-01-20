package websocket

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	PongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer
	MaxMessageSize = 512
)

// GetUpgrader returns a WebSocket upgrader with origin validation
func GetUpgrader() websocket.Upgrader {
	allowedOrigins := viper.GetStringSlice("CORS_ALLOWED_ORIGINS")
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			// Allow if origin is in allowed list or if wildcard is set
			for _, allowed := range allowedOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		},
	}
}

