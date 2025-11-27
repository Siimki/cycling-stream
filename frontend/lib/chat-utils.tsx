/**
 * Chat utility functions
 */

import * as React from 'react';
import { USER_COLORS } from '@/constants/colors';

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

/**
 * Gets a badge element for a username based on hash
 * Returns null if no badge should be shown
 */
export function getUserBadge(username: string): React.ReactElement | null {
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }
  const rand = Math.abs(hash) % 10;
  
  const baseClass = "px-1.5 py-0.5 text-xs font-bold uppercase rounded mr-1.5 align-middle inline-block";
  
  if (rand === 0) return <span className={`${baseClass} bg-green-500/20 text-green-400`}>mod</span>;
  if (rand === 1) return <span className={`${baseClass} bg-pink-500/20 text-pink-400`}>vip</span>;
  if (rand === 2) return <span className={`${baseClass} bg-primary/20 text-primary`}>sub</span>;
  if (rand === 3) return <span className={`${baseClass} bg-amber-500/20 text-amber-400`}>og</span>;
  
  return null;
}

