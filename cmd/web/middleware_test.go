package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/notkisi/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	// init dummy http.ResponseWriter and http.Request objects
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedBody := "whatever"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedBody))
	})

	secureHeaders(next).ServeHTTP(w, r)

	rs := w.Result()

	tests := []struct {
		name       string
		expected   string
		headerName string
	}{
		{
			name:       "CSP",
			expected:   "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
			headerName: "Content-Security-Policy",
		},
		{
			name:       "Referrer",
			expected:   "origin-when-cross-origin",
			headerName: "Referrer-Policy",
		},
		{
			name:       "XContentTypeOpts",
			expected:   "nosniff",
			headerName: "X-Content-Type-Options",
		},
		{
			name:       "XFrameOpts",
			expected:   "deny",
			headerName: "X-Frame-Options",
		},
		{
			name:       "XSS",
			expected:   "0",
			headerName: "X-XSS-Protection",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, rs.Header.Get(testCase.headerName), testCase.expected)
		})
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(body), expectedBody)
}
