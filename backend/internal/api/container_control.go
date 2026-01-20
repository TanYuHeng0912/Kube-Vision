package api

import (
	"context"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kubevision/kubevision/internal/utils"
)

// ContainerControlHandler handles container control operations
type ContainerControlHandler struct {
	dockerClient interface {
		ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
		ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
		ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
		ContainerPause(ctx context.Context, containerID string) error
		ContainerUnpause(ctx context.Context, containerID string) error
	}
	logger *zap.Logger
}

// NewContainerControlHandler creates a new container control handler
func NewContainerControlHandler(dockerClient interface {
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerPause(ctx context.Context, containerID string) error
	ContainerUnpause(ctx context.Context, containerID string) error
}, logger *zap.Logger) *ContainerControlHandler {
	return &ContainerControlHandler{
		dockerClient: dockerClient,
		logger:       logger,
	}
}

// StartContainer handles POST /api/containers/:id/start
func (h *ContainerControlHandler) StartContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		BadRequest(c, "Container ID is required")
		return
	}

	if !utils.ValidateContainerID(containerID) {
		BadRequest(c, "Invalid container ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.dockerClient.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		h.logger.Error("Failed to start container",
			zap.String("container_id", containerID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Failed to start container",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      gin.H{"message": "Container started successfully"},
		Timestamp: time.Now(),
	})
}

// StopContainer handles POST /api/containers/:id/stop
func (h *ContainerControlHandler) StopContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		BadRequest(c, "Container ID is required")
		return
	}

	if !utils.ValidateContainerID(containerID) {
		BadRequest(c, "Invalid container ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	timeout := 10 // seconds
	if err := h.dockerClient.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		h.logger.Error("Failed to stop container",
			zap.String("container_id", containerID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Failed to stop container",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      gin.H{"message": "Container stopped successfully"},
		Timestamp: time.Now(),
	})
}

// RestartContainer handles POST /api/containers/:id/restart
func (h *ContainerControlHandler) RestartContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		BadRequest(c, "Container ID is required")
		return
	}

	if !utils.ValidateContainerID(containerID) {
		BadRequest(c, "Invalid container ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	timeout := 10 // seconds
	if err := h.dockerClient.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		h.logger.Error("Failed to restart container",
			zap.String("container_id", containerID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Failed to restart container",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      gin.H{"message": "Container restarted successfully"},
		Timestamp: time.Now(),
	})
}

// PauseContainer handles POST /api/containers/:id/pause
func (h *ContainerControlHandler) PauseContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		BadRequest(c, "Container ID is required")
		return
	}

	if !utils.ValidateContainerID(containerID) {
		BadRequest(c, "Invalid container ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.dockerClient.ContainerPause(ctx, containerID); err != nil {
		h.logger.Error("Failed to pause container",
			zap.String("container_id", containerID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Failed to pause container",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      gin.H{"message": "Container paused successfully"},
		Timestamp: time.Now(),
	})
}

// UnpauseContainer handles POST /api/containers/:id/unpause
func (h *ContainerControlHandler) UnpauseContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		BadRequest(c, "Container ID is required")
		return
	}

	if !utils.ValidateContainerID(containerID) {
		BadRequest(c, "Invalid container ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.dockerClient.ContainerUnpause(ctx, containerID); err != nil {
		h.logger.Error("Failed to unpause container",
			zap.String("container_id", containerID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     "Failed to unpause container",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      gin.H{"message": "Container unpaused successfully"},
		Timestamp: time.Now(),
	})
}

