import Link from 'next/link';

import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';
import ErrorMessage from '@/components/ErrorMessage';
import RaceCard from '@/components/race/RaceCard';
import { Button } from '@/components/ui/button';
import { getRaces, Race } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import { isRaceReplay } from '@/lib/raceFilters';

export const metadata = {
  title: 'Replays | PelotonLive',
  description: 'Watch past races and criteriums on-demand.',
};

export default async function ReplaysPage() {
  let races: Race[] = [];
  let error: string | null = null;

  try {
    races = await getRaces();
  } catch (err) {
    error = APIErrorHandler.getErrorMessage(err);
    races = [];
  }

  const replayRaces = races.filter((race) => isRaceReplay(race));

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-7xl mx-auto px-6 lg:px-8 pt-8 pb-6 sm:pb-8 w-full">
        <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl sm:text-4xl font-bold text-foreground mb-2">Replays</h1>
            <p className="text-muted-foreground text-base sm:text-lg">
              Watch recently finished races on-demand. The Fuji Criterium doubleheader now lives here along with other past events.
            </p>
          </div>
          <Link href="/races">
            <Button variant="outline">View upcoming &amp; live races</Button>
          </Link>
        </div>

        {error ? (
          <ErrorMessage message={error} />
        ) : replayRaces.length === 0 ? (
          <div className="text-center py-12 px-4">
            <p className="text-muted-foreground text-base sm:text-lg">No replays are available yet.</p>
            <p className="text-muted-foreground/70 mt-2 text-sm sm:text-base">
              Check the live schedule while we add more finished races.
            </p>
            <div className="mt-6 flex justify-center">
              <Link href="/races">
                <Button>Browse live races</Button>
              </Link>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
            {replayRaces.map((race) => (
              <RaceCard key={race.id} race={race} />
            ))}
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}
