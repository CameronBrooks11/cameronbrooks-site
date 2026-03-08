package handlers

import (
	"net/http"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
)

// Home handles GET /.
// Renders the home page with featured projects and recent posts.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// The stdlib mux routes "/" as a catch-all for unmatched paths.
	// Return 404 for everything except exactly "/".
	if r.URL.Path != "/" {
		notFound(w, r)
		return
	}

	render(w, r, "home", http.StatusOK, PageData{
		Description: "Cameron Brooks - engineer & builder. Writing.",
		ActivePath:  "/",
		Data: HomeData{
			Recent: services.GetRecentPosts(5),
		},
	})
}
