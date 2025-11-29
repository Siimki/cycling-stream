/**
 * Video player controls component
 */

import { memo } from 'react';
import { Play, Pause, Volume2, VolumeX, Maximize, Settings } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Slider } from '@/components/ui/slider';
import { formatTimeDetailed } from '@/lib/formatters';
import { VideoSettingsMenu } from './VideoSettingsMenu';
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
  showSettingsMenu: boolean;
  onTogglePlay: () => void;
  onToggleMute: () => void;
  onVolumeChange: (volume: number[]) => void;
  onPlaybackSpeedChange: (speed: number) => void;
  onQualityChange: (level: number) => void;
  onToggleFullscreen: () => Promise<void>;
  onToggleSettingsMenu: () => void;
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
  showSettingsMenu,
  onTogglePlay,
  onToggleMute,
  onVolumeChange,
  onPlaybackSpeedChange,
  onQualityChange,
  onToggleFullscreen,
  onToggleSettingsMenu,
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
            <div className="relative">
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/10 w-8 h-8"
                onClick={onToggleSettingsMenu}
              >
                <Settings className="w-4 h-4" />
              </Button>
              <VideoSettingsMenu
                isOpen={showSettingsMenu}
                qualityLevels={qualityLevels}
                currentQuality={currentQuality}
                playbackSpeed={playbackSpeed}
                onQualityChange={onQualityChange}
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
