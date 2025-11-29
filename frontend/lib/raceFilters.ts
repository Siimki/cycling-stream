import { Race } from './api';

export function isRaceReplay(race: Race): boolean {
  // Offline races are replays
  return race.stream_status === 'offline';
}

export function isRaceUpcomingOrLive(race: Race): boolean {
  // Live or upcoming races show in main races page
  return race.stream_status === 'live' || race.stream_status === 'upcoming';
}
