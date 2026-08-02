# ADR-DES.DATA.storage-strategy

**Статус:** ПРИНЯТО
**Дата:** 2026-08-02
**Контекст:** Выбор технологии хранения и подход к проектированию схемы (M0.1, T4)

EduTrack — сервис-надстройка над VEDO Hub: он **не хранит онтологии** (граф знаний живёт в Hub, читается через REST/SPARQL/MCP, `REQ-FR-api.hub.read-ontology`). EduTrack хранит данные обучения: профили учеников, зафиксированные планы (снэпшоты маршрута с таймлайном, `REQ-FR-plan.fixation.snapshot`), записи прогресса, диагностики лакун, отчёты покрытия. Маршрут — **функция, а не документ**: вычисляется в памяти из скопированного подграфа онтологии (F0.2, `REQ-FR-api.hub.copy-subgraph`); БД хранит только снэпшоты (планы) и записи прогресса, а не маршруты как данные.

Архитектурный стиль — модульный монолит (`ADR-DES.INFRA.modular-monolith-approach`): 10 bounded contexts в одном процессе, in-process событийная шина + outbox для внешних интеграций. Внутренняя структура — Clean Architecture (`ADR-DES.INFRA.clean-architecture-adoption`): доступ к данным — через порты, адаптеры на периферии. Стек выбран (T3): Go + chi (бэкенд), React + TS (фронт), OpenAPI-first.

Ключевые требования, влияющие на выбор:
- **p95 ≤ 200 мс при 1000 concurrent** (`REQ-NFR-api.performance.latency-p95`) — read-heavy визуализация (F4), интерактивный пересчёт.
- **Горизонтальное масштабирование** (`REQ-NFR-ops.performance.scalability`): stateless-реплики, 10× пик, автоскейлинг ≤ 5 мин — БД не должна быть узким местом записи.
- **Консистентность записи** (`REQ-NFR-data.availability.write-consistency`): optimistic locking, 0 lost updates при конкурентной записи прогресса (родитель/учитель/HR/LMS-webhook пишут одновременно).
- **RPO ≤ 1 ч / RTO ≤ 4 ч** (`REQ-NFR-data.availability.backup-rpo`), **обратимые миграции** (`REQ-NFR-data.availability.migration-rollback`: up→down→up в CI, автооткат ≤ 15 мин).
- **152-ФЗ / 242-ФЗ** (`REQ-NFR-security.compliance.pii-152-fz`, `REQ-NFR-data.compliance.data-residency`): шифрование PII at-rest (AES-256), локализация данных в РФ, Enterprise on-prem.
- **Идемпотентность webhook** (`REQ-FR-api.webhooks.idempotency`, `REQ-NFR-api.availability.webhook-idempotency`): уникальный constraint на `event_id` как механизм дедупликации.
- **Enterprise on-prem для сторонних заказчиков** (`REQ-NFR-ops.compliance.support-sla`): один компонент хранения (без обязательных внешних сервисов) радикально снижает стоимость установки и поддержки.

**Требование-источник:**
- `REQ-NFR-api.performance.latency-p95`
- `REQ-NFR-ops.performance.scalability`
- `REQ-NFR-data.availability.write-consistency`
- `REQ-NFR-data.availability.backup-rpo`
- `REQ-NFR-data.availability.migration-rollback`
- `REQ-NFR-security.compliance.pii-152-fz`
- `REQ-NFR-data.compliance.data-residency`
- `REQ-FR-plan.fixation.snapshot`
- `REQ-FR-api.hub.copy-subgraph`
- `REQ-FR-api.webhooks.idempotency`
- `REQ-NFR-api.availability.webhook-idempotency`
- `REQ-NFR-ops.compliance.support-sla`
- `REQ-NFR-security.compliance.owasp-application-security`

**Решение:**

Принять **PostgreSQL** как единственную систему хранения EduTrack (MVP), со схемой по модулям и следующими принципами:

1. **PostgreSQL — единственный datastore (MVP)**. Реляционная модель, транзакции ACID, mature-экосистема. Кэш подграфа — in-memory в Go-процессе (иммутабелен по `ontologyVersion`); Redis (или иной брокер) — пост-MVP адаптер по триггеру (см. ADR модульного монолита). Enterprise on-prem: один артефакт + PostgreSQL — без внешних зависимостей.

2. **Схема-на-модуль (schema-per-module)**. Каждый bounded context — собственная схема PostgreSQL (`route_planning`, `execution_progress`, `gap_coverage`, `plan_management`, `ontology_port`, `resources`, `practice_life`, `visualization`, `identity_access`, `integrations`). Изоляция границ контекстов на уровне хранилища: модуль не имеет доступа к чужим таблицам (гранты на схему), данные пересекают границы только событиями. Контуры Community/Enterprise — те же схемы в отдельных БД/кластерах или tenant-колонка (параметр развёртывания, `REQ-NFR-infra.compliance.community-enterprise-isolation`).

3. **Агрегат = кластер таблиц**. Корень агрегата — таблица-агрегат, внутренние сущности и VO — дочерние таблицы; инварианты агрегата охраняются транзакцией на уровне агрегата (единая строка-корень). Таблицы называются по терминам глоссария (не «маршруты»): `learning_plans` (снэпшоты), `plan_steps`, `progress_log`, `trajectories`, `gap_diagnoses`, `coverage_reports`, `learners`, `accounts`.

4. **Снэпшоты плана — неизменяемые строки**. `learning_plans` — версионируемые неизменяемые снэпшоты: состав модулей, плановые даты, контрольные точки, `ontology_version`, хэш снэпшота. Изменение плана = новая версия (строка), не UPDATE (`REQ-FR-plan.fixation.snapshot`: неизменяемость после `PlanFixed`).

5. **Optimistic locking**. Каждая записываемая сущность (позиция, прогресс, план) содержит `version`; запись с устаревшей версией → `409 Conflict` (`REQ-NFR-data.availability.write-consistency`). Идемпотентность webhook: уникальный constraint на `(source, event_id)` в `integrations.event_dedup` — второй insert конкурентной доставки отклоняется (`REQ-FR-api.webhooks.idempotency`, `REQ-NFR-api.availability.webhook-idempotency`).

6. **Outbox-таблица в PostgreSQL**. `integrations.outbox` — in-DB очередь для внешних событий (webhooks, LMS, MCP): транзакционная публикация «бизнес-изменение + outbox-запись» атомарна; выгрузка — отдельным worker-процессом (polling), доставка идемпотентная (`event_id`). Redis/Broker — пост-MVP по триггеру (см. ADR модульного монолита).

7. **Read-модели для визуализации**. Схема `visualization` содержит материализованные проекции (read-модели, CQRS-light): таблицы дашбордов, покрытия, лакун — обновляются подписчиками событий ядра, не прямыми запросами в схемы ядра. Read-нагрузка F4 не бьёт по транзакционным таблицам.

8. **Схема данных не хранит маршруты как документ**. Никакой таблицы «routes» как данных: таблица `route_computation_results` (при необходимости) — только кэш последнего результата вычисления с `ontology_version` и TTL; маршрут всегда пересчитывается из подграфа (`REQ-FR-api.hub.copy-subgraph`).

9. **Безопасность данных**. PII-колонки шифруются на уровне приложения (AES-256-GCM, ключи вне БД — KMS/секрет-менеджер) или column-level encryption; at-rest шифрование диска БД; TLS ≥ 1.2 in-transit (`REQ-NFR-security.compliance.pii-152-fz`). Локализация: БД Community и Enterprise (для РФ) — в ЦОД на территории РФ (`REQ-NFR-data.compliance.data-residency`). Доступ через порты Clean Architecture — параметризованные запросы (sqlc), никакой конкатенации пользовательского ввода (`REQ-NFR-security.compliance.owasp-application-security`).

10. **Миграции — обратимые, через Atlas** (drift-детекция, из DESCRIPTION.md): каждая миграция up/down, up→down→up в CI, автооткат ≤ 15 мин при сбое (`REQ-NFR-data.availability.migration-rollback`). Схема-на-модуль упрощает независимый ревью миграций.

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **MySQL** | ⚠️ | Зрелая реляционная БД, дешевле в малых командах; но слабее в продвинутых типах (JSONB, диапазонные типы для дат плана, GENERATED ALWAYS, частичные уникальные индексы для dedup), слабее observability-экосистема, исторически хуже с миграциями на лету; PostgreSQL выбран как функционально более полный вариант с тем же паттерном эксплуатации |
| **MongoDB** | ❌ | Документная модель не соответствует агрегатам с инвариантами и межмодульной ссылочной целостностью; транзакции multi-document возможны, но слабее; отсутствие строгой схемы ломает контракты вайбкодинга и тестов; снэпшоты/отчёты — реляционная структура; optimistic locking вручную |
| **Redis как основной datastore** | ❌ | In-memory — не выдерживает RPO ≤ 1 ч / RTO ≤ 4 ч как единственный источник истины; персистентность (RDB/AOF) слабее ACID; нужен durable storage в любом случае. Redis остаётся пост-MVP адаптером (кэш, rate limiting, fan-out) |
| **Graph DB (Neo4j) как основное хранилище** | ❌ | Граф знаний живёт в VEDO Hub (не наша БД); EduTrack хранит обучение (снэпшоты, прогресс, отчёты) — реляционные данные; копия подграфа — in-memory структура для вычислений, не персистентное хранилище. Graph DB добавила бы вторую БД без выигрыша |
| **Event Store / Event Sourcing** | ⚠️ | Event Sourcing — механизм, а не необходимость: домен фиксирует план как снэпшот, а не как журнал событий; журнал событий для траектории ведём в реляционной модели (`progress_log`, `trajectories`); ES добавил бы сложность (проекции, переигрывание) без требований аудита полного журнала |
| **Key-Value (DynamoDB/Cassandra)** | ❌ | Операционный ад для on-prem Enterprise; слабая поддержка запросов отчётов/покрытия (агрегации по доменам); scale-out не нужен на MVP (1000 concurrent, stateless-реплики) |
| **Без БД (файловое/JSON-хранение)** | ❌ | Не выдерживает конкурентность, ACID, RPO/RTO, многопользовательский доступ; нет места в архитектуре модульного монолита |

**Последствия:**

*Положительные:*
- Один компонент хранения (PostgreSQL) на оба контура: простой on-prem для Enterprise (один контейнер БД + артефакт), простые бэкапы и restore (RPO/RTO), дешевле поддержка (SLA Enterprise ≤ 1 ч).
- Реляционная модель точно отражает агрегаты и инварианты; транзакции на уровне агрегата охраняют «неизменяемый снэпшот» и дедупликацию webhook.
- Schema-per-module даёт изоляцию bounded contexts на уровне хранилища — границы не обходятся случайным JOIN; арх-тесты могут проверять отсутствие кросс-схемных обращений.
- Обратимые миграции (Atlas, up→down→up в CI) — прямое выполнение `REQ-NFR-data.availability.migration-rollback` (автооткат ≤ 15 мин).
- PostgreSQL — огромный AI-корпус и экосистема (sqlc, Atlas, OTel-интеграции): стабильный вайбкодинг и надёжные интеграционные тесты (testcontainers-go).
- in-memory кэш подграфа (иммутабелен по версии) убирает Redis с критического пути MVP — меньше движущихся частей для Enterprise.

*Отрицательные и смягчение:*
- **PostgreSQL — единая точка отказа** → смягчение: managed-PG в Community (multi-AZ ≥ 2 зон, `REQ-NFR-infra.availability.multi-az-geo-dr`), hot standby + PITR для Enterprise; RPO ≤ 1 ч через WAL-архивирование (PITR), RTO ≤ 4 ч через restore drill.
- **Горизонтальная запись ограничена одним мастером** → смягчение: на MVP/10× пике запись не является узким местом (stateless-реплики, read-heavy нагрузка F4 уходит на read-модели); при необходимости — partition by tenant / read replicas (PostgreSQL 16+ logical replication) пост-MVP.
- **Schema-per-module увеличивает число схем** → смягчение: конвенция нейминга таблиц `<модуль>_<сущность>` внутри схемы, гранты по схемам, документация схем в репозитории (SQL-файлы Atlas).
- **Шифрование PII на уровне приложения** → смягчение: единый crypto-адаптер за портом (ключи в KMS/секрет-менеджере), вращение ключей, тесты шифрования (поиск известных значений в сырых файлах БД = 0 совпадений, `REQ-NFR-security.compliance.pii-152-fz`).
- **Outbox-polling добавляет латентность доставки событий** → смягчение: на MVP допустимо (события не критичны к субсекундной доставке наружу); при необходимости — LISTEN/NOTIFY или брокер пост-MVP.
- **Миграции требуют дисциплины** → смягчение: down-миграции обязательны (гейт CI), drift-детекция Atlas, окно сопровождения (`REQ-NFR-ops.release.maintenance-windows`).

**Связанные артефакты:**
- [ADR-DES.INFRA.modular-monolith-approach](ADR-DES.INFRA.modular-monolith-approach.md) — outbox, in-process шина, кэш за портами
- [ADR-DES.INFRA.clean-architecture-adoption](ADR-DES.INFRA.clean-architecture-adoption.md) — БД как адаптер на периферии
- [ADR-DES.API.communication-patterns](ADR-DES.API.communication-patterns.md) — REST + async-события, идемпотентность webhook (T4)
- [ADR-DES.STACK.framework-vs-vs](ADR-DES.STACK.framework-vs-vs.md) — Go/chi, sqlc + Atlas, pnpm-монорепо (T3)
- [Агрегаты, сущности, VO](../ddd/aggregates.md) — агрегаты `Account`, `Learner`, `LearningPlan`, `ResourceCatalog`
- [Каталог доменных событий](../ddd/domain-events.md) — события для outbox и read-моделей
- C4-диаграммы (T6): контейнер «PostgreSQL» — единственное хранилище EduTrack
