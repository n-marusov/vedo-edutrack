# MVP Acceptance Criteria — VEDO EduTrack

> Единый документ приёмочных критериев MVP с приоритизацией MoSCoW. Каждый критерий ссылается на конкретные UC, FR и US идентификаторы. Источник: `specs/vision.md` §2.5 (MVP-скоуп, строки 500–568) и `.ai-factory/DESCRIPTION.md` (MVP-секция).

**Статус:** принято (M0.0 baseline)
**Версия:** 1.0
**Дата:** 2026-08-02

---

## Сводка MoSCoW

| Категория | Кол-во | Критерии |
|-----------|--------|----------|
| **Must have** | 10 | Вычисление маршрута, диагностика лакун, карта знаний, план vs факт, покрытие ФГОС, REST API, бинарный прогноз, дефициты, аттестационная готовность, фиксация плана |
| **Should have** | 8 | Каталог ресурсов, истории/проектные идеи, панель группы, SPARQL, вебхуки, карта лакун, дашборды, триггеры пересчёта |
| **Could have** | 5 | Педконцепции, сложный прогноз, LMS-коннекторы, подбор по стилю, MCP-сервер |
| **Won't have (MVP)** | 6 | Полные 5 типов связей, каскад при обновлении онтологии, карта качеств, семантическое сравнение, DOI, онбординг-приложение |

**Логирование:** `INFO [mvp] acceptance-criteria: {must: 10, should: 8, could: 5, wont: 6, total: 29}`

---

## Must have (MVP-блокер — без них MVP не выпускается)

| # | Критерий | UC | FR | US | Условие приёмки |
|---|----------|-----|-----|-----|-----------------|
| M1 | Вычисление кратчайшего маршрута с учётом strict-связей | `UC-plan.compute.shortest-path-to-goal` | `REQ-FR-plan.compute.shortest-path` | `US-plan.compute.shortest-path` | Маршрут строится ≤ 5 минут от входа; все strict-пререквизиты соблюдены |
| M2 | Диагностика корневой лакуны | `UC-execute.gap.diagnose-root-cause` | `REQ-FR-execute.gap.diagnose-root-cause` | `US-execute.gap.diagnose-root-cause` | Подъём по strict-связям до первого неосвоенного модуля; каскадное влияние показано |
| M3 | Карта знаний с цветовой индикацией прогресса | `UC-viz.map.view-knowledge-graph-with-progress` | `REQ-FR-viz.map.progress-colors` | `US-viz.map.knowledge-graph`, `US-viz.map.color-progress` | 5 состояний цвета (зелёный/жёлтый/синий/серый/красный); красные пререквизиты с каскадными стрелками |
| M4 | План vs факт | `UC-execute.progress.plan-vs-actual-comparison` | `REQ-FR-execute.progress.plan-vs-actual` | `US-execute.progress.plan-vs-actual` | Плановая дата vs реальная, отклонение в днях, причина |
| M5 | Покрытие ФГОС в реальном времени | `UC-execute.coverage.fgos-coverage-live` | `REQ-FR-execute.coverage.fgos-live` | `US-execute.coverage.fgos-live` | coverage = освоенные-с-привязкой / всего; обновление в реальном времени |
| M6 | Бинарный прогноз «успевает / не успевает» | `UC-execute.forecast.binary-readiness-forecast` | `REQ-FR-execute.forecast.binary-readiness` | `US-execute.forecast.binary-readiness` | Прогноз к контрольной точке; точность ±10% за 30 дней |
| M7 | Фиксация плана на контрольной точке | `UC-plan.fixation.snapshot-plan` | `REQ-FR-plan.fixation.snapshot` | `US-plan.fixation.snapshot` | Снэпшот маршрута с таймлайном; маршрут продолжает пересчитываться независимо |
| M8 | REST API для EdTech: маршруты, прогресс, покрытие | `UC-api.rest.compute-route`, `UC-api.rest.query-progress`, `UC-api.rest.query-coverage` | `REQ-FR-api.rest.compute-route`, `REQ-FR-api.rest.query-progress`, `REQ-FR-api.rest.query-coverage` | `US-api.rest.compute-route`, `US-api.rest.query-progress`, `US-api.rest.query-coverage` | `GET /routes/compute`, `GET /progress/{learner_id}`, `GET /fgos/coverage/{learner_id}`; p95 ≤ 200 мс при 1000 параллельных |
| M9 | Дефициты с приоритетом | `UC-execute.coverage.deficit-list-with-priority` | `REQ-FR-execute.coverage.deficit-list` | `US-execute.coverage.deficit-list` | Незакрытые требования ранжированы (strict > ядро > опциональное) |
| M10 | Аттестационная готовность (базовый отчёт) | `UC-execute.attestation.attestation-readiness-report` | `REQ-FR-execute.attestation.readiness-report` | `US-execute.attestation.readiness-report` | Отчёт за ≤ 1 с: покрытие по доменам, дефициты, критический путь |

---

## Should have (ключевые потребности — включены в MVP-скоуп)

| # | Критерий | UC | FR | US | Условие приёмки |
|---|----------|-----|-----|-----|-----------------|
| S1 | Базовый каталог ресурсов | `UC-resource.catalog.filter-by-format` | `REQ-FR-resource.catalog.bind-to-module`, `REQ-FR-resource.catalog.filter-by-format` | `US-resource.catalog.filter-by-format` | Все типы ресурсов; фильтр по формату/источнику |
| S2 | Истории (50+) и проектные идеи (30+) | `UC-practice.stories.recommend-stories-at-mastery`, `UC-practice.projects.suggest-cross-subject-projects` | `REQ-FR-practice.stories.recommend-at-mastery`, `REQ-FR-practice.projects.suggest-cross-subject` | `US-practice.stories.recommend-at-mastery`, `US-practice.projects.suggest-cross-subject` | Рекомендация в момент освоения через appliesTo/enriches |
| S3 | Панель управления группой учащихся | `UC-viz.panel.group-management-panel` | `REQ-FR-viz.panel.group-management` | `US-viz.panel.group-management` | Мини-карточки, ролевая видимость, сводка группы |
| S4 | SPARQL endpoint (read-only) | `UC-api.sparql.read-only` | `REQ-FR-api.sparql.read-only` | `US-api.sparql.read-only` | Аутентификация, rate limiting, только чтение |
| S5 | Вебхуки `module.mastered`, `plan.deviated` | `UC-api.webhooks.module-mastered`, `UC-api.webhooks.plan-deviated` | `REQ-FR-api.webhooks.module-mastered`, `REQ-FR-api.webhooks.plan-deviated`, `REQ-FR-api.webhooks.idempotency` | `US-api.webhooks.module-mastered`, `US-api.webhooks.plan-deviated` | Идемпотентная доставка |
| S6 | Карта лакун (диагностический вид) | `UC-viz.map.view-gap-diagnostic-map` | `REQ-FR-viz.map.gap-diagnostic-view` | `US-viz.map.gap-diagnostic-view` | Только корневые лакуны + каскадные стрелки |
| S7 | Дашборды ученика / родителя | `UC-viz.dashboard.learner-dashboard`, `UC-viz.dashboard.parent-hr-dashboard` | `REQ-FR-viz.dashboard.learner`, `REQ-FR-viz.dashboard.parent-hr` | `US-viz.dashboard.learner`, `US-viz.dashboard.parent-hr` | Позиция, горизонты, план-факт, покрытие |
| S8 | Триггеры пересчёта: прогресс, смена цели | `UC-plan.compute.recompute-on-progress`, `UC-plan.compute.recompute-on-goal-change` | `REQ-FR-plan.trigger.recompute` | `US-plan.trigger.recompute-progress`, `US-plan.trigger.recompute-goal` | Событие → пересчёт маршрута (каскад) |

---

## Could have (улучшает опыт — приоритет после MVP)

| # | Критерий | UC | FR | US | Условие приёмки |
|---|----------|-----|-----|-----|-----------------|
| C1 | Маршрут с учётом педагогической концепции | — (F1.7) | `REQ-FR-plan.compute.pedagogy-concept` | — | Спиральные витки, проектное группирование (этап 2) |
| C2 | Сложный трехуровневый прогноз (±10%) | — | `REQ-FR-execute.forecast.binary-readiness` (расширение) | — | Успевает / под риском / не успевает (этап 2) |
| C3 | LMS-коннекторы WebTutor / iSpring | — (F6.3) | — | — | Двусторонний обмен с LMS (этап 2) |
| C4 | Подбор ресурсов по стилю обучения | `UC-resource.match.match-resources-to-learner` | `REQ-FR-resource.match.by-style-difficulty` | `US-resource.match.by-style-and-difficulty` | Фильтр по стилю/сложности/длительности (этап 2) |
| C5 | MCP-сервер для AI-агентов | — (F6.6) | `REQ-FR-api.mcp.server` | — | Семантический поиск, навигация по графу (этап 2) |

---

## Won't have (явные исключения MVP — vision.md §2.5 строки 558–568)

| # | Исключение | Источник | Причина |
|---|------------|----------|---------|
| W1 | Все 5 типов связей (MVP: strict + soft + appliesTo) | vision 2.5, строка 560 | Этап 2: enrich, isAnalogousTo |
| W2 | Каскадный пересчёт при обновлении онтологии | vision 2.5, строка 561 | MVP: пересчёт только при прогрессе и смене цели |
| W3 | Карта качеств и покрытие программы воспитания (F5.4) | vision 2.5, строка 566 | Этап 2 |
| W4 | Семантическое сравнение при запросе на слияние | vision 2.5, строка 564 | Зона VEDO Hub |
| W5 | DOI-регистрация | vision 2.5, строка 565 | Зона VEDO Hub |
| W6 | Корпоративные сценарии (онбординг, карьерные треки) | vision 2.5, строка 568 | Отдельное приложение «Вектор Компетенций» |

---

## Входные/выходные критерии вехи M0.0

### Входные критерии (для начала M0.1)
- [x] US-истории определены для plan, execute, resource, viz, practice, api, a11y (47 шт.)
- [x] UC определены для построения маршрута, исполнения плана, визуализации, ресурсов, практики, API, a11y (42 шт.)
- [x] FR и NFR определены со стабильными ID (60 FR + 48 NFR)
- [x] MVP acceptance criteria и MoSCoW задокументированы (данный документ)
- [x] `specs/requirements/README.md` определяет соглашения именования и трассируемости

### Выходные критерии (завершение M0.0)
- [ ] Quality-matrix: 0 критических + 0 высоких лакун в MVP-скоупе (проверяется в Phase 2/3)
- [ ] traceability.ttl: все US → UC → FR цепочки полны, нет сломанных orphans (Task 3.2)
- [ ] ROADMAP.md: M0.0 отмечен как завершённый (Task 3.3)
