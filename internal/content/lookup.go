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
// Returns false for drafts - they are never reachable via slug lookup.
func PostBySlug(slug string) (Post, bool) {
	for _, p := range Posts {
		if p.Slug == slug && p.Published {
			return p, true
		}
	}
	return Post{}, false
}
