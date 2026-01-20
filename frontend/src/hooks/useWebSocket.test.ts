import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useWebSocket } from './useWebSocket';

// Mock WebSocket
class MockWebSocket {
  url: string;
  readyState: number = WebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;

  constructor(url: string) {
    this.url = url;
    // Simulate connection after a short delay
    setTimeout(() => {
      this.readyState = WebSocket.OPEN;
      if (this.onopen) {
        this.onopen(new Event('open'));
      }
    }, 10);
  }

  close() {
    this.readyState = WebSocket.CLOSED;
    if (this.onclose) {
      this.onclose(new CloseEvent('close'));
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  send(_data: string) {
    // Mock send
  }
}

describe('useWebSocket', () => {
  let originalWebSocket: typeof WebSocket;

  beforeEach(() => {
    // Store original WebSocket
    originalWebSocket = globalThis.WebSocket;
    // Replace global WebSocket
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (globalThis as any).WebSocket = MockWebSocket;
    vi.useFakeTimers();
  });

  afterEach(() => {
    // Restore original WebSocket
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (globalThis as any).WebSocket = originalWebSocket;
    vi.useRealTimers();
  });

  it('should initialize with disconnected status when URL is null', () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: null,
        onMessage: vi.fn(),
      })
    );

    expect(result.current.status).toBe('disconnected');
  });

  it('should provide connect and disconnect functions', () => {
    const { result } = renderHook(() =>
      useWebSocket({
        url: null,
        onMessage: vi.fn(),
      })
    );

    expect(typeof result.current.connect).toBe('function');
    expect(typeof result.current.disconnect).toBe('function');
  });

  it('should handle URL changes', () => {
    const { result, rerender } = renderHook(
      ({ url }) =>
        useWebSocket({
          url,
          onMessage: vi.fn(),
        }),
      {
        initialProps: { url: null as string | null },
      }
    );

    expect(result.current.status).toBe('disconnected');

    rerender({ url: 'ws://localhost:8080/ws/stats/test-id' });
    // Status should change to connecting
    expect(['connecting', 'connected']).toContain(result.current.status);

    rerender({ url: null });
    expect(result.current.status).toBe('disconnected');
  });
});


