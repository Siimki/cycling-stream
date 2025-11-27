# CyclingStream API Documentation

This document provides comprehensive API documentation for the CyclingStream platform.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://api.cyclingstream.com` (example)

## Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### Getting a Token

1. Register a new account: `POST /auth/register`
2. Login: `POST /auth/login`

Both endpoints return a `token` in the response.

## Rate Limiting

- **Public endpoints**: Lenient rate limiting
- **Auth endpoints**: Strict rate limiting (prevents brute force)
- **User endpoints**: Standard rate limiting
- **Admin endpoints**: Standard rate limiting

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message here"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

---

## Public Endpoints

### Health Check

**GET** `/health`

Check if the API is running.

**Response:**
```json
{
  "status": "ok",
  "database": "connected"
}
```

---

### Get All Races

**GET** `/races`

Get a list of all races.

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "Tour de France",
    "description": "Annual cycling race",
    "start_date": "2024-07-01T00:00:00Z",
    "end_date": "2024-07-23T00:00:00Z",
    "location": "France",
    "category": "Grand Tour",
    "is_free": false,
    "price_cents": 999,
    "requires_login": false,
    "stage_name": "Stage 17",
    "stage_type": "Mountain",
    "elevation_meters": 4800,
    "stage_length_km": 166,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

---

### Get Race by ID

**GET** `/races/:id`

Get details for a specific race.

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "id": "uuid",
  "name": "Tour de France",
  "description": "Annual cycling race",
  "start_date": "2024-07-01T00:00:00Z",
  "end_date": "2024-07-23T00:00:00Z",
  "location": "France",
  "category": "Grand Tour",
  "is_free": false,
  "price_cents": 999,
  "requires_login": false,
  "stage_name": "Stage 17",
  "stage_type": "Mountain",
  "elevation_meters": 4800,
  "stage_length_km": 166,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### Get Public User Profile

**GET** `/profiles/:id`

Get public profile information for a user.

**Parameters:**
- `id` (path, required) - User UUID

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "bio": "Cycling enthusiast",
  "points": 150,
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

## Authentication Endpoints

### Register

**POST** `/auth/register`

Create a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe",
  "bio": "Cycling enthusiast"
}
```

**Response:**
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "bio": "Cycling enthusiast",
    "points": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### Login

**POST** `/auth/login`

Authenticate and get a JWT token.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "bio": "Cycling enthusiast",
    "points": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Stream Endpoints

### Get Race Stream

**GET** `/races/:id/stream`

Get stream information for a race. Requires authentication for paid races and races with `requires_login = true`.

**Parameters:**
- `id` (path, required) - Race UUID

**Authentication:** Optional (required for paid races or races with `requires_login = true`)

**Response:**
```json
{
  "status": "live",
  "stream_type": "hls",
  "source_id": "youtube-video-id",
  "cdn_url": "https://cdn.example.com/hls/stream.m3u8",
  "origin_url": "http://origin.example.com/hls/stream.m3u8"
}
```

**Error Responses:**
- `401` - Authentication required (for paid races or races with `requires_login = true`)
- `403` - Payment required to access this race
- `404` - Stream not found for this race

**Note:** If a race has `requires_login = true`, authentication is required to access the stream, even if the race is free. The `401` error will be returned with the message "Authentication required to access this stream".

---

### Get Stream Status

**GET** `/races/:id/stream/status`

Get the current status of a race stream.

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "status": "live"
}
```

Status values: `planned`, `live`, `ended`, `offline`

---

## Viewer Tracking Endpoints

### Start Viewer Session

**POST** `/viewers/sessions/start`

Start tracking a viewer session (for analytics).

**Request:**
```json
{
  "race_id": "uuid"
}
```

**Authentication:** Optional (authenticated users are tracked separately)

**Response:**
```json
{
  "session_id": "uuid",
  "session_token": "token-for-anonymous"
}
```

---

### End Viewer Session

**POST** `/viewers/sessions/end`

End a viewer session.

**Request:**
```json
{
  "session_id": "uuid"
}
```

**Authentication:** Optional

**Response:**
```json
{
  "message": "Session ended"
}
```

---

### Viewer Session Heartbeat

**POST** `/viewers/sessions/heartbeat`

Update the last seen timestamp for an active viewer session.

**Request:**
```json
{
  "session_id": "uuid"
}
```

**Authentication:** Optional

**Response:**
```json
{
  "message": "Heartbeat received"
}
```

---

### Get Concurrent Viewers

**GET** `/races/:id/viewers/concurrent`

Get the current number of concurrent viewers for a race.

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "race_id": "uuid",
  "concurrent_count": 150,
  "authenticated_count": 120,
  "anonymous_count": 30
}
```

---

### Get Unique Viewers

**GET** `/races/:id/viewers/unique`

Get the total number of unique viewers for a race (all time).

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "race_id": "uuid",
  "unique_viewer_count": 500,
  "unique_authenticated_count": 400,
  "unique_anonymous_count": 100
}
```

---

## Chat Endpoints

### Get Chat History

**GET** `/races/:id/chat/history`

Get chat message history for a race.

**Parameters:**
- `id` (path, required) - Race UUID
- `limit` (query, optional) - Number of messages to return (default: 50)
- `offset` (query, optional) - Offset for pagination (default: 0)

**Response:**
```json
{
  "messages": [
    {
      "id": "uuid",
      "race_id": "uuid",
      "user_id": "uuid",
      "username": "John Doe",
      "message": "Great race!",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 150
}
```

---

### Get Chat Stats

**GET** `/races/:id/chat/stats`

Get chat statistics for a race.

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "race_id": "uuid",
  "total_messages": 1500,
  "unique_users": 200,
  "messages_per_minute": 5.2
}
```

---

### Chat WebSocket

**GET** `/races/:id/chat/ws`

WebSocket endpoint for real-time chat. Requires authentication.

**Parameters:**
- `id` (path, required) - Race UUID

**Authentication:** Required

**WebSocket Protocol:**
- Connect with JWT token in query parameter: `?token=<jwt-token>`
- Send messages as JSON: `{"type": "message", "data": {"message": "Hello"}}`
- Receive messages: `{"type": "message", "data": {...}}`

---

## User Endpoints (Authenticated)

### Get Profile

**GET** `/users/me`

Get the authenticated user's profile.

**Authentication:** Required

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "bio": "Cycling enthusiast",
  "points": 150,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### Change Password

**POST** `/users/me/password`

Change the user's password.

**Authentication:** Required

**Request:**
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword"
}
```

**Response:**
```json
{
  "message": "Password updated successfully"
}
```

---

### Award Watch Points (Automatic Tick)

**POST** `/users/me/points/tick`

Award points for watching (called automatically every 10 seconds while watching).
Awards 10 points per tick.

**Authentication:** Required

**Request:** No body required

**Response:**
```json
{
  "message": "Watch points awarded",
  "points": 10,
  "total_points": 150
}
```

---

### Award Bonus Points (Manual Claim)

**POST** `/users/me/points/bonus`

Award bonus points to the user (manual claim button in UI).
Awards a fixed 50 points per claim.

**Authentication:** Required

**Request:** No body required

**Response:**
```json
{
  "message": "Bonus points awarded",
  "bonus_points": 50,
  "total_points": 200
}
```

---

### Create Checkout Session

**POST** `/users/payments/create-checkout`

Create a Stripe checkout session for purchasing race access.

**Authentication:** Required

**Request:**
```json
{
  "race_id": "uuid"
}
```

**Response:**
```json
{
  "checkout_url": "https://checkout.stripe.com/...",
  "session_id": "cs_..."
}
```

---

### Start Watch Session

**POST** `/users/watch/sessions/start`

Start tracking watch time for a race.

**Authentication:** Required

**Request:**
```json
{
  "race_id": "uuid"
}
```

**Response:**
```json
{
  "session_id": "uuid",
  "started_at": "2024-01-01T12:00:00Z"
}
```

---

### End Watch Session

**POST** `/users/watch/sessions/end`

End a watch session and record watch time.

**Authentication:** Required

**Request:**
```json
{
  "session_id": "uuid"
}
```

**Response:**
```json
{
  "session_id": "uuid",
  "duration_seconds": 3600,
  "duration_minutes": 60
}
```

---

### Get Watch Stats

**GET** `/users/watch/sessions/stats/:race_id`

Get watch time statistics for a user and race.

**Authentication:** Required

**Parameters:**
- `race_id` (path, required) - Race UUID

**Response:**
```json
{
  "race_id": "uuid",
  "total_sessions": 5,
  "total_seconds": 18000,
  "total_minutes": 300,
  "first_watched": "2024-01-01T10:00:00Z",
  "last_watched": "2024-01-01T15:00:00Z"
}
```

---

## Webhook Endpoints

### Stripe Webhook

**POST** `/webhooks/stripe`

Handle Stripe webhook events. Validates Stripe signature.

**Headers:**
- `Stripe-Signature` (required) - Stripe webhook signature

**Request Body:** Raw Stripe webhook payload

**Response:**
```json
{
  "received": true
}
```

---

## Admin Endpoints

All admin endpoints require admin authentication (JWT token with `is_admin: true`).

### Create Race

**POST** `/admin/races`

Create a new race.

**Authentication:** Admin required

**Request:**
```json
{
  "name": "Tour de France",
  "description": "Annual cycling race",
  "start_date": "2024-07-01T00:00:00Z",
  "end_date": "2024-07-23T00:00:00Z",
  "location": "France",
  "category": "Grand Tour",
  "is_free": false,
  "price_cents": 999
}
```

**Response:**
```json
{
  "id": "uuid",
  "name": "Tour de France",
  ...
}
```

---

### Update Race

**PUT** `/admin/races/:id`

Update an existing race.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Race UUID

**Request:** Same as Create Race (all fields optional)

**Response:** Updated race object

---

### Delete Race

**DELETE** `/admin/races/:id`

Delete a race.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "message": "Race deleted"
}
```

---

### Update Stream

**POST** `/admin/races/:id/stream`

Create or update stream information for a race.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Race UUID

**Request:**
```json
{
  "status": "live",
  "stream_type": "hls",
  "source_id": "youtube-video-id",
  "origin_url": "http://origin.example.com/hls/stream.m3u8",
  "cdn_url": "https://cdn.example.com/hls/stream.m3u8",
  "stream_key": "secret-stream-key"
}
```

**Stream Types:** `hls` (default), `youtube`

**Response:** Stream object

---

### Update Stream Status

**PUT** `/admin/races/:id/stream/status`

Update only the stream status.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Race UUID

**Request:**
```json
{
  "status": "live"
}
```

**Response:** Stream object

---

### Get Revenue

**GET** `/admin/revenue`

Get all revenue data with optional filters.

**Authentication:** Admin required

**Query Parameters:**
- `race_id` (optional) - Filter by race
- `year` (optional) - Filter by year
- `month` (optional) - Filter by month

**Response:**
```json
[
  {
    "id": "uuid",
    "race_id": "uuid",
    "race_name": "Tour de France",
    "year": 2024,
    "month": 7,
    "total_revenue_cents": 50000,
    "total_revenue_dollars": 500.00,
    "total_watch_minutes": 10000.5,
    "platform_share_cents": 25000,
    "platform_share_dollars": 250.00,
    "organizer_share_cents": 25000,
    "organizer_share_dollars": 250.00,
    "calculated_at": "2024-08-01T00:00:00Z"
  }
]
```

---

### Get Revenue by Race

**GET** `/admin/revenue/races/:id`

Get revenue data for a specific race.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Race UUID

**Response:** Array of revenue objects (same format as above)

---

### Get Revenue Summary by Race

**GET** `/admin/revenue/races/:id/summary`

Get aggregated revenue summary for a race.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Race UUID

**Response:**
```json
{
  "race_id": "uuid",
  "race_name": "Tour de France",
  "total_revenue_cents": 100000,
  "total_revenue_dollars": 1000.00,
  "total_watch_minutes": 20000.5,
  "total_platform_share_cents": 50000,
  "total_platform_share_dollars": 500.00,
  "total_organizer_share_cents": 50000,
  "total_organizer_share_dollars": 500.00
}
```

---

### Recalculate Revenue

**POST** `/admin/revenue/recalculate`

Recalculate all revenue data.

**Authentication:** Admin required

**Response:**
```json
{
  "message": "Revenue recalculation started"
}
```

---

### Recalculate Revenue for Period

**POST** `/admin/revenue/recalculate/:year/:month`

Recalculate revenue for a specific month.

**Authentication:** Admin required

**Parameters:**
- `year` (path, required) - Year (e.g., 2024)
- `month` (path, required) - Month (1-12)

**Response:**
```json
{
  "message": "Revenue recalculated for 2024-07"
}
```

---

### Get Race Analytics

**GET** `/admin/analytics/races`

Get analytics data for all races.

**Authentication:** Admin required

**Response:**
```json
[
  {
    "race_id": "uuid",
    "race_name": "Tour de France",
    "concurrent_viewers": 150,
    "authenticated_viewers": 120,
    "anonymous_viewers": 30,
    "unique_viewers": 500,
    "unique_authenticated": 400,
    "unique_anonymous": 100
  }
]
```

---

### Get Watch Time Analytics

**GET** `/admin/analytics/watch-time`

Get watch time analytics.

**Authentication:** Admin required

**Response:**
```json
[
  {
    "race_id": "uuid",
    "race_name": "Tour de France",
    "total_watch_minutes": 10000.5,
    "unique_watchers": 200,
    "average_watch_minutes": 50.0
  }
]
```

---

### Get Revenue Analytics

**GET** `/admin/analytics/revenue`

Get revenue analytics.

**Authentication:** Admin required

**Response:**
```json
[
  {
    "race_id": "uuid",
    "race_name": "Tour de France",
    "total_revenue_cents": 50000,
    "total_revenue_dollars": 500.00,
    "total_tickets_sold": 100
  }
]
```

---

### Create Cost

**POST** `/admin/costs`

Create a new cost entry.

**Authentication:** Admin required

**Request:**
```json
{
  "race_id": "uuid",
  "cost_type": "cdn",
  "amount_cents": 10000,
  "year": 2024,
  "month": 7,
  "description": "CDN bandwidth costs"
}
```

**Cost Types:** `cdn`, `server`, `storage`, `bandwidth`, `other`

**Response:** Cost object

---

### Get Costs

**GET** `/admin/costs`

Get all costs with optional filters.

**Authentication:** Admin required

**Query Parameters:**
- `race_id` (optional) - Filter by race
- `year` (optional) - Filter by year
- `month` (optional) - Filter by month
- `cost_type` (optional) - Filter by cost type

**Response:** Array of cost objects

---

### Get Cost Summary

**GET** `/admin/costs/summary`

Get aggregated cost summary.

**Authentication:** Admin required

**Query Parameters:**
- `race_id` (optional) - Filter by race
- `year` (optional) - Filter by year
- `month` (optional) - Filter by month

**Response:**
```json
{
  "race_id": "uuid",
  "year": 2024,
  "month": 7,
  "cdn_cents": 5000,
  "server_cents": 3000,
  "storage_cents": 1000,
  "bandwidth_cents": 2000,
  "other_cents": 0,
  "total_cents": 11000,
  "total_dollars": 110.00
}
```

---

### Get Cost by ID

**GET** `/admin/costs/:id`

Get a specific cost entry.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Cost UUID

**Response:** Cost object

---

### Update Cost

**PUT** `/admin/costs/:id`

Update a cost entry.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Cost UUID

**Request:** Same as Create Cost (all fields optional)

**Response:** Updated cost object

---

### Delete Cost

**DELETE** `/admin/costs/:id`

Delete a cost entry.

**Authentication:** Admin required

**Parameters:**
- `id` (path, required) - Cost UUID

**Response:**
```json
{
  "message": "Cost deleted"
}
```

---

### Get Costs by Race

**GET** `/admin/costs/races/:race_id`

Get all costs for a specific race.

**Authentication:** Admin required

**Parameters:**
- `race_id` (path, required) - Race UUID

**Response:** Array of cost objects

---

## Notes

- All timestamps are in ISO 8601 format with timezone (UTC)
- All UUIDs are in standard UUID v4 format
- All monetary amounts are in cents (integer) to avoid floating-point precision issues
- Pagination is not implemented for most list endpoints (consider adding for production)
- Rate limiting may vary by endpoint type
- CSRF protection is enabled for state-changing operations (POST, PUT, DELETE)

