'use client';

import Link from 'next/link';
import { Race } from '@/lib/api';
import RaceCard from '@/components/race/RaceCard';

interface ContinueWatchingSectionProps {
  races: Race[];
}

export function ContinueWatchingSection({ races }: ContinueWatchingSectionProps) {
  const racesList = races || [];

  if (racesList.length === 0) {
    return null;
  }

  return (
    <div className="mb-12">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Continue Watching</h2>
        <Link 
          href="/for-you" 
          className="text-muted-foreground hover:text-foreground text-sm font-medium transition"
        >
          View All
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {racesList.slice(0, 3).map((race) => (
          <RaceCard key={race.id} race={race} />
        ))}
      </div>
    </div>
  );
}

