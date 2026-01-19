import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

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
    const response = await apiClient.get<APIResponse<ContainerInfo[]>>('/api/containers');
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || 'Failed to fetch containers');
  },

  // Get container by ID
  async getContainer(id: string): Promise<any> {
    const response = await apiClient.get<APIResponse<any>>(`/api/containers/${id}`);
    if (response.data.success && response.data.data) {
      return response.data.data;
    }
    throw new Error(response.data.error || 'Failed to fetch container');
  },
};

export default api;

