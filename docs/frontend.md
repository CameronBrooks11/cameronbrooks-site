# Frontend Contracts

## Rendering model

- Server-rendered templates using `html/template`
- No client-side state store
- HTMX used for progressive navigation only

## Template contract

Each page template defines:

```gohtml
{{define "content"}} ... {{end}}
```

Layout wraps it:

```gohtml
<main id="main">{{template "content" .}}</main>
```

Full requests render layout; HTMX requests render content block only.

## HTMX contract

Layout-level attributes:

- `hx-boost="true"`
- `hx-target="#main"`
- `hx-push-url="true"`

Do not add `hx-select="#main"` with current partial-render strategy.

## Route-to-template map

- `/` -> `home.gohtml`
- `/projects` -> `projects.gohtml`
- `/projects/{slug}` -> `project.gohtml`
- `/writing` -> `writing.gohtml`
- `/writing/{slug}` -> `post.gohtml`
- `/about` -> `about.gohtml`
- `/contact` -> `contact.gohtml`
- `404` -> `notFound.gohtml`
- `500` -> `error.gohtml`

## Page data contract

`PageData` carries:

- `Title`
- `Description`
- `Year`
- `ActivePath`
- `Data`

`ActivePath` drives nav active-state highlighting.

## Style system

- Single stylesheet: `static/css/main.css`
- Token-first CSS custom properties
- System font stack only
- Single-column content layout
- Readable widths:
  - `680px` article width
  - `900px` container width

## Accessibility baseline

- Semantic landmarks (`header/nav/main/article/footer`)
- One `h1` per page
- Skip link to `#main`
- Visible focus styles
- WCAG AA contrast targets

## No-JS requirement

All core flows (navigation and content) must work with JavaScript disabled.
