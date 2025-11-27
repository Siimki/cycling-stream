'use client';

import ErrorMessage from '@/components/ErrorMessage';
import Link from 'next/link';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="min-h-screen bg-background p-8">
      <ErrorMessage error={error} onRetry={reset} variant="full" />
      <div className="text-center mt-6">
        <Link href="/" className="text-primary hover:underline">
          ← Back to races
        </Link>
      </div>
    </div>
  );
}
