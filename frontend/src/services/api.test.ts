import { describe, it, expect, vi, beforeEach } from 'vitest';

// Use vi.hoisted to create the mock function before the factory runs
const { mockGet } = vi.hoisted(() => {
  return {
    mockGet: vi.fn(),
  };
});

vi.mock('axios', () => {
  return {
    default: {
      create: () => ({
        get: mockGet,
      }),
    },
  };
});

import api from './api';

describe('API Service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('api.getContainers', () => {
    it('should fetch containers successfully', async () => {
      const mockResponse = {
        success: true,
        data: [
          {
            id: 'container-1',
            name: 'test-container',
            image: 'nginx:latest',
            status: 'running',
            state: 'running',
          },
        ],
        timestamp: new Date().toISOString(),
      };

      mockGet.mockResolvedValue({ data: mockResponse });

      const result = await api.getContainers();

      expect(mockGet).toHaveBeenCalledWith('/api/containers');
      expect(result).toEqual(mockResponse.data);
    });

    it('should handle errors', async () => {
      const errorResponse = {
        success: false,
        error: 'Network error',
        timestamp: new Date().toISOString(),
      };
      mockGet.mockResolvedValue({ data: errorResponse });

      await expect(api.getContainers()).rejects.toThrow('Network error');
    });
  });

  describe('api.getContainer', () => {
    it('should fetch a single container', async () => {
      const mockResponse = {
        success: true,
        data: {
          id: 'container-1',
          name: 'test-container',
          image: 'nginx:latest',
        },
        timestamp: new Date().toISOString(),
      };

      mockGet.mockResolvedValue({ data: mockResponse });

      const result = await api.getContainer('container-1');

      expect(mockGet).toHaveBeenCalledWith('/api/containers/container-1');
      expect(result).toEqual(mockResponse.data);
    });

    it('should handle container not found', async () => {
      const errorResponse = {
        success: false,
        error: 'Container not found',
        timestamp: new Date().toISOString(),
      };
      mockGet.mockResolvedValue({ data: errorResponse });

      await expect(api.getContainer('nonexistent')).rejects.toThrow('Container not found');
    });
  });

  describe('api.healthCheck', () => {
    it('should check health status', async () => {
      const mockResponse = {
        status: 'ok',
        message: 'Server is healthy',
      };

      mockGet.mockResolvedValue({ data: mockResponse });

      const result = await api.healthCheck();

      expect(mockGet).toHaveBeenCalledWith('/api/health');
      expect(result).toEqual(mockResponse);
    });
  });
});


