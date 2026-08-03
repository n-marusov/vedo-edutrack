# C4 Level 4: Deployment Diagram — Enterprise On-Prem

> Уровень 4 модели C4: физическое развёртывание. Сценарий — **Enterprise on-prem / private cloud**: минимальный контур (единый бинарник с embedded SPA + PostgreSQL), соответствует 242-ФЗ (данные в периметре). Первичный источник: `deploy/README.md` §Enterprise On-Prem (T8), `backend/Dockerfile` (T4–T5), ADR-DES.INFRA.modular-monolith-approach §7–8.

## Диаграмма

```mermaid
C4Deployment
    title Deployment — Enterprise On-Prem (single binary + PostgreSQL)

    Deployment_Node(browser, "Корпоративный браузер", "Chrome / Firefox / Safari") {
        Container(spa, "SPA (embedded)", "React (Go embed)", "Веб-интерфейс, отдаётся единым бинарником — не отдельный деплой")
    }

    Deployment_Node(entHost, "Enterprise host (on-prem / private cloud)", "Любая ОС с контейнерами (Docker) или systemd") {

        Deployment_Node(entNet, "Внутренняя сеть (периметр 242-ФЗ)") {
            Deployment_Node(backendNode, "backend", "vedo-edutrack (distroless, single artifact)") {
                Container(api, "API-сервер + SPA (Go embed)", "Go (distroless nonroot)", "10 bounded contexts + embedded SPA (M0.3); порт 8080; healthcheck /health; лёгкий TLS (termination на месте) или прямой доступ")
            }

            Deployment_Node(pgNode, "postgres", "postgres:16-alpine") {
                ContainerDb(pg, "PostgreSQL", "SQL", "Данные EduTrack; volume postgres_data; Atlas-миграции (drift = 0); резервное копирование pg_dump")
            }

            Deployment_Node(redisNode, "redis (опционально)", "redis:7-alpine") {
                Container(redis, "Redis", "—", "Пост-MVP: распределённый rate limiting, кэш read-моделей, token blacklist — не обязателен на MVP")
            }
        }
    }

    System_Ext(hub, "VEDO Hub", "Внешняя платформа онтологий (доступ по сети)")
    System_Ext(idp, "IdP (Keycloak)", "SSO/SAML — Enterprise (F6.5)")

    Rel(spa, api, "SPA вызывает API (один процесс)", "HTTP (in-process)")
    Rel(api, pg, "Читает и пишет", "SQL :5432")
    Rel(api, hub, "Читает онтологию, копирует подграф (F0.2)", "REST / MCP / SPARQL (read-only)")
    Rel(api, idp, "SSO/SAML, JWT (F6.5)", "SAML/OIDC")
    Rel(api, redis, "Кэш / rate limiting (пост-MVP)", "RESP :6379")
```

## Легенда

| Узел | Технология | Роль |
|------|-----------|------|
| **backend** | vedo-edutrack (distroless) | **Единственный артефакт**: SPA встроена в бинарник через Go embed (M0.3) — на one-artifact деплой без orchestration |
| **postgres** | postgres:16-alpine | Единственное обязательное хранилище; `pg_dump` перед миграциями; откат ≤ 15 мин |
| **redis** | redis:7-alpine | Опционально (пост-MVP): rate limiting при 2+ репликах, кэш read-моделей, token blacklist |
| **Traefik / OTel / Grafana** | — | **Не входят** в минимальный контур: телеметрия в периметре (242-ФЗ), TLS — на месте |

**Отличия от SaaS:** нет Traefik-эджа (или опциональный), нет nginx-контейнера (SPA embedded), observability-стек — по требованию Enterprise, без внешней телеметрии. Данные и логи остаются в периметре (242-ФЗ).

## Контекст

Контур Enterprise выводится из container strategy (`deploy/README.md`, T8): «single Go binary (embedded SPA) + PostgreSQL only; Redis optional; no orchestration required on MVP; K8s/Helm post-MVP». Соответствует `REQ-NFR-ops.compliance.support-sla` (один артефакт упрощает поддержку) и `REQ-NFR-infra.compliance.community-enterprise-isolation`. Требования SSO/SAML — через Keycloak (пост-MVP адаптер, F6.5), телеметрия — в периметре.

## Связи с артефактами

| Артефакт | Роль |
|----------|------|
| `backend/Dockerfile` | Distroless-образ с embedded SPA (T4–T5) |
| `deploy/postgres/init.sql` | Расширения и схема (T3) |
| [Container diagram](container-overview.md) | Логические контейнеры |

## Связанные артефакты

- [Deployment: Dev](deployment-dev.md) — локальный контур
- [Deployment: SaaS](deployment-saas.md) — Community-контур
- [Container overview](container-overview.md) — уровень 2
- [Container strategy](../../deploy/README.md) — Enterprise-контур (M0.2, T8)
- [ADR-DES.INFRA.modular-monolith-approach](../adr/ADR-DES.INFRA.modular-monolith-approach.md) §7–8 — Go embed
- [REQ-NFR-ops.compliance.support-sla](../requirements/REQ-NFR-ops.compliance.support-sla.md)
