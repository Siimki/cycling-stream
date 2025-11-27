import { getRace, getRaceStream } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import Link from 'next/link';
import ErrorMessage from '@/components/ErrorMessage';
import dynamic from 'next/dynamic';
import WatchTrackingProvider from '@/components/WatchTrackingProvider';
import DynamicVideoPlayer from '@/components/video/DynamicVideoPlayer';
import { StreamHeader } from '@/components/race/StreamHeader';
import { RaceStats } from '@/components/race/RaceStats';
import { PointsDisplay } from '@/components/user/PointsDisplay';
import { HudStatsProvider } from '@/components/user/HudStatsProvider';
import { AuthRequiredWrapper } from '@/components/race/AuthRequiredWrapper';
import { StreamProvider } from '@/components/race/StreamProvider';
import { ChatWrapper } from '@/components/race/ChatWrapper';
import Footer from '@/components/layout/Footer';
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
    <div className="min-h-screen bg-background flex flex-col">
      <StreamHeader />

      <main className="flex-1 flex flex-col lg:flex-row">
        <WatchTrackingProvider raceId={id}>
          <HudStatsProvider>
            {/* Main content area - Flex Col for Video + Stats + Points */}
            <div className="flex-1 flex flex-col bg-background relative min-w-0">
              
              {/* Video Player Section - Proper height on all devices */}
              <div className="w-full bg-black flex items-center justify-center relative">
                <div className="w-full aspect-video max-w-full">
                  <AuthRequiredWrapper requiresLogin={requiresLogin} raceId={id}>
                    <StreamProvider
                      raceId={id}
                      requiresLogin={requiresLogin}
                      initialStream={stream}
                    />
                  </AuthRequiredWrapper>
                </div>
              </div>

              {/* Race Stats Section - Collapsible */}
              <div className="shrink-0">
                 <RaceStats race={race} />
              </div>

              {/* Points Display Section - Collapsible */}
              <div className="shrink-0">
                 <PointsDisplay />
              </div>

            </div>

            {/* Chat sidebar - Fixed height, scrollable */}
            <ChatWrapper raceId={id} requiresLogin={requiresLogin} isLive={isLive} />
          </HudStatsProvider>
        </WatchTrackingProvider>
      </main>
      <Footer />
    </div>
  );
}
