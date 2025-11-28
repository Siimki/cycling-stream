'use client';

import { useAuth } from '@/contexts/AuthContext';

interface WelcomeHeaderProps {
  liveRacesCount?: number;
  upcomingText?: string;
}

export function WelcomeHeader({ liveRacesCount = 0, upcomingText }: WelcomeHeaderProps) {
  const { user } = useAuth();
  const userName = user?.name || user?.email?.split('@')[0] || 'there';

  return (
    <div className="mb-8">
      <h1 className="text-4xl md:text-5xl font-black mb-2">
        Welcome back, {userName}{' '}
        <span className="text-primary">👋</span>
      </h1>
      <p className="text-muted-foreground text-lg">
        {liveRacesCount > 0 ? (
          <>
            <span className="text-primary font-semibold">{liveRacesCount}</span> race{liveRacesCount !== 1 ? 's' : ''} live now
            {upcomingText && <> • {upcomingText}</>}
          </>
        ) : (
          upcomingText || 'Check out the latest races and continue watching'
        )}
      </p>
    </div>
  );
}

