'use client';

import { useAuth } from '@/contexts/AuthContext';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

interface AuthRequiredWrapperProps {
  requiresLogin: boolean;
  children: React.ReactNode;
  raceId: string;
}

export function AuthRequiredWrapper({ requiresLogin, children, raceId }: AuthRequiredWrapperProps) {
  const { isAuthenticated, isLoading } = useAuth();

  // If race doesn't require login, show children
  if (!requiresLogin) {
    return <>{children}</>;
  }

  // If still loading, show loading state
  if (isLoading) {
    return (
      <div className="bg-card aspect-video flex items-center justify-center rounded-lg relative overflow-hidden border border-border">
        <div className="absolute inset-0 bg-gradient-to-br from-background to-card"></div>
        <div className="relative text-center text-foreground z-10 px-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  // If not authenticated and race requires login, show login message
  if (!isAuthenticated) {
    return (
      <div className="bg-card aspect-video flex items-center justify-center rounded-lg relative overflow-hidden border border-border">
        <div className="absolute inset-0 bg-gradient-to-br from-background to-card"></div>
        <div className="relative text-center text-foreground z-10 px-4">
          <div className="text-6xl mb-4">🔒</div>
          <p className="text-2xl font-semibold mb-2">Stream is Only for Registered Users</p>
          <p className="text-muted-foreground mb-4">
            This stream is only available for logged-in users. Please log in to watch.
          </p>
          <Link href="/auth/login">
            <Button className="bg-primary hover:bg-primary/90 text-primary-foreground">
              Log In
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  // User is authenticated, show children
  return <>{children}</>;
}


