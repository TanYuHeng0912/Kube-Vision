package docker

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

var (
	once           sync.Once
	dockerClient   *client.Client
	clientInitErr  error
)

// DockerClient wraps the Docker API client with singleton pattern
type DockerClient struct {
	client *client.Client
	logger *zap.Logger
}

// GetClient returns a singleton Docker client instance
func GetClient(logger *zap.Logger) (*DockerClient, error) {
	once.Do(func() {
		// Use default client which auto-detects Docker connection on Windows/Linux/Mac
		// This will use DOCKER_HOST env var if set, otherwise auto-detect
		opts := []client.Opt{
			client.WithAPIVersionNegotiation(),
		}

		// Only set custom host if DOCKER_HOST is explicitly set
		if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost != "" {
			opts = append(opts, client.WithHost(dockerHost))
			logger.Info("Using custom DOCKER_HOST", zap.String("host", dockerHost))
		}

		cli, err := client.NewClientWithOpts(opts...)
		if err != nil {
			clientInitErr = fmt.Errorf("failed to create Docker client: %w", err)
			return
		}

		// Test connection
		ctx := context.Background()
		_, err = cli.Ping(ctx)
		if err != nil {
			clientInitErr = fmt.Errorf("failed to ping Docker daemon: %w", err)
			return
		}

		dockerClient = cli
		logger.Info("Docker client initialized successfully")
	})

	if clientInitErr != nil {
		return nil, clientInitErr
	}

	return &DockerClient{
		client: dockerClient,
		logger: logger,
	}, nil
}

// GetRawClient returns the underlying Docker client for direct API access
func (dc *DockerClient) GetRawClient() *client.Client {
	return dc.client
}

// getDockerHost retrieves Docker host from environment or defaults to Unix socket
func getDockerHost() string {
	// Check environment variable first
	if host := os.Getenv("DOCKER_HOST"); host != "" {
		return host
	}
	
	// Default to Unix socket
	return "unix:///var/run/docker.sock"
}

// Close closes the Docker client connection
func (dc *DockerClient) Close() error {
	if dc.client != nil {
		return dc.client.Close()
	}
	return nil
}

// HealthCheck verifies the Docker connection is still alive
func (dc *DockerClient) HealthCheck(ctx context.Context) error {
	_, err := dc.client.Ping(ctx)
	return err
}

