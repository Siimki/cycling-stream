"use client"

import { Mountain, MapPin, Timer, Users, ChevronUp, Flag } from "lucide-react"

interface RaceStatsProps {
  compact?: boolean
}

const riders = [
  { pos: 1, name: "T. Pogačar", team: "UAE", time: "—", gap: "—" },
  { pos: 2, name: "J. Vingegaard", team: "VIS", time: "+0:47", gap: "+0:47" },
  { pos: 3, name: "R. Evenepoel", team: "SOQ", time: "+1:32", gap: "+0:45" },
  { pos: 4, name: "P. Roglic", team: "RBH", time: "+2:15", gap: "+0:43" },
  { pos: 5, name: "C. Rodriguez", team: "LID", time: "+3:08", gap: "+0:53" },
]

export function RaceStats({ compact = false }: RaceStatsProps) {
  if (compact) {
    return (
      <div className="px-4 py-3 flex items-center justify-between gap-6 overflow-x-auto">
        <div className="flex items-center gap-6 text-xs">
          <div className="flex items-center gap-1.5">
            <Flag className="w-3.5 h-3.5 text-primary" />
            <span className="text-muted-foreground">Stage 17</span>
            <span className="font-medium">Col du Galibier</span>
          </div>
          <div className="flex items-center gap-1.5">
            <MapPin className="w-3.5 h-3.5 text-muted-foreground" />
            <span className="tabular-nums font-medium">124.3 / 166 km</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Mountain className="w-3.5 h-3.5 text-amber-500" />
            <span className="text-muted-foreground">4,800m</span>
          </div>
        </div>
        <div className="flex items-center gap-3 text-[11px]">
          <span className="text-primary font-medium">Pogačar</span>
          <span className="text-muted-foreground">leads by</span>
          <span className="font-medium text-foreground">+0:47</span>
        </div>
      </div>
    )
  }

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-4 py-3 border-b border-border/50">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-[10px] text-primary uppercase tracking-wide font-medium">Tour de France 2025</p>
            <h3 className="text-sm font-semibold">Stage 17 — Mountain</h3>
          </div>
          <div className="text-right">
            <p className="text-sm font-bold tabular-nums">166 km</p>
            <p className="text-[10px] text-muted-foreground">Distance</p>
          </div>
        </div>
      </div>

      {/* Stage info */}
      <div className="px-4 py-3 grid grid-cols-3 gap-2 border-b border-border/50">
        <div>
          <div className="flex items-center gap-1 text-muted-foreground mb-0.5">
            <Mountain className="w-3 h-3" />
            <span className="text-[10px]">Elevation</span>
          </div>
          <p className="text-xs font-semibold">4,800m</p>
        </div>
        <div>
          <div className="flex items-center gap-1 text-muted-foreground mb-0.5">
            <ChevronUp className="w-3 h-3" />
            <span className="text-[10px]">Max Grade</span>
          </div>
          <p className="text-xs font-semibold">12.4%</p>
        </div>
        <div>
          <div className="flex items-center gap-1 text-muted-foreground mb-0.5">
            <Timer className="w-3 h-3" />
            <span className="text-[10px]">Est. Finish</span>
          </div>
          <p className="text-xs font-semibold">~17:45</p>
        </div>
      </div>

      {/* Progress */}
      <div className="px-4 py-3 border-b border-border/50">
        <div className="flex items-center justify-between text-[10px] text-muted-foreground mb-1.5">
          <span>Race Progress</span>
          <span className="font-medium text-foreground">124.3 km</span>
        </div>
        <div className="h-1.5 bg-muted rounded-full overflow-hidden">
          <div className="h-full w-[75%] bg-primary rounded-full" />
        </div>
        <div className="flex justify-between mt-1 text-[9px] text-muted-foreground">
          <span>Start</span>
          <span>Galibier</span>
        </div>
      </div>

      {/* Live standings */}
      <div className="flex-1 flex flex-col min-h-0">
        <div className="px-4 py-2 flex items-center justify-between border-b border-border/50">
          <div className="flex items-center gap-1.5">
            <Users className="w-3 h-3 text-muted-foreground" />
            <span className="text-[10px] font-medium uppercase tracking-wide">Live GC</span>
          </div>
          <span className="text-[9px] text-muted-foreground">Virtual Standings</span>
        </div>

        <div className="flex-1 overflow-y-auto">
          {riders.map((rider, i) => (
            <div
              key={rider.pos}
              className={`px-4 py-2 flex items-center justify-between text-xs border-b border-border/30 last:border-0 ${
                i === 0 ? "bg-primary/5" : ""
              }`}
            >
              <div className="flex items-center gap-2.5">
                <span className={`w-4 text-center font-medium ${i === 0 ? "text-primary" : "text-muted-foreground"}`}>
                  {rider.pos}
                </span>
                <div>
                  <p className={`font-medium ${i === 0 ? "text-primary" : ""}`}>{rider.name}</p>
                  <p className="text-[9px] text-muted-foreground">{rider.team}</p>
                </div>
              </div>
              <span className={`tabular-nums ${i === 0 ? "text-primary font-medium" : "text-muted-foreground"}`}>
                {rider.time}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
