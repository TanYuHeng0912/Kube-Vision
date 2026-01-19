package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		allowedOrigins []string
		origin         string
		expectedStatus int
		expectedHeader string
	}{
		{
			name:           "allowed origin",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			expectedHeader: "http://localhost:3000",
		},
		{
			name:           "wildcard allows all",
			allowedOrigins: []string{"*"},
			origin:         "http://any-origin.com",
			expectedStatus: http.StatusOK,
			expectedHeader: "http://any-origin.com",
		},
		{
			name:           "no origin header",
			allowedOrigins: []string{"http://localhost:3000"},
			origin:         "",
			expectedStatus: http.StatusOK,
			expectedHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(CORSMiddleware(tt.allowedOrigins))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedHeader != "" {
				actualHeader := w.Header().Get("Access-Control-Allow-Origin")
				if actualHeader != tt.expectedHeader {
					t.Errorf("Expected header %s, got %s", tt.expectedHeader, actualHeader)
				}
			}
		})
	}
}

func TestCORSMiddleware_OptionsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSMiddleware([]string{"http://localhost:3000"}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusNoContent, w.Code)
	}
}


