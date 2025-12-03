"use client"

import { useState } from "react"
import { MessageSquare } from "lucide-react"
import WatchTrackingProvider from "@/components/WatchTrackingProvider"
import { AuthRequiredWrapper } from "@/components/race/AuthRequiredWrapper"
import { ChatWrapper } from "@/components/race/ChatWrapper"
import { StreamProvider } from "@/components/race/StreamProvider"
import { RaceStats } from "@/components/race/RaceStats"
import { FloatingReactions } from "@/components/race/FloatingReactions"
import Footer from "@/components/layout/Footer"
import { HudStatsProvider } from "@/components/user/HudStatsProvider"
import { cn } from "@/lib/utils"
import { Race, StreamResponse } from "@/lib/api"

interface WatchExperienceLayoutProps {
  raceId: string
  race: Race
  stream: StreamResponse | null
  requiresLogin: boolean
  isLive: boolean
}

export function WatchExperienceLayout({
  raceId,
  race,
  stream,
  requiresLogin,
  isLive,
}: WatchExperienceLayoutProps) {
  const [isChatOpen, setIsChatOpen] = useState(true)

  return (
    <WatchTrackingProvider raceId={raceId}>
      <HudStatsProvider>
        <div className="flex-1 flex flex-col min-h-0">
          <div
            className={cn(
              "flex-1 grid min-h-0 transition-[grid-template-columns] duration-300 ease-out lg:h-[calc(100vh-76px)]",
              isChatOpen ? "lg:grid-cols-[minmax(0,1fr)_clamp(380px,24vw,460px)]" : "lg:grid-cols-[minmax(0,1fr)]",
            )}
          >
            <div className="relative flex flex-col min-h-0 bg-black">
              <div className="flex-1 overflow-y-auto scrollbar-hide pb-12">
                <div className="relative w-full px-3 sm:px-6 lg:px-8 pt-4 pb-8">
                  <div className="relative w-full aspect-video bg-black rounded-none lg:rounded-2xl overflow-hidden border border-border/30 shadow-[0_30px_80px_rgba(0,0,0,0.55)]">
                    <AuthRequiredWrapper requiresLogin={requiresLogin} raceId={raceId}>
                      <StreamProvider raceId={raceId} requiresLogin={requiresLogin} initialStream={stream} />
                    </AuthRequiredWrapper>

                    <FloatingReactions enabled={isLive} />

                    {!isChatOpen && (
                      <button
                        onClick={() => setIsChatOpen(true)}
                        className="hidden lg:flex items-center gap-2 absolute top-4 right-4 rounded-full bg-primary/20 border border-primary/40 px-3 py-2 text-sm font-semibold text-primary-foreground backdrop-blur pointer-events-auto hover:bg-primary/30 transition-colors"
                        aria-label="Open chat"
                      >
                        <MessageSquare className="w-4 h-4" />
                        Open chat
                      </button>
                    )}
                  </div>
                </div>

                <div className="w-full bg-background/90 border-t border-border/20 backdrop-blur-sm">
                  <div className="max-w-[1600px] w-full mx-auto px-4 sm:px-6 lg:px-8">
                    <RaceStats race={race} />
                  </div>
                </div>

                <Footer />
              </div>
            </div>

            {isChatOpen && (
              <div className="hidden lg:flex flex-col min-h-0 border-l border-border/20 bg-background/95 h-full overflow-hidden">
                <ChatWrapper
                  raceId={raceId}
                  requiresLogin={requiresLogin}
                  isLive={isLive}
                  onCollapse={() => setIsChatOpen(false)}
                  className="h-[calc(100vh-76px)] max-h-[calc(100vh-76px)] border-0 overflow-hidden"
                />
              </div>
            )}
          </div>

          <div className="lg:hidden flex-none h-[52vh] border-t border-border/20 bg-background/95 overflow-hidden">
            <ChatWrapper raceId={raceId} requiresLogin={requiresLogin} isLive={isLive} className="h-full border-0" />
          </div>
        </div>
      </HudStatsProvider>
    </WatchTrackingProvider>
  )
}
