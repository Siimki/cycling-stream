'use client';

import { useEffect } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { createContextLogger } from '@/lib/logger';

const logger = createContextLogger('ErrorBoundary');

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Log error to console for debugging
    logger.error('Application error:', error);
  }, [error]);

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="max-w-md w-full text-center">
        <div className="mb-8">
          <div className="text-6xl mb-4">⚠️</div>
          <h1 className="text-4xl font-bold text-foreground/95 mb-2">Something went wrong</h1>
          <p className="text-muted-foreground mb-4">
            We encountered an unexpected error. Please try again.
          </p>
          {process.env.NODE_ENV === 'development' && error.message && (
            <div className="mt-4 p-4 bg-destructive/10 border border-destructive/20 rounded-lg text-left">
              <p className="text-sm font-semibold text-destructive mb-1">Error Details:</p>
              <p className="text-sm text-destructive/80 font-mono break-all">{error.message}</p>
            </div>
          )}
        </div>
        
        <div className="space-y-4">
          <Button
            onClick={reset}
            className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold"
          >
            Try Again
          </Button>
          <div>
            <Link
              href="/"
              className="text-primary hover:underline"
            >
              Go to Homepage
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

