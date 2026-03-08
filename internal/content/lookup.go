package content

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
// Returns false for drafts - they are never reachable via slug lookup.
func PostBySlug(slug string) (Post, bool) {
	for _, p := range Posts {
		if p.Slug == slug && p.Published {
			return p, true
		}
	}
	return Post{}, false
}
