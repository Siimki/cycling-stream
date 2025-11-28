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

        <main className="flex-1">
          <WatchTrackingProvider raceId={id}>
            <HudStatsProvider>
              {/* Consistent spacing between nav and content (32px) */}
              <div className="pt-8">
                {/* 12-column responsive grid layout with gutters - same padding as nav */}
                <div className="grid grid-cols-12 gap-4 lg:gap-6 px-6 lg:px-8">
                {/* Main content area - 12 columns on mobile, 8 on desktop, 9 on large screens */}
                <div className="col-span-12 lg:col-span-8 xl:col-span-9 flex flex-col bg-background relative min-w-0">
                  
                  {/* Video Player Section - Proper height on all devices with max-width constraint */}
                  <div className="w-full bg-background flex items-center justify-center relative py-4 lg:py-6">
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

                {/* Chat sidebar - 12 columns on mobile, 4 on desktop, 3 on large screens */}
                <div className="col-span-12 lg:col-span-4 xl:col-span-3">
                  <ChatWrapper raceId={id} requiresLogin={requiresLogin} isLive={isLive} />
                </div>
                </div>
              </div>
            </HudStatsProvider>
          </WatchTrackingProvider>
        </main>
        <Footer />
      </div>
  );
}
