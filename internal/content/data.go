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
		Summary: "The first post - a placeholder to exercise the writing routes and templates.",
		Body: `<p>This is a placeholder post. Replace with real content before go-live.</p>
<p>The date, tags, and published flag are all exercised by this entry.</p>`,
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
