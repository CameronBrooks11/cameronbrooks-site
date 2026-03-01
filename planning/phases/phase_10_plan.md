# Phase 10 — HTMX & Progress Bar

**Goal:** Vendor a pinned HTMX release into `static/htmx.min.js`, write `static/js/progress.js` to animate the progress bar using HTMX lifecycle events, and verify end-to-end that HTMX navigation works (partial `<main>` swaps, URL history pushes, no full reloads), the progress bar animates correctly, and the site is fully functional with JavaScript disabled.

**Exit gate:** Navigating between pages in browser shows no full reload (verify in DevTools Network tab — only a `fetch` for `<main>` content, no full document requests); URL updates in the address bar; progress bar animates and fades; disabling JavaScript in DevTools produces a fully functional plain HTML site.

---

## Prerequisites

- Phase 09 complete (CSS with `#progress-bar` `.loading` / `.complete` classes present)
- Phase 08 complete (`render()` checks `HX-Request` header and returns partial or full response)
- `layout.gohtml` has HTMX attributes on `<body>` and `<script src="/static/htmx.min.js">` reference
- `static/htmx.min.js` currently contains the Phase 02 stub comment (will be replaced)
- `static/js/progress.js` currently contains the Phase 02 stub comment (will be replaced)

---

## Files to write in this phase

```
static/htmx.min.js       — REPLACE stub with vendored HTMX (pinned version)
static/js/progress.js    — REPLACE stub with progress bar event handler
```

No Go code changes. No template changes.

---

## Step 1 — Vendor `htmx.min.js`

Download the minified HTMX file at a pinned version. Do not use a CDN reference — the file must be local.

**Pinned version: 2.0.4**

```sh
# PowerShell — download directly
Invoke-WebRequest `
  -Uri "https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js" `
  -OutFile "static/htmx.min.js"
```

Or from WSL/Git Bash:

```sh
curl -L -o static/htmx.min.js "https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js"
```

**Verify the download:**

```sh
# File should be ~50–60KB
(Get-Item static/htmx.min.js).Length / 1KB
```

Expected: approximately 50–60 KB. If the file is a few hundred bytes, the download likely returned an HTML error page — re-download.

```sh
# First line should be a minified JS comment, not HTML
Get-Content static/htmx.min.js -TotalCount 1
```

Expected: starts with `/* htmx.org` or a minified JS expression — not `<!DOCTYPE` or `<html`.

> **Why pin the version?** HTMX `hx-boost` behaviour and event names are stable within a major version but can change between majors. Version 2.x is used here. Do not upgrade without reviewing the HTMX v2 migration guide.

> **Do not add `static/htmx.min.js` to `.gitignore`.** It is a vendored dependency and must be committed — it is served from the embedded FS at runtime and must be in the binary.

---

## Step 2 — Write `static/js/progress.js`

The progress bar is driven by two HTMX events:

- `htmx:beforeRequest` — fired when HTMX begins a request; start the animation
- `htmx:afterRequest` — fired when the response arrives (success or error); complete and fade

The CSS classes `.loading` and `.complete` are already defined in Phase 09's `#progress-bar` block. This script only adds and removes those classes.

**File: `static/js/progress.js`**

```js
(function () {
  "use strict";

  var bar = document.getElementById("progress-bar");

  if (!bar) {
    // Progress bar element not found — nothing to do.
    // This should not happen in production; layout.gohtml always includes #progress-bar.
    return;
  }

  var completeTimer = null;

  // htmx:beforeRequest — HTMX is about to make a request.
  // Animate bar to 80% (CSS transition handles the movement).
  document.addEventListener("htmx:beforeRequest", function () {
    if (completeTimer) {
      clearTimeout(completeTimer);
      completeTimer = null;
    }
    bar.classList.remove("complete");
    // Force reflow so removing "complete" and adding "loading" are distinct frames.
    void bar.offsetWidth;
    bar.classList.add("loading");
  });

  // htmx:afterRequest — HTMX request completed (success or error).
  // Snap bar to 100%, then fade it out.
  document.addEventListener("htmx:afterRequest", function () {
    bar.classList.remove("loading");
    bar.classList.add("complete");
    // Remove "complete" after the CSS opacity transition finishes (0.3s defined in Phase 09).
    completeTimer = setTimeout(function () {
      bar.classList.remove("complete");
      completeTimer = null;
    }, 400);
  });
})();
```

**How this works with the Phase 09 CSS:**

| State            | Classes     | CSS result                                                                     |
| ---------------- | ----------- | ------------------------------------------------------------------------------ |
| Idle             | (none)      | `width: 0%; opacity: 0` — invisible                                            |
| Request started  | `.loading`  | `width: 80%; opacity: 1` — visible, animated via `transition: width 0.2s ease` |
| Request complete | `.complete` | `width: 100%; opacity: 0` — slides to full width and fades out                 |
| After fade       | (none)      | returns to idle state                                                          |

The `void bar.offsetWidth` line forces a browser reflow between removing `.complete` and adding `.loading`. Without it, rapid back-to-back navigations can appear to skip the animation.

---

## Step 3 — Verify layout.gohtml HTMX attributes

Open `internal/views/layout.gohtml` and confirm these four attributes are present on `<body>` (they were set in Phase 05 — this is a confirmation step, not an edit):

```html
<body
  hx-boost="true"
  hx-target="#main"
  hx-select="#main"
  hx-push-url="true"
></body>
```

| Attribute            | Effect                                                                                       |
| -------------------- | -------------------------------------------------------------------------------------------- |
| `hx-boost="true"`    | Intercepts all `<a>` clicks and converts them to HTMX fetch requests                         |
| `hx-target="#main"`  | Replaces the `<main id="main">` element with the response                                    |
| `hx-select="#main"`  | Extracts only `<main id="main">` from the full-page response (if partial render is not used) |
| `hx-push-url="true"` | Updates the browser URL and history on each navigation                                       |

Also confirm `<script src="/static/htmx.min.js">` appears **before** `</body>` (it does if Phase 05 was followed). HTMX must be loaded before any `htmx:*` event listeners.

Also confirm `<script src="/static/js/progress.js">` appears **after** the HTMX script tag.

---

## Step 4 — Verify the `render()` partial path

The HTMX navigation only works correctly if the server returns just the `<main>` inner content (not the full page) when `HX-Request: true` is sent.

Open `internal/handlers/render.go` and confirm this block is present (it was written in Phase 05):

```go
if r.Header.Get("HX-Request") == "true" {
    // HTMX partial: always implicit 200 so the fragment is swapped into #main.
    if err := tmplPart[page].ExecuteTemplate(w, "content", data); err != nil {
        slog.Error("partial render failed", "page", page, "err", err)
    }
    return
}
```

`tmplPart[page]` executes the `"content"` block only — it does not include the layout shell. The response body is just the inner HTML of `<main>`. HTMX then swaps that into the existing `<main id="main">` element, leaving the nav and footer untouched.

The `hx-select="#main"` attribute on `<body>` provides a fallback: if the server returns a full page (e.g. when JavaScript is disabled and the browser makes a normal GET), HTMX extracts `<main id="main">` from it. This is a belt-and-suspenders safety net — the partial render path is the primary path.

---

## Step 5 — End-to-end browser verification

Start the server:

```sh
$env:SITE_ENV="dev"; go run ./cmd/site
```

Open `http://localhost:8080` in Chrome or Firefox. Open DevTools → Network tab.

**HTMX navigation test:**

1. Click "Projects" in the nav
   - Network tab: one request to `/projects` with header `HX-Request: true`
   - Request type: `fetch` (not `document`)
   - Response: HTML fragment (no `<html>`, no `<head>`, no `<nav>`) — just the content block
   - URL in address bar: updates to `/projects`
   - Nav "Projects" link becomes active (teal)
   - Page content swaps **without a full reload** — nav and footer do not flicker

2. Click "Writing"
   - Same pattern — fetch, partial response, URL update

3. Click a project card to go to `/projects/cameronbrooks-site`
   - Same pattern

4. Click the "← Projects" back-link
   - Same pattern — not a full reload

5. Use browser Back button
   - Page navigates back using the pushed history entry
   - Network: one fetch request for the previous URL

**Progress bar test:**

On each navigation, the progress bar should:

1. Appear at the top (thin teal line, 2px, growing from left to ~80%)
2. Complete to full width
3. Fade out

If the network is local (localhost) the animation may be too fast to see clearly. Use DevTools → Network → throttle to "Slow 3G" temporarily to observe the animation.

---

## Step 6 — Verify back-link `<a>` tags work with hx-boost

`hx-boost` on `<body>` intercepts all `<a>` clicks on the same origin automatically. The back-links in `project.gohtml` and `post.gohtml` (`← Projects`, `← Writing`) should also be HTMX-boosted — confirm they do not trigger full reloads.

External links (`target="_blank"`) are not intercepted by `hx-boost` — correct behaviour.

---

## Step 7 — JavaScript disabled test

In DevTools → Settings → Preferences → Debugger → **Disable JavaScript**. (Or use a browser extension.)

With JS disabled:

- [ ] All pages still load fully (full-page renders, no HTMX)
- [ ] Nav links work as normal `<a>` tags
- [ ] Back-links work
- [ ] No broken UI, no spinner, no "JavaScript required" message
- [ ] `#progress-bar` remains invisible (no JS to trigger it — correct)

Re-enable JavaScript when done.

---

## Step 8 — Confirm `htmx.min.js` is served from embedded FS

With the server running:

```sh
curl -I http://localhost:8080/static/htmx.min.js
```

Expected response:

```
HTTP/1.1 200 OK
Content-Type: application/javascript
...
```

Status must be `200`, not `404`. If it is `404`, the file was not picked up by the embed directive — verify `static/htmx.min.js` exists on disk and `go build ./...` has been run since the file was downloaded (the embed happens at compile time; `go run` re-compiles automatically).

---

## Step 9 — Commit

```sh
git add static/htmx.min.js static/js/progress.js
git commit -m "phase 10: HTMX and progress bar"
```

---

## Exit gate checklist

- [ ] `static/htmx.min.js` is the real vendored file (~50–60KB), not the stub
- [ ] `static/htmx.min.js` first line starts with JS, not HTML
- [ ] `static/js/progress.js` contains the IIFE with `htmx:beforeRequest` and `htmx:afterRequest` listeners
- [ ] `layout.gohtml` has all four HTMX attributes on `<body>` (`hx-boost`, `hx-target`, `hx-select`, `hx-push-url`)
- [ ] `layout.gohtml` loads `htmx.min.js` before `progress.js`
- [ ] DevTools Network: clicking nav links sends `fetch` requests with `HX-Request: true` header
- [ ] DevTools Network: HTMX responses are HTML fragments (no `<html>` wrapper)
- [ ] URL in address bar updates on each navigation
- [ ] Browser Back button navigates correctly
- [ ] Progress bar animates on each navigation (visible with network throttling)
- [ ] JavaScript disabled: all pages load and navigate as plain HTML
- [ ] `GET /static/htmx.min.js` returns `200`
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 11.
