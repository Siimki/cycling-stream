'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { getToken } from '@/lib/auth';
import { ChatMessage } from '@/lib/api';
import { WS_URL } from '@/lib/config';
import { WEBSOCKET_PING_INTERVAL_MS, WEBSOCKET_RECONNECT_DELAY_MS, WEBSOCKET_MAX_RECONNECT_DELAY_MS } from '@/constants/intervals';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('Chat');

interface WSMessage {
  type: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data?: any; // WebSocket message data can be any shape
}

interface UseChatReturn {
  messages: ChatMessage[];
  sendMessage: (message: string) => void;
  isConnected: boolean;
  error: string | null;
  reconnect: () => void;
}

export function useChat(raceId: string, enabled: boolean): UseChatReturn {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // We use a ref for the WebSocket to access it in cleanup/callbacks without re-triggering effects
  const wsRef = useRef<WebSocket | null>(null);
  const errorTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  
  // This state is used to force a reconnection
  const [reconnectTrigger, setReconnectTrigger] = useState(0);

  // Helper to set error with optional delay
  const setChatError = useCallback((msg: string | null, delayMs: number = 0) => {
    if (errorTimeoutRef.current) {
      clearTimeout(errorTimeoutRef.current);
      errorTimeoutRef.current = null;
    }

    if (msg === null) {
      setError(null);
      return;
    }

    if (delayMs > 0) {
      errorTimeoutRef.current = setTimeout(() => {
        setError(msg);
      }, delayMs);
    } else {
      setError(msg);
    }
  }, []);

  // Clean up error timeout on unmount
  useEffect(() => {
    return () => {
      if (errorTimeoutRef.current) clearTimeout(errorTimeoutRef.current);
    };
  }, []);

  useEffect(() => {
    if (!enabled || !raceId) {
      return;
    }

    let ws: WebSocket | null = null;
    let pingInterval: NodeJS.Timeout;
    let reconnectTimeout: NodeJS.Timeout;
    let isCleanup = false;
    let reconnectDelay = WEBSOCKET_RECONNECT_DELAY_MS;
    const maxReconnectDelay = WEBSOCKET_MAX_RECONNECT_DELAY_MS;

    const connect = () => {
      if (isCleanup) return;

      try {
        const token = getToken();
        const wsUrl = token
          ? `${WS_URL}/races/${raceId}/chat/ws?token=${encodeURIComponent(token)}`
          : `${WS_URL}/races/${raceId}/chat/ws`;

        logger.debug('Connecting to chat:', { raceId });
        ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
          if (isCleanup) {
            ws?.close();
            return;
          }
          logger.debug('Chat connected');
          setIsConnected(true);
          setChatError(null);
          reconnectDelay = WEBSOCKET_RECONNECT_DELAY_MS; // Reset backoff

          // Setup ping
          pingInterval = setInterval(() => {
            if (ws?.readyState === WebSocket.OPEN) {
              ws.send(JSON.stringify({ type: 'ping' }));
            }
          }, WEBSOCKET_PING_INTERVAL_MS);
        };

        ws.onmessage = (event) => {
          if (isCleanup) return;
          
          try {
            const msg: WSMessage = JSON.parse(event.data);
            
            switch (msg.type) {
              case 'message':
                if (msg.data) {
                  const chatMsg: ChatMessage = {
                    id: msg.data.id,
                    race_id: msg.data.race_id,
                    user_id: msg.data.user_id,
                    username: msg.data.username,
                    message: msg.data.message,
                    created_at: msg.data.created_at,
                  };
                  setMessages((prev) => {
                    if (prev.some(m => m.id === chatMsg.id)) return prev;
                    return [...prev, chatMsg];
                  });
                }
                break;
              case 'error':
                logger.warn('Chat server error:', msg.data?.message);
                if (msg.data?.message) {
                    // Delay showing error by 500ms to avoid flashing during transient connection issues
                    setChatError(msg.data.message, 500);
                }
                break;
            }
          } catch (e) {
            logger.error('Failed to parse chat message:', e);
          }
        };

        ws.onclose = (event) => {
          if (isCleanup) return;
          
          setIsConnected(false);
          wsRef.current = null;
          clearInterval(pingInterval);

          logger.debug('Chat disconnected:', event.code, event.reason);

          // If strictly normal closure (1000), don't reconnect
          if (event.code === 1000) {
            return;
          }

          // Otherwise, retry with exponential backoff
          logger.debug(`Reconnecting in ${reconnectDelay}ms...`);
          reconnectTimeout = setTimeout(() => {
            reconnectDelay = Math.min(reconnectDelay * 1.5, maxReconnectDelay);
            connect();
          }, reconnectDelay);
        };

        ws.onerror = (e) => {
           // Error details usually appear in onclose
           logger.debug('WebSocket error:', e);
        };

      } catch (e) {
        logger.error('Failed to create WebSocket:', e);
        // Retry if creation fails (e.g. URL error)
        reconnectTimeout = setTimeout(connect, 5000);
      }
    };

    connect();

    return () => {
      isCleanup = true;
      logger.debug('Cleaning up chat connection');
      if (ws) {
        // Remove listeners to prevent state updates after unmount
        ws.onclose = null;
        ws.onmessage = null;
        ws.onopen = null;
        ws.onerror = null;
        ws.close();
      }
      wsRef.current = null;
      clearInterval(pingInterval);
      clearTimeout(reconnectTimeout);
      // Note: We don't clear errorTimeout here as it's handled by the parent effect
    };
  }, [raceId, enabled, reconnectTrigger, setChatError]);

  const sendMessage = useCallback((message: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({
        type: 'send_message',
        data: { message }
      }));
    } else {
      logger.warn('Cannot send message: Chat not connected');
      setChatError('Not connected');
    }
  }, [setChatError]);

  const reconnect = useCallback(() => {
    setReconnectTrigger(c => c + 1);
  }, []);

  return {
    messages,
    sendMessage,
    isConnected,
    error,
    reconnect
  };
}
