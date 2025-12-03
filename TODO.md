# TODO

## Codebase Cleanup (Dec 2025)
- ✅ Removed redundant audit documentation files (audit-opus.md, audit-cursor.md, audit-gemini.md, audit-codex.md)
- ✅ Removed test/plan documentation files (real-time-chat-plan.md, real-time-chat-test.md, chat.md, frontend-refactor-sum1.md, PHASE_6_MISSIONS_STATUS.md, ANALYTICS_STATUS_SUMMARY.md)
- ✅ Removed test artifacts (chat-test-output.log, chat-edge-case-output.log, *.json test results, test-chat-*.js scripts)
- ✅ Removed root package.json/package-lock.json (dependencies moved to frontend/)
- ✅ Removed compiled binaries (backend/api, backend/main)
- ✅ Removed unused frontend code (designerplan/, VideoPlayer.tsx wrapper, Chat.test.tsx, useChat.test.ts, lib/user-segmentation.ts, hooks/useRipple.ts duplicate)
- ✅ Updated .gitignore for build artifacts and test outputs

## Security & Configuration
- Remove hardcoded admin bypass in `backend/internal/handlers/auth.go` (lines ~170-185); use stored admins or guard behind a dev-only flag.
- Fix `/users/me` panic risk in `backend/internal/handlers/auth.go` by safely reading `user_id` from locals and returning 401 when absent.
- Rework CSRF middleware (`backend/internal/middleware/csrf.go`): correct `Expiration` to a `time.Duration`, and either remove the always-skip `Next` or clearly drop the middleware if JWT-only.
- Tighten rate limits in `backend/internal/middleware/ratelimit.go`: lower per-IP limits (especially auth/payments), consider honoring `X-Forwarded-For`, and avoid the current 5k–10k req/min prod-unsafe defaults.
- Chat spam throttle regression: `backend/internal/chat/ratelimit.go` bumped to 100 msgs/min without env gating; restore a stricter production default or make higher limit test-only.

## Missions System (85% complete)
- Enhanced Points Display - Add missions section to `frontend/components/user/PointsDisplay.tsx`
- Dynamic Watch Bonuses - Context-based bonuses (requires race event data)
- Season Progression System - Season pass/tour card UI

## Analytics System
- Fix ingestion endpoint auth issue (POST /analytics/events returns 401)
- Add scheduled jobs for aggregator and Bunny sync
- Update API documentation for analytics endpoints
- Investigate peak concurrent viewers calculation
- Add Bunny stats to admin UI
