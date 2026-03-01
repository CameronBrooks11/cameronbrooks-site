package handlers

// Handler holds application-level dependencies shared across all page handlers.
// Currently empty; may hold config or service interfaces later.
type Handler struct{}

// New returns an initialized Handler ready to register on a mux.
func New() *Handler {
	return &Handler{}
}
