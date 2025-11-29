/**
 * Combined video settings menu for quality and playback speed.
 */

import { memo } from 'react';
import { PLAYBACK_SPEEDS } from '@/constants/video';
import { QualityLevel } from '@/hooks/useVideoPlayer';

interface VideoSettingsMenuProps {
  isOpen: boolean;
  qualityLevels: QualityLevel[];
  currentQuality: number;
  playbackSpeed: number;
  onQualityChange: (level: number) => void;
  onSpeedChange: (speed: number) => void;
  onClose: () => void;
}

export const VideoSettingsMenu = memo(function VideoSettingsMenu({
  isOpen,
  qualityLevels,
  currentQuality,
  playbackSpeed,
  onQualityChange,
  onSpeedChange,
  onClose,
}: VideoSettingsMenuProps) {
  if (!isOpen) return null;

  return (
    <div className="absolute bottom-full right-0 mb-2 bg-black/90 rounded-lg overflow-hidden min-w-[180px] z-10 border border-white/10 shadow-lg shadow-black/30">
      <div className="px-4 py-3 border-b border-white/10">
        <p className="text-[11px] uppercase tracking-[0.08em] text-white/60">Playback Speed</p>
        <div className="mt-2 grid grid-cols-3 gap-1">
          {PLAYBACK_SPEEDS.map((speed) => (
            <button
              key={speed}
              onClick={() => {
                onSpeedChange(speed);
                onClose();
              }}
              className={`text-xs px-2 py-2 rounded-md transition-colors ${
                playbackSpeed === speed
                  ? 'bg-primary/20 text-primary border border-primary/30'
                  : 'text-white hover:bg-white/10 border border-transparent'
              }`}
            >
              {speed}x
            </button>
          ))}
        </div>
      </div>

      <div className="px-4 py-3">
        <p className="text-[11px] uppercase tracking-[0.08em] text-white/60 mb-2">Quality</p>
        {qualityLevels.length === 0 ? (
          <p className="text-xs text-white/60">Auto (only stream quality available)</p>
        ) : (
          <div className="space-y-1">
            <button
              onClick={() => {
                onQualityChange(-1);
                onClose();
              }}
              className={`w-full text-left px-3 py-2 text-xs rounded-md transition-colors ${
                currentQuality === -1
                  ? 'bg-primary/20 text-primary border border-primary/30'
                  : 'text-white hover:bg-white/10 border border-transparent'
              }`}
            >
              Auto
            </button>
            {qualityLevels.map((level, index) => (
              <button
                key={index}
                onClick={() => {
                  onQualityChange(index);
                  onClose();
                }}
                className={`w-full text-left px-3 py-2 text-xs rounded-md transition-colors ${
                  currentQuality === index
                    ? 'bg-primary/20 text-primary border border-primary/30'
                    : 'text-white hover:bg-white/10 border border-transparent'
                }`}
              >
                {level.name}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
});
