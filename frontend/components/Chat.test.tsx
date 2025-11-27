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

  it('should show character count', () => {
    render(<Chat raceId={raceId} enabled={enabled} />);

    const input = screen.getByPlaceholderText('Type a message...');
    fireEvent.change(input, { target: { value: 'Test' } });

    expect(screen.getByText(/4\/500 characters/i)).toBeInTheDocument();
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
});

