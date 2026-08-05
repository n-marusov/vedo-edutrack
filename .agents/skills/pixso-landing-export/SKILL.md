---
name: pixso-landing-export
description: >-
  Export a Pixso landing page design to a production-ready React component with
  Tailwind CSS v4, responsive layout, and Pixso CSS variables. Uses DSL+SVG
  extraction (NOT design_to_code) for pixel-accurate code generation with
  section-by-section automated processing via Python script. Use when the user
  says "export landing page", "export design to code", "convert to tailwind",
  "pixso to react", or after finishing a landing page design in Pixso.
argument-hint: "[guid] [--dark-guid <guid>] [--section <name>] [--no-route] [--no-assets]"
user-invocable: true
allowed-tools: Read Write Edit Fetch CreateDirectory PixsoMCP(get_node_dsl) PixsoMCP(get_export_image) PixsoMCP(fetch_context) PixsoMCP(pixso_mcp_get_screenshot) PixsoMCP(get_variables) PixsoMCP(get_top_level_frames) PixsoMCP(query_nodes) Terminal(pnpm *) Terminal(mkdir *) Terminal(cp *) Terminal(mv *) Terminal(python *) Terminal(curl *) SpawnAgent
metadata:
  author: VEDO EduTrack
  version: "3.0"
  category: design-to-code
---

# Pixso Landing Page Export (DSL+SVG Method)

**DO NOT use `design_to_code`.** It produces unreliable output: Pixso-* classes,
slot_4_258 props, separate 182KB CSS files, invented filler content, and wrong
text/data. Instead, use **DSL extraction + SVG verification + Python code generation**.

**Color strategy is non-negotiable: Pixso CSS variables (`var(--primary-warm)`)**
are kept as-is. See `specs/adr/ADR-IMPL.UI.pixso-variables-approach.md`.

## Arguments

| Argument | Description |
|----------|-------------|
| `<guid>` | Pixso node GUID for the light theme frame (e.g. `3:1309`). Omit to use current selection. |
| `--dark-guid <guid>` | Pixso node GUID for the dark theme frame. Exports only CSS variables. |
| `--section <name>` | Process only one section (e.g. `Pricing`, `Benefits`). |
| `--no-route` | Skip route registration. |
| `--no-assets` | Skip asset download. |

## Workflow

### Step 1: Identify the target frame

```
fetch_context({ include_map: true })
```

Confirm the frame is correct. Collect the light frame GUID (required) and
dark frame GUID (optional, for theme variables).

### Step 2: Extract variables (design tokens)

```
get_variables({ variableSetId: "<id>" })
```

This gives the color tokens for light/dark modes. Resolve the Pixso variable
set IDs from the frame context. The variable set typically has modes:
- `light` — light theme colors
- `dark` — dark theme colors

### Step 3: Get DSL per section — NOT design_to_code

For EACH major section of the page, get the structured DSL:

```
get_node_dsl({ guid: "<SECTION_GUID>", simplify: true })
```

The compact DSL (simplify:true) returns:
- `roots` — tree of nodes with exact text content, sizes, autoLayout, fills
- `refsIndex` — resolved variable names, component references, icon IDs
- `variableMap` — Pixso variable ID → CSS variable name mapping

**Save each section's DSL to a temp JSON file** for the Python script.

### Step 4: Verify geometry with SVG (optional, for critical sections)

For sections where exact positioning matters (Pricing cards, grids, etc.):

```
get_export_image({ guid: "<SECTION_GUID>", exportSettings: { imageType: 3, constraint: { type: 1, value: 1 } } })
```

Download the SVG to inspect:
- **Radii** — check `rx="..."` on `<rect>` elements for exact corner radius
- **Positions** — check `x="..." y="..."` for element placement
- **Colors** — check `fill="rgb(...)"` for exact color values
- **Layout** — verify gaps, padding, alignment

SVG is large (200-300KB per card) but contains pixel-exact geometry. Use it
for verification, not for text extraction (text renders as glyph paths).

### Step 5: Run the DSL-to-Code Python script

Save the DSL JSON to a temp file, then run:

```bash
python .agents/skills/pixso-landing-export/scripts/dsl-to-code.py <dsl.json> --section-name "<Name>"
```

The script outputs:
- **JSX structure** — recursive frame→div, text→span, autoLayout→flex/grid mapping
- **Data arrays** — extracted structured data (pricing plans, testimonials, FAQs, etc.)
- **Icon references** — Pixso vector → lucide-react mapping
- **Variable resolution** — fills → CSS var() references

For structured data-heavy sections, also run with `--data-only`:

```bash
python .agents/skills/pixso-landing-export/scripts/dsl-to-code.py <dsl.json> --data-only
```

This outputs JSON with arrays like `plans`, `testimonials`, `faqs`, `principles`,
`problems`, `benefits`, `metrics` — ready to paste into the React component.

### Step 6: Generate code section by section

For EACH section, generate a focused React component:

1. Save the section's DSL JSON from Step 3
2. Run the Python script to extract data + structure
3. Feed the output to a **spawned code-gen agent** with instructions:
   - Write clean React+Tailwind code
   - Use `var(--*)` for all colors
   - Use lucide-react icons
   - Match the DSL structure exactly
   - Verify radius/position from SVG if exported
4. Integrate the section into the master `Landing.tsx`

**NEVER generate the whole page at once** — section-by-section is more accurate
because the agent has enough context for each section.

### Step 7: Handle known component types

The script detects these component instances and generates proper code:

| Pixso Component | React Component |
|----------------|-----------------|
| `Component / Section Header` | `<SectionHeader badge="..." title="..." subtitle="..." />` |
| `Component / CTA Button` | `<OrangeButton>...</OrangeButton>` |
| `Component / Eyebrow` | `<SectionBadge>...</SectionBadge>` |

For the project, these reusable components are defined in `Landing.tsx`.

### Step 8: Save CSS variables

**File:** `frontend/src/styles/pixso-variables.css`

Contains ALL variables from BOTH light and dark frames. The dark theme values
are extracted from the dark frame's DSL (same structure, different fill values).

```css
:root,
[data-collection-3-4-mode="light"] {
  --background: rgba(255, 249, 240, 1);
  --foreground: rgba(28, 25, 23, 1);
  ...
}

[data-collection-3-4-mode="dark"] {
  --background: rgba(12, 10, 9, 1);
  --foreground: rgba(250, 250, 249, 1);
  ...
}
```

### Step 9: Write the theme hook

```tsx
import { useEffect, useState } from 'react';

const THEME_KEY = 'pixso-theme';

function useTheme() {
  const [theme, setTheme] = useState<string>(() => {
    const stored = localStorage.getItem(THEME_KEY);
    if (stored === 'dark' || stored === 'light') return stored;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  });

  useEffect(() => {
    document.documentElement.setAttribute('data-collection-3-4-mode', theme);
    localStorage.setItem(THEME_KEY, theme);
  }, [theme]);

  useEffect(() => {
    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = (e: MediaQueryListEvent) => {
      setTheme(e.matches ? 'dark' : 'light');
    };
    mq.addEventListener('change', handler);
    return () => mq.removeEventListener('change', handler);
  }, []);

  return { theme, setTheme };
}
```

### Step 10: Register route (unless --no-route)

The route already points to `./pages/Landing`. No changes needed when
overwriting the existing file.

### Step 11: Take screenshot and verify

```
pixso_mcp_get_screenshot({ guid: "<LIGHT_GUID>" })
```

Compare with the rendered component in the browser:
- All sections present
- Layout preserved
- Colors match
- Responsive at 768px and 480px
- Dark mode switches correctly

## Section types and their data extraction patterns

### Pricing Section
The script extracts `plans` array with fields: name, price, period, badge,
features (with highlight flag), annual_price, annual_old, annual_savings, cta.

### Testimonials Section
The script extracts `testimonials` array with fields: text, author, role, result.

### FAQ Section
The script extracts `faqs` array with fields: question, answer.

### Principles/Philosophy Section
The script extracts `principles` array with fields: title, subtitle, desc, link.

### Benefits Section
The script extracts `benefits` array with fields: title, desc, instead.

### Metrics Section
The script extracts `metrics` array with fields: value, label.

## Example

```
# Full export (all sections)
/pixso-landing-export 3:1309 --dark-guid 3:1677

# Single section
/pixso-landing-export 3:1309 --section Pricing

# Export without route registration
/pixso-landing-export 3:1309 --no-route
```

## Python Script Reference

The `dsl-to-code.py` script supports these modes:

```bash
# Full JSX generation
python dsl-to-code.py section.json -n "Pricing" -o pricing.tsx

# Data-only extraction (JSON arrays)
python dsl-to-code.py section.json --data-only

# List icons found in the design
python dsl-to-code.py section.json --list-icons

# Structure-only (no data extraction)
python dsl-to-code.py section.json --structure
```

## Supporting Files

- [scripts/dsl-to-code.py](scripts/dsl-to-code.py) — Python DSL→code generator
- [references/DSL-TO-CODE-GUIDE.md](references/DSL-TO-CODE-GUIDE.md) — Detailed DSL→Tailwind mapping
- [references/COLOR-STRATEGY.md](references/COLOR-STRATEGY.md) — Pixso variables strategy (Approach 1)
- [ADR-IMPL.UI.pixso-variables-approach.md](../../specs/adr/ADR-IMPL.UI.pixso-variables-approach.md) — ADR for CSS variables

## Edge Cases

| Situation | Handling |
|-----------|----------|
| Node is a page, not a frame | Warn user, ask for specific frame |
| DSL has no children | Frame is empty — generate empty div |
| Component instance not recognized | Generate as generic div with children |
| Icon not in lucide-react | Use inline SVG from Pixso export |
| SVG export returns 404 | Skip SVG verification, use DSL only |
| Python script not found | Generate code manually from DSL data |
| Route file not found | Skip route registration, log warning |
| Empty selection | Fetch context to show available frames |
| Section has no data arrays | Generate JSX structure only |
| Dark theme identical to light | Log warning; verify manually |