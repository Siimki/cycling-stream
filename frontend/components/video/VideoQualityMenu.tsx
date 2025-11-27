/**
 * Video quality selection menu component
 */

import { memo } from 'react';
import { QualityLevel } from '@/hooks/useVideoPlayer';

interface VideoQualityMenuProps {
  isOpen: boolean;
  qualityLevels: QualityLevel[];
  currentQuality: number;
  onQualityChange: (level: number) => void;
  onClose: () => void;
}

export const VideoQualityMenu = memo(function VideoQualityMenu({
  isOpen,
  qualityLevels,
  currentQuality,
  onQualityChange,
  onClose,
}: VideoQualityMenuProps) {
  if (!isOpen || qualityLevels.length === 0) return null;

  return (
    <div className="absolute bottom-full right-0 mb-2 bg-black/90 rounded overflow-hidden min-w-[120px] z-10">
      <button
        onClick={() => {
          onQualityChange(-1);
          onClose();
        }}
        className={`block w-full text-left px-4 py-2 text-white hover:bg-white/20 transition-colors text-xs ${
          currentQuality === -1 ? 'text-primary' : ''
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
          className={`block w-full text-left px-4 py-2 text-white hover:bg-white/20 transition-colors text-xs ${
            currentQuality === index ? 'text-primary' : ''
          }`}
        >
          {level.name}
        </button>
      ))}
    </div>
  );
});

