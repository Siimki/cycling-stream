'use client';

import dynamic from 'next/dynamic';

const VideoPlayer = dynamic(() => import('./VideoPlayer'), {
  ssr: false,
  loading: () => <div className="aspect-video bg-black" />,
});

interface DynamicVideoPlayerProps {
  streamUrl?: string;
  status: string;
  streamType?: string;
  sourceId?: string;
}

export default function DynamicVideoPlayer({ streamUrl, status, streamType, sourceId }: DynamicVideoPlayerProps) {
  return <VideoPlayer streamUrl={streamUrl} status={status} streamType={streamType} sourceId={sourceId} />;
}

