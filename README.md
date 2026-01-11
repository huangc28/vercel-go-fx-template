# Vercel Go + FX + Chi template

Conservative template for Go services deployed to Vercel serverless functions:
- Per-domain entrypoints: `api/<domain>/core.go`
- Per-request DI via Uber FX
- HTTP routing via `chi`
- Minimal `/health` example and optional `/api/inngest` entrypoint

## Create a new service (recommended: `gonew`)

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/huangc28/vercel-go-fx-template github.com/<org>/<new-service>
cd <new-service>
```

## Adopt into an existing repo (non-destructive)

To add the guidance files into an existing repo (without overwriting its layout), run the adopter CLI from inside that repo:

```bash
GOPROXY=direct go run github.com/huangc28/vercel-go-fx-template/cmd/adopt@main --dir .
```

This writes:
- `AGENTS.md`
- `architecture/go-vercel-reusable-template-plan.md`
- `codex/skills/adopt/SKILL.md`

To enable `/adopt` in Codex, install the skill to your Codex home (commonly `~/.codex/skills/adopt`).

## Run locally

```bash
make start/vercel APP_PORT=3010
```

Requires Vercel CLI auth (`vercel login`) and a linked project (`vercel link`).

Then:
- `GET /health`

## Add a new domain

1) Create `api/<domain>/core.go` (copy `api/health/core.go`).
2) Create handler(s) under `lib/app/<domain>/...` implementing `lib/router.Handler`.
3) Register handler constructors in `api/<domain>/core.go` via `router.AsRoute(...)`.
4) Add `vercel.json` rewrite, e.g. `/v2/<domain>/*` -> `/api/<domain>/core`.

## Env vars

Set env vars in Vercel (or locally) to match `config/config.go`:
- `APP_NAME` (default: `vercel-go-service`)
- `APP_ENV` (default: `development`)
- `APP_PORT` (default: `3010`)
- `LOG_LEVEL` (default: `info`)
- `PG_URL` (optional; if empty Postgres is disabled)
- `REDIS_URL` (optional; if empty Redis is disabled)
- `INNGEST_APP_ID` (default: `vercel-go-service`)
- `INNGEST_EVENT_KEY` (optional; defaults to env lookup in Inngest SDK)
- `INNGEST_SIGNING_KEY` (optional; defaults to env lookup in Inngest SDK)
- `INNGEST_SIGNING_KEY_FALLBACK` (optional)

## sqlc (manual)

1) Put your Supabase-exported schema at `supabase/schema.sql`.
2) Install sqlc:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3) Run:

```bash
sqlc generate
```

Generated code lands in `lib/app/models` (see `sqlc.yaml`).
