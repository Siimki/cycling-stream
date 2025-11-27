'use client';

import { Button } from '@/components/ui/button';
import { Bell, BellOff } from 'lucide-react';

interface StepNotificationsProps {
  value: {
    challenges: boolean;
    points: boolean;
    races: boolean;
  };
  onChange: (value: { challenges: boolean; points: boolean; races: boolean }) => void;
}

export function StepNotifications({ value, onChange }: StepNotificationsProps) {
  const toggle = (key: keyof typeof value) => {
    onChange({ ...value, [key]: !value[key] });
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-foreground mb-2">Opt in to challenges and points?</h2>
      <p className="text-muted-foreground mb-8">
        Choose what notifications you'd like to receive. You can change these anytime.
      </p>
      <div className="space-y-4">
        <Button
          variant={value.challenges ? 'default' : 'outline'}
          onClick={() => toggle('challenges')}
          className="w-full h-auto py-4 flex items-center justify-between"
        >
          <div className="flex items-center gap-3">
            {value.challenges ? <Bell className="w-5 h-5" /> : <BellOff className="w-5 h-5" />}
            <div className="text-left">
              <div className="font-semibold">Challenges & Missions</div>
              <div className="text-sm text-muted-foreground">Get notified about new challenges</div>
            </div>
          </div>
        </Button>
        <Button
          variant={value.points ? 'default' : 'outline'}
          onClick={() => toggle('points')}
          className="w-full h-auto py-4 flex items-center justify-between"
        >
          <div className="flex items-center gap-3">
            {value.points ? <Bell className="w-5 h-5" /> : <BellOff className="w-5 h-5" />}
            <div className="text-left">
              <div className="font-semibold">Points & Rewards</div>
              <div className="text-sm text-muted-foreground">Get notified about points milestones</div>
            </div>
          </div>
        </Button>
        <Button
          variant={value.races ? 'default' : 'outline'}
          onClick={() => toggle('races')}
          className="w-full h-auto py-4 flex items-center justify-between"
        >
          <div className="flex items-center gap-3">
            {value.races ? <Bell className="w-5 h-5" /> : <BellOff className="w-5 h-5" />}
            <div className="text-left">
              <div className="font-semibold">Race Updates</div>
              <div className="text-sm text-muted-foreground">Get notified about upcoming races</div>
            </div>
          </div>
        </Button>
      </div>
    </div>
  );
}

