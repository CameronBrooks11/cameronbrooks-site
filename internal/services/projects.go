package services

import (
	"html/template"

	"github.com/CameronBrooks11/cameronbrooks-site/internal/content"
)

// ProjectView is the view model for a single project, safe for template execution.
// Body is template.HTML; auto-escaping is disabled. This conversion is made
// here and nowhere else in the application (see ARCHITECTURE.md invariant #4).
type ProjectView struct {
	Slug        string
	Title       string
	Description string
	Body        template.HTML // trust grant: raw HTML from content.Project.Body
	Tags        []string
	Date        string // pre-formatted for display; use content.Project.Date.Format in the service
	Links       []content.Link
	Featured    bool
}

// toProjectView converts a storage Project to a view-ready ProjectView.
// This is the only place template.HTML conversion occurs for projects.
func toProjectView(p content.Project) ProjectView {
	return ProjectView{
		Slug:        p.Slug,
		Title:       p.Title,
		Description: p.Description,
		Body:        template.HTML(p.Body), // #nosec - trust grant; Body is authored content
		Tags:        p.Tags,
		Date:        p.Date.Format("January 2006"),
		Links:       p.Links,
		Featured:    p.Featured,
	}
}

// GetProjects returns all projects as view models.
// Order follows declaration order in content/data.go.
func GetProjects() []ProjectView {
	projects := content.Projects
	out := make([]ProjectView, len(projects))
	for i, p := range projects {
		out[i] = toProjectView(p)
	}
	return out
}

// GetFeaturedProjects returns featured projects as view models, for the home page.
func GetFeaturedProjects() []ProjectView {
	featured := content.FeaturedProjects()
	out := make([]ProjectView, len(featured))
	for i, p := range featured {
		out[i] = toProjectView(p)
	}
	return out
}

// GetProjectBySlug returns a single project view model and a found flag.
// Returns false if the slug does not match any project.
func GetProjectBySlug(slug string) (ProjectView, bool) {
	p, ok := content.ProjectBySlug(slug)
	if !ok {
		return ProjectView{}, false
	}
	return toProjectView(p), true
}
