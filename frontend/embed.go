package frontend

import (
	"embed"
	"html/template"
)

//go:embed static templates
var FS embed.FS

func LoadTemplates() (*template.Template, error) {
	return template.ParseFS(FS, "templates/*.html")
}
