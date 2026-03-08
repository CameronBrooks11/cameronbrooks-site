package content

import "embed"

// writingFS holds all Markdown post files embedded at compile time.
// Add a new post by creating a .md file in writing/ — no Go code changes needed.
//
//go:embed writing
var writingFS embed.FS
