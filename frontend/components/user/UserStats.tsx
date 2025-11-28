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
      <div className="flex items-center gap-2 text-sm text-gray-500">
        <div className="h-4 w-16 animate-pulse bg-gray-200 rounded"></div>
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
        <div className="flex items-center gap-1.5 px-2 py-1 bg-blue-100 dark:bg-blue-900 rounded-md">
          <span className="font-semibold text-blue-700 dark:text-blue-300">Lv.{xpData.level}</span>
        </div>

        {/* XP Bar */}
        <div className="flex items-center gap-2 min-w-[120px]">
          <div className="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500 transition-all duration-300"
              style={{ width: `${Math.min(progressPercent, 100)}%` }}
            />
          </div>
          <span className="text-xs text-gray-600 dark:text-gray-400 whitespace-nowrap">
            {xpData.progress_in_current_level}/{xpData.xp_for_next_level}
          </span>
        </div>

        {/* Points */}
        {user && (
          <div className="flex items-center gap-1 px-2 py-1 bg-green-100 dark:bg-green-900 rounded-md">
            <span className="text-green-700 dark:text-green-300 font-medium">
              {user.points || 0} pts
            </span>
          </div>
        )}

        {/* Streak */}
        {weeklyData && weeklyData.current_streak_weeks > 0 && (
          <div className="flex items-center gap-1 px-2 py-1 bg-orange-100 dark:bg-orange-900 rounded-md">
            <span className="text-orange-600 dark:text-orange-400">🔥</span>
            <span className="text-orange-700 dark:text-orange-300 font-medium">
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
      <div className="bg-white dark:bg-gray-800 rounded-lg p-4 shadow">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-lg font-semibold">Level {xpData.level}</h3>
          <div className="text-sm text-gray-600 dark:text-gray-400">
            {xpData.xp_total} XP
          </div>
        </div>
        <div className="space-y-1">
          <div className="flex justify-between text-sm text-gray-600 dark:text-gray-400">
            <span>Progress to Level {xpData.level + 1}</span>
            <span>{xpData.progress_in_current_level} / {xpData.xp_for_next_level} XP</span>
          </div>
          <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500 transition-all duration-300"
              style={{ width: `${Math.min(progressPercent, 100)}%` }}
            />
          </div>
        </div>
      </div>

      {/* Weekly Goal */}
      {weeklyData && (
        <div className="bg-white dark:bg-gray-800 rounded-lg p-4 shadow">
          <h3 className="text-lg font-semibold mb-3">This Week</h3>
          <div className="space-y-3">
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600 dark:text-gray-400">Watch Time</span>
                <span className="font-medium">{weeklyData.watch_minutes} / 30 min</span>
              </div>
              <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                <div
                  className="h-full bg-blue-500"
                  style={{ width: `${Math.min((weeklyData.watch_minutes / 30) * 100, 100)}%` }}
                />
              </div>
            </div>
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600 dark:text-gray-400">Chat Messages</span>
                <span className="font-medium">{weeklyData.chat_messages} / 3</span>
              </div>
              <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                <div
                  className="h-full bg-green-500"
                  style={{ width: `${Math.min((weeklyData.chat_messages / 3) * 100, 100)}%` }}
                />
              </div>
            </div>
            {weeklyData.weekly_goal_completed && (
              <div className="pt-2 border-t border-gray-200 dark:border-gray-700">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-green-600 dark:text-green-400">
                    ✓ Weekly Goal Completed!
                  </span>
                  <div className="text-xs text-gray-500">
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
        <div className="bg-white dark:bg-gray-800 rounded-lg p-4 shadow">
          <h3 className="text-lg font-semibold mb-2">Streak</h3>
          <div className="flex items-center gap-4">
            <div>
              <div className="text-2xl font-bold text-orange-600 dark:text-orange-400">
                🔥 {weeklyData.current_streak_weeks}
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-400">Current Streak</div>
            </div>
            <div>
              <div className="text-2xl font-bold text-gray-700 dark:text-gray-300">
                {weeklyData.best_streak_weeks}
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-400">Best Streak</div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

