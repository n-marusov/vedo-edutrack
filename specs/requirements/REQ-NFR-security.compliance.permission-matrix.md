# REQ-NFR-security.compliance.permission-matrix

**Приоритет:** P0
**Ключевая функция:** cross-cutting (авторизация)
**Источник:** ADR-DES.SECURITY.rbac-model (T5); конкретизирует REQ-NFR-security.compliance.role-based-access

**Описание:** Разрешения заданы матрицей «архетип × функциональная область × уровень доступа (C/R/U/D) × scope» и применяются deny-by-default на серверной стороне. Матрица покрывает 11 функциональных областей (route-compute, plan-view, plan-manage, progress-track, gap-diagnose, coverage-view, resource-manage, visualization, user-manage, ontology-read, webhook-configure) для всех **архетипов модели** (`self`, `dependents-owner`, `staff`, `management`, `integration`, `admin`, `ops`). Роли-экземпляры каталога наследуют разрешения своего архетипа (сид-данные, см. `role-catalog`); допустимы точечные дельты (ограничение, реже — расширение). Scope-предикат (`own` / `dependents` / `unit` / `all`) проверяется для каждого ресурса; нарушение → `403`, несуществующий ресурс → `404` без раскрытия информации.

## Матрица: архетип × функциональная область × действие

| Архетип | `route-compute` | `plan-view` | `plan-manage` | `progress-track` | `gap-diagnose` | `coverage-view` | `resource-manage` | `visualization` | `user-manage` | `ontology-read` | `webhook-configure` |
|---------|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `self` | **R** | **R** | **—** | **R** | **R** | **R** | **R** | **R** | **R** (профиль) | **R** | **—** |
| `dependents-owner` | **R** | **R** | **—** | **R** | **R** | **R** | **R** | **R** | **R** (подопечные) | **R** | **—** |
| `staff` | **R** | **R** | **U** (группа) | **R** | **R** | **R** | **—** | **R** | **R** (unit) | **R** | **—** |
| `management` | **R** | **R** | **U** | **R** | **R** | **R** | **U** | **R** | **R** (unit) | **R** | **—** |
| `integration` | **R** | **R** | **—** | **R** | **R** | **R** | **—** | **—** | **—** | **R** | **CRUD** (свои) |
| `admin` | **R** | **R** | **CRUD** | **R** | **R** | **R** | **CRUD** | **R** | **CRUD** | **R** | **CRUD** |
| `ops` | **—** | **—** | **—** | **—** | **—** | **—** | **—** | **—** | **—** | **—** | **—** |

> **Принципы матрицы:**
> - `route-compute`, `plan-view`, `progress-track`, `gap-diagnose`, `coverage-view`, `ontology-read` — **read-only для всех продуктовых архетипов**: маршрут — функция (вычисление, не документ), онтология — read-only (граница `ontology-port`, EduTrack не редактирует онтологии).
> - `plan-manage` — только управляющие архетипы (`staff` в рамках группы, `management`) и `admin`; `self`/`dependents-owner` фиксируют план через управляющего (если назначен) или только читают его.
> - `user-manage` — сам субъект (свой профиль, `self`) + подопечные (`dependents-owner`) + `unit`-архетипы в границах scope + `admin` (полный CRUD).
> - `resource-manage` — read для `self`/`dependents-owner`, write для `management`/`admin`.
> - `webhook-configure` — только `integration` (свои подписки) и `admin`; `ops` — только инфраструктура (JIT), вне продуктовых ресурсов.
> - `ops` — **никаких продуктовых разрешений**: производственный доступ (деплой, БД, инфраструктура) с 2-person rule и JIT (`REQ-NFR-ops.access.production-control`), вне матрицы продуктовых ресурсов.

## Критерии приёмки

- Матрица «архетип × функциональная область × действие × scope» реализована (данные `role_permissions`): каждая ячейка архетипа материализуется в строки `(role, resource, action, scope)` для всех ролей-экземпляров архетипа — 0 ячеек матрицы без реализации (проверка тестом соответствия матрицы и каталога).
- Deny-by-default: доступ разрешён только при явном grant; любая ячейка без grant → `403` (позитивные тесты — только для разрешённых ячеек; негативные тесты — для всех остальных).
- Scope-предикаты: `own` (свои ресурсы — `self`), `dependents` (подопечные — `dependents-owner`), `unit` (класс/школа/департамент — `staff`/`management`), `all` (только `admin`/`ops` по JIT) — проверяются на серверной стороне для каждого ресурса.
- Обход scope отклоняется: `dependents-owner` → чужие подопечные, `management` → чужой unit (школа/департамент), `staff` → чужой класс/группа — 100% попыток → `403/404` (тест-сьют scope-обхода).
- `ontology-read` — read-only для всех архетипов; `webhook-configure` — только `integration` (свои подписки) и `admin`; `ops` — только инфраструктура (JIT), вне продуктовых ресурсов.
- Тест-матрица enforcement в CI покрывает 100% ячеек архетип × ресурс (позитив/негатив), все негативные завершаются запретом.

**Связанные артефакты:** [ADR-DES.SECURITY.rbac-model](../adr/ADR-DES.SECURITY.rbac-model.md) (T5), [REQ-NFR-security.compliance.role-catalog](REQ-NFR-security.compliance.role-catalog.md) (роли-экземпляры архетипов), [REQ-NFR-security.compliance.role-based-access](REQ-NFR-security.compliance.role-based-access.md)
