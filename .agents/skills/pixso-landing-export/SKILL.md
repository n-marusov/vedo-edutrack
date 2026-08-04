---
name: pixso-landing-export
description: >-
  Export a Pixso landing page design to a production-ready React component with
  Tailwind CSS v4, responsive layout, Pixso CSS variables, and asset management.
  Use when the user says "export landing page", "export design to code",
  "convert to tailwind", "pixso to react", or after finishing a landing page
  design in Pixso. Always applies Approach 1 (Pixso CSS variables as-is).
argument-hint: "[guid] [--dark-guid <guid>] [--no-route]"
user-invocable: true
allowed-tools: Read Write Edit Fetch CreateDirectory PixsoMCP(design_to_code) PixsoMCP(refine_generated_code) PixsoMCP(fetch_context) PixsoMCP(pixso_mcp_get_screenshot) Terminal(pnpm *) Terminal(mkdir *) Terminal(cp *) Terminal(mv *)
metadata:
  author: VEDO EduTrack
  version: "2.0"
  category: design-to-code
---

# Pixso Landing Page Export

Export a Pixso landing page design to a production-ready React component in the
VEDO EduTrack project. The pipeline applies **A+B+C** refinement automatically:
responsive layout, CSS variables, and Tailwind CSS v4 utility classes.

**Color strategy is non-negotiable: Pixso CSS variables (`var(--primary-warm)`)**
are kept as-is. This is formally documented in
`specs/adr/ADR-IMPL.UI.pixso-variables-approach.md`. See also
`references/COLOR-STRATEGY.md` for detailed rationale.

## When to use

- After finishing a landing page design in Pixso
- When the user says "export design to code", "convert to tailwind", "pixso to react"
- When the user asks to turn a Pixso frame into a working React component

## Arguments

| Argument | Description |
|----------|-------------|
| `<guid>` | Pixso node GUID for the light theme frame (e.g. `3:1309`). Omit to use current selection. |
| `--dark-guid <guid>` | Pixso node GUID for the dark theme frame (e.g. `3:1677`). Exports only CSS variables for dark mode. |
| `--no-route` | Skip route registration |
| `--no-assets` | Skip asset download |

## Workflow

### Step 1: Identify the target node(s)

Two frames needed:

1. **Light frame** (required) — the main design, provides structure + light variables
2. **Dark frame** (optional, `--dark-guid`) — same layout, dark palette, provides dark variables

Accept GUIDs in order of priority:
- Pixso URL → extract `item-id`
- GUID string → use directly
- No argument → use current selection

### Step 2: Fetch context

```
fetch_context({ include_map: true })
```

Confirm the light frame is the right one. If it's a page rather than a frame,
warn the user.

### Step 3: Generate code from light frame

```
design_to_code({
  guids: ["<LIGHT_GUID>"],
  clientFrameworks: "react"
})
```

### Step 4: Apply A+B+C refinement

```
refine_generated_code({ refinementTags: ["A", "B", "C"] })
```

### Step 5: Extract dark theme variables (if --dark-guid provided)

```
design_to_code({
  guids: ["<DARK_GUID>"],
  clientFrameworks: "react"
})
```

From this output, extract ONLY the `:root` / `[data-collection-...]` CSS
variable blocks for dark mode. Discard the HTML body — it's identical to
the light frame. Merge the dark variables into `pixso-variables.css`.

### Step 6: Fetch and transform the generated code

1. `fetch()` the code URL from Step 3
2. Apply refinement guides from Step 4:
   - Convert all `Pixso-*` CSS classes to Tailwind utility classes
   - Replace fixed widths with responsive equivalents
   - Add media queries at `768px` and `480px`
   - Ensure all colors use `var(--*)` CSS custom properties
   - Remove all `<style>` blocks and inline `style=""`
   - Remove empty `<script>` blocks
   - Change `<html lang="zh-CN">` to `<html lang="ru">`
   - Change `<title>` to meaningful title

### Step 7: Download assets

For each unique asset URL at `http://localhost:PORT/assets/BATCH_TS/...`:

1. `fetch()` the URL to get binary content
2. Save to `frontend/public/assets/landing/<basename>`
3. Update reference in code from full URL to `/assets/landing/<basename>`
4. Font files to `frontend/public/assets/landing/fonts/`

### Step 8: Create the React component

**File:** `frontend/src/pages/Landing.tsx`

The component includes system theme auto-detection:

```tsx
import { useEffect } from 'react';
import './pixso-variables.css';

const THEME_KEY = 'pixso-theme';

function getInitialTheme(): string {
  const stored = localStorage.getItem(THEME_KEY);
  if (stored === 'dark' || stored === 'light') return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

export function Landing() {
  useEffect(() => {
    const theme = getInitialTheme();
    document.documentElement.setAttribute('data-collection-3-4-mode', theme);

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
      {/* exported sections */}
    </div>
  );
}
```

Key details:
- Reads `localStorage` for manual override before system preference
- Falls back to `prefers-color-scheme` media query
- Listens for live OS theme changes
- Sets `data-collection-3-4-mode` attribute (matching Pixso's dark mode selector)
- Cleans up listener on unmount

### Step 9: Save CSS variables

**File:** `frontend/src/styles/pixso-variables.css`

Contains ALL variables from BOTH light and dark frames:

```css
/* ===== LIGHT THEME (default) ===== */
:root,
[data-collection-3-1301-mode="default"],
[data-collection-3-4-mode="light"] {
  --secondary-calm: rgba(20, 184, 166, 1);
  --primary-warm: rgba(234, 88, 12, 1);
  --accent-warm: rgba(245, 158, 11, 1);
  --accent-foreground: rgba(67, 20, 7, 1);
  --primary-foreground: rgba(255, 255, 255, 1);
  --border: rgba(245, 230, 216, 1);
  --card: rgba(255, 255, 255, 1);
  --muted-foreground: rgba(87, 83, 78, 1);
  --muted: rgba(253, 243, 231, 1);
  --foreground: rgba(28, 25, 23, 1);
  --background: rgba(255, 249, 240, 1);
}

/* ===== DARK THEME ===== */
[data-collection-3-4-mode="dark"] {
  --secondary-calm: rgba(45, 212, 191, 1);
  --primary-warm: rgba(251, 146, 60, 1);
  --accent-warm: rgba(251, 191, 36, 1);
  --accent-foreground: rgba(69, 26, 3, 1);
  --primary-foreground: rgba(67, 20, 7, 1);
  --border: rgba(68, 64, 60, 1);
  --card: rgba(41, 37, 36, 1);
  --muted-foreground: rgba(168, 162, 158, 1);
  --muted: rgba(28, 25, 23, 1);
  --foreground: rgba(250, 250, 249, 1);
  --background: rgba(12, 10, 9, 1);
}
```

### Step 10: Register route (unless --no-route)

The route already points to `./pages/Landing`. No changes needed in `routes.tsx`
when overwriting the existing file. If saving to a different path, update the
import accordingly.

### Step 11: Take screenshot and verify

```
pixso_mcp_get_screenshot({ guid: "<LIGHT_GUID>" })
pixso_mcp_get_screenshot({ guid: "<DARK_GUID>" })
```

Compare screenshots with the generated component to verify:
- All sections present
- Layout preserved
- Colors match
- Dark mode switches correctly

## Example

```
/pixso-landing-export 3:1309 --dark-guid 3:1677
/pixso-landing-export 3:1309 --no-route
/pixso-landing-export --dark-guid 3:1677
```

## Supporting Files

- [references/COLOR-STRATEGY.md](references/COLOR-STRATEGY.md) — Pixso variables vs Tailwind @theme tokens comparison (Approach 1 is the default)
- [references/A-B-C-GUIDE.md](references/A-B-C-GUIDE.md) — Detailed conversion tables for Pixso to Tailwind, responsive breakpoints, CSS variable mapping, and asset handling
- [ADR-IMPL.UI.pixso-variables-approach.md](../../specs/adr/ADR-IMPL.UI.pixso-variables-approach.md) — Formal Architecture Decision Record for the CSS variables strategy

## Edge Cases

| Situation | Handling |
|-----------|----------|
| Node is a page, not a frame | Warn user, ask for specific frame |
| Server returns 404 for assets | Skip missing asset, log warning |
| Component already exists | Ask user: overwrite, merge, or skip? |
| No route file found | Skip route registration, log warning |
| Empty selection | Fetch context to show available frames |
| Font @font-face URLs are localhost | Download and save to `frontend/public/assets/landing/fonts/`, update src |
| SVG data:image URLs | Keep as embedded data URIs (no download needed) |
| Gradient backgrounds | Keep as `background-image` or inline SVG data URIs |
| Component has interactive state | Pixso export is static; add interactivity manually after export |
| --dark-guid points to same colors as light | Log warning: dark theme may be identical to light; verify |