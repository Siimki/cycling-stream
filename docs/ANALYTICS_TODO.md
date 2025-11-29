# Analytics TODO (Bunny Stream + Custom Tracking)
Step-by-step implementation guide for the Go backend (`backend/`) and Next.js frontend (`frontend/`). Keep handlers thin (helpers in `internal/handlers/http_helpers.go`), name files by concern, and run `gofmt`/lint where applicable.

## Phase 1 – Tracking MVP (fast path)
- Data model: add migrations for `streams`, `stream_providers`, `playback_events` (schema only, no sessions yet). Seed minimal dev streams/providers in a script under `scripts/` if helpful.
- Backend ingestion: implement `POST /analytics/events` (batch) in `backend/cmd/api` stack. Validate `stream_id`, constrain batch size, derive country/device (stub device via UA, country as placeholder until GeoIP available), and bulk insert events.
- Frontend player: build `<StreamPlayer>` that selects provider (Bunny/YouTube) from `GET /api/streams/:id`; embed Bunny script and YouTube IFrame API; generate `clientId` (localStorage) and send batched `play/heartbeat/ended` events (heartbeat ~15s, debounced sends). Keep network base URLs in env (`NEXT_PUBLIC_API_URL`).
- Admin stopgap: simple internal admin table/view or SQL query to show per-stream `distinct clientId` and rough watch time (heartbeat count × interval).
- Test plan: manual playback in Next.js hitting Bunny/YouTube; inspect network calls to `/analytics/events`; Go unit test for payload validation; run `make run-backend` + `make run-frontend` smoke; SQL sanity query for counts.

## Phase 2 – Sessions & Aggregates
- Data model: add migrations for `viewer_sessions` and `stream_stats` with indexes on `stream_id`, `client_id`, `created_at`.
- Sessioning in ingestion: reuse/create session on `(stream_id, client_id)` with inactivity timeout (e.g., 30 min); store country/device on session; update `ended_at` and provisional `total_watch_seconds` on new batches; keep ingestion lightweight.
- Aggregators (worker/cron in `cmd` or `scripts/`): finalize stale sessions; compute `stream_stats` (unique viewers, total/avg watch seconds, peak concurrent via sweep over session intervals, top countries/device breakdown); stamp `last_calculated_at`. Add Make target to run jobs locally.
- Admin APIs: `GET /admin/analytics/streams` (filters: date range, organizer) and `GET /admin/analytics/summary`. Protect with existing admin auth.
- Admin UI: real table with KPIs + filters + CSV export (either client-side CSV or server CSV endpoint); keep components in `frontend/components/` with PascalCase; hooks/utilities camelCase.
- Test plan: Go unit tests for session reuse/timeout and aggregation math (including peak concurrent); run aggregator locally with fixture data; frontend e2e smoke for filters/CSV; run scoped `make test-backend` on analytics packages.

## Phase 3 – Bunny Analytics & QoE
- Data model: add `bunny_video_stats` (per day) with optional `raw_payload` JSONB; extend `playback_events.extra` if needed for QoE fields.
- Bunny integration: scheduled job to pull analytics for `provider = bunny_stream` (map provider IDs from `stream_providers`); upsert stats; env for Bunny API key/URL; store watch time/views/geo breakdown.
- Tracking extensions: capture `error`, `buffer_start`, `buffer_end`; send `extra` payloads (error codes, network info) from the player; update ingestion to persist.
- QoE metrics: compute % sessions with errors and buffering ratio (`buffer_seconds/watch_seconds`) into `stream_stats` or a sibling table; expose via admin APIs.
- Admin UI: optional Bunny stats block/link per stream and QoE percentages; ensure Bunny data only renders for Bunny streams.
- Test plan: mock Bunny API client unit tests; dry-run job against sample payload; UI check that Bunny stats gate on provider; unit tests for QoE calculations.

## Phase 4 – Hardening & Scale
- Performance: tune batch inserts and DB pool; consider monthly partitioning for `playback_events`; add indexes as volume grows.
- Archival/retention: job to archive/delete raw events older than N months once aggregates exist; document retention policy in `docs/`.
- Privacy/compliance: drop/avoid storing raw IP after geo derivation; confirm cookie/localStorage notice plan; ensure anonymous `clientId` usage.
- Observability/ops: add logging/metrics around ingestion failures and job runs; dashboards if available; Make/cron wiring for jobs and update `README.md`/`.cursorrules` if commands change.
- Test plan: load test ingestion with batched events; verify partitions/archival jobs execute; lint/format passes; smoke `make docker-up` + run aggregations end-to-end.

## Progress log
- Added schema for `stream_providers` and `playback_events` (v1 event log).
- Backend: new `POST /analytics/events` ingestion with batching, stream existence check, device/country stubs, and optional auth context.
- Frontend: analytics hook to batch player events and wired into `StreamProvider` + video player (HLS + minimal YouTube heartbeat). Stream responses now include `stream_id` for analytics payloads.
- Still pending: scheduled jobs, admin analytics UI/CSV, and richer YouTube event tracking (needs iframe API).
- Phase 2 start: added `stream_stats` table and stream-level columns on `viewer_sessions`, stream stats repo, and an aggregator that derives sessions/stats from `playback_events`; admin routes expose `/admin/analytics/streams` and `/admin/analytics/streams/summary`. Still need scheduled job wiring, CSV export, and UI.
- Phase 3: added `bunny_video_stats` table, Bunny client/importer (admin-triggered `/admin/analytics/streams/bunny-sync`), QoE fields on `stream_stats` (buffer ratio, error rate), and extended player tracking to send buffer events; aggregator now computes QoE. Bunny env vars documented in `backend/.env.example`.
- Phase 4 groundwork: no partitions/archival wired yet; logging/cron wiring still pending; need to add admin UI/CSV and scheduled jobs.
- Admin UI: added `/admin/analytics/streams` page for stream stats (table + CSV export) and link from admin home; Bunny remains gated by env (sync endpoint returns 501 if not configured).

## End-to-end test checklist (after all phases)
- Phase 1: run `make migrate-up`; start backend/frontend; play a stream (HLS & YouTube) and confirm `/analytics/events` receives play/heartbeat/ended; verify playback_events rows.
- Phase 2: call `GET /admin/analytics/streams` after events; ensure counts/peak concurrency make sense; `GET /admin/analytics/streams/summary` returns aggregates; run `go test ./internal/...` for analytics packages.
- Phase 3: set Bunny envs (`BUNNY_API_KEY`, `BUNNY_LIBRARY_ID`); seed `stream_providers` with `provider=bunny_stream`; hit `POST /admin/analytics/streams/bunny-sync` and confirm `bunny_video_stats` rows and `stream_stats` updates; induce buffering/errors in player to see QoE fields populate.
- Phase 4: (once wired) run retention/partition job dry-run, load-test `/analytics/events`, and confirm logs/metrics captured; verify admin CSV/export/UI once added.

## Detailed test plan (execute when ready)
### Prereqs
- Migrations applied (`make migrate-up`).
- Backend `.env` with API base/DB and optional Bunny vars; frontend `NEXT_PUBLIC_API_URL` set.
- At least one stream with `streams.id`, race set to live, and (if Bunny) a `stream_providers` row with `provider=bunny_stream`.

### Backend unit/integration
- Ingestion validation: reject missing `streamId`, missing `clientId`, >100 events, negative `videoTime`, invalid event types; accept allowed types (play/pause/heartbeat/ended/error/buffer_start/buffer_end).
- Aggregator math: with fixture events, assert unique viewers, total/avg watch seconds, peak CCV sweep, buffer ratio, error rate; idle timeout splits sessions.
- Bunny client/importer: mock HTTP to return analytics payload; importer upserts `bunny_video_stats` and updates `stream_stats`; importer returns 501 when Bunny disabled.
- Retention helper: delete older playback events returns expected count (use temp table/tx).

### API smokes
- `/analytics/events`: send batched events for a known stream/client; expect 202 + rows in `playback_events`.
- `/admin/analytics/streams`: returns stats including QoE fields; `/summary` returns aggregate counts.
- `/admin/analytics/streams/bunny-sync`: returns 501 when Bunny env missing; with env + provider row, returns success and creates `bunny_video_stats`.

### Frontend/e2e
- Player tracking: play/pause/end HLS stream; verify network POSTs with expected event types and streamId; buffer (throttle network) to see buffer_start/end events; introduce error to see error event.
- Admin stream analytics page (`/admin/analytics/streams`): table renders, summary cards show values, CSV export downloads with rows, error state if API fails.
- Existing admin analytics page: ensure added “Stream Analytics” link works.

### Non-functional
- Load: fire concurrent batched `/analytics/events` (e.g., k6/JMeter) to ensure ingestion stays responsive within limits; monitor DB writes.
- Retention (when enabled): dry-run delete on a staging copy to confirm only old rows are removed.
- Auth/rate limits: confirm admin endpoints require token; ingestion respects existing rate-limit middleware.

### Expected outputs (baseline)
- Ingestion: 202 with `{ ingested: N }`; `playback_events` row count matches events sent.
- Stream stats: for a single client sending 4 heartbeats @15s, total_watch_seconds ≈ 60, avg_watch_seconds ≈ 60, peak_concurrent_viewers = 1.
- QoE: if buffering for ~5s once, buffer_ratio ≈ 5/60; if one error sent, error_rate ≈ 1.0 for single viewer.
- Bunny: after sync, `bunny_video_stats` row for today with views/watch_time from mock; `stream_stats.total_watch_seconds` updated to max of existing vs Bunny watch time.
