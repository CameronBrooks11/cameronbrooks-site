package handlers

import (
	"net/http"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/services"
)

// Projects handles GET /projects.
func (h *Handler) Projects(w http.ResponseWriter, r *http.Request) {
	render(w, r, "projects", http.StatusOK, PageData{
		Title:       "Projects",
		Description: "A selection of things I have built.",
		ActivePath:  "/projects",
		Data:        services.GetProjects(),
	})
}

// Project handles GET /projects/{slug}.
func (h *Handler) Project(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	project, ok := services.GetProjectBySlug(slug)
	if !ok {
		notFound(w, r)
		return
	}

	render(w, r, "project", http.StatusOK, PageData{
		Title:       project.Title,
		Description: project.Description,
		ActivePath:  "/projects",
		Data:        project,
	})
}
