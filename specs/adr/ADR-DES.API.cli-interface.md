# ADR-DES.API.cli-interface

**Статус:** ПРИНЯТО
**Дата:** 2026-08-03
**Контекст:**

- Бэкенд — модульный монолит, **один процесс, один артефакт** (`ADR-DES.INFRA.modular-monolith-approach`): единственный бинарник собирается из `cmd/server` (`ADR-IMPL.PROCESS.repository-structure` §4, принцип 4 «точка входа минимальна»).
- Инструментальный тулинг размазан по внешним утилитам: `make migrate` вызывает `atlas` напрямую, сиды и диагностика — будущие скрипты. Всё это **вне единого артефакта**: не воспроизводится в контейнере (distroless), не переиспользует доменную логику, не имеет типов.
- `REQ-FR-api.mcp.server` (F6.6) требует MCP-сервер **по stdio и по SSE/HTTP**: AI-клиенты (Claude Desktop, Cursor и т.п.) спавнят бинарник как subprocess и общаются JSON-RPC через stdin/stdout — нужен **spawnable-режим** того же бинарника.
- Нужен инструмент для **разработки, поддержки и тестирования**: дебаг маршрута конкретного ученика, диагностика лакун, отчёты по FGOS/аттестации, тестирование движка маршрутов без инфраструктуры (Postgres/Hub).
- NFR, которые затрагивает решение: `REQ-NFR-process.dev.test-coverage` и `REQ-NFR-process.dev.engineering-gates` (CLI как тест-инструмент и тестируемый компонент), `REQ-NFR-security.compliance.ops-admin-separation` (будущие админ-команды), `REQ-NFR-ops.compliance.support-sla` (поддержка без UI).

**Требование-источник:**

- `REQ-FR-api.mcp.server` (F6.6) — MCP stdio + SSE/HTTP
- `ADR-IMPL.PROCESS.repository-structure` (§1, §4 — `cmd/`, принцип 4 — минимальная точка входа)
- `ADR-DES.API.communication-patterns` (§7 — MCP-сервер как входной адаптер над Application-слоем)
- `REQ-NFR-process.dev.test-coverage`, `REQ-NFR-process.dev.engineering-gates` — тестирование
- `REQ-NFR-security.compliance.ops-admin-separation` — будущее администрирование
- `REQ-NFR-ops.compliance.support-sla` — поддержка

**Решение:**

Единый бинарник `vedo-edutrack` (Go, сборка из `backend/cmd/vedo-edutrack/`) с деревом подкоманд на **Cobra**. Бинарник — единственный артефакт бэкенда; подкоманды — входные адаптеры над Application-слоем (тот же паттерн, что MCP-сервер в `ADR-DES.API.communication-patterns` §7: «нет второго пути к данным»).

```
vedo-edutrack
├── server          HTTP API + SPARQL + webhooks + MCP(SSE) — долгоживущий процесс
├── mcp             MCP-сервер по stdio для AI-агентов (F6.6) — spawnable subprocess
├── migrate         up / down / validate — миграции Atlas (drift = 0)
├── seed            каталог ролей RBAC + демо-данные
├── ontology sync   копирование подграфа из VEDO Hub (F0.2)
├── route compute   вычисление маршрута (--stub | из БД)
├── plan get        чтение плана / прогресса
├── gap diagnose    диагностика корневой причины отставания
└── report          attestation / coverage — отчёты в файл (batch)
```

**Принципы:**

1. **CLI — входной адаптер**: команды вызывают те же use cases модулей через wire-провайдеры; никакой второй логики и второго пути к данным. Дерево команд — в `backend/internal/cli/` (не в `cmd/` — точка входа остаётся минимальной; не в `platform/` — `platform` не импортирует модули).
2. **Per-command lazy wire**: каждая команда собирает свой минимальный граф (`migrate` = конфиг + Postgres; `server` = весь граф). Старт CLI-команды — сотни миллисекунд, не секунды.
3. **Queries + точечные commands**: CLI в основном читает (дебаг, отчёты); мутирующие команды (`seed`, `ontology sync`) **не запускают событийные каскады сервера** (без эмиссии `RouteRecalculated` и т.п.) — побочные эффекты по явному флагу.
4. **Скриптуемость**: без интерактивных промптов (cron, CI, Kubernetes Job); подтверждения — только флаг `--yes`.
5. **Формат вывода** `--output json|table|csv`: JSON для машинного разбора (поддержка, CI), table — для человека.
6. **Аудит-логирование**: каждая команда пишет структурированную запись (команда, аргументы, результат, actor) через zap — фундамент будущих админ-команд (`ops-admin-separation`).
7. **CLI для тестирования**: `route compute --stub` (движок на стаб-онтологии без Postgres/Hub) — инструмент TDD движка маршрутов («маршрут — функция» проверяется из терминала); сами CLI-команды покрыты unit-тестами с моками портов (`REQ-NFR-process.dev.test-coverage`).
8. **`server` — явная подкоманда**: `vedo-edutrack` без аргументов → help; бинарник не становится сервером случайно (паттерн `cockroach`/`grafana`).

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Два бинарника** (`cmd/server` + `cmd/cli`) | ⚠️ | Чистое разделение ролей, но MCP stdio и без того требует spawnable-режим; два артефакта усложняют поставку (distroless, SBOM, версии, on-prem единый артефакт из `repository-structure` §4), а переиспользование домена одинаково в обоих вариантах |
| **Тонкий CLI-клиент** (ходит в REST, паттерн `kubectl`/`gh`) | ❌ | Требует живой сервер — бесполезен для дебага/поддержки/тестирования движка; дублирует уже запланированный OpenAPI-клиент фронта; MCP stdio так не реализовать |
| **Только Makefile + внешние утилиты** (atlas, psql, скрипты) | ⚠️ | Работает для миграций, но тулинг вне единого артефакта: не воспроизводится в контейнере, нет типизации и переиспользования домена, нет аудита, скрипты дублируют логику |
| **CLI в `cmd/server/main.go`** (один файл, без `internal/cli`) | ❌ | Нарушает принцип минимальной точки входа; дерево команд — отдельный слой, тестируемый изолированно |
| **CLI в `internal/platform/cli`** | ❌ | `platform` запрещено импортировать модули (`repository-structure` §6, принцип 5) — CLI-команды обязаны звать use cases модулей |

**Последствия:**

*Положительные:*
- **Один артефакт**: поставка (distroless, on-prem единый бинарник с Go embed), SBOM, версионирование — одно место.
- **MCP stdio тривиален**: конфигурация AI-клиента — `{"command": "vedo-edutrack", "args": ["mcp"]}`.
- **Переиспользование домена**: движок маршрутов/диагностика/отчёты доступны из терминала — поддержка без UI, тесты без HTTP.
- **Тулинг воспроизводим в контейнере**: миграции/сиды/синк запускаются в том же образе (`exec` / entrypoint).
- **Аудит + скриптуемость с первого дня** — фундамент админ-команд пост-MVP (развитие в сторону администрирования).
- **Makefile становится тонкой обёрткой**: `make migrate` → `vedo-edutrack migrate up`.

*Отрицательные и смягчение:*
- **Размер бинарника растёт** (Cobra + команды) → смягчение: lazy wire (команды не тянут граф сервера), `-ldflags="-s -w"`, команды подключаются по мере готовности модулей.
- **CLI обходит JWT/RBAC** (прямой вызов Application-слоя) → смягчение: CLI — **доверенный операторский инструмент** (требует креды БД, как `psql`), не входит в пользовательскую поверхность; аудит-логирование; будущие админ-команды — под `ops-admin-separation` (2-person rule, JIT).
- **Команды в distroless-контейнере** → смягчение: контейнер хранит миграции и сиды, `docker exec vedo-edutrack migrate up`; при необходимости — отдельная entrypoint-команда в compose.
- **Cobra — новая зависимость** → смягчение: стабильная, де-факто стандарт CLI для Go-монолитов (kubectl, docker, gh); Apache-2.0, без транзитивных конфликтов; фиксируется в `go.mod` с проверкой лицензии (`REQ-NFR-infra.compliance.oss-licensing`).

**Связанные артефакты:**

- [ADR-IMPL.PROCESS.repository-structure](ADR-IMPL.PROCESS.repository-structure.md) — `cmd/vedo-edutrack/` + `internal/cli/`
- [ADR-DES.API.communication-patterns](ADR-DES.API.communication-patterns.md) — MCP-сервер (F6.6) как входной адаптер
- [ADR-IMPL.PROCESS.development-tooling](ADR-IMPL.PROCESS.development-tooling.md) — Makefile-обёртки, Cobra
- [REQ-FR-api.mcp.server](../requirements/REQ-FR-api.mcp.server.md) — MCP stdio/SSE
- [REQ-NFR-security.compliance.ops-admin-separation](../requirements/REQ-NFR-security.compliance.ops-admin-separation.md) — будущие админ-команды
- [REQ-NFR-process.dev.test-coverage](../requirements/REQ-NFR-process.dev.test-coverage.md) — CLI как тест-инструмент
