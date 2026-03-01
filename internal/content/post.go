package content

import "time"

// Post represents a writing entry (blog post, article, note).
// Body is stored as raw HTML string. The explicit trust conversion to a safe
// HTML type is performed in internal/services/ only.
type Post struct {
	Slug      string
	Title     string
	Summary   string    // one sentence; used in list view and meta description
	Body      string    // raw HTML; do NOT perform trust conversion here or in handlers
	Tags      []string
	Date      time.Time // publication date; used for sorting and display
	Published bool      // false = draft; excluded from all routes and lists
}
