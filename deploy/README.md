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
| Backend | `backend/Dockerfile` | `golang:1.26-alpine` → `gcr.io/distroless/static-debian12:nonroot` | ~15 MB | Production runtime (healthcheck, nonroot, OCI labels) |
| SPA (embed) | `frontend/Dockerfile.embed` | node build → Go embed server → distroless | small | Proves the on-prem embed path (moves into the backend binary in M0.3) |
| SPA (nginx) | `frontend/Dockerfile.nginx` | node build → `nginx:1.27-alpine` | ~50 MB | SaaS/CDN-friendly variant B |
| Frontend dev | — (no build) | `node:24-alpine` + Vite dev server | — | Hot-reload HMR, no Docker build needed |

Build contexts: `backend/Dockerfile` → `backend/`; frontend Dockerfiles → **repo root**
(the pnpm workspace lockfile lives at the root). `.dockerignore` files keep contexts lean.

## Dev Environment

```bash
docker compose -f deploy/docker-compose.yml up -d --wait   # or: make up
```

Starts 9 services: `backend` (air hot-reload), `frontend` (Vite HMR), `postgres`,
`otel-collector`, `prometheus`, `loki`, `tempo`, `grafana`, `traefik`.
Networks: `edutrack-net` (services) + `edutrack-public` (Traefik ingress).
Volumes: `postgres_data`, `grafana_data`, `loki_data`, `tempo_data`,
`frontend_node_modules` (container-owned, avoids host/container binary clashes).

Ports: backend `8080`, frontend `5173`, postgres `5432`, Grafana `3000`,
Prometheus `9090`, Loki `3100`, Tempo `3200`, OTLP `4317/4318`, Traefik `80/443`
(dashboard `8082`, dev only). Defaults are overridable via `deploy/.env`
(template: `deploy/.env.example`).

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

- nonroot runtime user (`USER nonroot:nonroot`), distroless (no shell/packages)
- CSP/HSTS/security headers + rate limiting + circuit breaker at the Traefik edge
- `gitleaks` secret scan + `gosec` SAST + `pnpm audit` + syft SBOM in CI
- no secrets in images (`.dockerignore` excludes `.env*`); secrets come from env/CI

## Deployment Contours

| Contour | Stack | Notes |
|---------|-------|-------|
| **Community (SaaS)** | Traefik + compose + full OTel stack | Blue-green, dashboards, alerting |
| **Enterprise (on-prem)** | single binary + PostgreSQL (+ optional Redis) | Minimal surface, 242-ФЗ, no external telemetry |
