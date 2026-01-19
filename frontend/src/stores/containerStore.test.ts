import { describe, it, expect, beforeEach } from 'vitest';
import { useContainerStore, ContainerMetrics } from './containerStore';

describe('ContainerStore', () => {
  beforeEach(() => {
    // Reset store before each test by clearing all container metrics
    const state = useContainerStore.getState();
    const containerIds = Object.keys(state.metrics);
    containerIds.forEach((id) => state.clearMetrics(id));
  });

  it('should add metrics for a container', () => {
    const metrics: ContainerMetrics = {
      container_id: 'test-container',
      timestamp: new Date().toISOString(),
      cpu_percent: 50.5,
      memory_usage: 1000000000,
      memory_limit: 2000000000,
      memory_percent: 50.0,
      network_rx: 1000,
      network_tx: 2000,
      block_read: 5000,
      block_write: 10000,
      pids: 10,
    };

    useContainerStore.getState().addMetrics('test-container', metrics);

    const storedMetrics = useContainerStore.getState().metrics['test-container'] || [];
    expect(storedMetrics).toHaveLength(1);
    expect(storedMetrics[0].cpu_percent).toBe(50.5);
  });

  it('should limit metrics to 60 data points', () => {
    const metrics: ContainerMetrics = {
      container_id: 'test-container',
      timestamp: new Date().toISOString(),
      cpu_percent: 50.5,
      memory_usage: 1000000000,
      memory_limit: 2000000000,
      memory_percent: 50.0,
      network_rx: 1000,
      network_tx: 2000,
      block_read: 5000,
      block_write: 10000,
      pids: 10,
    };

    // Add 70 metrics
    for (let i = 0; i < 70; i++) {
      useContainerStore.getState().addMetrics('test-container', {
        ...metrics,
        cpu_percent: i,
      });
    }

    const storedMetrics = useContainerStore.getState().metrics['test-container'] || [];
    expect(storedMetrics).toHaveLength(60);
    // Should keep the most recent 60
    expect(storedMetrics[0].cpu_percent).toBe(10); // First of the last 60
    expect(storedMetrics[59].cpu_percent).toBe(69); // Last one
  });

  it('should clear metrics for a specific container', () => {
    const metrics: ContainerMetrics = {
      container_id: 'test-container',
      timestamp: new Date().toISOString(),
      cpu_percent: 50.5,
      memory_usage: 1000000000,
      memory_limit: 2000000000,
      memory_percent: 50.0,
      network_rx: 1000,
      network_tx: 2000,
      block_read: 5000,
      block_write: 10000,
      pids: 10,
    };

    useContainerStore.getState().addMetrics('container-1', metrics);
    useContainerStore.getState().addMetrics('container-2', metrics);

    // clearMetrics accepts containerId parameter
    useContainerStore.getState().clearMetrics('container-1');

    expect(useContainerStore.getState().metrics['container-1']).toBeUndefined();
    expect(useContainerStore.getState().metrics['container-2']).toHaveLength(1);
  });

  it('should manage multiple containers independently', () => {
    const metrics1: ContainerMetrics = {
      container_id: 'container-1',
      timestamp: new Date().toISOString(),
      cpu_percent: 25.0,
      memory_usage: 500000000,
      memory_limit: 1000000000,
      memory_percent: 50.0,
      network_rx: 1000,
      network_tx: 2000,
      block_read: 5000,
      block_write: 10000,
      pids: 5,
    };

    const metrics2: ContainerMetrics = {
      container_id: 'container-2',
      timestamp: new Date().toISOString(),
      cpu_percent: 75.0,
      memory_usage: 1500000000,
      memory_limit: 2000000000,
      memory_percent: 75.0,
      network_rx: 2000,
      network_tx: 4000,
      block_read: 10000,
      block_write: 20000,
      pids: 10,
    };

    useContainerStore.getState().addMetrics('container-1', metrics1);
    useContainerStore.getState().addMetrics('container-2', metrics2);

    const stored1 = useContainerStore.getState().metrics['container-1'] || [];
    const stored2 = useContainerStore.getState().metrics['container-2'] || [];

    expect(stored1).toHaveLength(1);
    expect(stored2).toHaveLength(1);
    expect(stored1[0].cpu_percent).toBe(25.0);
    expect(stored2[0].cpu_percent).toBe(75.0);
  });
});


