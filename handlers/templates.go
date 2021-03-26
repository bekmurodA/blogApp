package handlers

import (
	"html/template"
	"path/filepath"
	"snippetbox/forms"
	"snippetbox/models"
	"time"
)

type TemplateData struct {
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	CurrentYear       int
	Form              *forms.Form
	Flash             string
	AuthenticatedUser *models.User
}

func HumanDate(t time.Time) string {
	if t.IsZero(){
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{"humanDate": HumanDate}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	//new map to act as a cache
	cache := map[string]*template.Template{}

	//*.page.tmpl templates
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	//loop through the pages one by one
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}
	return cache, nil
}
