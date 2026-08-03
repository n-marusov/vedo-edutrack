# VEDO EduTrack — Container & Deployment Strategy

> M0.2 exit criterion: *container strategy documented and aligned with selected stack*.
> Sources: `ADR-DES.INFRA.modular-monolith-approach`, `ADR-IMPL.PROCESS.development-tooling` §7–8.

## Architecture Overview

VEDO EduTrack is a **modular monolith**: one Go binary (`vedo-edutrack`) that hosts
all bounded contexts behind a single HTTP surface. Deployment follows the same shape:

- **One deployable artifact** — the backend binary. From M0.3 the SPA is **embedded**
  into it via `//go:embed` (single artifact for on-prem Enterprise simplicity).
- **`docker-compose`** drives dev and staging (all 9 services, one command).
- **Traefik** is the edge reverse-proxy: TLS, routing, rate limiting, circuit breaker.

## Image Strategy

| Image | Dockerfile | Base → Runtime | Size | Purpose |
|-------|------------|----------------|------|---------|
| Unified backend (API + embedded SPA) | `backend/Dockerfile` | node SPA build → `golang:1.26-alpine` → `gcr.io/distroless/static-debian12:nonroot` | ~20 MB | Single binary serving API and SPA from one port (M0.3); multi-arch amd64+arm64 via buildx |
| SPA (nginx) | `frontend/Dockerfile` | node build → `nginxinc/nginx-unprivileged:1.27-alpine-slim` | ~7 MB | SaaS/CDN-friendly variant B, non-root (UID 101), port 8080 |
| Frontend dev | — (no build) | `node:24-alpine` + Vite dev server | — | Hot-reload HMR, no Docker build needed |

Build contexts: `backend/Dockerfile` → **repo root** (the pnpm workspace lockfile and
`frontend/` live there; the backend image builds the SPA in a node stage and embeds
the dist). `frontend/Dockerfile.embed` was removed in M0.3 — the embed mechanism
moved into the backend binary. `.dockerignore` files keep contexts lean.

## Dev Environment

```bash
cp .env.dev.example .env.dev   # optional: copy dev env template (root)
docker compose -f deploy/docker-compose.yml up -d --wait   # or: make up
```

Starts 9 services: `backend` (air hot-reload), `frontend` (Vite HMR), `postgres`,
`otel-collector`, `prometheus`, `loki`, `tempo`, `grafana`, `traefik`.
Networks: `edutrack-net` (services) + `edutrack-public` (Traefik ingress).
Volumes: `postgres_data`, `grafana_data`, `loki_data`, `tempo_data`,
`frontend_node_modules` (container-owned, avoids host/container binary clashes).

Ports: backend `8080`, frontend `5173`, postgres `5432`, Grafana `3000`,
Prometheus `9090`, Loki `3100`, Tempo `3200`, OTLP `4317/4318`, Traefik `80/443`
(dashboard `8082`, dev only). Defaults are overridable via **`.env.dev`** (root,
`ENV_FILE` in the Makefile): `make up` passes it to compose via
`--env-file .env.dev` (template: `.env.dev.example`). When the file is absent,
compose falls back to the process environment and per-variable defaults.
Direct `docker compose` calls interpolate from `deploy/.env` (template:
`deploy/.env.example`).

## SaaS Deployment

- **Edge**: Traefik (TLS via Let's Encrypt, rate limiting, security headers).
- **Routing**: `api.edutrack.localhost → backend:8080`, `edutrack.localhost → frontend`.
- **Blue-green**: post-MVP via Traefik weighted services (two backend replicas,
  weighted round-robin, canary + kill switch ≤ 5 min).
- **CI**: GitHub Actions builds and pushes images to GHCR (`:latest`, `:sha-<commit>`);
  deployment via SSH + compose (see `.github/workflows/`).

## Enterprise On-Prem

- **Single Go binary** (embedded SPA) + **PostgreSQL** only.
- Redis optional (distributed rate limiting, read-model cache, token blacklist).
- No orchestration required on MVP; K8s/Helm post-MVP for customers with their own K8s.
- Runtime: distroless, nonroot, no shell — minimal attack surface (242-ФЗ perimeter).

## Observability Stack

OTel SDKs (Go + Web) → OTLP → **OTel Collector** → Prometheus (metrics) + Loki (logs)
+ Tempo (traces) + Grafana (dashboards, provisioned as-code in
`deploy/observability/`). Sampling: 100% error traces + 10% successful.
PII redaction (152-ФЗ) happens in the collector pipeline.

## CI Integration

`make ci` mirrors the GitHub Actions pipeline (`bash deploy/ci/run-gates.sh --tier delivery --trigger ci`).
Images are built and pushed by CI on `main`; the delivery gate checks the distroless
image stays ≤ 20 MB.

## Security

- nonroot runtime user (backend: distroless `nonroot`; SPA nginx: `nginxinc/nginx-unprivileged`, UID 101)
- CSP/HSTS/security headers + rate limiting + circuit breaker at the Traefik edge
- `gitleaks` secret scan + `gosec` SAST + `pnpm audit` + syft SBOM in CI
- no secrets in images (`.dockerignore` excludes `.env*`); secrets come from env/CI

## Deployment Contours

| Contour | Stack | Notes |
|---------|-------|-------|
| **Community (SaaS)** | Traefik + compose + full OTel stack | Blue-green, dashboards, alerting |
| **Enterprise (on-prem)** | single binary + PostgreSQL (+ optional Redis) | Minimal surface, 242-ФЗ, no external telemetry |
