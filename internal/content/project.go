package content

import "time"

// Link is an external reference attached to a Project (source repo, live demo, write-up, etc.).
type Link struct {
	Label string // e.g. "Source", "Demo", "Write-up"
	URL   string
}

// Project represents a portfolio item.
// Body is stored as raw HTML string. The explicit trust conversion to a safe
// HTML type is performed in internal/services/ only.
type Project struct {
	Slug        string
	Title       string
	Description string    // one sentence; used in list/card view and meta description
	Body        string    // raw HTML; do NOT perform trust conversion here or in handlers
	Tags        []string
	Date        time.Time // publication/completion date; used for sorting and display
	Links       []Link    // external references: source, demo, paper, etc.
	Featured    bool      // true = include on home page featured section
}
