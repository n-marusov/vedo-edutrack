# REQ-NFR-security.compliance.permission-matrix

**Приоритет:** P0
**Ключевая функция:** cross-cutting (авторизация)
**Источник:** specs/rbac-matrix.md (T8), ADR-DES.SECURITY.rbac-model (T5); конкретизирует REQ-NFR-security.compliance.role-based-access

**Описание:** Разрешения ролей заданы матрицей «роль × функциональная область × уровень доступа (C/R/U/D) × scope» и применяются deny-by-default на серверной стороне. Матрица покрывает 11 функциональных областей (route-compute, plan-view, plan-manage, progress-track, gap-diagnose, coverage-view, resource-manage, visualization, user-manage, ontology-read, webhook-configure) для всех ролей каталога. Scope-предикат (`own` / `dependents` / `unit` / `all`) проверяется для каждого ресурса; нарушение → `403`, несуществующий ресурс → `404` без раскрытия информации.

**Критерии приёмки:**
- Матрица «роль × функциональная область × действие × scope» реализована (данные `role_permissions`): каждая ячейка матрицы соответствует строке `(role, resource, action, scope)` — 0 ячеек матрицы без реализации (проверка тестом соответствия матрицы и каталога).
- Deny-by-default: доступ разрешён только при явном grant; любая ячейка без grant → `403` (позитивные тесты — только для разрешённых ячеек; негативные тесты — для всех остальных).
- Scope-предикаты: `own` (свои ресурсы), `dependents` (подопечные — parent), `unit` (класс/школа/департамент — teacher/methodologist/school-director/hr-manager), `all` (только admin/ops) — проверяются на серверной стороне для каждого ресурса.
- Обход scope отклоняется: parent → чужие дети, school-director → чужая школа, hr-manager → чужой департамент, teacher → чужой класс — 100% попыток → `403/404` (тест-сьют scope-обхода).
- `ontology-read` — read-only для всех ролей; `webhook-configure` — только platform-integrator (свои подписки) и admin; `ops` — только инфраструктура (JIT), вне продуктовых ресурсов.
- Тест-матрица enforcement в CI покрывает 100% ячеек роль × ресурс (позитив/негатив), все негативные завершаются запретом.

**Связанные артефакты:** [Роль-разрешение матрица](../rbac-matrix.md) (T8), [ADR-DES.SECURITY.rbac-model](../adr/ADR-DES.SECURITY.rbac-model.md) (T5), [REQ-NFR-security.compliance.role-based-access](REQ-NFR-security.compliance.role-based-access.md)
