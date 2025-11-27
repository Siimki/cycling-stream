/**
 * Video playback speed selection menu component
 */

import { memo } from 'react';
import { PLAYBACK_SPEEDS } from '@/constants/video';

interface VideoSpeedMenuProps {
  isOpen: boolean;
  currentSpeed: number;
  onSpeedChange: (speed: number) => void;
  onClose: () => void;
}

export const VideoSpeedMenu = memo(function VideoSpeedMenu({
  isOpen,
  currentSpeed,
  onSpeedChange,
  onClose,
}: VideoSpeedMenuProps) {
  if (!isOpen) return null;

  return (
    <div className="absolute bottom-full right-0 mb-2 bg-black/90 rounded overflow-hidden min-w-[100px] z-10">
      {PLAYBACK_SPEEDS.map((speed) => (
        <button
          key={speed}
          onClick={() => {
            onSpeedChange(speed);
            onClose();
          }}
          className={`block w-full text-left px-4 py-2 text-white hover:bg-white/20 transition-colors text-xs ${
            currentSpeed === speed ? 'text-primary' : ''
          }`}
        >
          {speed}x
        </button>
      ))}
    </div>
  );
});

