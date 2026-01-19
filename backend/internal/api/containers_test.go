package api

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

// Mock Docker client for testing
type mockDockerClient struct {
	containers []types.Container
	container  types.ContainerJSON
	err        error
}

func (m *mockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.containers, nil
}

func (m *mockDockerClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	if m.err != nil {
		return types.ContainerJSON{}, m.err
	}
	return m.container, nil
}

func TestNewContainerHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	mockClient := &mockDockerClient{
		containers: []types.Container{
			{
				ID:      "container-1",
				Names:   []string{"/test-container"},
				Image:   "nginx:latest",
				Status:  "Up 5 minutes",
				State:   "running",
				Created: time.Now().Unix(),
				Ports: []types.Port{
					{
						PrivatePort: 80,
						PublicPort:  8080,
						Type:        "tcp",
					},
				},
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
	}

	handler := NewContainerHandler(mockClient, logger)

	if handler == nil {
		t.Fatal("NewContainerHandler returned nil")
	}

	if handler.dockerClient == nil {
		t.Error("Handler dockerClient is nil")
	}

	if handler.logger == nil {
		t.Error("Handler logger is nil")
	}
}

func TestContainerInfo_Conversion(t *testing.T) {
	// Test container name extraction
	testCases := []struct {
		name     string
		names    []string
		expected string
	}{
		{
			name:     "normal name",
			names:    []string{"/my-container"},
			expected: "my-container",
		},
		{
			name:     "name without slash",
			names:    []string{"my-container"},
			expected: "my-container",
		},
		{
			name:     "empty names",
			names:    []string{},
			expected: "test-id", // When no names, use full ID (or first 12 chars if longer)
		},
		{
			name:     "multiple names",
			names:    []string{"/first", "/second"},
			expected: "first",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := types.Container{
				ID:    "test-id",
				Names: tc.names,
			}

			var result string
			if len(container.Names) > 0 && len(container.Names[0]) > 0 {
				result = container.Names[0]
				if result[0] == '/' {
					result = result[1:]
				}
			} else {
				// Use full ID or first 12 characters, whichever is shorter
				if len(container.ID) >= 12 {
					result = container.ID[:12]
				} else {
					result = container.ID
				}
			}

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}


