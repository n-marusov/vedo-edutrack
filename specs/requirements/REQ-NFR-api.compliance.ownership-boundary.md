# REQ-NFR-api.compliance.ownership-boundary

**Приоритет:** P0
**Ключевая функция:** cross-cutting (границы системы)
**Источник:** vision.md §1.7 (идеология платформы), specs/glossary.md §4, RULES.md (не дублировать Hub-ответственность)

**Описание:** Владение данными, вычислениями, API, событиями и развёртыванием строго разделено между VEDO EduTrack и VEDO Hub без дублирования ответственности. EduTrack владеет учебными данными и механиками (learner/plan/progress, маршруты, лакуны, покрытие); Hub владеет онтологией и её жизненным циклом (модули, связи, ФГОС, ресурсы, истории, концепции, версионирование, форки, социальные механики). EduTrack не реализует Hub-ответственность (редактор онтологий, версионирование, ABox, Git-модель, социальный хаб, LLM-извлечение).

**Критерии приёмки:**
- Владение данными: EduTrack хранит профили учеников, планы, прогресс, траектории, диагнозы лакун, отчёты покрытия; Hub хранит онтологию, ФГОС-привязки, ресурсы, истории, проектные идеи, педагогические концепции (0 пересечений: EduTrack не хранит сущности онтологии, Hub не хранит учебные данные EduTrack).
- Владение вычислениями: EduTrack вычисляет маршруты, лакуны, прогнозы, покрытие, подбор ресурсов, рекомендации; Hub обслуживает онтологические запросы, версионирование, форки/merge, LLM-извлечение (0 дублирования вычислений: EduTrack не реализует форки/дифф-мерж/социальные механики).
- Владение API: EduTrack отдаёт API маршрутов/планов/прогресса/покрытия/визуализации, read-only SPARQL-прокси, webhooks, MCP, а также GUI-виджеты (iframe, postMessage, Client SDK) для встраивания в интерфейсы EdTech-платформ; Hub отдаёт REST/MCP/SPARQL онтологии и CRUD/форки/merge (0 дублирования: SPARQL-эндпоинт EduTrack — прокси, не отдельное хранилище).
- Владение событиями: EduTrack эмитирует учебные события (`ModuleMastered`, `GoalChanged`, `RouteRecalculated`, `PlanFixed`, `PlanDeviationDetected`, `StandardDeficitDetected`, ...); Hub эмитирует онтологические события (`OntologyUpdated`, `OntologyContributionSubmitted/Merged`) (0 перекрёстной эмиссии).
- Владение развёртыванием: EduTrack разворачивает свои контейнеры (веб, API-сервер, PostgreSQL) отдельно от Hub; Hub — отдельная платформа со своим SLA (согласуется с REQ-NFR-api.availability.hub-dependency-sla).
- Отсутствие Hub-ответственности в EduTrack: в коде/архитектуре EduTrack отсутствуют редактор онтологий, версионирование онтологии, ABox, Git-модель онтологии, социальный хаб, LLM-извлечение (0 компонентов, проверка архитектурным обзором).

**Связанные артефакты:** [ADR-DES.API.communication-patterns](../adr/ADR-DES.API.communication-patterns.md) (T4), [ADR-DES.UI.eduplatform-gui-integration](../adr/ADR-DES.UI.eduplatform-gui-integration.md) (GUI-интеграция в EdTech-платформы), [REQ-NFR-api.compliance.ontology-read-only](REQ-NFR-api.compliance.ontology-read-only.md), [REQ-NFR-api.availability.hub-dependency-sla](REQ-NFR-api.availability.hub-dependency-sla.md)
