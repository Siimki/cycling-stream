'use client';

import Link from 'next/link';
import { Race } from '@/lib/api';
import { formatDate } from '@/lib/formatters';

interface LiveRaceCardProps {
  race: Race;
  status?: 'live' | 'planned' | 'offline';
  viewerCount?: number;
  progressPercent?: number;
  timeRemaining?: string;
}

export function LiveRaceCard({ race, status = 'offline', viewerCount, progressPercent, timeRemaining }: LiveRaceCardProps) {
  const isLive = status === 'live';
  const isUpcoming = status === 'planned';

  return (
    <Link href={`/races/${race.id}`}>
      <div className="group relative h-64 bg-card/80 backdrop-blur-sm rounded-xl overflow-hidden cursor-pointer border border-border/50 hover:border-border transition-all hover:-translate-y-1 hover:shadow-lg hover:shadow-primary/20">
        {/* Placeholder image - can be replaced with race image later */}
        <div className="absolute inset-0 bg-gradient-to-br from-primary/20 to-primary/5 opacity-60 group-hover:opacity-80 transition-opacity" />
        
        {/* Status badges */}
        <div className="absolute top-3 left-3 flex gap-2 z-10">
          {isLive && (
            <span className="bg-red-600 text-white text-[10px] font-bold px-2 py-1 rounded uppercase tracking-wider">
              Live
            </span>
          )}
          {isUpcoming && (
            <span className="bg-muted text-foreground text-[10px] font-bold px-2 py-1 rounded uppercase tracking-wider">
              {timeRemaining || 'Upcoming'}
            </span>
          )}
          {viewerCount !== undefined && viewerCount > 0 && (
            <span className="bg-black/60 backdrop-blur-sm text-white text-[10px] font-bold px-2 py-1 rounded">
              {viewerCount >= 1000 ? `${(viewerCount / 1000).toFixed(1)}K` : viewerCount} viewers
            </span>
          )}
        </div>

        {/* Content overlay */}
        <div className="absolute bottom-0 left-0 w-full p-5 bg-gradient-to-t from-black via-black/90 to-transparent">
          <h3 className="text-white text-lg font-bold mb-2">{race.name}</h3>
          <div className="flex items-center gap-4 text-xs text-muted-foreground mb-2">
            {race.stage_name && <span>{race.stage_name}</span>}
            {race.stage_name && <span>•</span>}
            {timeRemaining && (
              <span className={isLive ? 'text-primary font-semibold' : ''}>
                {timeRemaining}
              </span>
            )}
          </div>
          {progressPercent !== undefined && progressPercent > 0 && (
            <div className="flex items-center gap-2">
              <div className="w-full bg-muted/50 rounded-full h-1">
                <div 
                  className="bg-primary h-1 rounded-full transition-all"
                  style={{ width: `${Math.min(progressPercent, 100)}%` }}
                />
              </div>
              <span className="text-[10px] text-muted-foreground font-mono">
                {Math.round(progressPercent)}%
              </span>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
}

