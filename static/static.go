package static

import "embed"

// FS holds all static assets embedded at compile time.
// Served via http.FileServer(http.FS(static.FS)) at /static/.
//
//go:embed css js images favicon.ico htmx.min.js
var FS embed.FS
