package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kubevision/kubevision/internal/utils"
)

// ContainerHandler handles container-related API endpoints
type ContainerHandler struct {
	dockerClient interface {
		ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
		ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	}
	logger *zap.Logger
}

// NewContainerHandler creates a new container handler
func NewContainerHandler(dockerClient interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}, logger *zap.Logger) *ContainerHandler {
	return &ContainerHandler{
		dockerClient: dockerClient,
		logger:       logger,
	}
}

// ContainerInfo represents a container in the API response
type ContainerInfo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	Status     string    `json:"status"`
	State      string    `json:"state"`
	Created    time.Time `json:"created"`
	Ports      []Port    `json:"ports"`
	Labels     map[string]string `json:"labels"`
}

// Port represents a container port mapping
type Port struct {
	PrivatePort uint16 `json:"private_port"`
	PublicPort  uint16 `json:"public_port"`
	Type        string `json:"type"`
}

// APIResponse represents a standardized API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Meta      *Meta       `json:"meta,omitempty"`
}

// Meta contains metadata about the response
type Meta struct {
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
}

// ListContainers handles GET /api/containers
func (h *ContainerHandler) ListContainers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all containers (including stopped)
	containers, err := h.dockerClient.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		h.logger.Error("Failed to list containers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success:   false,
			Error:     fmt.Sprintf("Failed to list containers: %v", err),
			Timestamp: time.Now(),
		})
		return
	}

	// Convert to API format
	containerInfos := make([]ContainerInfo, 0, len(containers))
	for _, container := range containers {
		// Get container name (first name in Names slice, remove leading /)
		name := container.ID[:12] // Default to short ID if no name
		if len(container.Names) > 0 && len(container.Names[0]) > 0 {
			name = container.Names[0]
			if name[0] == '/' {
				name = name[1:]
			}
		}

		// Convert ports
		ports := make([]Port, 0, len(container.Ports))
		for _, p := range container.Ports {
			ports = append(ports, Port{
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}

		containerInfos = append(containerInfos, ContainerInfo{
			ID:      container.ID,
			Name:    name,
			Image:   container.Image,
			Status:  container.Status,
			State:   container.State,
			Created: time.Unix(container.Created, 0),
			Ports:   ports,
			Labels:  container.Labels,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      containerInfos,
		Timestamp: time.Now(),
		Meta: &Meta{
			Total: len(containerInfos),
		},
	})
}

// GetContainer handles GET /api/containers/:id
func (h *ContainerHandler) GetContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		BadRequest(c, "Container ID is required")
		return
	}

	// Validate container ID to prevent path traversal
	if !utils.ValidateContainerID(containerID) {
		BadRequest(c, "Invalid container ID format")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	container, err := h.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		h.logger.Error("Failed to inspect container", zap.String("container_id", containerID), zap.Error(err))
		c.JSON(http.StatusNotFound, APIResponse{
			Success:   false,
			Error:     "Container not found",
			Timestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      container,
		Timestamp: time.Now(),
	})
}

