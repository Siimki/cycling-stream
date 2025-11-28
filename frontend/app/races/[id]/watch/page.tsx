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
    <div className="h-screen flex flex-col bg-black overflow-hidden">
        {/* Header - Fixed height */}
        <div className="flex-none">
          <StreamHeader />
        </div>

        <main className="flex-1 flex min-h-0">
          <WatchTrackingProvider raceId={id}>
            <HudStatsProvider>
              {/* Left Column: Video & Stats - Scrollable */}
              <div className="flex-1 overflow-y-auto scrollbar-hide bg-black relative flex flex-col">
                <div className="flex-1 flex flex-col min-h-[50vh]">
                  {/* Video Player Section - Centered */}
                  <div className="w-full flex-1 flex items-center justify-center bg-black">
                    <div className="w-full max-w-[1600px] aspect-video mx-auto">
                      <AuthRequiredWrapper requiresLogin={requiresLogin} raceId={id}>
                        <StreamProvider
                          raceId={id}
                          requiresLogin={requiresLogin}
                          initialStream={stream}
                        />
                      </AuthRequiredWrapper>
                    </div>
                  </div>

                  {/* Race Stats & Info - Below video */}
                  <div className="w-full bg-background border-t border-border/20">
                    <div className="max-w-[1600px] mx-auto">
                      <RaceStats race={race} />
                      <PointsDisplay />
                    </div>
                  </div>
                  
                  {/* Footer inside scrollable area */}
                  <Footer />
                </div>
              </div>

              {/* Right Column: Chat - Fixed width on desktop */}
              <div className="hidden lg:flex flex-none w-80 xl:w-96 flex-col border-l border-border/20 bg-background h-full">
                <ChatWrapper 
                  raceId={id} 
                  requiresLogin={requiresLogin} 
                  isLive={isLive} 
                  className="h-full border-0"
                />
              </div>
            </HudStatsProvider>

            {/* Mobile Chat Placeholder - visible only on mobile/tablet, potentially togglable or stacked */}
            <div className="lg:hidden flex-none h-[50vh] border-t border-border/20">
                <HudStatsProvider>
                    <ChatWrapper 
                        raceId={id} 
                        requiresLogin={requiresLogin} 
                        isLive={isLive} 
                        className="h-full border-0"
                    />
                </HudStatsProvider>
            </div>
          </WatchTrackingProvider>
        </main>
      </div>
  );
}
