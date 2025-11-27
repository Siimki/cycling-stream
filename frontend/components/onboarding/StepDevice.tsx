'use client';

import { Button } from '@/components/ui/button';
import { Tv, Monitor, Smartphone, Tablet } from 'lucide-react';

interface StepDeviceProps {
  value: ('tv' | 'laptop' | 'phone' | 'tablet')[];
  onChange: (value: ('tv' | 'laptop' | 'phone' | 'tablet')[]) => void;
}

export function StepDevice({ value, onChange }: StepDeviceProps) {
  // Ensure value is always an array
  const devices = Array.isArray(value) ? value : [];
  
  const toggleDevice = (device: 'tv' | 'laptop' | 'phone' | 'tablet') => {
    if (devices.includes(device)) {
      onChange(devices.filter(d => d !== device));
    } else {
      onChange([...devices, device]);
    }
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-foreground mb-2">What do you watch on?</h2>
      <p className="text-muted-foreground mb-8">
        Select all devices you use. We'll optimize the interface for your preferred devices.
      </p>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Button
          variant={devices.includes('tv') ? 'default' : 'outline'}
          onClick={() => toggleDevice('tv')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <Tv className="w-8 h-8" />
          <div className="font-semibold">TV</div>
        </Button>
        <Button
          variant={devices.includes('laptop') ? 'default' : 'outline'}
          onClick={() => toggleDevice('laptop')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <Monitor className="w-8 h-8" />
          <div className="font-semibold">Laptop</div>
        </Button>
        <Button
          variant={devices.includes('phone') ? 'default' : 'outline'}
          onClick={() => toggleDevice('phone')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <Smartphone className="w-8 h-8" />
          <div className="font-semibold">Phone</div>
        </Button>
        <Button
          variant={devices.includes('tablet') ? 'default' : 'outline'}
          onClick={() => toggleDevice('tablet')}
          className="h-auto py-6 flex flex-col items-center gap-3"
        >
          <Tablet className="w-8 h-8" />
          <div className="font-semibold">Tablet</div>
        </Button>
      </div>
      {devices.length > 0 && (
        <p className="mt-4 text-sm text-muted-foreground text-center">
          Selected: {devices.join(', ')}
        </p>
      )}
    </div>
  );
}

