'use client';

import { Button } from '@/components/ui/button';
import { Bike, User, Zap } from 'lucide-react';

interface StepCyclingLevelProps {
  value: 'new' | 'casual' | 'superfan' | null;
  onChange: (value: 'new' | 'casual' | 'superfan') => void;
}

export function StepCyclingLevel({ value, onChange }: StepCyclingLevelProps) {
  return (
    <div>
      <h2 className="text-2xl font-bold text-foreground mb-2">How into cycling are you?</h2>
      <p className="text-muted-foreground mb-8">
        Help us personalize your experience by telling us about your cycling interest level.
      </p>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Button
          variant={value === 'new' ? 'default' : 'outline'}
          onClick={() => onChange('new')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <User className="w-8 h-8" />
          <div className="text-center">
            <div className="font-semibold">New</div>
            <div className="text-xs text-muted-foreground mt-1">Just getting started</div>
          </div>
        </Button>
        <Button
          variant={value === 'casual' ? 'default' : 'outline'}
          onClick={() => onChange('casual')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <Bike className="w-8 h-8" />
          <div className="text-center">
            <div className="font-semibold">Casual Fan</div>
            <div className="text-xs text-muted-foreground mt-1">Watch occasionally</div>
          </div>
        </Button>
        <Button
          variant={value === 'superfan' ? 'default' : 'outline'}
          onClick={() => onChange('superfan')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <Zap className="w-8 h-8" />
          <div className="text-center">
            <div className="font-semibold">Superfan</div>
            <div className="text-xs text-muted-foreground mt-1">Can't get enough</div>
          </div>
        </Button>
      </div>
    </div>
  );
}

