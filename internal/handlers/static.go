package handlers

import "net/http"

// About handles GET /about.
func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	render(w, r, "about", http.StatusOK, PageData{
		Title:       "About",
		Description: "A little about Cameron Brooks.",
		ActivePath:  "/about",
	})
}

// Contact handles GET /contact.
func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	render(w, r, "contact", http.StatusOK, PageData{
		Title:       "Contact",
		Description: "How to reach Cameron Brooks.",
		ActivePath:  "/contact",
	})
}
