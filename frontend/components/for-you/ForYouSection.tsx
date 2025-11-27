'use client';

import { ReactNode } from 'react';

interface ForYouSectionProps {
  title: string;
  children: ReactNode;
  emptyMessage?: string;
  isEmpty?: boolean;
}

export function ForYouSection({ title, children, emptyMessage, isEmpty }: ForYouSectionProps) {
  return (
    <div className="mb-8">
      <h2 className="text-xl font-bold text-foreground mb-4">{title}</h2>
      {isEmpty ? (
        <p className="text-muted-foreground text-sm">{emptyMessage || 'No items to display'}</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {children}
        </div>
      )}
    </div>
  );
}

