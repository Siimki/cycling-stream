/**
 * Chat utility functions
 */

import * as React from 'react';
import { USER_COLORS } from '@/constants/colors';
import { cn } from '@/lib/utils';

/**
 * Gets a consistent color for a username based on hash
 */
export function getUserColor(username: string): string {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }
  const index = Math.abs(hash) % USER_COLORS.length;
  return USER_COLORS[index];
}

const BADGE_STYLES: Record<string, string> = {
  mod: 'bg-emerald-500/15 text-emerald-300',
  vip: 'bg-purple-500/20 text-purple-300',
  sub: 'bg-primary/20 text-primary',
  og: 'bg-amber-500/20 text-amber-300',
};

export function renderUserBadges(badges?: string[]): React.ReactElement[] {
  if (!badges || badges.length === 0) {
    return [];
  }

  return badges.map((badge, index) => (
    <span
      key={`${badge}-${index}`}
      className={cn(
        'px-1.5 py-0.5 text-[10px] font-bold uppercase rounded mr-1.5 align-middle inline-block tracking-wide',
        BADGE_STYLES[badge] ?? 'bg-muted/60 text-muted-foreground'
      )}
    >
      {badge}
    </span>
  ));
}

