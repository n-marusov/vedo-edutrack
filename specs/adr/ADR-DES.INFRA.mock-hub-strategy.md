# ADR-DES.INFRA.mock-hub-strategy

**Статус:** ПРЕДЛОЖЕНО
**Дата:** 2026-08-03
**Контекст:** Mock VEDO Hub для dev/тест/CI (M0.3, T20; RESEARCH session 2026-08-03)

## Контекст

- В dev-стеке нет VEDO Hub: `VEDO_HUB_API_URL=http://localhost:8081` указывает в никуда; внутри контейнера `localhost` ≠ хост.
- Контрактные тесты границы Hub (M1) и CI нуждаются в контролируемом стенде-заменителе.
- Онтология-сервис Hub предоставляет **read-only GraphQL-интерфейс**, созданный для навигации по графу (`graphNeighborhood`, `classDescendants`), тогда как REST API не имеет traversal-эндпоинта.

## Источники требований

- `REQ-FR-api.hub.read-ontology` (F0.1)
- `REQ-FR-api.hub.copy-subgraph` (F0.2)
- `REQ-NFR-api.availability.hub-dependency-sla`
- `REQ-NFR-process.dev.test-coverage`
- `ADR-DES.API.communication-patterns` §6

## Решение

Отдельный Docker-контейнер (`hub-mock`) в том же Go-модуле:

- **GraphQL-only** (`POST /graphql`), схема — копия SDL онтология-сервиса vedo-hub (embed как строка).
- Онтология в памяти из произвольного `.ttl` (`ONTOLOGY_FILE`; первый сид — `traceability.ttl`).
- Парсер — мини-Turtle (собственный, 0 зависимостей; `knakk/rdf` как fallback), `gqlparser/v2` — резолверы.
- Тот же хендлер доступен как `httptest.Server` для контрактных тестов M1 (`backend/internal/testing/mockhub`).
- Документированное исключение из ADR единого бинарника (`ADR-DES.API.cli-interface`): test-only, не поставляется, не в SBOM.

## Альтернативы (кратко)

| Альтернатива | Отклонение |
|--------------|-----------|
| REST-only mock | Не покрывает GraphQL-контракт (traversal-эндпоинты) |
| WireMock / Prism / MockServer | Не GraphQL-first, дублирование схемы в DSL |
| Реальный трипльстор (Fuseki/GraphDB) | Тяжёлый для CI, лишняя инфраструктура |
| gqlgen | Генерация кода для тестового инструмента избыточна |
| Hand-rolled GraphQL executor | gqlparser/v2 уже даёт парсинг/валидацию |
| In-process fake only | Не закрывает лакуну dev-стека (`localhost:8081`) и CI |

## Последствия и митигации

- **Пиннинг зависимостей:** `gqlparser/v2` — единственная новая (test-only) зависимость; фиксируется в go.mod.
- **Синхронизация контракта:** SDL копируется из vedo-hub; расхождение ловится контрактными тестами M1.
- **Второй `cmd/` entry:** `backend/cmd/mockhub` — осознанное исключение (test-only), помечено в AGENTS.md и ADR.

**Связанные артефакты:** [ADR-DES.API.communication-patterns](ADR-DES.API.communication-patterns.md) §6, [REQ-FR-api.hub.read-ontology](../requirements/REQ-FR-api.hub.read-ontology.md), [REQ-FR-api.hub.copy-subgraph](../requirements/REQ-FR-api.hub.copy-subgraph.md), `specs/c4/deployment-dev.md`, `specs/c4/deployment-test.md`
