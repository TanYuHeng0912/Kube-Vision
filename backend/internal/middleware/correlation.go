package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CorrelationIDHeader = "X-Correlation-ID"
const CorrelationIDKey = "correlation_id"

// CorrelationIDMiddleware adds a correlation ID to each request for tracing
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get correlation ID from header or generate new one
		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Store in context for use in handlers
		c.Set(CorrelationIDKey, correlationID)

		// Add to response header
		c.Writer.Header().Set(CorrelationIDHeader, correlationID)

		c.Next()
	}
}

