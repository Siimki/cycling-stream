'use client';

import { useEffect, useState } from 'react';
import { getUserXP, getUserWeekly, type XPProgress, type WeeklyGoalProgress } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';

interface UserStatsProps {
  compact?: boolean; // For header display
}

export default function UserStats({ compact = false }: UserStatsProps) {
  const { user, isAuthenticated } = useAuth();
  const [xpData, setXpData] = useState<XPProgress | null>(null);
  const [weeklyData, setWeeklyData] = useState<WeeklyGoalProgress | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Only fetch stats if user is authenticated
    if (!isAuthenticated || !user) {
      setLoading(false);
      return;
    }

    async function fetchStats() {
      try {
        const [xp, weekly] = await Promise.all([
          getUserXP().catch((err: any) => {
            // Silently fail for auth errors or 404s (endpoint might not exist yet)
            if (err?.status === 401 || err?.status === 404) return null;
            // Only log unexpected errors
            if (err?.status !== 401 && err?.status !== 404) {
              console.error('Failed to fetch XP:', err);
            }
            return null;
          }),
          getUserWeekly().catch((err: any) => {
            // Silently fail for auth errors or 404s (endpoint might not exist yet)
            if (err?.status === 401 || err?.status === 404) return null;
            // Only log unexpected errors
            if (err?.status !== 401 && err?.status !== 404) {
              console.error('Failed to fetch weekly stats:', err);
            }
            return null;
          }),
        ]);
        setXpData(xp);
        setWeeklyData(weekly);
      } catch (error: any) {
        // Only log unexpected errors
        if (error?.status !== 401 && error?.status !== 404) {
          console.error('Failed to fetch user stats:', error);
        }
      } finally {
        setLoading(false);
      }
    }

    fetchStats();
    // Refresh stats every 30 seconds
    const interval = setInterval(fetchStats, 30000);
    return () => clearInterval(interval);
  }, [isAuthenticated, user]);

  if (loading) {
    return (
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <div className="h-4 w-16 animate-pulse bg-muted rounded"></div>
      </div>
    );
  }

  if (!xpData) {
    return null;
  }

  // Calculate progress percentage: current XP in level / total XP needed for next level
  const progressPercent = xpData.xp_for_next_level > 0
    ? ((xpData.progress_in_current_level / xpData.xp_for_next_level) * 100)
    : 100;

  if (compact) {
    // Compact header display
    return (
      <div className="flex items-center gap-3 text-sm">
        {/* Level Badge */}
        <div className="flex items-center gap-1.5 px-2 py-1 bg-primary/20 rounded-md">
          <span className="font-semibold text-primary">Lv.{xpData.level}</span>
        </div>

        {/* XP Bar */}
        <div className="flex items-center gap-2 min-w-[120px]">
          <div className="flex-1 h-2 bg-muted rounded-full overflow-hidden">
            <div
              className="h-full bg-success transition-all duration-300"
              style={{ width: `${Math.min(progressPercent, 100)}%` }}
            />
          </div>
          <span className="text-xs text-muted-foreground whitespace-nowrap">
            {xpData.progress_in_current_level}/{xpData.xp_for_next_level}
          </span>
        </div>

        {/* Points */}
        {user && (
          <div className="flex items-center gap-1 px-2 py-1 bg-primary/20 rounded-md">
            <span className="text-primary font-medium">
              {user.points || 0} pts
            </span>
          </div>
        )}

        {/* Streak */}
        {weeklyData && weeklyData.current_streak_weeks > 0 && (
          <div className="flex items-center gap-1 px-2 py-1 bg-warning/20 rounded-md">
            <span className="text-warning">🔥</span>
            <span className="text-warning font-medium">
              {weeklyData.current_streak_weeks}w
            </span>
          </div>
        )}
      </div>
    );
  }

  // Full display (for profile page)
  return (
    <div className="space-y-4">
      {/* Level and XP */}
      <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-lg font-semibold text-foreground/95">Level {xpData.level}</h3>
          <div className="text-sm text-muted-foreground">
            {xpData.xp_total} XP
          </div>
        </div>
        <div className="space-y-1">
          <div className="flex justify-between text-sm text-muted-foreground">
            <span>Progress to Level {xpData.level + 1}</span>
            <span>{xpData.progress_in_current_level} / {xpData.xp_for_next_level} XP</span>
          </div>
          <div className="h-3 bg-muted rounded-full overflow-hidden">
            <div
              className="h-full bg-success transition-all duration-300"
              style={{ width: `${Math.min(progressPercent, 100)}%` }}
            />
          </div>
        </div>
      </div>

      {/* Weekly Goal */}
      {weeklyData && (
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4">
          <h3 className="text-lg font-semibold text-foreground/95 mb-3">This Week</h3>
          <div className="space-y-3">
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-muted-foreground">Watch Time</span>
                <span className="font-medium text-foreground">{weeklyData.watch_minutes} / 30 min</span>
              </div>
              <div className="h-2 bg-muted rounded-full overflow-hidden">
                <div
                  className="h-full bg-success"
                  style={{ width: `${Math.min((weeklyData.watch_minutes / 30) * 100, 100)}%` }}
                />
              </div>
            </div>
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-muted-foreground">Chat Messages</span>
                <span className="font-medium text-foreground">{weeklyData.chat_messages} / 3</span>
              </div>
              <div className="h-2 bg-muted rounded-full overflow-hidden">
                <div
                  className="h-full bg-success"
                  style={{ width: `${Math.min((weeklyData.chat_messages / 3) * 100, 100)}%` }}
                />
              </div>
            </div>
            {weeklyData.weekly_goal_completed && (
              <div className="pt-2 border-t border-border/50">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-success">
                    ✓ Weekly Goal Completed!
                  </span>
                  <div className="text-xs text-muted-foreground">
                    +{weeklyData.reward_xp} XP, +{weeklyData.reward_points} pts
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Streak */}
      {weeklyData && (
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4">
          <h3 className="text-lg font-semibold text-foreground/95 mb-2">Streak</h3>
          <div className="flex items-center gap-4">
            <div>
              <div className="text-2xl font-bold text-warning">
                🔥 {weeklyData.current_streak_weeks}
              </div>
              <div className="text-sm text-muted-foreground">Current Streak</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-foreground">
                {weeklyData.best_streak_weeks}
              </div>
              <div className="text-sm text-muted-foreground">Best Streak</div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

