import { useCallback, useEffect, useRef, useState } from 'react';

export type WebSocketStatus = 'connecting' | 'connected' | 'disconnected' | 'error';

interface UseWebSocketOptions {
  url: string | null;
  onMessage?: (data: unknown) => void;
  onError?: (error: Event) => void;
  onOpen?: () => void;
  onClose?: () => void;
  reconnect?: boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

export function useWebSocket({
  url,
  onMessage,
  onError,
  onOpen,
  onClose,
  reconnect = true,
  reconnectInterval = 1000,
  maxReconnectAttempts = 10,
}: UseWebSocketOptions) {
  const [status, setStatus] = useState<WebSocketStatus>('disconnected');
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Keep latest handlers without forcing reconnect when parent re-renders (new fn identity).
  const onMessageRef = useRef(onMessage);
  const onErrorRef = useRef(onError);
  const onOpenRef = useRef(onOpen);
  const onCloseRef = useRef(onClose);
  onMessageRef.current = onMessage;
  onErrorRef.current = onError;
  onOpenRef.current = onOpen;
  onCloseRef.current = onClose;

  const connect = useCallback(() => {
    if (!url) {
      setStatus('disconnected');
      return;
    }

    try {
      setStatus('connecting');
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        setStatus('connected');
        reconnectAttemptsRef.current = 0;
        onOpenRef.current?.();
      };

      ws.onmessage = (event) => {
        // Handle both text and binary messages
        if (typeof event.data === 'string') {
          try {
            const data = JSON.parse(event.data);
            onMessageRef.current?.(data);
          } catch {
            // If not JSON, pass as string (for logs)
            onMessageRef.current?.(event.data);
          }
        } else {
          // Binary data (for logs)
          const reader = new FileReader();
          reader.onload = () => {
            onMessageRef.current?.(reader.result);
          };
          reader.readAsText(event.data);
        }
      };

      ws.onerror = (error) => {
        setStatus('error');
        onErrorRef.current?.(error);
      };

      ws.onclose = () => {
        setStatus('disconnected');
        onCloseRef.current?.();
        wsRef.current = null;

        // Attempt reconnection
        if (reconnect && reconnectAttemptsRef.current < maxReconnectAttempts) {
          const delay = Math.min(
            reconnectInterval * Math.pow(2, reconnectAttemptsRef.current),
            30000 // Max 30 seconds
          );
          
          reconnectAttemptsRef.current++;
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, delay);
        }
      };
    } catch (error) {
      setStatus('error');
      console.error('WebSocket connection error:', error);
    }
  }, [url, reconnect, reconnectInterval, maxReconnectAttempts]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    
    setStatus('disconnected');
  }, []);

  useEffect(() => {
    connect();
    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    status,
    connect,
    disconnect,
  };
}

