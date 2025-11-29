/**
 * Video player controls component
 */

import { memo } from 'react';
import { Play, Pause, Volume2, VolumeX, Maximize, Settings } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Slider } from '@/components/ui/slider';
import { VideoSettingsMenu } from './VideoSettingsMenu';
import { QualityLevel } from '@/hooks/useVideoPlayer';
import { useMemo } from 'react';

interface VideoControlsProps {
  showControls: boolean;
  isPlaying: boolean;
  isMuted: boolean;
  volume: number[];
  playbackSpeed: number;
  watchTime: number;
  currentTime: number;
  duration: number | null;
  isLive: boolean;
  qualityLevels: QualityLevel[];
  currentQuality: number;
  showSettingsMenu: boolean;
  onTogglePlay: () => void;
  onToggleMute: () => void;
  onVolumeChange: (volume: number[]) => void;
  onPlaybackSpeedChange: (speed: number) => void;
  onQualityChange: (level: number) => void;
  onSeek: (time: number) => void;
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
  currentTime,
  duration,
  isLive,
  qualityLevels,
  currentQuality,
  showSettingsMenu,
  onTogglePlay,
  onToggleMute,
  onVolumeChange,
  onPlaybackSpeedChange,
  onQualityChange,
  onSeek,
  onToggleFullscreen,
  onToggleSettingsMenu,
  onCloseMenus,
}: VideoControlsProps) {
  const sliderValue = useMemo(() => {
    if (!duration || !Number.isFinite(duration)) return [0];
    const clamped = Math.min(Math.max(currentTime, 0), duration);
    return [clamped];
  }, [currentTime, duration]);

  const formatClock = (seconds: number) => {
    if (!Number.isFinite(seconds)) return '00:00';
    const total = Math.max(0, Math.floor(seconds));
    const hrs = Math.floor(total / 3600);
    const mins = Math.floor((total % 3600) / 60);
    const secs = total % 60;
    const padded = `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    if (hrs > 0) {
      return `${hrs}:${padded}`;
    }
    return padded;
  };

  const timeLabel = isLive
    ? formatClock(Math.max(currentTime, watchTime))
    : `${formatClock(currentTime)} / ${formatClock(duration ?? 0)}`;

  return (
    <div
      className={`absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/70 via-black/35 to-transparent backdrop-blur-sm border-t border-white/10 transition-opacity duration-200 ${
        showControls ? 'opacity-100' : 'opacity-0'
      }`}
    >
      <div className="px-2 sm:px-3 pt-1.5 pb-2 flex flex-col gap-2 sm:gap-2.5">
        <div className="flex items-center gap-3">
          <div className="flex-1">
            {isLive || !duration ? (
              <div className="flex items-center gap-2">
                <div className="flex-1 h-1.5 bg-white/15 rounded-full overflow-hidden">
                  <div className="h-full bg-primary w-full animate-pulse" />
                </div>
                <div className="flex items-center gap-2 text-[11px] sm:text-xs text-white/90 font-mono">
                  <span className="px-2.5 py-1 rounded-full bg-primary text-white text-xs sm:text-sm font-semibold leading-tight">
                    LIVE
                  </span>
                  <span>{timeLabel}</span>
                </div>
              </div>
            ) : (
              <div className="flex items-center gap-3">
                <Slider
                  value={sliderValue}
                  max={duration || 0}
                  step={0.5}
                  onValueChange={(val) => onSeek(val[0] || 0)}
                  className="flex-1 cursor-pointer"
                />
                <div className="flex items-center gap-2 text-xs sm:text-sm text-white/90 font-mono">
                  <span>{timeLabel}</span>
                </div>
              </div>
            )}
          </div>
        </div>

          <div className="flex items-center justify-between gap-4">
          <div className="flex items-center gap-2 sm:gap-3">
            <Button
              variant="ghost"
              size="icon"
              className="text-white hover:bg-white/10 w-[2.125rem] h-[2.125rem] sm:w-9 sm:h-9"
              onClick={onTogglePlay}
            >
              {isPlaying ? <Pause className="w-5 h-5 sm:w-6 sm:h-6" /> : <Play className="w-5 h-5 sm:w-6 sm:h-6" />}
            </Button>

            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/10 w-[2.125rem] h-[2.125rem] sm:w-9 sm:h-9"
                onClick={onToggleMute}
              >
                {isMuted ? <VolumeX className="w-5 h-5 sm:w-6 sm:h-6" /> : <Volume2 className="w-5 h-5 sm:w-6 sm:h-6" />}
              </Button>

              <Slider
                value={volume}
                onValueChange={onVolumeChange}
                max={100}
                step={1}
                className="w-16 sm:w-24 md:w-28 cursor-pointer"
              />
            </div>

            {!isLive && (
              <span className="hidden sm:inline text-xs sm:text-sm text-white/80 font-mono tabular-nums">
                {timeLabel}
              </span>
            )}
          </div>

          <div className="flex items-center gap-2 sm:gap-3">
            <div className="relative">
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/10 w-[2.125rem] h-[2.125rem] sm:w-9 sm:h-9 flex items-center justify-center"
                onClick={onToggleSettingsMenu}
                >
                <Settings className="w-5 h-5 sm:w-6 sm:h-6" />
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
              className="text-white hover:bg-white/10 w-9 h-9 sm:w-10 sm:h-10"
              onClick={onToggleFullscreen}
            >
              <Maximize className="w-5 h-5 sm:w-6 sm:h-6" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
});
