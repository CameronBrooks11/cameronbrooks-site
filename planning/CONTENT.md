# Content Model

Defines the data structures for all site content. This is the source of truth for what fields exist, how they are typed, and how they are stored.

---

## Principles

- No database in v1. No file I/O at request time.
- All content is Go struct literals in `internal/content/`. Change content = change code = redeploy. Accepted trade-off for a low-volume personal site.
- The `Body` field is stored as `string` (raw HTML). The explicit conversion to `template.HTML` — which disables Go's auto-escaping — happens in `internal/services/` at the service boundary, not in the storage struct. This keeps XSS guarantees intact and makes the trust grant auditable.
- Every content type has a `Slug` field that matches the URL path segment exactly.

---

## Types

### Project

```go
// internal/content/project.go

package content

import "time"

type Link struct {
    Label string // "Source", "Demo", "Write-up", etc.
    URL   string
}

type Project struct {
    Slug        string
    Title       string
    Description string    // one sentence — used in list/card view
    Body        string    // raw HTML; trust conversion to template.HTML in internal/services/
    Tags        []string
    Date        time.Time // used for sorting and display; format in template with .Format
    Links       []Link    // external links: source, demo, paper, etc.
    Featured    bool      // whether to show on home page
}
```

### Post (writing)

```go
// internal/content/post.go

package content

import "time"

type Post struct {
    Slug        string
    Title       string
    Summary     string    // one sentence — used in list view
    Body        string    // raw HTML; trust conversion to template.HTML in internal/services/
    Tags        []string
    Date        time.Time // used for sorting and display; format in template with .Format
    Published   bool      // false = draft, excluded from all lists and routes
}
```

---

## Storage

### v1: hardcoded Go slice literals

```go
// internal/content/data.go

package content

import "time"

var Projects = []Project{
    {
        Slug:        "my-project",
        Title:       "My Project",
        Description: "A one-sentence description.",
        Body:        "<p>Full description as HTML.</p>",
        Tags:        []string{"go", "cli"},
        Date:        time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC),
        Links:       []Link{{Label: "Source", URL: "https://github.com/..."}},
        Featured:    true,
    },
}

var Posts = []Post{
    {
        Slug:      "my-first-post",
        Title:     "My First Post",
        Summary:   "A one-sentence summary.",
        Body:      "<p>Full post body as HTML.</p>",
        Tags:      []string{"go"},
        Date:      time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
        Published: true,
    },
}
```

### Lookup helpers

```go
// internal/content/lookup.go

package content

func ProjectBySlug(slug string) (Project, bool) {
    for _, p := range Projects {
        if p.Slug == slug {
            return p, true
        }
    }
    return Project{}, false
}

func FeaturedProjects() []Project {
    var out []Project
    for _, p := range Projects {
        if p.Featured {
            out = append(out, p)
        }
    }
    return out
}

func PublishedPosts() []Post {
    var out []Post
    for _, p := range Posts {
        if p.Published {
            out = append(out, p)
        }
    }
    return out
}

func PostBySlug(slug string) (Post, bool) {
    for _, p := range Posts {
        if p.Slug == slug && p.Published {
            return p, true
        }
    }
    return Post{}, false
}
```

These functions are called from `internal/services/`, not directly from handlers. The service layer is responsible for constructing view models from these types, including the explicit `Body string` → `template.HTML` conversion.

---

## Future: markdown-backed content

When content volume makes HTML literals impractical, the migration path is:

1. Add `.md` files to `internal/content/posts/` and `internal/content/projects/`
2. Embed them with `//go:embed` — no runtime file I/O
3. Parse and render markdown to HTML string once at startup (e.g. with `goldmark`), store in `Body`
4. Populate `Posts` and `Projects` slices at startup instead of literal declaration

The `Post`, `Project`, and lookup helper signatures remain unchanged. No handler or template changes required.

---

## Content field reference

| Field                     | Type        | Required | Notes                                                                     |
| ------------------------- | ----------- | -------- | ------------------------------------------------------------------------- |
| `Slug`                    | `string`    | Yes      | URL-safe, lowercase, hyphens, no slashes                                  |
| `Title`                   | `string`    | Yes      | Used in `<h1>` and `<title>`                                              |
| `Description` / `Summary` | `string`    | Yes      | Used in list views and `<meta name="description">`                        |
| `Body`                    | `string`    | Yes      | Raw HTML in storage; converted to `template.HTML` in `internal/services/` |
| `Tags`                    | `[]string`  | No       | Displayed as labels; not linked in v1                                     |
| `Date`                    | `time.Time` | Yes      | Sorting and display; format in template with `.Format("Jan 2006")`        |
| `Links`                   | `[]Link`    | No       | Project only; external references                                         |
| `Featured`                | `bool`      | No       | Project only; controls home page inclusion                                |
| `Published`               | `bool`      | Yes      | Post only; false = draft, never exposed                                   |
