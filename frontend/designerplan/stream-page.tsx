"use client"

import { useState } from "react"
import { VideoPlayer } from "./video-player"
import { LiveChat } from "./live-chat"
import { PointsDisplay } from "./points-display"
import { RaceStats } from "./race-stats"
import { StreamHeader } from "./stream-header"

export function StreamPage() {
  const [points, setPoints] = useState(1250)
  const [watchTime, setWatchTime] = useState(47)

  const addPoints = (amount: number) => {
    setPoints((prev) => prev + amount)
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <StreamHeader points={points} />

      <main className="flex-1 flex flex-col lg:flex-row">
        {/* Main content area */}
        <div className="flex-1 flex flex-col">
          {/* Video + Stats row */}
          <div className="flex flex-col xl:flex-row">
            {/* Video player */}
            <div className="flex-1 relative">
              <VideoPlayer watchTime={watchTime} setWatchTime={setWatchTime} addPoints={addPoints} />
            </div>

            {/* Race stats sidebar - visible on xl screens */}
            <div className="hidden xl:block w-80 border-l border-border">
              <RaceStats />
            </div>
          </div>

          {/* Race stats - visible below video on smaller screens */}
          <div className="xl:hidden border-t border-border">
            <RaceStats compact />
          </div>

          {/* Points display below video */}
          <div className="border-t border-border">
            <PointsDisplay points={points} watchTime={watchTime} addPoints={addPoints} />
          </div>
        </div>

        {/* Chat sidebar */}
        <div className="lg:w-80 xl:w-96 border-t lg:border-t-0 lg:border-l border-border flex flex-col h-[400px] lg:h-auto">
          <LiveChat />
        </div>
      </main>
    </div>
  )
}
