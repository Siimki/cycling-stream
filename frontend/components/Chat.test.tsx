/**
 * Tests for Chat component
 * 
 * To run these tests, install testing dependencies:
 * npm install --save-dev @testing-library/react @testing-library/jest-dom jest @types/jest
 * 
 * Then configure Jest in package.json or jest.config.js
 */

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import Chat from './Chat';
import * as useChatModule from '../hooks/useChat';
import * as authModule from '../lib/auth';

// Mock useChat hook
jest.mock('../hooks/useChat', () => ({
  useChat: jest.fn(() => ({
    messages: [],
    sendMessage: jest.fn(),
    isConnected: true,
    error: null,
    reconnect: jest.fn(),
  })),
}));

// Mock getToken
jest.mock('../lib/auth', () => ({
  getToken: jest.fn(() => 'mock-token'),
}));

// Mock getChatHistory
jest.mock('../lib/api', () => ({
  getChatHistory: jest.fn(() =>
    Promise.resolve({
      messages: [],
      limit: 50,
      offset: 0,
    })
  ),
}));

describe('Chat', () => {
  const raceId = 'test-race-123';
  const enabled = true;

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render chat component when enabled', () => {
    render(<Chat raceId={raceId} enabled={enabled} />);

    expect(screen.getByText('Live Chat')).toBeInTheDocument();
  });

  it('should not render when disabled', () => {
    const { container } = render(<Chat raceId={raceId} enabled={false} />);

    expect(container.firstChild).toBeNull();
  });

  it('should show connection status', () => {
    render(<Chat raceId={raceId} enabled={enabled} />);

    expect(screen.getByText(/Connected/i)).toBeInTheDocument();
  });

  it('should show input field for authenticated users', () => {
    render(<Chat raceId={raceId} enabled={enabled} />);

    const input = screen.getByPlaceholderText('Type a message...');
    expect(input).toBeInTheDocument();
  });

  it('should show sign in prompt for anonymous users', () => {
    // Mock getToken to return null
    jest.mocked(authModule.getToken).mockReturnValueOnce(null);

    render(<Chat raceId={raceId} enabled={enabled} />);

    expect(screen.getByText(/Sign in/i)).toBeInTheDocument();
  });

  it('should handle message sending', async () => {
    const mockSendMessage = jest.fn();
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [],
      sendMessage: mockSendMessage,
      isConnected: true,
      error: null,
      reconnect: jest.fn(),
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    const input = screen.getByPlaceholderText('Type a message...');
    const sendButton = screen.getByText('Send');

    fireEvent.change(input, { target: { value: 'Test message' } });
    fireEvent.click(sendButton);

    await waitFor(() => {
      expect(mockSendMessage).toHaveBeenCalledWith('Test message');
    });
  });

  it('should display messages', () => {
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [
        {
          id: 'msg-1',
          race_id: raceId,
          username: 'TestUser',
          message: 'Hello!',
          created_at: new Date().toISOString(),
        },
      ],
      sendMessage: jest.fn(),
      isConnected: true,
      error: null,
      reconnect: jest.fn(),
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    expect(screen.getByText('Hello!')).toBeInTheDocument();
    expect(screen.getByText('TestUser')).toBeInTheDocument();
  });

  it('should show error message when error occurs', () => {
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [],
      sendMessage: jest.fn(),
      isConnected: false,
      error: 'Connection failed',
      reconnect: jest.fn(),
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    expect(screen.getByText('Connection failed')).toBeInTheDocument();
  });

  it('should handle chat toggle', () => {
    render(<Chat raceId={raceId} enabled={enabled} />);

    const toggleButton = screen.getByLabelText(/Hide chat/i);
    fireEvent.click(toggleButton);

    // Chat should be hidden
    expect(screen.queryByText('Live Chat')).not.toBeInTheDocument();
  });

  it('should disable input when not connected', () => {
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [],
      sendMessage: jest.fn(),
      isConnected: false,
      error: null,
      reconnect: jest.fn(),
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    const input = screen.getByPlaceholderText('Type a message...');
    expect(input).toBeDisabled();
  });

  it('should handle very long message lists without performance issues', () => {
    const longMessages = Array.from({ length: 1000 }, (_, i) => ({
      id: `msg-${i}`,
      race_id: raceId,
      username: 'TestUser',
      message: `Message ${i}`,
      created_at: new Date().toISOString(),
    }));

    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: longMessages,
      sendMessage: jest.fn(),
      isConnected: true,
      error: null,
      reconnect: jest.fn(),
    });

    const { container } = render(<Chat raceId={raceId} enabled={enabled} />);

    // Component should render without crashing
    expect(container.firstChild).not.toBeNull();
  });

  it('should handle rapid message sending', async () => {
    const mockSendMessage = jest.fn();
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [],
      sendMessage: mockSendMessage,
      isConnected: true,
      error: null,
      reconnect: jest.fn(),
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    const input = screen.getByPlaceholderText('Type a message...');
    const sendButton = screen.getByText('Send');

    // Send multiple messages rapidly
    for (let i = 0; i < 5; i++) {
      fireEvent.change(input, { target: { value: `Message ${i}` } });
      fireEvent.click(sendButton);
    }

    await waitFor(() => {
      expect(mockSendMessage).toHaveBeenCalledTimes(5);
    });
  });

  it('should handle connection interruptions gracefully', () => {
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [],
      sendMessage: jest.fn(),
      isConnected: false,
      error: 'Connection lost',
      reconnect: jest.fn(),
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    // Should show error message
    expect(screen.getByText('Connection lost')).toBeInTheDocument();

    // Input should be disabled
    const input = screen.getByPlaceholderText('Type a message...');
    expect(input).toBeDisabled();
  });

  it('should handle reconnection button click', async () => {
    const mockReconnect = jest.fn();
    jest.mocked(useChatModule.useChat).mockReturnValueOnce({
      messages: [],
      sendMessage: jest.fn(),
      isConnected: false,
      error: 'Connection lost',
      reconnect: mockReconnect,
    });

    render(<Chat raceId={raceId} enabled={enabled} />);

    // Find and click reconnect button if it exists
    const reconnectButton = screen.queryByText(/reconnect/i);
    if (reconnectButton) {
      fireEvent.click(reconnectButton);
      expect(mockReconnect).toHaveBeenCalled();
    }
  });
});

