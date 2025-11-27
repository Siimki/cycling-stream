import { getRace } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import Link from 'next/link';
import ErrorMessage from '@/components/ErrorMessage';
import { Navigation } from '@/components/layout/Navigation';
import { Button } from '@/components/ui/button';
import Footer from '@/components/layout/Footer';
import { formatDate } from '@/lib/formatters';
import { notFound } from 'next/navigation';

interface RaceDetailPageProps {
  params: Promise<{ id: string }>;
}

export default async function RaceDetailPage({ params }: RaceDetailPageProps) {
  const { id } = await params;
  let race = null;
  let error = null;

  try {
    race = await getRace(id);
  } catch (err) {
    error = APIErrorHandler.getErrorMessage(err);
  }

  if (error && error.includes('not found')) {
    notFound();
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex flex-col">
        <Navigation variant="full" />
        <div className="flex-1 p-8">
          <div className="max-w-4xl mx-auto">
            <ErrorMessage message={error} />
            <Link href="/" className="mt-4 inline-block text-primary hover:underline">
              ← Back to races
            </Link>
          </div>
        </div>
        <Footer />
      </div>
    );
  }

  if (!race) {
    notFound();
  }


  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6 lg:p-8">
          <h1 className="text-2xl sm:text-3xl font-bold text-foreground/95 mb-4">{race.name}</h1>

          {race.description && (
            <p className="text-muted-foreground mb-6 text-sm sm:text-base">{race.description}</p>
          )}

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
            {race.location && (
              <div>
                <h3 className="text-xs sm:text-sm font-semibold text-muted-foreground uppercase tracking-wider">Location</h3>
                <p className="text-foreground/95 text-sm sm:text-base mt-1">{race.location}</p>
              </div>
            )}
            {race.start_date && (
              <div>
                <h3 className="text-xs sm:text-sm font-semibold text-muted-foreground uppercase tracking-wider">Start Date</h3>
                <p className="text-foreground/95 text-sm sm:text-base mt-1">{formatDate(race.start_date, { includeTime: true })}</p>
              </div>
            )}
            {race.category && (
              <div>
                <h3 className="text-xs sm:text-sm font-semibold text-muted-foreground uppercase tracking-wider">Category</h3>
                <p className="text-foreground/95 text-sm sm:text-base mt-1">{race.category}</p>
              </div>
            )}
            <div>
              <h3 className="text-xs sm:text-sm font-semibold text-muted-foreground uppercase tracking-wider">Price</h3>
              <p className="text-foreground/95 text-sm sm:text-base mt-1">
                {race.is_free ? (
                  <span className="text-primary font-semibold">Free</span>
                ) : (
                  <span>${(race.price_cents / 100).toFixed(2)}</span>
                )}
              </p>
            </div>
          </div>

          <Link href={`/races/${race.id}/watch`}>
            <Button className="w-full sm:w-auto bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold">
              Watch Race
            </Button>
          </Link>
        </div>
      </main>
      <Footer />
    </div>
  );
}

