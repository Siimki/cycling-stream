'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import Link from 'next/link';
import { useChatContext } from '@/components/chat/ChatProvider';
import { useAuth } from '@/contexts/AuthContext';
import { getChatHistory, ChatMessage } from '@/lib/api';
import { CHAT_HISTORY_LIMIT, CHAT_MESSAGE_MAX_LENGTH } from '@/constants/intervals';
import { Send, Settings } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useHudStats } from "@/components/user/HudStatsProvider"
import { getUserColor, getUserBadge } from '@/lib/chat-utils';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('Chat');

export default function Chat() {
  const { messages: wsMessages, sendMessage, isConnected, error, reconnect, raceId, enabled } = useChatContext();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [historyLoaded, setHistoryLoaded] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);
  const { isAuthenticated } = useAuth();
  const { bonusReady, claimBonus } = useHudStats()

  // Load chat history on mount (only once)
  useEffect(() => {
    if (enabled && raceId && !historyLoaded) {
      getChatHistory(raceId, CHAT_HISTORY_LIMIT, 0)
        .then((response) => {
          // Messages are already in chronological order from API
          setMessages(response.messages);
          setHistoryLoaded(true);
          setIsLoading(false);
        })
        .catch((err) => {
          logger.error('Failed to load chat history:', err);
          setHistoryLoaded(true);
          setIsLoading(false);
        });
    }
  }, [enabled, raceId, historyLoaded]);

  // Merge WebSocket messages with loaded history
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {
    if (historyLoaded && wsMessages.length > 0) {
      setMessages((prev) => {
        const existingIds = new Set(prev.map((m) => m.id));
        const newMessages = wsMessages.filter((m) => !existingIds.has(m.id));
        if (newMessages.length > 0) {
          logger.debug('Merging new WebSocket messages:', newMessages.length);
          return [...prev, ...newMessages].sort(
            (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
          );
        }
        return prev;
      });
    }
    // wsMessages is intentionally not in deps - we want to merge when it changes
  }, [historyLoaded, wsMessages]);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSend = useCallback(
    () => { // Changed signature to match button onClick
      if (!inputValue.trim() || !isAuthenticated) {
        return;
      }

      sendMessage(inputValue.trim());
      setInputValue('');
    },
    [inputValue, sendMessage, isAuthenticated]
  );


  if (!enabled) {
    return null;
  }

  return (
    <div className="flex flex-col h-full" style={{ backgroundColor: 'var(--design-surface)' }}>
      {/* Header */}
      <div className="px-5 py-4 border-b border-border/50 flex items-center justify-between shrink-0">
        <div className="flex items-center gap-2.5">
          <span className="text-base font-semibold text-foreground tracking-tight">Stream Chat</span>
          <div className="flex items-center gap-2">
            {isConnected ? (
                <span className="w-2 h-2 bg-connected rounded-full animate-pulse shadow-glow-green" title="Connected"></span>
            ) : (
                <span className="w-2 h-2 bg-disconnected rounded-full" title="Disconnected"></span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-foreground">
            <Settings className="w-4 h-4" />
          </Button>
        </div>
      </div>

      {/* Messages - Fixed height with scroll */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto chat-scroll px-5 py-4 min-h-0">
        {isLoading ? (
            <div className="flex items-center justify-center h-full text-muted-foreground text-base">Loading...</div>
        ) : (
            <div className="space-y-4">
            {messages.map((msg) => (
                <div key={msg.id} className="leading-relaxed break-words transition-colors" style={{ lineHeight: '1.6' }}>
                {getUserBadge(msg.username)}
                {msg.user_id ? (
                    <Link 
                      href={`/users/${msg.user_id}`}
                      className="text-sm font-semibold cursor-pointer hover:underline mr-1.5 text-primary" 
                      style={{ color: getUserColor(msg.username) }}
                    >
                      {msg.username}
                    </Link>
                ) : (
                    <span className="text-sm font-semibold mr-1.5 text-primary" style={{ color: getUserColor(msg.username) }}>
                      {msg.username}
                    </span>
                )}
                <span className="text-muted-foreground text-sm mr-1.5">:</span>
                <span className="text-base text-foreground font-normal">{msg.message}</span>
                </div>
            ))}
            {messages.length === 0 && (
                <div className="text-center text-muted-foreground text-base mt-8">No messages yet. Say hello! 👋</div>
            )}
            </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="px-5 py-3 bg-destructive/10 border-t border-destructive/20 text-sm font-medium text-destructive flex justify-between items-center shrink-0">
           <span>{error}</span>
           {!isConnected && (
              <button onClick={reconnect} className="underline hover:text-destructive/80 font-bold">Reconnect</button>
           )}
        </div>
      )}

      {/* Input */}
      <div className="px-5 py-4 border-t border-border/50 shrink-0">
        {!isAuthenticated ? (
            <div className="text-center py-3 text-base text-muted-foreground">
                 <a href="/auth/login" className="text-primary hover:underline font-medium">Sign in</a> to chat
            </div>
        ) : (
            <>
                <div className="flex items-center gap-2">
                <Input
                    placeholder="Send a message..."
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleSend()}
                    disabled={!isConnected}
                    className="flex-1 h-10 bg-background border border-border/50 rounded-lg text-base placeholder:text-muted-foreground/70 focus-visible:ring-2 focus-visible:ring-primary/50 focus-visible:border-primary/50"
                />
                <Button 
                    size="icon" 
                    className="h-10 w-10 bg-primary text-primary-foreground hover:bg-primary/90 shrink-0 rounded-lg" 
                    onClick={handleSend}
                    disabled={!isConnected || !inputValue.trim()}
                >
                    <Send className="w-4 h-4" />
                </Button>
                </div>
                <div className="flex justify-end mt-2 px-0.5">
                  <p className="text-xs text-muted-foreground/60 font-medium">
                      {inputValue.length} / {CHAT_MESSAGE_MAX_LENGTH}
                  </p>
                </div>
                {bonusReady && (
                  <div className="mt-3 rounded-lg border border-border/50 bg-muted/30 p-3">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <p className="text-sm font-semibold text-foreground">Watch bonus</p>
                        <p className="text-xs text-muted-foreground/80">Ready to collect</p>
                      </div>
                      <Button
                        size="sm"
                        className="h-9 px-4 text-xs font-semibold bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-600 hover:to-amber-700 text-black shadow-lg shadow-amber-500/25 border-0"
                        onClick={claimBonus}
                      >
                        Claim +50
                      </Button>
                    </div>
                  </div>
                )}
            </>
        )}
      </div>
    </div>
  );
}
