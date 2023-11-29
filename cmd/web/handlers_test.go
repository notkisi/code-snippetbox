package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/notkisi/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	httpRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(responseRecorder, httpRequest)
	result := responseRecorder.Result()

	assert.Equal(t, result.StatusCode, http.StatusOK)

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	trimmedBody := bytes.TrimSpace(body)
	assert.Equal(t, string(trimmedBody), "OK")
}

func TestPingE2E(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name         string
		urlPath      string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid ID",
			urlPath:      "/snippet/view/1",
			expectedCode: http.StatusOK,
			expectedBody: "An old silent pond...",
		},
		{
			name:         "Non-existent ID",
			urlPath:      "/snippet/view/2",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Negative ID",
			urlPath:      "/snippet/view/-1",
			expectedCode: http.StatusNotFound,
		},
		{name: "Decimal ID",
			urlPath:      "/snippet/view/1.23",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "String ID",
			urlPath:      "/snippet/view/foo",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Empty ID",
			urlPath:      "/snippet/view/",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code, _, body := ts.get(t, tc.urlPath)

			assert.Equal(t, code, tc.expectedCode)

			if tc.expectedBody != "" {
				assert.StringContains(t, body, tc.expectedBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCsrfToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)
	tests := []struct {
		name            string
		username        string
		email           string
		password        string
		csrfToken       string
		expectedCode    int
		expectedFormTag string
	}{
		{
			name:         "Valid submission",
			username:     validName,
			email:        validEmail,
			password:     validPassword,
			csrfToken:    csrfToken,
			expectedCode: http.StatusSeeOther,
		},
		{
			name:         "inValid submission, wrong token",
			username:     validName,
			email:        validEmail,
			csrfToken:    "badToken",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:            "Empty name",
			username:        "",
			password:        validPassword,
			email:           validEmail,
			csrfToken:       csrfToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Empty email",
			username:        validName,
			email:           "",
			password:        validPassword,
			csrfToken:       csrfToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Short password",
			username:        validName,
			email:           validEmail,
			password:        "shorty",
			csrfToken:       csrfToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
		{
			name:            "Duplicate email",
			username:        validName,
			email:           "dupe@example.com", // take a look at mock
			password:        validPassword,
			csrfToken:       csrfToken,
			expectedCode:    http.StatusUnprocessableEntity,
			expectedFormTag: formTag,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tc.username)
			form.Add("email", tc.email)
			form.Add("password", tc.password)
			form.Add("csrf_token", tc.csrfToken)

			code, _, body := ts.post(t, "/user/signup", form)
			assert.Equal(t, code, tc.expectedCode)
			if tc.expectedFormTag != "" {
				assert.StringContains(t, body, tc.expectedFormTag)
			}
		})
	}

}
