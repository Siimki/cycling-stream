'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { getToken } from '@/lib/auth';
import { ChatMessage } from '@/lib/api';
import { WS_URL } from '@/lib/config';
import { WEBSOCKET_PING_INTERVAL_MS, WEBSOCKET_RECONNECT_DELAY_MS, WEBSOCKET_MAX_RECONNECT_DELAY_MS } from '@/constants/intervals';
import { createContextLogger } from '@/lib/logger';
import { useAuth } from '@/contexts/AuthContext';

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
  
  // Use AuthContext to detect user changes (more reliable than polling)
  const { user, token } = useAuth();
  const userIdRef = useRef<string | null>(user?.id || null);

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

  // Watch for user changes (when user logs in/out or switches accounts)
  useEffect(() => {
    const currentUserId = user?.id || null;
    const previousUserId = userIdRef.current;

    // If user ID changed, we need to reconnect
    if (currentUserId !== previousUserId) {
      logger.debug('User changed, reconnecting WebSocket', { 
        previousUserId, 
        currentUserId 
      });
      
      userIdRef.current = currentUserId;
      
      // Clear messages when user changes (different user = different chat context)
      setMessages([]);
      
      // Close existing connection if any
      if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
        wsRef.current.onclose = null; // Prevent automatic reconnection
        wsRef.current.close(1000, 'User changed');
        wsRef.current = null;
        setIsConnected(false);
      }
      
      // Force reconnection by updating reconnectTrigger
      setReconnectTrigger(c => c + 1);
    }
  }, [user?.id, token]);

  useEffect(() => {
    if (!enabled || !raceId) {
      return;
    }

    let ws: WebSocket | null = null;
    let pingInterval: NodeJS.Timeout;
    let reconnectTimeout: NodeJS.Timeout;
    let isCleanup = false;
    let reconnectAttempts = 0;
    // Retry schedule: 3s, 5s, 10s, 30s, then stop
    const retryDelays = [3000, 5000, 10000, 30000];
    const maxReconnectAttempts = retryDelays.length;

    const connect = () => {
      if (isCleanup) return;

      // Close existing connection if any (e.g., when token changes)
      if (wsRef.current && wsRef.current.readyState !== WebSocket.CLOSED) {
        logger.debug('Closing existing WebSocket connection before reconnecting');
        wsRef.current.onclose = null; // Prevent reconnection logic
        wsRef.current.close(1000, 'Reconnecting with new token');
        wsRef.current = null;
      }

      try {
        const token = getToken();
        const wsUrl = token
          ? `${WS_URL}/races/${raceId}/chat/ws?token=${encodeURIComponent(token)}`
          : `${WS_URL}/races/${raceId}/chat/ws`;

        logger.debug('Connecting to chat:', { raceId, hasToken: !!token });
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
          reconnectAttempts = 0; // Reset retry attempts on successful connection

          // Setup ping
          pingInterval = setInterval(() => {
            if (ws?.readyState === WebSocket.OPEN) {
              ws.send(JSON.stringify({ type: 'ping' }));
            }
          }, WEBSOCKET_PING_INTERVAL_MS);
        };

        ws.onerror = (event) => {
          // Only log errors if connection was actually attempted (not during cleanup)
          if (!isCleanup && ws?.readyState !== WebSocket.CLOSED) {
            logger.debug('WebSocket error (may be harmless during navigation):', event);
          }
        };

        ws.onmessage = (event) => {
          if (isCleanup) return;
          
          try {
            // Handle non-string data (Blob, ArrayBuffer, etc.)
            let dataString: string;
            if (typeof event.data === 'string') {
              dataString = event.data;
            } else if (event.data instanceof Blob) {
              // Skip binary data (might be ping/pong or other non-JSON messages)
              logger.debug('Received binary WebSocket message, skipping');
              return;
            } else {
              logger.debug('Received non-string WebSocket message, skipping', { type: typeof event.data });
              return;
            }

            // Backend sends multiple JSON messages separated by newlines in a single frame
            // Split by newline and process each message
            const messages = dataString.split('\n').filter(line => line.trim().length > 0);

            for (const messageStr of messages) {
              const trimmed = messageStr.trim();
              
              // Skip empty messages
              if (!trimmed) {
                continue;
              }

              // Handle ping/pong responses (might not be JSON)
              if (trimmed === 'pong' || trimmed === 'ping') {
                logger.debug('Received ping/pong message');
                continue;
              }

              try {
                const msg: WSMessage = JSON.parse(trimmed);
                
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
                  case 'pong':
                    // Silently handle pong responses
                    logger.debug('Received pong response');
                    break;
                  case 'joined':
                  case 'left':
                    // Silently handle join/leave messages (they're informational)
                    logger.debug(`User ${msg.type}:`, msg.data?.username);
                    break;
                }
              } catch (parseError) {
                // Log parse error for individual message, but continue processing others
                logger.error('Failed to parse individual chat message:', {
                  error: parseError instanceof Error ? parseError.message : String(parseError),
                  messagePreview: trimmed.substring(0, 100),
                  messageLength: trimmed.length
                });
              }
            }
          } catch (e) {
            // Log the actual data for debugging, but truncate if too long
            const dataPreview = typeof event.data === 'string' 
              ? event.data.substring(0, 100) 
              : `[${typeof event.data}]`;
            logger.error('Failed to process WebSocket message:', {
              error: e instanceof Error ? e.message : String(e),
              dataPreview,
              dataLength: typeof event.data === 'string' ? event.data.length : 'N/A'
            });
          }
        };

        ws.onclose = (event) => {
          if (isCleanup) {
            // Component is unmounting, don't try to reconnect or log errors
            return;
          }
          
          setIsConnected(false);
          wsRef.current = null;
          clearInterval(pingInterval);

          // Don't log normal closures (code 1000) or going away (1001) as errors
          if (event.code === 1000 || event.code === 1001) {
            logger.debug('Chat disconnected normally');
            return;
          }

          logger.debug('Chat disconnected:', event.code, event.reason);

          // Retry with scheduled delays: 3s, 5s, 10s, 30s, then stop
          if (reconnectAttempts < maxReconnectAttempts) {
            const delay = retryDelays[reconnectAttempts];
            reconnectAttempts++;
            logger.debug(`Reconnecting in ${delay}ms... (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);
            reconnectTimeout = setTimeout(() => {
              connect();
            }, delay);
          } else {
            logger.warn('Max reconnection attempts reached. Stopping retry.');
            setChatError('Chat connection failed. Please refresh the page to reconnect.');
          }
        };

        ws.onerror = (e) => {
          // Suppress errors during cleanup (component unmounting)
          if (isCleanup) {
            return;
          }
          // Error details usually appear in onclose
          // Only log if connection was actually attempted
          if (ws?.readyState === WebSocket.CONNECTING || ws?.readyState === WebSocket.OPEN) {
            logger.debug('WebSocket error (may be harmless):', e);
          }
        };

      } catch (e) {
        logger.error('Failed to create WebSocket:', e);
        // Retry if creation fails (e.g. URL error) using the same retry schedule
        if (reconnectAttempts < maxReconnectAttempts) {
          const delay = retryDelays[reconnectAttempts];
          reconnectAttempts++;
          logger.debug(`Retrying WebSocket creation in ${delay}ms... (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);
          reconnectTimeout = setTimeout(connect, delay);
        } else {
          logger.warn('Max reconnection attempts reached. Stopping retry.');
          setChatError('Chat connection failed. Please refresh the page to reconnect.');
        }
      }
    };

    connect();

    return () => {
      isCleanup = true;
      logger.debug('Cleaning up chat connection');
      
      // Clear any pending timeouts
      if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
      }
      if (pingInterval) {
        clearInterval(pingInterval);
      }
      
      // Close WebSocket connections gracefully
      const closeSocket = (socket: WebSocket | null) => {
        if (!socket) return;
        
        // Suppress all error events before closing
        socket.onerror = () => {}; // Empty handler to suppress errors
        socket.onclose = () => {}; // Empty handler to suppress close events
        
        // Only close if not already closed or closing
        if (socket.readyState === WebSocket.CONNECTING || socket.readyState === WebSocket.OPEN) {
          // Remove listeners to prevent state updates after unmount
          socket.onmessage = null;
          socket.onopen = null;
          try {
            socket.close(1000, 'Component unmounting'); // Normal closure
          } catch (e) {
            // Ignore errors during cleanup (socket may already be closing)
          }
        }
      };
      
      closeSocket(ws);
      closeSocket(wsRef.current);
      wsRef.current = null;
      wsRef.current = null;
      clearInterval(pingInterval);
      clearTimeout(reconnectTimeout);
      // Note: We don't clear errorTimeout here as it's handled by the parent effect
    };
  }, [raceId, enabled, reconnectTrigger, setChatError, token]);

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
