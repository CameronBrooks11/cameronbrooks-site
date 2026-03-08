package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
	"github.com/CameronBrooks11/cameronbrooks-site/internal/views"
)

// PageData is passed as the root data object to every template execution.
type PageData struct {
	Title       string // used in <title> and <h1>; leave empty on home page
	Description string // used in <meta name="description">
	Year        string // current year for footer copyright; injected automatically by render()
	ActivePath  string // current URL path, e.g. "/projects"; used for nav active state
	Data        any    // page-specific payload; see per-page types defined below
}

// HomeData is the payload for the home page handler.
// Defined here because it is only used by handlers.Home.
type HomeData struct {
	Recent []services.PostView
}

// tmplFull maps page name -> template parsed with layout (used for full-page renders).
// tmplPart maps page name -> template parsed alone (used for HTMX partial renders).
var (
	tmplFull map[string]*template.Template
	tmplPart map[string]*template.Template
)

// pages lists every template name that must be present in the cache.
// Each name corresponds to a <name>.gohtml file in internal/views/.
var pages = []string{
	"home", "writing", "post",
	"about", "contact", "notFound", "error",
}

// InitTemplates parses all page templates at startup and populates tmplFull and tmplPart.
// Must be called once from main() before the HTTP server starts.
// Returns an error if any template file is missing or unparseable; treat as fatal.
func InitTemplates() error {
	tmplFull = make(map[string]*template.Template)
	tmplPart = make(map[string]*template.Template)

	layoutSrc, err := views.FS.ReadFile("layout.gohtml")
	if err != nil {
		return err
	}

	for _, name := range pages {
		pageSrc, err := views.FS.ReadFile(name + ".gohtml")
		if err != nil {
			return err
		}

		// Full: parse layout first, then associate the page template with it.
		full, err := template.New("layout").Parse(string(layoutSrc))
		if err != nil {
			return err
		}
		if _, err = full.New(name).Parse(string(pageSrc)); err != nil {
			return err
		}
		tmplFull[name] = full

		// Partial: parse the page template alone (no layout wrapper).
		part, err := template.New(name).Parse(string(pageSrc))
		if err != nil {
			return err
		}
		tmplPart[name] = part
	}
	return nil
}

// render executes the named template and writes the response.
// status is the HTTP status code for full-page renders (e.g. 200, 404).
// If the request carries an HX-Request header, only the "content" block is
// executed and status remains 200 so HTMX swaps the fragment into #main.
// Year is injected automatically; callers set Title, Description, ActivePath, and Data.
func render(w http.ResponseWriter, r *http.Request, page string, status int, data PageData) {
	data.Year = strconv.Itoa(time.Now().Year())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		if err := tmplPart[page].ExecuteTemplate(w, "content", data); err != nil {
			slog.Error("partial render failed", "page", page, "err", err)
		}
		return
	}

	w.WriteHeader(status)
	if err := tmplFull[page].ExecuteTemplate(w, "layout", data); err != nil {
		slog.Error("full render failed", "page", page, "err", err)
	}
}
