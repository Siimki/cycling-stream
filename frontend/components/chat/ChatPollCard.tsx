'use client';

import { useMemo, useState } from 'react';
import { ChatPoll } from '@/lib/api';
import { cn } from '@/lib/utils';

interface ChatPollCardProps {
  poll: ChatPoll;
  onVote: (optionId: string) => Promise<void>;
  disabled?: boolean;
  loading?: boolean;
  animate?: boolean;
  onDismiss?: () => void;
  headline?: string;
}

export default function ChatPollCard({
  poll,
  onVote,
  disabled,
  loading,
  animate,
  onDismiss,
  headline,
}: ChatPollCardProps) {
  const [pendingOption, setPendingOption] = useState<string | null>(null);
  const [localError, setLocalError] = useState<string | null>(null);

  const totalVotes = useMemo(() => {
    if (poll.total_votes) {
      return poll.total_votes;
    }
    return poll.options.reduce((sum, option) => sum + option.votes, 0);
  }, [poll]);

  const handleVote = async (optionId: string) => {
    if (disabled || poll.closed) return;
    try {
      setPendingOption(optionId);
      setLocalError(null);
      await onVote(optionId);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to submit vote';
      setLocalError(message);
    } finally {
      setPendingOption(null);
    }
  };

  const closesLabel = useMemo(() => {
    if (!poll.closes_at || poll.closed) {
      return null;
    }
    const closesAt = new Date(poll.closes_at);
    const diffMs = closesAt.getTime() - Date.now();
    if (diffMs <= 0) {
      return 'closing…';
    }
    const seconds = Math.max(5, Math.round(diffMs / 1000));
    return `closing in ${seconds}s`;
  }, [poll]);

  return (
    <div
      className={cn(
        'bg-card/70 border border-border/50 rounded-lg p-4 space-y-3 shadow-sm',
        animate && 'motion-slide-in-up'
      )}
    >
      <div className="flex items-start justify-between gap-3">
        <div>
          {headline && (
            <p className="text-xs uppercase tracking-widest text-muted-foreground mb-1">{headline}</p>
          )}
          <h4 className="text-sm font-semibold text-foreground">{poll.question}</h4>
          <p className="text-xs text-muted-foreground">
            {poll.closed ? 'Poll closed' : closesLabel || 'Live poll'}
          </p>
        </div>
        {onDismiss && (
          <button
            type="button"
            aria-label="Dismiss poll"
            onClick={onDismiss}
            className="text-muted-foreground/70 hover:text-foreground text-xs"
          >
            ✕
          </button>
        )}
      </div>

      <div className="space-y-2">
        {poll.options.map((option) => {
          const percent = totalVotes === 0 ? 0 : Math.round((option.votes / totalVotes) * 100);
          const width = Math.max(percent, option.votes > 0 ? 6 : 0);
          const isPending = pendingOption === option.id || loading;

          return (
            <button
              key={option.id}
              type="button"
              onClick={() => handleVote(option.id)}
              disabled={poll.closed || disabled || loading}
              className={cn(
                'w-full text-left border border-border/40 rounded-md p-2 transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40',
                poll.closed
                  ? 'bg-muted/40 cursor-default'
                  : 'hover:border-primary/60 hover:bg-primary/5'
              )}
            >
              <div className="flex items-center justify-between text-sm font-medium text-foreground">
                <span>{option.label}</span>
                <span className="text-xs text-muted-foreground">{percent}%</span>
              </div>
              <div className="mt-2 h-2 bg-muted rounded-full overflow-hidden">
                <div
                  className={cn(
                    'h-full bg-primary',
                    animate && 'transition-all duration-500 ease-out'
                  )}
                  style={{ width: `${width}%` }}
                />
              </div>
              {isPending && !poll.closed && (
                <p className="text-[11px] text-primary mt-1">Submitting…</p>
              )}
            </button>
          );
        })}
      </div>

      <div className="flex items-center justify-between text-xs text-muted-foreground">
        <span>{totalVotes} votes</span>
        {poll.closed && <span>Thanks for voting!</span>}
      </div>

      {localError && (
        <p className="text-xs text-destructive">{localError}</p>
      )}
    </div>
  );
}

