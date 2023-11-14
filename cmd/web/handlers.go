package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/notkisi/snippetbox/internal/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
	}

	a.render(w, http.StatusOK, "home.tmpl", &templateData{
		Snippets: snippets,
	})
}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}
		return
	}
	a.render(w, http.StatusOK, "view.tmpl", &templateData{
		Snippet: snippet,
	})
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
