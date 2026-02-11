package main

import (
    "html/template" // New import
    "path/filepath" // New import
	"time"
    "snippetbox.net/internal/models"
)

type templateData struct {
	CurrentYear int
    Snippet *models.Snippet
	Snippets []*models.Snippet
	Form        any
	Flash string 
	IsAuthenticated bool
	 CSRFToken       string
}

func humanDate(t time.Time) string {
    return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("/home/aartchik/project/golang/snippetbox/ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("/home/aartchik/project/golang/snippetbox/ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("/home/aartchik/project/golang/snippetbox/ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
    }
	return cache, nil
}
