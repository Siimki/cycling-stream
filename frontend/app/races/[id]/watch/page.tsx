import { getRace, getRaceStream } from '@/lib/api';
import { APIErrorHandler } from '@/lib/error-handler';
import Link from 'next/link';
import ErrorMessage from '@/components/ErrorMessage';
import dynamic from 'next/dynamic';
import WatchTrackingProvider from '@/components/WatchTrackingProvider';
import { ChatProvider } from '@/components/chat/ChatProvider';
import DynamicVideoPlayer from '@/components/video/DynamicVideoPlayer';
import { StreamHeader } from '@/components/race/StreamHeader';
import { RaceStats } from '@/components/race/RaceStats';
import { PointsDisplay } from '@/components/user/PointsDisplay';
import { HudStatsProvider } from '@/components/user/HudStatsProvider';
import Footer from '@/components/layout/Footer';
import { notFound } from 'next/navigation';

const Chat = dynamic(() => import('@/components/chat/Chat'), {
  loading: () => (
    <div className="flex flex-col h-full bg-card/50 border-l border-border min-h-0">
      {/* Header skeleton */}
      <div className="px-4 py-3 border-b border-border/50 flex items-center justify-between shrink-0 h-12">
        <div className="h-4 w-24 bg-muted animate-pulse rounded" />
        <div className="h-8 w-8 bg-muted animate-pulse rounded" />
      </div>
      {/* Messages skeleton */}
      <div className="flex-1 overflow-hidden px-4 py-3 space-y-2">
        <div className="h-4 w-3/4 bg-muted animate-pulse rounded" />
        <div className="h-4 w-full bg-muted animate-pulse rounded" />
        <div className="h-4 w-5/6 bg-muted animate-pulse rounded" />
        <div className="h-4 w-2/3 bg-muted animate-pulse rounded" />
      </div>
      {/* Input skeleton */}
      <div className="p-3 border-t border-border/50 shrink-0 bg-card/30">
        <div className="h-10 w-full bg-muted animate-pulse rounded" />
      </div>
    </div>
  ),
});

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
    } catch {
      stream = { status: 'offline' };
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

  const streamUrl = stream?.cdn_url || stream?.origin_url;
  const isLive = stream?.status === 'live';

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <StreamHeader />

      <main className="flex-1 flex flex-col lg:flex-row">
        <WatchTrackingProvider raceId={id}>
          <HudStatsProvider>
            <ChatProvider raceId={id} enabled={isLive}>
            {/* Main content area - Flex Col for Video + Stats + Points */}
            <div className="flex-1 flex flex-col bg-background relative min-w-0">
              
              {/* Video Player Section - Proper height on all devices */}
              <div className="w-full bg-black flex items-center justify-center relative">
                <div className="w-full aspect-video max-w-full">
                  <DynamicVideoPlayer 
                    streamUrl={streamUrl} 
                    status={stream?.status || 'offline'} 
                    streamType={stream?.stream_type}
                    sourceId={stream?.source_id}
                  />
                </div>
              </div>

              {/* Race Stats Section - Collapsible */}
              <div className="shrink-0">
                 <RaceStats />
              </div>

              {/* Points Display Section - Collapsible */}
              <div className="shrink-0">
                 <PointsDisplay />
              </div>

            </div>

            {/* Chat sidebar - Fixed height, scrollable */}
            <div className="lg:w-80 xl:w-96 2xl:w-[400px] border-t lg:border-t-0 lg:border-l border-border flex flex-col h-[300px] sm:h-[350px] lg:h-[calc(100vh-4rem)] shrink-0 bg-background">
               <Chat />
            </div>
            </ChatProvider>
          </HudStatsProvider>
        </WatchTrackingProvider>
      </main>
      <Footer />
    </div>
  );
}
