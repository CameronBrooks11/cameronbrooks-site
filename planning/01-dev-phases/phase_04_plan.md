# Phase 04 — Services Layer

**Goal:** Build `internal/services/` — the single place where content lookups are called, view models are constructed, and the `Body string` → `template.HTML` trust conversion is made. Handlers in Phase 07 will call only these functions; they will never import `internal/content/` directly.

**Exit gate:** `go vet ./internal/services/...` passes; `grep -r "template.HTML(" internal/` returns matches only in `internal/services/`; `go build ./...` exits 0.

---

## Prerequisites

- Phase 03 complete (`internal/content/` types and lookup helpers exist)
- `internal/services/` directory exists (created in Phase 01); `.gitkeep` can be deleted

---

## Why a service layer

Handlers must be thin (ARCHITECTURE.md invariant #6). The trust conversion from `Body string` → `template.HTML` must happen at a single auditable location (invariant #4). Without a service layer, both of those concerns end up in handlers, which violates both invariants.

The service layer is also the right place to add future concerns — sorting, filtering, pagination, viewmodel shaping — without touching handlers or content types.

---

## Files to create in this phase

```
internal/services/projects.go   — project view models and lookup wrappers
internal/services/posts.go      — post view models and lookup wrappers
```

---

## Step 1 — Delete `.gitkeep`

```sh
Remove-Item internal/services/.gitkeep
```

---

## Step 2 — View model types

Each service function returns a **view model** — a struct shaped for template consumption, not raw storage. The key difference from the storage type is that `Body` is `template.HTML`, not `string`. Auto-escaping is disabled for this field at the point of conversion.

The view models are defined in the same files as the service functions. They do not need their own package.

---

## Step 3 — `internal/services/projects.go`

**File: `internal/services/projects.go`**

```go
package services

import (
	"html/template"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

// ProjectView is the view model for a single project, safe for template execution.
// Body is template.HTML — auto-escaping is disabled. This conversion is made
// here and nowhere else in the application (see ARCHITECTURE.md invariant #4).
type ProjectView struct {
	Slug        string
	Title       string
	Description string
	Body        template.HTML // trust grant: raw HTML from content.Project.Body
	Tags        []string
	Date        string // pre-formatted for display; use content.Project.Date.Format in the service
	Links       []content.Link
	Featured    bool
}

// toProjectView converts a storage Project to a view-ready ProjectView.
// This is the only place template.HTML conversion occurs for projects.
func toProjectView(p content.Project) ProjectView {
	return ProjectView{
		Slug:        p.Slug,
		Title:       p.Title,
		Description: p.Description,
		Body:        template.HTML(p.Body), // #nosec — trust grant; Body is authored content
		Tags:        p.Tags,
		Date:        p.Date.Format("January 2006"),
		Links:       p.Links,
		Featured:    p.Featured,
	}
}

// GetProjects returns all projects as view models.
// Order follows declaration order in content/data.go.
func GetProjects() []ProjectView {
	projects := content.Projects
	out := make([]ProjectView, len(projects))
	for i, p := range projects {
		out[i] = toProjectView(p)
	}
	return out
}

// GetFeaturedProjects returns featured projects as view models, for the home page.
func GetFeaturedProjects() []ProjectView {
	featured := content.FeaturedProjects()
	out := make([]ProjectView, len(featured))
	for i, p := range featured {
		out[i] = toProjectView(p)
	}
	return out
}

// GetProjectBySlug returns a single project view model and a found flag.
// Returns false if the slug does not match any project.
func GetProjectBySlug(slug string) (ProjectView, bool) {
	p, ok := content.ProjectBySlug(slug)
	if !ok {
		return ProjectView{}, false
	}
	return toProjectView(p), true
}
```

---

## Step 4 — `internal/services/posts.go`

**File: `internal/services/posts.go`**

```go
package services

import (
	"html/template"
	"slices"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

// PostView is the view model for a single post, safe for template execution.
// Body is template.HTML — auto-escaping is disabled. This conversion is made
// here and nowhere else in the application (see ARCHITECTURE.md invariant #4).
type PostView struct {
	Slug      string
	Title     string
	Summary   string
	Body      template.HTML // trust grant: raw HTML from content.Post.Body
	Tags      []string
	Date      string // pre-formatted for display
	Published bool
}

// toPostView converts a storage Post to a view-ready PostView.
// This is the only place template.HTML conversion occurs for posts.
func toPostView(p content.Post) PostView {
	return PostView{
		Slug:      p.Slug,
		Title:     p.Title,
		Summary:   p.Summary,
		Body:      template.HTML(p.Body), // #nosec — trust grant; Body is authored content
		Tags:      p.Tags,
		Date:      p.Date.Format("January 2006"),
		Published: p.Published,
	}
}

// GetPosts returns all published posts as view models, sorted newest first.
func GetPosts() []PostView {
	published := content.PublishedPosts()
	// Sort by date descending before converting, so the view models are ordered correctly.
	slices.SortFunc(published, func(a, b content.Post) int {
		return b.Date.Compare(a.Date) // newest first
	})
	out := make([]PostView, len(published))
	for i, p := range published {
		out[i] = toPostView(p)
	}
	return out
}

// GetRecentPosts returns the n most recent published posts, for the home page.
func GetRecentPosts(n int) []PostView {
	all := GetPosts()
	if n > len(all) {
		n = len(all)
	}
	return all[:n]
}

// GetPostBySlug returns a single post view model and a found flag.
// Returns false if the slug does not match any published post.
func GetPostBySlug(slug string) (PostView, bool) {
	p, ok := content.PostBySlug(slug)
	if !ok {
		return PostView{}, false
	}
	return toPostView(p), true
}

```

---

## Step 5 — Verify

```sh
go vet ./internal/services/...
```

Expected: exits 0.

```sh
go build ./...
```

Expected: exits 0.

**Trust boundary audit — the critical check for this phase:**

```sh
# Must return matches only in internal/services/
grep -rn "template\.HTML(" internal/
```

Expected output should contain only lines from `internal/services/projects.go` and `internal/services/posts.go`. If any match appears in `internal/content/` or `internal/handlers/`, the trust boundary is violated.

---

## Step 6 — Commit

```sh
git add internal/services/
git commit -m "phase 04: services layer"
```

---

## Exit gate checklist

- [ ] `go vet ./internal/services/...` exits 0
- [ ] `go build ./...` exits 0
- [ ] `grep -rn "template.HTML(" internal/` returns matches **only** in `internal/services/`
- [ ] `GetProjectBySlug("does-not-exist")` returns `false` (no panic)
- [ ] `GetPostBySlug("draft-post")` returns `false` (drafts are not exposed)
- [ ] `GetPosts()` and `GetRecentPosts(n)` return only published posts
- [ ] `.gitkeep` removed from `internal/services/`
- [ ] `git status` is clean after commit

All boxes checked → proceed to Phase 05.
