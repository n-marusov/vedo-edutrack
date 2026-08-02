# ADR-DES.STACK.e2e-testing-playwright-vs-cypress

**Статус:** ПРИНЯТО
**Дата:** 2026-08-03
**Контекст:** Выбор инструмента E2E-тестирования для VEDO EduTrack (M0.1, follow-up к T3)

Стек зафиксирован (`ADR-DES.STACK.language-vs-vs`, `ADR-DES.STACK.framework-vs-vs`, оба ПРИНЯТО, T3): **Go + chi** (бэкенд), **React + TypeScript** (фронтенд), PostgreSQL, REST/OpenAPI. Стратегия проекта — вайбкодинг + компенсация тестированием и документацией. Тест-пирамида зафиксирована в `ADR-IMPL.PROCESS.development-tooling` (§6): unit → integration (testcontainers) → контракты (OpenAPI) → компонентные (Vitest + RTL) → **E2E (10 Must-сценариев MVP) Playwright**. Однако выбор E2E-инструмента не оформлен отдельным обоснованным решением: Playwright упомянут в стеке (T3) и инструментах разработки вскользь, без рассмотрения альтернатив и последствий — настоящий ADR закрывает этот пробел.

Продукт — тяжёлая SPA-визуализация (F4): карта знаний на canvas (Cytoscape.js, F4.1), конструктор маршрутов (React Flow, F4.3), дашборды. Критичные сценарии MVP кросс-слойные: UI → REST API (M8) → PostgreSQL → VEDO Hub. E2E в реальном браузере — верхний уровень стратегии компенсации: он проверяет целостные пользовательские пути, которые не покрываются unit/integration.

**Ключевые драйверы:**
- **e2e-покрытие 100% Must-критериев MVP** (M1–M10, 10 сценариев из `specs/requirements/MVP-ACCEPTANCE-CRITERIA.md`) — `REQ-NFR-process.dev.test-coverage`; урок опыта Knight Capital 2012 (дефект не был пойман e2e-проверкой) указан в источнике NFR.
- **Тяжёлая SPA-визуализация**: графы на canvas/WebGL и интерактивные редакторы — E2E обязан работать в реальном браузере, а не в jsdom; рендеринг асинхронный (layout-воркеры, `Web Workers`) — нужны авто-ожидания.
- **Кросс-слойность**: проверка «клиент ходит только через серверный API» (`REQ-NFR-infra.compliance.client-server-web-app`) и REST-контракты (M8) — E2E-сценарии покрывают цепочку целиком.
- **WCAG 2.1 AA, axe-core = 0 critical** (`REQ-NFR-ui.compliance.wcag-level`) — a11y-проверки в E2E-контуре.
- **i18n RU+EN** (`REQ-NFR-ops.compliance.i18n-readiness`) — ключевые сценарии прогоняются в обеих локалях.
- **CI-бюджет MTR ≤ 2 ч** (`REQ-NFR-infra.compliance.cicd-supply-chain-security`) — параллелизм и шардинг без внешнего платного SaaS.
- **Вайбкодинг**: инструмент с максимальным AI-корпусом и стабильной генерацией тестов (Page Objects, фикстуры, авто-ожидания).

**Требование-источник:**
- `REQ-NFR-process.dev.test-coverage` (e2e покрывает 100% Must-критериев M1–M10)
- `REQ-NFR-process.dev.engineering-gates` (E2E — часть CI-гейтов)
- `REQ-NFR-infra.compliance.client-server-web-app` (клиент → только через API)
- `REQ-NFR-ui.compliance.wcag-level` (axe-core в E2E)
- `REQ-NFR-ops.compliance.i18n-readiness` (RU/EN-прогоны)
- `REQ-NFR-infra.compliance.cicd-supply-chain-security` (пиннинг браузеров, CI-бюджет)
- `specs/requirements/MVP-ACCEPTANCE-CRITERIA.md` (M1–M10)

**Решение:**

Принять **Playwright** (`@playwright/test`) как инструмент E2E-тестирования VEDO EduTrack.

| Аспект | Решение |
|--------|---------|
| **Раннер** | `@playwright/test` (Playwright Test): авто-ожидания (auto-waiting), web-first assertions, trace viewer, video/screenshots, параллельные workers, шардинг `--shard` |
| **Язык** | **TypeScript** — тот же, что и фронтенд; переиспользование типов (`openapi-typescript`) и Page Objects |
| **Браузеры** | **Chromium** — базовый (headless); **Firefox и WebKit** — smoke-прогон ключевых сценариев (кросс-браузерность для Enterprise-контуров, разные окружения) |
| **Область покрытия** | 10 Must-сценариев MVP (M1–M10) + критичные Should-have-пути (панель группы, вебхуки, SPARQL-слой); E2E венчает пирамиду, не заменяет unit/integration/component |
| **Размещение** | `tests/e2e/gui/` (браузерные сценарии) и `tests/e2e/api/` (API-флоу) — корневой `tests/` (pnpm-воркспейс); структура — `ADR-IMPL.PROCESS.repository-structure` §4–5 |
| **API-флоу** | `request` fixture (`APIRequestContext`) — проверки REST API (M8) и подготовка данных (seed через API); проверка «клиент → только через API-шлюз» в тех же тестах |
| **Auth** | JWT: логин через UI или инъекция токена через `storageState` (без повторных логинов в каждом сценарии) |
| **a11y** | `@axe-core/playwright` — WCAG 2.1 AA-проверки в ключевых сценариях |
| **i18n** | Проектные конфиги на RU и EN; ключевые сценарии прогоняются в обеих локалях |
| **CI** | GitHub Actions: `npx playwright install --with-deps` (пиннинг версии), шардинг по workers, `webServer`-конфиг (подъём frontend/backend против dev-стека `make up`), HTML-отчёт и трассы как артефакты |
| **Локально** | `--ui`-режим для разработки и дебага; `webServer` + `reuseExistingServer` (fit с `make dev`); интеграция с `make test`/`make ci` (§11 `ADR-IMPL.PROCESS.development-tooling`) |
| **Визуальные регрессии** | Опционально, точечно (`toHaveScreenshot`) — только стабильные экраны; не в MVP-гейте |
| **Supply-chain** | Пиннинг версии `playwright` и браузеров (`playwright install` с фиксированной версией в CI); обновление — отдельным PR |

**Аргументы за Playwright:**

| # | Аргумент | Обоснование |
|---|----------|-------------|
| 1 | UI + API в одном инструменте | `request` fixture — браузерные и API-флоу в одном раннере: M8 (REST API) и NFR «клиент → только через API» проверяются без второго инструмента |
| 2 | Авто-ожидания и web-first assertions | `expect(locator).toBeVisible()` с retry — стабильность против flaky; критично для canvas-визуализации (Cytoscape.js) и асинхронных обновлений (план vs факт, live-покрытие ФГОС) |
| 3 | Параллелизм и шардинг из коробки | Workers + `--shard` в GitHub Actions без внешнего сервиса (в отличие от Cypress Cloud) — бюджет MTR ≤ 2 ч |
| 4 | Trace viewer / video / screenshots | Отчёт с трассами — быстрая диагностика флейков; для вайбкодинга: дебаг по трассам, а не «глазами» по скриншотам |
| 5 | Один язык с фронтом (TypeScript) | Тот же TS, что и SPA; переиспользование типов из `openapi-typescript` и Page Objects; ниже порог входа |
| 6 | a11y-интеграция | `@axe-core/playwright` — WCAG 2.1 AA (0 critical) в E2E-контуре без отдельного инструмента |
| 7 | Мультибраузерность одним API | Chromium/Firefox/WebKit — один API и один раннер; Enterprise-контуры с разными окружениями |
| 8 | Вайбкодинг: AI-корпус | Playwright — доминирующий E2E-инструмент в обучающих данных ИИ: стабильная генерация тестов, Page Objects, фикстур |
| 9 | Управление окружением | `webServer` + `reuseExistingServer` — подъём frontend/backend из конфига теста; fit с `make up`/`make dev` |
| 10 | Open source, активное развитие | Microsoft-backed, быстрые релизы, крупное сообщество; zero-cost, без платных фич для CI |

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Cypress** | ⚠️ | Зрелый DX и хорошая документация, но: браузерный E2E без полноценных API-фикстур (`cy.request` слабее `request` fixture), мультитаб и нативные события ограничены архитектурой, параллелизм и шардинг — через платный Cypress Cloud (CI-бюджет, supply-chain), AI-корпус меньше Playwright. Отклонён: Playwright закрывает те же задачи нативно, без SaaS |
| **Selenium WebDriver** | ❌ | Устаревший протокол (WebDriver JSON), медленный, нет авто-ожиданий из коробки — flaky; оправдан только для legacy-скриптов (в проекте отсутствуют) |
| **WebdriverIO** | ⚠️ | Гибкий и модульный, но UI-центричный: API-флоу и параллелизм слабее, нужна связка с Mocha/Jasmine и дополнительная конфигурация; меньше fit с нашим стеком |
| **TestCafe** | ❌ | Меньше возможностей (нет trace viewer, слабее параллелизм), медленнее развитие, малый AI-корпус |
| **Vitest + RTL (браузерный режим)** | ⚠️ | Отличный выбор для компонентных тестов (остаётся в пирамиде), но не заменяет E2E: кросс-слойные M1–M10 требуют реального браузера + сервера + PostgreSQL + Hub |
| **Jest + Puppeteer** | ❌ | Puppeteer без раннера: нет авто-ожиданий, отчётов, шардинга — экосистему собирать руками; устаревший паттерн |
| **Katalon Studio / UFT и др.** | ❌ | Проприетарные, тяжёлые, коммерческие лицензии — не fit с open-source-стеком и вайбкодингом |

**Последствия:**

*Положительные:*
- 100% Must-критериев MVP (M1–M10) покрыты браузерными E2E-сценариями — `REQ-NFR-process.dev.test-coverage` выполнен на верхнем уровне пирамиды.
- Стабильность: авто-ожидания + web-first assertions + retries + trace viewer — минимизация времени на дебаг флейков.
- CI-бюджет MTR ≤ 2 ч: параллельные workers и шардинг без внешнего SaaS (в отличие от Cypress Cloud).
- API-флоу в том же раннере (`request` fixture): проверка «клиент ходит только через API-шлюз» (client-server NFR) и подготовка данных без отдельного инструмента.
- Один язык (TypeScript) на фронт и E2E — переиспользование типов и Page Objects, ниже порог входа для команды.
- a11y (axe-core) и i18n (RU/EN) проверяются в том же E2E-контуре — экономия на отдельных инструментах.

*Отрицательные и смягчение:*
- E2E медленнее unit/integration → держим число сценариев минимальным (10 Must + критичные Should), параллелизм/шардинг; E2E не заменяют пирамиду, а венчают её.
- Canvas-рендеринг (Cytoscape.js, React Flow) — риск flaky-селекторов → стабильные локаторы (роли, `data-testid`), авто-ожидания, retries; трассы и видео для диагностики.
- Мультибраузерность (Firefox/WebKit) дороже в CI → базовый Chromium; Firefox/WebKit — smoke-прогон ключевых сценариев (по расписанию или в релизном пайплайне).
- Браузеры — внешние бинарники (supply-chain) → пиннинг версии `playwright` и `playwright install` с фиксированной версией в CI; обновление — отдельным PR.
- Визуальные регрессии потенциально flaky → вне MVP-гейта; точечно для стабильных экранов, при необходимости — pixel-сравнение с порогом.
- Новый движок (не WebDriver) — нет поддержки legacy-браузеров → не релевантно: продукт — современная SPA (Vite), целевые браузеры — свежие версии Chromium-семейства, Safari, Firefox.

**Связанные артефакты:**
- [Стек: фреймворки (chi + React + oapi-codegen)](ADR-DES.STACK.framework-vs-vs.md) — T3, Playwright для e2e в тест-стратегии
- [Инструменты разработки](ADR-IMPL.PROCESS.development-tooling.md) — §6 тест-пирамида, §7 CI/CD, §11 Makefile
- [Репозиторная структура](ADR-IMPL.PROCESS.repository-structure.md) — §5 тестовая структура (E2E: Playwright, 10 Must-сценариев)
- [Дизайн-процесс (Pencil)](ADR-DES.PROCESS.pencil-design-adoption.md) — генерируемые компоненты проходят те же гейты
- `specs/requirements/MVP-ACCEPTANCE-CRITERIA.md` — M1–M10 (10 Must-сценариев)
- `specs/requirements/REQ-NFR-process.dev.test-coverage.md` — e2e-покрытие 100% Must-критериев
