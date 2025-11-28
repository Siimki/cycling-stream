import { getRaces, Race } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import RaceCard from '@/components/race/RaceCard';
import ErrorMessage from '@/components/ErrorMessage';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';

export const metadata = {
  title: 'Races | PelotonLive',
  description: 'Browse all cycling races available on PelotonLive',
};

export default async function RacesPage() {
  let races: Race[] = [];
  let error: string | null = null;

  try {
    const result = await getRaces();
    races = result || [];
  } catch (err) {
    error = APIErrorHandler.getErrorMessage(err);
    races = [];
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-7xl mx-auto px-6 lg:px-8 pt-8 pb-6 sm:pb-8 w-full">
        <div className="mb-8">
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground mb-2">All Races</h1>
          <p className="text-muted-foreground text-base sm:text-lg">
            Browse all cycling races available on PelotonLive
          </p>
        </div>

        {error ? (
          <ErrorMessage message={error} />
        ) : !races || races.length === 0 ? (
          <div className="text-center py-12 px-4">
            <p className="text-muted-foreground text-base sm:text-lg">No races available at the moment.</p>
            <p className="text-muted-foreground/70 mt-2 text-sm sm:text-base">Check back soon for upcoming events!</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
            {races.map((race) => (
              <RaceCard key={race.id} race={race} />
            ))}
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}

