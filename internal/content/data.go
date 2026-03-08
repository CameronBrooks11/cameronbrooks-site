package content

import "time"

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
