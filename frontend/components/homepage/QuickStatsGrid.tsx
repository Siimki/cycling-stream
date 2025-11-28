'use client';

import { StatCard } from './StatCard';
import { formatTime } from '@/lib/formatters';

interface QuickStatsGridProps {
  watchTimeHours?: number;
  watchTimeMinutes?: number;
  watchTimeChange?: string;
  streakDays?: number;
  longestStreak?: number;
  racesWatched?: number;
  racesWatchedLabel?: string;
  impactAmount?: number;
}

export function QuickStatsGrid({
  watchTimeHours = 0,
  watchTimeMinutes = 0,
  watchTimeChange,
  streakDays = 0,
  longestStreak,
  racesWatched = 0,
  racesWatchedLabel,
  impactAmount = 0,
}: QuickStatsGridProps) {
  // Format watch time
  const totalMinutes = watchTimeHours * 60 + watchTimeMinutes;
  const totalHours = watchTimeHours + (watchTimeMinutes / 60);
  const watchTimeDisplay = totalHours >= 1 
    ? `${totalHours.toFixed(1)}h`
    : `${Math.round(totalMinutes)}m`;

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-12">
      {/* Watch Time */}
      <StatCard
        label="Watch Time"
        value={watchTimeDisplay}
        icon={<span className="text-2xl">⏱️</span>}
        change={watchTimeChange ? `+${watchTimeChange} this week` : undefined}
        changeType={watchTimeChange ? 'positive' : 'neutral'}
      />

      {/* Current Streak */}
      <StatCard
        label="Streak"
        value={streakDays > 0 ? `${streakDays} Days` : '0 Days'}
        icon={<span className="text-2xl">🔥</span>}
        change={longestStreak ? `Longest: ${longestStreak} days` : undefined}
      />

      {/* Races Watched */}
      <StatCard
        label="Races Watched"
        value={racesWatched}
        icon={<span className="text-2xl">🏁</span>}
        change={racesWatchedLabel}
      />

      {/* Impact */}
      <StatCard
        label="Your Impact"
        value={`$${impactAmount.toFixed(2)}`}
        icon={<span className="text-2xl">💚</span>}
        change="To race organizers"
        changeType="positive"
        highlight={true}
      />
    </div>
  );
}

