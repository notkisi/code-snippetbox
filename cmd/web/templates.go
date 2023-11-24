package main

import (
	"html/template"
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/notkisi/snippetbox/internal/models"
	"github.com/notkisi/snippetbox/ui"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templCache struct {
	templateCache map[string]*template.Template
}

func (t *templCache) Update() {
	// todo properly handle error
	var err error
	t.templateCache, err = newTemplateCache()
	if err != nil {
		log.Fatal(err)
		return
	}
}

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//strip full path
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:14")
}
