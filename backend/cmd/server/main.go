package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/kubevision/kubevision/internal/api"
	"github.com/kubevision/kubevision/internal/docker"
	"github.com/kubevision/kubevision/internal/metrics"
	"github.com/kubevision/kubevision/internal/middleware"
	"github.com/kubevision/kubevision/internal/websocket"
)

func main() {
	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Load configuration
	if err := loadConfig(); err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Docker client
	dockerClient, err := docker.GetClient(logger)
	if err != nil {
		logger.Fatal("Failed to initialize Docker client", zap.Error(err))
	}
	defer dockerClient.Close()

	// Initialize stats calculator
	statsCalculator := docker.NewStatsCalculator(logger)

	// Initialize Gin router
	if viper.GetString("LOG_LEVEL") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Correlation ID middleware (must be first)
	router.Use(middleware.CorrelationIDMiddleware())

	// CORS middleware
	allowedOrigins := viper.GetStringSlice("CORS_ALLOWED_ORIGINS")
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	router.Use(middleware.CORSMiddleware(allowedOrigins))

	// Rate limiting middleware
	rateLimitEnabled := viper.GetBool("RATE_LIMIT_ENABLED")
	if rateLimitEnabled {
		rateLimit := viper.GetInt("RATE_LIMIT_REQUESTS")
		if rateLimit == 0 {
			rateLimit = 100 // Default: 100 requests per minute
		}
		rateLimitDuration := viper.GetDuration("RATE_LIMIT_DURATION")
		if rateLimitDuration == 0 {
			rateLimitDuration = 1 * time.Minute
		}
		router.Use(middleware.RateLimitMiddleware(rateLimit, rateLimitDuration))
	}

	// Metrics endpoint (Prometheus format)
	router.GET("/metrics", metrics.MetricsHandler)

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := dockerClient.HealthCheck(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"message": "Docker connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Service is running",
		})
	})

	// API routes
	apiGroup := router.Group("/api")
	{
		// Container routes
		containerHandler := api.NewContainerHandler(dockerClient.GetRawClient(), logger)
		apiGroup.GET("/containers", containerHandler.ListContainers)
		apiGroup.GET("/containers/:id", containerHandler.GetContainer)

		// Container control routes (require auth)
		authEnabled := viper.GetBool("AUTH_ENABLED")
		authToken := viper.GetString("AUTH_TOKEN")
		controlHandler := api.NewContainerControlHandler(dockerClient.GetRawClient(), logger)
		controlGroup := apiGroup.Group("/containers/:id")
		controlGroup.Use(middleware.AuthMiddleware(authEnabled, authToken))
		{
			controlGroup.POST("/start", controlHandler.StartContainer)
			controlGroup.POST("/stop", controlHandler.StopContainer)
			controlGroup.POST("/restart", controlHandler.RestartContainer)
			controlGroup.POST("/pause", controlHandler.PauseContainer)
			controlGroup.POST("/unpause", controlHandler.UnpauseContainer)
		}

		// Image routes
		imageHandler := api.NewImageHandler(dockerClient.GetRawClient(), logger)
		apiGroup.GET("/images", imageHandler.ListImages)
		apiGroup.GET("/images/:id", imageHandler.GetImage)
		imageControlGroup := apiGroup.Group("/images/:id")
		imageControlGroup.Use(middleware.AuthMiddleware(authEnabled, authToken))
		{
			imageControlGroup.DELETE("", imageHandler.RemoveImage)
		}
	}

	// WebSocket routes (must be before static files)
	wsGroup := router.Group("/ws")
	{
		wsGroup.GET("/stats/:id", websocket.StatsHandler(
			dockerClient.GetRawClient(),
			statsCalculator,
			logger,
		))
		wsGroup.GET("/logs/:id", websocket.LogsHandler(
			dockerClient.GetRawClient(),
			logger,
		))
		wsGroup.GET("/events", websocket.EventsHandler(
			dockerClient.GetRawClient(),
			logger,
		))
	}

	// Serve static files (frontend) - simple direct approach
	// Serve assets directory
	router.StaticFS("/assets", http.Dir("web/dist/assets"))
	
	// Serve index.html for root
	router.GET("/", func(c *gin.Context) {
		c.File("web/dist/index.html")
	})
	
	// Fallback for SPA routes - serve index.html for all non-API routes
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/ws") || path == "/metrics" {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		// For SPA, serve index.html
		c.File("web/dist/index.html")
	})

	// Start server
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}
	host := viper.GetString("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Starting server", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	
	logLevel := viper.GetString("LOG_LEVEL")
	switch logLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	return config.Build()
}

func loadConfig() error {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DOCKER_HOST", "unix:///var/run/docker.sock")
	viper.SetDefault("AUTH_ENABLED", false)
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_DURATION", "1m")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")

	// Read from environment variables
	viper.AutomaticEnv()

	// Try to read config file (optional)
	viper.ReadInConfig()

	return nil
}
