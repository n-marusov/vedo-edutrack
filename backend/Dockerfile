# syntax=docker/dockerfile:1
# VEDO EduTrack backend image — single vedo-edutrack binary (modular monolith).
#
# Multi-stage build: Go build (alpine) -> distroless static runtime (nonroot).
# See ADR-IMPL.PROCESS.development-tooling §8 and
# ADR-DES.INFRA.modular-monolith-approach (single artifact, on-prem friendly).
#
# NOTE: go.mod pins `go 1.26.3` — the plan's golang:1.24-alpine is superseded
# by the module's toolchain version (T1 wrote go.mod first).

# --- Build stage ---
FROM golang:1.26-alpine AS build
WORKDIR /src

# Cache module downloads. go.sum appears when the first dependency lands (M0.3).
COPY go.mod ./go.mod
RUN go mod download

# Build the single binary (cobra CLI tree — ADR-DES.API.cli-interface).
COPY . .
ARG VERSION=0.1.0
# Dynamic build-version injection (ADR-DES.INFRA.dynamic-config-injection):
# the version is baked via ldflags and reported by `vedo-edutrack version`.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X vedo-edutrack/backend/internal/platform/config.Version=${VERSION}" \
    -o /out/vedo-edutrack \
    ./cmd/vedo-edutrack

# --- Runtime stage (distroless static, nonroot) ---
FROM gcr.io/distroless/static-debian12:nonroot AS runtime
ARG VERSION=0.1.0
ARG SOURCE=https://github.com/vedo-edutrack/vedo-edutrack

LABEL org.opencontainers.image.title="vedo-edutrack" \
      org.opencontainers.image.description="VEDO EduTrack — educational trajectory service (modular monolith, single binary)" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.licenses="MIT"

COPY --from=build /out/vedo-edutrack /usr/local/bin/vedo-edutrack

USER nonroot:nonroot
EXPOSE 8080

# Liveness probe via the binary's own `health` subcommand (distroless runtime
# has no shell/curl/wget): it GETs /healthz over HTTP and exits 0 on 200.
# Liveness/readiness split: /healthz + /readyz (PostgreSQL) — development-tooling §8.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/usr/local/bin/vedo-edutrack", "health"]

ENTRYPOINT ["/usr/local/bin/vedo-edutrack"]
CMD ["server"]
