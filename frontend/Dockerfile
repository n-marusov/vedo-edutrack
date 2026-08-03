# syntax=docker/dockerfile:1
# VEDO EduTrack SPA — nginx variant (SaaS/CDN-friendly).
#
# Build context: repository root (pnpm workspace lockfile lives at the root).
#   docker build -f frontend/Dockerfile .
#
# Runtime: nginxinc/nginx-unprivileged (non-root, UID 101) + alpine-slim
# (5.2 MB base vs 20 MB nginx:alpine — no docs/default html/debug tooling).
# Non-root matches the project security posture (backend: distroless nonroot;
# REQ-NFR-security.compliance.owasp-application-security). The unprivileged
# image binds 8080 (not 80) and keeps pid/temp paths under /tmp.
#
# Variant B from ADR-IMPL.PROCESS.development-tooling §8. The docker-compose
# dev stack does NOT use this image — it runs the Vite dev server directly
# (hot-reload, no Docker build needed). This image targets staging/SaaS.

# =============================================================================
# Stage 1: SPA build (Vite + pnpm)
# =============================================================================
#
# Задача: собрать production-бандл React SPA во frontend/dist.
# Кэширование: package.json/lockfile → pnpm install → исходники.
# Такой порядок максимизирует попадание в кэш при PR-билдах.

FROM node:24-alpine AS spa-build
WORKDIR /app

# Шаг 1. Файлы зависимостей (кешируются независимо от исходников).
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
COPY frontend/package.json ./frontend/package.json

# Шаг 2. Установка зависимостей (frozen-lockfile — строгая фиксация версий).
# --mount=type=cache: pnpm-store переживает между билдами (ускорение CI).
RUN --mount=type=cache,id=pnpm-store,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile

# Шаг 3. Исходники — самый часто меняющийся слой.
COPY --link frontend/ ./frontend/

# Шаг 4. Production-сборка Vite.
# VERSION — обязательный build-arg, прошивается в SPA-бандл.
ARG VERSION
RUN VITE_APP_VERSION=${VERSION} pnpm --filter @vedo-edutrack/frontend build

# =============================================================================
# Stage 2: nginx runtime (non-root, alpine-slim)
# =============================================================================
#
# Задача: раздавать статику SPA через nginx (non-root, минимальный образ).
# SaaS/CDN-путь (variant B): SPA отделена от бэкенда — кэшируется на edge,
# масштабируется независимо.
#
# Безопасность:
#   - nginxinc/nginx-unprivileged — non-root (UID 101), порт 8080 (не 80)
#   - alpine-slim — без docs/default-html/debug-тулинга (база 5.2 МБ)
#   - Динамический конфиг через envsubst-шаблоны — один билд, много окружений

FROM nginxinc/nginx-unprivileged:1.27-alpine-slim

# Версия передаётся через --build-arg, без дефолта — сборка без явной
# версии упадёт. Прошивается в OCI-лейбл и ENV APP_VERSION.
ARG VERSION
ARG SOURCE=https://github.com/vedo-edutrack/vedo-edutrack

LABEL org.opencontainers.image.title="vedo-edutrack-spa" \
      org.opencontainers.image.description="VEDO EduTrack SPA (nginx variant — non-root, alpine-slim)" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.licenses="MIT"

# Конфигурация nginx + шаблон для динамического config.js.
COPY --link frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --link frontend/nginx-config.js.template /etc/nginx/templates/config.js.template
COPY --link --from=spa-build /app/frontend/dist /usr/share/nginx/html

# entrypoint-скрипт 20-envsubst-on-templates.sh должен ПИСАТЬ
# сгенерированный /usr/share/nginx/html/config.js — делаем директорию
# доступной для non-root пользователя nginx (UID 101), затем возвращаемся.
USER root
RUN chown -R nginx:nginx /usr/share/nginx/html
USER nginx

# Динамический рантайм-конфиг (ADR-DES.INFRA.dynamic-config-injection):
# entrypoint генерирует /usr/share/nginx/html/config.js из этих env-переменных
# при старте контейнера — один билд, много окружений.
# Переопределяются в рантайме (compose -e / CI-секреты).
ENV NGINX_ENVSUBST_OUTPUT_DIR=/usr/share/nginx/html \
    APP_VERSION=${VERSION} \
    APP_ENV=production \
    API_BASE_URL=/api

EXPOSE 8080

# HEALTHCHECK: busybox wget уже есть в alpine; nginx-unprivileged — UID 101.
# 127.0.0.1 (не localhost): nginx слушает только IPv4, а busybox wget
# не умеет fallback с ::1 на 127.0.0.1 при IPv6-резолве localhost.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["wget", "-q", "-O", "/dev/null", "http://127.0.0.1:8080/"]
