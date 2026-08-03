# Архитектурные решения (ADR) — VEDO EduTrack

В этой директории хранятся записи архитектурных решений (Architecture Decision Records) проекта VEDO EduTrack — сервиса образовательных маршрутов на базе VEDO Hub.

> **Статус:** язык и фреймворки зафиксированы (`ADR-DES.STACK.*`, ПРИНЯТО, T3): Go + chi (бэкенд), React + TypeScript (фронтенд). БД и коммуникационные паттерны зафиксированы (`ADR-DES.DATA.storage-strategy`, `ADR-DES.API.communication-patterns`, ПРИНЯТО, T4): PostgreSQL (единственный datastore MVP, схема-на-модуль, optimistic locking, outbox) + REST/OpenAPI (URL-версионирование) + in-process события/webhooks (идемпотентность по `event_id`). Репозиторная структура и RBAC-модель зафиксированы (`ADR-IMPL.PROCESS.repository-structure`, `ADR-DES.SECURITY.rbac-model`, ПРИНЯТО, T5): монорепо (backend/frontend/specs) + двухслойная RBAC-модель — абстрактный движок (Role/Permission/Scope/Assignment, роли-как-данные) + контрактный слой архетипов (self/dependents-owner/staff/management/integration/admin/ops) + каталог ролей EduTrack как сид-данные (экземпляры архетипов: learner/parent/teacher/methodologist/school-director/hr-manager/employee/platform-integrator/admin/ops) с scope-ограничениями и серверным энфорсментом. E2E-тестирование зафиксировано (`ADR-DES.STACK.e2e-testing-playwright-vs-cypress`, ПРИНЯТО): Playwright — 10 Must-сценариев MVP (M1–M10), тест-пирамида — в `ADR-IMPL.PROCESS.development-tooling` §6. Стратегия генерации Docker-образов зафиксирована (`ADR-DES.INFRA.docker-images-environments`, ПРИНЯТО): dev без сборки (Vite/air в compose), SaaS/Community — backend (distroless) + SPA на nginx unprivileged (CDN), Enterprise on-prem — единый артефакт с embedded SPA (Go embed, M0.3); embed-подход обоснован как целевая форма поставки on-prem (минимум компонентов, nonroot distroless, same-origin `/api`, runtime-конфиг `APP_CONFIG`). ADR уровня `DES` фиксируют инварианты, не зависящие от стека.

## Правила именования файлов

Файлы именуются по шаблону `<ADR-ID>.md`, где `ADR-ID` — уникальный идентификатор решения.

**Формат идентификатора:**
```
ADR-<LEVEL>.<AREA>.<semantic-tag>
```

Где:
- `LEVEL` = `BIZ` | `DES` | `IMPL`
- `AREA` = `API` | `DATA` | `INFRA` | `SECURITY` | `UI` | `PROCESS` | `INTEGRATION` | `STACK` | `OPS` | `DOC`
- `semantic-tag` — короткая англоязычная метка в kebab-case (2–5 слов), отражающая суть выбора

## Структура ADR

Каждый ADR должен содержать следующие обязательные поля:

- **Статус:** `[ЧЕРНОВИК | ПРЕДЛОЖЕНО | ПРИНЯТО | УСТАРЕЛО | ЗАМЕНЕНО]`
- **Дата:** `ГГГГ-ММ-ДД`
- **Контекст:** описание контекста проблемы
- **Требование-источник:** ссылки на файлы или ID требований
- **Решение:** краткое описание выбора (должно соответствовать semantic-tag)
- **Рассмотренные альтернативы:** перечисление рассмотренных вариантов (если применимо)
- **Последствия:** положительные и отрицательные последствия + способы смягчения

## Паттерны semantic-tag

| Паттерн | Описание | Пример |
|---------|----------|--------|
| `-vs-` / `-vs-...-vs-` | Выбор между альтернативами | `postgres-vs-mysql` |
| `-or-` | Равнозначные варианты | `redis-or-memcached` |
| `-tradeoff` | Компромисс между качествами | `latency-vs-consistency-tradeoff` |
| `-adoption` | Внедрение технологии без альтернатив | `kubernetes-adoption` |
| `-mandate` | Вынужденное решение | `gdpr-mandate` |
| `-strategy` / `-approach` / `-pattern` | Выбор подхода | `cache-strategy`, `saga-pattern` |
| `-evolution` / `-migration` | Изменение существующего | `microservices-evolution` |
| `-scope` / `-boundary` | Определение границ | `context-boundary` |

## Правила

1. **Уникальность ID** — каждый ADR-ID должен быть уникальным в рамках проекта
2. **LEVEL соответствует типу решения:** `BIZ` — бизнес-решения, `DES` — проектные, `IMPL` — реализационные
3. **Максимальная длина semantic-tag** — 40 символов
4. **Запрещены пробелы** в ID, только дефисы и точки
5. **Перед созданием нового ADR** проверьте существующие на пересечение темы
6. **Изменения принятых ADR** оформляются через обновление статуса (заменено/устарело) и создание нового ADR

## Специфика EduTrack

- **Сервис-надстройка над VEDO Hub**: EduTrack не хранит онтологии — читает их через API VEDO Hub. ADR не должны дублировать решения VEDO Core (редактор онтологий, версионирование, ABox, Git-модель) — они принимаются в рамках проекта VEDO Core.
- **Два контура развёртывания**: Community (публичная онтология) и Enterprise (корпоративная, изолированная). ADR уровня `INFRA` и `SECURITY` должны учитывать оба контура.
- **Детерминированное ядро + LLM-модуль**: вычисление маршрута — детерминированная функция; LLM используется для генерации с валидацией (истории, проектные идеи, проверочные задания). ADR, затрагивающие LLM, фиксируют модель валидации.
- **Стек — выбран (T3–T5)**: язык (Go + TypeScript), фреймворки (chi + React), БД/коммуникация (PostgreSQL, REST/OpenAPI + события), репозиторная структура (монорепо), RBAC-модель и E2E-тестирование (Playwright, 10 Must-сценариев MVP) зафиксированы в `ADR-DES.STACK.*`, `ADR-DES.DATA.*`, `ADR-DES.API.*`, `ADR-IMPL.PROCESS.*`, `ADR-DES.SECURITY.rbac-model`. `DES`-уровень фиксирует инвариантные решения (границы сервиса, модель данных, API-контракт) независимо от оставшихся IMPL-деталей.

## Связанные артефакты

- [Видение продукта](../vision.md)
- [Прецеденты использования](../use-cases/README.md)
- [Пользовательские истории](../user-stories/README.md)
- [Глоссарий](../glossary.md)
