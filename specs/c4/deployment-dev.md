# C4 Level 4: Deployment Diagram — Dev Environment (local)

> Уровень 4 модели C4: физическое развёртывание. Сценарий — **локальная разработка** (`docker compose -f deploy/docker-compose.yml up -d --wait`, `make dev`). Первичный источник: `deploy/docker-compose.yml` (T6), `deploy/README.md` §Dev Environment (T8), ADR-IMPL.PROCESS.development-tooling §8.

## Диаграмма

```mermaid
C4Deployment
    title Deployment — Dev Environment (docker-compose, localhost)

    Deployment_Node(devHost, "Host (developer machine)", "macOS / Linux / Windows (Docker Desktop)") {
        Deployment_Node(docker, "Docker Engine", "Docker Desktop / daemon") {

            Deployment_Node(publicNet, "edutrack-public (bridge)") {
                Deployment_Node(traefikNode, "traefik", "traefik:v3.1.2") {
                    Container(traefik, "Traefik (edge)", "Go", "Reverse proxy: entrypoints 80/443, dashboard :8080 (host 8082); dynamic.yml — routers api/spa, rate-limit, security-headers, circuit-breaker")
                }
            }

            Deployment_Node(devNet, "edutrack-net (bridge)") {

                Deployment_Node(frontendNode, "frontend", "node:24-alpine") {
                    Container(vite, "Vite dev server", "Node.js + Vite", "SPA hot-reload (HMR), порт 5173; прокси /api → backend:8080 (VITE_PROXY_TARGET)")
                }

                Deployment_Node(backendNode, "backend", "golang:1.26-alpine") {
                    Container(air, "air (hot-reload)", "Go", "Пересборка при изменении .go; запускает бинарник vedo-edutrack (M0.3: server)")
                    Container(api, "API-сервер (монолит)", "Go", "Бизнес-логика 10 bounded contexts; порт 8080")
                }

                Deployment_Node(pgNode, "postgres", "postgres:16-alpine") {
                    ContainerDb(pg, "PostgreSQL", "SQL", "Данные EduTrack (volume postgres_data); init.sql: uuid-ossp, pg_trgm, schema edutrack")
                }

                Deployment_Node(hubMockNode, "hub-mock", "backend/Dockerfile.mockhub (build)") {
                    Container(hubMock, "VEDO Hub mock (GraphQL)", "Go", "Стенд-заменитель VEDO Hub: POST /graphql (classes, graphNeighborhood, classDescendants…), онтология в памяти из .ttl (ONTOLOGY_FILE), :8081; healthcheck /healthz; ADR-DES.INFRA.mock-hub-strategy")
                }

                Deployment_Node(otelNode, "otel-collector", "otel/opentelemetry-collector-contrib:0.118.0") {
                    Container(otel, "OTel Collector", "Go", "Приём OTLP (4317/4318), PII-redaction, экспорт в Prometheus/Loki/Tempo")
                }

                Deployment_Node(obsNet, "Observability (edutrack-net)") {
                    Deployment_Node(promNode, "prometheus", "prom/prometheus:v3.4.0") {
                        Container(prom, "Prometheus", "Go", "Метрики, :9090; scrape backend/metrics + otel-collector")
                    }
                    Deployment_Node(lokiNode, "loki", "grafana/loki:3.5.0") {
                        Container(loki, "Loki", "Go", "Логи, :3100 (volume loki_data)")
                    }
                    Deployment_Node(tempoNode, "tempo", "grafana/tempo:2.6.1") {
                        Container(tempo, "Tempo", "Go", "Трейсы, :3200 (volume tempo_data)")
                    }
                    Deployment_Node(grafanaNode, "grafana", "grafana/grafana:11.4.0") {
                        Container(grafana, "Grafana", "Go", "Дашборды, :3000; datasources provisioned as-code (grafana-datasources.yml)")
                    }
                }
            }
        }
    }

    Rel(traefik, vite, "Роутинг edutrack.localhost → SPA (dev, frontend-dev profile)", "HTTP :5173")
    Rel(traefik, api, "Роутинг api.edutrack.localhost → API", "HTTP :8080")
    Rel(vite, api, "Прокси /api (dev, минуя edge)", "HTTP :8080")
    Rel(air, api, "Перезапускает при изменениях", "process")
    Rel(api, pg, "Читает и пишет", "SQL :5432")
    Rel(api, hubMock, "Читает онтологию (F0), GraphQL", "HTTP :8081")
    Rel(api, otel, "OTLP (трейсы/метрики/логи)", "gRPC/HTTP :4317/:4318")
    Rel(otel, prom, "Экспорт метрик", "HTTP :8889")
    Rel(otel, loki, "Экспорт логов", "HTTP :3100")
    Rel(otel, tempo, "Экспорт трейсов", "OTLP gRPC :4317")
    Rel(prom, grafana, "Источник данных", "HTTP")
    Rel(loki, grafana, "Источник данных", "HTTP")
    Rel(tempo, grafana, "Источник данных", "HTTP")
```

## Легенда

| Узел | Технология | Роль |
|------|-----------|------|
| **traefik** | traefik:v3.1.2 | Edge (опционально в dev): TLS, маршруты api/spa, rate limit, security headers, circuit breaker; дашборд на host :8082 |
| **frontend / Vite** | node:24-alpine | SPA dev-сервер с HMR; `pnpm dev --host`; прокси `/api` → backend |
| **backend / air** | golang:1.26-alpine | Модульный монолит; air пересобирает при изменении `.go` (`backend/.air.toml`) |
| **postgres** | postgres:16-alpine | Единственное хранилище MVP; volume `postgres_data`, init-скрипт `deploy/postgres/init.sql` |
| **hub-mock** | Go (сборка из `backend/cmd/mockhub`) | Стенд VEDO Hub: GraphQL read-only, онтология из `.ttl` (`traceability.ttl` → `/data/ontology.ttl`); `VEDO_HUB_API_URL=http://hub-mock:8081` |
| **otel-collector** | otel/opentelemetry-collector-contrib | Приём OTLP, redaction PII (152-ФЗ), экспорт в 3 бэкенда |
| **prometheus / loki / tempo / grafana** | — | Observability-стек provisioned as-code (`deploy/observability/`) |

**Особенности dev-контура:** браузер ходит напрямую на Vite (`localhost:5173`) и API через его прокси — Traefik не является обязательной точкой входа в dev; конфигурация `deploy/traefik/` актуальна для SaaS/staging. Сеть `edutrack-public` содержит только Traefik.

## Контекст

Диаграмма отражает инженерную платформу M0.2 (T6): один набор сервисов поднимается одной командой (`make up` / `docker compose up -d --wait`). Внутренняя сеть `edutrack-net` — служебная, `edutrack-public` — Traefik-вход. Тома: `postgres_data`, `grafana_data`, `loki_data`, `tempo_data`, `frontend_node_modules` (контейнерный node_modules). Переменные — из `deploy/.env` с дефолтами в compose.

## Связи с артефактами

| Артефакт | Роль |
|----------|------|
| `deploy/docker-compose.yml` | Определяет весь dev-стек (T6) |
| `deploy/observability/` | Конфиги collector/prometheus/loki/tempo/grafana (T3) |
| `deploy/traefik/` | Edge-конфиг для SaaS (T7) |
| `deploy/postgres/init.sql` | Расширения и схема при первом старте (T3) |
| `backend/Dockerfile` | Production-образ (в dev не используется — hot-reload) |
| [Container diagram](container-overview.md) | Логические контейнеры этого развёртывания |

## Связанные артефакты

- [Container overview](container-overview.md) — логический уровень 2
- [Контуры развёртывания](../adr/../deploy/README.md) — Community vs Enterprise (M0.2, T8)
- [ADR-IMPL.PROCESS.development-tooling](../adr/ADR-IMPL.PROCESS.development-tooling.md) §8 — dev-env
- [REQ-NFR-infra.compliance.environment-isolation](../requirements/REQ-NFR-infra.compliance.environment-isolation.md) — сети compose
