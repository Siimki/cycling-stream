'use client';

import { Race } from '@/lib/api';
import RaceCard from '@/components/race/RaceCard';
import { ForYouSection } from './ForYouSection';

interface ContinueWatchingProps {
  races: Race[];
}

export function ContinueWatching({ races }: ContinueWatchingProps) {
  const racesList = races || [];
  return (
    <ForYouSection
      title="Continue Watching"
      emptyMessage="No races to continue watching. Start watching a race to see it here!"
      isEmpty={racesList.length === 0}
    >
      {racesList.map((race) => (
        <RaceCard key={race.id} race={race} />
      ))}
    </ForYouSection>
  );
}

