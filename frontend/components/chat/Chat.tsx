'use client';

import { useEffect, useRef, useState, useCallback, useMemo } from 'react';
import { useChatContext } from '@/components/chat/ChatProvider';
import { useAuth } from '@/contexts/AuthContext';
import { useExperience } from '@/contexts/ExperienceContext';
import { useSound } from '@/components/providers/SoundProvider';
import { getChatHistory, ChatMessage } from '@/lib/api';
import { CHAT_HISTORY_LIMIT, CHAT_MESSAGE_MAX_LENGTH } from '@/constants/intervals';
import { Send, Settings } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useHudStats } from '@/components/user/HudStatsProvider';
import ChatMessageRow from '@/components/chat/ChatMessageRow';
import ChatPollCard from '@/components/chat/ChatPollCard';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('Chat');

export default function Chat() {
  const {
    messages: wsMessages,
    sendMessage,
    isConnected,
    error,
    reconnect,
    raceId,
    enabled,
    activePoll,
    lastClosedPoll,
    voteInPoll,
    pollVoteLoading,
  } = useChatContext();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [historyLoaded, setHistoryLoaded] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);
  const { isAuthenticated, user } = useAuth();
  const { bonusReady, claimBonus } = useHudStats();
  const { resolvedUIPreferences, uiPreferences, updateUIPreferences } = useExperience();
  const animationsEnabled = resolvedUIPreferences.chat_animations && !resolvedUIPreferences.reduced_motion;
  const pollAnimationsEnabled = resolvedUIPreferences.poll_animations && !resolvedUIPreferences.reduced_motion;
  const [dismissedPollId, setDismissedPollId] = useState<string | null>(null);
  const { play } = useSound();
  const pollSoundRef = useRef<string | null>(null);
  const mentionSoundRef = useRef<string | null>(null);
  const mentionTargets = useMemo(() => {
    const targets = new Set<string>();
    if (user?.name) {
      targets.add(user.name.replace(/\s+/g, '').toLowerCase());
    }
    if (user?.email) {
      targets.add(user.email.split('@')[0].toLowerCase());
    }
    return Array.from(targets).filter(Boolean);
  }, [user]);

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

  useEffect(() => {
    if (historyLoaded && wsMessages.length > 0) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
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
  }, [historyLoaded, wsMessages]);

  useEffect(() => {
    if (activePoll) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setDismissedPollId(null);
    }
  }, [activePoll]);

  useEffect(() => {
    if (!lastClosedPoll) {
      return;
    }
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setDismissedPollId(null);
    const timer = setTimeout(() => {
      setDismissedPollId(lastClosedPoll.id);
    }, 8000);
    return () => clearTimeout(timer);
  }, [lastClosedPoll]);

  useEffect(() => {
    if (!activePoll || pollSoundRef.current === activePoll.id) {
      return;
    }
    play('notification');
    pollSoundRef.current = activePoll.id;
  }, [activePoll, play]);

  useEffect(() => {
    if (!messages.length || mentionTargets.length === 0) {
      return;
    }
    const latest = messages[messages.length - 1];
    if (!latest || (user?.id && latest.user_id === user.id)) {
      return;
    }
    const content = latest.message.toLowerCase();
    const mentioned = mentionTargets.some((target) => target && content.includes(`@${target}`));
    if (mentioned && mentionSoundRef.current !== latest.id) {
      play('chat-mention');
      mentionSoundRef.current = latest.id;
    }
  }, [messages, mentionTargets, play, user?.id]);

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

  const getPulseClass = useCallback(
    (msg: ChatMessage) => {
      if (!animationsEnabled || !msg.role) {
        return undefined;
      }
      const createdAt = new Date(msg.created_at).getTime();
      if (Date.now() - createdAt > 2000) {
        return undefined;
      }
      switch (msg.role) {
        case 'mod':
          return 'chat-role-mod';
        case 'vip':
          return 'chat-role-vip';
        case 'subscriber':
          return 'chat-role-sub';
        default:
          return undefined;
      }
    },
    [animationsEnabled]
  );

  const handlePollVote = useCallback(
    async (optionId: string) => {
      if (!activePoll) {
        return;
      }
      await voteInPoll(activePoll.id, optionId);
    },
    [activePoll, voteInPoll]
  );

  const toggleAnimations = useCallback(() => {
    updateUIPreferences({
      chat_animations: !uiPreferences.chat_animations,
    });
  }, [uiPreferences.chat_animations, updateUIPreferences]);

  const visibleClosedPoll = useMemo(() => {
    if (!lastClosedPoll) {
      return null;
    }
    if (dismissedPollId && dismissedPollId === lastClosedPoll.id) {
      return null;
    }
    return lastClosedPoll;
  }, [lastClosedPoll, dismissedPollId]);


  if (!enabled) {
    return null;
  }

  return (
    <div className="flex flex-col h-full" style={{ backgroundColor: 'var(--design-surface)' }}>
      {/* Header */}
      <div className="px-5 py-4 border-b border-border/50 flex items-center justify-between shrink-0">
        <div className="flex items-center gap-2.5">
          <span className="text-lg font-semibold text-foreground tracking-tight">Stream Chat</span>
          <div className="flex items-center gap-2">
            {isConnected ? (
                <span className="w-2 h-2 bg-connected rounded-full animate-pulse shadow-glow-green" title="Connected"></span>
            ) : (
                <span className="w-2 h-2 bg-disconnected rounded-full" title="Disconnected"></span>
            )}
          </div>
        </div>
        <Button
          variant="outline"
          size="sm"
          className="h-8 gap-2 text-xs text-muted-foreground hover:text-foreground"
          onClick={toggleAnimations}
        >
          <Settings className="w-4 h-4" />
          {uiPreferences.chat_animations ? 'Animations On' : 'Animations Off'}
        </Button>
      </div>

      {/* Messages - Fixed height with scroll */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto chat-scroll px-5 py-4 min-h-0 space-y-3">
        {isLoading ? (
            <div className="flex items-center justify-center h-full text-muted-foreground text-base">Loading...</div>
        ) : (
            <div className="space-y-3">
              {activePoll && (
                <ChatPollCard
                  poll={activePoll}
                  onVote={handlePollVote}
                  loading={pollVoteLoading}
                  animate={pollAnimationsEnabled}
                  headline="Live poll"
                />
              )}
              {!activePoll && visibleClosedPoll && (
                <ChatPollCard
                  poll={visibleClosedPoll}
                  onVote={async () => {}}
                  disabled
                  animate={pollAnimationsEnabled}
                  headline="Poll results"
                  onDismiss={() => setDismissedPollId(visibleClosedPoll.id)}
                />
              )}
              {messages.map((msg) => (
                <ChatMessageRow
                  key={msg.id}
                  message={msg}
                  animate={animationsEnabled}
                  pulseClass={getPulseClass(msg)}
                />
              ))}
              {messages.length === 0 && (
                <div className="text-center text-muted-foreground text-lg mt-8">No messages yet. Say hello! 👋</div>
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
                {/* Bonus section - Always visible (enabled or disabled state) */}
                <div className="mt-3 rounded-lg border border-border/50 bg-muted/30 p-3">
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="text-sm font-semibold text-foreground">Watch bonus</p>
                      <p className="text-xs text-muted-foreground/80">
                        {bonusReady ? "Ready to collect" : "Keep watching to unlock"}
                      </p>
                    </div>
                    <Button
                      size="sm"
                      className={`h-9 px-4 text-xs font-semibold border-0 transition-all duration-300 ${
                        bonusReady 
                          ? "bg-success hover:bg-success/90 text-success-foreground shadow-lg shadow-success/20" 
                          : "bg-muted text-muted-foreground cursor-not-allowed opacity-70"
                      }`}
                      onClick={claimBonus}
                      disabled={!bonusReady}
                    >
                      {bonusReady ? "Claim +50" : "Locked"}
                    </Button>
                  </div>
                </div>
            </>
        )}
      </div>
    </div>
  );
}
