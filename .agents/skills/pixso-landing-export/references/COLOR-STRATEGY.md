# Color Strategy: Pixso Variables as Source of Truth

## Decision

**Approach 1 (Pixso CSS variables) is the standard.** This is non-negotiable for
the current project stage. The agent MUST apply this approach without discussion.

## Rationale

The project's `@theme` tokens in `frontend/src/styles/index.css` are documented
as **placeholders**:

```css
/* VEDO EduTrack design tokens (placeholder — populated later via Pixso design process).
   See ADR-IMPL.PROCESS.development-tooling §11 (Pixso → production code). */
```

The actual design uses an orange/warm palette (`--primary-warm: rgba(234, 88, 12, 1)`),
while the project tokens are blue (`--color-primary-500: #3b82f6`). Overwriting
project tokens would break all existing components (Dashboard, auth, etc.).

## Implementation

### CSS Variables File

All Pixso colors are stored in `frontend/src/styles/pixso-variables.css`:

```css
/* ===== LIGHT THEME (default) ===== */
:root,
[data-collection-3-4-mode="light"] {
  --background: rgba(255, 249, 240, 1);
  --foreground: rgba(28, 25, 23, 1);
  --muted: rgba(253, 243, 231, 1);
  --muted-foreground: rgba(87, 83, 78, 1);
  --card: rgba(255, 255, 255, 1);
  --border: rgba(245, 230, 216, 1);
  --primary-warm: rgba(234, 88, 12, 1);
  --secondary-calm: rgba(20, 184, 166, 1);
  --accent-warm: rgba(245, 158, 11, 1);
  --accent-foreground: rgba(67, 20, 7, 1);
  --primary-foreground: rgba(255, 255, 255, 1);
}

/* ===== DARK THEME ===== */
[data-collection-3-4-mode="dark"] {
  --background: rgba(12, 10, 9, 1);
  --foreground: rgba(250, 250, 249, 1);
  --muted: rgba(28, 25, 23, 1);
  --muted-foreground: rgba(168, 162, 158, 1);
  --card: rgba(41, 37, 36, 1);
  --border: rgba(68, 64, 60, 1);
  --primary-warm: rgba(251, 146, 60, 1);
  --secondary-calm: rgba(45, 212, 191, 1);
  --accent-warm: rgba(251, 191, 36, 1);
  --accent-foreground: rgba(69, 26, 3, 1);
  --primary-foreground: rgba(67, 20, 7, 1);
}
```

### Usage in Tailwind

```tsx
<div className="bg-[var(--background)] text-[var(--foreground)]" />
<button className="bg-[var(--primary-warm)] text-[var(--accent-foreground)]" />
```

### Theme Switching

The `data-collection-3-4-mode` attribute on `<html>` switches all variables
between light and dark. The React component handles this automatically:

1. Checks `localStorage` for manual override
2. Falls back to `window.matchMedia('(prefers-color-scheme: dark)')`
3. Listens for live OS theme changes via `change` event

## Dark Theme Export

The dark theme is exported from a separate Pixso frame (e.g. "Дай пять / Desktop / Dark").
Only the CSS variables are extracted — the HTML structure is identical to the
light frame and is discarded.

### Export Steps

1. Export light frame → get structure + light variables
2. Export dark frame → extract only `:root` / `[data-collection-...]` CSS blocks
3. Merge dark variables into `pixso-variables.css` under `[data-collection-3-4-mode="dark"]`
4. Discard dark frame's HTML body (identical to light)

## Migration Path

When the project adopts the orange palette globally:

```
Phase 1 (now):     Pixso var(--primary-warm) in landing component
                   ↓
Phase 2 (later):   @theme { --color-primary: #ea580c }
                   + convert landing to use `bg-primary` classes
                   ↓
Phase 3 (future):  Migrate other components from blue → orange
                   + remove pixso-variables.css
```