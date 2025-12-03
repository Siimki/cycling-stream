'use client';

import Link from 'next/link';
import { useMemo } from 'react';
import { ChatMessage } from '@/lib/api';
import { getUserColor, renderUserBadges } from '@/lib/chat-utils';
import { cn } from '@/lib/utils';
import ChatEmote from './ChatEmote';

interface ChatMessageRowProps {
  message: ChatMessage;
  animate?: boolean;
  pulseClass?: string;
  index?: number;
}

const KNOWN_EMOTES = new Set([
  ':bike:',
  ':fire:',
  ':zap:',
  ':bolt:',
  ':crown:',
  ':clap:',
  ':rocket:',
  ':podium:',
  ':star:',
]);

type Segment = {
  type: 'text' | 'emote';
  value: string;
};

export default function ChatMessageRow({ message, animate, pulseClass, index }: ChatMessageRowProps) {
  const userColor = useMemo(() => getUserColor(message.username), [message.username]);
  const badges = useMemo(() => renderUserBadges(message.badges), [message.badges]);
  const primaryBadge = badges?.[0];

  const segments = useMemo(() => parseSegments(message.message), [message.message]);
  const isAltRow = typeof index === 'number' && index % 2 === 1;

  const usernameNode = message.user_id ? (
    <Link
      href={`/users/${message.user_id}`}
      className="font-semibold transition-colors"
      style={{ color: userColor }}
    >
      {message.username}
    </Link>
  ) : (
    <span className="font-semibold" style={{ color: userColor }}>
      {message.username}
    </span>
  );

  const usernameWrapperClass = cn('inline-flex items-center gap-1');

  const roleStyles: Record<string, string> = {
    mod: 'bg-emerald-500/15 text-emerald-200 border border-emerald-400/35',
    vip: 'bg-purple-500/15 text-purple-200 border border-purple-400/35',
    subscriber: 'bg-amber-500/15 text-amber-100 border border-amber-300/35',
  };

  const roleLabels: Record<string, string> = {
    mod: 'Mod',
    vip: 'Pro',
    subscriber: 'Member',
  };

  const rolePill = message.role ? (
    <span
      className={cn(
        'text-[11px] px-2 py-0.5 rounded-full font-semibold tracking-wide uppercase leading-tight',
        roleStyles[message.role] || 'bg-muted/20 text-foreground border border-border/50'
      )}
    >
      {roleLabels[message.role] || message.role}
    </span>
  ) : null;

  return (
    <div className={cn(animate && 'chat-message-enter')}>
      <div
        className={cn(
          'break-words rounded-xl border border-border/30 px-3.5 py-2.5 text-sm shadow-[0_1px_0_rgba(255,255,255,0.02)] transition-colors',
          isAltRow ? 'bg-muted/10' : 'bg-background/60',
          pulseClass
        )}
        style={{ lineHeight: 1.45 }}
      >
        <div className="flex flex-wrap items-center gap-1.5 text-[15px] font-medium leading-relaxed text-foreground/90">
          {primaryBadge}
          <span className={usernameWrapperClass}>{usernameNode}</span>
          {rolePill}
          <span className="text-muted-foreground font-semibold">:</span>
          {segments.map((segment, idx) =>
            segment.type === 'emote' ? (
              <ChatEmote
                key={`${message.id}-emote-${idx}`}
                text={segment.value}
                special={message.special_emote}
                disabled={!animate}
              />
            ) : (
              <span key={`${message.id}-text-${idx}`} className="whitespace-pre-wrap">
                {segment.value}
              </span>
            )
          )}
        </div>
      </div>
    </div>
  );
}

function parseSegments(message: string): Segment[] {
  if (!message) {
    return [];
  }

  const tokens = message.split(/(\s+)/);
  const segments: Segment[] = [];

  tokens.forEach((token) => {
    if (!token) {
      return;
    }

    if (token.trim() === '') {
      segments.push({ type: 'text', value: token });
      return;
    }

    const normalized = token.trim().toLowerCase();
    if (KNOWN_EMOTES.has(normalized)) {
      segments.push({ type: 'emote', value: token.trim() });
    } else {
      segments.push({ type: 'text', value: token });
    }
  });

  return segments;
}
