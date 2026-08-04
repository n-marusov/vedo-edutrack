# ADR-IMPL.UI.pixso-variables-approach

**Статус:** ПРИНЯТО
**Дата:** 2026-08-04 (обновлено: добавлены Developer Guidelines, тёмная тема, white-labeling, нецветовые токены, конвергенция)
**Контекст:** Выбор способа интеграции цветов из Pixso-дизайна в CSS фронтенда

Проект использует Tailwind CSS v4 с `@theme`-токенами (`ADR-DES.STACK.framework-vs-vs`, `ADR-DES.PROCESS.pixso-design-adoption`). В `frontend/src/styles/index.css` определён базовый набор токенов:

```css
@theme {
  --color-primary-50: #eff6ff;
  --color-primary-500: #3b82f6;
  /* ... голубая палитра */
}
```

Эти токены — **плейсхолдеры** (явно указано: «placeholder — populated later via Pixso design process», `ADR-IMPL.PROCESS.development-tooling` §11).

Дизайн публичной посадочной страницы «Дай пять» (M0.3, T14) разработан в Pixso и использует другую цветовую гамму — **оранжево-тёплую** (`--primary-warm: rgba(234, 88, 12, 1)`). Экспорт дизайна в код через Pixso MCP (`design_to_code` + `refine_generated_code`) генерирует CSS-переменные в `:root` с цветами из дизайна.

Возникает ряд вопросов:

**Вопрос 1: Две цветовые системы.** `@theme`-токены проекта (голубые) несовместимы с дизайном лендинга (оранжевый). Замена проектных токенов сломает существующие компоненты (Dashboard, авторизация и др.). Как интегрировать?

**Вопрос 2: Тёмная тема.** Pixso экспортирует тёмную тему через `data-collection-*-mode`. Как синхронизировать её с проектной тёмной темой и переключением в SPA?

**Вопрос 3: Developer Guidelines.** Когда разработчику использовать `bg-[var(--...)]`, когда `bg-primary`, а когда создавать новый токен в `@theme`?

**Вопрос 4: White-labeling для EdTech-платформ.** Встраиваемые iframe-виджеты (`ADR-DES.UI.eduplatform-gui-integration`) должны поддерживать брендирование. Как EdTech-платформа переопределяет токены внешнего вида?

**Вопрос 5: Нецветовые токены.** Pixso экспортирует не только цвета, но и типографику, радиусы, отступы, тени. Как их интегрировать в Tailwind?

**Вопрос 6: Нейминг переменных.** Pixso генерирует имена `--primary-warm`, проектные токены — `--color-primary-500`. Какая конвенция для новых переменных?

**Вопрос 7: Конвергенция.** Описаны 3 фазы (var → @theme → миграция), но нет критериев перехода между фазами.

**Ключевые драйверы:**
- Цветовая гамма лендинга — **источник истины** для посадочной страницы (решение Product Owner)
- Проектные `@theme`-токены — плейсхолдеры, но на них уже ссылаются существующие компоненты
- Процесс экспорта Pixso → код должен быть бесшовным: без ручного маппинга цветов при каждом экспорте
- Тёмная тема light/dark должна работать из коробки (Pixso уже экспортирует обе темы)
- При смене дизайна (ребрендинг, A/B-тест) цвета должны обновляться без изменения кода компонентов
- EdTech-платформы должны кастомизировать внешний вид без копирования кода

**Требование-источник:**
- `REQ-NFR-ui.compliance.wcag-level` — WCAG 2.1 AA, контраст через токены
- `REQ-NFR-process.dev.developer-documentation` — единый словарь, воспроизводимый процесс
- `ADR-DES.PROCESS.pixso-design-adoption` — дизайн-токены как источник CSS
- `ADR-DES.UI.landing-page-design` — дизайн лендинга
- `ADR-DES.UI.eduplatform-gui-integration` — white-labeling для iframe-виджетов

**Решение:**

Принять **подход гибридной токен-системы** с чёткими правилами выбора между `var(--...)` и `@theme`, с формализованной тёмной темой, с конвенцией нейминга и с путём конвергенции.

### 1. Иерархия токенов (что и когда использовать)

В проекте действуют три уровня токенов. Каждый уровень имеет строго определённую область применения:

| Уровень | Источник | Способ использования | Область применения |
|---------|----------|---------------------|--------------------|
| **`@theme`-токены** | Tailwind v4 (`--color-*`, `--font-*`, `--radius-*`, `--spacing-*`) | `bg-primary`, `text-secondary`, `rounded-lg`, `p-4` | Ручные (hand-coded) компоненты: Dashboard, авторизация, shared-компоненты (Button, Card, Input) |
| **Pixso CSS-переменные** (`var(--...)`) | Pixso MCP (design_to_code) → `pixso-variables.css` | `bg-[var(--primary-warm)]`, `text-[var(--foreground)]` | Экспортированные из Pixso компоненты: лендинг, iframe-виджеты, страницы, созданные через design_to_code |
| **CSS-переменные iframe** (`var(--edutrack-*)`) | EdTech-платформа (переопределение через postMessage `setTheme`) | `bg-[var(--edutrack-primary)]` | Встраиваемые iframe-виджеты (`ADR-DES.UI.eduplatform-gui-integration`): белый список переменных для кастомизации |

**Правило выбора уровня:**
```typescript
// Если компонент написан вручную (hand-coded) → используй @theme-токены
<Button className="bg-primary text-white rounded-lg" />

// Если компонент сгенерирован из Pixso → используй var(--...)
<div className="bg-[var(--primary-warm)] text-[var(--foreground)]" />

// Если компонент — iframe-виджет для EdTech-платформы → используй var(--edutrack-*)
<div className="bg-[var(--edutrack-primary)] text-[var(--edutrack-foreground)]" />
```

### 2. Pixso CSS-переменные: структура файла

Цвета из дизайна сохраняются в `frontend/src/styles/pixso-variables.css`:

```css
/* Светлая тема (default) */
:root {
  /* Цвета */
  --primary-warm: #ea580c;
  --primary-warm-hover: #c2410c;
  --foreground: #1c1917;
  --background: #fafaf9;
  --muted: #78716c;
  --border: #e7e5e4;
  --card-bg: #ffffff;
  --accent: #fdba74;

  /* Типографика */
  --font-heading: 800 2.5rem/1.2 'Inter', sans-serif;
  --font-body: 400 1rem/1.5 'Inter', sans-serif;
  --font-muted: 400 0.875rem/1.5 'Inter', sans-serif;

  /* Радиусы */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 20px;
  --radius-xl: 28px;
  --radius-full: 999px;

  /* Отступы */
  --spacing-section: 80px;
  --spacing-card: 24px;
  --spacing-element: 16px;

  /* Тени */
  --shadow-card: 0 12px 48px -16px rgba(124, 45, 18, 0.15);
  --shadow-elevated: 0 32px 4px -16px rgba(0, 0, 0, 0.15);
}

/* Тёмная тема */
[data-theme="dark"] {
  --primary-warm: #fdba74;
  --primary-warm-hover: #fb923c;
  --foreground: #fafaf9;
  --background: #1c1917;
  --muted: #a8a29e;
  --border: #292524;
  --card-bg: #292524;
  --accent: #ea580c;

  --shadow-card: 0 12px 48px -16px rgba(0, 0, 0, 0.4);
  --shadow-elevated: 0 32px 4px -16px rgba(0, 0, 0, 0.4);
}
```

### 3. Тёмная тема: структура и переключение

**Единый атрибут `data-theme`** на `<html>` управляет тёмной темой для **всех** компонентов — и `@theme`, и Pixso, и iframe-виджетов.

```css
/* @theme-токены для тёмной темы */
@theme {
  --color-primary-50: #eff6ff;
  --color-primary-500: #3b82f6;
}

/* Data-theme переопределяет @theme */
[data-theme="dark"] {
  --color-primary-50: rgba(59, 130, 246, 0.2);
  --color-primary-500: #60a5fa;
}
```

**Переключение темы** — через единый хук `useTheme` (Zustand store или React Context):

```typescript
// frontend/src/store/useTheme.ts
type Theme = 'light' | 'dark' | 'system';

function useTheme() {
  const [theme, setTheme] = useState<Theme>('system');

  useEffect(() => {
    const applied = theme === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
      : theme;
    document.documentElement.setAttribute('data-theme', applied);
  }, [theme]);

  return { theme, setTheme };
}
```

**Приоритет переопределения:**
1. `localStorage` (выбор пользователя) — выше всех
2. `prefers-color-scheme` (системная тема) — если `data-theme` не установлен
3. `data-theme` на `<html>` (программное переключение)

**Для iframe-виджетов** тёмная тема наследуется от родительского окна через postMessage:
```
Родитель → iframe: { type: "setTheme", theme: "dark" }
```

### 4. White-labeling для EdTech-платформ

EdTech-платформа может кастомизировать встраиваемые виджеты через **переопределение `--edutrack-*` переменных**. Контракт белого списка:

| Переменная | Назначение | Значение по умолчанию |
|------------|------------|----------------------|
| `--edutrack-primary` | Основной цвет бренда | `var(--primary-warm)` |
| `--edutrack-primary-hover` | Ховер основного цвета | `var(--primary-warm-hover)` |
| `--edutrack-background` | Фон виджета | `var(--background)` |
| `--edutrack-foreground` | Цвет текста | `var(--foreground)` |
| `--edutrack-border` | Цвет границ | `var(--border)` |
| `--edutrack-radius` | Радиус скругления | `var(--radius-md)` |
| `--edutrack-font-family` | Шрифт (безопасный список) | `'Inter', sans-serif` |

**Механизм:** SDK `@vedo/edutrack-embed` принимает опцию `themeOverrides` → передаёт через postMessage → iframe применяет в CSS:

```typescript
new EmbeddableKnowledgeMap({
  token: '...',
  learnerId: '123',
  themeOverrides: {
    '--edutrack-primary': '#7c3aed',    // фиолетовый бренд платформы
    '--edutrack-radius': '8px',
  },
});
```

В CSS виджета переменные используются с fallback-цепочкой:
```css
.widget-button {
  background: var(--edutrack-primary, var(--primary-warm));
  border-radius: var(--edutrack-radius, var(--radius-md));
}
```

### 5. Нейминг переменных: конвенция

| Категория | Префикс | Пример | Источник |
|-----------|---------|--------|----------|
| Цвета `@theme` | `--color-*` | `--color-primary-500`, `--color-secondary` | Tailwind v4 slot |
| Шрифты `@theme` | `--font-*` | `--font-sans`, `--font-heading` | Tailwind v4 slot |
| Радиусы `@theme` | `--radius-*` | `--radius-md`, `--radius-lg` | Tailwind v4 slot |
| Pixso-цвета | `--[семантика]` | `--primary-warm`, `--foreground`, `--muted` | Исход из назначения, не из страницы |
| Pixso-типографика | `--font-*` | `--font-heading`, `--font-body`, `--font-muted` | Шорткат для текстовых стилей |
| Pixso-радиусы | `--radius-*` | `--radius-sm`, `--radius-xl`, `--radius-full` | Соответствие шкале Tailwind |
| Pixso-отступы | `--spacing-*` | `--spacing-section`, `--spacing-card` | Семантические, не арифметические |
| Pixso-тени | `--shadow-*` | `--shadow-card`, `--shadow-elevated` | Семантические (назначение, а не размер) |
| iframe white-label | `--edutrack-*` | `--edutrack-primary`, `--edutrack-radius` | Закреплённый контракт SDK |

**Правила именования:**
- **Семантика, а не страница**: `--primary-warm`, а не `--landing-orange`. Одна переменная переиспользуется между страницами одного дизайна.
- **Назначение, а не значение**: `--foreground`, `--muted`, а не `--text-black`, `--text-gray-500`.
- **Шкала Tailwind для радиусов**: `--radius-sm` (8px), `--radius-md` (12px), `--radius-lg` (20px), `--radius-xl` (28px), `--radius-full` (999px).
- **Без генерации для каждой страницы**: если новая страница использует те же цвета — переиспользовать существующие переменные.

### 6. Нецветовые токены из Pixso

Pixso экспортирует не только цвета, но и другие дизайн-токены. **Все они помещаются в тот же `pixso-variables.css`** и используются через `var(--...)`.

**Типографика:**
```css
--font-heading: 800 2.5rem/1.2 'Inter', sans-serif;
--font-body: 400 1rem/1.5 'Inter', sans-serif;
--font-muted: 400 0.875rem/1.5 'Inter', sans-serif;
```

Использование: `style="font: var(--font-heading)"` или через Tailwind с arbitrary value:
```tsx
<h1 className="text-[length:var(--font-heading)]">...</h1>
```
> **Ограничение**: `font` shorthand не поддерживается Tailwind utility напрямую. Для типографики из Pixso используйте `style`-проп или отдельные утилиты (`text-2xl font-bold leading-tight`).

**Радиусы и отступы — напрямую в Tailwind через arbitrary values:**
```tsx
<div className="rounded-[var(--radius-xl)] p-[var(--spacing-card)]" />
```

**Тени:**
```tsx
<div className="shadow-[var(--shadow-card)]" />
```

**Правило**: если дизайн-токен повторяется в 3+ компонентах одной страницы — вынести в `pixso-variables.css`. Если используется однократно — inline `style` (без ущерба для темизации).

### 7. Путь конвергенции: триггеры и критерии

| Фаза | Состояние | Триггер перехода |
|------|-----------|------------------|
| **Phase 1 (текущая)** | Pixso-компоненты: `var(--...)` в `pixso-variables.css`. Hand-coded компоненты: `@theme` (голубая палитра) | — |
| **Phase 2 (унификация палитры)** | Проект принимает единую цветовую палитру (оранжевую/любую другую). Pixso-переменные переносятся в `@theme`: `@theme { --color-primary: #ea580c }`. `var(--primary-warm)` → `bg-primary`. `pixso-variables.css` сокращается до минимума. | **(1)** Решение Product Owner о фиксации брендовой палитры **И** **(2)** появление 3+ страниц/компонентов, использующих оранжевую тему |
| **Phase 3 (полная миграция)** | Все компоненты (включая Dashboard, auth) переведены на единую палитру. `pixso-variables.css` удалён. Единый `@theme` — единственный источник цветов. | **(1)** Редизайн существующих экранов под брендовую палитру **И** **(2)** выделение отдельного спринта на миграцию <br>**(Блокирующие критерии:)** все интеграционные тесты проходят; visual regression тесты подтверждают соответствие макетам |

**Пока триггеры не выполнены — Phase 1 остаётся в силе. Не форсировать конвергенцию.**

### 8. Developer Guidelines (свод правил)

| Ситуация | Делать | Не делать |
|----------|--------|-----------|
| **Пишу новый компонент вручную** (Button, Card, Dashboard) | Использовать `@theme`-токены: `bg-primary`, `text-secondary`, `rounded-lg` | Не использовать `var(--...)` — нарушает консистентность ручных компонентов |
| **Экспортирую из Pixso** (design_to_code) | Сохранить `var(--...)` как есть. Поместить в `pixso-variables.css` если переменная новая | Не переводить в `@theme` вручную — дождаться Phase 2 |
| **Переиспользую переменную** между ручным и Pixso-компонентом | Добавить токен в `@theme` и использовать в обоих местах | Не дублировать одну сущность и в `@theme`, и в `pixso-variables.css` |
| **Добавляю iframe-виджет** для EdTech-платформы | Использовать `var(--edutrack-*)` с fallback на `var(--primary-warm)` | Не хардкодить цвета бренда платформы |
| **Меняю дизайн** (ребрендинг, A/B-тест) | Обновить значения в `pixso-variables.css` или `@theme` — код компонентов не трогать | Не заменять `var(--x)` на новые значения в каждом компоненте |
| **Нужна новая CSS-переменная** | Определить назначение → выбрать префикс из конвенции (п.5) → добавить в `pixso-variables.css` или `@theme` | Не называть по странице (`--landing-orange`) — назвать по семантике (`--primary-warm`) |
| **Работаю с тёмной темой** | Использовать `[data-theme="dark"]` — добавлять переменные под этот селектор | Не писать `@media (prefers-color-scheme: dark)` отдельно — централизованный механизм важнее |

**Рассмотренные альтернативы:**

| Альтернатива | Оценка | Причина отклонения |
|--------------|--------|--------------------|
| **Замена проектных `@theme`-токенов на цвета из Pixso** | ❌ | Сломает все существующие компоненты (Dashboard, auth, etc.), которые используют голубую `--color-primary-500`. Требует массового рефакторинга без немедленной выгоды. |
| **Добавление `landing-` префикса в `@theme` (гибрид)** | ⚠️ | Даёт Tailwind-утилиты (`bg-landing-primary`), но: (1) требует ручного маппинга каждого цвета при каждом экспорте, (2) Pixso не генерирует `@theme`-токены автоматически, (3) два лендинга — недостаточная нагрузка для оправдания сложности. Отложено до появления 2+ страниц в оранжевой теме. |
| **Только inline-стили (без CSS-переменных)** | ❌ | Невозможность темизации, дублирование цветов в каждом элементе, сложность обновления дизайна. |
| **CSS-переменные в `@theme` через `@property`** | ⚠️ | Tailwind v4 `@theme` не поддерживает кастомные CSS-переменные напрямую — только предопределённые слоты (`--color-*`, `--font-*`). Попытка использовать `@theme` для произвольных переменных не даст преимуществ перед прямыми `var(--...)`. |
| **Единая тема (один `data-theme` для всех)** | ✅ **ПРИНЯТО** | Описан в п.3: `data-theme` на `<html>` управляет и `@theme`, и Pixso, и iframe — единый источник темы |
| **Раздельные переключатели темы** (Pixso-компоненты отдельно, @theme-компоненты отдельно) | ❌ | Пользовательский опыт страдает: половина страницы в светлой теме, половина в тёмной; сложность поддержки двух переключателей |

**Последствия:**

*Положительные:*
- **Бесшовный экспорт**: Pixso MCP генерирует CSS-переменные автоматически — никакого ручного маппинга цветов.
- **Точность дизайна**: цвета в коде идентичны цветам в макете — zero translation loss.
- **Единая тёмная тема**: один атрибут `data-theme` управляет всеми компонентами — и ручными, и из Pixso, и iframe-виджетами.
- **White-labeling из коробки**: EdTech-платформа переопределяет `--edutrack-*` без копирования кода.
- **Изоляция**: изменение дизайна лендинга не затрагивает существующие компоненты проекта.
- **Developer Guidelines**: каждый разработчик знает, какой подход применить в конкретной ситуации.
- **Чёткий путь конвергенции**: триггеры и критерии, не форсировать без необходимости.

*Отрицательные и смягчение:*
- **Три уровня токенов** (`@theme` + `var(--...)` + `var(--edutrack-*)`) — сложнее, чем одна система. Смягчение: уровни строго разграничены по области применения (п.1), Developer Guidelines (п.8) дают однозначный ответ для каждой ситуации.
- **Нет автодополнения IDE для `var(--primary-warm)`** — в отличие от `bg-primary`. Смягчение: типичная плата за кастомные переменные; при миграции на `@theme` (Phase 2) автодополнение появится.
- **Дополнительный CSS-файл**: `pixso-variables.css` — ещё один запрос при загрузке. Смягчение: файл кэшируется браузером, размер < 1 KB.
- **Риск разрастания `pixso-variables.css`**: каждая новая страница из Pixso может добавлять свои переменные. Смягчение: конвенция нейминга (п.5) предписывает семантические имена, переиспользуемые между страницами; код-ревью контролирует отсутствие дублирования.
- **Ограничение `font` shorthand** — Tailwind не поддерживает `font` shorthand через arbitrary values. Смягчение: использовать отдельные Tailwind-утилиты (`text-2xl font-bold`) или `style`-проп для точного совпадения с дизайном.

**Связанные артефакты:**
- [Дизайн как код (Pixso)](ADR-DES.PROCESS.pixso-design-adoption.md) — процесс дизайна, токены, темы
- [Дизайн лендинга](ADR-DES.UI.landing-page-design.md) — структура посадочной страницы
- [Инструменты разработки](ADR-IMPL.PROCESS.development-tooling.md) — §11 дизайн-процесс Pixso
- [Стек: фреймворки (React + Tailwind)](ADR-DES.STACK.framework-vs-vs.md) — T3
- [GUI-интеграция в EdTech-платформы](ADR-DES.UI.eduplatform-gui-integration.md) — white-labeling, iframe-виджеты
- Скилл экспорта: `.agents/skills/pixso-landing-export/SKILL.md`
- Референс цветовой стратегии: `.agents/skills/pixso-landing-export/references/COLOR-STRATEGY.md`