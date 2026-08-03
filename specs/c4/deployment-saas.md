# C4 Level 4: Deployment Diagram — Community SaaS / Staging

> Уровень 4 модели C4: физическое развёртывание. Сценарий — **SaaS (Community) / staging**: Traefik-эдж + compose-стек, образы из GHCR. Первичный источник: `deploy/README.md` §SaaS Deployment (T8), `deploy/traefik/` (T7), `deploy/docker-compose.yml` (T6), `.github/workflows/ci.yml` (T12).

## Диаграмма

```mermaid
C4Deployment
    title Deployment — Community SaaS / Staging (Traefik edge + GHCR)

    Deployment_Node(browser, "Пользовательский браузер", "Chrome / Firefox / Safari") {
        Container(spa, "SPA", "React (Vite build)", "Веб-приложение EduTrack")
    }

    Deployment_Node(vps, "VPS / VM (SaaS-хост)", "Ubuntu LTS, Docker Engine") {

        Deployment_Node(publicNet, "edutrack-public") {
            Deployment_Node(traefikNode, "traefik", "traefik:v3.1.2 (GHCR)") {
                Container(traefik, "Traefik (edge)", "Go", "TLS (Let's Encrypt, ACME-резолвер), роутеры api/spa, rate-limit, security-headers (CSP/HSTS), circuit-breaker; blue-green веса (post-MVP)")
            }
        }

        Deployment_Node(prodNet, "edutrack-net") {

            Deployment_Node(frontendNode, "frontend (nginx)", "nginxinc/nginx-unprivileged:1.27-alpine-slim") {
                Container(nginx, "nginx (SPA-статика)", "nginx", "Отдаёт dist/, SPA-fallback, gzip, non-root UID 101, порт 8080")
            }

            Deployment_Node(backendNode, "backend", "vedo-edutrack (distroless, GHCR)") {
                Container(api, "API-сервер (монолит)", "Go (distroless nonroot)", "10 bounded contexts; порт 8080; liveness /healthz + readiness /readyz (PostgreSQL); health-проба `vedo-edutrack health`; из M0.3 — SPA через Go embed (single artifact)")
            }

            Deployment_Node(pgNode, "postgres", "postgres:16-alpine") {
                ContainerDb(pg, "PostgreSQL", "SQL", "Данные EduTrack; volume postgres_data; Atlas-миграции (drift = 0)")
            }

            Deployment_Node(otelNode, "otel-collector", "otel/opentelemetry-collector-contrib") {
                Container(otel, "OTel Collector", "Go", "OTLP-приём, PII-redaction (152-ФЗ), экспорт в Prometheus/Loki/Tempo")
            }

            Deployment_Node(obs, "Observability") {
                Deployment_Node(promNode, "prometheus", "prom/prometheus:v3.4.0") {
                    Container(prom, "Prometheus", "Go", "Метрики :9090")
                }
                Deployment_Node(lokiNode, "loki", "grafana/loki:3.5.0") {
                    Container(loki, "Loki", "Go", "Логи :3100")
                }
                Deployment_Node(tempoNode, "tempo", "grafana/tempo:2.6.1") {
                    Container(tempo, "Tempo", "Go", "Трейсы :3200")
                }
                Deployment_Node(grafanaNode, "grafana", "grafana/grafana:11.4.0") {
                    Container(grafana, "Grafana", "Go", "Дашборды :3000, datasources as-code")
                }
            }
        }
    }

    System_Ext(hub, "VEDO Hub", "Внешняя платформа онтологий")
    System_Ext(idp, "IdP (Keycloak)", "Enterprise SSO — пост-MVP")

    Rel(spa, traefik, "HTTPS", "443")
    Rel(traefik, nginx, "Отдаёт SPA (edutrack.localhost)", "HTTP :8080")
    Rel(traefik, api, "Маршрутизирует API (api.edutrack.localhost, /api)", "HTTP :8080")
    Rel(nginx, api, "SPA вызывает API", "HTTP (логический контракт)")
    Rel(api, pg, "Читает и пишет", "SQL :5432")
    Rel(api, hub, "Читает онтологию, копирует подграф (F0.2)", "REST / MCP / SPARQL (read-only)")
    Rel(api, idp, "SSO/SAML, JWT (F6.5, пост-MVP)", "SAML/OIDC")
    Rel(api, otel, "OTLP", "gRPC/HTTP")
    Rel(otel, prom, "Метрики", "HTTP :8889")
    Rel(otel, loki, "Логи", "HTTP :3100")
    Rel(otel, tempo, "Трейсы", "OTLP gRPC :4317")
    Rel(prom, grafana, "Datasource", "HTTP")
    Rel(loki, grafana, "Datasource", "HTTP")
    Rel(tempo, grafana, "Datasource", "HTTP")
```

## Легенда

| Узел | Технология | Роль |
|------|-----------|------|
| **traefik** | traefik:v3.1.2 | Единственная точка входа: TLS (Let's Encrypt staging→prod), роутинг SPA/API, rate limiting, CSP/HSTS, circuit breaker; blue-green weighted (post-MVP) |
| **frontend (nginx)** | nginxinc/nginx-unprivileged:1.27-alpine-slim | SPA-статика, non-root UID 101, SPA-fallback, 404 для отсутствующих файлов; вариант B (SaaS/CDN) |
| **backend** | vedo-edutrack (distroless) | Модульный монолит ~1.4 МБ, nonroot, healthcheck; M0.3 — embedded SPA (single artifact) |
| **postgres** | postgres:16-alpine | Единственное хранилище; Atlas-миграции, drift = 0 |
| **Observability** | OTel → Prom/Loki/Tempo + Grafana | Golden-signals, provisioning as-code, PII-redaction |

**GHCR-поставка:** образы собираются в CI (T12) и пушатся на `main` с тегами `:latest` и `:sha-<commit>` (`deploy/ci/push-images.sh`). Размер distroless-образа контролируется гейтом `docker-build` (≤ 20 МБ).

## Контекст

Диаграмма соответствует SaaS-контуру Community из `deploy/README.md` (T8): Traefik-эдж + полный OTel-стек + nginx-вариант SPA. Blue-green деплой — пост-MVP через weighted-сервисы Traefik (canary + kill switch ≤ 5 мин, `REQ-NFR-ops.release.canary-kill-switch`). В отличие от dev — нет air/Vite; образы из реестра. Enterprise-контур — отдельная диаграмма (`deployment-enterprise.md`).

## Связи с артефактами

| Артефакт | Роль |
|----------|------|
| `deploy/traefik/traefik.yml` | Static: entrypoints, ACME-резолвер, docker/file-провайдеры (T7) |
| `deploy/traefik/dynamic.yml` | Routers api/spa, middlewares, weighted blue-green (T7) |
| `backend/Dockerfile` | Distroless-образ бэкенда (T4) |
| `frontend/Dockerfile` | nginx-вариант SPA (T5) |
| `.github/workflows/ci.yml` + `deploy/ci/` | Сборка и публикация образов (T12, T16) |
| [Container diagram](container-overview.md) | Логические контейнеры |

## Связанные артефакты

- [Deployment: Dev](deployment-dev.md) — локальный контур
- [Deployment: Enterprise](deployment-enterprise.md) — on-prem минимальный
- [Container overview](container-overview.md) — уровень 2
- [Container strategy](../../deploy/README.md) — SaaS-контур (M0.2, T8)
- [ADR-IMPL.PROCESS.development-tooling](../adr/ADR-IMPL.PROCESS.development-tooling.md) §7–8
- [REQ-NFR-ops.release.canary-kill-switch](../requirements/REQ-NFR-ops.release.canary-kill-switch.md)
