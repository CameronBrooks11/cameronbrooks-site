# UI/UX Plan

## Principles

1. **Content first.** Layout exists to serve the text and work, not the other way around.
2. **No JS required.** Every page is fully readable and navigable without JavaScript. HTMX is an enhancement, not a dependency.
3. **No external requests at page load.** No web fonts, no CDN CSS, no third-party scripts on the critical path. HTMX is vendored into `static/`.
4. **System fonts.** No CLS, no flash, no font load. Looks native on every OS.
5. **Flat specificity.** CSS is a single file layered as: tokens → reset → base → layout → components. No `!important`, no deep nesting.
6. **Design tokens first.** Every color, spacing, and type size is a CSS custom property. This makes future changes (dark mode, rebranding) a one-file edit.

---

## Visual Language

### Color tokens

Semantic token names throughout — no hue-based names (`--color-blue-500`). Tokens are the only place color values appear; nothing else in the CSS uses raw hex.

**Light mode:**

```css
:root {
  --color-bg: #ffffff;
  --color-surface: #f8fafc; /* card backgrounds, code blocks */
  --color-border: #e2e8f0;

  --color-text: #0f172a;
  --color-text-muted: #475569;

  --color-accent: #0f766e; /* links, active nav, buttons — restrained teal */
  --color-accent-hover: #115e59;
  --color-accent-soft: #ccfbf1; /* subtle teal tint for hover backgrounds */

  --color-success: #059669; /* emerald — semantic emphasis only */
  --color-highlight: #84cc16; /* lime — focus rings and micro-interactions only */
}
```

**Dark mode (planned, tokens defined now):**

```css
@media (prefers-color-scheme: dark) {
  :root {
    --color-bg: #020617;
    --color-surface: #0f172a;
    --color-border: #1e293b;

    --color-text: #e2e8f0;
    --color-text-muted: #94a3b8;

    --color-accent: #2dd4bf;
    --color-accent-hover: #5eead4;
    --color-accent-soft: #134e4a;

    --color-success: #34d399;
    --color-highlight: #a3e635;
  }
}
```

Dark mode tokens are defined now so any future activation is a display toggle, not a token redesign. Rendering in dark mode is deferred to post-v1.

The identity uses restrained teal as primary accent with emerald as a semantic secondary. Lime is used only for focus rings and micro-interaction highlights — never as a surface color or dominant visual element.

> **Visual intent:** the site should feel closer to high-quality engineering documentation or tooling interfaces (GitHub, Linear) — not a startup landing page or marketing palette. Color must never compete with content.

### Color usage principles

- **Neutral colors (~70% of surface area):** bg, surface, border, text — these dominate every page
- **Teal accent (~25%):** `--color-accent` for links, active nav state, interactive elements
- **Emerald (`--color-success`):** semantic emphasis only — success states, confirmations; not decorative
- **Lime (`--color-highlight`):** focus rings and micro-interaction highlights only; never as a background or block-level color
- No large colored backgrounds
- No gradients (except the optional progress bar)
- No sections combining multiple accent colors

### Typography

System font stack — no web fonts:

```css
--font-sans: ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif;
--font-mono: ui-monospace, "Cascadia Code", "Fira Code", monospace;
```

Type scale (unitless ratios, base 16px):

| Token         | Size     | Use                     |
| ------------- | -------- | ----------------------- |
| `--text-sm`   | 0.875rem | metadata, captions, nav |
| `--text-base` | 1rem     | body copy               |
| `--text-lg`   | 1.125rem | lead / intro paragraph  |
| `--text-xl`   | 1.25rem  | subheadings             |
| `--text-2xl`  | 1.5rem   | h2                      |
| `--text-3xl`  | 1.875rem | h1, page titles         |

Line height: `1.7` for body copy, `1.2` for headings. Letter-spacing: normal for body, slight negative (`-0.02em`) for large headings.

### Spacing

8px base unit. Not every integer — only the steps that have real use cases:

| Token        | Value | Primary use                     |
| ------------ | ----- | ------------------------------- |
| `--space-1`  | 8px   | Icon gaps, tight inline spacing |
| `--space-2`  | 16px  | Default padding, list item gaps |
| `--space-3`  | 24px  | Card padding, input gaps        |
| `--space-4`  | 32px  | Section internal padding        |
| `--space-6`  | 48px  | Major section breaks            |
| `--space-8`  | 64px  | Between page sections           |
| `--space-10` | 80px  | Page-level vertical rhythm      |

### Max-widths

| Purpose                                  | Width   |
| ---------------------------------------- | ------- |
| Reading content (articles, about)        | `680px` |
| Full site container (nav, footer, lists) | `900px` |

Centered with `margin: 0 auto` and horizontal padding of `--space-3` (24px) for breathing room on small viewports.

---

## Layout

### Structure

```txt
┌─────────────────────────────────────┐
│  <header> nav                       │
├─────────────────────────────────────┤
│  <main id="main">                   │
│    page content                     │
│                                     │
│                                     │
└─────────────────────────────────────┘
│  <footer>                           │
└─────────────────────────────────────┘
```

No sidebar. No grid. Single column throughout. Footer is minimal and not sticky.

### Navigation

- Sticky at top (`position: sticky; top: 0`), full width, subtle bottom border (`--color-border`)
- Sticky is preferred over `position: fixed` for a content site: no body padding-top hack, no anchor link offset bugs, the nav scrolls away on short pages and returns on scroll-up naturally
- Inner container constrained to `900px`
- Left: name / logo as text link to `/`
- Right: `Projects  Writing  About  Contact`
- Font: `--text-sm`, uppercase or regular weight — decide during implementation
- Active page link: `--color-accent`, no underline
- Hover: subtle text color shift to `--color-text` or a short underline; no colored nav backgrounds
- All other nav links: `--color-text-muted` idle

Five links fit inline at any realistic viewport width. No hamburger menu needed for this link count.

### Footer

Three lines maximum:

- Name + year
- Primary links (same as nav) or just a subset
- Optional: a "built with Go" line or nothing at all

---

## Page Patterns

### Home (`/`)

- Short intro paragraph — who you are, one or two sentences
- "Selected projects" section: 2–3 featured project cards
- "Recent writing" section: 3–5 list items (title + date)
- No hero image, no banner. Prose and links.

### Projects list (`/projects`)

- Page heading: "Projects"
- Optional one-line description
- List of project cards, stacked vertically
- Each card: title, one-sentence description, tag(s), link to slug
- No images required initially; add if a project warrants it

### Project detail (`/projects/:slug`)

- `<article>` element
- H1 title, metadata line (date, tags, links to source/demo)
- Prose body (Go struct field or markdown later)
- Back link to `/projects` at both top and bottom

### Writing list (`/writing`)

- Same pattern as projects list
- Each item: title, date, one-sentence summary
- Sorted newest first

### Writing detail (`/writing/:slug`)

- `<article>` element, reading width (`680px`)
- H1, date
- Body prose — this is where good line-height and measure matter most
- Code blocks use `--font-mono`, `--color-surface` background, horizontal scroll on overflow
- Back link to `/writing`

### About (`/about`)

- Single `<article>`, reading width
- Prose only — no resume/CV layout for now
- Optional: a photo (plain `<img>`, no special styling needed)

### Contact (`/contact`)

- Short paragraph
- Email as a `mailto:` link
- Links to GitHub, LinkedIn, or wherever relevant
- No form in v1

---

## HTMX Navigation

`hx-boost="true"` is set on `<body>` (or the layout wrapper). On every link click:

- HTMX fetches the new URL
- Swaps only `<main id="main">` (via `hx-select="#main" hx-target="#main"`)
- Pushes the URL to history

Handlers check the `HX-Request` header:

- If present: render only the `main` block partial
- If absent: render the full page (nav + main + footer)

This means **one template per page**, no separate "partial" templates. The layout template wraps or unwraps based on the request type.

**Loading state:** a minimal top progress bar using HTMX's `htmx:beforeRequest` / `htmx:afterRequest` events and a CSS transition. No library needed — a thin `<div id="progress-bar">` with a CSS width transition. Styling: `background: var(--color-accent)`, max height `2px`, optional subtle gradient to `--color-success` allowed. Must not exceed 2px height.

---

## CSS File Structure

Single file `static/css/main.css`, layered in this order:

```txt
1. Tokens          — custom properties
2. Reset           — box-sizing, margin/padding reset, img max-width
3. Base            — body, a, p, h1-h6, code, pre, ul/ol defaults
4. Layout          — .container, header, nav, main, footer
5. Components      — .card, .post-list, .tag, .btn, .back-link, #progress-bar
6. Utilities       — .sr-only (screen reader only), .text-muted, etc.
```

No scoping, no nesting beyond one level, no preprocessor. Keep the file as concise as practical; the current Phase 09 reference stylesheet is longer due explicit layering/comments and accessibility states.

---

## Accessibility baseline

- Semantic HTML throughout: `<header>`, `<nav>`, `<main>`, `<article>`, `<footer>`
- Heading hierarchy: one `<h1>` per page, logical `h2`/`h3` nesting
- Skip-to-content link as the first focusable element (visually hidden until focused)
- All images have `alt` text
- Focus styles are visible — do not suppress the default outline; use `--color-highlight` (lime) for focus rings to ensure distinction from accent and meet WCAG AA on both light and dark surfaces
- Color contrast: meet WCAG AA (4.5:1 for body text) for both light mode tokens against `--color-bg: #ffffff` and dark mode tokens against `--color-bg: #020617`; verify `--color-accent` remains readable on both surfaces

---

## What is explicitly deferred

| Item                           | Reason                                                    |
| ------------------------------ | --------------------------------------------------------- |
| Dark mode rendering            | Tokens are defined; activation deferred until post-v1     |
| Animations beyond progress bar | Unnecessary complexity for v1                             |
| Image galleries / lightboxes   | Not needed unless a project page warrants it              |
| Search                         | No content volume yet                                     |
| Comments                       | No need for a personal site initially                     |
| RSS feed                       | Worth adding soon after writing is live; it's one handler |
