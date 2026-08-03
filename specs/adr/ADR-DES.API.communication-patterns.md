# ADR-DES.API.communication-patterns

**Статус:** ПРИНЯТО
**Дата:** 2026-08-02
**Контекст:** Выбор паттернов межсервисного взаимодействия и API-контрактов (M0.1, T4)

EduTrack — клиент-серверное веб-приложение (`REQ-NFR-infra.compliance.client-server-web-app`) с внешними интеграторами: EdTech-платформы через REST API (F6.1), read-only SPARQL endpoint (F6.2), LMS-коннекторы (WebTutor, iSpring, SAP SuccessFactors — F6.3), webhooks (F6.4), SSO/SAML через Keycloak (F6.6), MCP-сервер для AI-агентов (F6.5). Два контура (Community / Enterprise) разделяют **один API-контракт** (`REQ-NFR-infra.compliance.community-enterprise-isolation`).

Внутри системы — модульный монолит (`ADR-DES.INFRA.modular-monolith-approach`): 10 bounded contexts в одном процессе; каскадные пересчёты (`ModuleMastered → RouteRecalculated → Resources/Stories → Plan`, каталог `specs/ddd/domain-events.md`) доставляются **in-process асинхронной шиной**, внешние события — через outbox. Стек (T3): Go + chi + oapi-codegen (OpenAPI-first), React + TS. Граница VEDO Hub — через `ontology-port` (ACL): REST API + MCP + SPARQL, контрактные тесты (`REQ-NFR-api.availability.hub-dependency-sla`).

Ключевые требования:
- **REST API**: p95 ≤ 200 мс при 1000 concurrent (`REQ-NFR-api.performance.latency-p95`), OpenAPI как источник истины (документация ≥ 90% эндпоинтов, `REQ-NFR-ops.compliance.user-documentation`).
- **Webhooks**: идемпотентность (`REQ-FR-api.webhooks.idempotency`, `REQ-NFR-api.availability.webhook-idempotency`): at-least-once доставка, dedup-ключ `event_id`, уникальный constraint в БД.
- **SPARQL**: read-only (`REQ-FR-api.sparql.read-only`), параметризация (`REQ-NFR-security.compliance.owasp-application-security`), rate limiting ≤ 10 запросов/мин, таймаут 30 с, truncation > 10 000 строк.
- **Надёжность внешних зависимостей**: circuit breaker на VEDO Hub, таймауты (`REQ-NFR-api.availability.hub-dependency-sla`).
- **Версионирование**: изменения API не должны ломать интеграторов (EdTech-платформы, LMS); предсказуемая эволюция контракта.

**Требование-источник:**
- `REQ-NFR-api.performance.latency-p95`
- `REQ-NFR-ops.performance.scalability`
- `REQ-NFR-api.availability.webhook-idempotency`
- `REQ-FR-api.webhooks.idempotency`
- `REQ-FR-api.sparql.read-only`
- `REQ-NFR-security.compliance.owasp-application-security`
- `REQ-NFR-api.availability.hub-dependency-sla`
- `REQ-NFR-ops.compliance.user-documentation`
- `REQ-NFR-infra.compliance.community-enterprise-isolation`
- `REQ-FR-api.mcp.server`
- `REQ-FR-api.hub.subscribe-updates`

**Решение:**

Принять **синхронный REST (OpenAPI-first) + асинхронные события (in-process шина + outbox + webhooks)** как коммуникационную модель системы.

1. **Синхронный канал — REST/JSON (OpenAPI-first)**:
   - Единый контракт: OpenAPI 3.1-спека в репозитории — источник истины (код сервера через `oapi-codegen`, клиент через `openapi-typescript`, документация через swagger-ui, контрактные тесты в CI).
   - REST — для запросов с ожиданием ответа: вычисление маршрута, запрос прогресса/покрытия, CRUD-ресурсов, read-модели визуализации.
   - p95 ≤ 200 мс: тонкие адаптеры, кэширование read-моделей, параллельные вызовы Hub с таймаутами.
   - Idempotent-safe методы (`PUT`, `DELETE`) + `If-Match`/`version` для optimistic locking (`REQ-NFR-data.availability.write-consistency`).

2. **Асинхронный канал — события**:
   - **In-process шина** (внутри монолита): каскады пересчёта (`RouteRecalculated`, `ModuleMastered`, `PlanDeviationDetected`, `StandardDeficitDetected`, …) — асинхронная доставка в пределах процесса, гарантия порядка для одного агрегата (из ADR модульного монолита).
   - **Outbox в PostgreSQL**: события, уходящие наружу (webhooks, LMS, MCP-нотификации), пишутся транзакционно с бизнес-изменением; worker доставляет с гарантией at-least-once.
   - **Внешний канал для интеграторов — webhooks**: `route.recalculated`, `module.mastered`, `plan.deviated`, `standard.risk_detected` (форматы из каталога событий) с подписью секретом, retry с экспоненциальной задержкой, dedup по `event_id`.

3. **Идемпотентность webhook** (приём и отправка):
   - Каждое событие несёт уникальный `event_id` (UUID v4), неизменный при повторах.
   - Приём: уникальный constraint `(source, event_id)` в `integrations.event_dedup`; дубликат → `200 OK` с `duplicate=true`, без изменения состояния; конкурентные доставки — ровно одно применение (второй insert отклоняется constraint'ом). (`REQ-FR-api.webhooks.idempotency`, `REQ-NFR-api.availability.webhook-idempotency`.)
   - Отправка: outbox-запись с `event_id`, идемпотентная доставка, ретраи после сбоя/рестарта (состояние дедупликации в БД, переживает рестарт).

4. **Версионирование API — URL-path versioning** (`/api/v1/...`):
   - Мажорные версии в URL (`/api/v1`, `/api/v2`); минорные изменения — аддитивные (новые поля, эндпоинты) без смены URL.
   - Контракт: OpenAPI-спека версионируется; breaking change → новая мажорная версия, старые версии живут по SLA (документированный период поддержки, deprecation-заголовок + warning).
   - Причина: URL-версия проста для интеграторов (EdTech/LMS), не требует контент-negotiation, хорошо поддерживается шлюзом, удобна для параллельного роутинга в chi.
   - Аддитивные изменения не требуют новой версии; удаление/переименование полей — мажорная версия.

5. **SPARQL endpoint — read-only, параметризованный** (F6.2):
   - Только `SELECT/ASK/DESCRIBE`; мутирующие формы (`INSERT DATA`, `DELETE`, `LOAD`, `CLEAR`, `CREATE`) → `403 Forbidden`.
   - Аутентификация (API-ключ/token) → `401`; rate limit ≤ 10 запросов/мин на клиента → `429` с `Retry-After` (`REQ-NFR-security.compliance.owasp-application-security`).
   - Таймаут выполнения 30 с → `504 Gateway Timeout`; результат > 10 000 строк обрезается с `truncated=true`.
   - Параметризация запросов: пользовательский ввод передаётся параметрами/предобработанными литералами, никогда не конкатенируется в SPARQL-строку (тест-сьют на injection, `REQ-FR-api.sparql.read-only`, OWASP).
   - Endpoint проксирует запросы к VEDO Hub (SPARQL-эндпоинт Hub) через `ontology-port` с circuit breaker и кэшированием результатов.

6. **Граница VEDO Hub — `ontology-port` (ACL)**:
   - Единственный путь к онтологии: REST API Hub + MCP + GraphQL-интерфейс онтология-сервиса + SPARQL через адаптеры `ontology-port` (F0). GraphQL (read-only, traversal-эндпоинты `graphNeighborhood`/`classDescendants`) — целевой канал для навигации по графу; контракт SDL копируется в репозиторий и покрывается стендом `hub-mock` (`ADR-DES.INFRA.mock-hub-strategy`).
   - Синхронные вызовы Hub: таймаут ≤ 3 с, circuit breaker, retry с backoff (`REQ-NFR-api.availability.hub-dependency-sla`).
   - Копирование подграфа (F0.2) — асинхронно, инкрементально, с версионированием по `ontology_version`.
   - Обновления онтологии — через подписку (webhook/events от Hub, F0.3, `REQ-FR-api.hub.subscribe-updates`): `OntologyUpdated` → каскадный пересчёт.

7. **MCP-сервер (F6.5)** — read-ориентированный доступ для AI-агентов: те же read-модели и REST-слой за MCP-адаптером; инструменты только на чтение (маршрут, прогресс, покрытие, лакуны) с теми же auth/rate-limit/параметризацией. Без прямого доступа к БД — через Application-слой.

8. **Коммуникация модулей между собой** — только события (in-process) и синхронные вызовы через порты Application-слоя в пределах одного процесса; прямых межмодульных REST-вызовов нет (монолит — один процесс). Вынос модуля в сервис (пост-MVP) — через те же порты: in-process шина → брокер, вызовы портов → gRPC/REST (слоты эволюции из ADR монолита).

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **GraphQL** | ⚠️ | Интеграторам нужен стабильный REST-контракт (F6.1) и SPARQL; GraphQL добавил бы второй контракт, сложнее кэширование read-моделей, слабее AI-корпус для вайбкодинга; не заменяет webhooks/events. Остаётся опцией пост-MVP для внутреннего BFF, не для внешнего API |
| **gRPC** | ⚠️ | Отличный для внутренних сервисов, но внешним интеграторам (EdTech/LMS/SPARQL-экосистема) нужен REST/JSON + OpenAPI; gRPC-web не зрелый для публичных API; для монолита in-process gRPC не нужен. Отклонён для внешнего контракта (возможен пост-MVP между вынесенными сервисами) |
| **Версионирование в заголовке (Accept: application/vnd.api+json;version=2)** | ⚠️ | Тоньше гранулярность, но сложнее для интеграторов (контент-negotiation), хуже кэшируется, не виден в URL/документации, сложнее дебаг; URL-версия проще и нагляднее для B2B-контракта |
| **Event-driven только (без REST)** | ❌ | Большинство сценариев — запрос-ответ (прогресс, покрытие, отчёты); события без запросов не покрывают read-модели и интеграторов |
| **Saga/распределённые транзакции между модулями** | ❌ | Монолит: транзакции на уровне агрегата внутри одного процесса; саги не нужны на MVP (нет распределённых сервисов). Слоты эволюции — на будущее |
| **Прямой доступ к БД внешними потребителями** | ❌ | Нарушает Clean Architecture и границы контекстов; все запросы — через API/порты |
| **WebSocket/SSE как основной канал** | ⚠️ | Реалтайм-обновления (дашборды) — да, но как дополнение к REST (fan-out из in-process шины на WebSocket/SSE для visualization); не основной контракт |
| **Message Broker с MVP** | ⚠️ | Outbox + polling достаточно для MVP-объёма событий; брокер (Kafka/RabbitMQ/NATS) — пост-MVP по триггеру (10× нагрузка, кросс-нодовый fan-out, тяжёлый outbox-трафик) — замена без изменения ядра |

**Последствия:**

*Положительные:*
- Один источник истины (OpenAPI) → документация, код, контрактные тесты генерируются из спеки; изменения API контролируются CI (дрейф спеки/кода = ошибка).
- URL-версионирование простое для B2B-интеграторов и совместимо с webhook/LMS-экосистемой.
- Идемпотентность webhook гарантирована на уровне хранилища (unique constraint) — нет двойного начисления прогресса; состояние дедупликации переживает рестарт.
- SPARQL read-only с параметризацией и rate limit закрывает F6.2 и OWASP (инъекции), защищает Hub от перегрузки.
- in-process события + outbox: каскады мгновенные (p95 ≤ 200 мс), внешние интеграции надёжные (at-least-once, retry, dedup).
- MCP-сервер переиспользует REST/Application-слой — нет второго пути к данным.
- Контуры Community/Enterprise разделяют один контракт — изоляция только конфигурацией (`REQ-NFR-infra.compliance.community-enterprise-isolation`).

*Отрицательные и смягчение:*
- **REST + события — две модели взаимодействия** → смягчение: чёткое правило выбора канала (запрос-ответ → REST; уведомление/каскад → событие); документируется в контракте и ревью.
- **URL-версионирование удваивает код поддержки старых версий** → смягчение: деплоятся обе версии на одном сервере (chi-роутер), старые версии закрываются по SLA; deprecation-политика.
- **At-least-once доставка webhook** → смягчение: потребитель обязан обрабатывать дубликаты (dedup по `event_id` — обязательный контракт webhook-события, тесты двойной/конкурентной доставки в CI).
- **SPARQL-эндпоинт — вектор атаки (инъекции, перегрузка)** → смягчение: read-only whitelist, параметризация, rate limit, таймауты, circuit breaker на Hub.
- **Outbox-polling — задержка внешних событий** → смягчение: на MVP допустимо; при потребности — LISTEN/NOTIFY или брокер пост-MVP.
- **Синхронные вызовы Hub (SPARQL/REST) добавляют латентность** → смягчение: таймауты, кэш подграфа (in-memory, иммутабелен), параллелизм, fallback-стратегии; p95 < 1 с для пересчёта из `REQ-NFR-api.performance.latency-p95`.

**Связанные артефакты:**
- [ADR-DES.INFRA.modular-monolith-approach](ADR-DES.INFRA.modular-monolith-approach.md) — in-process шина, outbox, слоты эволюции
- [ADR-DES.DATA.storage-strategy](ADR-DES.DATA.storage-strategy.md) — event_dedup, outbox-таблица в PostgreSQL (T4)
- [ADR-DES.STACK.framework-vs-vs](ADR-DES.STACK.framework-vs-vs.md) — OpenAPI-first (oapi-codegen, openapi-typescript, swagger-ui)
- [Каталог доменных событий](../ddd/domain-events.md) — события ядра и webhook-представления, матрица каскадов
- [Карта контекстов](../ddd/context-map.md) — `ontology-port` как ACL к VEDO Hub
- C4-диаграммы (T6): контейнеры «API-сервер», «SPARQL-эндпоинт», «VEDO Hub (external)»
- ADR `DES.SECURITY.rbac-model` (T5) — auth enforcement points на API-границе
