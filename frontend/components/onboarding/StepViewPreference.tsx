'use client';

import { Button } from '@/components/ui/button';
import { Eye, BarChart3 } from 'lucide-react';

interface StepViewPreferenceProps {
  value: 'clean' | 'data-rich' | null;
  onChange: (value: 'clean' | 'data-rich') => void;
}

export function StepViewPreference({ value, onChange }: StepViewPreferenceProps) {
  return (
    <div>
      <h2 className="text-2xl font-bold text-foreground mb-2">Which do you prefer?</h2>
      <p className="text-muted-foreground mb-8">
        Choose your preferred viewing experience. You can always change this later.
      </p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Button
          variant={value === 'clean' ? 'default' : 'outline'}
          onClick={() => onChange('clean')}
          className="h-auto py-8 flex flex-col items-center gap-4 px-4 w-full whitespace-normal"
          style={{ wordBreak: 'break-word', overflowWrap: 'break-word' }}
        >
          <Eye className="w-10 h-10 shrink-0" />
          <div className="text-center w-full" style={{ minWidth: 0, maxWidth: '100%', wordBreak: 'break-word', overflowWrap: 'break-word' }}>
            <div className="font-semibold text-lg mb-2 whitespace-normal">Clean Screen</div>
            <div className="text-sm text-muted-foreground whitespace-normal break-words" style={{ wordBreak: 'break-word', overflowWrap: 'break-word' }}>
              Minimal on-screen data. Perfect for casual watching.
            </div>
          </div>
        </Button>
        <Button
          variant={value === 'data-rich' ? 'default' : 'outline'}
          onClick={() => onChange('data-rich')}
          className="h-auto py-8 flex flex-col items-center gap-4 px-4 w-full whitespace-normal"
          style={{ wordBreak: 'break-word', overflowWrap: 'break-word' }}
        >
          <BarChart3 className="w-10 h-10 shrink-0" />
          <div className="text-center w-full" style={{ minWidth: 0, maxWidth: '100%', wordBreak: 'break-word', overflowWrap: 'break-word' }}>
            <div className="font-semibold text-lg mb-2 whitespace-normal">Data-Rich Screen</div>
            <div className="text-sm text-muted-foreground whitespace-normal break-words" style={{ wordBreak: 'break-word', overflowWrap: 'break-word' }}>
              All the stats, maps, and analytics you want.
            </div>
          </div>
        </Button>
      </div>
    </div>
  );
}

