package services

import (
	"slices"
	"testing"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

func TestGetPostsReturnsPublishedSortedNewestFirst(t *testing.T) {
	want := content.PublishedPosts()
	slices.SortFunc(want, func(a, b content.Post) int {
		return b.Date.Compare(a.Date)
	})

	got := GetPosts()
	if len(got) != len(want) {
		t.Fatalf("unexpected post count: got %d want %d", len(got), len(want))
	}

	for i := range got {
		if got[i].Slug != want[i].Slug {
			t.Fatalf("unexpected order at %d: got %q want %q", i, got[i].Slug, want[i].Slug)
		}
		if !got[i].Published {
			t.Fatalf("expected only published posts, found slug %q", got[i].Slug)
		}
	}
}

func TestGetRecentPostsClampsToAvailablePosts(t *testing.T) {
	all := GetPosts()

	recentOne := GetRecentPosts(1)
	if len(recentOne) != 1 {
		t.Fatalf("unexpected recent count: got %d want 1", len(recentOne))
	}
	if recentOne[0].Slug != all[0].Slug {
		t.Fatalf("unexpected recent item: got %q want %q", recentOne[0].Slug, all[0].Slug)
	}

	recentMany := GetRecentPosts(999)
	if len(recentMany) != len(all) {
		t.Fatalf("unexpected clamped count: got %d want %d", len(recentMany), len(all))
	}
}

func TestGetPostBySlugExcludesDrafts(t *testing.T) {
	got, ok := GetPostBySlug("hello-world")
	if !ok {
		t.Fatal("expected published post to be found")
	}
	if got.Slug != "hello-world" {
		t.Fatalf("unexpected slug: got %q", got.Slug)
	}

	_, ok = GetPostBySlug("draft-post")
	if ok {
		t.Fatal("expected draft post to return ok=false")
	}
}
