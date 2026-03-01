package services

import (
	"html/template"
	"slices"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

// PostView is the view model for a single post, safe for template execution.
// Body is template.HTML; auto-escaping is disabled. This conversion is made
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
		Body:      template.HTML(p.Body), // #nosec - trust grant; Body is authored content
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
	if n <= 0 {
		return []PostView{}
	}
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
