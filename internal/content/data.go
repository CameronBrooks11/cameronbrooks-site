package content

import "time"

// Projects is the full list of portfolio projects.
// Order does not matter - use FeaturedProjects() and sorting helpers as needed.
var Projects = []Project{
	{
		Slug:        "cameronbrooks-site",
		Title:       "cameronbrooks-site",
		Description: "This site - a minimal personal site built with Go, HTMX, and a single Debian VPS.",
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
		Slug:        "go-request-middleware-kit",
		Title:       "Go Request Middleware Kit",
		Description: "A small middleware package for request IDs, structured logging, and production-safe defaults.",
		Body: `<p>A reusable middleware set built while iterating on personal and
internal services. The package focuses on predictable request tracing and
clean log output without adding a framework dependency.</p>
<p>It includes request ID propagation, status-aware logging, and a simple
composition helper for standard library handlers.</p>`,
		Tags:        []string{"go", "middleware", "observability"},
		Date:        time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
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
		Summary: "Why this site exists, what it is built with, and how I plan to use it.",
		Body: `<p>This is the first post on the site. It outlines the goals for this
space: documenting projects, sharing implementation notes, and writing clearly
about tradeoffs encountered while building.</p>
<p>The site intentionally stays simple: Go templates, HTMX for progressive
enhancement, and a deployment model that is easy to operate solo.</p>`,
		Tags:      []string{"meta"},
		Date:      time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
		Published: true,
	},
	{
		Slug:      "draft-post",
		Title:     "Draft Post",
		Summary:   "A draft post - should never appear in any route or list.",
		Body:      `<p>This post is a draft and must not be reachable via any URL.</p>`,
		Tags:      []string{"draft"},
		Date:      time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC),
		Published: false,
	},
}
