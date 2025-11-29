'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { useVideoPlayer } from '@/hooks/useVideoPlayer';
import { useVideoKeyboardShortcuts } from '@/hooks/useVideoKeyboardShortcuts';
import { VideoOverlay } from './VideoOverlay';
import { VideoControls } from './VideoControls';
import { VIDEO_CONTROLS_HIDE_DELAY_MS } from '@/constants/intervals';
import { useAnalyticsTracking } from '@/hooks/useAnalyticsTracking';

interface VideoPlayerProps {
  streamUrl?: string;
  status: string;
  streamType?: string;
  sourceId?: string;
  streamId?: string;
  provider?: string;
  requiresLogin?: boolean; // Kept for backward compatibility but not used
  raceId?: string;
}

export default function VideoPlayer({ streamUrl, status, streamType, sourceId, streamId }: VideoPlayerProps) {
  const [showControls, setShowControls] = useState(false);
  const [showSettingsMenu, setShowSettingsMenu] = useState(false);
  const controlsTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastErrorRef = useRef<string | null>(null);
  const lastBufferingRef = useRef<boolean>(false);

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

  const { trackPlay, trackPause, trackHeartbeat, trackEnded, trackError, trackBufferStart, trackBufferEnd } =
    useAnalyticsTracking(streamId);
  const isYouTube = status === 'live' && streamType === 'youtube' && !!sourceId;
  // Check if this is a Bunny Stream embed URL (player.mediadelivery.net/embed)
  // vs HLS URL (stream.mediadelivery.net/hls)
  const isBunnyStreamEmbed = streamUrl && streamUrl.includes('player.mediadelivery.net/embed');
  const isBunnyStreamHLS = streamUrl && (streamUrl.includes('stream.mediadelivery.net/hls') || streamUrl.includes('.m3u8'));

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
      setShowSettingsMenu(false);
    },
    [handlePlaybackSpeedChange]
  );

  const handleQualityChangeWithClose = useCallback(
    (level: number) => {
      handleQualityChange(level);
      setShowSettingsMenu(false);
    },
    [handleQualityChange]
  );

  const handleCloseMenus = useCallback(() => {
    setShowSettingsMenu(false);
  }, []);

  // Fire tracking events for YouTube (limited without iframe API)
  useEffect(() => {
    if (!isYouTube || !streamId) return;
    trackPlay(0);
    const interval = window.setInterval(() => {
      trackHeartbeat();
    }, 15000);
    return () => {
      window.clearInterval(interval);
    };
  }, [isYouTube, streamId, trackHeartbeat, trackPlay]);

  // Send errors once per occurrence
  useEffect(() => {
    if (!streamId || !error) return;
    if (error === lastErrorRef.current) return;
    lastErrorRef.current = error;
    trackError(undefined, { message: error });
  }, [error, streamId, trackError]);

  // Hook up HTML5 video events for analytics
  useEffect(() => {
    const videoEl = videoRef.current;
    if (!videoEl || !streamId || isYouTube) {
      return;
    }

    const handlePlay = () => trackPlay(Math.floor(videoEl.currentTime || 0));
    const handlePause = () => trackPause(Math.floor(videoEl.currentTime || 0));
    const handleEnded = () => trackEnded(Math.floor(videoEl.currentTime || 0));
    const handleError = () => trackError(Math.floor(videoEl.currentTime || 0), { message: 'video_error' });

    videoEl.addEventListener('play', handlePlay);
    videoEl.addEventListener('pause', handlePause);
    videoEl.addEventListener('ended', handleEnded);
    videoEl.addEventListener('error', handleError);

    return () => {
      videoEl.removeEventListener('play', handlePlay);
      videoEl.removeEventListener('pause', handlePause);
      videoEl.removeEventListener('ended', handleEnded);
      videoEl.removeEventListener('error', handleError);
    };
  }, [isYouTube, streamId, trackEnded, trackError, trackPause, trackPlay, videoRef]);

  // Heartbeat loop for HTML5 playback
  useEffect(() => {
    if (!streamId || isYouTube) return;
    const interval = window.setInterval(() => {
      const videoEl = videoRef.current;
      if (!videoEl || videoEl.paused || status !== 'live') {
        return;
      }
      trackHeartbeat(Math.floor(videoEl.currentTime || 0));
    }, 15000);

    return () => {
      window.clearInterval(interval);
    };
  }, [isYouTube, status, streamId, trackHeartbeat, videoRef]);

  // Buffer tracking for HTML5 playback
  useEffect(() => {
    if (!streamId || isYouTube) return;
    const videoEl = videoRef.current;
    const wasBuffering = lastBufferingRef.current;
    if (isBuffering && !wasBuffering) {
      const position = Math.floor(videoEl?.currentTime || 0);
      trackBufferStart(position);
    }
    if (!isBuffering && wasBuffering) {
      const position = Math.floor(videoEl?.currentTime || 0);
      trackBufferEnd(position);
    }
    lastBufferingRef.current = isBuffering;
  }, [isBuffering, isYouTube, streamId, trackBufferStart, trackBufferEnd, videoRef]);

  // YouTube Player
  if (isYouTube && sourceId) {
    return (
      <div className="aspect-video w-full h-full bg-black rounded-lg overflow-hidden border border-border/50 shadow-lg shadow-black/20">
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

  // Bunny Stream Embed Player (for embed player only)
  // Only use embed player if the URL is explicitly an embed URL
  if (isBunnyStreamEmbed && streamUrl) {
    // Extract the embed URL - if it's already a full URL, use it directly
    // Otherwise construct it from sourceId if available
    let embedUrl = streamUrl;
    
    // If the URL contains query parameters, ensure they're included
    // The embed URL should already be complete from the backend
    return (
      <div className="aspect-video w-full h-full bg-black rounded-lg overflow-hidden border border-border/50 shadow-lg shadow-black/20">
        <div style={{ position: 'relative', paddingTop: '56.25%' }}>
          <iframe
            src={embedUrl}
            loading="lazy"
            style={{ border: 0, position: 'absolute', top: 0, height: '100%', width: '100%' }}
            allow="accelerometer;gyroscope;autoplay;encrypted-media;picture-in-picture;"
            allowFullScreen={true}
            title="Bunny Stream video player"
            className="w-full h-full"
          ></iframe>
        </div>
      </div>
    );
  }

  // Show offline message only if there's no stream URL available
  // Allow playback for both live and offline/replay streams if URL exists
  if (!streamUrl) {
    return (
      <div className="bg-card aspect-video flex items-center justify-center rounded-lg relative overflow-hidden border border-border/50 shadow-lg shadow-black/20">
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
      <div className="bg-card aspect-video flex items-center justify-center rounded-lg relative overflow-hidden border border-border/50 shadow-lg shadow-black/20">
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
      className="relative aspect-video bg-black group rounded-lg overflow-hidden border border-border/50 shadow-lg shadow-black/20"
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
        showSettingsMenu={showSettingsMenu}
        onTogglePlay={togglePlay}
        onToggleMute={toggleMute}
        onVolumeChange={handleVolumeChange}
        onPlaybackSpeedChange={handlePlaybackSpeedChangeWithClose}
        onQualityChange={handleQualityChangeWithClose}
        onToggleFullscreen={toggleFullscreen}
        onToggleSettingsMenu={() => setShowSettingsMenu(!showSettingsMenu)}
        onCloseMenus={handleCloseMenus}
      />

      {/* Click outside to close menus */}
      {showSettingsMenu && (
        <div className="absolute inset-0 z-[5]" onClick={handleCloseMenus} />
      )}
    </div>
  );
}
