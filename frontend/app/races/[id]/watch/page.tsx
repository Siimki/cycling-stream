import { getRace, getRaceStream } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import Link from 'next/link';
import ErrorMessage from '@/components/ErrorMessage';
import { StreamHeader } from '@/components/race/StreamHeader';
import { WatchExperienceLayout } from '@/components/race/WatchExperienceLayout';
import { notFound } from 'next/navigation';

interface WatchPageProps {
  params: Promise<{ id: string }>;
}

export default async function WatchPage({ params }: WatchPageProps) {
  const { id } = await params;
  let race = null;
  let stream = null;
  let error = null;

  try {
    race = await getRace(id);
    try {
      stream = await getRaceStream(id);
    } catch (streamErr: any) {
      // If stream fetch fails with 401 and race requires login, that's expected
      // Set stream to offline so the AuthRequiredWrapper can handle it
      // Don't treat this as an error - it's expected behavior
      if (streamErr?.status === 401 && race?.requires_login) {
        stream = { status: 'offline' };
      } else {
        // For other errors, still set to offline but don't log expected auth errors
        stream = { status: 'offline' };
      }
    }
  } catch (err) {
    error = APIErrorHandler.getErrorMessage(err);
  }

  if (error && error.includes('not found')) {
    notFound();
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background p-8">
        <ErrorMessage message={error} />
        <Link href="/" className="mt-4 inline-block text-primary hover:underline">
          ← Back to races
        </Link>
      </div>
    );
  }

  if (!race) {
    notFound();
  }

  const isLive = stream?.status === 'live';
  const requiresLogin = race.requires_login || false;

  return (
    <div className="min-h-screen flex flex-col bg-black">
      <div className="flex-none">
        <StreamHeader />
      </div>

      <main className="flex-1 flex min-h-0">
        <WatchExperienceLayout
          raceId={id}
          race={race}
          stream={stream}
          requiresLogin={requiresLogin}
          isLive={isLive}
        />
      </main>
    </div>
  );
}
