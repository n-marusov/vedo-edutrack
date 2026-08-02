# REQ-NFR-security.compliance.role-catalog

**Приоритет:** P0
**Ключевая функция:** cross-cutting (авторизация)
**Источник:** specs/vision.md §3.1 (персоны), specs/rbac-matrix.md (T8), ADR-DES.SECURITY.rbac-model (T5); дополняет REQ-NFR-security.compliance.role-based-access

**Описание:** Каталог ролей EduTrack (данные-конфигурация, а не типы кода) содержит полный набор ролей, покрывающий всех персон `vision.md` §3.1: learner, parent, teacher, methodologist, school-director, hr-manager, employee, platform-integrator, admin, ops. Каталог различается для контуров Community и Enterprise (конфигурация, не код); обязательные роли NFR (learner, parent, school-director, methodologist, hr-manager) гарантированы конфиг-гейтом. Контрибьюторские роли (методист-контрибьютор, сопровождающий) живут в VEDO Hub, не в EduTrack.

**Критерии приёмки:**
- Каталог ролей содержит все 10 ролей: `learner`, `parent`, `teacher`, `methodologist`, `school-director`, `hr-manager`, `employee`, `platform-integrator`, `admin`, `ops` — проверяется конфиг-гейтом/тестом каталога (0 недостающих ролей).
- Роли — данные (`identity_access.roles`, `role_permissions`), не типы кода: добавление роли/изменение прав не требует перекомпиляции и перезапуска сервиса (сид-миграция), применение — в рамках TTL кэша прав ≤ 15 мин.
- Контур Community: каталог `learner, parent, teacher, methodologist, school-director, platform-integrator, admin, ops`; контур Enterprise: `learner(=employee), methodologist, hr-manager, admin, ops` (+ маппинг ролей IdP через Keycloak) — оба каталога валидны, различие — конфигурация, не код.
- Персоны vision.md §3.1 трассируются в роли: parent→`parent`, learner→`learner`, учитель/методист→`teacher`+`methodologist`, директор→`school-director`, CTO/CPO→`platform-integrator`, HR→`hr-manager`, сотрудник→`employee` (100% персон покрыто; проверка тестом каталога).
- Роли контрибьюторов (методист-контрибьютор, сопровождающий) отсутствуют в каталоге EduTrack — они принадлежат VEDO Hub; EduTrack предоставляет им только read-only доступ к онтологии через `ontology-port`.

**Связанные артефакты:** [Роль-разрешение матрица](../rbac-matrix.md) (T8), [ADR-DES.SECURITY.rbac-model](../adr/ADR-DES.SECURITY.rbac-model.md) (T5), [REQ-NFR-security.compliance.role-based-access](REQ-NFR-security.compliance.role-based-access.md)
