# Vercel Go + FX + Chi template

Internal template for Go services deployed to Vercel serverless functions.

Baseline contract:
- Per-domain entrypoints: `api/<domain>/core.go`
- Per-request dependency injection via Uber FX
- HTTP routing via `chi`
- Minimal `/health` example only
- Optional Postgres and Redis wiring through env vars

This template is intentionally small. It does not ship business logic, Inngest, sqlc, Supabase, Turso, or NotebookLM bootstrap.

## Create a new service

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/huangc28/vercel-go-fx-template github.com/<org>/<new-service>
cd <new-service>
```

## Run locally

```bash
make start/vercel APP_PORT=3010
```

Requires Vercel CLI auth (`vercel login`) and a linked project (`vercel link`).

Smoke test:
- `GET /health`

## Add a new domain

1. Create `api/<domain>/core.go` by copying `api/health/core.go`.
2. Create handler(s) under `lib/app/<domain>/...` implementing `lib/router.Handler`.
3. Register handler constructors in `api/<domain>/core.go` via `router.AsRoute(...)`.
4. Add a `vercel.json` rewrite such as `/v1/<domain>/(.*)` -> `/api/<domain>/core`.

## Env vars

Set env vars in Vercel or locally to match `config/config.go`:
- `APP_NAME` default: `vercel-go-service`
- `APP_ENV` default: `development`
- `APP_PORT` default: `3010`
- `LOG_LEVEL` default: `info`
- `PG_URL` optional; if empty Postgres is disabled
- `REDIS_URL` optional; if empty Redis is disabled

## Non-goals

This template does not include:
- existing-repo adoption tools
- long-running local app entrypoints
- workflow engines such as Inngest
- database code generation or schema management
- product-specific domains or third-party integrations
