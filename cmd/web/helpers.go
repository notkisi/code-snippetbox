package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/notkisi/snippetbox/internal/validator"
)

func (a *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := a.templateCache.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		a.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		a.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (a *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       a.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (a *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = a.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
	}
	return nil
}

func (a *application) validateSnippetForm(form *snippetCreateForm) error {
	form.CheckField(validator.NotBlank(form.Title), "title", "this field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "this field cant have more than 100 chars")
	form.CheckField(validator.NotBlank(form.Content), "content", "this field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "this field must be one of 1 7 365")

	if !form.Valid() {
		return errors.New("Form is not valid")
	}
	return nil
}

func (a *application) validateSignupForm(form *userSignupForm) error {
	form.CheckField(validator.NotBlank(form.Name), "name", "this field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "this field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "this field cannot be blank")
	form.CheckField(validator.MaxBytes(form.Password), "password", "password too long")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "must be min 8 chars long")
	form.CheckField(validator.Matches(form.Email, validator.EmailRE), "email", "this field must be a valid email address")

	if !form.Valid() {
		return errors.New("Form is not valid")
	}
	return nil
}
