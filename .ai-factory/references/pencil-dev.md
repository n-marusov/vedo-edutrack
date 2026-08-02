# Pencil (pen.dev) Reference

> Source:
> - https://pencil.dev/ (official site)
> - https://docs.pen.dev/ (official docs; retrieved via web.archive.org snapshots, 2026-07-20)
>   - .pen Files, Design as Code, Design ↔ Code, AI Integration, Variables, Components, Slots, Design Libraries, Import and Export, Keyboard Shortcuts, Installation, The .pen Format
> - https://betterstack.com/community/guides/ai/pencil-ai/ (guide, 2026-02-22)
> - https://raw.githubusercontent.com/unliftedq/skills/main/skills/pencil-dev/SKILL.md (community Claude Code skill)
> Created: 2026-08-02
> Updated: 2026-08-02

## Overview

Pencil (pen.dev) is an **agent-driven vector design tool that lives inside your IDE** (VS Code, Cursor, or standalone desktop app for macOS/Windows/Linux). Its tagline is "Design on canvas. Land in code." It is positioned as an alternative to the classic Figma→dev handoff: designs are stored as **`.pen` files — pure JSON** that live in your Git repository next to the code, so design and code are version-controlled, diffable, and readable by both humans and AI agents.

The key differentiator is a **local MCP (Model Context Protocol) server**: when pen.dev is running, AI assistants (Claude Code, Claude Desktop, Cursor, Windsurf, Codex CLI, Antigravity, OpenCode CLI) get read **and write** tools to manipulate the design canvas programmatically. The AI can create frames, insert components, apply styles, and arrange layouts; the user stays in control (AI suggests, user approves). This collapses the design→code handoff into a single version-controlled workflow.

Use cases for frontend design: rapid prototyping from text prompts, design-system construction (components, variables, light/dark themes), importing existing code as designs, and generating production code (React, Next.js, Vue, Svelte, Tailwind, plain HTML/CSS) from designs with high fidelity.

## Core Concepts

- **`.pen` file**: the design file format. JSON-based, version-control friendly (Git diffs on text), portable. Keep `.pen` files in the project workspace alongside code. There is **no auto-save yet** — save frequently (`Cmd/Ctrl + S`) and commit to Git. Example names: `dashboard.pen`, `components.pen`.
- **Object tree**: a `.pen` document is a JSON structure describing an *object tree*, similar to HTML or SVG. Each object is a graphical entity on an infinite 2D canvas with a unique `id` and a `type` (e.g. `rectangle`, `frame`, `text`, `ellipse`, `path`, `icon`, `script`, `ref`, etc.). Top-level objects have `x`/`y`; nested objects are positioned relative to their parent.
- **Flexbox-style layout**: parents can control children sizing/positioning via `layout` (`none`/`vertical`/`horizontal`), `gap`, `justifyContent`, `alignItems`, `padding`. Children can `fill` their parent or use fixed `width`/`height`; parents can fit content or be fixed size. Layout containers should carry structure — avoid overusing absolute positioning.
- **Components & instances**: an object with `"reusable": true` becomes a reusable component; `type: "ref"` creates an instance of it. Instances can **override** properties (e.g. `fill`) and customize descendants via a `descendants` map (keyed by descendant id-path, e.g. `"ok-button/label"`). Descendants can be property-overridden or fully replaced (presence of `type` = replacement). Children of a nested instance can be replaced — ideal for container components (panels, cards, sidebars).
- **Slots**: a frame inside a component marked with `slot` is a designated drop area. Only empty frames in component origins can become slots. Components in the same file can be marked as "suggested slot components" to guide what belongs in a slot.
- **Variables (design tokens)**: work like CSS custom properties. Define HEX colors, numbers (spacing, radii, sizes), strings (font names) once; changing a variable updates everywhere it is used. Variables can be created manually, from CSS (`globals.css`) via AI, or pasted from Figma. **Theming**: a variable can hold multiple values keyed by theme (e.g. `mode: light/dark`, `spacing: regular/condensed`); the last value whose theme is satisfied wins. Objects can set a `theme` that propagates to all descendants.
- **Design libraries**: a `.pen` file can be turned into a library (`.lib.pen` suffix — cannot be undone). Library components can be imported into other `.pen` files; edits to the library propagate to all usages. Default libraries (e.g. Shadcn UI, Lunaris, Halo, Nitro UI kits) ship with the app.
- **Design as Code**: commit `.pen` files like code, view text-based diffs in Git, branch and merge designs with code. Best practice: frequent commits with descriptive messages.
- **MCP server**: runs locally and automatically when pen.dev is open. No cloud dependency for design operations; design files stay local. Tool list is inspectable in IDE settings / `pencil` server in Codex `/mcp`.

## API / Interface

### MCP Tools (as exposed to AI assistants)

| Tool | Purpose |
|------|---------|
| `get_editor_state` | Current editor context: active file, selection, schema |
| `open_document` | Open an existing `.pen` file or create a new one |
| `get_guidelines` | Task-specific design guidance |
| `get_style_guide_tags` / `get_style_guide` | Discover and load style guides by tags/name |
| `get_variables` / `set_variables` | Read/write design tokens and theme values |
| `batch_get` | Read components, node trees, structured design data; search by patterns |
| `batch_design` | Create, modify, manipulate elements (insert, copy, update, replace, move, delete); set variables/themes; generate and place images |
| `find_empty_space_on_canvas` | Locate available canvas space for new content |
| `snapshot_layout` | Structural layout QA: detect positioning issues, overlapping elements |
| `get_screenshot` | Render design preview; visual QA only — never the primary source of structure |
| `search_all_unique_properties` / `replace_all_matching_properties` | Search and batch-replace property values across node trees |
| `export_nodes` | Export nodes to images (PNG/JPEG/WEBP/PDF) |

### .pen document format (key elements)

- `Document`: `version` (`"2.14"` at time of writing), optional `themes` (map of theme axis → allowed values, first value = default), `imports` (alias → relative `.pen` URI), `variables` (map of `$name` → value or theme-keyed value list), `children` (array of node objects).
- `Entity` (common): `id` (unique, must NOT contain `/`; auto-generated if omitted), `name`, `context`, `reusable`, `theme`, `enabled`, `opacity`, `flipX/flipY`, `rotation`, `layoutPosition` (`auto`/`absolute`), `metadata`.
- Node types: `frame`, `group`, `rectangle`, `ellipse`, `path`, `polygon`, `text`, `note`, `prompt`, `context`, `icon`, `script`, `ref`.
- `Variable` = dollar-prefixed string binding, e.g. `"fill": "$color.background"`. `NumberOrVariable`, `ColorOrVariable`, `BooleanOrVariable`, `StringOrVariable` union types.
- Graphics: `fill` (color / gradient / image / shader / mesh_gradient; multiple fills painted in order), `stroke` (single stroke, multiple fills), `effect` (blur, background_blur, shadow). Colors: `#RGB`, `#RRGGBB`, `#RRGGBBAA`.
- `Layout`: `layout: "none" | "vertical" | "horizontal"`; `gap`; `padding` (all / [v,h] / [t,r,b,l]); `justifyContent: start|center|end|space_between|space_around`; `alignItems: start|center|end`.
- `SizingBehavior`: `fit_content` / `fill_container` with optional fallback, e.g. `fit_content(100)`.
- `Text`: `content`, font props, `textGrowth: "auto" | "fixed-width" | "fixed-width-height"` (required before width/height take effect), `textAlign`, `lineHeight`.
- `Icon`: `library` in `'lucide' | 'feather' | 'Material Symbols Outlined' | 'Material Symbols Rounded' | 'Material Symbols Sharp' | 'phosphor'`; `icon`, `weight`.
- `Script`: generates nested children from JavaScript; `scriptUri` (relative JS file), `inputs`.
- `Ref`: `ref` (target id), `descendants` (id-path → overrides or replacements).

> Full authoritative TypeScript schema: https://docs.pen.dev/for-developers/the-pen-format. The format is live — breaking changes are possible.

## Usage Patterns

### 1. Prompt → design (AI generation)

1. Open/create a `.pen` file in pen.dev.
2. Create a frame (tool `A`/`F`), select it in the Layers panel and rename it (e.g. "Step 3 Frame") — prompts can target frames by name.
3. Prompt the AI (in terminal with `claude` or in the built-in chat `Cmd/Ctrl + K`), e.g.:
   ```
   Design a dashboard in the "Step 3 Frame" for a rover management platform using the components.
   - add a sidebar for navigation
   - add rover stats and table with available rovers for rent in the content body
   use the pencil mcp server
   ```
4. Watch the AI build the design step by step, then refine manually in the Properties inspector (padding, spacing, etc.).

Effective prompts are specific: state where to build (frame name), what to build, which components/kits to use, and layout requirements.

### 2. Design → Code

- Save the `.pen` file in the project workspace.
- Open AI chat (`Cmd/Ctrl + K`) and ask for code, e.g.:
  ```
  Create a React component for this button
  Generate a Next.js page from this design
  Export this dashboard as a React component with Tailwind CSS
  Generate code using Shadcn UI components
  Generate this design using Lucide icons
  ```
- The AI parses the `.pen` JSON structure (layout, styles, text, images) and writes matching code.

### 3. Code → Design (import existing UI)

- Keep the `.pen` file in the same workspace as the code so the agent can read both.
- Ask:
  ```
  Recreate the Button component from src/components/Button.tsx
  Import the LoginForm from my codebase into this design
  ```
- Imports component structure/hierarchy, layout/positioning, and styling (colors, typography, spacing).

### 4. Two-way sync loop

Start with code → import into pen.dev → design improvements → sync changes back to code → iterate. Design-system maintenance: define variables in pen.dev, sync to CSS, use variables in both design and code, update once.

### 5. Design tokens & CSS sync

- CSS → pen.dev: `Create pen.dev variables from my globals.css` / `Import design tokens from src/styles/tokens.css`.
- pen.dev → CSS: `Update globals.css with these pen.dev variables` / `Sync these design tokens to my CSS`.

### 6. Design system workflows

- `Create a button component with variants` / `Build a typography scale` / `Generate a color palette based on #3b82f6`.
- Consistency: `Ensure all buttons use the primary color variable` / `Apply 8px spacing grid to all elements`.
- Batch: `Create 5 variations of this button component` / `Design an entire landing page with hero, features, pricing, and footer`.

## Configuration

### Installation

| Platform | How |
|----------|-----|
| VS Code | Extensions (`Cmd/Ctrl + Shift + X`) → search "pen.dev" → Install. Verify: open a `.pen` file, look for the pen.dev icon top-right |
| Cursor | Extensions → search "pen.dev" → Install. May require Cursor Pro for some features |
| Desktop macOS | Download `.dmg` from pencil.dev → drag to Applications |
| Desktop Windows | Runs natively; or use the VS Code/Cursor extension |
| Desktop Linux | `.deb` (`sudo dpkg -i pencil-*.deb`) or `.AppImage` (`chmod +x pencil-*.AppImage; ./pencil-*.AppImage`); X11 more stable than Wayland/Hyprland |

### Prerequisites for AI features

- **Claude Code CLI**: `npm install -g @anthropic-ai/claude-code-cli` or `curl https://claude.ai/cli/install.sh | sh`; then `claude` to authenticate; verify with `claude --version`.
- **Activation**: complete pen.dev activation with your email after install.
- **MCP server**: starts automatically with pen.dev, runs locally. Verify: Cursor → Settings → Tools & MCP (look for `pencil`); Codex → `/mcp`.

### Supported AI assistants

Claude Code (CLI and IDE), Claude Desktop, Cursor, Windsurf IDE (Codeium), Codex CLI (OpenAI), Antigravity IDE, OpenCode CLI.

### Built-in icon libraries

Material Symbols (Outlined, Rounded, Sharp), Lucide Icons, Feather, Phosphor. Custom SVG icons can be imported like images. For generated code, specify library in the prompt (Lucide, Heroicons, FontAwesome, React Icons).

### Import formats

Complete Figma files (toolbar chevron → `Import Figma`; desktop `File > Import Image/SVG/Figma...`), individual Figma layers (copy/paste; images not supported via paste), images (drag-drop / clipboard / toolbar; PNG, JPEG, SVG).

### Export

Individual elements → PNG, JPEG, WEBP, PDF (select element → Properties panel bottom → choose size/format → Export layer). Code export via AI chat.

## Best Practices

1. **Keep `.pen` files in the repo next to code** — the AI agent can then see both design and code, version control tracks both, and sync stays easy.
2. **Save frequently and commit** — there is no auto-save yet; use Git for version history (`.pen` is text, diffs are readable).
3. **Start with variables, not literals** — typography, color, spacing, radius, shadow as a system. Reuse existing variables first, extend second, fall back to raw values only when necessary.
4. **Design light and dark together** — put theme differences at the variable layer early so components and pages inherit them.
5. **Extract reusable components before composing pages** — promote repeated structures (cards, nav, form rows, tables) into components; prefer instance overrides over detached copies.
6. **Keep shared shells consistent** — navigation, sidebar, top bar, filter bars, modal framing should not drift page by page.
7. **Let layout containers carry structure** — use `layout`/`gap`/`padding`/alignment instead of manual absolute positioning where possible.
8. **Prompt precisely** — "Increase the button padding to 16px and change color to blue" beats "Make it better"; provide context ("Add a login form with email, password, remember me checkbox, submit button").
9. **Iterate broad → refined** — start with a layout, then add components, then style states, then polish spacing (e.g. 8px grid).
10. **Validate structurally, then visually** — for implementation, read structured nodes/instances/variables as source of truth; use screenshots only for visual review, never to recover hierarchy/spacing/tokens.
11. **Use structured data, not screenshots, for code generation** — inspect `.pen` nodes, instances, and layout data; treat screenshots as visual references only.

## Common Pitfalls

- **Hard-coding font sizes, colors, spacing, or radius across page nodes** instead of using variables — breaks theming and makes light/dark and rebranding expensive.
- **Deferring dark mode** ("finish light first, promise dark later") — put both themes at the variable layer from the start.
- **Copy-pasting repeated UI fragments** instead of extracting components — duplicates drift apart.
- **Reading screenshots to infer node hierarchy, spacing, or token usage** when structured `.pen` data is available — recover structure from the JSON, not pixels.
- **Breaking shared components into detached copies** for page-specific tweaks — use instance overrides/descendants.
- **Shared page shells (nav, sidebar, top bar) diverging across screens** — abstract into reusable components/layout shells.
- **Mixing token-driven styling with arbitrary literals** without a clear reason.
- **Forgetting to save** — no auto-save; lost work is unrecoverable.
- **MCP connection issues** — "Claude Code not connected": ensure `claude` is logged in, restart pen.dev, run `claude` in the project directory. "MCP server not appearing": verify pen.dev is running, check IDE MCP settings, restart both. "Invalid API key": re-authenticate, check conflicting auth configs/env vars.
- **Codex config.toml modifications** — pen.dev may modify or duplicate the config (acknowledged issue); back up config before first use.
- **Unspecific prompts** ("Make it better") — AI makes unexpected changes; be specific and use version control to revert.

## Version Notes

- `.pen` document `version: "2.14"` (as of 2026-07-20 docs snapshot). The format is **live** — breaking changes possible; consult the TypeScript schema at docs.pen.dev/for-developers/the-pen-format for the authoritative reference.
- Design library files use `.lib.pen` suffix; conversion to library **cannot be undone**.
- Auto-save not yet available (per 2026-07-20 docs).
- Codex `config.toml` modification issue acknowledged and under investigation.
- Linux Wayland/Hyprland UI issues reported; X11 more stable.
