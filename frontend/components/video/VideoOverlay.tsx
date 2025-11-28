/**
 * Video overlay component showing live badge and other overlay elements
 */

import { memo } from 'react';
import { cn } from '@/lib/utils';
import { useMotionPreset, useSweepHighlight } from '@/motion';

interface VideoOverlayProps {
  showControls: boolean;
}

export const VideoOverlay = memo(function VideoOverlay({ showControls }: VideoOverlayProps) {
  const overlayFade = useMotionPreset('overlay-fade', { disabled: !showControls });
  const liveBadgeMotion = useMotionPreset('stat-emphasis', { disabled: !showControls });
  const livePulse = useMotionPreset('vip-ring', { disabled: !showControls });
  const sweepHighlight = useSweepHighlight({ disabled: !showControls });

  return (
    <>
      {/* Top gradient + sweep line */}
      <div className="absolute inset-x-0 top-0 h-24 pointer-events-none">
        <div
          className={cn(
            'absolute inset-x-0 top-0 h-full bg-gradient-to-b from-black/60 via-black/40 to-transparent transition-opacity duration-200',
            showControls ? 'opacity-100' : 'opacity-0',
            overlayFade
          )}
        />
        <div
          className={cn(
            'absolute inset-x-0 top-0 h-full opacity-0',
            showControls ? 'opacity-60' : 'opacity-0',
            sweepHighlight
          )}
        />
      </div>

      {/* Live Badge */}
      <div
        className={cn(
          'absolute top-3 left-3 flex items-center gap-2 transition-opacity duration-200',
          showControls ? 'opacity-100' : 'opacity-0',
          liveBadgeMotion
        )}
      >
        <div className="flex items-center gap-1.5 bg-live/90 pl-1.5 pr-2 py-0.5 rounded text-white shadow-lg shadow-red-500/25">
          <span className={cn('w-1.5 h-1.5 rounded-full bg-white', livePulse)} />
          <span className="text-xs font-bold uppercase tracking-wide">Live</span>
        </div>
      </div>
    </>
  );
});

