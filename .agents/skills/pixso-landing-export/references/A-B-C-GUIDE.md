# A+B+C Refinement Guide for Pixso Export

This reference documents the exact conversion patterns used to transform
Pixso-generated HTML into a Tailwind CSS v4 React component.

## Dark Theme Export

Pixso provides separate frames for light and dark themes (e.g.
"Дай пять / Desktop / Light" and "Дай пять / Desktop / Dark"). The structure
is identical — only CSS variables differ.

### Export Process

1. Export **light** frame via `design_to_code` → provides HTML structure + light variables
2. Export **dark** frame via `design_to_code` → extract only the CSS variable block
3. Merge dark variables into `pixso-variables.css` under `[data-collection-3-4-mode="dark"]`
4. Discard dark frame's HTML body (identical structure)

### Theme Switching in the Component

The React component handles theme automatically:

```tsx
import { useEffect } from 'react';

const THEME_KEY = 'pixso-theme';

function getInitialTheme(): string {
  const stored = localStorage.getItem(THEME_KEY);
  if (stored === 'dark' || stored === 'light') return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

export function Landing() {
  useEffect(() => {
    // Apply initial theme
    const theme = getInitialTheme();
    document.documentElement.setAttribute('data-collection-3-4-mode', theme);

    // Listen for OS theme changes in real time
    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = (e: MediaQueryListEvent) => {
      const next = e.matches ? 'dark' : 'light';
      document.documentElement.setAttribute('data-collection-3-4-mode', next);
      localStorage.setItem(THEME_KEY, next);
    };
    mq.addEventListener('change', handler);
    return () => mq.removeEventListener('change', handler);
  }, []);

  return (
    <div className="w-full relative flex flex-col items-start bg-[var(--background)]">
      {/* sections */}
    </div>
  );
}
```

Key features:
- `localStorage` override priority (user can manually set theme)
- Falls back to `prefers-color-scheme` media query
- Live listener for OS theme changes
- Sets `data-collection-3-4-mode` on `<html>` (matching Pixso's dark mode selector)
- Cleanup on component unmount

## A — Responsive Layout

### Breakpoints

| Screen | Breakpoint | Changes |
|--------|-----------|---------|
| Desktop | ≥ 769px | Full multi-column layout |
| Tablet | 481–768px | Stack some grids, reduce padding |
| Mobile | ≤ 480px | Single column, minimal padding, hide decorations |

### Conversion Table

| Pixso Pattern | Tailwind Equivalent |
|--------------|-------------------|
| `width: 720px` | `max-w-[720px] w-full` |
| `min-width: 1920px` | Remove; use `w-full` |
| `height: 10602px` | Remove; use auto height |
| `display: flex; flex-direction: row` | `flex flex-row` |
| `flex-direction: column` | `flex-col` |
| `justify-content: space-between` | `justify-between` |
| `align-items: center` | `items-center` |
| `gap: 32px` | `gap-8` |
| `gap: 24px` | `gap-6` |
| `gap: 16px` | `gap-4` |
| `gap: 12px` | `gap-3` |
| `gap: 8px` | `gap-2` |
| `padding: 80px` | `p-20` (desktop) → `md:p-5` → `p-3` (mobile) |
| `padding: 40px` | `p-10` |
| `padding: 32px` | `p-8` |
| `padding: 24px` | `p-6` |
| `padding: 16px` | `p-4` |
| `padding: 140px 80px` | `px-5 md:px-20 py-[140px]` → mobile: `px-3 py-12` |

### Media Query Template

```css
@media (max-width: 768px) {
  .responsive-row { flex-direction: column !important; }
  .responsive-card { width: 100% !important; max-width: 100% !important; }
  .hide-mobile { display: none !important; }
  .section-padding { padding: 3rem 1rem !important; }
  .hero-text { font-size: clamp(2rem, 10vw, 3rem) !important; }
}

@media (max-width: 480px) {
  .section-padding { padding: 2rem 0.75rem !important; }
  .hero-text { font-size: clamp(1.5rem, 8vw, 2rem) !important; }
}
```

## B — CSS Variables

### Pixso Variable Mapping

Pixso exports color tokens as CSS custom properties automatically. These are
the standard variables produced by Pixso's design_to_code:

| Pixso Variable | Light Value | Dark Value | Usage |
|---------------|-------------|-----------|-------|
| `--background` | `rgba(255, 249, 240, 1)` | `rgba(12, 10, 9, 1)` | Page background |
| `--foreground` | `rgba(28, 25, 23, 1)` | `rgba(250, 250, 249, 1)` | Primary text |
| `--muted` | `rgba(253, 243, 231, 1)` | `rgba(28, 25, 23, 1)` | Muted background |
| `--muted-foreground` | `rgba(87, 83, 78, 1)` | `rgba(168, 162, 158, 1)` | Secondary text |
| `--card` | `rgba(255, 255, 255, 1)` | `rgba(41, 37, 36, 1)` | Card background |
| `--border` | `rgba(245, 230, 216, 1)` | `rgba(68, 64, 60, 1)` | Borders |
| `--primary-warm` | `rgba(234, 88, 12, 1)` | `rgba(251, 146, 60, 1)` | Primary accent |
| `--secondary-calm` | `rgba(20, 184, 166, 1)` | `rgba(45, 212, 191, 1)` | Secondary accent |
| `--accent-warm` | `rgba(245, 158, 11, 1)` | `rgba(251, 191, 36, 1)` | Accent |
| `--accent-foreground` | `rgba(67, 20, 7, 1)` | `rgba(69, 26, 3, 1)` | Text on accent |
| `--primary-foreground` | `rgba(255, 255, 255, 1)` | `rgba(67, 20, 7, 1)` | Text on primary |

### Usage in Tailwind

Use arbitrary value syntax to reference CSS variables.
The same `var(--*)` works for both light and dark themes because the
`data-collection-3-4-mode` attribute on `<html>` switches their values:

```tsx
{/* Background */}
<div className="bg-[var(--background)]">
<div className="bg-[var(--card)]">

{/* Text */}
<p className="text-[var(--foreground)]">
<p className="text-[var(--muted-foreground)]">

{/* Accent */}
<button className="bg-[var(--primary-warm)] text-[var(--accent-foreground)]">
<div className="bg-[var(--muted)]">

{/* Borders */}
<div className="border border-[var(--border)]">
```

## C — Tailwind Conversion

### Common Pixso → Tailwind Patterns

#### Layout Containers

```
Pixso: class="Pixso-frame-3_1309"
       style="width:100%; display:flex; flex-direction:column; background-color:var(--background)"
→ Tailwind: class="w-full flex flex-col bg-[var(--background)]"
```

#### Text Elements

```
Pixso: class="Pixso-paragraph-3_1317"
       style="font-size:22px; font-family:Inter-Bold; color:var(--foreground)"
→ Tailwind: class="text-[clamp(1rem,2.5vw,1.375rem)] font-['Inter'] font-bold text-[var(--foreground)]"
```

#### Cards with Borders

```
Pixso: class="stroke-wrapper-3_1352" + inner div + "stroke-3_1352"
→ Tailwind: class="border border-[var(--border)] rounded-[28px]"
```

#### Icons (SVG)

```
Pixso: class="Pixso-vector-3_1312"
       style="background-image:url(assets/hand.svg); width:28px; height:28px"
→ Tailwind: class="w-7 h-7 bg-[url('/assets/landing/hand.svg')] bg-contain bg-no-repeat bg-center"
```

#### Buttons / CTAs

```
Pixso: class="Pixso-instance-7_291"
       style="display:flex; padding:14px 24px; border-radius:999px; background-color:var(--primary-warm)"
→ Tailwind: class="flex items-center gap-3 px-6 py-[14px] rounded-[999px] bg-[var(--primary-warm)] cursor-pointer"
```

### Font Families

| Pixso Font Name | Tailwind Class |
|----------------|----------------|
| `Inter-Regular` (400) | `font-['Inter'] font-normal` |
| `Inter-Medium` (500) | `font-['Inter'] font-medium` |
| `Inter-Semi Bold` (600) | `font-['Inter'] font-semibold` |
| `Inter-Bold` (700) | `font-['Inter'] font-bold` |
| `Inter-Extra Bold` (800) | `font-['Inter'] font-extrabold` |

### Border Radius

| Pixso Value | Tailwind |
|-------------|----------|
| `border-radius: 999px` | `rounded-[999px]` (full pill) |
| `border-radius: 32px` | `rounded-[32px]` |
| `border-radius: 28px` | `rounded-[28px]` |
| `border-radius: 24px` | `rounded-[24px]` |
| `border-radius: 20px` | `rounded-[20px]` |
| `border-radius: 12px` | `rounded-[12px]` |
| `border-radius: 8px` | `rounded-lg` or `rounded-[8px]` |

### Shadows

| Pixso | Tailwind |
|-------|----------|
| `box-shadow: 0px 12px 48px -16px rgba(124,45,18,0.15)` | `shadow-[0_12px_48px_-16px_rgba(124,45,18,0.15)]` |
| `box-shadow: 0px 32px 4px -16px rgba(0,0,0,0.15)` | `shadow-[0_32px_4px_-16px_rgba(0,0,0,0.15)]` |

## File Structure After Export

```
frontend/
├── public/
│   └── assets/
│       └── landing/           # Downloaded SVGs and images
│           ├── hand.svg
│           ├── arrowright.svg
│           ├── sparkles.svg
│           ├── ...
│           └── fonts/
│               ├── Inter_1.ttf
│               ├── Inter_Medium.ttf
│               └── ...
├── src/
│   ├── pages/
│   │   └── Landing.tsx         # Generated React component
│   └── styles/
│       └── pixso-variables.css # Pixso color tokens
```
