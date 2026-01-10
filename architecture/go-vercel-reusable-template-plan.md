# Vercel Go Structure Template Plan (Based on Current Repo)

This document captures the current Go + Vercel architecture in this repo and defines a conservative, reusable boilerplate/template plan that keeps the same structure and patterns (per-domain `api/<domain>/core.go`, FX wiring, Chi router, Vercel rewrites).

## 1) Current Project Structure (Inventory)

Key top-level folders in this repo:

```
api/                   # Vercel serverless function entrypoints (per domain)
cmd/                   # Local/utility executables (non-Vercel)
config/                # Viper-based typed config + defaults
cache/                 # Redis client init
db/                    # SQLX Postgres init + FX lifecycle hooks
lib/                   # Core app modules, routes, services, pkg integrations
http/                  # Local HTTP client test files (not runtime)
architecture/          # Architecture + plans
vercel.json            # Vercel function config + rewrites
sqlc.yaml              # sqlc codegen configuration (queries + schema -> Go)
```

### 1.1 Vercel entrypoints (serverless functions)

- `api/<domain>/core.go` is the unit of deployment.
- Each `core.go` exports `Handler(w http.ResponseWriter, r *http.Request)`.
- Each request bootstraps an FX app (per-request DI).
- Shared dependencies are included via:
  - `lib/app/fx/CoreAppOptions` (config, logger, DB, Redis)
  - `lib/router/fx/CoreRouterOptions` (Chi router + grouped handlers)
- Vercel routing is configured via `vercel.json` `rewrites` to map public routes to `api/<domain>/core`.

### 1.2 Dependency injection + router pattern

- `lib/router/core.go` defines:
  - `router.Handler` interface (`RegisterRoute(*chi.Mux)`, `Handle(w, r)`)
  - `router.AsRoute(...)` which registers handlers into the `group:"handlers"` FX group
  - shared middleware setup + handler registration loop
- `lib/router/fx/options.go` provides the router constructor with FX group injection.

### 1.3 Config + infra modules

- `config/config.go` uses Viper + defaults and exposes `NewViper()` and `NewConfig(*viper.Viper)`.
- `db/db.go` provides a DB module with close hooks via `fx.Lifecycle`.
- `db/tx.go` provides a small transaction helper (`db.Tx`) for SQLX.
- `cache/redis.go` provides Redis init (serverless-friendly client settings).
- `lib/logs/*` provides the logger.
- `lib/interfaces/*` should expose cross-domain interfaces (and only minimal shared types) to avoid cyclic imports.

### 1.4 sqlc (SQL-to-Go codegen)

- `sqlc.yaml` configures `sqlc generate`:
  - queries: `db/query/*.sql`
  - schema: `supabase/schema.sql` (this is the default convention for the boilerplate; users pull/export it from Supabase)
  - output: choose a stable package (example: `lib/app/models`) if/when you generate code
- This is useful in the boilerplate as a standard, typed DB access layer for each new service.

### 1.5 Default response helpers (Chi)

- `lib/pkg/render` should provide the default response helpers for handlers:
  - `render.ChiJSON(...)` for unwrapped success responses
  - `render.ChiErr(...)` for unwrapped error responses

### 1.6 Non-goal: non-Vercel entrypoints

- This boilerplate is tailored for Go + Vercel serverless functions only.
- Do not include a local/monolith entrypoint like `cmd/app/main.go` in this template.

## 2) Template Goal (What we want to extract)

Create a reusable template that mirrors this architecture:

- Keep `api/<domain>/core.go` as Vercel entrypoints.
- Keep shared DI modules (`lib/app/fx`, `lib/router/fx`) and router conventions (`lib/router`).
- Keep infra packages `config/`, `db/`, `cache/`, and logging.
- Keep `vercel.json` rewrite style the same.
- Provide a minimal example domain (`health`) and optional Inngest support.

Non-goals:
- No large refactors (no `internal/` migration, no renaming core folders).
- No peasydeal-specific business/domain code in the boilerplate.

## 3) Template Structure (Conservative, same as current)

```
.
├── api/
│   ├── health/
│   │   └── core.go
│   ├── inngest/                # optional
│   │   └── core.go
│   └── <domain>/
│       └── core.go
├── cache/
│   └── redis.go
├── config/
│   └── config.go
├── db/
│   └── db.go
│   └── tx.go                   # SQLX transaction helper
│   └── query/                  # sqlc queries
│       └── *.sql
├── lib/
│   ├── app/
│   │   ├── fx/
│   │   │   └── core.go
│   │   ├── health/
│   │   │   └── core.go          # example handler
│   │   └── <domain>/            # domain-specific handlers
│   ├── logs/
│   │   └── logs.go
│   ├── pkg/
│   │   └── inngestclient/       # optional, Inngest client wrapper (neutral name)
│   │   └── render/              # default response helpers (unwrapped JSON)
│   ├── interfaces/              # cross-domain interfaces to avoid cyclic imports
│   │   └── README.md
│   └── router/
│       ├── core.go
│       └── fx/
│           └── options.go
├── vercel.json
```

## 4) Template Build Plan (Extraction Rules)

### Step 1 — Identify what’s boilerplate vs. product-specific

Boilerplate:
- Vercel entrypoints (the `api/<domain>/core.go` pattern)
- FX wiring (CoreApp + CoreRouter options)
- Router + handler registration patterns
- Infra (config/db/cache/logging)
- Standard JSON responses
- Minimal example handler + route
- `vercel.json` rewrites convention

Product-specific:
- Any domain handlers tied to business logic
- Product-specific config defaults and helper packages
- Anything that can’t compile/run without secrets or vendor-specific services

### Step 2 — Trim to the minimal reusable template

- Keep only a minimal example domain (`health`) with a simple route.
- Keep the Inngest entrypoint optional (include with minimal setup, and document how to remove it if unused).
- Remove product-specific handlers, DAOs, and any domain-specific services.
- Remove peasydeal-specific config helpers (e.g. `bcc_list.go`) and any email-list defaults.
- Keep sqlc “wired in” as part of the template (config + example query), but do not require generated code for the template to compile out-of-the-box.
- Include default response helpers (`lib/pkg/render`) and use them in example handlers.

### Step 3 — Parameterize names and paths

- The template should compile as its own Go module (with its own module path).
- When instantiating, rewrite module/import paths to the new service module path (see `gonew` section below).

### Step 4 — Prepare Vercel rewrites

- Provide a minimal `vercel.json` with:
  - `/health` -> `/api/health/core`
  - optional `/api/inngest` -> `/api/inngest/core`
- Document how to add new domain rewrites (pattern: `/v2/<domain>/*` -> `/api/<domain>/core`).

### Step 5 — Documentation

- Template README should include:
  - How to instantiate a new service
  - How to add a new domain (new `api/<domain>/core.go` + handler in `lib/app/<domain>`)
  - Required env vars (from `config/config.go`)
  - How to pull `supabase/schema.sql` and then run sqlc (`sqlc generate`)
 - Add a minimal `Makefile` for common dev actions (see below).

### Step 6 — Validation checklist

- `go test ./...` (should not require a live DB/Redis for the template baseline)
- `go build ./...`
- `vercel dev` smoke test: ensure `/health` works
- If Inngest included: ensure `/api/inngest` responds (or registers functions) with correct env
- sqlc: document as a manual post-init step (requires `supabase/schema.sql`), not part of baseline validation

## 5) Concrete Actions (Template Build Checklist)

- [ ] Create a dedicated template repo (copy only boilerplate).
- [ ] Keep exact folder structure for `api/`, `lib/`, `config/`, `db/`, `cache/`.
- [ ] Keep FX modules as-is (`CoreAppOptions`, `CoreRouterOptions`, `AsRoute`).
- [ ] Add minimal `health` handler + route to validate wiring.
- [ ] Decide on Inngest inclusion (default: include with TODO note) and keep the wrapper name generic (e.g. `lib/pkg/inngestclient`).
- [ ] Trim config to essentials; exclude peasydeal-specific helpers like `bcc_list.go`.
- [ ] Include sqlc: add `sqlc.yaml`, minimal `db/query/*.sql`, and a template-owned `db/schema.sql`.
- [ ] Add a minimal `Makefile` with `help`, `start/vercel`, `start/inngest`.
- [ ] Add `lib/interfaces/` for cross-domain interfaces (to prevent cyclic imports).
- [ ] Add default response helpers in `lib/pkg/render` and standardize example handlers on them.
- [ ] Ensure `vercel.json` rewrites are minimal and clear.
- [ ] Add README and Codex-focused AI instruction file (`AGENTS.md`).

## 6) Implementation TODO List (Copy/Paste for Execution)

Use this as the step-by-step TODO list to implement the boilerplate template repo.

### 6.1 Create template repo
- [ ] Create a new repo: `github.com/<org>/vercel-go-service-template`
- [ ] Initialize `go.mod` with the template module path
- [ ] Copy in only the boilerplate directories/files (see structure below)
- [ ] Ensure `go build ./...` works without any secrets or DB connections

### 6.2 Vercel entrypoints
- [ ] Add `api/health/core.go` Vercel handler (FX bootstrap + router)
- [ ] Add minimal health route handler in `lib/app/health` that returns JSON via `lib/pkg/render`
- [ ] Configure `vercel.json` rewrites for `/health` -> `/api/health/core`
- [ ] (Optional) Add `api/inngest/core.go` and `/api/inngest` rewrite if Inngest is enabled by default

### 6.3 Core DI + router
- [ ] Add `lib/app/fx/core.go` (`CoreAppOptions`) for logger/config/db/redis
- [ ] Add `lib/router/core.go` (chi router + middleware + handler group registration)
- [ ] Add `lib/router/fx/options.go` (`CoreRouterOptions`) to construct the router from grouped handlers

### 6.4 Infra packages
- [ ] Add `config/config.go` (Viper defaults + typed config; remove peasydeal-only fields)
- [ ] Add `db/db.go` (SQLX Postgres module + lifecycle close; serverless-friendly pool settings)
- [ ] Add `db/tx.go` (`db.Tx` transaction helper)
- [ ] Add `cache/redis.go` (Redis client init with serverless-friendly options)
- [ ] Add `lib/logs/*` (logger setup)

### 6.5 Default response helpers
- [ ] Add `lib/pkg/render/render.go` (unwrapped `ChiJSON`/`ChiErr`)
- [ ] Ensure example handlers use `render.ChiJSON`/`render.ChiErr` instead of ad-hoc JSON writes

### 6.6 Cross-domain interfaces
- [ ] Add `lib/interfaces/README.md`
- [ ] Add at least one example subpackage (optional) to demonstrate pattern, e.g. `lib/interfaces/example`

### 6.7 sqlc (manual generation)
- [ ] Add `sqlc.yaml` wired to:
  - `db/query/*.sql`
  - schema: `supabase/schema.sql`
  - output package/path (choose and document)
- [ ] Add `db/query/example.sql` (a minimal example query)
- [ ] Add `supabase/` directory (empty except `.gitkeep`) and document that `supabase/schema.sql` is pulled manually before running `sqlc generate`
- [ ] Document the required steps to pull/export `supabase/schema.sql` and run:
  - `go install github.com/sqlc-dev/sqlc/cmd/sqlc@<version>`
  - `sqlc generate`

### 6.8 Dev ergonomics
- [ ] Add minimal `Makefile` (`help`, `start/vercel`, `start/inngest`)
- [ ] Add `README.md` with:
  - `gonew` instantiation commands
  - how to add a domain (new `api/<domain>/core.go` + `lib/app/<domain>` handler)
  - required env vars (and where they live)
  - sqlc manual generation steps
- [ ] Add root `AGENTS.md` (Codex rules from this plan)

### 6.9 Safety + cleanliness
- [ ] Confirm no secrets/keys/dumps are included (no `.env*`, no DB dumps, no service account JSON)
- [ ] Confirm `go test ./...` does not require a live Postgres/Redis for baseline
- [ ] Smoke test with `vercel dev` and verify `/health` responds

## 7) Decisions (Before Extraction)

Decisions:
- Vercel-only: do not include `cmd/app/main.go`.
- sqlc schema convention: use `supabase/schema.sql`.

## 8) Implement Using `gonew`

`gonew` (from `golang.org/x/tools/cmd/gonew`) can instantiate a new project from a template module and rewrite the Go module path/imports for you. This is a good fit for this boilerplate because the main “template substitution” you need is module path rewriting.

### 7.1 Create and publish the template as a Go module

- Create a new repo for the template, for example: `github.com/<org>/vercel-go-service-template`.
- Set `go.mod` to that module path (example: `module github.com/<org>/vercel-go-service-template`).
- Ensure all internal imports in the template use the same module path.
- Keep the repo “template-clean”:
  - Do not include secrets, dumps, credentials, or private keys.
  - Keep only minimal example domains (ex: `health`) + optional Inngest.
  - Keep `vercel.json` and docs that apply broadly.
  - Keep `sqlc.yaml` + minimal schema/queries for codegen.

### 7.2 Install `gonew`

```
go install golang.org/x/tools/cmd/gonew@latest
```

### 7.3 Instantiate a new service

Run in the parent directory where you want the new folder created:

```
gonew github.com/<org>/vercel-go-service-template github.com/<org>/<new-service>
cd <new-service>
```

### 7.4 Private template notes (if applicable)

If the template repo is private:

- Set `GOPRIVATE=github.com/<org>/*` (or your git host/org pattern).
- Ensure your git credentials allow `go` to fetch the module (SSH keys or HTTPS tokens).

### 7.5 After `gonew` (manual edits you still do)

- Update `vercel.json` rewrites for the domains you add under `api/<domain>/core.go`.
- Set Vercel env vars to satisfy `config/config.go`.
- Add/modify sqlc queries under `db/query/` and update schema path in `sqlc.yaml` as needed.
- Optionally search/replace any remaining template repo name in README/docs (non-Go files).

### 7.6 Pinning `sqlc` (recommended)

`gonew` rewrites module paths, but it doesn’t install external binaries. For a consistent dev experience, the template should pin a `sqlc` version and document one command to install it, for example:

- `go install github.com/sqlc-dev/sqlc/cmd/sqlc@<version>`

Optionally, add a `tools.go` (with `//go:build tools`) to keep the tool dependency recorded in `go.mod`, and add Make targets like:

- `make tools` (installs `sqlc`)
- `make gen/sqlc` (runs `sqlc generate`)

## 9) Minimal Makefile (Template)

The boilerplate should include a small `Makefile` that mirrors the commands you already use, but without the repo-specific Docker/deploy targets.

Recommended contents:

```makefile
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## start/vercel: start vercel dev server
.PHONY: start/vercel
start/vercel:
	@vercel dev --debug --listen $(APP_PORT)

## start/inngest: start the inngest dev server
.PHONY: start/inngest
start/inngest:
	PORT=3011 npx inngest-cli@latest dev \
		--no-discovery \
		--poll-interval 10000 \
		-u http://localhost:3010/api/inngest
```

## 10) Codex Instructions (`AGENTS.md`)

This boilerplate is primarily targeting developers using Codex. Include an `AGENTS.md` at the template repo root and treat it as the single source of truth for “where things go” and “how to wire them”.

Notes:
- `.github/copilot-instructions.md` is not required for this boilerplate; omit it unless you explicitly want Copilot-specific guidance.

Recommended `AGENTS.md` contents (template-quality, keep it short and strict):

```md
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
- Postgres uses SQLX via `db.SQLXPostgresDBModule` and must be closed via `fx.Lifecycle` hooks.
- Redis uses `cache.NewRedis`.
- Logging uses `lib/logs.NewLogger`.

## Inngest (optional)
- If included, keep the wrapper package generic: `lib/pkg/inngestclient`.
- Inngest Vercel handler lives in `api/inngest/core.go`.

## sqlc
- sqlc config lives at `sqlc.yaml`; queries under `db/query/`.
- Generated code should go into a stable package (documented in `sqlc.yaml`), and must not be edited by hand.
- Schema source lives at `supabase/schema.sql` (Supabase convention), and is pulled/exported manually from Supabase before running `sqlc generate`.

## Responses
- Use `lib/pkg/render.ChiJSON` and `lib/pkg/render.ChiErr` as the default response helpers (unwrapped JSON).

## Output expectations for scaffolding tasks
When asked to “expose an existing package as a service”:
- Prefer minimal changes to the existing package; wrap it with constructors and FX providers.
- Create the minimal Vercel entrypoint + one example route (plus health if missing).
---
- Create the minimal Vercel entrypoint + one example route (plus health if missing).
- Update `vercel.json` rewrites and add a short README snippet with run steps.
```

