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
  } catch (error) {
    if (suppressAuthErrorLog && typeof error === 'object' && error !== null && 'status' in error && (error as { status?: number }).status === 401) {
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
  role?: 'viewer' | 'mod' | 'vip' | 'subscriber';
  badges?: string[];
  special_emote?: boolean;
}

export interface ChatHistoryResponse {
  messages: ChatMessage[];
  limit: number;
  offset: number;
}

export interface ChatPollOption {
  id: string;
  label: string;
  votes: number;
}

export interface ChatPoll {
  id: string;
  race_id: string;
  question: string;
  options: ChatPollOption[];
  total_votes: number;
  created_at: string;
  closes_at?: string;
  closed: boolean;
}

export interface CreateChatPollRequest {
  question: string;
  options: string[];
  duration_seconds?: number;
}

export async function createChatPoll(raceId: string, payload: CreateChatPollRequest): Promise<ChatPoll> {
  return fetchAuthenticatedAPI<ChatPoll>(`/races/${raceId}/chat/polls`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });
}

export async function closeChatPoll(raceId: string, pollId: string): Promise<ChatPoll> {
  return fetchAuthenticatedAPI<ChatPoll>(`/races/${raceId}/chat/polls/${pollId}/close`, {
    method: 'POST',
  });
}

export async function voteInChatPoll(raceId: string, pollId: string, optionId: string): Promise<ChatPoll> {
  return fetchAuthenticatedAPI<ChatPoll>(`/races/${raceId}/chat/polls/${pollId}/vote`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ option_id: optionId }),
  });
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
  xp_total: number;
  level: number;
  total_watch_minutes: number;
}

export async function getLeaderboard(): Promise<LeaderboardEntry[]> {
  return fetchAPI<LeaderboardEntry[]>('/leaderboard');
}

// User Preferences Types and API
export interface UIPreferences {
  chat_animations: boolean;
  reduced_motion: boolean;
  button_pulse: boolean;
  poll_animations: boolean;
}

export interface AudioPreferences {
  button_clicks: boolean;
  notification_sounds: boolean;
  mention_pings: boolean;
  master_volume: number;
}

export interface UserPreferences {
  id: string;
  user_id: string;
  data_mode: 'casual' | 'standard' | 'pro';
  preferred_units: 'metric' | 'imperial';
  theme: 'light' | 'dark' | 'auto';
  accent_color?: string;
  device_type?: 'tv' | 'desktop' | 'mobile' | 'tablet';
  notification_preferences: Record<string, boolean>;
  ui_preferences: UIPreferences;
  audio_preferences: AudioPreferences;
  onboarding_completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpdateUIPreferencesRequest {
  chat_animations?: boolean;
  reduced_motion?: boolean;
  button_pulse?: boolean;
  poll_animations?: boolean;
}

export interface UpdateAudioPreferencesRequest {
  button_clicks?: boolean;
  notification_sounds?: boolean;
  mention_pings?: boolean;
  master_volume?: number;
}

export interface UpdatePreferencesRequest {
  data_mode?: 'casual' | 'standard' | 'pro';
  preferred_units?: 'metric' | 'imperial';
  theme?: 'light' | 'dark' | 'auto';
  accent_color?: string;
  device_type?: 'tv' | 'desktop' | 'mobile' | 'tablet';
  notification_preferences?: Record<string, boolean>;
  onboarding_completed?: boolean;
  ui_preferences?: UpdateUIPreferencesRequest;
  audio_preferences?: UpdateAudioPreferencesRequest;
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
type FetchError = Error & { status?: number };

async function fetchAuthenticatedAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
  if (!token) {
    const authError = new Error('Authentication required') as FetchError;
    authError.status = 401;
    throw authError;
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
  } catch (error) {
    const errorStatus = (error as FetchError)?.status;
    // Handle 404 (preferences not found) by returning defaults
    // Also handle 500 errors which might occur for unauthenticated users
    if (errorStatus === 404 || errorStatus === 500) {
      return {
        id: '',
        user_id: '',
        data_mode: 'standard',
        preferred_units: 'metric',
        theme: 'auto',
        notification_preferences: {},
        ui_preferences: {
          chat_animations: true,
          reduced_motion: false,
          button_pulse: true,
          poll_animations: true,
        },
        audio_preferences: {
          button_clicks: true,
          notification_sounds: true,
          mention_pings: true,
          master_volume: 0.15,
        },
        onboarding_completed: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
    }
    // Re-throw 401/500 errors so providers can handle them (user not authenticated)
    // Re-throw other errors as well
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

// Mission types and interfaces
export type MissionType = 'watch_time' | 'chat_message' | 'watch_race' | 'follow_series' | 'streak' | 'predict_winner';

export interface Mission {
  id: string;
  mission_type: MissionType;
  title: string;
  description?: string;
  points_reward: number;
  target_value: number;
  valid_from: string;
  valid_until?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserMission {
  id: string;
  user_id: string;
  mission_id: string;
  progress: number;
  completed_at?: string;
  claimed_at?: string;
  created_at: string;
  updated_at: string;
}

// UserMissionWithDetails has all fields flattened (Go embedded structs serialize as flat JSON)
// Note: When Go embeds structs, duplicate fields from Mission overwrite UserMission fields
// So 'id', 'created_at', 'updated_at' come from Mission, not UserMission
export interface UserMissionWithDetails {
  // From UserMission (but id/created_at/updated_at are overwritten by Mission)
  id: string; // This is actually Mission.id, not UserMission.id
  user_id: string;
  mission_id: string;
  progress: number;
  completed_at?: string;
  claimed_at?: string;
  created_at: string; // This is actually Mission.created_at
  updated_at: string; // This is actually Mission.updated_at
  
  // From Mission
  mission_type: MissionType;
  title: string;
  description?: string;
  points_reward: number;
  target_value: number;
  valid_from: string;
  valid_until?: string;
  is_active: boolean;
}

export async function getUserMissions(): Promise<UserMissionWithDetails[]> {
  return fetchAuthenticatedAPI<UserMissionWithDetails[]>('/users/me/missions');
}

export async function getActiveMissions(): Promise<Mission[]> {
  return fetchAPI<Mission[]>('/missions/active');
}

export async function claimMissionReward(missionId: string): Promise<{ message: string }> {
  return fetchAuthenticatedAPI<{ message: string }>(`/users/me/missions/${missionId}/claim`, {
    method: 'POST',
  });
}

// XP and Level types and API
export interface XPProgress {
  user_id: string;
  xp_total: number;
  level: number;
  xp_for_current_level_start: number;
  xp_for_next_level: number;
  xp_to_next_level: number;
  progress_in_current_level: number;
}

export async function getUserXP(): Promise<XPProgress> {
  return fetchAuthenticatedAPI<XPProgress>('/users/me/xp');
}

export interface UserAchievement {
  id: string;
  slug: string;
  title: string;
  description?: string;
  icon?: string;
  points: number;
  unlocked_at: string;
  metadata?: Record<string, unknown>;
}

interface UserAchievementsResponse {
  achievements: UserAchievement[];
}

export async function getUserAchievements(): Promise<UserAchievement[]> {
  if (typeof window === 'undefined') {
    return [];
  }
  const token = localStorage.getItem('auth_token');
  if (!token) {
    return [];
  }

  const response = await fetch(`${API_URL}/users/me/achievements`, {
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });

  if (response.status === 404) {
    return [];
  }

  if (response.status === 401 || response.status === 500) {
    // 401 = unauthorized, 500 = likely invalid token or server error
    // Return empty array instead of throwing to allow graceful degradation
    const authError = new Error('Authentication required') as FetchError;
    authError.status = response.status;
    throw authError;
  }

  if (!response.ok) {
    throw new Error(`Failed to load achievements (${response.status})`);
  }

  const data: UserAchievementsResponse = await response.json();
  return data.achievements ?? [];
}

// Weekly stats types and API
export interface WeeklyGoalProgress {
  user_id: string;
  week_number: string;
  watch_minutes: number;
  chat_messages: number;
  weekly_goal_completed: boolean;
  current_streak_weeks: number;
  best_streak_weeks: number;
  can_claim_reward: boolean;
  reward_xp: number;
  reward_points: number;
}

export async function getUserWeekly(): Promise<WeeklyGoalProgress> {
  return fetchAuthenticatedAPI<WeeklyGoalProgress>('/users/me/weekly');
}

export async function claimWeeklyReward(): Promise<{ message: string }> {
  return fetchAuthenticatedAPI<{ message: string }>('/users/me/weekly/claim', {
    method: 'POST',
  });
}

