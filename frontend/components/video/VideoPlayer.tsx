'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { useVideoPlayer } from '@/hooks/useVideoPlayer';
import { useVideoKeyboardShortcuts } from '@/hooks/useVideoKeyboardShortcuts';
import { VideoOverlay } from './VideoOverlay';
import { VideoControls } from './VideoControls';
import { VIDEO_CONTROLS_HIDE_DELAY_MS } from '@/constants/intervals';

interface VideoPlayerProps {
  streamUrl?: string;
  status: string;
  streamType?: string;
  sourceId?: string;
  requiresLogin?: boolean; // Kept for backward compatibility but not used
}

export default function VideoPlayer({ streamUrl, status, streamType, sourceId }: VideoPlayerProps) {
  const [showControls, setShowControls] = useState(false);
  const [showQualityMenu, setShowQualityMenu] = useState(false);
  const [showSpeedMenu, setShowSpeedMenu] = useState(false);
  const controlsTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const {
    videoRef,
    containerRef,
    error,
    isBuffering,
    isPlaying,
    isMuted,
    volume,
    playbackSpeed,
    qualityLevels,
    currentQuality,
    watchTime,
    togglePlay,
    toggleMute,
    handleVolumeChange,
    handlePlaybackSpeedChange,
    handleQualityChange,
    toggleFullscreen,
  } = useVideoPlayer(streamUrl, status);

  // Keyboard shortcuts - must be called before any conditional returns
  useVideoKeyboardShortcuts({
    videoRef,
    togglePlay,
    toggleFullscreen,
    toggleMute,
  });

  // Auto-hide controls - must be called before any conditional returns
  useEffect(() => {
    const resetControlsTimeout = () => {
      if (controlsTimeoutRef.current) {
        clearTimeout(controlsTimeoutRef.current);
      }
      setShowControls(true);
      controlsTimeoutRef.current = setTimeout(() => {
        if (!videoRef.current?.paused) {
          setShowControls(false);
        }
      }, VIDEO_CONTROLS_HIDE_DELAY_MS);
    };

    const container = containerRef.current;
    if (container) {
      container.addEventListener('mousemove', resetControlsTimeout);
      container.addEventListener('mouseenter', resetControlsTimeout);
    }

    return () => {
      if (controlsTimeoutRef.current) {
        clearTimeout(controlsTimeoutRef.current);
      }
      if (container) {
        container.removeEventListener('mousemove', resetControlsTimeout);
        container.removeEventListener('mouseenter', resetControlsTimeout);
      }
    };
  }, [containerRef, videoRef]);

  // Callbacks - must be called before any conditional returns
  const handlePlaybackSpeedChangeWithClose = useCallback(
    (speed: number) => {
      handlePlaybackSpeedChange(speed);
      setShowSpeedMenu(false);
    },
    [handlePlaybackSpeedChange]
  );

  const handleQualityChangeWithClose = useCallback(
    (level: number) => {
      handleQualityChange(level);
      setShowQualityMenu(false);
    },
    [handleQualityChange]
  );

  const handleCloseMenus = useCallback(() => {
    setShowQualityMenu(false);
    setShowSpeedMenu(false);
  }, []);

  // YouTube Player
  if (status === 'live' && streamType === 'youtube' && sourceId) {
    return (
      <div className="aspect-video w-full h-full bg-black rounded-lg overflow-hidden border border-border">
         <iframe
            width="100%"
            height="100%"
            src={`https://www.youtube.com/embed/${sourceId}?autoplay=1`}
            title="YouTube video player"
            frameBorder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            allowFullScreen
            className="w-full h-full"
          ></iframe>
      </div>
    );
  }

  if (status !== 'live' || !streamUrl) {
    // Note: Authentication checks are handled by AuthRequiredWrapper
    // This component only shows the stream or offline message
    return (
      <div className="bg-card aspect-video flex items-center justify-center rounded-lg relative overflow-hidden border border-border">
        <div className="absolute inset-0 bg-gradient-to-br from-background to-card"></div>
        <div className="relative text-center text-foreground z-10 px-4">
          <div className="text-6xl mb-4">📺</div>
          <p className="text-2xl font-semibold mb-2">Stream Offline</p>
          <p className="text-muted-foreground">
            The stream is not currently available. Please check back later.
          </p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-card aspect-video flex items-center justify-center rounded-lg relative overflow-hidden border border-border">
        <div className="absolute inset-0 bg-gradient-to-br from-destructive/20 to-card"></div>
        <div className="relative text-center text-foreground z-10 px-4">
          <div className="text-6xl mb-4">⚠️</div>
          <p className="text-2xl font-semibold mb-2 text-destructive">Stream Error</p>
          <p className="text-muted-foreground">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div
      ref={containerRef}
      className="relative aspect-video bg-black group rounded-lg overflow-hidden border border-border"
    >
      {/* Video element */}
      <video ref={videoRef} className="w-full h-full object-cover" playsInline />

      {/* Buffering indicator */}
      {isBuffering && (
        <div className="absolute inset-0 flex items-center justify-center bg-black/50 z-20">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        </div>
      )}

      {/* Overlay (Live badge, gradients) */}
      <VideoOverlay showControls={showControls} />

      {/* Video Controls */}
      <VideoControls
        showControls={showControls}
        isPlaying={isPlaying}
        isMuted={isMuted}
        volume={volume}
        playbackSpeed={playbackSpeed}
        watchTime={watchTime}
        qualityLevels={qualityLevels}
        currentQuality={currentQuality}
        showSpeedMenu={showSpeedMenu}
        showQualityMenu={showQualityMenu}
        onTogglePlay={togglePlay}
        onToggleMute={toggleMute}
        onVolumeChange={handleVolumeChange}
        onPlaybackSpeedChange={handlePlaybackSpeedChangeWithClose}
        onQualityChange={handleQualityChangeWithClose}
        onToggleFullscreen={toggleFullscreen}
        onToggleSpeedMenu={() => setShowSpeedMenu(!showSpeedMenu)}
        onToggleQualityMenu={() => setShowQualityMenu(!showQualityMenu)}
        onCloseMenus={handleCloseMenus}
      />

      {/* Click outside to close menus */}
      {(showQualityMenu || showSpeedMenu) && (
        <div className="absolute inset-0 z-[5]" onClick={handleCloseMenus} />
      )}
    </div>
  );
}
