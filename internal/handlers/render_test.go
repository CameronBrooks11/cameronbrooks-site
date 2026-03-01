package handlers

import "testing"

func TestInitTemplatesLoadsAllPages(t *testing.T) {
	if err := InitTemplates(); err != nil {
		t.Fatalf("InitTemplates returned error: %v", err)
	}

	if len(tmplFull) != len(pages) {
		t.Fatalf("unexpected tmplFull size: got %d want %d", len(tmplFull), len(pages))
	}
	if len(tmplPart) != len(pages) {
		t.Fatalf("unexpected tmplPart size: got %d want %d", len(tmplPart), len(pages))
	}

	for _, name := range pages {
		if tmplFull[name] == nil {
			t.Fatalf("tmplFull missing key %q", name)
		}
		if tmplPart[name] == nil {
			t.Fatalf("tmplPart missing key %q", name)
		}
	}
}
