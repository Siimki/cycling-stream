'use client';

import { Race } from '@/lib/api';
import RaceCard from '@/components/race/RaceCard';
import { ForYouSection } from './ForYouSection';

interface UpcomingRacesProps {
  races: Race[];
}

export function UpcomingRaces({ races }: UpcomingRacesProps) {
  const racesList = races || [];
  return (
    <ForYouSection
      title="Upcoming Races"
      emptyMessage="No upcoming races matching your preferences. Check back soon!"
      isEmpty={racesList.length === 0}
    >
      {racesList.map((race) => (
        <RaceCard key={race.id} race={race} />
      ))}
    </ForYouSection>
  );
}

