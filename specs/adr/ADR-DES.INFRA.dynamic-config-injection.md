# ADR-DES.INFRA.dynamic-config-injection

**Статус:** ПРИНЯТО
**Дата:** 2026-08-03
**Контекст:** Динамический проброс переменных окружения (версия сборки, URL-ы) во фронтенд и бэкенд (follow-up к M0.2)

Инженерная платформа M0.2 зафиксировала две формы поставки SPA: **nginx-вариант** (`frontend/Dockerfile.nginx`, SaaS/CDN) и **Go-embed** (`frontend/Dockerfile.embed` → M0.3: единый бинарник с embedded SPA, on-prem Enterprise). Бэкенд — единый бинарник `vedo-edutrack` (distroless). Возникла потребность:

- **Версия сборки**: показывать в UI (футер), логировать, использовать в операциях. Vite статически подставляет `import.meta.env.VITE_*` на этапе сборки — одна сборка не может обслуживать несколько окружений без пересборки.
- **URL-ы**: базовый URL API (`/api` через edge), внешний базовый URL (`PUBLIC_BASE_URL` — webhook/абсолютные ссылки), URL VEDO Hub (`VEDO_HUB_API_URL`) — зависят от окружения (dev/staging/prod, контуры Community/Enterprise).
- **Ограничения**: `REQ-NFR-infra.compliance.client-server-web-app` (клиент ходит только через серверный API), `REQ-NFR-infra.compliance.environment-isolation` (0 prod-данных в dev, конфигурация по окружениям), `REQ-NFR-ops.release.deployment-verification` (drift = 0, воспроизводимые конфиги), `REQ-NFR-ops.compliance.support-sla` (один артефакт для on-prem).

Ключевое требование — **динамичность**: одна сборка (образ) разворачивается в разные окружения без пересборки; значения инжектятся при старте контейнера / запуске бинарника.

**Требование-источник:**
- `REQ-NFR-infra.compliance.environment-isolation` (конфигурация по окружениям, 0 prod-данных в dev)
- `REQ-NFR-ops.release.deployment-verification` (drift = 0, воспроизводимость)
- `REQ-NFR-ops.compliance.support-sla` (единый артефакт, версия в одном месте)
- `REQ-NFR-ops.observability.log-level-config` (runtime-конфиг через env — существующий паттерн)
- `REQ-NFR-security.compliance.owasp-application-security` (никаких секретов в рантайм-конфиге/UI)

**Решение:**

Принять **рантайм-инъекцию через `window.APP_CONFIG`** (файл `config.js`, отдаётся/генерируется при старте) для фронтенда и **env-конфиг + ldflags-версию** для бэкенда. Прецедент разрешения: рантайм-значение > build-time fallback (`VITE_*`) > дефолт.

| Слой | Механизм | Версия сборки | URL-ы |
|------|----------|---------------|-------|
| **Frontend (dev, Vite)** | `frontend/public/config.js` — статика из `public/` (Vite отдаёт с диска) | `VITE_APP_VERSION` / "dev" | `apiBaseUrl: "/api"` (прокси Vite) |
| **Frontend (nginx, SaaS)** | `frontend/nginx-config.js.template` → `/etc/nginx/templates/`; entrypoint `20-envsubst-on-templates.sh` генерирует `/usr/share/nginx/html/config.js` из `APP_VERSION`/`APP_ENV`/`API_BASE_URL` при старте контейнера | `APP_VERSION` (ENV-дефолт = ARG VERSION образа) | `API_BASE_URL` (дефолт `/api`), `APP_ENV` |
| **Frontend (Go embed, on-prem)** | embed-сервер генерирует `GET /config.js` из env при рантайме (хендлер имеет приоритет над статикой dist/) | `APP_VERSION` | `API_BASE_URL`, `APP_ENV` |
| **Backend (Go)** | `internal/platform/config` — env-загрузка (12-factor); версия — `-ldflags "-X …/config.Version=…"` | `config.Version` (ldflags, дефолт "dev"), подкоманда `vedo-edutrack version` | `PUBLIC_BASE_URL`, `VEDO_HUB_API_URL`, `DATABASE_URL`, `OTEL_*` (env при запуске) |

**Чтение на фронтенде:** `src/config.ts` — типизированный доступ `appConfig = { version, env, apiBaseUrl }` с приоритетом `window.APP_CONFIG` → `import.meta.env.VITE_*` → дефолт. `index.html` подключает `/config.js` **до** модуля приложения.

**Аргументы за (рантайм-конфиг, а не build-time):**

| # | Аргумент | Обоснование |
|---|----------|-------------|
| 1 | Одна сборка → много окружений | nginx/embed образы собираются в CI один раз; версия/URL-ы меняются переменными окружения при старте — без пересборки и повторного push |
| 2 | Версия не «протухает» | build-time `VITE_*` бейкется в бандл; рантайм-конфиг показывает актуальную версию в каждом окружении |
| 3 | Безопасно | в `window.APP_CONFIG` — только несекретные значения (версия, env, URL); секреты остаются на сервере (OWASP) |
| 4 | Стандартный паттерн | envsubst-шаблоны nginx — официальный механизм официального nginx-образа; Go embed — генерация на лету |
| 5 | Совместимо с контурами | dev (Vite/статик), SaaS (nginx), on-prem (embed) — один контракт `APP_CONFIG` |
| 6 | Бэкенд: 12-factor + ldflags | URL-ы из env (runtime), версия — атрибут бинарника (ldflags) — стандарт для Go |

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Build-time `VITE_*` только** | ⚠️ | Работает, но статично: пересборка на каждое окружение, версия бейкется в бандл. Остаётся как fallback (п.2 решения) |
| **`GET /config.json` + fetch при старте SPA** | ⚠️ | Работает, но добавляет async-загрузку до рендера (flash), гонку с бандлом; `config.js` синхронен и не требует изменения App-кода |
| **Env-переменные в рантайме напрямую (без файла)** | ❌ | Браузер не имеет доступа к env контейнера — нужен механизм доставки (файл/эндпоинт) |
| **Хардкод в коде + пересборка** | ❌ | Противоположность динамичности; риск дрейфа окружений (`deployment-verification`) |
| **K8s ConfigMap + inject** | ⚠️ | Пост-MVP (K8s не в MVP); механизм контракта `APP_CONFIG` останется тем же |

**Последствия:**

*Положительные:*
- Один артефакт SPA разворачивается в dev/staging/prod без пересборки — `environment-isolation` и `deployment-verification` соблюдены.
- Версия видна в UI (футер) и через `vedo-edutrack version` — поддержка on-prem упрощается (`support-sla`).
- Единый контракт `APP_CONFIG` для всех форм поставки — nginx/embed/dev не расходятся.
- Бэкенд-URL-ы (Hub, PUBLIC_BASE_URL, OTel) конфигурируются env — контуры Community/Enterprise из одного бинарника.

*Отрицательные и смягчение:*
- Рантайм-значения не проверяются компилятором → типизация контракта в `src/config.ts` (interface AppConfig) + дефолты.
- envsubst-шаблон чувствителен к `$` в JS → шаблон содержит только `$APP_VERSION/$APP_ENV/$API_BASE_URL`, без template literals; при расширении — фильтр `NGINX_ENVSUBST_FILTER`.
- Non-root nginx не может писать в `/usr/share/nginx/html` → в образе выполнен `chown -R nginx:nginx` (USER root → chown → USER nginx).
- MSYS/Windows мапит `/api` в `-e API_BASE_URL=/api` → в CI/Linux не проявляется; документировано в compose/.env.example.
- Версия в UI может расходиться с версией API → версия API доступна через `vedo-edutrack version`; M0.3: `/healthz` вернёт обе версии.

**Связанные артефакты:**
- [Модульный монолит (Go embed)](ADR-DES.INFRA.modular-monolith-approach.md) — §7–8, on-prem single artifact
- [Инструменты разработки](ADR-IMPL.PROCESS.development-tooling.md) — §8 развёртывание, образы
- [Хранение данных](ADR-DES.DATA.storage-strategy.md) — env-конфиг `DATABASE_URL`
- `frontend/src/config.ts`, `frontend/public/config.js`, `frontend/nginx-config.js.template`, `frontend/Dockerfile.nginx`, `frontend/Dockerfile.embed`, `backend/internal/platform/config/config.go`, `backend/cmd/vedo-edutrack/main.go`, `Makefile`, `deploy/docker-compose.yml`
- [C4 Container](c4/container-overview.md) — Traefik edge (same-origin `/api`)
- [Container strategy](../../deploy/README.md) — контуры Community/Enterprise
