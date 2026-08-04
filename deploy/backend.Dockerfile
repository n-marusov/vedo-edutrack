# syntax=docker/dockerfile:1
#
# =============================================================================
# VEDO EduTrack — dev-образ с air hot-reload
# =============================================================================
#
# Dev-образ для локальной разработки (ADR-IMPL.PROCESS.development-tooling).
# Расширяет golang:1.26-alpine с предустановленным air v1.67.4 для горячей
# перезагрузки Go-кода. Не имеет ENTRYPOINT/CMD — команда задаётся в
# docker-compose.yml (deploy/docker-compose.yml).
#
# Сборка:
#   docker build -f deploy/backend.Dockerfile -t vedo-edutrack-backend:dev .
#
# Использование:
#   docker compose -f deploy/docker-compose.yml up backend
#
# См. deploy/README.md (container strategy),
#     ADR-DES.INFRA.docker-images-environments.

FROM golang:1.26-alpine

ARG SOURCE=https://github.com/vedo-edutrack/vedo-edutrack

LABEL org.opencontainers.image.title="vedo-edutrack-backend" \
      org.opencontainers.image.description="VEDO EduTrack — dev backend with air hot-reload" \
      org.opencontainers.image.version="dev" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.licenses="MIT"

# Pre-install air v1.67.4 for hot-reload (no ENTRYPOINT/CMD — set by compose).
RUN go install github.com/air-verse/air@v1.67.4

# Pre-download Go module dependencies at build time so that air does not need
# DNS at runtime (common issue on Windows Docker Desktop with VPN/corporate
# proxy). WORKDIR must be set to a non-GOPATH directory (default GOPATH is /go)
# — Go ignores go.mod files inside GOPATH.
WORKDIR /build-cache
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# WORKDIR for the runtime volume mount (backend/ -> /app). The module cache
# lives in /root/go/pkg/mod (outside the mount) and is inherited from the
# pre-download above.
WORKDIR /app