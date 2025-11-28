# Phase 6: Missions System - Status

**Status:** 85% complete

## What's Not Done

1. **Enhanced Points Display** - Missions not shown in `PointsDisplay` component
   - Add missions section to `frontend/components/user/PointsDisplay.tsx`
   - Show 2-3 active missions with progress bars and claim buttons

2. **Dynamic Watch Bonuses** - Context-based bonuses (e.g., "+100 if you stay until the next climb")
   - Requires race event data (climbs, sprints, etc.)

3. **Season Progression System** - Season pass/tour card UI with unlock tiers and rewards

## Key Files

**Backend:**
- `backend/internal/services/mission_service.go` - Core mission logic
- `backend/internal/services/mission_triggers.go` - Auto progress tracking
- `backend/internal/handlers/missions.go` - API handlers
- `backend/migrations/024_create_missions.sql`
- `backend/migrations/025_create_user_missions.sql`

**Frontend:**
- `frontend/components/missions/` - All mission components
- `frontend/app/missions/page.tsx` - Missions page
- `frontend/components/user/PointsDisplay.tsx` - Needs missions integration

**API Endpoints:**
- `GET /users/me/missions` - Get user missions
- `GET /missions/active` - Get active missions
- `POST /users/me/missions/:missionId/claim` - Claim reward

## Notes

- Mission triggers automatically track watch time, chat messages, and race watches
- `predict_winner` mission type defined but not implemented (no prediction system)
- Mission rewards automatically award points when claimed
