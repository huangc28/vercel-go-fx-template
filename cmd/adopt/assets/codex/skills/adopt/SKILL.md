---
name: adopt
description: Non-destructive adoption of the Vercel Go + Uber FX + chi architecture in an existing Go repository. Use when asked to “/adopt”, “initialize this repo to the boilerplate structure”, “add api/<domain>/core.go Vercel entrypoints”, “introduce FX wiring + chi router”, or “migrate existing routes into the new per-domain serverless architecture without deleting old code”.
---

# Adopt

## Overview

Add the “Vercel serverless entrypoints + per-request FX + chi router” structure to an existing Go repo without breaking or deleting the current implementation. Start by adding a parallel `/health` endpoint, then optionally bridge one existing domain behind a new `vercel.json` rewrite.

## Workflow

### 0) Hard rules

- Do not delete or rename existing code unless explicitly asked.
- Prefer additive changes and adapters/wrappers.
- Keep baseline `go test ./...` and `go build ./...` working without requiring live Postgres/Redis (treat infra as optional/stubbed).
- Keep the architecture conventions:
  - Vercel entrypoints at `api/<domain>/core.go` exporting `Handler(w, r)`
  - Request bootstraps an `fx.App` and serves an `http.Handler` (chi mux or compatible)
  - Handlers implement `lib/router.Handler` and are registered via `router.AsRoute(...)`

### 1) Preflight discovery (no edits yet)

Confirm:
- `go.mod` exists; record the module path.
- Current entrypoint/router locations (search for `http.ListenAndServe`, `chi.NewRouter`, `gin.New`, `echo.New`, `mux.NewRouter`, etc.).
- Whether `vercel.json` already exists and what rewrites/redirects it contains.
- Whether `api/` already exists (avoid collisions).

If you cannot infer the “first domain to bridge”, ask at most 3 questions:
1) Which route prefix should be migrated first (e.g. `/v1/foo/*`)?
2) Where is the current router/handler that serves that prefix?
3) Is the repo deployed on Vercel today (i.e. do we edit `vercel.json` now)?

### 2) Add the core scaffold (parallel, minimal)

Add (or adapt to match existing packages):
- `lib/router/core.go`: `Handler` interface + `AsRoute(...)` + `NewMux(...)` that registers all grouped handlers into a `chi.Mux`.
- `lib/router/fx/options.go`: `CoreRouterOptions` providing the mux.
- `lib/pkg/render/render.go`: `ChiJSON` and `ChiErr` helpers (unwrapped JSON).

If the repo doesn’t already provide a zap logger/config via constructors, add minimal providers so FX can build:
- Logger: provide `*zap.Logger` (ok to use `zap.NewNop()` as a baseline)
- Config: optional; only required if your handlers depend on it

### 3) Prove wiring with `/health`

Add:
- `lib/app/health/...` handler implementing `lib/router.Handler` and returning `{ "ok": true }`.
- `api/health/core.go` Vercel entrypoint that:
  - builds `fx.New(...)` with core app options + router options + `router.AsRoute(health.NewHandler)`
  - `fx.Populate(&mux)` and `mux.ServeHTTP(w, r)`

Update `vercel.json`:
- Add rewrite `/health` -> `/api/health/core` (preserve existing rewrites; merge conservatively).

### 4) Bridge one existing domain (optional but recommended)

Goal: route requests to `/api/<domain>/core` while still using the existing domain implementation internally.

Add:
- `lib/app/<domain>/legacy.go` implementing `lib/router.Handler`.
  - `RegisterRoute`: register only the target prefix(es) you’re bridging.
  - `Handle`: delegate to the existing handler/router (call into your existing code; do not duplicate business logic).

Add:
- `api/<domain>/core.go` entrypoint wiring FX + router + `router.AsRoute(legacy.NewHandler)`.

Update `vercel.json`:
- Add one rewrite for the prefix (e.g. `/v1/<domain>/(.*)` -> `/api/<domain>/core`).

### 5) Validation

Run:
- `gofmt` on changed files
- `go mod tidy`
- `go test ./...`
- `go build ./...`

If infra dependencies require live services, gate them behind env vars and return `nil` clients when not configured, so tests/builds don’t fail.

### 6) Handoff notes

Leave the repo in a state where:
- Old entrypoints still work (unless rewrites were switched).
- New entrypoints work (`/health` at minimum).
- There’s a clear “bridge handler” pattern to migrate further domains route-by-route.

