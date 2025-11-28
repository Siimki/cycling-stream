'use client';

import { useMemo } from 'react';
import type { UserAchievement } from '@/lib/api';

interface AchievementToastProps {
  achievement: UserAchievement;
  onDismiss?: () => void;
}

export function AchievementToast({ achievement, onDismiss }: AchievementToastProps) {
  const icon = useMemo(() => achievement.icon ?? '⭐', [achievement.icon]);
  return (
    <div className="achievement-toast motion-slide-in-up">
      <div className="achievement-toast-icon" aria-hidden="true">
        {icon}
      </div>
      <div className="achievement-toast-body">
        <p className="achievement-toast-label">Achievement unlocked</p>
        <p className="achievement-toast-title">{achievement.title}</p>
        {achievement.description && (
          <p className="achievement-toast-desc">{achievement.description}</p>
        )}
        <p className="achievement-toast-points">+{achievement.points} pts</p>
      </div>
      {onDismiss && (
        <button
          type="button"
          className="achievement-toast-close"
          aria-label="Dismiss achievement"
          onClick={onDismiss}
        >
          ✕
        </button>
      )}
    </div>
  );
}

