/**
 * Video overlay component showing live badge and other overlay elements
 */

import { memo } from 'react';

interface VideoOverlayProps {
  showControls: boolean;
}

export const VideoOverlay = memo(function VideoOverlay({ showControls }: VideoOverlayProps) {
  return (
    <>
      {/* Top gradient */}
      <div
        className={`absolute inset-x-0 top-0 h-24 bg-gradient-to-b from-black/60 to-transparent pointer-events-none transition-opacity duration-200 ${
          showControls ? 'opacity-100' : 'opacity-0'
        }`}
      />

      {/* Live Badge */}
      <div
        className={`absolute top-3 left-3 flex items-center gap-2 transition-opacity duration-200 ${
          showControls ? 'opacity-100' : 'opacity-0'
        }`}
      >
        <div className="flex items-center gap-1.5 bg-red-600 pl-1.5 pr-2 py-0.5 rounded text-white">
          <span className="w-1.5 h-1.5 rounded-full bg-white animate-live-pulse" />
          <span className="text-xs font-bold uppercase tracking-wide">Live</span>
        </div>
      </div>
    </>
  );
});

