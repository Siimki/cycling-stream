'use client';

import Link from 'next/link';
import { Race, StreamResponse } from '@/lib/api';
import { LiveRaceCard } from './LiveRaceCard';
import { ArrowRight } from 'lucide-react';

interface LiveRacesSectionProps {
  races: Array<{
    race: Race;
    stream?: StreamResponse;
    viewerCount?: number;
    progressPercent?: number;
    timeRemaining?: string;
  }>;
}

export function LiveRacesSection({ races }: LiveRacesSectionProps) {
  const liveRaces = races.filter(item => item.stream?.status === 'live');
  const upcomingRaces = races.filter(item => item.stream?.status === 'planned');
  const displayRaces = [...liveRaces, ...upcomingRaces].slice(0, 3);

  if (displayRaces.length === 0) {
    return null;
  }

  return (
    <section className="py-12 mb-12">
      <div className="flex items-center justify-between mb-10">
        <div>
          <h2 className="text-3xl sm:text-4xl font-bold mb-2">Happening Now</h2>
          <p className="text-muted-foreground">Jump into live action from around the world</p>
        </div>
        <Link 
          href="/races" 
          className="text-primary font-medium flex items-center gap-2 hover:gap-3 transition-all"
        >
          View All Races
          <ArrowRight className="w-5 h-5" />
        </Link>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {displayRaces.map((item) => (
          <LiveRaceCard
            key={item.race.id}
            race={item.race}
            status={item.stream?.status as 'live' | 'planned' | 'offline'}
            viewerCount={item.viewerCount}
            progressPercent={item.progressPercent}
            timeRemaining={item.timeRemaining}
          />
        ))}
      </div>
    </section>
  );
}

