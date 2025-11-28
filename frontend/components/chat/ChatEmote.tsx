'use client';

import { cn } from '@/lib/utils';

interface ChatEmoteProps {
  text: string;
  special?: boolean;
  disabled?: boolean;
}

export default function ChatEmote({ text, special, disabled }: ChatEmoteProps) {
  return (
    <span
      className={cn(
        'chat-emote px-1 font-semibold uppercase tracking-wide',
        special ? 'chat-emote-special' : 'chat-emote-basic',
        disabled && 'chat-emote-static'
      )}
    >
      {text}
    </span>
  );
}

