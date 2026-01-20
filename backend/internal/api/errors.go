package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// APIError represents a structured error response
type APIError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

// ErrorResponse sends a structured error response
func ErrorResponse(c *gin.Context, statusCode int, message string, details ...string) {
	correlationID, _ := c.Get("correlation_id")
	correlationIDStr := ""
	if id, ok := correlationID.(string); ok {
		correlationIDStr = id
	}

	error := APIError{
		Code:      statusCode,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		RequestID: correlationIDStr,
	}

	if len(details) > 0 {
		error.Details = details[0]
	}

	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   error,
	})
}

// BadRequest sends a 400 Bad Request error
func BadRequest(c *gin.Context, message string, details ...string) {
	ErrorResponse(c, http.StatusBadRequest, message, details...)
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// Forbidden sends a 403 Forbidden error
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

// NotFound sends a 404 Not Found error
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerError sends a 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, details ...string) {
	ErrorResponse(c, http.StatusInternalServerError, message, details...)
}

