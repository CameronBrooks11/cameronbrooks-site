package content

import "testing"

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
