# ADR-IMPL.PROCESS.development-tooling

**Статус:** ПРИНЯТО
**Дата:** 2026-08-02
**Контекст:** Инструменты разработки и полный стек VEDO EduTrack (M0.1, T3 + follow-up)

Язык (`ADR-DES.STACK.language-vs-vs`: Go + TS) и фреймворки (`ADR-DES.STACK.framework-vs-vs`: chi + oapi-codegen + React) зафиксированы. Настоящий ADR фиксирует **инструменты разработки** (линтеры, форматтеры, pre-commit) и **все инженерные решения по стеку**, обсуждавшиеся при выборе: доступ к данным, DI, логирование, аутентификация, тестирование, CI/CD, развёртывание, наблюдаемость. Стратегия проекта: вайбкодинг (AI-assisted) + компенсация качественным тестированием и документацией.

**Требование-источник:**
- `REQ-NFR-process.dev.test-coverage` (покрытие ≥ 90% ядра, mutation testing ≤ 15%)
- `REQ-NFR-process.dev.engineering-gates` (линтеры/ревью/CI-гейты)
- `REQ-NFR-process.dev.code-complexity` (CC ≤ 10)
- `REQ-NFR-infra.compliance.cicd-supply-chain-security` (SBOM, сканер секретов, pinning)
- `REQ-NFR-ops.observability.log-level-config` (JSON-логи, `LOG_LEVEL`)
- `REQ-NFR-ops.observability.distributed-tracing` (OTLP/Jaeger)
- `REQ-NFR-ops.observability.golden-signals-dashboards` (Prometheus)
- `REQ-NFR-ops.release.deployment-verification` (drift = 0)
- `REQ-NFR-ops.release.canary-kill-switch` (canary, kill switch ≤ 5 мин)
- `REQ-NFR-security.compliance.owasp-application-security` (JWT RS256/JWKS, rate limiting)
- `REQ-NFR-api.performance.latency-p95`

**Решение:**

## 1. Линтеры и форматтеры

| Слой | Форматтер | Линтер | Обоснование |
|------|-----------|--------|-------------|
| **Фронтенд (TS/React)** | **Biome** | **Biome** (встроенный linter) | Один Rust-инструмент вместо ESLint + Prettier: formatter + linter + import-sorter; быстрый, zero-config, без плагинов. Фит: вайбкодинг (одна зависимость, меньше конфигов) |
| **Бэкенд (Go)** | **gofmt** (stdlib) | **golangci-lint** (агрегатор: govet, staticcheck, errcheck, revive) | Стандарт Go: gofmt обязателен, golangci-lint — единая точка конфигурации линтеров; покрывает errcheck (явные ошибки — наш стиль) |
| **Markdown/спеки** | markdownlint (опционально) | — | Спеки — важный артефакт (стратегия документации) |

## 2. Git-хуки: Lefthook

Используем **Lefthook** (https://lefthook.dev/) как менеджер git-хуков — быстрый, кросс-платформенный, без Python-зависимости. Хуки переиспользуют проектные установки инструментов (biome из frontend/, gofmt + golangci-lint системные) — без двойного версионирования. Референс: `references/lefthook.md`.

```yaml
# lefthook.yml
assert_lefthook_installed: true
min_version: 2.1.10

pre-commit:
  parallel: true
  jobs:
    - name: biome-check        # frontend: формат + линт + импорты, авто-фикс
      root: frontend/
      glob: "frontend/**/*.{js,jsx,ts,tsx,json,jsonc,css}"
      run: pnpm exec biome check --write --no-errors-on-unmatched --files-ignore-unknown=true {staged_files}
      stage_fixed: true
    - name: gofmt              # backend: формат только staged .go
      glob: "backend/**/*.go"
      run: gofmt -l -w {staged_files}
      stage_fixed: true
    - name: golangci-lint      # backend: анализ изменённых строк целыми пакетами
      glob: "backend/**/*.go"
      run: cd backend && golangci-lint run --new-from-rev=HEAD~1
```

**Правила biome (из референса `references/biome-precommit.md`):**
- **Пиннинг версии**: `npm i -D -E @biomejs/biome` — Biome меняет CLI между релизами; незакреплённая версия ломает хуки/CI
- **Локально**: `biome check --write` по `{staged_files}` (Lefthook передаёт только staged-файлы из `frontend/`, `root:` снимает префикс) — формат + линт + сортировка импортов, безопасные фиксы; `stage_fixed: true` перестейдживает исправленное
- **В CI**: `biome ci .` (read-only, GitHub-аннотации) — CI ловит то, что хуки автофиксили
- **Всегда**: `--no-errors-on-unmatched` (коммит только markdown не падает) + `--files-ignore-unknown=true` (работает с новыми типами файлов)
- **`vcs.useIgnoreFile: true`** в `biome.json` — `.gitignore`-файлы (dist, сборка) не проверяются
- Обновление: `biome migrate` после апгрейда (v1→v2 сменил конфиг-семантику)

**Правила Lefthook (из референса `references/lefthook.md`):**
- Установка: `pnpm add -D lefthook` (одобрение postinstall: `onlyBuiltDependencies` в `.npmrc` + `allowBuilds` в `pnpm-workspace.yaml`); `lefthook install` вешает хуки в `.git/hooks`
- `assert_lefthook_installed: true` — падать громко, если бинарник отсутствует; `min_version` — защита от старых версий
- `parallel: true` — джобы идут конкурентно; `stage_fixed` работает только в `pre-commit`
- Джобы через `{staged_files}` + `glob` — линтеры трогают только изменённые файлы; инструменты, требующие анализа целого пакета/индекса (golangci-lint, biome --staged), гейтятся `glob` без шаблона файлов
- `root: <dir>/` меняет CWD команды и снимает префикс с путей `{staged_files}`; glob-ы считаются от корня git-репозитория

## 3. Доступ к данным и миграции

| Решение | Выбор | Обоснование |
|---------|-------|-------------|
| **Доступ к данным** | `sqlc` | SQL → типизированный Go-код; явный, без магии ORM; чистые границы для модулей Clean Architecture; SQL — AI-дружелюбен. (Spike M0.2: sqlc vs Ent vs GORM) |
| **Миграции** | `Atlas` | Декларативные миграции + **детекция дрейфа схемы** (drift = 0 — NFR deployment-verification); обратимость (откат ≤ 15 мин) |
| **Политика миграций** | Expand → deploy → contract | Custom migration linter: запрет DROP COLUMN/TABLE без 2-этапного deprecated-периода; `pg_dump`-бэкап перед миграцией в CI |

## 4. DI, логирование, i18n

| Слой | Выбор | Обоснование |
|------|-------|-------------|
| **DI** | `wire` (compile-time) | 10 модулей со статическим графом зависимостей: генерация DI на компиляции — быстрее сборки, проще дебаг, без runtime-магии (в отличие от uber-fx); per-command lazy wire для CLI-команд (`ADR-DES.API.cli-interface`) |
| **CLI** | `spf13/cobra` | Де-факто стандарт CLI для Go-монолитов (kubectl, docker, gh); единый бинарник `vedo-edutrack` с подкомандами (server, mcp, migrate, seed, ontology sync, route compute, plan get, gap diagnose, report) — `ADR-DES.API.cli-interface` |
| **Логирование** | `zap` + `otelzapbridge` | Структурированные JSON-логи, `LOG_LEVEL`, `trace_id`/`span_id`/`request_id` (NFR log-level-config), связка с OTel |
| **i18n (бэкенд)** | `go-i18n` + `golang.org/x/text` | Ошибки/уведомления/webhook; основная масса строк — на фронте (react-i18next) |

## 5. Аутентификация

- **JWT RS256 + JWKS** через `lestrrat-go/jwx` (NFR owasp: `/.well-known/jwks.json`) — stateless
- **Короткий TTL (15–60 мин) + refresh-ротация** — инвалидация без серверного состояния
- **Ступенчатая ротация ключей**: новый ключ публикуется в JWKS до истечения старого (перекрытие ≥ max TTL) — старые токены валидны весь срок жизни
- **Token blacklist** — только для logout/отзыва, опциональный Redis-кейс (не обязателен на MVP)
- **Keycloak** — пост-MVP адаптер для Enterprise SSO/SAML (F6.5), за портом identity

## 6. Тестирование (стратегия компенсации)

```
Unit (ядро F1/F2 ≥ 90%)     Go testing + testify
Integration (≥ 70% API)     testcontainers-go (реальный PostgreSQL в CI)
Контракты                   oapi-codegen + OpenAPI-проверки (дрейф = CI-ошибка)
Мутационное (≤ 15% выжило)  gremlins / go-mutesting → SPIKE на M0.2
E2E (10 Must-сценариев MVP) Playwright — tests/e2e/gui (React-флоу) + tests/e2e/api (API-флоу)
Компонентные (фронт)        Vitest + React Testing Library (Domain/Application ≥ 90%, без браузера)
```

## 7. CI/CD

**GitHub Actions** (репо на GitHub):

```
push/PR → lint (biome ci + golangci-lint + gofmt)
        → typecheck (tsc --noEmit)
        → unit → integration (testcontainers) → e2e (playwright) → mutation (spike)
        → coverage-гейт (ядро ≥ 90%)
        → SAST/DAST + SBOM (syft) + сканер секретов (supply-chain NFR)
        → build (distroless) → push
        → deploy (SSH + compose на VPS) → smoke (deployment-verification)
```

| NFR | Механизм |
|-----|----------|
| MTR ≤ 2 ч | Кэш, параллельные джобы |
| engineering-gates | Линт + тесты + coverage + security в каждом PR |
| deployment-verification | Smoke после деплоя + Atlas drift-проверка |
| canary + kill switch ≤ 5 мин | Blue-green через Traefik (см. §8) |
| supply-chain | SBOM (syft), пиннинг, сканер секретов |

## 8. Развёртывание

| Компонент | Решение |
|-----------|---------|
| **Образы** | Docker: Go-бэкенд — distroless (10–20 МБ), SPA — **Go embed** (один артефакт для on-prem) |
| **Dev-env** | docker-compose (одна команда: backend + frontend + PostgreSQL + OTel-стек) — M0.2 |
| **SaaS MVP / on-prem** | docker-compose + **Traefik** (reverse-proxy, blue-green, TLS) |
| **БД-бэкап** | `pg_dump` перед миграциями (автоматически в CI) |
| **K8s** | Пост-MVP: managed K8s (HPA, canary) при росте Community / опциональный Helm для Enterprise-клиентов со своим K8s |
| **Graceful shutdown** | Signal-handling + drain inflight-запросов (контекст с таймаутом) |
| **Health-checks** | Liveness (лёгкий) + readiness (PostgreSQL, JWKS-endpoint) — M0.3 |

## 9. Наблюдаемость

```
OTel (Go SDK + Web SDK) → OTLP → OTel Collector
  → Prometheus (метрики) + Loki (логи) + Tempo (трейсы)
Grafana: дашборды + alerting (P1–P4, шум ≤ 20%) + корреляция (метрика→лог→трейс)
PostgreSQL как data source для продуктовых метрик (NPS, точность прогноза ±10%)
```

- **Sampling**: 100% ошибочных трасс + 10% успешных (NFR distributed-tracing)
- **Redaction PII (152-ФЗ)**: фильтры на атрибуты (без email/ФИО/тел запросов); телеметрия псевдонимизирована
- **Provisioning as-code**: docker-compose + grafana-provisioning + collector-config в репозитории
- **Контуры**: Community — полный стек; Enterprise — минимальный (телеметрия в периметре, 242-ФЗ)

## 10. Кросс-функциональные MVP-требования

| Решение | Статус |
|---------|--------|
| **Feature flags** | MVP: env/config-флаги (LLM/GPU-фичи выключаются на on-prem без пересборки); Flagsmith — пост-MVP |
| **Rate limiting** | Обязательно с MVP: in-memory token bucket на ноду (chi-middleware); Redis — при 2+ репликах (распределённый) |
| **Redis** | До production: распределённый rate limiting + кэш read-моделей (+опц. token blacklist). Подграф — in-memory (иммутабелен по `ontologyVersion`). On-prem Enterprise — опция, не обязательство |
| **CSP-заголовки** | Обязательно с MVP (OWASP): CSP, HSTS, X-Content-Type-Options на SPA |

## 11. Build automation (Makefile)

**Makefile в корне репозитория** — единая точка входа для всех workflow (M0.2: «build automation covers up/down/build/test/lint»):

| Target | Команда | Назначение |
|--------|---------|-------------|
| `make up` / `make down` | `docker compose up -d` / `down` | Dev-окружение одной командой (backend + frontend + PostgreSQL + OTel-стек) |
| `make dev` | compose + hot-reload | Режим разработки (air для Go, Vite dev для React) |
| `make build` | сборка Go-бинарника + фронт | Продакшн-сборка (distroless-образ) |
| `make test` | `go test ./...` + vitest + playwright | Все тесты (unit + integration + фронт + e2e) |
| `make test:e2e` | `playwright test` (tests/e2e/gui + tests/e2e/api) | E2E-сценарии (M1–M10 Must) |
| `make lint` | `golangci-lint run` + `biome ci` | Линт обоих концов |
| `make format` | `gofmt -l -w` + `biome check --write` | Форматирование |
| `make migrate` / `make migrate-down` | `vedo-edutrack migrate up` / `down` (обёртка над Atlas через CLI, `ADR-DES.API.cli-interface`) | Миграции БД (drift = 0) |
| `make hooks` | `lefthook install` | Установка Lefthook git-хуков |
| `make ci` | полный локальный CI-прогон (lint → test → coverage) | Воспроизведение CI локально |

**CLI-обёртки:** инструментальные цели делегируют в единый бинарник `vedo-edutrack` (cobra-дерево, `ADR-DES.API.cli-interface`): `migrate` → `vedo-edutrack migrate`, сиды → `vedo-edutrack seed`, синк онтологии → `vedo-edutrack ontology sync`, отчёты → `vedo-edutrack report`. Makefile остаётся тонкой обёрткой; тулинг воспроизводим в контейнере (тот же образ).

**Правила:** цели идемпотентны; `make ci` — зеркало GitHub Actions (одна команда для локальной проверки и CI); Makefile генерируется/поддерживается `/aif-build-automation`. Windows: GNU Make через Git Bash/WSL (см. альтернативу Taskfile ниже).

### Фронтенд-тулчейн: Vite + pnpm

| Решение | Выбор | Обоснование |
|---------|-------|-------------|
| **Бандлер/dev-server** | `Vite` | Стандарт для React SPA: быстрый dev-server + HMR (hot reload), быстрая production-сборка; fit с Makefile `make dev` (Vite dev для React) |
| **Пакетный менеджер** | `pnpm` | Быстрее и экономнее npm (content-addressable store, симлинки), строгая изоляция зависимостей; монорепо-дружелюбен (Go + React в одном репо) |
| **Версия Node** | `.nvmrc` / volta (fnm) | Фиксация Node-версии для фронта и CI |

Правила: пиннинг зависимостей (lockfile `pnpm-lock.yaml` в репо), `pnpm approve-builds` для postinstall-скриптов, единая версия Node через `.nvmrc` (CI и локально).

### Дизайн-процесс: Pencil (.pen → production-код)

**Дизайн фронтенда ведётся в Pencil** (`.pen`-файлы) с генерацией production-кода — дизайн-процесс как часть тулчейна:

| Аспект | Решение |
|--------|---------|
| **Инструмент дизайна** | Pencil (`.pen`-файлы, design-as-code) — экраны и компоненты проектируются в Pencil, а не «в коде сразу» |
| **Дизайн-система** | Переиспользуемые компоненты + дизайн-токены (variables): цвета, радиусы, типографика, отступы — единый источник для всех экранов |
| **Генерация кода** | Pencil → React/TSX: компоненты мапятся на дизайн-систему (Button, Card, Input…); семантические Tailwind-классы (`bg-primary`, `rounded-md`), токены → `@theme` (Tailwind v4) |
| **Верификация** | Скриншоты секций + layout-проверки (overflow/наложения) на каждом этапе; финальный визуальный контроль |
| **Стиль** | `frontend-design`-скилл — направление эстетики (типографика, цвет, композиция), без «generic AI»-шаблонов |
| **Связь со стеком** | Генерируется React-код (fit с `ADR-DES.STACK.framework-vs-vs`); Tailwind v4 (semantic utilities); Lucide-иконки |

Правила: токены дизайна = источник истины (не хардкодить цвета/радиусы в коде); компоненты переиспользуются, не пересоздаются; каждый экран верифицируется визуально перед генерацией кода; сгенерированный код проходит те же гейты (§1 линтеры, §6 тесты). Дизайн-артефакты (`.pen`) — в репозитории рядом с кодом.

## 12. AI Factory (вайбкодинг-процесс)

Проект ведётся в режиме **вайбкодинга через AI Factory** — это не «инструмент-помощник», а основной процесс разработки, зафиксированный в `.ai-factory/`:

| Слой AI Factory | Роль в процессе |
|-----------------|-----------------|
| **Контекст** | `.ai-factory/DESCRIPTION.md` (стек/скоуп), `ARCHITECTURE.md` (паттерны), `RULES.md` (конвенции), `ROADMAP.md` (вехи), `RESEARCH.md` (решения/разведка) — единый источник контекста для AI |
| **Планирование** | `/aif-plan` — планы (fast/full), привязка к вехам ROADMAP; `specs/` — требования/UC/US/ADR как формализованный вход |
| **Реализация** | `/aif-implement` — пошаговое исполнение планов, чекбоксы, коммит-план; AI генерирует код в стиле Clean Architecture (ADR clean-architecture-adoption) |
| **Качество** | `/aif-verify`, `/aif-review`, `/aif-qa`, `/aif-test-quality` — верификация, ревью, QA; тесты — стратегия компенсации (§6) |
| **Коммиты** | `/aif-commit` — conventional commits по группам Commit Plan |
| **Документация** | `/aif-docs` — README + docs/; ADR как артефакты решения |
| **Самосовершенствование** | `/aif-evolve` — патчи/skill-context: накопление проектных конвенций из опыта |
| **Трассируемость** | `traceability.ttl` — цепочка Vision → UC → FR → ADR → COMP → TEST (правило проекта: любое изменение артефакта требует прохода трассируемости) |

**Последствия для инструментария:**
- AI Factory-артефакты (`.ai-factory/*`) — часть репозитория и деплоя знаний; они обновляются вместе с кодом (план, DESCRIPTION, RULES, RESEARCH).
- ADR/спеки — не «бумага», а **входные данные для AI**: генерация кода опирается на ADR (Clean Architecture, монолит) и спеки (UC/FR/NFR) — поэтому они обязаны быть актуальными и машинно-читаемыми.
- Конвенции языка: `language.ui: ru` (общение), `language.artifacts: en` (артефакты) — из `config.yaml`.
- Вайбкодинг НЕ отменяет качество: AI генерирует, а инженерный контур (тесты §6, линтеры §1, ревью `/aif-review`, CI §7) страхует результат — это и есть зафиксированная стратегия «вайбкодинг + компенсация тестированием и документацией».

## Рассмотренные альтернативы

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **ESLint + Prettier (фронт)** | ⚠️ | Два инструмента + конфиги; Biome заменяет оба одним Rust-бинарём (скорость, zero-config) |
| **uber-fx (DI)** | ⚠️ | Runtime-магия, дольше сборка, сложнее дебаг; wire — compile-time для статического графа из 10 модулей |
| **GORM / Ent (доступ к данным)** | ⚠️ | GORM — магия, слабые границы; Ent — своя DSL; sqlc — явный типизированный SQL |
| **golang-migrate** | ⚠️ | SQL-версии, но без drift-детекции; Atlas даёт drift = 0 (NFR) |
| **Husky + lint-staged** | ⚠️ | Husky один не умеет staged-файлы (нужен lint-staged); Lefthook — единый конфиг для Go+TS хуков без JS-зависимости хуков |
| **Jaeger вместо Tempo** | ⚠️ | Оба принимают OTLP; Tempo — Grafana-native (TraceQL, единый UI) |
| **Webpack / CRA (сборка фронта)** | ❌ | Медленная сборка, устаревший CRA; Vite — быстрее dev-server + HMR, современный стандарт React-SPA |
| **npm / yarn (пакетный менеджер)** | ⚠️ | npm — медленнее установки и больше диска; yarn классик устарел. pnpm — быстрее, экономнее, строгая изоляция |
| **Только ручной процесс / без AI Factory** | ❌ | Стратегия проекта — вайбкодинг через AI Factory (§12); без формализованного контекста (спеки/ADR/RULES) AI генерирует вразнобой — качество падает |
| **LaunchDarkly (feature flags)** | ⚠️ | На MVP достаточно env/config-флагов; внешний сервис — когда понадобится управление вне кода |
| **Taskfile (go-task) / Justfile** | ⚠️ | Кросс-платформенные альтернативы Makefile (YAML, один бинарник, лучше Windows-native). Принят Makefile как де-факто стандарт + `make ci`-зеркало CI; Windows — через Git Bash/WSL. Переход на Taskfile возможен по факту боли (single binary, без GNU Make) |

**Последствия:**

*Положительные:*
- Один hook-менеджер (Lefthook, `lefthook.yml`) покрывает Go + TS хуки — единый вход для вайбкодинга и ревью.
- Biome: одна зависимость вместо ESLint+Prettier, `biome ci` даёт GitHub-аннотации — CI-гейты из коробки (engineering-gates NFR).
- sqlc + wire: явный типизированный код — чистые границы монолита, надёжная AI-генерация.
- Observability as-code: наблюдаемость воспроизводима, телеметрия псевдонимизирована (152-ФЗ).

*Отрицательные и смягчение:*
- Biome требует пиннинга версии и периодического `biome migrate` → зафиксировано в правилах, хуки падают при дрейфе версии.
- gremlins менее зрел, чем Stryker → spike M0.2 + компенсация арх-тестами/property-based.
- Traefik + compose — не K8s: ручное масштабирование до роста → переход на managed K8s по триггеру (зафиксирован).
- Wire генерирует код → сгенерированный `wire_gen.go` в репозитории, гейт на актуальность в CI.
- Makefile требует GNU Make (Windows — Git Bash/WSL) → при реальной боли кросс-платформенности — миграция на Taskfile (зафиксирована как альтернатива).
- AI Factory как основной процесс создаёт зависимость от качества контекста (спеки/ADR/RULES) → смягчение: контекст-артефакты обновляются вместе с кодом, `/aif-verify`/`/aif-review` страхуют, трассируемость обязательна.

**Связанные артефакты:**
- [Язык (Go + TS)](ADR-DES.STACK.language-vs-vs.md)
- [Фреймворк (chi + React + oapi-codegen)](ADR-DES.STACK.framework-vs-vs.md)
- [Модульный монолит](ADR-DES.INFRA.modular-monolith-approach.md)
- [Clean Architecture](ADR-DES.INFRA.clean-architecture-adoption.md)
- Референс: `.ai-factory/references/biome-precommit.md`
- AI Factory-контекст: `.ai-factory/config.yaml`, `DESCRIPTION.md`, `RULES.md`, `ROADMAP.md`
- ADR `DES.DATA.storage-strategy`, `DES.API.communication-patterns` — T4
