"use client"

import { useMemo, useState, type ReactNode } from "react"
import { Mountain, Timer, MapPin } from "lucide-react"
import { Race } from "@/lib/api"
import { POINTS_TIERS } from "@/constants/tiers"
import { useHudStats } from "@/components/user/HudStatsProvider"
import { useAuth } from "@/contexts/AuthContext"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

interface RaceStatsProps {
  race?: Race
  compact?: boolean
}

function formatNumber(num: number | undefined | null): string {
  if (num === undefined || num === null) return "N/A"
  return num.toLocaleString()
}

function formatTime(time: string | undefined | null): string {
  if (!time) return "N/A"
  const match = time.match(/(\d{1,2}):(\d{2})/)
  if (match) {
    return `~${match[1]}:${match[2]}`
  }
  return time
}

export function RaceStats({ race, compact = false }: RaceStatsProps) {
  const [isOpen, setIsOpen] = useState(true)
  const { isAuthenticated } = useAuth()
  const { points } = useHudStats()

  if (compact) return null

  const raceName = race?.name || "Race"
  const stageName = race?.stage_name || ""
  const stageType = race?.stage_type || ""
  const stageDisplay = stageName && stageType ? `${stageName} — ${stageType}` : stageName || stageType || "Stage"
  const elevation = formatNumber(race?.elevation_meters)
  const finishTime = formatTime(race?.estimated_finish_time)
  const stageLength = formatNumber(race?.stage_length_km)

  const profilePoints = useMemo(() => createProfileData(race?.stage_length_km, race?.elevation_meters), [
    race?.stage_length_km,
    race?.elevation_meters,
  ])

  const { currentTier, nextTier, progressToNext } = useMemo(() => {
    const tiers = POINTS_TIERS
    const current = tiers.reduce((acc, tier) => (points >= tier.min ? tier : acc), tiers[0])
    const next = tiers[tiers.indexOf(current) + 1]
    const progress = next ? ((points - current.min) / (next.min - current.min)) * 100 : 100
    return { currentTier: current, nextTier: next, progressToNext: progress }
  }, [points])

  return (
    <div className="relative w-full bg-[#121212] border-t border-[#222] rounded-b-xl overflow-hidden shadow-[0_-12px_40px_rgba(0,0,0,0.35)] sticky top-[68px] z-10">
      <div className="flex items-center justify-between px-4 sm:px-6 py-2 border-b border-[#1c1c1c] bg-[#141414]">
        <div className="text-xs font-semibold uppercase tracking-wide text-muted-foreground/80">Cockpit</div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setIsOpen((prev) => !prev)}
          className="h-8 text-xs"
          aria-expanded={isOpen}
        >
          {isOpen ? "Hide telemetry" : "Show telemetry"}
        </Button>
      </div>

      <div
        className={cn(
          "transition-all duration-300 ease-in-out",
          isOpen ? "max-h-[360px] opacity-100" : "max-h-0 opacity-0 overflow-hidden",
        )}
      >
        <div className="flex flex-col gap-4 p-4 sm:p-5 lg:p-6">
          <div className="flex flex-col lg:flex-row items-stretch gap-4 lg:gap-6">
            {/* Left Module */}
            <div className="flex-1 min-w-[240px] bg-black/25 border border-border/30 rounded-lg px-4 py-3 flex flex-col gap-2">
              <div className="space-y-1">
                <p className="text-[12px] text-primary/80 font-semibold uppercase tracking-wide">Live context</p>
                <h3 className="text-lg font-bold text-foreground leading-tight truncate">{raceName}</h3>
                <p className="text-sm text-muted-foreground/80 truncate">{stageDisplay}</p>
              </div>
              <div className="h-20 w-full bg-black/30 border border-border/30 rounded-md px-3 py-2 flex items-center">
                <svg viewBox="0 0 100 40" className="w-full h-full">
                  <defs>
                    <linearGradient id="elevGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" stopColor="rgb(74,222,128)" stopOpacity="0.9" />
                      <stop offset="100%" stopColor="rgb(74,222,128)" stopOpacity="0.05" />
                    </linearGradient>
                  </defs>
                  <path
                    d={buildPath(profilePoints)}
                    fill="url(#elevGradient)"
                    stroke="rgb(74,222,128)"
                    strokeWidth="1.5"
                    strokeLinejoin="round"
                    strokeLinecap="round"
                    opacity={0.9}
                  />
                  <circle
                    cx={60}
                    cy={yFromPoint(profilePoints, 0.6)}
                    r="2.6"
                    fill="rgb(74,222,128)"
                    className="animate-pulse"
                  />
                </svg>
              </div>
            </div>

            {/* Center Module */}
            <div className="flex flex-1 min-w-[260px] items-stretch gap-3 bg-black/20 border border-border/30 rounded-lg px-3 py-3">
              <TelemetryStat label="Elevation" value={`${elevation} m`} icon={<Mountain className="w-4 h-4" />} dense />
              <div className="w-px bg-border/40 self-stretch" />
              <TelemetryStat label="Distance" value={`${stageLength} km`} icon={<MapPin className="w-4 h-4" />} dense />
            </div>

            {/* Right Module */}
            <div className="flex-1 min-w-[260px] bg-black/25 border border-border/30 rounded-lg p-4 flex flex-col justify-between gap-3">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-xs uppercase tracking-wide text-muted-foreground/80">Points</p>
                  <p className="text-3xl font-mono font-semibold text-foreground tabular-nums">
                    {isAuthenticated ? points.toLocaleString() : "—"}
                  </p>
                  {nextTier && isAuthenticated && (
                    <p className="text-xs text-muted-foreground/80">+{nextTier.min - points} to {nextTier.name}</p>
                  )}
                </div>
                {nextTier && isAuthenticated && (
                  <span className="text-xs font-semibold text-foreground/70 text-right">
                    Watch 15 min for Gold Reward
                  </span>
                )}
              </div>
              {nextTier && isAuthenticated && (
                <div className="space-y-2">
                  <div className="h-2 rounded-full bg-black/50 border border-border/40 overflow-hidden shadow-[0_0_18px_rgba(34,197,94,0.25)]">
                    <div
                      className="h-full bg-gradient-to-r from-primary to-primary/70 rounded-full shadow-[0_0_18px_rgba(34,197,94,0.55)] transition-all duration-500"
                      style={{ width: `${Math.min(progressToNext, 100)}%` }}
                    />
                  </div>
                  <div className="text-xs text-muted-foreground/80 flex items-center justify-between">
                    <span className="font-semibold text-foreground/80">{currentTier.name}</span>
                    <span className="font-semibold text-foreground/80">{nextTier.name}</span>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function createProfileData(length?: number | null, elevation?: number | null): number[] {
  const base = [20, 30, 24, 36, 28, 32, 40, 26, 22, 30, 34, 25]
  if (!length || !elevation) return base
  const scale = Math.max(0.6, Math.min(1.4, elevation / 1500))
  return base.map((v, i) => v + Math.sin(i * 0.7) * 4 * scale)
}

function buildPath(points: number[]): string {
  const maxVal = Math.max(...points)
  const minVal = Math.min(...points)
  const normalized = points.map((p) => ((maxVal - p) / (maxVal - minVal || 1)) * 30 + 5)
  const step = 100 / (points.length - 1)
  let d = `M 0 ${normalized[0]}`
  normalized.forEach((y, idx) => {
    if (idx === 0) return
    d += ` L ${idx * step} ${y}`
  })
  d += ` L 100 40 L 0 40 Z`
  return d
}

function yFromPoint(points: number[], ratio: number): number {
  const maxVal = Math.max(...points)
  const minVal = Math.min(...points)
  const normalized = points.map((p) => ((maxVal - p) / (maxVal - minVal || 1)) * 30 + 5)
  const idx = Math.min(normalized.length - 1, Math.floor(normalized.length * ratio))
  return normalized[idx] ?? 20
}

function TelemetryStat({
  label,
  value,
  icon,
  dense,
}: {
  label: string
  value: string
  icon?: ReactNode
  dense?: boolean
}) {
  return (
    <div className="flex flex-col justify-center px-4 py-2.5 border-r last:border-r-0 border-border/25">
      <div className="flex items-center gap-2 text-[11px] uppercase tracking-wide text-muted-foreground/80">
        {icon}
        <span>{label}</span>
      </div>
      <div
        className={cn(
          "mt-2 font-mono font-bold text-foreground tabular-nums",
          dense ? "text-[22px] lg:text-[26px]" : "text-2xl lg:text-3xl",
        )}
      >
        {value}
      </div>
    </div>
  )
}
