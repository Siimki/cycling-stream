'use client';

import { ReactNode } from 'react';

interface StatCardProps {
  label: string;
  value: string | number;
  icon?: ReactNode;
  change?: string;
  changeType?: 'positive' | 'negative' | 'neutral';
  highlight?: boolean;
}

export function StatCard({ label, value, icon, change, changeType = 'neutral', highlight = false }: StatCardProps) {
  const changeColorClass = 
    changeType === 'positive' ? 'text-primary' :
    changeType === 'negative' ? 'text-red-500' :
    'text-muted-foreground';

  const borderClass = highlight
    ? 'border-primary glow-green'
    : 'border-border/50 hover:border-border';

  return (
    <div className={`bg-card/80 backdrop-blur-sm border ${borderClass} rounded-xl p-5 transition-all hover:border-border`}>
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm text-muted-foreground font-medium">{label}</span>
        {icon && <span className="text-2xl">{icon}</span>}
      </div>
      <div className="text-3xl font-black text-foreground mb-1">{value}</div>
      {change && (
        <div className={`text-xs font-semibold ${changeColorClass}`}>{change}</div>
      )}
    </div>
  );
}

