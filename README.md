Here is the professional English version of the README.md file for your project root.

Create a file named README.md and paste this content in.

Markdown

# Pro Cycling Streaming Platform 🚴‍♂️💨

This is an open-source streaming platform (MVP) designed to provide fair monetization and distribution channels for smaller cycling races (Tier 2/3 events). The project is inspired by the "Fairmus" model (user-centric revenue share).

## 🏗 Architecture & Tech Stack

The project follows a **monorepo** structure.

| Component | Technology | Description |
|-----------|-------------|-----------|
| **Backend** | **Go** (Fiber) | REST API, business logic, database interaction. |
| **Database** | **PostgreSQL** | Users, races, ticketing/entitlements. |
| **Frontend** | **Next.js** (App Router) | React, TypeScript, Tailwind CSS. User Interface. |
| **Streaming** | **Owncast / Nginx** | RTMP ingest (from OBS) and HLS output. |
| **CDN** | **BunnyCDN** | Global distribution of HLS streams (planned). |
| **Infra** | **Docker** & **Hetzner** | Containerization and server hosting. |

## 📂 Project Structure

```text
cycling-stream-platform/
├── backend/             # Go API source code
│   ├── cmd/api/         # Server entrypoint
│   ├── internal/        # Business logic (handlers, models, db)
│   └── migrations/      # SQL files for database schema changes
├── frontend/            # Next.js application
│   ├── app/             # Pages and routing
│   └── components/      # UI components
├── stream/              # Streaming server config (Docker/Owncast)
├── docker-compose.yml   # Local development DB and services
├── Makefile             # Shortcuts (for running app, migrations, etc.)
└── README.md            # This file

## 🐹 Backend Patterns & Conventions

The Go backend follows a few helper-based patterns to keep handlers small and consistent:

- `backend/internal/handlers/http_helpers.go`
  - `parseBody(c, &req)` wraps `c.BodyParser` and returns a standard `400 {"error":"Invalid request body"}` on failure.
  - `requireParam(c, "id", "Race ID is required")` validates required path params and sends a `400` with the provided message if missing.
  - `requireUserID(c, "Authentication required")` reads `user_id` from context and sends a `401` if not present.
  - `APIError` is the shared error payload (`{"error": "..."}`) used in new/updated handlers.
- `backend/internal/handlers/race_utils.go`
  - `loadRaceOr404(c, raceRepo, id)` fetches a race and centralizes `500`/`404` responses for race lookups.
  - `loadStreamOr404(c, streamRepo, raceID, notFoundMessage)` does the same for streams by `race_id`.
- `backend/internal/handlers/session_utils.go`
  - `verifyViewerSessionOwnership(c, viewerSessionRepo, sessionID)` checks that a viewer session belongs to the authenticated user (when present) and returns `403` on mismatch.

**When adding new handlers**, prefer:

- Using `parseBody` instead of inlining `c.BodyParser` + `"Invalid request body"`.
- Using `requireParam` for required path params (`id`, `race_id`, `session_id`, etc.).
- Using `requireUserID` for endpoints that require authentication.
- Using `APIError` for simple error responses, and the helper utilities (`loadRaceOr404`, `loadStreamOr404`, `verifyViewerSessionOwnership`) where applicable.
## 🎯 Key Features

- **Live Race Streaming**: Watch live cycling races with HLS streaming support
- **User Accounts & Authentication**: Secure user registration and login
- **Pay-Per-View Tickets**: Stripe integration for race access
- **Personalization**: 
  - Onboarding wizard to customize your experience
  - Personalized "For You" page with recommendations
  - Favorite riders, teams, and races
  - Watch history tracking
- **Real-time Chat**: WebSocket-based chat during live races
- **Watch Time Tracking**: Track viewing statistics per user per race
- **Admin Dashboard**: Manage races, streams, and view analytics

## 🚀 Quickstart: Run the app locally

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- Make

### 1. Clone the repo

```bash
git clone https://github.com/yourusername/cycling-stream-platform.git
cd cycling-stream-platform
```

### 2. Start Postgres (Docker)

Postgres runs in Docker and is exposed on host port **5434** (see `docker-compose.yml`).

```bash
make docker-up
```

This will start:
- Postgres on `localhost:5434`
- pgAdmin on `http://localhost:5050`

### 3. Run database migrations

Preferred (if you have the `migrate` CLI installed):

```bash
make migrate-up
```

Alternative (using Docker + psql, if `migrate` is not installed):

```bash
cat backend/migrations/0*.sql | docker exec -i cyclingstream_postgres \
  psql -U cyclingstream -d cyclingstream
```

The backend is already configured to use:
- DB host: `localhost`
- DB port: `5434` (via `DB_PORT`, default is `5434`)

### 4. Start the backend API

```bash
make run-backend
```

- API base URL: `http://localhost:8080`
- Health check: `http://localhost:8080/health`

### 5. Start the frontend (Next.js)

In another terminal:

```bash
make run-frontend
```

- Frontend URL: `http://localhost:3000`
- The frontend talks to the backend via `NEXT_PUBLIC_API_URL` and defaults to `http://localhost:8080`.

### 6. Verify everything

- Visit `http://localhost:3000` in your browser.
- Hit `http://localhost:8080/health` and confirm it returns a JSON payload with `"status": "ok"`.
- Use `/auth/register` and `/auth/login` in the UI to create an account and log in.
🗺️ Project Roadmap (MVP)
Mark as [x] when the phase is complete.

P🌱 Phase 0 – Project Initialization & Dev Environment

Goal: One repo, boots up locally with backend + DB + frontend.

Create monorepo structure (backend/, frontend/, stream/, docker-compose.yml, Makefile, README.md).

Initialize:

Go module + Fiber in backend/

Next.js + Tailwind in frontend/

Docker Compose for Postgres (and maybe pgAdmin).

Simple .env convention (backend/.env, frontend/.env.local).

Basic CI:

Run go test ./...

Run npm test or npm run lint for frontend.

✅ Done when:
make docker-up, make run-backend, make run-frontend → app opens at localhost:3000 and backend /health returns OK.

🧠 Phase 1 – Backend Core (API Skeleton + Data Model)

Goal: Stable backend that knows about races & streams.

Config loading (env, ports, DB URL).

DB layer (Postgres connection, migrations).

Core tables:

users (even if basic)

races

streams

entitlements / tickets (even as minimal stub)

Public API:

GET /health

GET /races

GET /races/:id

GET /races/:id/stream (status + URL)

✅ Done when:
You can curl the API and see real data from Postgres, not hardcoded JSON.

🎨 Phase 2 – Frontend Core (Next.js UI Skeleton)

Goal: Users can browse races and open a “watch” page (even if video is dummy).

Layout, basic design system (Tailwind, global styles).

Pages:

Home: list of upcoming races (from backend).

Race detail: info page per race.

Watch page: /races/[id]/watch with placeholder video player.

Simple data fetching (REST → fetch/axios).

Basic error/loading states.

✅ Done when:
A friend can open your site, see race list, click into a race, and land on a watch page that looks ready for streaming (even if video is just a sample HLS URL).

📡 Phase 3 – Streaming Origin Setup (Owncast/Nginx on Hetzner)

Goal: You can stream from OBS → server → browser, without CDN yet.

Spin up Hetzner server.

Deploy Owncast or nginx-rtmp there (via Docker).

Configure:

RTMP ingest (rtmp://.../live/streamkey)

HLS output (public .m3u8 URL).

Test:

Push OBS → RTMP.

Open .m3u8 in VLC / browser player and see your stream.

✅ Done when:
You can stream from your laptop (OBS) to Hetzner and watch it live from anywhere via raw HLS URL.

🌍 Phase 4 – CDN Integration (BunnyCDN) + End-to-End MVP

Goal: Real viewer traffic flows through CDN, not straight from your server.

Create BunnyCDN Pull Zone, origin = Hetzner HLS path.

Get CDN HLS URL (https://yourzone.b-cdn.net/...m3u8).

Update backend:

streams.cdn_url = Bunny URL

GET /races/:id/stream returns real cdn_url.

Update frontend Watch page:

Use HLS-capable player (video.js, hls.js, Shaka …).

If status === 'live', play the Bunny URL.

Do live end-to-end test:

OBS → Hetzner → Bunny → Watch page in browser.

✅ Done when:
You can send one link to a friend (“watch this race”) and they see a real live stream through BunnyCDN.

🛠 Phase 5 – Admin & Content Operations

Goal: You don’t have to manually hack the DB to add races/streams.

Backend:

Basic auth for admin (JWT or simple password).

Admin endpoints:

POST /admin/races

PUT /admin/races/:id

POST /admin/races/:id/stream (attach URLs, toggle status planned/live/ended)

Frontend:

Simple admin page (could be behind HTTP basic auth at first).

Forms for creating/editing races.

Button/field to attach origin_url and cdn_url.

Helpers:

Seed scripts (create example races).

Logs for when races change status.

✅ Done when:
You can create a new race and its stream entirely via UI/admin, no psql or manual inserts needed.

💳 Phase 6 – Users, Tickets & Fair Monetization Model (MVP Level) ✅ (Mostly Complete)

Goal: Start moving from "toy project" → actual platform with users & payments.

**Completed:**
- ✅ Minimal user accounts (email + password)
- ✅ Sign up / login flow on frontend (/auth/register, /auth/login)
- ✅ Entitlements table and access checking
- ✅ Backend check on GET /races/:id/stream → only return URL if user has access (or race is free)
- ✅ Stripe Checkout integration for single race "ticket" (PPV)
- ✅ Webhook → when payment succeeds, create entitlement rows
- ✅ Watch time tracking - store per-user watch minutes per race
- ✅ Watch time statistics endpoint

**Remaining:**
- ⏳ Frontend payment flow UI (payment button/component)
- ⏳ End of month revenue split calculation (Phase 6.4)

🎯 Phase 6.5 – Personalization Platform ✅ (Complete)

Goal: Make each viewer feel the stream is built for them.

**Completed:**
- ✅ User preferences system (data mode, units, theme, device type, notifications)
- ✅ User favorites (riders, teams, races, series)
- ✅ Watch history tracking and aggregation
- ✅ 5-step onboarding wizard (cycling level, view preference, favorites, device, notifications)
- ✅ User segmentation logic (Casual Viewer, Hardcore Fan, Data Nerd, etc.)
- ✅ "For You" personalized page with:
  - Continue watching section
  - Upcoming races from favorite series
  - Recommended replays based on watch history
  - Personal stats (watch time, points, races watched)
- ✅ Rules-based recommendation service
- ✅ Personalized navigation ("For You" tab for logged-in users)
- ✅ WebSocket chat with improved retry logic (3s, 5s, 10s, 30s then stop)

✅ Done when:
Someone can create an account, pay for access, and then watch a paid race legitimately on your platform. Watching generates usable data for later revenue share.

📊 Phase 7 – Analytics, Monitoring & Cost Visibility

Goal: You understand what’s going on: viewers, costs, errors.

Metrics:

Viewer counts per race (at least “peak” and “total viewers”).

Watch time per user per race.

Admin dashboards:

Simple charts/tables: list races with viewer counts and total minutes.

Infra observability:

Basic logs for backend (errors, slow queries).

Hetzner monitoring (CPU, RAM, bandwidth).

BunnyCDN usage dashboards (GB/zone/day).

Cost tracking:

Rough cost per race (CDN GB + server).

Store these numbers in DB or simple Notion/Sheet at first.

✅ Done when:
If someone asks “how many people watched Race X, and what did it cost?” — you can answer from your data, not guess.

🚀 Phase 8 – Polish, Hardening & First “Real” Season

Goal: Turn the MVP into something you’re not ashamed to show organizers.

UX polish:

Proper responsive design (mobile/tablet).

Nicer race cards, flags, categories.

Robustness:

Better error pages (stream offline, network errors).

Retries, simple fallback messages in the player.

Legal/operational basics:

Terms of Service, Privacy Policy.

Contact page for organizers & takedown requests.

Content:

3–10 real low-tier UCI or regional races onboarded.

Ask for feedback from organizers and a few beta users.

✅ Done when:
You can run a small “season” of live streams and VODs without constantly firefighting basic bugs.
💡 Business Model Notes
Revenue Share: Ticket revenue is split 50/50 between the platform and the race organizer.

User-Centric: If a user only watches one specific race, their subscription money (after costs) goes only to that organizer.

Cost Optimization: We utilize "Basic" AWS/CDN tiers and 720p/60fps streaming to keep bandwidth costs sustainable.

Author: Siim