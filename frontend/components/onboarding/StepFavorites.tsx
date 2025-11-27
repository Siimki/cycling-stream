'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { X, Plus } from 'lucide-react';
import type { AddFavoriteRequest } from '@/lib/api';

interface StepFavoritesProps {
  favorites: AddFavoriteRequest[];
  onChange: (favorites: AddFavoriteRequest[]) => void;
}

export function StepFavorites({ favorites, onChange }: StepFavoritesProps) {
  const [favoriteType, setFavoriteType] = useState<'rider' | 'team' | 'race' | 'series'>('rider');
  const [favoriteId, setFavoriteId] = useState('');

  const handleAdd = () => {
    if (favoriteId.trim() && favorites.length < 3) {
      onChange([...favorites, { favorite_type: favoriteType, favorite_id: favoriteId.trim() }]);
      setFavoriteId('');
    }
  };

  const handleRemove = (index: number) => {
    onChange(favorites.filter((_, i) => i !== index));
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-foreground mb-2">Pick up to 3 favorite riders/teams</h2>
      <p className="text-muted-foreground mb-8">
        This helps us recommend races and highlight what matters to you. You can skip this step.
      </p>

      <div className="space-y-4">
        <div className="flex gap-2">
          <select
            value={favoriteType}
            onChange={(e) => setFavoriteType(e.target.value as any)}
            className="px-3 py-2 bg-background border border-input rounded-md text-sm"
          >
            <option value="rider">Rider</option>
            <option value="team">Team</option>
            <option value="race">Race</option>
            <option value="series">Series</option>
          </select>
          <Input
            type="text"
            placeholder="Enter name or ID"
            value={favoriteId}
            onChange={(e) => setFavoriteId(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && handleAdd()}
            className="flex-1"
          />
          <Button
            onClick={handleAdd}
            disabled={!favoriteId.trim() || favorites.length >= 3}
            size="icon"
          >
            <Plus className="w-4 h-4" />
          </Button>
        </div>

        {favorites.length > 0 && (
          <div className="space-y-2">
            {favorites.map((fav, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-3 bg-muted/50 rounded-md"
              >
                <div>
                  <span className="text-sm font-medium capitalize">{fav.favorite_type}:</span>{' '}
                  <span className="text-sm">{fav.favorite_id}</span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => handleRemove(index)}
                  className="h-8 w-8"
                >
                  <X className="w-4 h-4" />
                </Button>
              </div>
            ))}
          </div>
        )}

        {favorites.length === 0 && (
          <p className="text-sm text-muted-foreground text-center py-4">
            No favorites added yet. You can add them later in your profile.
          </p>
        )}
      </div>
    </div>
  );
}

