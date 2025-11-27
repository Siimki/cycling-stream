"use client"

import { Zap, Search, Bell, User, ChevronDown } from "lucide-react"
import { Button } from "@/components/ui/button"

interface StreamHeaderProps {
  points: number
}

export function StreamHeader({ points }: StreamHeaderProps) {
  return (
    <header className="sticky top-0 z-50 bg-background/80 backdrop-blur-xl border-b border-border/50">
      <div className="px-4 lg:px-6">
        <div className="flex items-center justify-between h-14">
          {/* Logo */}
          <div className="flex items-center gap-8">
            <div className="flex items-center gap-2">
              <div className="w-7 h-7 rounded-md bg-primary flex items-center justify-center">
                <Zap className="w-4 h-4 text-primary-foreground" />
              </div>
              <span className="text-base font-semibold tracking-tight">
                Peloton<span className="text-primary">Live</span>
              </span>
            </div>

            {/* Navigation */}
            <nav className="hidden md:flex items-center gap-1">
              <Button variant="ghost" size="sm" className="text-foreground text-xs font-medium h-8 px-3">
                Live
                <span className="ml-1.5 w-1.5 h-1.5 rounded-full bg-red-500 animate-live-pulse" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="text-muted-foreground hover:text-foreground text-xs font-medium h-8 px-3"
              >
                Schedule
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="text-muted-foreground hover:text-foreground text-xs font-medium h-8 px-3"
              >
                Replays
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="text-muted-foreground hover:text-foreground text-xs font-medium h-8 px-3"
              >
                Teams
              </Button>
            </nav>
          </div>

          {/* Right Side */}
          <div className="flex items-center gap-2">
            {/* Points Display */}
            <div className="hidden sm:flex items-center gap-1.5 bg-primary/10 hover:bg-primary/15 transition-colors cursor-pointer pl-2 pr-3 py-1.5 rounded-full">
              <div className="w-5 h-5 rounded-full bg-primary/20 flex items-center justify-center">
                <Zap className="w-3 h-3 text-primary" />
              </div>
              <span className="text-xs font-semibold text-primary tabular-nums">{points.toLocaleString()}</span>
            </div>

            <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-foreground w-8 h-8">
              <Search className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="text-muted-foreground hover:text-foreground w-8 h-8 relative"
            >
              <Bell className="w-4 h-4" />
              <span className="absolute top-1 right-1 w-2 h-2 bg-primary rounded-full" />
            </Button>

            <div className="flex items-center gap-1 ml-1 cursor-pointer hover:opacity-80 transition-opacity">
              <div className="w-7 h-7 rounded-full bg-gradient-to-br from-primary/60 to-primary flex items-center justify-center">
                <User className="w-3.5 h-3.5 text-primary-foreground" />
              </div>
              <ChevronDown className="w-3 h-3 text-muted-foreground" />
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}
