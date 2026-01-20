import axios from 'axios';
import { API_URL } from '../config';
import { retry, isRetryableError } from '../utils/retry';

const API_BASE_URL = API_URL;

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface ContainerInfo {
  id: string;
  name: string;
  image: string;
  status: string;
  state: string;
  created: string;
  ports: Port[];
  labels: Record<string, string>;
}

export interface Port {
  private_port: number;
  public_port: number;
  type: string;
}

export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  timestamp: string;
  meta?: {
    total?: number;
    page?: number;
  };
}

export const api = {
  // Health check
  async healthCheck(): Promise<{ status: string; message: string }> {
    const response = await apiClient.get('/api/health');
    return response.data;
  },

  // Get all containers
  async getContainers(): Promise<ContainerInfo[]> {
    return retry(
      async () => {
        const response = await apiClient.get<APIResponse<ContainerInfo[]>>('/api/containers');
        if (response.data.success && response.data.data) {
          return response.data.data;
        }
        throw new Error(response.data.error || 'Failed to fetch containers');
      },
      { retryable: isRetryableError }
    );
  },

  // Get container by ID
  async getContainer(id: string): Promise<Record<string, unknown>> {
    return retry(
      async () => {
        const response = await apiClient.get<APIResponse<Record<string, unknown>>>(`/api/containers/${id}`);
        if (response.data.success && response.data.data) {
          return response.data.data;
        }
        throw new Error(response.data.error || 'Failed to fetch container');
      },
      { retryable: isRetryableError }
    );
  },

  // Image management
  async getImages(): Promise<ImageInfo[]> {
    return retry(
      async () => {
        const response = await apiClient.get<APIResponse<ImageInfo[]>>('/api/images');
        if (response.data.success && response.data.data) {
          return response.data.data;
        }
        throw new Error(response.data.error || 'Failed to fetch images');
      },
      { retryable: isRetryableError }
    );
  },

  async getImage(id: string): Promise<Record<string, unknown>> {
    return retry(
      async () => {
        const response = await apiClient.get<APIResponse<Record<string, unknown>>>(`/api/images/${id}`);
        if (response.data.success && response.data.data) {
          return response.data.data;
        }
        throw new Error(response.data.error || 'Failed to fetch image');
      },
      { retryable: isRetryableError }
    );
  },

  async removeImage(id: string, force: boolean = false): Promise<void> {
    return retry(
      async () => {
        await apiClient.delete(`/api/images/${id}?force=${force}`);
      },
      { retryable: isRetryableError }
    );
  },
};

export interface ImageInfo {
  id: string;
  repo_tags: string[];
  repo_digests: string[];
  size: number;
  created: string;
  labels: Record<string, string>;
}

export default api;

