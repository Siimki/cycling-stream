import * as React from "react"

import { cn } from "@/lib/utils"

// InputProps extends all standard HTML input attributes
export type InputProps = React.InputHTMLAttributes<HTMLInputElement>

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex h-[var(--control-height-md)] min-h-[var(--control-height-md)] w-full rounded-[var(--control-radius)] border border-input bg-background px-[var(--space-3)] py-[var(--space-2)] text-[var(--font-size-sm)] ring-offset-background file:border-0 file:bg-transparent file:text-[var(--font-size-sm)] file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Input.displayName = "Input"

export { Input }
