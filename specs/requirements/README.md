# Требования — VEDO EduTrack

В этой директории хранятся формализованные требования к проекту VEDO EduTrack — сервису образовательных маршрутов на базе VEDO Hub.

## Иерархия требований

Требования организованы по нисходящему принципу — от бизнес-целей до атомарных пользовательских историй:

```
vision.md (бизнес-требования, SMART-цели, MVP, roadmap)
  └─ use-cases/ (прецеденты использования)
       └─ requirements/ (функциональные требования — FR, генерируются на основе UC)
            └─ user-stories/ (пользовательские истории, Gherkin)
```

Каждый уровень детализирует предыдущий, сохраняя трассируемость «UC → FR → US». Связи формализованы в `traceability.ttl` (OWL Turtle, корень проекта).

## Формат идентификатора требования

```
REQ-<TYPE>-<domain-or-area>.<qualifier>.<action|attribute>
```

Где:

| Часть | Описание |
|-------|----------|
| **TYPE** | Тип требования: `FR` (функциональное) · `NFR` (нефункциональное) |
| **domain / area** | Для FR — функциональный домен: `plan` · `execute` · `resource` · `viz` · `practice` · `api` · `a11y`. Для NFR — область ограничения: `api` · `security` · `data` · `infra` · `ui` · `ops` · `integration` · `process` · `doc` |
| **qualifier** | Для FR — поддомен (см. [use-cases/README.md](../use-cases/README.md)). Для NFR — категория качества: `performance` · `availability` · `observability` · `compliance` · `maintainability` |
| **action / attribute** | Для FR — семантическая метка в kebab-case (действие). Для NFR — конкретный атрибут качества в kebab-case (ограничение) |

### Примеры FR

| REQ-ID | Описание |
|--------|----------|
| `REQ-FR-plan.compute.shortest-path` | Вычисление кратчайшего пути до цели |
| `REQ-FR-execute.gap.root-cause` | Диагностика корневой лакуны |
| `REQ-FR-viz.map.progress-colors` | Цветовая схема карты знаний |
| `REQ-FR-api.rest.route-endpoint` | REST-эндпоинт вычисления маршрута |

### Примеры NFR

| REQ-ID | Описание |
|--------|----------|
| `REQ-NFR-api.performance.latency-p95` | P95 latency API ≤ 200 мс при 1000 одновременных запросов |
| `REQ-NFR-api.availability.webhook-idempotency` | Идемпотентность вебхуков: дублирование событий не ломает состояние |
| `REQ-NFR-security.compliance.role-based-access` | Ролевая модель: learner / parent / school / methodologist / HR |
| `REQ-NFR-security.compliance.pii-152-fz` | Защита ПДн в рамках 152-ФЗ (Enterprise-контур) |
| `REQ-NFR-data.availability.backup-rpo` | RPO бэкапов ≤ 1 час |
| `REQ-NFR-infra.compliance.community-enterprise-isolation` | Изоляция контуров Community и Enterprise |
| `REQ-NFR-ui.compliance.wcag-level` | Доступность по WCAG 2.1 AA |
| `REQ-NFR-ops.observability.log-level-config` | Уровень логирования через `LOG_LEVEL` |

## Правила оформления

1. **Один файл — одно атомарное требование.**
2. **Имя файла**: `REQ-<TYPE>-<domain-or-area>.<qualifier>.<action-or-attribute>.md`.
3. **Язык**: русский. Технические термины — на английском.
4. **Обязательные поля**:
   - **Приоритет:** `P0` / `P1` / `P2`
   - **Ключевая функция:** ссылка на F1–F6 из `vision.md`
   - **Источник:** ссылка на породивший UC и бизнес-требование/проблему (P1–P22) из `vision.md`
   - **Описание:** суть требования
   - **Критерии приёмки:** измеримые условия выполнения
5. **Трассируемость**: связи FR → UC и US → FR ведутся в `traceability.ttl` (OWL Turtle, корень проекта).

## Приоритеты

- **P0** — критическое: без него MVP невозможен (например, вычисление маршрута, базовая карта знаний)
- **P1** — важное: ключевая пользовательская потребность (например, план-факт, покрытие ФГОС)
- **P2** — желательное: улучшает опыт, но не блокирует выпуск (например, подбор ресурсов по стилю обучения)

## Связанные артефакты

- [Видение продукта](../vision.md) — бизнес-требования, SMART-цели, проблемы
- [Прецеденты использования](../use-cases/README.md) — UC
- [Пользовательские истории](../user-stories/README.md) — US в формате Gherkin
- [Глоссарий](../glossary.md) — термины предметной области
- [Роль-разрешение матрица](REQ-NFR-security.compliance.role-catalog.md) — каталог ролей (архетипы как контракт, роли — экземпляры); [permission-matrix](REQ-NFR-security.compliance.permission-matrix.md) — архетипы × области × CRUD; [ops-admin-separation](REQ-NFR-security.compliance.ops-admin-separation.md) — разделение admin/ops (T8)
- [Граница ответственности](REQ-NFR-api.compliance.ownership-boundary.md) — VEDO Hub ↔ EduTrack; [ontology-read-only](REQ-NFR-api.compliance.ontology-read-only.md) — порт онтологии (T10)
