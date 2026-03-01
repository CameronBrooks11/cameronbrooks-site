package handlers

import (
	"fmt"
	"net/http"
)

// RobotsTxt handles GET /robots.txt.
func (h *Handler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User-agent: *\nAllow: /\n")
}

// SecurityTxt handles GET /.well-known/security.txt.
func (h *Handler) SecurityTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Contact: mailto:cambrooks3393@gmail.com\n")
	fmt.Fprint(w, "Preferred-Languages: en\n")
	fmt.Fprint(w, "Canonical: https://cameronbrooks.net/.well-known/security.txt\n")
}
