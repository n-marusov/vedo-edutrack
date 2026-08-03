# ADR-DES.INFRA.docker-images-environments

**Статус:** ПРИНЯТО
**Дата:** 2026-08-03
**Контекст:**

VEDO EduTrack поставляется в двух контурах развёртывания — Community (SaaS) и Enterprise (on-premise / private cloud, `REQ-NFR-infra.compliance.community-enterprise-isolation`) — и проходит окружения ЖЦ dev/staging/prod (`REQ-NFR-infra.compliance.environment-isolation`). Архитектурный стиль — модульный монолит: один разворачиваемый артефакт бэкенда `vedo-edutrack` (`ADR-DES.INFRA.modular-monolith-approach`). Стек зафиксирован (T3–T5): Go + chi (бэкенд), React + TypeScript (SPA). M0.2 задокументировал container strategy (`deploy/README.md`) и ввёл две формы поставки SPA — nginx-вариант (SaaS/CDN) и Go-embed (on-prem) — с единым контрактом рантайм-конфига `APP_CONFIG` (`ADR-DES.INFRA.dynamic-config-injection`).

Необходимо зафиксировать решение о том, **какие Docker-образы и как генерируются для каждого окружения/контура**: что собирается в CI, какие рантаймы используются, чем отличаются контуры, и обосновать **embed-подход** (SPA вшивается в бинарник бэкенда) как целевую форму поставки для Enterprise on-prem (M0.3).

**Ключевые драйверы решения:**
- Enterprise on-prem разворачивается силами стороннего заказчика (или с ограниченной поддержкой вендора): один артефакт с минимумом внешних зависимостей радикально снижает стоимость установки, эксплуатации и поддержки (SLA Enterprise ≤ 1 ч, `REQ-NFR-ops.compliance.support-sla`).
- Минимальная поверхность атаки: distroless nonroot без shell (OWASP, периметр 242-ФЗ) — чем меньше компонентов в рантайме, тем меньше CVE для сопровождения (`REQ-NFR-security.compliance.owasp-application-security`).
- Клиент общается с сервисом только через серверный API (same-origin `/api`) (`REQ-NFR-infra.compliance.client-server-web-app`) — SPA и API могут отдаваться одним сервером без CORS и раздельных origin.
- Одна сборка → много окружений, без пересборки (`REQ-NFR-ops.release.deployment-verification`, drift = 0, воспроизводимость).
- SaaS-контур (Community): read-heavy визуализация (F4) и CDN — статику выгодно отдавать с edge/CDN-кэшированием.
- Масштабирование — stateless-реплики за балансировщиком (Traefik), автоскейлинг ≤ 5 мин (`REQ-NFR-ops.performance.scalability`) — применимо к обоим контурам независимо от формы поставки SPA.

**Требование-источник:**
- `REQ-NFR-infra.compliance.community-enterprise-isolation` (изоляция контуров Community/Enterprise)
- `REQ-NFR-infra.compliance.environment-isolation` (изоляция dev/staging/prod)
- `REQ-NFR-infra.compliance.client-server-web-app` (клиент — только через серверный API)
- `REQ-NFR-ops.compliance.support-sla` (тиры поддержки; Enterprise ≤ 1 ч — простая система дешевле обслуживать)
- `REQ-NFR-ops.release.deployment-verification` (drift = 0, воспроизводимость)
- `REQ-NFR-ops.performance.scalability` (stateless, автоскейлинг)
- `REQ-NFR-security.compliance.owasp-application-security` (nonroot, минимальная поверхность)
- `REQ-NFR-infra.compliance.cicd-supply-chain-security` (SBOM, пиннинг, сканер секретов)

**Решение:**

Зафиксировать **матрицу генерации образов по окружениям/контурам** — dev-контур без сборки и два производственных образа, собираемых в CI:

| Окружение / контур | Образ | Источник сборки | Рантайм | Назначение |
|---------------------|-------|-----------------|---------|------------|
| **Dev** | — (сборка не нужна) | `deploy/docker-compose.yml` (сервисы `backend`, `frontend`) | `golang:1.26-alpine` (air hot-reload) + `node:24-alpine` (Vite HMR) | Разработка одной командой `make up`; без production-образов |
| **Staging / SaaS (Community)** | `vedo-edutrack` + `vedo-edutrack-nginx` | `backend/Dockerfile` + `frontend/Dockerfile.nginx` | distroless static nonroot + `nginxinc/nginx-unprivileged:1.27-alpine-slim` | SPA на edge/CDN: статика кэшируется и масштабируется отдельно |
| **Prod / Enterprise (on-prem)** | один артефакт: `vedo-edutrack` с embedded SPA | `backend/Dockerfile` (SPA-стейдж embed-сборки) | distroless static nonroot | single binary + PostgreSQL — минимум компонентов для заказчика |

**Эволюция и текущий статус (M0.2 → M0.3):**
- **M0.2:** `frontend/Dockerfile.embed` — standalone embed-сервер (`frontend/cmd/spa-embed`, distroless nonroot, HEALTHCHECK через подкоманду `health`) как доказательство механизма. Makefile-таргеты: `docker-build` (backend + embed), `docker-build-all` (+ nginx).
- **M0.3:** embed-сервер переносится в бинарник бэкенда — хендлер SPA встраивается в chi-роутер `backend/cmd/vedo-edutrack`; `frontend/Dockerfile.embed` становится **SPA-стейджем** backend-образа (его финальный distroless-стейдж удаляется). Параллельные пути сохраняются только до M0.3.
- **CI:** гейт `docker-build` собирает backend-образ и контролирует размер ≤ 20 МБ (`deploy/ci/docker-build-check.sh`); на main образ пушится в GHCR (`:latest`, `:sha-<commit>`, `deploy/ci/push-images.sh`). С M0.3 backend-образ включает SPA — push-гейт не меняется, меняется состав образа.

**Embed-подход — обоснование:**

| # | Аргумент | Обоснование |
|---|----------|-------------|
| 1 | **Один артефакт для on-prem** | Установка Enterprise = один контейнер + PostgreSQL. Нет отдельного nginx/статик-сервера, который нужно настраивать, поддерживать и патчить; поддержка диагностирует один процесс — SLA ≤ 1 ч достижим (`support-sla`) |
| 2 | **Минимальная поверхность атаки** | Distroless nonroot, без shell и без второго веб-сервера: меньше компонентов в рантайме → меньше уязвимостей для сопровождения; один образ в Supply Chain/SBOM вместо двух (`cicd-supply-chain-security`, OWASP) |
| 3 | **Same-origin `/api`** | SPA и API отдаются одним сервером на одном порту: нет CORS, нет раздельных origin; полностью согласуется с `client-server-web-app` (клиент ходит только через серверный API) |
| 4 | **Динамический runtime-конфиг из коробки** | Embed-сервер генерирует `/config.js` из env при старте (контракт `APP_CONFIG`, `dynamic-config-injection`): одна сборка → все окружения, версия не «протухает», секреты не попадают в UI (OWASP) |
| 5 | **Версия в одном месте** | ldflags-версия бинарника = версия API и SPA одновременно (`vedo-edutrack version`; M0.3: `/healthz` отдаёт обе версии); поддержка on-prem отвечает «какая версия стоит» одной командой |
| 6 | **Health-check без shell** | Встроенная подкоманда `health` для HEALTHCHECK в distroless (нет curl/wget) — тот же паттерн, что и у backend-образа |
| 7 | **Единый паттерн деплоя** | Оба контура — «один деплой» (модульный монолит): контуры различаются конфигурацией и тирами, а не структурой; blue-green/canary через Traefik работают одинаково (`canary-kill-switch`) |

**Граница применения embed:** статика отдаётся Go-сервером — нет CDN/edge-кэширования уровня nginx. Для SaaS-контура (Community, read-heavy визуализация F4) это ограничение существенно, поэтому **nginx-вариант остаётся целевым для SaaS** (CDN-кэширование, независимое масштабирование статики, envsubst-шаблоны на старте контейнера). Для on-prem объёмы трафика несопоставимы, а хешированные immutable-ассеты отдаются через `http.ServeContent` с `Cache-Control` — Go-сервера достаточно.

Таким образом, **embed — не «эксперимент», а целевая форма поставки on-prem Enterprise; nginx — целевая форма поставки SaaS Community.** Оба варианта собираются из одного SPA-бандла и разделяют контракт `APP_CONFIG` — функциональная идентичность гарантируется одним набором E2E-тестов.

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Только nginx для SPA (оба контура)** | ❌ | On-prem: два контейнера (API + статика), два компонента для настройки/патчинга, больше поверхность атаки; противоречит «одному артефакту» modular-monolith ADR и удорожает SLA-поддержку |
| **Только embed (оба контура)** | ⚠️ | SaaS: без CDN/edge-кэширования статики — проигрыш на read-нагрузке визуализации (F4); nginx-вариант даёт CDN и независимый деплой/масштабирование статики |
| **nginx «всё-в-одном» (статик + прокси `/api`) для on-prem** | ⚠️ | Один контейнер, но добавляет прокси-хоп для API и отдельный компонент (конфиг nginx, CVE-обновления); embed делает то же самое в одном процессе без прокси и nginx |
| **Объектное хранилище/CDN (S3) для статики** | ⚠️ | Для Community возможна пост-MVP эволюция nginx-варианта; для Enterprise on-prem недопустима (автономность развёртывания, данные и артефакты в периметре — 242-ФЗ) |
| **Vite dev-сервер в проде** | ❌ | Dev-инструмент: без production-сборки/оптимизаций, без nonroot, без health-check; только для разработки (dev-контур) |

**Последствия:**

*Положительные:*
- Enterprise: установка одной командой (контейнер + PostgreSQL), один образ в SBOM/supply-chain, диагностика одного процесса — SLA ≤ 1 ч реалистичен (`support-sla`).
- Безопасность: distroless nonroot, без shell, без второго веб-сервера — минимальная поверхность атаки (OWASP, 242-ФЗ).
- Динамичность: одна сборка → dev/staging/prod через env (`APP_CONFIG`); версия SPA и API из одного источника (ldflags) — `deployment-verification` соблюдён.
- SaaS: CDN-кэширование статики через nginx-вариант; независимый деплой SPA — кэш CDN не сбрасывается при релизе API.
- Оба варианта SPA собираются из одного бандла → функциональная идентичность гарантируется одним набором E2E-тестов.

*Отрицательные и смягчение:*
- Два Dockerfile для SPA (embed/nginx) — дублирование build-стейджа → общий node-build стейдж в обоих файлах (pnpm `--frozen-lockfile`, пиннинг Node 24 через `.nvmrc`, кэш pnpm-store); расхождение конфигов контролируется гейтом `compose-health` и докеризованной проверкой.
- Обновление SPA в embed = релиз бинарника → для on-prem приемлемо (релизы редки, RTO ≤ 4 ч); версия видна в UI (футер) и через `vedo-edutrack version` — расхождение версий исключено.
- Размер backend-образа растёт на размер SPA-бандла (~1–2 МБ) → гейт ≤ 20 МБ сохраняется (`docker-build-check.sh`), бандл минифицирован (Vite).
- M0.2–M0.3: два параллельных embed-пути (standalone-сервер → встраивание в backend) → временная поддержка обоих до M0.3; `frontend/cmd/spa-embed` покрыт тестами (`server_test.go`), миграция — встраивание того же хендлера.
- Статика без CDN в on-prem → immutable hashed-ассеты + `Cache-Control`; SaaS использует nginx/CDN.

**Связанные артефакты:**
- [Модульный монолит](ADR-DES.INFRA.modular-monolith-approach.md) — один артефакт, on-prem Enterprise
- [Динамический конфиг](ADR-DES.INFRA.dynamic-config-injection.md) — контракт `APP_CONFIG`, embed-генерация `/config.js`
- [Инструменты разработки](ADR-IMPL.PROCESS.development-tooling.md) — §7 CI/CD, §8 развёртывание (образы, health-checks)
- [Container strategy](../../deploy/README.md) — матрица образов, контуры Community/Enterprise
- `backend/Dockerfile`, `frontend/Dockerfile.embed`, `frontend/Dockerfile.nginx`, `frontend/cmd/spa-embed/`, `frontend/nginx.conf`, `frontend/nginx-config.js.template`
- CI/build: `deploy/ci/docker-build-check.sh`, `deploy/ci/push-images.sh`, `deploy/ci/gates.yaml`, `.github/workflows/ci.yml`, Makefile (`docker-build*`)
- C4 Deployment: `specs/c4/deployment-*.md` (dev/saas/enterprise)
