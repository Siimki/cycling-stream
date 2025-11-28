'use client';

import { Race } from '@/lib/api';
import RaceCard from '@/components/race/RaceCard';

interface UpcomingFromFavoritesProps {
  races: Race[];
}

export function UpcomingFromFavorites({ races }: UpcomingFromFavoritesProps) {
  const racesList = races || [];

  if (racesList.length === 0) {
    return null;
  }

  return (
    <div className="mb-12">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold">Upcoming from Your Favorites</h2>
          <p className="text-muted-foreground text-sm mt-1">Based on your followed series and teams</p>
        </div>
        <a 
          href="/favorites" 
          className="text-muted-foreground hover:text-foreground text-sm font-medium transition"
        >
          Manage Favorites
        </a>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {racesList.slice(0, 3).map((race) => (
          <RaceCard key={race.id} race={race} />
        ))}
      </div>
    </div>
  );
}

