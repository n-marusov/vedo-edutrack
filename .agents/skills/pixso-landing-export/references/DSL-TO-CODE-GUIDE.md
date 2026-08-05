# DSL → Tailwind Mapping Guide

This reference documents how the `dsl-to-code.py` script converts Pixso compact DSL properties to Tailwind CSS v4 classes. Use this when manually reviewing or fixing the generated code.

## AutoLayout → Flex/Grid

### Direction

| DSL value | Tailwind class |
|-----------|---------------|
| `autoLayout.direction: "VERTICAL"` | `flex-col` |
| `autoLayout.direction: "HORIZONTAL"` | `flex-row` |

### Gap

| DSL value | Tailwind class |
|-----------|---------------|
| `autoLayout.gap: 4` | `gap-1` |
| `autoLayout.gap: 8` | `gap-2` |
| `autoLayout.gap: 12` | `gap-3` |
| `autoLayout.gap: 16` | `gap-4` |
| `autoLayout.gap: 20` | `gap-5` |
| `autoLayout.gap: 24` | `gap-6` |
| `autoLayout.gap: 32` | `gap-8` |
| `autoLayout.gap: 40` | `gap-10` |
| `autoLayout.gap: 48` | `gap-12` |
| `autoLayout.gap: 64` | `gap-16` |
| `autoLayout.gap: 72` | `gap-18` |
| `autoLayout.gap: 80` | `gap-20` |
| Other | `gap-[<value>px]` |

### Alignment (Horizontal direction)

| DSL value | Tailwind class |
|-----------|---------------|
| `align: "space-between"` or `primaryAlign: "space-between"` | `justify-between` |
| `align: "center"` or `counterAlign: "center"` | `items-center` |
| `counterAlign: "flex-start"` | `items-start` |
| `counterAlign: "flex-end"` | `items-end` |

### Alignment (Vertical direction)

| DSL value | Tailwind class |
|-----------|---------------|
| `align: "center"` or `primaryAlign: "center"` | `justify-center` |
| `counterAlign: "center"` | `items-center` |
| `counterAlign: "flex-start"` | `items-start` |

### Sizing

| DSL value | Tailwind class |
|-----------|---------------|
| `widthResize: 2` (FILL) | `w-full` |
| `widthResize: 1` (HUG) | `w-fit` |
| `heightResize: 2` (FILL) | `h-full` |

### Padding

| DSL value | Tailwind class |
|-----------|---------------|
| `padding: 40` | `p-10` |
| `padding: [20, 32, 20, 32]` | `py-5 px-8` |
| `padding: [140, 80, 140, 80]` | `py-[140px] px-20` |

Common pixel-to-Tailwind padding mapping:
| px | Tailwind |
|----|----------|
| 8 | `p-2` |
| 12 | `p-3` |
| 16 | `p-4` |
| 20 | `p-5` |
| 24 | `p-6` |
| 32 | `p-8` |
| 40 | `p-10` |
| 48 | `p-12` |
| 64 | `p-16` |
| 80 | `p-20` |
| 100 | `p-25` |
| 140 | `p-[140px]` |

## Fill → CSS Variable

### Variable References

DSL fills contain `variable:<id>` references. These resolve to:

| Variable ID | Name | Tailwind usage |
|-------------|------|---------------|
| `3:11` | `background` | `bg-[var(--background)]` |
| `3:12` | `foreground` | `text-[var(--foreground)]` |
| `3:13` | `muted` | `bg-[var(--muted)]` |
| `3:14` | `muted-foreground` | `text-[var(--muted-foreground)]` |
| `3:15` | `card` | `bg-[var(--card)]` |
| `3:17` | `border` | `border-[var(--border)]` |
| `3:1303` | `primary-warm` | `bg-[var(--primary-warm)]` |
| `3:1305` | `secondary-calm` | `text-[var(--secondary-calm)]` |
| `3:1302` | `accent-warm` | `bg-[var(--accent-warm)]` |
| `3:22` | `accent-foreground` | `text-[var(--accent-foreground)]` |
| `3:20` | `primary-foreground` | `text-[var(--primary-foreground)]` |

### Gradient Fills

For gradient fills, the script generates CSS gradient strings:
- `GRADIENT_LINEAR` → `linear-gradient(to right, rgba(...) 0%, rgba(...) 100%)`
- `GRADIENT_RADIAL` → `radial-gradient(circle, rgba(...) 0%, rgba(...) 100%)`

These are used as `bg-[<gradient>]` in Tailwind.

## Corner Radius

| DSL value | Tailwind class |
|-----------|---------------|
| `cornerRadius: 4` | `rounded-sm` |
| `cornerRadius: 8` | `rounded` |
| `cornerRadius: 12` | `rounded-xl` |
| `cornerRadius: 16` | `rounded-2xl` |
| `cornerRadius: 20` | `rounded-[20px]` |
| `cornerRadius: 24` | `rounded-[24px]` |
| `cornerRadius: 28` | `rounded-[28px]` |
| `cornerRadius: 32` | `rounded-[32px]` |
| `cornerRadius: 40` | `rounded-[40px]` |
| `cornerRadius: 999` | `rounded-full` |

## Effects → Shadow

Drop shadows in the DSL have this structure:
```json
{
  "type": "DROP_SHADOW",
  "color": {"r": 124, "g": 45, "b": 18, "a": 0.14902},
  "offset": {"x": 0, "y": 16},
  "radius": 64,
  "spread": -20
}
```

Generated Tailwind: `shadow-[0_16px_64px_-20px_rgba(124,45,18,0.149)]`

Common shadow patterns in the design:

| Description | Tailwind |
|-------------|----------|
| Card shadow | `shadow-[0_12px_48px_-16px_rgba(124,45,18,0.149)]` |
| Large card shadow | `shadow-[0_16px_64px_-20px_rgba(124,45,18,0.149)]` |
| Hero card shadow | `shadow-[0_24px_64px_-20px_rgba(124,45,18,0.149)]` |

## Text Properties

### Font Size

| DSL fontSize | Tailwind class |
|-------------|---------------|
| 11 | `text-[11px]` |
| 12 | `text-xs` |
| 13 | `text-[13px]` |
| 14 | `text-sm` |
| 15 | `text-[15px]` |
| 16 | `text-base` |
| 17 | `text-[17px]` |
| 18 | `text-lg` |
| 20 | `text-xl` |
| 22 | `text-[22px]` |
| 24 | `text-2xl` |
| 28 | `text-3xl` |
| 32 | `text-4xl` |
| 48 | `text-5xl` |
| 64 | `text-7xl` |
| 72 | `text-8xl` |

### Font Weight

| DSL fontWeight | fontStyle | Tailwind class |
|--------------|-----------|---------------|
| 400 | "Regular" | `font-normal` |
| 500 | "Medium" | `font-medium` |
| 600 | "Semi Bold" | `font-semibold` |
| 700 | "Bold" | `font-bold` |
| 800 | "Extra Bold" | `font-extrabold` |

### Line Height

| DSL lineHeightUnit + lineHeightNumber | Tailwind class |
|---------------------------------------|---------------|
| PERCENT 100 | `leading-none` |
| PERCENT 110 | `leading-none` |
| PERCENT 130 | `leading-tight` |
| PERCENT 140 | `leading-snug` |
| PERCENT 150 | `leading-normal` |
| PERCENT 160 | `leading-relaxed` |
| PERCENT 170 | `leading-loose` |
| PIXELS <value> | `leading-[<value>px]` |

## Icon Mapping

The script maps Pixso vector names to lucide-react icons using the `ICON_MAP`
dictionary. Key mappings:

| Pixso name | lucide-react icon |
|------------|------------------|
| `hand` | `Hand` |
| `sparkles` | `Sparkles` |
| `arrow-right` | `ArrowRight` |
| `rocket` | `Rocket` |
| `target` | `Target` |
| `clock` | `Clock` |
| `search` | `Search` |
| `unlink` | `Unlink` |
| `question` | `HelpCircle` |
| `alert-circle` | `AlertCircle` |
| `check` | `Check` |
| `plus` | `Plus` |
| `minus` | `Minus` |
| `quote` | `Quote` |
| `compass` | `Compass` |
| `layers` | `Layers` |
| `link` | `Link2` |
| `users` | `Users` |
| `eye` | `Eye` |
| `badge-check` | `BadgeCheck` |
| `map` | `Map` |
| `navigation` | `Navigation` |
| `info` | `Info` |
| `briefcase` | `Briefcase` |
| `external-link` | `ExternalLink` |
| `send` | `Send` |
| `play` | `Play` |
| `flask-conical` | `FlaskConical` |
| `folder-kanban` | `FolderKanban` |
| `refresh-cw` | `RefreshCw` |

## Data Extraction Patterns

The `--data-only` flag extracts these data arrays from the DSL:

### Pricing Plans
```json
{
  "plans": [
    {
      "name": "Карта",
      "price": "0 ₽",
      "period": "/ 7 дней",
      "description": "Диагностика пробелов и карта тем на 7 дней.",
      "features": [
        {"text": "Карта всех тем по предметам", "highlight": false},
        {"text": "Истории и контекст «зачем это знать»", "highlight": true}
      ],
      "annual_price": "4 410 ₽",
      "annual_old": "5 880 ₽",
      "annual_savings": "Экономия 1 470 ₽",
      "cta": "Попробовать 7 дней"
    }
  ]
}
```

### Testimonials
```json
{
  "testimonials": [
    {
      "text": "«До «Дай пять» я тратила 3 часа в неделю...»",
      "author": "Елена, 38 лет",
      "role": "Мама Маши (6 класс) и Пети (3 класс)",
      "result": "Результат: аттестация сдана на «отлично», пробелов нет"
    }
  ]
}
```

### FAQ
```json
{
  "faqs": [
    {
      "question": "Что, если мы не успеем подготовиться к аттестации?",
      "answer": "Система заранее показывает пробелы..."
    }
  ]
}
```

## Verifying with SVG

When the Python output looks off on positioning, export the section as SVG:

```
get_export_image({ guid: "<SECTION_GUID>", exportSettings: { imageType: 3, constraint: { type: 1, value: 1 } } })
```

Download with curl and inspect key attributes:

```bash
curl -s "<URL>" -o section.svg
# Check radii:
grep -oP 'rx="[^"]*"' section.svg
# Check padding/gaps from rect positions:
grep -oP 'x="[^"]*" y="[^"]*"' section.svg | head -20
# Check colors:
grep -oP 'fill="[^"]*"' section.svg | sort -u
```

Key SVG elements to look for:
- `<rect>` with `rx="..."` — corner radius of cards
- `<rect>` with `x="..." y="..."` — absolute position of elements
- `fill="rgb(...)"` — exact colors used
- `<g>` grouping — indicates nested containers
