# REQ-NFR-security.compliance.role-catalog

**Приоритет:** P0
**Ключевая функция:** cross-cutting (авторизация)
**Источник:** specs/vision.md §3.1 (персоны), ADR-DES.SECURITY.rbac-model (T5); дополняет REQ-NFR-security.compliance.role-based-access

**Описание:** Каталог ролей EduTrack — данные-конфигурация (`identity_access.roles`, `role_permissions`), а не типы кода. Контракт строится на **универсальном слое архетипов** (ADR T5): набор архетипов фиксирован моделью и не зависит от продукта; конкретные роли продукта — **экземпляры архетипов** (сид-данные каталога). Движок RBAC и архетипы переиспользуются для любого продукта/контура платформы — EduTrack-роли лишь один из каталогов. Каталог различается для контуров Community и Enterprise (конфигурация, не код); обязательное покрытие ролями-экземплярами гарантируется конфиг-гейтом.

## Архетипы (контрактный слой)

Архетип — паттерн доступа: набор разрешений (permission set) + scope по умолчанию. Роль-экземпляр наследует возможности своего архетипа; допустимы точечные дельты (ограничение, реже — расширение). Полная матрица «архетип × функциональная область × действие» — в `permission-matrix.md`.

| Архетип | Семантика | Scope по умолчанию | Ключевые возможности (функциональные области) |
|---------|-----------|--------------------|-----------------------------------------------|
| `self` | Доступ к собственным данным (профиль, маршрут, прогресс) | `own` | route-compute R, plan-view R, progress-track R, gap-diagnose R, coverage-view R, resource-manage R, visualization R, user-manage R (профиль), ontology-read R |
| `dependents-owner` | Доступ к данным подопечных (дети, подчинённые) | `dependents` | Возможности `self` + user-manage R (подопечные) |
| `staff` | Работа с группами/классами в рамках `unit` | `unit` | plan-manage U (класс/группа), user-manage R (unit) |
| `management` | Управление `unit` (школа/департамент): отчётность, назначения, покрытие | `unit` | plan-manage U, resource-manage U, user-manage R (unit) |
| `integration` | API/service-доступ: ключи, webhooks, лимиты | `own` (API) | route-compute R, webhook-configure CRUD (свои подписки) |
| `admin` | Администрирование системы: пользователи, роли, конфигурация, аудит | `all` | Полный CRUD по продуктовым ресурсам |
| `ops` | Производственный доступ: инфраструктура, деплой, БД (2-person rule, JIT) | `all` (инфраструктура) | Вне продуктовых ресурсов |

> **Инвариант:** набор архетипов фиксирован моделью (`self`, `dependents-owner`, `staff`, `management`, `integration`, `admin`, `ops`) — новые продукты платформы комбинируют архетипы, не добавляют их в модель без решения T5-уровня.

## Начальный каталог EduTrack (сид-данные, иллюстративно)

Конкретные имена — **данные каталога, контракт к ним не привязан**. Роли-экземпляры архетипов, покрывающие персон `vision.md` §3.1 (реальный каталог живёт в `identity_access.roles`, сид-миграции):

| Роль (сид) | Архетип | Персона (vision.md §3.1) | Контур |
|------------|---------|--------------------------|--------|
| `learner` | `self` | Ученик (семейное обучение), Сотрудник (корпорация) | Community + Enterprise |
| `parent` | `dependents-owner` | Родитель (семейное обучение) | Community |
| `teacher` | `staff` | Учитель / методист | Community (Pro) |
| `methodologist` | `staff` + `management` | Учитель / методист | Community + Enterprise |
| `school-director` | `management` | Директор частной школы | Community (Pro) |
| `hr-manager` | `management` | HR-директор / L&D-руководитель | Enterprise |
| `employee` | `self` | Сотрудник (корпорация) | Enterprise |
| `platform-integrator` | `integration` | CTO / CPO EdTech-платформы | Community (Pro) |
| `admin` | `admin` | — (операционная) | Оба |
| `ops` | `ops` | — (операционная) | Оба |

Контрибьюторские роли (методист-контрибьютор, сопровождающий публичной онтологии) — **роли VEDO Hub, не EduTrack**: EduTrack не управляет их правами, предоставляет только read-only `ontology-read` через `ontology-port`.

## Критерии приёмки

- **Покрытие архетипов:** каталог ролей покрывает все архетипы модели — `self`, `dependents-owner`, `staff`, `management`, `integration`, `admin`, `ops`: по каждому архетипу существует ≥ 1 роль-экземпляр в применимом контуре (конфиг-гейт/тест каталога, 0 непокрытых архетипов).
- **Роли — данные:** роли (`identity_access.roles`, `role_permissions`) — экземпляры архетипов, не типы кода: добавление роли/изменение прав не требует перекомпиляции и перезапуска сервиса (сид-миграция), применение — в рамках TTL кэша прав ≤ 15 мин.
- **Полнота контуров:** контур Community — архетипы `self, dependents-owner, staff, management, integration, admin, ops`; контур Enterprise — `self, management, admin, ops` (+ маппинг ролей IdP через Keycloak) — оба каталога валидны, различие — конфигурация, не код; конкретные имена ролей в контуре — данные каталога.
- **Трассируемость персон:** каждая персона `vision.md` §3.1 покрыта ролью-экземпляром архетипа с требуемым набором разрешений из `permission-matrix.md` (100% персон; проверка тестом каталога).
- **Роли контрибьюторов** (методист-контрибьютор, сопровождающий) отсутствуют в каталоге EduTrack — принадлежат VEDO Hub; EduTrack предоставляет read-only доступ через `ontology-port`.

**Связанные артефакты:** [ADR-DES.SECURITY.rbac-model](../adr/ADR-DES.SECURITY.rbac-model.md) (T5 — модель и архетипы), [REQ-NFR-security.compliance.permission-matrix](REQ-NFR-security.compliance.permission-matrix.md) (архетип × область × действие), [REQ-NFR-security.compliance.role-based-access](REQ-NFR-security.compliance.role-based-access.md)
