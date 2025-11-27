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

export async function getRaceStream(id: string): Promise<StreamResponse> {
  return fetchAPI<StreamResponse>(`/races/${id}/stream`);
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

