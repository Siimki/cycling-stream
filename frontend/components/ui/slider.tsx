"use client"

import * as React from "react"
import { cn } from "@/lib/utils"

interface SliderProps {
  value: number[]
  onValueChange: (value: number[]) => void
  max?: number
  step?: number
  className?: string
}

export function Slider({ value, onValueChange, max = 100, step = 1, className }: SliderProps) {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onValueChange([parseFloat(e.target.value)])
  }

  return (
    <div className={cn("relative flex w-full touch-none select-none items-center", className)}>
      <input
        type="range"
        min={0}
        max={max}
        step={step}
        value={value[0]}
        onChange={handleChange}
        className="h-1.5 w-full cursor-pointer appearance-none rounded-full bg-secondary accent-primary"
      />
    </div>
  )
}

