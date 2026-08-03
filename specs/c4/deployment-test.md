# C4 Level 4: Deployment Diagram — CI Test Environment (GitHub Actions)

> Уровень 4 модели C4: физическое развёртывание. Сценарий — **непрерывная интеграция** (GitHub Actions, `deploy/ci/run-gates.sh --tier delivery`). Первичный источник: `.github/workflows/ci.yml` (T12), `deploy/ci/gates.yaml` (T16), ADR-DES.INFRA.mock-hub-strategy (T20).

## Диаграмма

```mermaid
C4Deployment
    title Deployment — CI (GitHub Actions runner)

    Deployment_Node(ciRunner, "GitHub Actions runner", "ubuntu-latest") {

        Deployment_Node(services, "Service containers (compose)") {
            Deployment_Node(pgNode, "postgres", "postgres:16-alpine") {
                ContainerDb(pg, "PostgreSQL", "SQL", "Тестовая БД EduTrack; миграции + seed применяются в test-джобе")
            }
            Deployment_Node(hubMockNode, "hub-mock", "backend/Dockerfile.mockhub (build)") {
                Container(hubMock, "VEDO Hub mock (GraphQL)", "Go", "POST /graphql; онтология из traceability.ttl; контрактные/e2e тесты ходят сюда вместо реального Hub")
            }
        }

        Deployment_Node(runner, "Runner tools", "go / pnpm / docker / golangci-lint / biome") {
            Container(gates, "Gate runner (run-gates.sh)", "bash", "fast/delivery тиры; lint → typecheck → test → coverage → security → build → smoke")
            Container(e2e, "E2E-тесты (Playwright)", "Node", "tests/e2e/gui + tests/e2e/api против backend (API+SPA, unified binary)")
            Container(contract, "Контрактные тесты Hub-границы", "Go (testcontainers)", "httptest.Server из backend/internal/testing/mockhub (M1)")
        }
    }

    Rel(gates, pg, "Миграции + seed (тестовая БД)", "SQL :5432")
    Rel(gates, hubMock, "Smoke: POST /graphql classes-query → 200", "HTTP :8081")
    Rel(e2e, pg, "Тестовые данные", "SQL")
    Rel(contract, hubMock, "Читает онтологию (F0), GraphQL", "HTTP :8081")
```

## Контекст

CI собирает unified backend-образ (API + embedded SPA, `backend/Dockerfile`) и образ `hub-mock` (`backend/Dockerfile.mockhub`), затем прогоняет гейты delivery-тира. Контрактные тесты границы Hub (M1) и e2e-сценарии ходят на `hub-mock` вместо реального VEDO Hub — контролируемый стенд закрывает зависимость `REQ-NFR-api.availability.hub-dependency-sla` в CI.

## Связи с артефактами

| Артефакт | Роль |
|----------|------|
| `.github/workflows/ci.yml` | CI-пайплайн: build обоих образов + smoke (T12, T25) |
| `deploy/ci/gates.yaml` | Манифест гейтов (T16) |
| `backend/Dockerfile.mockhub` | Образ стенда (T24) |
| `backend/internal/testing/mockhub` | Тот же хендлер как `httptest.Server` (T22) |
| `traceability.ttl` | Первый сид онтологии стенда |
| [ADR-DES.INFRA.mock-hub-strategy](../adr/ADR-DES.INFRA.mock-hub-strategy.md) | Дизайн стенда (T20) |
