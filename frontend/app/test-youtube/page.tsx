import DynamicVideoPlayer from '@/components/video/DynamicVideoPlayer';

export default function TestYoutubePage() {
  const streamType = 'youtube';
  const sourceId = 'jfKfPfyJRdk'; // ID from user provided iframe
  const status = 'live';

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-8">
      <h1 className="text-2xl font-bold mb-4">YouTube Stream Test</h1>
      <div className="w-full max-w-4xl aspect-video bg-black rounded-lg overflow-hidden border border-border">
        <DynamicVideoPlayer
          status={status}
          streamType={streamType}
          sourceId={sourceId}
        />
      </div>
      <p className="mt-4 text-muted-foreground">
        Testing YouTube Embed: {sourceId}
      </p>
    </div>
  );
}

