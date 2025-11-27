/**
 * Custom hook for managing HLS video player state and lifecycle
 */

import { useEffect, useRef, useState, useCallback } from 'react';
import Hls from 'hls.js';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('useVideoPlayer');

export interface QualityLevel {
  height: number;
  width: number;
  bitrate: number;
  name: string;
}

export interface UseVideoPlayerReturn {
  videoRef: React.RefObject<HTMLVideoElement | null>;
  containerRef: React.RefObject<HTMLDivElement | null>;
  error: string | null;
  isBuffering: boolean;
  isPlaying: boolean;
  isMuted: boolean;
  volume: number[];
  isFullscreen: boolean;
  playbackSpeed: number;
  qualityLevels: QualityLevel[];
  currentQuality: number;
  hasHls: boolean;
  watchTime: number;
  setIsPlaying: (playing: boolean) => void;
  setIsMuted: (muted: boolean) => void;
  setVolume: (volume: number[]) => void;
  setPlaybackSpeed: (speed: number) => void;
  setCurrentQuality: (quality: number) => void;
  togglePlay: () => void;
  toggleMute: () => void;
  handleVolumeChange: (newVolume: number[]) => void;
  handlePlaybackSpeedChange: (speed: number) => void;
  handleQualityChange: (level: number) => void;
  toggleFullscreen: () => Promise<void>;
}

export function useVideoPlayer(streamUrl?: string, status: string = 'offline'): UseVideoPlayerReturn {
  const videoRef = useRef<HTMLVideoElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  
  const [error, setError] = useState<string | null>(null);
  const [isBuffering, setIsBuffering] = useState(false);
  const [isPlaying, setIsPlaying] = useState(true);
  const [isMuted, setIsMuted] = useState(false);
  const [volume, setVolume] = useState([100]);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [playbackSpeed, setPlaybackSpeed] = useState(1);
  const [qualityLevels, setQualityLevels] = useState<QualityLevel[]>([]);
  const [currentQuality, setCurrentQuality] = useState<number>(-1);
  const [hasHls, setHasHls] = useState(false);
  const [watchTime, setWatchTime] = useState(0);

  // Watch time counter
  useEffect(() => {
    if (!isPlaying || isBuffering) return;
    const interval = setInterval(() => {
      setWatchTime(prev => prev + 1);
    }, 1000);
    return () => clearInterval(interval);
  }, [isPlaying, isBuffering]);

  // Handle fullscreen changes
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  const toggleFullscreen = useCallback(async () => {
    if (!containerRef.current) return;

    try {
      if (!document.fullscreenElement) {
        await containerRef.current.requestFullscreen();
      } else {
        await document.exitFullscreen();
      }
    } catch (err) {
      logger.error('Fullscreen error:', err);
    }
  }, []);

  // Handle Play/Pause
  const togglePlay = useCallback(() => {
    if (videoRef.current) {
      if (videoRef.current.paused) {
        videoRef.current.play();
      } else {
        videoRef.current.pause();
      }
    }
  }, []);

  // Handle Volume
  const handleVolumeChange = useCallback((newVolume: number[]) => {
    if (videoRef.current) {
      const vol = newVolume[0] / 100;
      videoRef.current.volume = vol;
      setVolume(newVolume);
      setIsMuted(vol === 0);
      videoRef.current.muted = vol === 0;
    }
  }, []);

  const toggleMute = useCallback(() => {
    if (videoRef.current) {
      videoRef.current.muted = !videoRef.current.muted;
      setIsMuted(videoRef.current.muted);
      if (videoRef.current.muted) {
        setVolume([0]);
      } else {
        setVolume([videoRef.current.volume * 100]);
      }
    }
  }, []);

  const handlePlaybackSpeedChange = useCallback((speed: number) => {
    if (videoRef.current) {
      videoRef.current.playbackRate = speed;
      setPlaybackSpeed(speed);
    }
  }, []);

  const handleQualityChange = useCallback((level: number) => {
    if (hlsRef.current && level >= 0) {
      hlsRef.current.currentLevel = level;
      setCurrentQuality(level);
    } else if (hlsRef.current && level === -1) {
      // Auto quality
      hlsRef.current.currentLevel = -1;
      setCurrentQuality(-1);
    }
  }, []);

  // HLS initialization and management
  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => {
    if (!videoRef.current) return;

    // Cleanup previous HLS instance
    if (hlsRef.current) {
      hlsRef.current.destroy();
      hlsRef.current = null;
    }

    // Buffering state - event handlers are fine, setState in handlers is valid
    const handleWaiting = () => setIsBuffering(true);
    const handleCanPlay = () => setIsBuffering(false);
    const handlePlaying = () => {
      setIsBuffering(false);
      setIsPlaying(true);
    };
    const handlePause = () => setIsPlaying(false);

    videoRef.current.addEventListener('waiting', handleWaiting);
    videoRef.current.addEventListener('canplay', handleCanPlay);
    videoRef.current.addEventListener('playing', handlePlaying);
    videoRef.current.addEventListener('pause', handlePause);

    if (status === 'live' && streamUrl) {
      if (Hls.isSupported()) {
        const hls = new Hls({
          enableWorker: true,
          lowLatencyMode: true,
          backBufferLength: 90,
        });

        hls.loadSource(streamUrl);
        hls.attachMedia(videoRef.current);

        hls.on(Hls.Events.MANIFEST_PARSED, () => {
          if (videoRef.current) {
            // Get quality levels
            const levels = hls.levels.map((level) => ({
              height: level.height,
              width: level.width,
              bitrate: level.bitrate,
              name: level.height ? `${level.height}p` : 'Auto',
            }));
            setQualityLevels(levels);
            setCurrentQuality(hls.currentLevel);

            videoRef.current.play().catch((err) => {
              logger.error('Error playing video:', err);
              setIsPlaying(false);
            });
          }
        });

        hls.on(Hls.Events.LEVEL_SWITCHED, () => {
          setCurrentQuality(hls.currentLevel);
        });

        hls.on(Hls.Events.ERROR, (event, data) => {
          if (data.fatal) {
            switch (data.type) {
              case Hls.ErrorTypes.NETWORK_ERROR:
                logger.error('Network error, trying to recover...');
                hls.startLoad();
                break;
              case Hls.ErrorTypes.MEDIA_ERROR:
                logger.error('Media error, trying to recover...');
                hls.recoverMediaError();
                break;
              default:
                logger.error('Fatal error, destroying HLS instance');
                hls.destroy();
                setTimeout(() => setError('Stream error occurred'), 0);
                break;
            }
          }
        });

        hlsRef.current = hls;
        // Set HLS state after initialization - valid pattern for setup
        setHasHls(true);
      } else if (videoRef.current.canPlayType('application/vnd.apple.mpegurl')) {
        // Native HLS support (Safari)
        videoRef.current.src = streamUrl;
        videoRef.current.play().catch((err) => {
          logger.error('Error playing video:', err);
          setIsPlaying(false);
        });
      } else {
        setTimeout(() => setError('HLS playback not supported in this browser'), 0);
      }
    }

    return () => {
      if (videoRef.current) {
        videoRef.current.removeEventListener('waiting', handleWaiting);
        videoRef.current.removeEventListener('canplay', handleCanPlay);
        videoRef.current.removeEventListener('playing', handlePlaying);
        videoRef.current.removeEventListener('pause', handlePause);
      }
      if (hlsRef.current) {
        hlsRef.current.destroy();
        hlsRef.current = null;
        // Cleanup state - valid pattern for cleanup
        setHasHls(false);
      }
    };
  }, [streamUrl, status]); // Event handlers are stable, no need in deps

  return {
    videoRef,
    containerRef,
    error,
    isBuffering,
    isPlaying,
    isMuted,
    volume,
    isFullscreen,
    playbackSpeed,
    qualityLevels,
    currentQuality,
    hasHls,
    watchTime,
    setIsPlaying,
    setIsMuted,
    setVolume,
    setPlaybackSpeed,
    setCurrentQuality,
    togglePlay,
    toggleMute,
    handleVolumeChange,
    handlePlaybackSpeedChange,
    handleQualityChange,
    toggleFullscreen,
  };
}

