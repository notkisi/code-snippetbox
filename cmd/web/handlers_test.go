package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
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
