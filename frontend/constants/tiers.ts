/**
 * Points tier definitions
 */

export interface PointsTier {
  name: string;
  min: number;
  color: string;
}

export const POINTS_TIERS: PointsTier[] = [
  { name: 'Bronze', min: 0, color: 'text-amber-600' },
  { name: 'Silver', min: 500, color: 'text-slate-400' },
  { name: 'Gold', min: 1500, color: 'text-yellow-500' },
  { name: 'Platinum', min: 3500, color: 'text-cyan-400' },
  { name: 'Diamond', min: 7500, color: 'text-violet-400' },
];

