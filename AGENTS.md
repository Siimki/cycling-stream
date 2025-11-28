# Repository Guidelines

## Project Structure & Module Organization
- Go backend in `backend/` (entry `cmd/api/main.go`, logic in `internal/`, SQL in `migrations/`, helper scripts in `scripts/`).
- Next.js frontend in `frontend/` (`app/`, `components/`, `lib/`, assets in `public/`); shared docs in `docs/`, `DESIGN.md`, and other root notes.
- Infra/streaming lives in `stream/`; local services orchestrated by `docker-compose.yml`.
- Copy env templates (`backend/.env.example`, `frontend/.env.local.example`) before running.

## Build, Test, and Development Commands
- Install deps: `make setup`.
- Local stack: `make docker-up` (Postgres + pgAdmin), `make run-backend`, `make run-frontend` in a second terminal.
- Migrations: `make migrate-up` (needs `migrate` CLI) or `make migrate-create NAME=...`.
- Quality gates: `make test`, or per side `make test-backend` / `make test-frontend`; lint with `make lint-backend` / `make lint-frontend`; build via `make build-backend` / `make build-frontend`.

## Coding Style & Naming Conventions
- Go: always `gofmt`; keep handlers thin with helpers in `internal/handlers/http_helpers.go` and related utils; name files by concern (`*_repository.go`, `*_service.go`).
- Frontend: TypeScript + functional components; PascalCase for components in `components/`, camelCase for hooks/utilities (`useChat`, `chatService`); lean on Tailwind utilities with `tailwind-merge`; shared values in `constants/` or `lib/`.
- Keep secrets in local `.env` files; browser calls should use `NEXT_PUBLIC_API_URL`.

## Testing Guidelines
- Go tests sit beside code as `*_test.go`; DB-touching suites need `make docker-up`, env vars from `.env.example`, and a Postgres role `cyclingstream` (per QA plan `go test ./...` fails without it). Use `testify`.
- Frontend has placeholder specs (`hooks/useChat.test.ts`, `components/Chat.test.tsx`) but Jest is not wired; add React Testing Library tests and update `npm test` when enabling. Repo-wide `npm run lint` currently surfaces existing debt—scope linting if needed until cleanup.
- For new endpoints or UI flows, include a happy-path test and note seed/data needs in the PR.

## Commit & Pull Request Guidelines
- Match existing conventional subjects (`feat: ...`, `refactor: ...`, `update: ...`); keep scopes small and imperative.
- PRs should list a concise summary, linked issue/mission ID, migrations/env steps, test results (`make test-backend`, `make lint-frontend`, etc.), and UI screenshots/GIFs when visuals change.
- Avoid mixing refactors with features; flag breaking changes and update docs/Makefile targets when commands shift.

## Security & Configuration Tips
- Never commit `.env*` or generated secrets; rotate `JWT_SECRET` and DB credentials per environment.
- Keep local services bound to localhost; verify `PGADMIN_PORT`/`DB_PORT` conflicts before `make docker-up`.
- Document port/CDN expectations when editing streaming config under `stream/`.
