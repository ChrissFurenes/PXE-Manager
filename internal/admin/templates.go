package admin

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateManager struct {
	templates *template.Template
}

func NewTemplateManager(templateDir string) (*TemplateManager, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &TemplateManager{
		templates: tmpl,
	}, nil
}

func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tm.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
