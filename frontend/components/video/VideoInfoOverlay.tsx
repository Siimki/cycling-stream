/**
 * Video Info Overlay Component
 * Displays key race information at the bottom of the video player
 * with improved text hierarchy (labels vs values)
 */

import { memo } from 'react';
import { Mountain, MapPin } from 'lucide-react';
import { Race } from '@/lib/api';
import { cn } from '@/lib/utils';
import { usePopAccent, useSlideIn } from '@/motion';

interface VideoInfoOverlayProps {
  race?: Race;
  showControls: boolean;
}

// Format number with commas
function formatNumber(num: number | undefined | null): string {
  if (num === undefined || num === null) return "N/A";
  return num.toLocaleString();
}

export const VideoInfoOverlay = memo(function VideoInfoOverlay({
  race,
  showControls,
}: VideoInfoOverlayProps) {
  const overlayMotion = useSlideIn('up', { disabled: !showControls });
  const statMotion = usePopAccent({ disabled: !showControls });

  // Only show when controls are visible to avoid cluttering
  if (!showControls || !race) {
    return null;
  }

  const stageName = race?.stage_name || "";
  const stageType = race?.stage_type || "";
  const stageDisplay = stageName && stageType 
    ? `${stageName} — ${stageType}`
    : stageName || stageType || "";
  const elevation = formatNumber(race?.elevation_meters);
  const stageLength = formatNumber(race?.stage_length_km);

  // Don't show if no data
  if (!stageDisplay && !elevation && !stageLength) {
    return null;
  }

  return (
    <div
      className={cn(
        'absolute inset-x-0 bottom-16 bg-gradient-to-t from-black/60 via-black/40 to-transparent transition-opacity duration-200',
        showControls ? 'opacity-100' : 'opacity-0',
        overlayMotion
      )}
    >
      <div className="px-4 py-3">
        <div className="flex items-center gap-4 flex-wrap">
          {/* Stage Name */}
          {stageDisplay && (
            <div className={cn('flex items-center gap-2', statMotion)}>
              <div>
                <div className="text-xs-label text-white/70 mb-0.5">STAGE</div>
                <div className="text-base-value text-white">{stageDisplay}</div>
              </div>
            </div>
          )}

          {/* Elevation */}
          {race?.elevation_meters && (
            <div className={cn('flex items-center gap-2', statMotion)}>
              <Mountain className="w-4 h-4 text-primary/70" />
              <div>
                <div className="text-xs-label text-white/70 mb-0.5">ELEVATION</div>
                <div className="text-base-value text-white">{elevation}m</div>
              </div>
            </div>
          )}

          {/* Stage Length */}
          {race?.stage_length_km && (
            <div className={cn('flex items-center gap-2', statMotion)}>
              <MapPin className="w-4 h-4 text-primary/70" />
              <div>
                <div className="text-xs-label text-white/70 mb-0.5">LENGTH</div>
                <div className="text-base-value text-white">{stageLength} km</div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
});

