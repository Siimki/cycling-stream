import { RaceDetailSkeleton } from '@/components/SkeletonLoader';

export default function Loading() {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="h-6 w-24 bg-gray-200 rounded animate-pulse"></div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <RaceDetailSkeleton />
      </main>
    </div>
  );
}

