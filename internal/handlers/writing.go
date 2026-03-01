package handlers

import (
	"net/http"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
)

// Writing handles GET /writing.
func (h *Handler) Writing(w http.ResponseWriter, r *http.Request) {
	render(w, r, "writing", http.StatusOK, PageData{
		Title:       "Writing",
		Description: "Notes and longer pieces.",
		ActivePath:  "/writing",
		Data:        services.GetPosts(),
	})
}

// Post handles GET /writing/{slug}.
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, ok := services.GetPostBySlug(slug)
	if !ok {
		notFound(w, r)
		return
	}

	render(w, r, "post", http.StatusOK, PageData{
		Title:       post.Title,
		Description: post.Summary,
		ActivePath:  "/writing",
		Data:        post,
	})
}
