package content

import "testing"

func TestProjectBySlug(t *testing.T) {
	got, ok := ProjectBySlug("cameronbrooks-site")
	if !ok {
		t.Fatal("expected project to be found")
	}
	if got.Slug != "cameronbrooks-site" {
		t.Fatalf("unexpected slug: got %q", got.Slug)
	}

	_, ok = ProjectBySlug("does-not-exist")
	if ok {
		t.Fatal("expected missing project lookup to return ok=false")
	}
}

func TestFeaturedProjects(t *testing.T) {
	got := FeaturedProjects()
	if len(got) == 0 {
		t.Fatal("expected at least one featured project")
	}
	for _, p := range got {
		if !p.Featured {
			t.Fatalf("expected only featured projects, found slug %q", p.Slug)
		}
	}
}

func TestPublishedPosts(t *testing.T) {
	got := PublishedPosts()
	if len(got) == 0 {
		t.Fatal("expected at least one published post")
	}
	for _, p := range got {
		if !p.Published {
			t.Fatalf("expected only published posts, found slug %q", p.Slug)
		}
	}
}

func TestPostBySlugExcludesDrafts(t *testing.T) {
	_, ok := PostBySlug("hello-world")
	if !ok {
		t.Fatal("expected published post to be found")
	}

	_, ok = PostBySlug("draft-post")
	if ok {
		t.Fatal("expected draft post lookup to return ok=false")
	}
}
