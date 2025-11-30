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

export default function ChatMessageRow({ message, animate, pulseClass }: ChatMessageRowProps) {
  const userColor = useMemo(() => getUserColor(message.username), [message.username]);
  const badges = useMemo(() => renderUserBadges(message.badges), [message.badges]);

  const segments = useMemo(() => parseSegments(message.message), [message.message]);

  const usernameNode = message.user_id ? (
    <Link
      href={`/users/${message.user_id}`}
      className="font-semibold hover:underline transition-colors"
      style={{ color: userColor }}
    >
      {message.username}
    </Link>
  ) : (
    <span className="font-semibold" style={{ color: userColor }}>
      {message.username}
    </span>
  );

  const usernameWrapperClass = cn('chat-username inline-flex items-center gap-1', {
    vip: message.role === 'vip',
    mod: message.role === 'mod',
    sub: message.role === 'subscriber',
  });

  return (
    <div className={cn(animate && 'chat-message-enter')}>
      <div
        className={cn(
          'leading-snug break-words py-0.5 text-sm text-foreground',
          pulseClass
        )}
        style={{ lineHeight: 1.4 }}
      >
        <div className="flex flex-wrap items-center gap-1.5">
          {badges}
          <span className={usernameWrapperClass}>
            {usernameNode}
          </span>
          <span className="text-muted-foreground font-medium">:</span>
          <span className="flex flex-wrap items-center gap-1 text-[15px] font-medium">
            {segments.map((segment, index) =>
              segment.type === 'emote' ? (
                <ChatEmote
                  key={`${message.id}-emote-${index}`}
                  text={segment.value}
                  special={message.special_emote}
                  disabled={!animate}
                />
              ) : (
                <span key={`${message.id}-text-${index}`} className="whitespace-pre-wrap">
                  {segment.value}
                </span>
              )
            )}
          </span>
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
