package services

import (
	"testing"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

func TestGetProjectsMapsContentFields(t *testing.T) {
	got := GetProjects()
	if len(got) != len(content.Projects) {
		t.Fatalf("unexpected project count: got %d want %d", len(got), len(content.Projects))
	}

	for i := range got {
		want := content.Projects[i]
		if got[i].Slug != want.Slug {
			t.Fatalf("unexpected slug at %d: got %q want %q", i, got[i].Slug, want.Slug)
		}
		if string(got[i].Body) != want.Body {
			t.Fatalf("unexpected body at %d", i)
		}
		if got[i].Date != want.Date.Format("January 2006") {
			t.Fatalf("unexpected date at %d: got %q", i, got[i].Date)
		}
	}
}

func TestGetFeaturedProjectsOnlyReturnsFeatured(t *testing.T) {
	got := GetFeaturedProjects()
	if len(got) == 0 {
		t.Fatal("expected at least one featured project")
	}
	for _, p := range got {
		if !p.Featured {
			t.Fatalf("expected only featured projects, found slug %q", p.Slug)
		}
	}
}

func TestGetProjectBySlugHandlesMissing(t *testing.T) {
	got, ok := GetProjectBySlug("cameronbrooks-site")
	if !ok {
		t.Fatal("expected existing project to be found")
	}
	if got.Slug != "cameronbrooks-site" {
		t.Fatalf("unexpected slug: got %q", got.Slug)
	}

	_, ok = GetProjectBySlug("does-not-exist")
	if ok {
		t.Fatal("expected missing project to return ok=false")
	}
}
