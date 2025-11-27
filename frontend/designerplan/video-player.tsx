"use client"

import { useState, useEffect } from "react"
import { Play, Pause, Volume2, VolumeX, Maximize, Settings, Users } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Slider } from "@/components/ui/slider"

interface VideoPlayerProps {
  watchTime: number
  setWatchTime: (time: number | ((prev: number) => number)) => void
  addPoints: (amount: number) => void
}

export function VideoPlayer({ watchTime, setWatchTime, addPoints }: VideoPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(true)
  const [isMuted, setIsMuted] = useState(false)
  const [volume, setVolume] = useState([75])
  const [showControls, setShowControls] = useState(false)
  const [viewers] = useState(24853)

  // Simulate watch time and points earning
  useEffect(() => {
    if (!isPlaying) return

    const interval = setInterval(() => {
      setWatchTime((prev: number) => {
        const newTime = prev + 1
        // Award 10 points every minute
        if (newTime % 60 === 0) {
          addPoints(10)
        }
        return newTime
      })
    }, 1000)

    return () => clearInterval(interval)
  }, [isPlaying, setWatchTime, addPoints])

  const formatTime = (seconds: number) => {
    const h = Math.floor(seconds / 3600)
    const m = Math.floor((seconds % 3600) / 60)
    const s = seconds % 60
    return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`
  }

  return (
    <div
      className="relative aspect-video bg-black group"
      onMouseEnter={() => setShowControls(true)}
      onMouseLeave={() => setShowControls(false)}
    >
      {/* Video Placeholder */}
      <img src="/professional-cycling-peloton-racing-through-mounta.jpg" alt="Live cycling race" className="w-full h-full object-cover" />

      {/* Top gradient */}
      <div className="absolute inset-x-0 top-0 h-24 bg-gradient-to-b from-black/60 to-transparent pointer-events-none" />

      {/* Live Badge + Viewers */}
      <div className="absolute top-3 left-3 flex items-center gap-2">
        <div className="flex items-center gap-1.5 bg-red-600 pl-1.5 pr-2 py-0.5 rounded text-white">
          <span className="w-1.5 h-1.5 rounded-full bg-white animate-live-pulse" />
          <span className="text-[10px] font-bold uppercase tracking-wide">Live</span>
        </div>
        <div className="flex items-center gap-1 bg-black/50 backdrop-blur-sm px-2 py-0.5 rounded text-white/90">
          <Users className="w-3 h-3" />
          <span className="text-[10px] font-medium tabular-nums">{viewers.toLocaleString()}</span>
        </div>
      </div>

      {/* Race Info Badge */}
      <div className="absolute top-3 right-3 bg-black/50 backdrop-blur-sm px-2.5 py-1.5 rounded text-white">
        <p className="text-[10px] text-white/60 uppercase tracking-wide">Stage 17 · Tour de France</p>
        <p className="text-xs font-medium">Col du Galibier Summit</p>
      </div>

      {/* Bottom controls gradient */}
      <div
        className={`absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent transition-opacity duration-200 ${
          showControls ? "opacity-100" : "opacity-0"
        }`}
      >
        <div className="p-3 pt-12">
          {/* Progress indicator */}
          <div className="mb-2.5 flex items-center gap-2">
            <div className="flex-1 h-0.5 bg-white/20 rounded-full overflow-hidden">
              <div className="h-full bg-primary w-[68%] rounded-full" />
            </div>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1">
              <Button
                variant="ghost"
                size="icon"
                className="text-white hover:bg-white/10 w-8 h-8"
                onClick={() => setIsPlaying(!isPlaying)}
              >
                {isPlaying ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
              </Button>

              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-white hover:bg-white/10 w-8 h-8"
                  onClick={() => setIsMuted(!isMuted)}
                >
                  {isMuted ? <VolumeX className="w-4 h-4" /> : <Volume2 className="w-4 h-4" />}
                </Button>
                <div className="w-16 hidden sm:block">
                  <Slider
                    value={isMuted ? [0] : volume}
                    onValueChange={(v) => {
                      setVolume(v)
                      setIsMuted(v[0] === 0)
                    }}
                    max={100}
                    step={1}
                    className="cursor-pointer"
                  />
                </div>
              </div>

              <span className="text-[11px] text-white/80 ml-2 font-mono tabular-nums">{formatTime(watchTime)}</span>
            </div>

            <div className="flex items-center gap-0.5">
              <Button variant="ghost" size="icon" className="text-white hover:bg-white/10 w-8 h-8">
                <Settings className="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="icon" className="text-white hover:bg-white/10 w-8 h-8">
                <Maximize className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
