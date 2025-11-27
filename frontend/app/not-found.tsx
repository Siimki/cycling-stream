import Link from 'next/link';
import { Button } from '@/components/ui/button';

export default function NotFound() {
  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="max-w-md w-full text-center">
        <div className="mb-8">
          <h1 className="text-9xl font-bold text-muted-foreground/20">404</h1>
          <div className="mt-4">
            <h2 className="text-3xl font-bold text-foreground/95 mb-2">Page Not Found</h2>
            <p className="text-muted-foreground">
              The page you&apos;re looking for doesn&apos;t exist or has been moved.
            </p>
          </div>
        </div>
        
        <div className="space-y-4">
          <Link href="/">
            <Button className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold">
              Go to Homepage
            </Button>
          </Link>
          <div>
            <Link
              href="/"
              className="text-primary hover:underline"
            >
              Browse Races
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

