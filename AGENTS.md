# AGENTS.md (Codex)

## Architecture (hard rules)
- This repo uses Uber FX for dependency injection. No global singletons for DB/Redis/Config/Logger.
- Vercel serverless entrypoints live in `api/<domain>/core.go` and must export `Handler(w, r)`.
- Every `api/<domain>/core.go` must bootstrap an FX app per request using:
  - `lib/app/fx.CoreAppOptions`
  - `lib/router/fx.CoreRouterOptions` (unless you intentionally bypass the router)
- HTTP routing uses `chi`. Handlers must implement `lib/router.Handler`:
  - `RegisterRoute(r *chi.Mux)`
  - `Handle(w http.ResponseWriter, r *http.Request)`
- New HTTP endpoints:
  - Create handler(s) in `lib/app/<domain>/...`
  - Register them in `api/<domain>/core.go` via `router.AsRoute(<constructor>)`
- Vercel routing is defined in `vercel.json` rewrites. When adding a domain, update rewrites to route to `/api/<domain>/core`.
- Avoid cyclic imports between domains:
  - Put cross-domain interfaces in `lib/interfaces/<area>` packages.
  - Depend on interfaces, not concrete implementations; wire concrete implementations via FX in the caller’s `api/<domain>/core.go`.

## Infra conventions
- Config is via Viper: `config.NewViper` and `config.NewConfig`.
- Postgres uses SQLX via `db.NewSQLXPostgresDB` and must be closed via `fx.Lifecycle` hooks.
- Redis uses `cache.NewRedis`.
- Logging uses `lib/logs.NewLogger`.

## Inngest (optional)
- If included, keep the wrapper package generic: `lib/pkg/inngestclient`.
- Inngest Vercel handler lives in `api/inngest/core.go`.

## sqlc
- sqlc config lives at `sqlc.yaml`; queries under `db/query/`.
- Generated code should go into a stable package (documented in `sqlc.yaml`), and must not be edited by hand.
- Schema source lives at `supabase/schema.sql` (Supabase convention), and is pulled/exported manually before running `sqlc generate`.

## Responses
- Use `lib/pkg/render.ChiJSON` and `lib/pkg/render.ChiErr` as the default response helpers (unwrapped JSON).

## Output expectations for scaffolding tasks
When asked to “expose an existing package as a service”:
- Prefer minimal changes to the existing package; wrap it with constructors and FX providers.
- Create the minimal Vercel entrypoint + one example route (plus health if missing).
- Update `vercel.json` rewrites and add a short README snippet with run steps.

