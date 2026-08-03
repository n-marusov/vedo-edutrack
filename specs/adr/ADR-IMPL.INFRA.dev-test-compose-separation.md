# ADR-IMPL.INFRA.dev-test-compose-separation

**Статус:** ПРИНЯТО
**Дата:** 2026-08-03

**Контекст:**

E2E-тесты (`tests/e2e/api`, `tests/e2e/gui`) и интеграционные прогоны используют docker-compose стек для окружения. Изначально существовал один стек `deploy/docker-compose.yml` (9 сервисов: backend с air hot-reload, frontend/Vite, postgres, otel-collector, prometheus, loki, tempo, grafana, traefik), и тестовые гейты поднимали именно его (`make up` / `make down`).

Это порождало проблемы:

1. **Тяжесть:** тесты не нуждаются в observability-стеке (otel-collector, prometheus, loki, tempo, grafana) и traefik — каждый их старт добавляет минуты к циклу E2E без пользы.
2. **Риск для данных:** `make down` выполняет `docker compose down --volumes` — при прогоне гейтов против dev-стека можно было удалить данные разработки (том `postgres_data`).
3. **Профиль-гейт frontend:** в dev-стеке сервис `frontend` за `profiles: ["frontend-dev"]` — `make up` его не поднимает, из-за чего `e2e-gui` не мог гарантированно получить SPA на `:5173` (наблюдался `warn`-статус гейта).
4. **Детерминизм:** air hot-reload в dev-режиме пересобирает бинарник при изменении исходников — недетерминированное окружение для тестов.

**Требование-источник:**
- `REQ-NFR-infra.compliance.environment-isolation` (изоляция окружений dev/staging/prod — распространяется и на test)
- `REQ-NFR-process.dev.test-coverage` (E2E-гейты: 100% Must-критериев MVP)
- `REQ-NFR-ops.release.deployment-verification` (воспроизводимость прогонов)

**Решение:**

Разделить dev и test стеки на два независимых compose-проекта:

| Аспект | Dev-стек | Test-стек |
|--------|----------|-----------|
| Файл | `deploy/docker-compose.yml` | `deploy/docker-compose.test.yml` |
| Compose project | `vedo-edutrack` | `vedo-edutrack-test` |
| Сервисы | 9: backend (air), frontend (Vite, profile), postgres, otel-collector, prometheus, loki, tempo, grafana, traefik | 4: postgres, backend (`go run`, без air, телеметрия off), hub-mock, frontend (Vite, **без profile-гейта**) |
| Host-порты | backend `8080`, frontend `5173`, postgres `5432`, hub-mock `8081` | backend `58080`, frontend `55173`, postgres `55432`, hub-mock `58081` (все +50000 от dev — стеки могут работать параллельно) |
| Lifecycle | `make up` / `make down` | `make test-up` / `make test-down` |
| Тома | `postgres_data`, `grafana_data`, … | `postgres_test_data` (изолированы — `test-down --volumes` физически не может задеть dev-данные) |
| Телеметрия | полный observability-стек | `OTEL_SAMPLING_RATIO=0`, экспортёров нет |

Тестовые E2E-гейты (`e2e-gui`, `e2e-api` в `deploy/ci/gates.yaml`) управляют **test-стеком** через `deploy/ci/e2e-run.sh <gui|api>`: probe → `make test-up` → Playwright → `make test-down` (trap на падение; стек гасится только если скрипт сам его поднял). Playwright `webServer` в обоих конфигах отключается флагом `E2E_STACK_MANAGED=1` — стеком управляет скрипт, а не Playwright.

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Один стек, включать/исключать сервисы профилями** | ⚠️ | Профили compose не изолируют тома/проект: `down --volumes` dev-стека всё равно мог удалить общие тома; фронт остаётся profile-гейтом, усложняя `e2e-gui`; команда `up` становится неоднозначной («dev или test?») |
| **Тесты против production-образа (distroless)** | ⚠️ | Правильно для smoke, но медленно для цикла разработки (build образа каждый прогон) и не покрывает Vite-путь GUI-сценариев; сохраняется как будущая эволюция (`e2e` против `vedo-edutrack` образа из CI) |
| **Тесты без compose (go run на хосте)** | ❌ | Ломает воспроизводимость (хостовые версии Go/Node), требует ручного подъёма postgres/hub-mock; compose даёт единый контракт окружения в dev/CI |

**Последствия:**

*Положительные:*
- E2E-цикл быстрее: 4 сервиса вместо 9, без observability и traefik.
- Изоляция данных: `make test-down --volumes` удаляет только `postgres_test_data`; dev-данные неприкосновенны (разные compose-проекты).
- `e2e-gui` получает frontend гарантированно: в test-стеке сервис `frontend` не за profile-гейтом.
- Детерминизм: `go run` без air, телеметрия выключена.
- Dev и test стеки могут работать одновременно (все host-порты разнесены: backend 8080/58080, frontend 5173/55173, postgres 5432/55432, hub-mock 8081/58081).

*Отрицательные и смягчение:*
- Дублирование конфигурации сервисов между двумя compose-файлами (postgres, backend, hub-mock, frontend) → сервисы остаются минимальными и параметризованными общими `deploy/.env`-дефолтами; расхождение контролируется гейтом `compose-health` и ревью.
- `go run` на каждый старт компилирует бинарник (~30–60 с) → приемлемо для тестового цикла; CI дополнительно собирает production-образ для smoke (не конфликтует).
- Два compose-проекта = два набора сетей/томов → осознанная плата за изоляцию, документирована в `deploy/README.md`.

**Связанные артефакты:**
- [Container strategy](../../deploy/README.md) — матрица стеков dev/test
- [Мок VEDO Hub](ADR-DES.INFRA.mock-hub-strategy.md) — hub-mock в test-стеке
- [Образы и окружения](ADR-DES.INFRA.docker-images-environments.md) — dev без сборки
- [Инструменты разработки](ADR-IMPL.PROCESS.development-tooling.md) — §8 развёртывание, §11 Makefile
- `deploy/docker-compose.yml`, `deploy/docker-compose.test.yml`, `deploy/ci/e2e-run.sh`, `deploy/ci/gates.yaml`, Makefile (`test-up`/`test-down`)
- C4 Deployment: `specs/c4/deployment-dev.md`, `specs/c4/deployment-test.md`
