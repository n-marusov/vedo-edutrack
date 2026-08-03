# syntax=docker/dockerfile:1
#
# =============================================================================
# VEDO EduTrack — единый production-образ
# =============================================================================
#
# Концепция (M0.3 embed unification):
#   ОДИН бинарник `vedo-edutrack` обслуживает и REST API, и SPA с одного порта.
#   SPA (React/Vite) встраивается в Go-бинарник через `//go:embed`, исключая
#   необходимость в отдельном веб-сервере или reverse-proxy для статики.
#   На выходе — минимальный distroless-образ (~20 МБ), готовый к деплою как
#   в SaaS-контуре, так и в Enterprise on-prem (242-ФЗ).
#
# Три стадии сборки (multi-stage):
#   Stage 1 (spa-build) — Node.js: pnpm install → Vite build → frontend/dist
#   Stage 2 (build)     — Go: go build со встроенным frontend/dist
#   Stage 3 (runtime)   — distroless static nonroot: только бинарник и HEALTHCHECK
#
# Версионирование:
#   VERSION передаётся снаружи через --build-arg, без дефолта — сборка без
#   явной версии упадёт. В CI версия берётся из git describe или GitHub ref;
#   в Makefile — VERSION ?= $(shell git describe --tags --always --dirty).
#   Значение прошивается в трёх местах:
#     • VITE_APP_VERSION  — env-переменная для Vite (попадает в SPA-бандл)
#     • ldflags -X Version — Go-константа config.Version (отдаётся в /healthz)
#     • OCI-лейбл version  — Docker-метаданные образа
#
# Multi-arch (buildx):
#   build-стадии работают на $BUILDPLATFORM (хост-архитектура билдера),
#   runtime — на $TARGETPLATFORM (целевая архитектура). GOARCH=$TARGETARCH
#   указывает компилятору Go целевую архитектуру. Поддерживаются linux/amd64
#   и linux/arm64.
#   Сборка: docker buildx build --platform linux/amd64,linux/arm64
#
# См. deploy/README.md (container strategy),
#     ADR-DES.INFRA.modular-monolith-approach,
#     ADR-DES.INFRA.docker-images-environments.

# =============================================================================
# Stage 1: SPA build (Vite + pnpm)
# =============================================================================
#
# Задача: собрать production-бандл React SPA во frontend/dist.
# Кэширование: слои упорядочены от редко меняющихся к часто меняющимся:
#   1. package.json / lockfile (меняется при добавлении зависимостей)
#   2. pnpm install (использует buildkit-кэш pnpm-store)
#   3. Исходники frontend/ (меняются чаще всего)
# Такой порядок позволяет переиспользовать слой pnpm install при каждом
# PR-билде, если lockfile не менялся.

FROM --platform=$BUILDPLATFORM node:24-alpine AS spa-build
WORKDIR /app

# Шаг 1. Файлы зависимостей (корень воркспейса + frontend).
# Копируются до исходников, чтобы слой установки кэшировался независимо.
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
COPY frontend/package.json ./frontend/package.json

# Шаг 2. Установка зависимостей (pnpm, frozen-lockfile — строгая фиксация версий).
# --mount=type=cache сохраняет pnpm-store между билдами — повторные билды
# без изменений в lockfile проходят установку мгновенно.
RUN --mount=type=cache,id=pnpm-store,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile

# Шаг 3. Исходники frontend/ — самый часто меняющийся слой.
COPY frontend/ ./frontend/

# Шаг 4. Production-сборка Vite.
# VERSION — обязательный build-arg, прошивается в SPA-бандл как VITE_APP_VERSION
# и доступна фронтенду через import.meta.env.VITE_APP_VERSION.
ARG VERSION
RUN VITE_APP_VERSION=${VERSION} pnpm --filter @vedo-edutrack/frontend build

# =============================================================================
# Stage 2: Go build (embedded SPA)
# =============================================================================
#
# Задача: скомпилировать Go-бинарник `vedo-edutrack` со встроенной SPA.
# Статическая линковка (CGO_ENABLED=0) — бинарник не зависит от системных
# библиотек и работает в distroless (где нет glibc).

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
WORKDIR /src

# TARGETARCH — целевая архитектура (amd64 / arm64), заполняется buildx.
ARG TARGETARCH
ARG VERSION

# Шаг 1. Скачивание модулей (кешируется при неизменном go.mod/go.sum).
COPY backend/go.mod ./go.mod
COPY backend/go.sum ./go.sum
RUN go mod download

# Шаг 2. Исходники бэкенда.
COPY backend/ ./

# Шаг 3. SPA dist — копируется из Stage 1 в путь, ожидаемый директивой
# //go:embed all:frontend/dist в backend/internal/platform/spa/embed.go.
COPY --from=spa-build /app/frontend/dist ./internal/platform/spa/frontend/dist

# Шаг 4. Компиляция Go-бинарника.
#   - CGO_ENABLED=0       — статическая линковка (distroless-совместимость)
#   - GOOS=linux          — целевая ОС
#   - GOARCH=${TARGETARCH} — целевая архитектура (из buildx)
#   - -trimpath            — убирает абсолютные пути из бинарника (воспроизводимость)
#   - -ldflags="-s -w"     — strip отладочной информации (‑3–5 МБ)
#   - -X config.Version    — прошивает версию в рантайм (отдаётся /healthz, логи)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-s -w -X vedo-edutrack/backend/internal/platform/config.Version=${VERSION}" \
    -o /out/vedo-edutrack \
    ./cmd/vedo-edutrack

# =============================================================================
# Stage 3: Runtime (distroless static, nonroot)
# =============================================================================
#
# Задача: минимальный runtime-образ без shell, пакетного менеджера и лишних
# утилит. Только бинарник и HEALTHCHECK. Атака-поверхность сведена к минимуму.
#
# Безопасность:
#   - nonroot-пользователь (UID 65532) — нет прав на запись в ФС
#   - Нет shell (нет /bin/sh) — невозможен RCE через shell-инъекции
#   - Нет пакетного менеджера — невозможно доустановить утилиты
#   - Подходит под требования 242-ФЗ (периметр, минимальная поверхность)
#
# HEALTHCHECK: в distroless нет curl/wget, поэтому проверка реализована
# через подкоманду `vedo-edutrack health` — бинарник сам делает HTTP GET /healthz
# и возвращает exit code 0 при 200 OK.

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# OCI-метаданные образа.
ARG VERSION
ARG SOURCE=https://github.com/vedo-edutrack/vedo-edutrack

LABEL org.opencontainers.image.title="vedo-edutrack" \
      org.opencontainers.image.description="VEDO EduTrack — educational trajectory service (API + embedded SPA, single binary)" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.source="${SOURCE}" \
      org.opencontainers.image.licenses="MIT"

# Перенос бинарника из build-стадии.
COPY --from=build /out/vedo-edutrack /usr/local/bin/vedo-edutrack

# Runtime-настройки.
USER nonroot:nonroot
EXPOSE 8080

# Liveness probe: бинарник сам делает HTTP GET /healthz и возвращает exit 0 при 200.
#   --interval=30s     — проверка каждые 30 секунд
#   --timeout=3s        — таймаут одного запроса
#   --start-period=5s   — задержка перед первой проверкой (даём бинарнику стартануть)
#   --retries=3         — три фейла подряд → контейнер unhealthy
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/usr/local/bin/vedo-edutrack", "health"]

# ENTRYPOINT — бинарник; CMD — подкоманда по умолчанию (запуск HTTP-сервера).
# Можно переопределить: docker run ... vedo-edutrack migrate
ENTRYPOINT ["/usr/local/bin/vedo-edutrack"]
CMD ["server"]
