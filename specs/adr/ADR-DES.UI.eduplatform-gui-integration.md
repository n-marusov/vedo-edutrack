# ADR-DES.UI.eduplatform-gui-integration

**Статус:** ПРИНЯТО
**Дата:** 2026-08-04
**Контекст:** Интеграция GUI EduTrack в интерфейсы EdTech-платформ (M0.3, развитие F6)

EduTrack предоставляет REST API, webhooks, SPARQL-эндпоинт и MCP-сервер для бэкенд-интеграции (`ADR-DES.API.communication-patterns`). Однако EdTech-платформы (онлайн-школы, LMS, корпоративные порталы) ожидают также **GUI-интеграции**: встраивание карты знаний, дашбордов прогресса, конструктора маршрутов и панели ученика непосредственно в собственный интерфейс. Это особенно актуально для B2B-контура Enterprise — корпоративные заказчики хотят видеть EduTrack-функциональность внутри своего портала без переключения между системами.

Продукт — тяжёлая SPA-визуализация (Cytoscape.js, React Flow), и простого REST API недостаточно для интеграции визуальных компонентов: платформам нужен способ встроить живые интерактивные элементы (карта знаний, дашборды, конструктор) в свою HTML-разметку, с единой аутентификацией и бесшовным пользовательским опытом.

**Ключевые драйверы:**
- **Developer Experience (DX)**: интеграция должна требовать минимум кода и конфигурации со стороны EdTech-платформы — подключил скрипт/виджет → получил работающий компонент.
- **Единая аутентификация**: пользователь EdTech-платформы не должен входить в EduTrack отдельно — SSO/OAuth 2.0 поверх существующей сессии платформы (`REQ-NFR-security.compliance.owasp-application-security`, F6.6 SSO/SAML).
- **Изоляция и безопасность**: встроенные компоненты не должны открывать доступ к данным других клиентов (tenant isolation для Enterprise-контура, `REQ-NFR-infra.compliance.community-enterprise-isolation`).
- **Кастомизация внешнего вида**: EdTech-платформа должна иметь возможность адаптировать тему EduTrack под свой бренд (white-labeling, темы light/dark).
- **Интерактивность в реальном времени**: дашборды и карта знаний должны обновляться при изменениях (прогресс, пересчёт маршрута) без перезагрузки страницы.
- **Поддержка мобильных устройств**: встраиваемые компоненты должны корректно работать на мобильных платформах, куда всё больше уходит образовательный контент.
- **Отсутствие vendor lock-in**: интеграция через открытые веб-стандарты (iframe, postMessage, OAuth 2.0, OpenAPI) — платформа может переиспользовать компоненты в любом стеке.

**Требование-источник:**
- `REQ-FR-api.onboarding.landing.explore-value-proposition` — публичная страница как точка входа для EdTech-платформ
- `REQ-NFR-infra.compliance.client-server-web-app` — веб-приложение, клиент-сервер
- `REQ-NFR-infra.compliance.community-enterprise-isolation` — изоляция контуров
- `REQ-NFR-security.compliance.owasp-application-security` — OWASP (Clickjacking, XSS, CSRF)
- `REQ-NFR-ops.compliance.support-sla` — поддержка интеграторов
- `REQ-NFR-ui.compliance.wcag-level` — доступность встраиваемых компонентов

**Решение:**

Принять **многоуровневую стратегию GUI-интеграции** с приоритетом Developer Experience:

### 1. Embeddable iframe-виджеты (приоритетный способ)

EduTrack предоставляет набор **iframe-встраиваемых виджетов** для ключевых компонентов:

| Виджет | Назначение | Размеры |
|--------|------------|---------|
| **Карта знаний** (`/embed/knowledge-map`) | Интерактивная карта знаний с прогрессом ученика | responsive, min 600×400 |
| **Дашборд прогресса** (`/embed/progress`) | План-факт, отклонения, прогноз для одного ученика | responsive, min 400×300 |
| **Панель группы** (`/embed/group-panel`) | Сводка по группе/классу: кто под риском, общий прогресс | responsive, min 800×400 |
| **Конструктор маршрутов** (`/embed/route-builder`) | Визуальный конструктор с графом и временной шкалой | responsive, min 900×500 |
| **Диагностика лакун** (`/embed/gap-diagnosis`) | Карта лакун и корневые причины отставания | responsive, min 600×400 |

**Принципы iframe-виджетов:**
- Каждый виджет — отдельный маршрут (`/embed/*`) с минимальным бандлом (только необходимые библиотеки, без шапки/футера/навигации SPA).
- Размеры — адаптивные через CSS `width: 100%; height: 100%` с опорой на `min-width`/`min-height`.
- Аутентификация — через OAuth 2.0 token, переданный в query-параметре `?token=` (HTTPS-only, короткоживущий, одноразовый). Альтернативно — через postMessage-рукопожатие.
- CORS: заголовки `Content-Security-Policy: frame-ancestors <origin>` для каждого клиента (Enterprise — белый список разрешённых доменов; Community — `*` с ограничением по токену).
- X-Frame-Options: `ALLOW-FROM <origin>` или `Content-Security-Policy` (SAMEORIGIN не подходит для iframe).

### 2. postMessage API для двусторонней связи

Стандартизированный протокол обмена сообщениями между родительским окном (EdTech-платформа) и iframe (EduTrack):

**Сообщения от родителя к iframe:**
| Сообщение | Назначение |
|-----------|------------|
| `{ type: "auth", token: "..." }` | Передать токен аутентификации (альтернатива query-параметру) |
| `{ type: "navigate", path: "/route-builder?learnerId=X" }` | Навигация внутри виджета |
| `{ type: "setTheme", theme: "light"|"dark", overrides: {...} }` | Смена темы / white-labeling |
| `{ type: "resize", width: N, height: M }` | Принудительное изменение размера iframe |
| `{ type: "setLanguage", locale: "ru"|"en" }` | Смена языка интерфейса |

**Сообщения от iframe к родителю:**
| Сообщение | Назначение |
|-----------|------------|
| `{ type: "ready" }` | Виджет загружен и готов к работе |
| `{ type: "heightChanged", height: N }` | Авто-изменение высоты (ресайз по контенту) |
| `{ type: "action", action: "routeClicked", data: {...} }` | Действие пользователя (клик по модулю, открытие модуля и т.п.) |
| `{ type: "error", code: "...", message: "..." }` | Ошибка виджета |

### 3. Client SDK (JavaScript/TypeScript)

Для платформ, которые хотят интеграцию без ручной работы с iframe и postMessage, предоставляется **npm-пакет `@vedo/edutrack-embed`**:

```typescript
import { EmbeddableKnowledgeMap, EmbeddableProgressDashboard } from '@vedo/edutrack-embed';

// Простая интеграция карты знаний
const map = new EmbeddableKnowledgeMap({
  container: document.getElementById('edu-map'),
  token: 'oauth-token',
  learnerId: '123',
  theme: 'light',
  onModuleClick: (module) => console.log('Clicked:', module),
});

// Или через HTML-атрибуты
// <div data-edutrack-widget="knowledge-map" data-token="..." data-learner-id="123"></div>
```

**Возможности SDK:**
- Автоматическое создание и управление iframe
- Типизированные TypeScript-определения всех сообщений и опций
- Авто-подстройка размера под контент
- Управление темой и языком
- Обработка ошибок и повторное подключение при сбое токена
- Zero external dependencies (только 5-10 KB gzip)

### 4. REST API для headless-интеграции

Для платформ, которые хотят построить собственный UI поверх EduTrack (не используя готовые виджеты), предоставляются те же данные через REST API (`ADR-DES.API.communication-patterns`, OpenAPI-first):

- Получение данных карты знаний в JSON-формате (ноды, связи, прогресс)
- Получение данных дашборда (план-факт, отклонения, прогноз)
- Получение диагностики лакун
- Все эндпоинты — с поддержкой CORS и пагинации

### 5. Webhook-уведомления для синхронизации

Для обновления данных в реальном времени без polling: платформа подписывается на webhook-события (`ADR-DES.API.communication-patterns`, §2) и обновляет свои копии данных или перезагружает iframe-виджет при получении уведомления.

### 6. White-labeling и кастомизация

- **Темы**: CSS-переменные (дизайн-токены) виджетов переопределяются через postMessage-сообщение `setTheme` или через CSS-класс на контейнере. Используется конвенция `--edutrack-*` для белого списка кастомизируемых переменных (`ADR-IMPL.UI.pixso-variables-approach`, п.4).
- **Логотип и брендинг**: настраиваемые элементы бренда (лого, цвета, шрифты) через конфигурацию SDK.
- **Локализация**: поддержка RU/EN + ICU через `setLanguage`; добавление нового языка без изменения кода виджета (`REQ-NFR-ops.compliance.i18n-readiness`).

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Только REST API (без виджетов)** | ⚠️ | Платформа строит свой UI — максимальная гибкость, но каждая платформа реализует сложную визуализацию (Cytoscape.js, React Flow) с нуля; дорого (2-6 месяцев на платформу), высокий порог входа; противоречит требованию DX |
| **Web Components (Custom Elements)** | ⚠️ | Нативная встраиваемость без iframe, единый DOM, простота стилизации — но: нет изоляции CSS/JS (конфликты с хост-приложением), сложнее обеспечить tenant isolation; риск XSS через атрибуты. Оставлен как опция пост-MVP при потребности в тесной интеграции без iframe |
| **Микросервис фронтенда (Module Federation / qiankun)** | ⚠️ | Максимальная интеграция (один SPA), но требует, чтобы EdTech-платформа использовала тот же бандлер и версию React; огромные накладные расходы на поддержку разных версий зависимостей; не подходит для разнородных стеков платформ |
| **Только iframe без postMessage** | ⚠️ | Просто, но нет двусторонней связи: iframe не может сообщить родителю о действиях пользователя, нет авто-ресайза, нет управления темой. Минимальный уровень — требуется postMessage |
| **Embedded React компонент напрямую** | ❌ | Платформа должна установить наш npm-пакет с React-зависимостями — конфликты версий React, сложности с бандлингом; неприемлемо для разнородных стеков |

**Последствия:**

*Положительные:*
- **Быстрая интеграция**: EdTech-платформа добавляет один `<iframe>` или npm-пакет → получает работающий интерактивный компонент за 1-2 часа, а не недели.
- **Изоляция и безопасность**: iframe обеспечивает tenant isolation (разные origin), CSS/JS не конфликтуют с хост-приложением; CORS и CSP контролируют разрешённые источники.
- **Единая аутентификация**: OAuth 2.0 + postMessage-рукопожатие — пользователь не покидает платформу.
- **Открытые стандарты**: iframe + postMessage + OAuth 2.0 работают в любом стеке (React, Vue, Angular, vanilla JS, Rails, PHP, WordPress, 1C-Битрикс).
- **Поддержка мобильных**: адаптивные виджеты корректно работают на мобильных устройствах.
- **White-labeling из коробки**: каждая платформа настраивает внешний вид под свой бренд.
- **Реалтайм-обновления**: webhooks + postMessage позволяют обновлять данные без перезагрузки.

*Отрицательные и смягчение:*
- **iframe — не серебряная пуля**: ограничения API (postMessage), сложность с глубокими ссылками внутри iframe, производительность (отдельный документ) → смягчение: документированные ограничения, SDK скрывает сложность postMessage; критичные виджеты (карта знаний) протестированы в iframe на всех целевых платформах.
- **OAuth 2.0 token в URL** (query-параметр) — риск утечки через Referer/logs → смягчение: HTTPS-only, короткоживущие токены (TTL ≤ 15 мин), одноразовые (token binding), альтернатива — postMessage-рукопожатие без URL-токена.
- **Clickjacking-атаки** через iframe → смягчение: CSP `frame-ancestors` для Enterprise (белый список origin), `Content-Security-Policy-Report-Only` для обнаружения неавторизованных встраиваний; для Community — проверка origin в postMessage и валидация токена.
- **Поддержка множества версий SDK** → смягчение: семантическое версионирование SDK, аддитивные изменения; мажорные версии — с периодом поддержки.
- **SEO-ограничения**: контент в iframe не индексируется поисковиками → смягчение: публичная информация (лендинг) не использует iframe; для Enterprise-контура SEO не актуален.

**Связанные артефакты:**
- [ADR-DES.API.communication-patterns](ADR-DES.API.communication-patterns.md) — REST API, webhooks, MCP (F6)
- [ADR-DES.STACK.framework-vs-vs](ADR-DES.STACK.framework-vs-vs.md) — React, Tailwind (T3)
- [ADR-DES.PROCESS.pixso-design-adoption](ADR-DES.PROCESS.pixso-design-adoption.md) — дизайн-токены, темы, white-labeling
- [ADR-IMPL.UI.pixso-variables-approach](ADR-IMPL.UI.pixso-variables-approach.md) — гибридная токен-система (`@theme` / `var(--...)` / `var(--edutrack-*)`), тёмная тема, Developer Guidelines
- [REQ-NFR-infra.compliance.community-enterprise-isolation](../requirements/REQ-NFR-infra.compliance.community-enterprise-isolation.md) — изоляция контуров
- [REQ-NFR-security.compliance.owasp-application-security](../requirements/REQ-NFR-security.compliance.owasp-application-security.md) — Clickjacking, XSS, CSRF
- [REQ-NFR-ops.compliance.support-sla](../requirements/REQ-NFR-ops.compliance.support-sla.md) — поддержка интеграторов