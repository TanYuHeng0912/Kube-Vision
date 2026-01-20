package metrics

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// HTTP metrics
	httpRequestsTotal   = make(map[string]int64)
	httpRequestsLatency = make(map[string][]time.Duration)
	httpRequestsMutex   sync.RWMutex

	// WebSocket metrics
	websocketConnectionsActive int64
	websocketConnectionsTotal   int64
	websocketConnectionsMutex   sync.RWMutex
)

// RecordHTTPRequest records an HTTP request
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	httpRequestsMutex.Lock()
	defer httpRequestsMutex.Unlock()

	key := fmt.Sprintf("%s %s %d", method, path, statusCode)
	httpRequestsTotal[key]++

	if httpRequestsLatency[key] == nil {
		httpRequestsLatency[key] = make([]time.Duration, 0, 100)
	}
	httpRequestsLatency[key] = append(httpRequestsLatency[key], duration)
	if len(httpRequestsLatency[key]) > 100 {
		httpRequestsLatency[key] = httpRequestsLatency[key][1:]
	}
}

// IncrementWebSocketConnections increments active WebSocket connections
func IncrementWebSocketConnections() {
	websocketConnectionsMutex.Lock()
	defer websocketConnectionsMutex.Unlock()
	websocketConnectionsActive++
	websocketConnectionsTotal++
}

// DecrementWebSocketConnections decrements active WebSocket connections
func DecrementWebSocketConnections() {
	websocketConnectionsMutex.Lock()
	defer websocketConnectionsMutex.Unlock()
	if websocketConnectionsActive > 0 {
		websocketConnectionsActive--
	}
}

// MetricsHandler returns Prometheus-compatible metrics
func MetricsHandler(c *gin.Context) {
	httpRequestsMutex.RLock()
	websocketConnectionsMutex.RLock()
	defer httpRequestsMutex.RUnlock()
	defer websocketConnectionsMutex.RUnlock()

	var output string

	// HTTP request metrics
	output += "# HELP http_requests_total Total number of HTTP requests\n"
	output += "# TYPE http_requests_total counter\n"
	for key, count := range httpRequestsTotal {
		output += fmt.Sprintf("http_requests_total{key=\"%s\"} %d\n", key, count)
	}

	// HTTP latency metrics
	output += "# HELP http_request_duration_seconds HTTP request duration in seconds\n"
	output += "# TYPE http_request_duration_seconds histogram\n"
	for key, latencies := range httpRequestsLatency {
		if len(latencies) > 0 {
			var sum time.Duration
			for _, lat := range latencies {
				sum += lat
			}
			avg := sum / time.Duration(len(latencies))
			output += fmt.Sprintf("http_request_duration_seconds{key=\"%s\"} %f\n", key, avg.Seconds())
		}
	}

	// WebSocket metrics
	output += "# HELP websocket_connections_active Current number of active WebSocket connections\n"
	output += "# TYPE websocket_connections_active gauge\n"
	output += fmt.Sprintf("websocket_connections_active %d\n", websocketConnectionsActive)

	output += "# HELP websocket_connections_total Total number of WebSocket connections\n"
	output += "# TYPE websocket_connections_total counter\n"
	output += fmt.Sprintf("websocket_connections_total %d\n", websocketConnectionsTotal)

	c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(output))
}

