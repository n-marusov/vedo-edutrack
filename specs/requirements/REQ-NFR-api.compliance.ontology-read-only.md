# REQ-NFR-api.compliance.ontology-read-only

**Приоритет:** P0
**Ключевая функция:** F0 (порт онтологии)
**Источник:** vision.md §1.7 (идеология платформы), REQ-FR-api.hub.read-ontology, ADR-DES.API.communication-patterns (T4)

**Описание:** EduTrack — сервис-слой поверх VEDO Hub: читает онтологию через REST API / MCP / SPARQL Hub (read-only) и никогда не хранит и не редактирует онтологии. Онтология (TBox/ABox, версии, форки, merge requests) принадлежит VEDO Hub; EduTrack хранит только кэш релевантного подграфа (in-memory, иммутабелен по `ontologyVersion`). Формальный стык — порт онтологии (F0, ACL) как единственный путь к данным Hub.

**Критерии приёмки:**
- 0 мутирующих вызовов к VEDO Hub из сервиса: ни один адаптер `ontology-port` не содержит операций создания/изменения/удаления онтологии (проверка статическим анализом + тестом — 0 write-операций).
- EduTrack не хранит онтологию в своей БД: в PostgreSQL отсутствуют таблицы онтологии (модули, связи, ФГОС-привязки, ресурсы, истории, концепции) — только кэш подграфа в памяти (иммутабелен по `ontologyVersion`) и read-модели.
- Единственный путь к данным Hub — `ontology-port` (ACL): ни один модуль, кроме адаптеров порта, не обращается к Hub напрямую (архитектурный тест — 0 обходов).
- Версия онтологии фиксируется при каждом чтении и сохраняется в результатах вычислений (воспроизводимость, согласуется с REQ-FR-api.hub.read-ontology и REQ-FR-plan.compute.shortest-path).
- Write-права на онтологию отсутствуют в API EduTrack: нет эндпоинтов создания/редактирования модулей/связей/рамок (0 мутирующих эндпоинтов в OpenAPI-спеке).

**Связанные артефакты:** [REQ-FR-api.hub.read-ontology](REQ-FR-api.hub.read-ontology.md), [ADR-DES.API.communication-patterns](../adr/ADR-DES.API.communication-patterns.md) (T4), [REQ-NFR-api.compliance.ownership-boundary](REQ-NFR-api.compliance.ownership-boundary.md)
