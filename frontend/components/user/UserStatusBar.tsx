'use client';

import { useEffect, useState } from 'react';
import { getUserXP, getUserWeekly, type XPProgress, type WeeklyGoalProgress } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Link } from 'lucide-react';

/**
 * Unified User Status Bar component for navigation
 * Displays: Avatar + Level, XP progress bar, Points, Weekly streak
 */
export default function UserStatusBar() {
  const { user, isAuthenticated } = useAuth();
  const [xpData, setXpData] = useState<XPProgress | null>(null);
  const [weeklyData, setWeeklyData] = useState<WeeklyGoalProgress | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated || !user) {
      setLoading(false);
      return;
    }

    async function fetchStats() {
      try {
        const [xp, weekly] = await Promise.all([
          getUserXP().catch((err: any) => {
            if (err?.status === 401 || err?.status === 404) return null;
            if (err?.status !== 401 && err?.status !== 404) {
              console.error('Failed to fetch XP:', err);
            }
            return null;
          }),
          getUserWeekly().catch((err: any) => {
            if (err?.status === 401 || err?.status === 404) return null;
            if (err?.status !== 401 && err?.status !== 404) {
              console.error('Failed to fetch weekly stats:', err);
            }
            return null;
          }),
        ]);
        setXpData(xp);
        setWeeklyData(weekly);
      } catch (error: any) {
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

  if (!isAuthenticated || loading || !xpData) {
    return null;
  }

  // Calculate progress percentage
  const progressPercent = xpData.xp_for_next_level > 0
    ? ((xpData.progress_in_current_level / xpData.xp_for_next_level) * 100)
    : 100;

  return (
    <div className="flex items-center gap-2 px-2.5 py-1.5 border border-border/30 rounded-lg bg-muted/20">
      {/* XP Bar with Level integrated - Main visual green accent */}
      <div className="flex flex-col gap-1 min-w-[100px]">
        <div className="flex items-center justify-between gap-2 mb-0.5">
          <span className="text-xs font-medium text-muted-foreground">Lv.{xpData.level}</span>
          <span className="text-xs text-muted-foreground whitespace-nowrap tabular-nums leading-none">
            {xpData.progress_in_current_level}/{xpData.xp_for_next_level}
          </span>
        </div>
        <div className="flex-1 h-2 bg-muted rounded-full overflow-hidden">
          <div
            className="h-full bg-primary transition-all duration-300"
            style={{ width: `${Math.min(progressPercent, 100)}%` }}
          />
        </div>
      </div>

      {/* Points and Streak - White/grey text, 8px gap from XP bar */}
      <div className="flex items-center gap-2 text-sm text-foreground">
        {user && (
          <span className="tabular-nums whitespace-nowrap">
            {user.points || 0} pts
          </span>
        )}
        {weeklyData && weeklyData.current_streak_weeks > 0 && (
          <>
            <span className="text-muted-foreground">·</span>
            <div className="flex items-center gap-1">
              <Link className="w-3.5 h-3.5 text-muted-foreground" />
              <span className="tabular-nums whitespace-nowrap">
                {weeklyData.current_streak_weeks}w
              </span>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

