'use client';

import { useEffect, useState } from 'react';
import { getUserWeekly, type WeeklyGoalProgress } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Link } from 'lucide-react';

export function StreakCard() {
  const { isAuthenticated, user } = useAuth();
  const [weeklyData, setWeeklyData] = useState<WeeklyGoalProgress | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated || !user) {
      setLoading(false);
      return;
    }

    async function fetchStats() {
      try {
        const weekly = await getUserWeekly().catch((err: any) => {
          if (err?.status === 401 || err?.status === 404) return null;
          if (err?.status !== 401 && err?.status !== 404) {
            console.error('Failed to fetch weekly stats:', err);
          }
          return null;
        });
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

  if (!weeklyData) {
    return null;
  }

  const isActive = weeklyData.current_streak_weeks > 0;

  return (
    <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
      <h3 className="text-lg font-semibold text-foreground/95 mb-4">Streak</h3>
      <div className="flex items-start gap-8">
        <div className="flex-1">
          <div className="flex items-baseline gap-2 mb-1">
            <Link className={`w-6 h-6 ${isActive ? 'text-primary' : 'text-muted-foreground'} flex-shrink-0`} />
            <span className={`text-3xl font-bold leading-none ${isActive ? 'text-primary' : 'text-foreground'}`}>
              {weeklyData.current_streak_weeks}
            </span>
          </div>
          <div className="text-xs text-muted-foreground ml-8">Current Streak</div>
        </div>
        <div className="flex-1">
          <div className="text-2xl font-bold text-foreground leading-none mb-1">
            {weeklyData.best_streak_weeks}
          </div>
          <div className="text-xs text-muted-foreground">Best Streak</div>
        </div>
      </div>
    </div>
  );
}

