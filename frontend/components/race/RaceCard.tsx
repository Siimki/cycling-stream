import { memo } from 'react';
import Link from 'next/link';
import { Race } from '@/lib/api';
import { formatDate } from '@/lib/formatters';

interface RaceCardProps {
  race: Race;
}

function RaceCard({ race }: RaceCardProps) {

  return (
    <Link href={`/races/${race.id}`}>
      <div className="bg-card/80 backdrop-blur-sm border-2 border-primary/40 rounded-lg p-4 sm:p-6 hover:border-primary hover:border-[3px] hover:bg-card/90 hover:shadow-[0_4px_20px_rgba(0,0,0,0.3)] hover:shadow-primary/20 hover:-translate-y-1 transition-all duration-200 cursor-pointer h-full flex flex-col group">
        <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mb-2 group-hover:text-primary transition-colors">{race.name}</h3>
        {race.description && (
          <p className="text-muted-foreground mb-4 line-clamp-2 text-sm sm:text-base flex-grow">
            {race.description}
          </p>
        )}
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 sm:gap-0 text-xs sm:text-sm text-muted-foreground mt-auto">
          <div className="flex flex-wrap items-center gap-1 sm:gap-2">
            {race.location && <span>{race.location}</span>}
            {race.start_date && (
              <>
                <span className="hidden sm:inline">•</span>
                <span>{formatDate(race.start_date, { format: 'short' })}</span>
              </>
            )}
          </div>
          <div className="flex items-center gap-2">
            {race.requires_login && (
              <span className="px-2 py-1 bg-primary/20 text-primary rounded-md text-xs sm:text-sm font-semibold whitespace-nowrap">
                🔒 Login Required
              </span>
            )}
            {race.is_free ? (
              <span className="px-2 py-1 bg-primary/20 text-primary rounded-md text-xs sm:text-sm font-semibold whitespace-nowrap">
                Free
              </span>
            ) : (
              <span className="px-2 py-1 bg-muted/50 text-foreground rounded-md text-xs sm:text-sm font-semibold whitespace-nowrap">
                ${(race.price_cents / 100).toFixed(2)}
              </span>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
}

export default memo(RaceCard);

