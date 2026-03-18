# AGENTS.md (Codex)

This repository is the canonical source for the minimal Vercel Go + FX + chi scaffold.

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

## Responses
- Use `lib/pkg/render.ChiJSON` and `lib/pkg/render.ChiErr` as the default response helpers (unwrapped JSON).

## Scope guardrails
- Keep this repo free of business logic and third-party product integrations.
- Do not add long-running server entrypoints such as `cmd/app/main.go`.
- Do not add Inngest, sqlc, Supabase, Turso, or NotebookLM bootstrap to the template baseline.
- Prefer a single example domain (`health`) and keep all other behavior as documentation, not prebuilt code.
