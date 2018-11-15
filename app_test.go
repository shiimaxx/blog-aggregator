package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("NewRequest failed: ", err.Error())
	}

	rec := httptest.NewRecorder()

	s := server{}
	handler := s.handleRoot()

	handler.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusMovedPermanently; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	if got, want := rec.HeaderMap.Get("Location"), "/api/v1/entries"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestHandleEntries(t *testing.T) {
	// TODO: Write test
}
