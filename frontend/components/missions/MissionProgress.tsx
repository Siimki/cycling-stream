'use client';

interface MissionProgressProps {
  progress: number;
  target: number;
  className?: string;
}

export function MissionProgress({ progress, target, className = '' }: MissionProgressProps) {
  const percentage = Math.min((progress / target) * 100, 100);

  return (
    <div className={`w-full ${className}`}>
      <div className="flex items-center justify-between text-xs sm:text-sm mb-1">
        <span className="text-muted-foreground">
          {progress} / {target}
        </span>
        <span className="text-muted-foreground">{Math.round(percentage)}%</span>
      </div>
      <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
        <div
          className="h-full bg-primary transition-all duration-300 ease-out"
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  );
}

