package api

import (
	"context"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ImageHandler handles image-related API endpoints
type ImageHandler struct {
	dockerClient interface {
		ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
		ImageInspectWithRaw(ctx context.Context, imageID string) (image.InspectResponse, []byte, error)
		ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	}
	logger *zap.Logger
}

// NewImageHandler creates a new image handler
func NewImageHandler(dockerClient interface {
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageInspectWithRaw(ctx context.Context, imageID string) (image.InspectResponse, []byte, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
}, logger *zap.Logger) *ImageHandler {
	return &ImageHandler{
		dockerClient: dockerClient,
		logger:       logger,
	}
}

// ImageInfo represents a Docker image in the API response
type ImageInfo struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repo_tags"`
	RepoDigests []string          `json:"repo_digests"`
	Size        int64             `json:"size"`
	Created     time.Time         `json:"created"`
	Labels      map[string]string `json:"labels"`
}

// ListImages handles GET /api/images
func (h *ImageHandler) ListImages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all images
	images, err := h.dockerClient.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		h.logger.Error("Failed to list images", zap.Error(err))
		InternalServerError(c, "Failed to list images", err.Error())
		return
	}

	// Convert to API format
	imageInfos := make([]ImageInfo, 0, len(images))
	for _, img := range images {
		imageInfos = append(imageInfos, ImageInfo{
			ID:          img.ID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Size:        img.Size,
			Created:     time.Unix(img.Created, 0),
			Labels:      img.Labels,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      imageInfos,
		Timestamp: time.Now(),
		Meta: &Meta{
			Total: len(imageInfos),
		},
	})
}

// GetImage handles GET /api/images/:id
func (h *ImageHandler) GetImage(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		BadRequest(c, "Image ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, _, err := h.dockerClient.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		h.logger.Error("Failed to inspect image", zap.String("image_id", imageID), zap.Error(err))
		NotFound(c, "Image not found")
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      image,
		Timestamp: time.Now(),
	})
}

// RemoveImage handles DELETE /api/images/:id
func (h *ImageHandler) RemoveImage(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		BadRequest(c, "Image ID is required")
		return
	}

	force := c.DefaultQuery("force", "false") == "true"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := h.dockerClient.ImageRemove(ctx, imageID, image.RemoveOptions{Force: force})
	if err != nil {
		h.logger.Error("Failed to remove image",
			zap.String("image_id", imageID),
			zap.Error(err))
		InternalServerError(c, "Failed to remove image", err.Error())
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      gin.H{"message": "Image removed successfully"},
		Timestamp: time.Now(),
	})
}


