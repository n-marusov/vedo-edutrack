# syntax=docker/dockerfile:1
# VEDO EduTrack unified image — single vedo-edutrack binary serving both the
# API and the embedded SPA from one port (M0.3 embed unification).
#
# Multi-stage build:
#   Stage 1: node:24-alpine  — build the SPA (Vite/pnpm)
#   Stage 2: golang:1.26-alpine — go build with the SPA dist embedded
#   Stage 3: distroless static (nonroot) — runtime
#
# Multi-arch: build stages run on $BUILDPLATFORM, runtime on $TARGETPLATFORM;
# GOARCH=$TARGETARCH for the Go compiler (buildx --platform linux/amd64,linux/arm64).
# See deploy/README.md container strategy and ADR-DES.INFRA.modular-monolith-approach.

# --- Stage 1: SPA build (Vite) ---
FROM --platform=$BUILDPLATFORM node:24-alpine AS spa-build
WORKDIR /app
# Workspace root files before frontend sources (cache-friendly layer).
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
COPY frontend/package.json ./frontend/package.json
RUN --mount=type=cache,id=pnpm-store,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile
COPY frontend/ ./frontend/
ARG VERSION=0.1.0
RUN VITE_APP_VERSION=${VERSION} pnpm --filter @vedo-edutrack/frontend build

# --- Stage 2: Go build (embedded SPA) ---
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
WORKDIR /src
ARG TARGETARCH
ARG VERSION=0.1.0

# Cache module downloads.
COPY backend/go.mod ./go.mod
COPY backend/go.sum ./go.sum
RUN go mod download

# Copy sources; then place the SPA dist at the embed path.
COPY backend/ ./
COPY --from=spa-build /app/frontend/dist ./internal/platform/spa/frontend/dist

# Build the single binary (cobra CLI tree — ADR-DES.API.cli-interface).
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-s -w -X vedo-edutrack/backend/internal/platform/config.Version=${VERSION}" \
    -o /out/vedo-edutrack \
    ./cmd/vedo-edutrack

# --- Stage 3: Runtime (distroless static, nonroot) ---
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
ARG VERSION=0.1.0
ARG SOURCE=https://github.com/vedo-edutrack/vedo-edutrack

LABEL org.opencontainers.image.title="vedo-edutrack" \
      org.opencontainers.image.description="VEDO EduTrack — educational trajectory service (API + embedded SPA, single binary)" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.licenses="MIT"

COPY --from=build /out/vedo-edutrack /usr/local/bin/vedo-edutrack

USER nonroot:nonroot
EXPOSE 8080

# Liveness probe via the binary's own `health` subcommand (distroless runtime
# has no shell/curl/wget): it GETs /healthz over HTTP and exits 0 on 200.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/usr/local/bin/vedo-edutrack", "health"]

ENTRYPOINT ["/usr/local/bin/vedo-edutrack"]
CMD ["server"]
