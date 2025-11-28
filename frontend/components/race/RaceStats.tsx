"use client"

import { useState } from "react"
import { Mountain, Timer, MapPin } from "lucide-react"
import { HudToggleButton } from "@/components/user/HudToggleButton"
import { Race } from "@/lib/api"

interface RaceStatsProps {
  race?: Race
  compact?: boolean
}

// Format number with commas (e.g., 4800 -> "4,800")
function formatNumber(num: number | undefined | null): string {
  if (num === undefined || num === null) return "N/A"
  return num.toLocaleString()
}

// Format time from "HH:MM:SS" or "HH:MM" to "~HH:MM"
function formatTime(time: string | undefined | null): string {
  if (!time) return "N/A"
  // Extract HH:MM from time string
  const match = time.match(/(\d{1,2}):(\d{2})/)
  if (match) {
    return `~${match[1]}:${match[2]}`
  }
  return time
}

export function RaceStats({ race, compact = false }: RaceStatsProps) {
  const [isVisible, setIsVisible] = useState(true);

  if (compact) return null;

  // Get display values with fallbacks
  const raceName = race?.name || "Race"
  const stageName = race?.stage_name || ""
  const stageType = race?.stage_type || ""
  const stageDisplay = stageName && stageType 
    ? `${stageName} — ${stageType}`
    : stageName || stageType || "Stage"
  const elevation = formatNumber(race?.elevation_meters)
  const finishTime = formatTime(race?.estimated_finish_time)
  const stageLength = formatNumber(race?.stage_length_km)

  return (
    <div className="relative flex flex-col w-full bg-card/80 backdrop-blur-sm border-t border-border/50 overflow-hidden min-h-[4.5rem]">
      {/* Placeholder bar when collapsed - Always maintains space */}
      <div className={`transition-all duration-300 ease-in-out ${!isVisible ? "h-14 opacity-100" : "h-0 opacity-0 overflow-hidden"}`}>
        <div className="h-14 w-full flex items-center px-4">
          <span className="text-xs text-muted-foreground/50">Race info collapsed</span>
        </div>
      </div>

      {/* Toggle Button - Always visible, positioned above content */}
      <HudToggleButton
        isActive={isVisible}
        label="Race info"
        onToggle={() => setIsVisible(!isVisible)}
        className="top-2 right-3"
      />

      {/* Collapsible Content */}
      <div
        className={`flex flex-col md:flex-row md:items-center w-full overflow-x-auto transition-all duration-300 ease-in-out ${
          isVisible ? "max-h-[300px] opacity-100" : "max-h-0 opacity-0 overflow-hidden"
        }`}
      >
      {/* Stage Name & Distance - Mobile: Full width, Desktop: Separate sections */}
      <div className="flex items-center justify-between px-5 md:px-6 py-4 md:py-2 md:border-r border-border/30 md:min-w-[240px] lg:min-w-[300px]">
        <div className="flex-1 min-w-0">
          <p className="text-sm-label text-primary mb-1.5 truncate">{raceName}</p>
          <h3 className="text-lg md:text-xl font-bold tracking-tight truncate text-foreground">{stageDisplay}</h3>
        </div>
        {/* Mobile: Show stage length inline */}
        {race?.stage_length_km && (
          <div className="md:hidden flex items-center gap-2 ml-3 bg-muted/50 px-3 py-1.5 rounded-md">
            <MapPin className="w-5 h-5 shrink-0 text-primary" />
            <span className="text-base font-bold tabular-nums text-foreground">{stageLength} km</span>
          </div>
        )}
      </div>

      {/* Stats Grid - Hidden on mobile, shown on tablet+ */}
      {(race?.elevation_meters || race?.estimated_finish_time) && (
        <div className="hidden md:flex items-center px-8 gap-8 lg:gap-12 border-r border-border/30">
          {race?.elevation_meters && (
            <div className="py-4">
              <div className="flex items-center gap-2.5 text-muted-foreground mb-1.5">
                <Mountain className="w-5 h-5 text-primary/70" />
                <span className="text-sm-label">Elevation</span>
              </div>
              <p className="text-xl-heading text-foreground tabular-nums">{elevation}m</p>
            </div>
          )}
          {race?.estimated_finish_time && (
            <div className="py-4">
              <div className="flex items-center gap-2.5 text-muted-foreground mb-1.5">
                <Timer className="w-5 h-5 text-primary/70" />
                <span className="text-sm-label">Est. Finish</span>
              </div>
              <p className="text-xl-heading text-foreground tabular-nums">{finishTime}</p>
            </div>
          )}
        </div>
      )}

      {/* Stage Length - Hidden on mobile (shown inline above), full on tablet+ */}
      {race?.stage_length_km && (
        <div className="hidden md:flex items-center px-8 py-4 min-w-[160px]">
          <div className="flex items-center gap-3">
            <MapPin className="w-5 h-5 text-primary/70" />
            <div>
              <div className="text-sm-label text-muted-foreground mb-1.5">
                Stage Length
              </div>
              <p className="text-xl-heading text-foreground tabular-nums">{stageLength} km</p>
            </div>
          </div>
        </div>
      )}

      {/* Mobile-only: Compact stats row */}
      {(race?.elevation_meters || race?.estimated_finish_time) && (
        <div className="flex md:hidden items-center justify-between px-5 py-4 border-t border-border/20 bg-muted/20">
          <div className="flex items-center gap-6">
            {race?.elevation_meters && (
              <div className="flex items-center gap-2.5">
                <Mountain className="w-5 h-5 text-primary/70" />
                <span className="text-base font-bold tabular-nums text-foreground">{elevation}m</span>
              </div>
            )}
            {race?.estimated_finish_time && (
              <div className="flex items-center gap-2.5">
                <Timer className="w-5 h-5 text-primary/70" />
                <span className="text-base font-bold tabular-nums text-foreground">{finishTime}</span>
              </div>
            )}
          </div>
        </div>
      )}
      </div>
    </div>
  )
}
