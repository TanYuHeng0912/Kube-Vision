import { create } from 'zustand';
import { ContainerInfo } from '../services/api';

export interface ContainerMetrics {
  container_id: string;
  timestamp: string;
  cpu_percent: number;
  memory_usage: number;
  memory_limit: number;
  memory_percent: number;
  network_rx: number;
  network_tx: number;
  block_read: number;
  block_write: number;
  pids: number;
}

interface ContainerState {
  containers: ContainerInfo[];
  metrics: Record<string, ContainerMetrics[]>;
  selectedContainer: string | null;
  isLoading: boolean;
  error: string | null;
  timeUnit: 'seconds' | 'minutes' | 'hours';
  
  // Actions
  setContainers: (containers: ContainerInfo[]) => void;
  addMetrics: (containerId: string, metrics: ContainerMetrics) => void;
  setSelectedContainer: (id: string | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearMetrics: (containerId: string) => void;
  setTimeUnit: (unit: 'seconds' | 'minutes' | 'hours') => void;
}

// Keep only last 10 data points (fixed window)
const MAX_METRICS_POINTS = 10;

export const useContainerStore = create<ContainerState>((set) => ({
  containers: [],
  metrics: {},
  selectedContainer: null,
  isLoading: false,
  error: null,
  timeUnit: 'seconds',

  setContainers: (containers) => set({ containers }),

  setTimeUnit: (unit) => set({ timeUnit: unit }),

  addMetrics: (containerId, newMetrics) =>
    set((state) => {
      const currentMetrics = state.metrics[containerId] || [];
      const updatedMetrics = [...currentMetrics, newMetrics].slice(-MAX_METRICS_POINTS);
      
      return {
        metrics: {
          ...state.metrics,
          [containerId]: updatedMetrics,
        },
      };
    }),

  setSelectedContainer: (id) => set({ selectedContainer: id }),

  setLoading: (loading) => set({ isLoading: loading }),

  setError: (error) => set({ error }),

  clearMetrics: (containerId) =>
    set((state) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { [containerId]: _, ...rest } = state.metrics;
      return { metrics: rest };
    }),
}));

