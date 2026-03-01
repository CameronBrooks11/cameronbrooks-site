# Phase 09 — CSS

**Goal:** Write the complete `static/css/main.css` as a single file with all six ordered layers: tokens, reset, base, layout, components, utilities. Every element the templates produce must be correctly styled. Navigation active state, skip link, focus rings, article body typography, code blocks, tags, back-links, and the progress bar placeholder are all in scope. Target: keep the file maintainable and close to the reference baseline size.

**Exit gate:** All six pages render correctly in a browser at `localhost:8080`; nav active state highlights on each page; skip link is visible on focus; focus rings are visible on interactive elements; `#progress-bar` is present in the DOM and styled; no raw-hex color values appear outside the tokens layer; file size remains within a practical baseline range (roughly 500-700 lines for this reference stylesheet).

---

## Prerequisites

- Phase 08 complete — the server is running and all routes respond
- `static/css/main.css` currently contains the Phase 02 one-line stub comment
- Browser DevTools available to inspect layout and contrast

---

## File to write in this phase

```
static/css/main.css   — REPLACE stub with full six-layer stylesheet
```

---

## Layer order rule

The six layers are written **top to bottom in this exact sequence** with a section comment header for each. No rule from a later layer references a variable not yet defined in the tokens layer. No `!important`. No nesting beyond one level (e.g. `nav a` is fine; `header nav ul li a:hover` is not).

---

## Complete `static/css/main.css`

**File: `static/css/main.css`**

```css
/* =============================================================================
   cameronbrooks-site — main.css
   Layers (in order): tokens → reset → base → layout → components → utilities
   Single file, no build step, no preprocessor. Target: concise and maintainable.
   ============================================================================= */

/* =============================================================================
   1. TOKENS — all design values as custom properties.
      Raw hex/values appear ONLY here. Everything else references a token.
   ============================================================================= */

:root {
  /* Color — light mode */
  --color-bg: #ffffff;
  --color-surface: #f8fafc; /* card backgrounds, code blocks, subtle fills */
  --color-border: #e2e8f0;

  --color-text: #0f172a;
  --color-text-muted: #475569;

  --color-accent: #0f766e; /* teal — links, active nav, interactive elements */
  --color-accent-hover: #115e59;
  --color-accent-soft: #ccfbf1; /* teal tint for hover backgrounds */

  --color-success: #059669; /* emerald — semantic emphasis only */
  --color-highlight: #84cc16; /* lime — focus rings only, never a surface color */

  /* Typography */
  --font-sans: ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif;
  --font-mono: ui-monospace, "Cascadia Code", "Fira Code", monospace;

  --text-sm: 0.875rem;
  --text-base: 1rem;
  --text-lg: 1.125rem;
  --text-xl: 1.25rem;
  --text-2xl: 1.5rem;
  --text-3xl: 1.875rem;

  /* Spacing — 8px base unit */
  --space-1: 8px;
  --space-2: 16px;
  --space-3: 24px;
  --space-4: 32px;
  --space-6: 48px;
  --space-8: 64px;
  --space-10: 80px;

  /* Max-widths */
  --width-content: 680px; /* articles, about, post detail */
  --width-container: 900px; /* nav, footer, list pages */
}

/* Dark mode — tokens defined now; rendering activation deferred to post-v1.
   To activate: remove this comment block and nothing else needs to change. */
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

/* =============================================================================
   2. RESET — minimal, purposeful. Only resets with a concrete reason.
   ============================================================================= */

*,
*::before,
*::after {
  box-sizing: border-box;
}

* {
  margin: 0;
  padding: 0;
}

img,
picture,
video,
canvas,
svg {
  display: block;
  max-width: 100%;
}

input,
button,
textarea,
select {
  font: inherit;
}

/* =============================================================================
   3. BASE — global element defaults. No classes.
   ============================================================================= */

html {
  font-size: 100%; /* respect browser default (typically 16px) */
  -webkit-text-size-adjust: 100%;
}

body {
  font-family: var(--font-sans);
  font-size: var(--text-base);
  line-height: 1.7;
  color: var(--color-text);
  background-color: var(--color-bg);
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
}

/* Links */
a {
  color: var(--color-accent);
  text-decoration: underline;
  text-underline-offset: 3px;
}

a:hover {
  color: var(--color-accent-hover);
}

a:focus-visible {
  outline: 2px solid var(--color-highlight);
  outline-offset: 3px;
  border-radius: 2px;
}

/* Headings */
h1,
h2,
h3,
h4,
h5,
h6 {
  line-height: 1.2;
  letter-spacing: -0.02em;
  color: var(--color-text);
}

h1 {
  font-size: var(--text-3xl);
  margin-bottom: var(--space-3);
}
h2 {
  font-size: var(--text-2xl);
  margin-bottom: var(--space-2);
}
h3 {
  font-size: var(--text-xl);
  margin-bottom: var(--space-1);
}

/* Body copy rhythm */
p {
  margin-bottom: var(--space-2);
}

p:last-child {
  margin-bottom: 0;
}

/* Lists */
ul,
ol {
  padding-left: var(--space-3);
  margin-bottom: var(--space-2);
}

li {
  margin-bottom: var(--space-1);
}

/* Inline code */
code {
  font-family: var(--font-mono);
  font-size: 0.9em;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 3px;
  padding: 0.1em 0.35em;
}

/* Code blocks */
pre {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: var(--space-3);
  overflow-x: auto;
  margin-bottom: var(--space-3);
  line-height: 1.5;
}

pre code {
  background: none;
  border: none;
  padding: 0;
  font-size: inherit;
}

hr {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: var(--space-6) 0;
}

/* =============================================================================
   4. LAYOUT — structural containers, header, nav, main, footer.
   ============================================================================= */

.container {
  max-width: var(--width-container);
  margin-inline: auto;
  padding-inline: var(--space-3);
}

/* Header & nav */
header {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: var(--color-bg);
  border-bottom: 1px solid var(--color-border);
}

header nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-block: var(--space-2);
}

.nav-home {
  font-weight: 600;
  font-size: var(--text-base);
  text-decoration: none;
  color: var(--color-text);
}

.nav-home:hover {
  color: var(--color-accent);
}

.nav-links {
  display: flex;
  gap: var(--space-3);
  list-style: none;
  padding: 0;
  margin: 0;
}

.nav-links a {
  font-size: var(--text-sm);
  text-decoration: none;
  color: var(--color-text-muted);
  transition: color 0.15s ease;
}

.nav-links a:hover {
  color: var(--color-text);
}

.nav-links a.active {
  color: var(--color-accent);
  font-weight: 500;
}

/* Main content area */
main {
  flex: 1;
  padding-block: var(--space-8);
}

/* Article reading width constraint */
article {
  max-width: var(--width-content);
}

.article-header {
  margin-bottom: var(--space-6);
}

.article-body {
  max-width: var(--width-content);
}

.article-body p,
.article-body li {
  line-height: 1.8;
}

.article-links {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
  margin-top: var(--space-1);
}

/* Footer */
footer {
  border-top: 1px solid var(--color-border);
  padding-block: var(--space-4);
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}

/* =============================================================================
   5. COMPONENTS — reusable UI patterns.
   ============================================================================= */

/* Card */
.card {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: var(--space-3);
  transition: border-color 0.15s ease;
}

.card:hover {
  border-color: var(--color-accent);
}

/* Post / project list */
.post-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* Tag pill */
.tag {
  display: inline-block;
  font-size: var(--text-sm);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  padding: 0.1em 0.5em;
  color: var(--color-text-muted);
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-1);
  margin-bottom: 0;
}

/* Back link */
.back-link {
  display: inline-block;
  font-size: var(--text-sm);
  text-decoration: none;
  color: var(--color-text-muted);
  margin-bottom: var(--space-4);
}

.back-link:hover {
  color: var(--color-accent);
}

/* Home page sections */
.home-intro {
  margin-bottom: var(--space-8);
}

.home-section {
  margin-bottom: var(--space-8);
}

.home-section h2 {
  margin-bottom: var(--space-3);
}

/* Progress bar — HTMX loading indicator (animated via JS in Phase 10) */
#progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  height: 2px;
  width: 0%;
  background-color: var(--color-accent);
  transition:
    width 0.2s ease,
    opacity 0.3s ease;
  opacity: 0;
  z-index: 9999;
  pointer-events: none;
}

#progress-bar.loading {
  opacity: 1;
  width: 80%;
}

#progress-bar.complete {
  width: 100%;
  opacity: 0;
}

/* =============================================================================
   6. UTILITIES — single-purpose helper classes.
   ============================================================================= */

/* Screen reader only — visually hidden, accessible to assistive technology */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

/* Skip-to-content link — sr-only until focused */
.skip-link {
  position: absolute;
  top: var(--space-2);
  left: var(--space-2);
  z-index: 100;
  padding: var(--space-1) var(--space-2);
  background-color: var(--color-bg);
  border: 2px solid var(--color-highlight);
  border-radius: 4px;
  font-size: var(--text-sm);
  text-decoration: none;
  color: var(--color-text);
  /* Hidden until focused */
  clip: rect(0 0 0 0);
  clip-path: inset(50%);
  overflow: hidden;
  white-space: nowrap;
}

.skip-link:focus {
  clip: auto;
  clip-path: none;
  overflow: visible;
  white-space: normal;
}

/* Text helpers */
.text-muted {
  color: var(--color-text-muted);
}

.text-lg {
  font-size: var(--text-lg);
  line-height: 1.6;
}
```

---

## Step 1 — Write the file

Replace `static/css/main.css` entirely with the content above.

Check line count after writing:

```sh
(Get-Content static/css/main.css).Count
```

For this reference stylesheet, expect roughly 500-700 lines. If it is far outside that range, verify nothing was accidentally duplicated or omitted.

---

## Step 2 — Verify in browser

Start the server (if not running):

```sh
$env:SITE_ENV="dev"; go run ./cmd/site
```

Open `http://localhost:8080` and walk through each check:

**Layout checks:**

- [ ] Nav is sticky — scroll down on a long page, nav stays at top
- [ ] Footer stays at the bottom even on short pages (flex column, `main` has `flex: 1`)
- [ ] Content is centered and does not stretch beyond `900px` container width
- [ ] Article pages (`/projects/cameronbrooks-site`, `/writing/hello-world`) have narrower body (`680px`)

**Nav checks:**

- [ ] `/projects` — "Projects" link is teal, others are muted
- [ ] `/writing` — "Writing" link is teal
- [ ] `/about` — "About" link is teal
- [ ] `/` — no nav link is active (home has `ActivePath: "/"` which does not match any `<li>` link)
- [ ] Hover on inactive nav link shifts text to darker; no background color change

**Typography checks:**

- [ ] Body text is comfortably readable, ~16px, line-height generous
- [ ] `h1` is visibly larger than `h2`; letter-spacing is slightly tight
- [ ] Code inline (`<code>`) has surface background and border
- [ ] Code block (`<pre><code>`) has surface background, scrolls horizontally if needed

**Component checks:**

- [ ] Project card on home page has surface background, border, border highlights teal on hover
- [ ] Tags render as small pills with border
- [ ] Back-link on project/post detail is small, muted, gains teal on hover
- [ ] `#progress-bar` div is in DOM with `height: 2px` (inspect with DevTools — should be present but invisible: `opacity: 0`, `width: 0%`)

**Accessibility checks:**

- [ ] Tab to the first element on any page — skip-to-content link becomes visible with lime border
- [ ] Tab through nav links — focus ring (lime outline) is visible on each
- [ ] Tab to any `<a>` in body copy — focus ring visible

---

## Step 3 — Contrast spot-check

Using browser DevTools or a contrast checker (e.g. `https://webaim.org/resources/contrastchecker/`):

| Pair                                                      | Ratio required | Check    |
| --------------------------------------------------------- | -------------- | -------- |
| `--color-text` (#0f172a) on `--color-bg` (#ffffff)        | 4.5:1 (AA)     | ~19:1 ✓  |
| `--color-text-muted` (#475569) on `--color-bg` (#ffffff)  | 4.5:1 (AA)     | ~5.9:1 ✓ |
| `--color-accent` (#0f766e) on `--color-bg` (#ffffff)      | 4.5:1 (AA)     | verify   |
| `--color-accent` (#0f766e) on `--color-surface` (#f8fafc) | 4.5:1 (AA)     | verify   |

If `--color-accent` on white fails AA, adjust the hex value (darker teal) before proceeding. The token change is one line in the file.

---

## Step 4 — Token purity check

No raw hex values should appear outside the tokens layer:

```sh
# Should return only lines within the tokens section (before the RESET comment)
Select-String -Path static/css/main.css -Pattern "#[0-9a-fA-F]{3,6}" | Select-Object LineNumber, Line
```

All matches should be within approximately lines 1–75 (the tokens layer). Any match after that is a violation — replace it with the appropriate `var(--color-*)` reference.

---

## Step 5 — Commit

```sh
git add static/css/main.css
git commit -m "phase 09: CSS"
```

---

## Exit gate checklist

- [ ] `static/css/main.css` exists and is within the expected baseline range (roughly 500-700 lines)
- [ ] All six pages render without visible layout breakage in browser
- [ ] Nav active state works on all nav-linked pages
- [ ] Skip-to-content link is invisible at rest, visible on Tab focus
- [ ] Focus rings (lime outline) visible on all interactive elements
- [ ] `#progress-bar` is in the DOM, `height: 2px`, `opacity: 0` at rest
- [ ] No raw hex values outside the tokens layer
- [ ] `--color-accent` (#0f766e) on white passes WCAG AA contrast
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 10.
