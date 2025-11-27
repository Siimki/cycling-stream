"use client"

import { useState } from "react"
import { Mountain, Timer, MapPin } from "lucide-react"
import { HudToggleButton } from "@/components/user/HudToggleButton"

interface RaceStatsProps {
  compact?: boolean
}

export function RaceStats({ compact = false }: RaceStatsProps) {
  const [isVisible, setIsVisible] = useState(true);

  if (compact) return null;

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
      <div className="flex items-center justify-between px-4 md:px-5 lg:px-6 py-2.5 md:py-0 md:border-r border-border/30 md:min-w-[180px] lg:min-w-[220px]">
        <div className="flex-1 min-w-0">
          <p className="text-[0.65rem] md:text-xs text-primary uppercase tracking-wider font-semibold mb-0.5 truncate">Tour de France 2025</p>
          <h3 className="text-sm md:text-base lg:text-lg font-bold tracking-tight truncate text-foreground/95">Stage 17 — Mountain</h3>
        </div>
        {/* Mobile: Show stage length inline */}
        <div className="md:hidden flex items-center gap-1.5 ml-3 bg-muted/50 px-2 py-1 rounded-md">
          <MapPin className="w-3.5 h-3.5 shrink-0 text-primary" />
          <span className="text-xs font-bold tabular-nums text-foreground">166 km</span>
        </div>
      </div>

      {/* Stats Grid - Hidden on mobile, shown on tablet+ */}
      <div className="hidden md:flex items-center px-4 lg:px-6 gap-5 lg:gap-8 border-r border-border/30">
        <div className="py-2">
          <div className="flex items-center gap-1.5 text-muted-foreground mb-0.5">
            <Mountain className="w-3.5 h-3.5 lg:w-4 lg:h-4 text-primary/70" />
            <span className="text-[0.65rem] lg:text-xs font-semibold uppercase tracking-wider">Elevation</span>
          </div>
          <p className="text-base lg:text-xl font-bold tabular-nums text-foreground/95">4,800m</p>
        </div>
        <div className="py-2">
          <div className="flex items-center gap-1.5 text-muted-foreground mb-0.5">
            <Timer className="w-3.5 h-3.5 lg:w-4 lg:h-4 text-primary/70" />
            <span className="text-[0.65rem] lg:text-xs font-semibold uppercase tracking-wider">Est. Finish</span>
          </div>
          <p className="text-base lg:text-xl font-bold tabular-nums text-foreground/95">~17:45</p>
        </div>
      </div>

      {/* Stage Length - Hidden on mobile (shown inline above), full on tablet+ */}
      <div className="hidden md:flex items-center px-4 lg:px-6 py-2 min-w-[120px]">
        <div className="flex items-center gap-1.5">
          <MapPin className="w-3.5 h-3.5 lg:w-4 lg:h-4 text-primary/70" />
          <div>
            <div className="text-[0.65rem] lg:text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-0.5">
              Stage Length
            </div>
            <p className="text-base lg:text-xl font-bold tabular-nums text-foreground/95">166 km</p>
          </div>
        </div>
      </div>

      {/* Mobile-only: Compact stats row */}
      <div className="flex md:hidden items-center justify-between px-4 py-2 border-t border-border/20 bg-muted/20">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1.5">
            <Mountain className="w-3.5 h-3.5 text-primary/70" />
            <span className="text-xs font-bold tabular-nums">4,800m</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Timer className="w-3.5 h-3.5 text-primary/70" />
            <span className="text-xs font-bold tabular-nums">~17:45</span>
          </div>
        </div>
      </div>
      </div>
    </div>
  )
}
