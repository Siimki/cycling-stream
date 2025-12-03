'use client';

import { useEffect, useState } from 'react';
import { getUserXP, type XPProgress } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';

/**
 * Unified User Status Bar component for navigation
 * Displays: Level, XP progress bar, Points
 * 
 * Redesigned to match compact pill-style status indicator with:
 * - Darker background with full rounded corners
 * - Thicker, more prominent progress bar
 * - Green-colored points text for emphasis
 * - Simplified layout (removed progress fraction and streak)
 * 
 * Input data contract unchanged for backward compatibility.
 */
export default function UserStatusBar() {
  const { user, isAuthenticated } = useAuth();
  const [xpData, setXpData] = useState<XPProgress | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated || !user) {
      setLoading(false);
      return;
    }

    async function fetchStats() {
      try {
        const xp = await getUserXP().catch((err: any) => {
          if (err?.status === 401 || err?.status === 404) return null;
          if (err?.status !== 401 && err?.status !== 404) {
            console.error('Failed to fetch XP:', err);
          }
          return null;
        });
        setXpData(xp);
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
    <div className="flex items-center gap-3 px-4 py-2 rounded-full bg-[oklch(0.16_0.005_260)] transition-all duration-200 hover:scale-[0.98] hover:bg-[oklch(0.18_0.005_260)] hover:shadow-lg hover:shadow-primary/20">
      {/* Level (left) */}
      <span className="text-sm font-medium text-muted-foreground whitespace-nowrap">
        Lvl {xpData.level}
      </span>

      {/* Progress Bar (center) - Thicker and more prominent */}
      <div className="min-w-[120px] max-w-[160px] flex-1">
        <div className="h-3 bg-[oklch(0.20_0.005_260)] rounded-full overflow-hidden">
          <div
            className="h-full bg-primary rounded-full transition-all duration-300"
            style={{ width: `${Math.min(progressPercent, 100)}%` }}
          />
        </div>
      </div>

      {/* Points (right) - Green colored for emphasis */}
      {user && (
        <span className="text-sm font-semibold text-primary tabular-nums whitespace-nowrap">
          {user.points || 0} pts
        </span>
      )}
    </div>
  );
}

