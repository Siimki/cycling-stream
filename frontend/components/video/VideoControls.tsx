/**
 * Video player controls component
 */

import { memo } from 'react';
import { Play, Pause, Volume2, VolumeX, Maximize, Settings } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Slider } from '@/components/ui/slider';
import { formatTimeDetailed } from '@/lib/formatters';
import { VideoSpeedMenu } from './VideoSpeedMenu';
import { VideoQualityMenu } from './VideoQualityMenu';
import { QualityLevel } from '@/hooks/useVideoPlayer';

interface VideoControlsProps {
  showControls: boolean;
  isPlaying: boolean;
  isMuted: boolean;
  volume: number[];
  playbackSpeed: number;
  watchTime: number;
  qualityLevels: QualityLevel[];
  currentQuality: number;
  showSpeedMenu: boolean;
  showQualityMenu: boolean;
  onTogglePlay: () => void;
  onToggleMute: () => void;
  onVolumeChange: (volume: number[]) => void;
  onPlaybackSpeedChange: (speed: number) => void;
  onQualityChange: (level: number) => void;
  onToggleFullscreen: () => Promise<void>;
  onToggleSpeedMenu: () => void;
  onToggleQualityMenu: () => void;
  onCloseMenus: () => void;
}

export const VideoControls = memo(function VideoControls({
  showControls,
  isPlaying,
  isMuted,
  volume,
  playbackSpeed,
  watchTime,
  qualityLevels,
  currentQuality,
  showSpeedMenu,
  showQualityMenu,
  onTogglePlay,
  onToggleMute,
  onVolumeChange,
  onPlaybackSpeedChange,
  onQualityChange,
  onToggleFullscreen,
  onToggleSpeedMenu,
  onToggleQualityMenu,
  onCloseMenus,
}: VideoControlsProps) {
  return (
    <div
      className={`absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent transition-opacity duration-200 ${
        showControls ? 'opacity-100' : 'opacity-0'
      }`}
    >
      <div className="p-3 pt-12">
        {/* Progress indicator (Live stream so full) */}
        <div className="mb-2.5 flex items-center gap-2">
          <div className="flex-1 h-0.5 bg-white/20 rounded-full overflow-hidden">
            <div className="h-full bg-primary w-full rounded-full" />
          </div>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="text-white hover:bg-white/10 w-8 h-8"
              onClick={onTogglePlay}
            >
              {isPlaying ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
            </Button>

            <div className="flex items-center gap-1 group/volume">
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/10 w-8 h-8"
                onClick={onToggleMute}
              >
                {isMuted ? <VolumeX className="w-4 h-4" /> : <Volume2 className="w-4 h-4" />}
              </Button>
              <div className="w-0 overflow-hidden group-hover/volume:w-16 transition-all duration-300">
                <Slider
                  value={volume}
                  onValueChange={onVolumeChange}
                  max={100}
                  step={1}
                  className="cursor-pointer"
                />
              </div>
            </div>

            <span className="text-xs text-white/80 ml-2 font-mono tabular-nums">
              {formatTimeDetailed(watchTime)}
            </span>
          </div>

          <div className="flex items-center gap-0.5">
            {/* Quality Settings */}
            {qualityLevels.length > 0 && (
              <div className="relative">
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-white hover:bg-white/10 w-8 h-8"
                  onClick={onToggleQualityMenu}
                >
                  <Settings className="w-4 h-4" />
                </Button>
                <VideoQualityMenu
                  isOpen={showQualityMenu}
                  qualityLevels={qualityLevels}
                  currentQuality={currentQuality}
                  onQualityChange={onQualityChange}
                  onClose={onCloseMenus}
                />
              </div>
            )}

            {/* Speed Settings */}
            <div className="relative">
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/10 w-8 h-8"
                onClick={onToggleSpeedMenu}
              >
                <Settings className="w-4 h-4" />
              </Button>
              <VideoSpeedMenu
                isOpen={showSpeedMenu}
                currentSpeed={playbackSpeed}
                onSpeedChange={onPlaybackSpeedChange}
                onClose={onCloseMenus}
              />
            </div>

            <Button
              variant="ghost"
              size="icon"
              className="text-white hover:bg-white/10 w-8 h-8"
              onClick={onToggleFullscreen}
            >
              <Maximize className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
});

