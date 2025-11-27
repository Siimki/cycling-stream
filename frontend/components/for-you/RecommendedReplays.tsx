'use client';

import { Race } from '@/lib/api';
import RaceCard from '@/components/race/RaceCard';
import { ForYouSection } from './ForYouSection';

interface RecommendedReplaysProps {
  races: Race[];
}

export function RecommendedReplays({ races }: RecommendedReplaysProps) {
  const racesList = races || [];
  return (
    <ForYouSection
      title="Recommended Replays"
      emptyMessage="Watch some races to get personalized recommendations!"
      isEmpty={racesList.length === 0}
    >
      {racesList.map((race) => (
        <RaceCard key={race.id} race={race} />
      ))}
    </ForYouSection>
  );
}

