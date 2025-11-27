'use client';

import { createContext, useContext, ReactNode } from 'react';
import { useChat } from '@/hooks/useChat';
import { ChatMessage } from '@/lib/api';

interface ChatContextType {
  messages: ChatMessage[];
  sendMessage: (message: string) => void;
  isConnected: boolean;
  error: string | null;
  reconnect: () => void;
  raceId: string;
  enabled: boolean;
}

const ChatContext = createContext<ChatContextType | undefined>(undefined);

interface ChatProviderProps {
  raceId: string;
  enabled: boolean;
  children: ReactNode;
}

export function ChatProvider({ raceId, enabled, children }: ChatProviderProps) {
  const chat = useChat(raceId, enabled);

  return (
    <ChatContext.Provider value={{ ...chat, raceId, enabled }}>
      {children}
    </ChatContext.Provider>
  );
}

export function useChatContext() {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error('useChatContext must be used within a ChatProvider');
  }
  return context;
}

