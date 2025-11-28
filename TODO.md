# TODO

- Remove hardcoded admin bypass in `backend/internal/handlers/auth.go` (lines ~170-185); use stored admins or guard behind a dev-only flag.
- Fix `/users/me` panic risk in `backend/internal/handlers/auth.go` by safely reading `user_id` from locals and returning 401 when absent.
- Rework CSRF middleware (`backend/internal/middleware/csrf.go`): correct `Expiration` to a `time.Duration`, and either remove the always-skip `Next` or clearly drop the middleware if JWT-only.
- Tighten rate limits in `backend/internal/middleware/ratelimit.go`: lower per-IP limits (especially auth/payments), consider honoring `X-Forwarded-For`, and avoid the current 5k–10k req/min prod-unsafe defaults.
- Chat spam throttle regression: `backend/internal/chat/ratelimit.go` bumped to 100 msgs/min without env gating; restore a stricter production default or make higher limit test-only.
