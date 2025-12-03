'use client';

import { useEffect, useState } from 'react';
import { getUserXP, getUserWeekly, claimWeeklyReward, type XPProgress, type WeeklyGoalProgress } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';

export function WeeklyOverview() {
  const { isAuthenticated, user, refreshUser } = useAuth();
  const [xpData, setXpData] = useState<XPProgress | null>(null);
  const [weeklyData, setWeeklyData] = useState<WeeklyGoalProgress | null>(null);
  const [loading, setLoading] = useState(true);
  const [claiming, setClaiming] = useState(false);

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
    const interval = setInterval(fetchStats, 30000);
    return () => clearInterval(interval);
  }, [isAuthenticated, user]);

  if (loading) {
    return (
      <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <div className="h-4 w-16 animate-pulse bg-muted rounded"></div>
        </div>
      </div>
    );
  }

  if (!xpData) {
    return (
      <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6 text-center">
        <p className="text-sm text-muted-foreground">Weekly stats will appear here once you start watching.</p>
      </div>
    );
  }

  const progressPercent = xpData.xp_for_next_level > 0
    ? ((xpData.progress_in_current_level / xpData.xp_for_next_level) * 100)
    : 100;

  const handleClaimWeekly = async () => {
    if (claiming || !weeklyData?.can_claim_reward) return;
    
    setClaiming(true);
    try {
      await claimWeeklyReward();
      // Refresh weekly data and user
      const [xp, weekly] = await Promise.all([
        getUserXP().catch(() => null),
        getUserWeekly().catch(() => null),
      ]);
      setXpData(xp);
      setWeeklyData(weekly);
      if (refreshUser) {
        await refreshUser();
      }
    } catch (error) {
      console.error('Failed to claim weekly reward:', error);
    } finally {
      setClaiming(false);
    }
  };

  return (
    <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
      {/* Level & XP Progress - Main progress bar */}
      <div className="mb-4">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-lg font-semibold text-foreground/95">
            Level {xpData.level}
          </h3>
          <div className="text-sm text-muted-foreground text-right min-w-[80px]">
            {xpData.xp_total} XP
          </div>
        </div>
        <div className="space-y-1">
          <div className="flex justify-between text-sm text-muted-foreground">
            <span>Progress to Level {xpData.level + 1}</span>
            <span className="text-right min-w-[80px]">{xpData.progress_in_current_level} / {xpData.xp_for_next_level} XP</span>
          </div>
          <div className="h-2 bg-muted rounded-full overflow-hidden">
            <div
              className="h-full bg-primary"
              style={{ width: `${Math.min(progressPercent, 100)}%` }}
            />
          </div>
        </div>
      </div>

      {/* Weekly Goals - Dynamic thresholds and rewards */}
      {weeklyData && (
        <div className="space-y-4 pt-4 border-t border-border/50">
          {/* Watch Time Goal */}
          <div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Watch Time</span>
              <span className="font-medium text-foreground text-right min-w-[80px]">
                {weeklyData.watch_minutes} / {weeklyData.watch_minutes_goal || 30} min
              </span>
            </div>
          </div>

          {/* Chat Messages Goal */}
          <div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Chat Messages</span>
              <span className="font-medium text-foreground text-right min-w-[80px]">
                {weeklyData.chat_messages} / {weeklyData.chat_messages_goal || 3}
              </span>
            </div>
          </div>

          {/* Weekly Reward Info */}
          {weeklyData.reward_xp > 0 && weeklyData.reward_points > 0 && (
            <div className="pt-2">
              {weeklyData.can_claim_reward ? (
                <button
                  onClick={handleClaimWeekly}
                  disabled={claiming}
                  className="w-full bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold py-2 px-4 rounded-lg transition disabled:opacity-50"
                >
                  {claiming ? 'Claiming...' : `Claim weekly reward (+${weeklyData.reward_xp} XP, +${weeklyData.reward_points} pts)`}
                </button>
              ) : weeklyData.weekly_goal_completed ? (
                <div className="text-xs text-muted-foreground text-center">
                  Weekly reward claimed! Come back next week.
                </div>
              ) : (
                <div className="text-xs text-muted-foreground text-center">
                  Weekly reward: +{weeklyData.reward_xp} XP, +{weeklyData.reward_points} pts
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

