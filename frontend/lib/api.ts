import { API_URL } from './config';
import { APIErrorHandler } from './error-handler';

export interface Race {
  id: string;
  name: string;
  description?: string;
  start_date?: string;
  end_date?: string;
  location?: string;
  category?: string;
  is_free: boolean;
  price_cents: number;
  requires_login?: boolean;
  stage_name?: string;
  stage_type?: string;
  elevation_meters?: number;
  estimated_finish_time?: string;
  stage_length_km?: number;
  created_at: string;
  updated_at: string;
}

export interface StreamResponse {
  status: string;
  stream_type?: string;
  source_id?: string;
  origin_url?: string;
  cdn_url?: string;
}

/**
 * Fetches data from the API with standardized error handling
 */
async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  try {
    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    return await APIErrorHandler.handleResponse<T>(response);
  } catch (error) {
    // Re-throw APIError as-is, wrap other errors
    if (error && typeof error === 'object' && 'status' in error) {
      throw error;
    }
    throw new Error(APIErrorHandler.getErrorMessage(error));
  }
}

export async function getRaces(): Promise<Race[]> {
  return fetchAPI<Race[]>('/races');
}

export async function getRace(id: string): Promise<Race> {
  return fetchAPI<Race>(`/races/${id}`);
}

export async function getRaceStream(id: string, suppressAuthErrorLog = false): Promise<StreamResponse> {
  try {
    return await fetchAPI<StreamResponse>(`/races/${id}/stream`);
  } catch (error: any) {
    // If this is an expected 401 error and we're suppressing logs, re-throw without logging
    if (suppressAuthErrorLog && error?.status === 401) {
      throw error;
    }
    throw error;
  }
}

export interface ChatMessage {
  id: string;
  race_id: string;
  user_id?: string;
  username: string;
  message: string;
  created_at: string;
}

export interface ChatHistoryResponse {
  messages: ChatMessage[];
  limit: number;
  offset: number;
}

export interface ChatStatsResponse {
  total_messages: number;
  concurrent_connections: number;
}

export async function getChatHistory(raceId: string, limit = 50, offset = 0): Promise<ChatHistoryResponse> {
  return fetchAPI<ChatHistoryResponse>(`/races/${raceId}/chat/history?limit=${limit}&offset=${offset}`);
}

export async function getChatStats(raceId: string): Promise<ChatStatsResponse> {
  return fetchAPI<ChatStatsResponse>(`/races/${raceId}/chat/stats`);
}

export interface PublicUser {
  id: string;
  name?: string;
  bio: string;
  points: number;
  total_watch_minutes: number;
  created_at: string;
}

export async function getPublicUser(userId: string): Promise<PublicUser> {
  return fetchAPI<PublicUser>(`/profiles/${userId}`);
}

export interface LeaderboardEntry {
  id: string;
  name?: string;
  points: number;
  total_watch_minutes: number;
}

export async function getLeaderboard(): Promise<LeaderboardEntry[]> {
  return fetchAPI<LeaderboardEntry[]>('/leaderboard');
}

// User Preferences Types and API
export interface UserPreferences {
  id: string;
  user_id: string;
  data_mode: 'casual' | 'standard' | 'pro';
  preferred_units: 'metric' | 'imperial';
  theme: 'light' | 'dark' | 'auto';
  accent_color?: string;
  device_type?: 'tv' | 'desktop' | 'mobile' | 'tablet';
  notification_preferences: Record<string, boolean>;
  onboarding_completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpdatePreferencesRequest {
  data_mode?: 'casual' | 'standard' | 'pro';
  preferred_units?: 'metric' | 'imperial';
  theme?: 'light' | 'dark' | 'auto';
  accent_color?: string;
  device_type?: 'tv' | 'desktop' | 'mobile' | 'tablet';
  notification_preferences?: Record<string, boolean>;
  onboarding_completed?: boolean;
}

export interface UserFavorite {
  id: string;
  user_id: string;
  favorite_type: 'rider' | 'team' | 'race' | 'series';
  favorite_id: string;
  created_at: string;
}

export interface AddFavoriteRequest {
  favorite_type: 'rider' | 'team' | 'race' | 'series';
  favorite_id: string;
}

export interface WatchHistoryEntry {
  user_id: string;
  race_id: string;
  race_name: string;
  race_category?: string;
  race_start_date?: string;
  session_count: number;
  total_seconds: number;
  total_minutes: number;
  first_watched: string;
  last_watched: string;
  likely_completed: boolean;
}

export interface WatchHistoryResponse {
  entries: WatchHistoryEntry[];
  total: number;
  limit: number;
  offset: number;
}

/**
 * Fetches data from authenticated API endpoints
 */
async function fetchAuthenticatedAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  if (!token) {
    const error: any = new Error('Authentication required');
    error.status = 401;
    throw error;
  }

  return fetchAPI<T>(endpoint, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      ...options?.headers,
    },
  });
}

export async function getUserPreferences(): Promise<UserPreferences> {
  try {
    return await fetchAuthenticatedAPI<UserPreferences>('/users/me/preferences');
  } catch (error: any) {
    // If 404, return defaults (preferences don't exist yet)
    if (error?.status === 404) {
      return {
        id: '',
        user_id: '',
        data_mode: 'standard',
        preferred_units: 'metric',
        theme: 'auto',
        notification_preferences: {},
        onboarding_completed: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
    }
    throw error;
  }
}

export async function updateUserPreferences(prefs: UpdatePreferencesRequest): Promise<UserPreferences> {
  return fetchAuthenticatedAPI<UserPreferences>('/users/me/preferences', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(prefs),
  });
}

export async function completeOnboarding(): Promise<UserPreferences> {
  return fetchAuthenticatedAPI<UserPreferences>('/users/me/onboarding/complete', {
    method: 'POST',
  });
}

export async function getUserFavorites(type?: 'rider' | 'team' | 'race' | 'series'): Promise<UserFavorite[]> {
  const endpoint = type ? `/users/me/favorites?type=${type}` : '/users/me/favorites';
  return fetchAuthenticatedAPI<UserFavorite[]>(endpoint);
}

export async function addFavorite(favorite: AddFavoriteRequest): Promise<UserFavorite> {
  return fetchAuthenticatedAPI<UserFavorite>('/users/me/favorites', {
    method: 'POST',
    body: JSON.stringify(favorite),
  });
}

export async function removeFavorite(type: 'rider' | 'team' | 'race' | 'series', id: string): Promise<void> {
  return fetchAuthenticatedAPI<void>(`/users/me/favorites/${type}/${id}`, {
    method: 'DELETE',
  });
}

export async function getWatchHistory(limit = 20, offset = 0): Promise<WatchHistoryResponse> {
  return fetchAuthenticatedAPI<WatchHistoryResponse>(`/users/me/watch-history?limit=${limit}&offset=${offset}`);
}

// Recommendations Types and API
export interface RecommendationsResponse {
  continue_watching: Race[];
  upcoming: Race[];
  replays: Race[];
}

export async function getRecommendations(): Promise<RecommendationsResponse> {
  return fetchAuthenticatedAPI<RecommendationsResponse>('/users/me/recommendations');
}

export async function getContinueWatching(): Promise<Race[]> {
  return fetchAuthenticatedAPI<Race[]>('/users/me/recommendations/continue-watching');
}

export async function getUpcomingRaces(): Promise<Race[]> {
  return fetchAuthenticatedAPI<Race[]>('/users/me/recommendations/upcoming');
}

export async function getRecommendedReplays(): Promise<Race[]> {
  return fetchAuthenticatedAPI<Race[]>('/users/me/recommendations/replays');
}

