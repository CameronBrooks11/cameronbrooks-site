# Phase 03 — Content Model & Data

**Goal:** Define the exact Go types for all site content, write 2–3 hardcoded placeholder entries for each type, and implement the four lookup helpers. The rest of the application is built on top of these types — they must be final before anything else touches them.

**Exit gate:** `go vet ./internal/content/...` passes; lookup functions return correct results when exercised manually; trust boundary is visible — `Body` is `string` throughout this entire package, no `template.HTML` anywhere.

---

## Prerequisites

- Phase 02 complete (`go build ./...` exits 0)
- `internal/content/` directory exists (created in Phase 01)
- `.gitkeep` in `internal/content/` can now be deleted — it will be replaced by real files

---

## Files to create in this phase

```
internal/content/project.go    — Link and Project types
internal/content/post.go       — Post type
internal/content/data.go       — hardcoded slice literals (placeholder content)
internal/content/lookup.go     — ProjectBySlug, FeaturedProjects, PublishedPosts, PostBySlug
```

---

## Step 1 — Delete `.gitkeep`

```sh
Remove-Item internal/content/.gitkeep
```

---

## Step 2 — `internal/content/project.go`

**File: `internal/content/project.go`**

```go
package content

import "time"

// Link is an external reference attached to a Project (source repo, live demo, write-up, etc.).
type Link struct {
	Label string // e.g. "Source", "Demo", "Write-up"
	URL   string
}

// Project represents a portfolio item.
// Body is stored as raw HTML string. The explicit conversion to template.HTML
// — which disables Go's auto-escaping — is performed in internal/services/ only.
type Project struct {
	Slug        string
	Title       string
	Description string    // one sentence; used in list/card view and meta description
	Body        string    // raw HTML; do NOT cast to template.HTML here or in handlers
	Tags        []string
	Date        time.Time // publication/completion date; used for sorting and display
	Links       []Link    // external references: source, demo, paper, etc.
	Featured    bool      // true = include on home page featured section
}
```

---

## Step 3 — `internal/content/post.go`

**File: `internal/content/post.go`**

```go
package content

import "time"

// Post represents a writing entry (blog post, article, note).
// Body is stored as raw HTML string. The explicit conversion to template.HTML
// — which disables Go's auto-escaping — is performed in internal/services/ only.
type Post struct {
	Slug      string
	Title     string
	Summary   string    // one sentence; used in list view and meta description
	Body      string    // raw HTML; do NOT cast to template.HTML here or in handlers
	Tags      []string
	Date      time.Time // publication date; used for sorting and display
	Published bool      // false = draft; excluded from all routes and lists
}
```

---

## Step 4 — `internal/content/data.go`

Placeholder entries. These exist so the site has something to render during development. Replace with real content before go-live (Phase 13).

Use at least **two projects** (one featured, one not) and **two posts** (one published, one draft) so every lookup path — featured filter, published filter, slug lookup — can be exercised.

**File: `internal/content/data.go`**

```go
package content

import "time"

// Projects is the full list of portfolio projects.
// Order does not matter — use FeaturedProjects() and sorting helpers as needed.
var Projects = []Project{
	{
		Slug:        "cameronbrooks-site",
		Title:       "cameronbrooks-site",
		Description: "This site — a minimal personal site built with Go, HTMX, and a single Debian VPS.",
		Body: `<p>A personal site built from scratch using Go's standard library, HTMX for
progressive enhancement, and hand-rolled CSS. No frameworks, no build steps,
no third-party dependencies.</p>
<p>Deployed as a single self-contained binary on a Debian VPS behind Caddy.</p>`,
		Tags:     []string{"go", "htmx", "css"},
		Date:     time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
		Links:    []Link{{Label: "Source", URL: "https://github.com/CameronBrooks11/cameronbrooks-site"}},
		Featured: true,
	},
	{
		Slug:        "placeholder-project",
		Title:       "Placeholder Project",
		Description: "A second project entry to exercise non-featured and list-view rendering.",
		Body:        `<p>Placeholder body. Replace with real content before go-live.</p>`,
		Tags:        []string{"placeholder"},
		Date:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Links:       nil,
		Featured:    false,
	},
}

// Posts is the full list of writing entries.
// Draft posts (Published: false) are never exposed via any route.
var Posts = []Post{
	{
		Slug:    "hello-world",
		Title:   "Hello, World",
		Summary: "The first post — a placeholder to exercise the writing routes and templates.",
		Body: `<p>This is a placeholder post. Replace with real content before go-live.</p>
<p>The date, tags, and published flag are all exercised by this entry.</p>`,
		Tags:      []string{"meta"},
		Date:      time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
		Published: true,
	},
	{
		Slug:      "draft-post",
		Title:     "Draft Post",
		Summary:   "A draft post — should never appear in any route or list.",
		Body:      `<p>This post is a draft and must not be reachable via any URL.</p>`,
		Tags:      []string{"draft"},
		Date:      time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC),
		Published: false,
	},
}
```

---

## Step 5 — `internal/content/lookup.go`

These are the only four surfaces the rest of the application calls. Handlers do not import this package directly — they go through `internal/services/`.

**File: `internal/content/lookup.go`**

```go
package content

// ProjectBySlug returns the project with the given slug, or false if not found.
// All projects are eligible regardless of Featured flag.
func ProjectBySlug(slug string) (Project, bool) {
	for _, p := range Projects {
		if p.Slug == slug {
			return p, true
		}
	}
	return Project{}, false
}

// FeaturedProjects returns all projects where Featured == true.
// Order follows the declaration order in data.go.
func FeaturedProjects() []Project {
	var out []Project
	for _, p := range Projects {
		if p.Featured {
			out = append(out, p)
		}
	}
	return out
}

// PublishedPosts returns all posts where Published == true.
// Order follows the declaration order in data.go; callers may re-sort by Date.
func PublishedPosts() []Post {
	var out []Post
	for _, p := range Posts {
		if p.Published {
			out = append(out, p)
		}
	}
	return out
}

// PostBySlug returns the post with the given slug if it exists and is published.
// Returns false for drafts — they are never reachable via slug lookup.
func PostBySlug(slug string) (Post, bool) {
	for _, p := range Posts {
		if p.Slug == slug && p.Published {
			return p, true
		}
	}
	return Post{}, false
}
```

---

## Step 6 — Verify

```sh
go vet ./internal/content/...
```

Expected: exits 0, no output.

```sh
go build ./...
```

Expected: exits 0. The `internal/content` package is not yet imported by anything — that is fine. `go build ./...` will compile it as part of the full build graph.

**Manually verify lookup logic** by running a quick inline test:

```sh
go run - <<'EOF'
package main

import (
    "fmt"
    "github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

func main() {
    featured := content.FeaturedProjects()
    fmt.Printf("featured: %d project(s)\n", len(featured))           // expect 1

    p, ok := content.ProjectBySlug("cameronbrooks-site")
    fmt.Printf("project by slug: %v, found=%v\n", p.Title, ok)       // expect found=true

    _, ok = content.ProjectBySlug("does-not-exist")
    fmt.Printf("missing slug: found=%v\n", ok)                       // expect found=false

    posts := content.PublishedPosts()
    fmt.Printf("published posts: %d\n", len(posts))                  // expect 1

    _, ok = content.PostBySlug("draft-post")
    fmt.Printf("draft post by slug: found=%v\n", ok)                 // expect found=false
}
EOF
```

All five assertions must match expectations before proceeding.

**Trust boundary audit:**

```sh
# Must return zero results — template.HTML must not appear in this package
grep -r "template.HTML" internal/content/
```

Expected: no output.

---

## Step 7 — Commit

```sh
git add internal/content/
git commit -m "phase 03: content model and data"
```

---

## Exit gate checklist

- [ ] `go vet ./internal/content/...` exits 0
- [ ] `go build ./...` exits 0
- [ ] `Projects` has ≥1 `Featured: true` entry and ≥1 `Featured: false` entry
- [ ] `Posts` has ≥1 `Published: true` entry and ≥1 `Published: false` entry
- [ ] Manual lookup verification: all five assertions pass
- [ ] `grep -r "template.HTML" internal/content/` returns no matches
- [ ] `.gitkeep` removed from `internal/content/`
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 04.
