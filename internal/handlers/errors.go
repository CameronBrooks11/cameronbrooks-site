package handlers

import "net/http"

// notFound writes a 404 response using the notFound template.
func notFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, "notFound", http.StatusNotFound, PageData{
		Title:       "Not Found",
		Description: "The page you're looking for doesn't exist.",
	})
}

// internalError writes a 500 response using the error template.
func internalError(w http.ResponseWriter, r *http.Request) {
	render(w, r, "error", http.StatusInternalServerError, PageData{
		Title:       "Error",
		Description: "An internal error occurred.",
	})
}
