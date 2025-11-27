/**
 * API endpoint paths
 */

export const API_ENDPOINTS = {
  // Auth
  AUTH_REGISTER: '/auth/register',
  AUTH_LOGIN: '/auth/login',
  
  // Users
  USERS_ME: '/users/me',
  USERS_ME_PASSWORD: '/users/me/password',
  USERS_ME_POINTS_BONUS: '/users/me/points/bonus',
  USERS_PROFILE: (userId: string) => `/profiles/${userId}`,
  
  // Races
  RACES: '/races',
  RACE_BY_ID: (id: string) => `/races/${id}`,
  RACE_STREAM: (id: string) => `/races/${id}/stream`,
  
  // Chat
  CHAT_HISTORY: (raceId: string) => `/races/${raceId}/chat/history`,
  CHAT_STATS: (raceId: string) => `/races/${raceId}/chat/stats`,
  CHAT_WS: (raceId: string) => `/races/${raceId}/chat/ws`,
  
  // Watch sessions
  WATCH_SESSIONS_START: '/users/watch/sessions/start',
  WATCH_SESSIONS_END: '/users/watch/sessions/end',
  WATCH_SESSIONS_STATS: (raceId: string) => `/users/watch/sessions/stats/${raceId}`,
  
  // Admin
  ADMIN_ANALYTICS_RACES: '/admin/analytics/races',
  ADMIN_ANALYTICS_WATCH_TIME: '/admin/analytics/watch-time',
  ADMIN_ANALYTICS_REVENUE: '/admin/analytics/revenue',
} as const;

