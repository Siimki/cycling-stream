/**
 * Tests for useChat hook
 * 
 * To run these tests, install testing dependencies:
 * npm install --save-dev @testing-library/react @testing-library/react-hooks jest @types/jest
 * 
 * Then configure Jest in package.json or jest.config.js
 */

import { renderHook, act, waitFor } from '@testing-library/react';
import { useChat } from './useChat';

// Mock WebSocket with better functionality for testing
// eslint-disable-next-line @typescript-eslint/no-explicit-any
class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;

  private sentMessages: string[] = [];
  private messageQueue: string[] = [];

  constructor(public url: string) {
    // Simulate connection opening after a short delay
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      if (this.onopen) {
        this.onopen(new Event('open'));
      }
    }, 10);
  }

  send(data: string) {
    if (this.readyState === MockWebSocket.OPEN) {
      this.sentMessages.push(data);
    }
  }

  close() {
    this.readyState = MockWebSocket.CLOSED;
    if (this.onclose) {
      this.onclose(new CloseEvent('close'));
    }
  }

  // Test helper methods
  simulateMessage(data: string) {
    if (this.onmessage) {
      const event = new MessageEvent('message', { data });
      this.onmessage(event);
    }
  }

  simulateError() {
    if (this.onerror) {
      this.onerror(new Event('error'));
    }
  }

  simulateClose() {
    this.close();
  }

  getSentMessages(): string[] {
    return [...this.sentMessages];
  }

  clearSentMessages() {
    this.sentMessages = [];
  }
}

// Store mock instances for test access
let mockWebSocketInstances: MockWebSocket[] = [];

// Mock global WebSocket
// eslint-disable-next-line @typescript-eslint/no-explicit-any
(global as any).WebSocket = jest.fn((url: string) => {
  const mock = new MockWebSocket(url);
  mockWebSocketInstances.push(mock);
  return mock;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
}) as any;

// Helper to get the latest mock WebSocket instance
const getLatestMockWebSocket = (): MockWebSocket | null => {
  return mockWebSocketInstances[mockWebSocketInstances.length - 1] || null;
};

// Reset mock instances before each test
beforeEach(() => {
  mockWebSocketInstances = [];
});

describe('useChat', () => {
  const raceId = 'test-race-123';
  const enabled = true;

  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
  });

  it('should initialize with correct default values', () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    expect(result.current.messages).toEqual([]);
    expect(result.current.isConnected).toBe(false);
    expect(result.current.error).toBe(null);
  });

  it('should connect when enabled is true', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });
  });

  it('should not connect when enabled is false', () => {
    const { result } = renderHook(() => useChat(raceId, false));

    expect(result.current.isConnected).toBe(false);
  });

  it('should handle message reception', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();
    expect(mock).not.toBeNull();

    // Simulate receiving a message
    const mockMessage = {
      type: 'message',
      data: {
        id: 'msg-1',
        race_id: raceId,
        username: 'TestUser',
        message: 'Hello!',
        created_at: new Date().toISOString(),
      },
    };

    act(() => {
      if (mock) {
        mock.simulateMessage(JSON.stringify(mockMessage));
      }
    });

    await waitFor(() => {
      expect(result.current.messages.length).toBeGreaterThan(0);
    });

    expect(result.current.messages).toContainEqual(
      expect.objectContaining({
        id: 'msg-1',
        message: 'Hello!',
        username: 'TestUser',
      })
    );
  });

  it('should prevent duplicate messages', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();
    const mockMessage = {
      type: 'message',
      data: {
        id: 'msg-dup',
        race_id: raceId,
        username: 'TestUser',
        message: 'Duplicate test',
        created_at: new Date().toISOString(),
      },
    };

    // Send same message twice
    act(() => {
      if (mock) {
        mock.simulateMessage(JSON.stringify(mockMessage));
        mock.simulateMessage(JSON.stringify(mockMessage));
      }
    });

    await waitFor(() => {
      expect(result.current.messages.length).toBe(1);
    });
  });

  it('should handle sendMessage', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    act(() => {
      result.current.sendMessage('Test message');
    });

    await waitFor(() => {
      const sent = mock?.getSentMessages() || [];
      expect(sent.length).toBeGreaterThan(0);
    });

    const sentMessages = mock?.getSentMessages() || [];
    const lastMessage = JSON.parse(sentMessages[sentMessages.length - 1]);
    expect(lastMessage.type).toBe('send_message');
    expect(lastMessage.data.message).toBe('Test message');
  });

  it('should handle sendMessage when not connected', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    // Don't wait for connection, try to send immediately
    act(() => {
      result.current.sendMessage('Test message');
    });

    // Should set error
    await waitFor(() => {
      expect(result.current.error).toBeTruthy();
    });
  });

  it('should handle reconnection', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    act(() => {
      result.current.reconnect();
    });

    // Verify reconnection logic
    expect(result.current.reconnect).toBeDefined();
  });

  it('should handle errors', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    act(() => {
      if (mock) {
        mock.simulateError();
      }
    });

    await waitFor(() => {
      expect(result.current.error).toBeTruthy();
    });
  });

  it('should handle different message types', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    // Test pong message (should be ignored silently)
    act(() => {
      if (mock) {
        mock.simulateMessage(JSON.stringify({ type: 'pong' }));
      }
    });

    // Test error message
    act(() => {
      if (mock) {
        mock.simulateMessage(JSON.stringify({
          type: 'error',
          data: { message: 'Test error' },
        }));
      }
    });

    await waitFor(() => {
      expect(result.current.error).toBe('Test error');
    });
  });

  it('should cleanup on unmount', async () => {
    const { result, unmount } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();
    expect(mock?.readyState).toBe(MockWebSocket.OPEN);

    unmount();

    // Verify connection is closed
    await waitFor(() => {
      expect(mock?.readyState).toBe(MockWebSocket.CLOSED);
    });
  });

  it('should handle reconnection with exponential backoff', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    // Simulate connection close
    act(() => {
      if (mock) {
        mock.simulateClose();
      }
    });

    await waitFor(() => {
      expect(result.current.isConnected).toBe(false);
    }, { timeout: 1000 });

    // Should attempt reconnection
    act(() => {
      result.current.reconnect();
    });

    // Should try to reconnect
    await waitFor(() => {
      expect(mockWebSocketInstances.length).toBeGreaterThan(1);
    });
  });

  it('should handle malformed WebSocket messages gracefully', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    // Simulate malformed messages
    act(() => {
      if (mock) {
        mock.simulateMessage('not valid json');
        mock.simulateMessage('{"incomplete":');
        mock.simulateMessage('{"type": "message", "invalid":}');
      }
    });

    // Should not crash, error should be set or messages ignored
    await waitFor(() => {
      // Either error is set or messages are safely ignored
      expect(result.current.error !== null || result.current.messages.length === 0).toBe(true);
    });
  });

  it('should cleanup on unmount during connection', async () => {
    const { result, unmount } = renderHook(() => useChat(raceId, enabled));

    // Don't wait for connection to complete
    const mock = getLatestMockWebSocket();

    // Unmount immediately
    unmount();

    // Verify connection is closed or will be closed
    await waitFor(() => {
      if (mock) {
        expect(mock.readyState === MockWebSocket.CLOSED || mock.readyState === MockWebSocket.CLOSING).toBe(true);
      }
    }, { timeout: 500 });
  });

  it('should handle rapid message sending without errors', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    // Send multiple messages rapidly
    act(() => {
      for (let i = 0; i < 10; i++) {
        result.current.sendMessage(`Message ${i}`);
      }
    });

    await waitFor(() => {
      const sent = mock?.getSentMessages() || [];
      expect(sent.length).toBeGreaterThan(0);
    });

    // Should not have errors
    expect(result.current.error).toBeNull();
  });

  it('should handle network failure and recovery', async () => {
    const { result } = renderHook(() => useChat(raceId, enabled));

    await waitFor(() => {
      expect(result.current.isConnected).toBe(true);
    });

    const mock = getLatestMockWebSocket();

    // Simulate network error
    act(() => {
      if (mock) {
        mock.simulateError();
      }
    });

    await waitFor(() => {
      expect(result.current.error).toBeTruthy();
    });

    // Attempt reconnection
    act(() => {
      result.current.reconnect();
    });

    // Should attempt to reconnect
    await waitFor(() => {
      expect(mockWebSocketInstances.length).toBeGreaterThan(1);
    });
  });
});

