'use client';

import { useAuth } from '@/contexts/AuthContext';
import { Clock, Zap, Trophy } from 'lucide-react';
import { formatTime } from '@/lib/formatters';

interface PersonalStatsProps {
  totalWatchMinutes?: number;
}

export function PersonalStats({ totalWatchMinutes = 0 }: PersonalStatsProps) {
  const { user } = useAuth();
  const watchTimeSeconds = totalWatchMinutes * 60;

  return (
    <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 mb-8">
      <h2 className="text-xl font-bold text-foreground mb-4">Your Stats</h2>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-primary/20 flex items-center justify-center">
            <Clock className="w-5 h-5 text-primary" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Watch Time</div>
            <div className="text-lg font-bold">{formatTime(watchTimeSeconds)}</div>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-primary/20 flex items-center justify-center">
            <Zap className="w-5 h-5 text-primary" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Points</div>
            <div className="text-lg font-bold">{(user?.points || 0).toLocaleString()}</div>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-primary/20 flex items-center justify-center">
            <Trophy className="w-5 h-5 text-primary" />
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Races Watched</div>
            <div className="text-lg font-bold">-</div>
          </div>
        </div>
      </div>
    </div>
  );
}

