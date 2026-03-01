package handlers

// Handler holds application-level dependencies shared across all page handlers.
type Handler struct {
	AppVersion   string // injected from main.Version (set via -ldflags at build time)
	AppBuildTime string // injected from main.BuildTime (set via -ldflags at build time)
}

// New returns an initialized Handler. Set AppVersion and AppBuildTime after construction
// if build metadata is needed (see cmd/site/main.go).
func New() *Handler {
	return &Handler{}
}
