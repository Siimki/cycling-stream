interface SkeletonLoaderProps {
  className?: string;
}

export function SkeletonBox({ className = '' }: SkeletonLoaderProps) {
  return (
    <div
      className={`animate-pulse bg-gray-200 rounded ${className}`}
      aria-hidden="true"
    ></div>
  );
}

export function RaceCardSkeleton() {
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <SkeletonBox className="h-6 w-3/4 mb-4" />
      <SkeletonBox className="h-4 w-full mb-2" />
      <SkeletonBox className="h-4 w-2/3 mb-4" />
      <div className="flex items-center justify-between">
        <SkeletonBox className="h-4 w-1/3" />
        <SkeletonBox className="h-6 w-16 rounded" />
      </div>
    </div>
  );
}

export function RaceDetailSkeleton() {
  return (
    <div className="bg-white rounded-lg shadow-md p-8">
      <SkeletonBox className="h-8 w-2/3 mb-4" />
      <SkeletonBox className="h-4 w-full mb-2" />
      <SkeletonBox className="h-4 w-5/6 mb-6" />
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <SkeletonBox className="h-20" />
        <SkeletonBox className="h-20" />
        <SkeletonBox className="h-20" />
        <SkeletonBox className="h-20" />
      </div>
      <SkeletonBox className="h-12 w-32 rounded-lg" />
    </div>
  );
}

