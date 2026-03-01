package views

import "embed"

// FS holds all .gohtml template files embedded at compile time.
// Parsed once at startup by handlers.InitTemplates().
//
//go:embed *.gohtml
var FS embed.FS
