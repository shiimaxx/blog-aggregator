package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/shiimaxx/blog-aggregator/structs"
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

func TestHandleEntries_CacheHit(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/entries", nil)
	if err != nil {
		t.Fatal("NewRequest failed: ", err.Error())
	}

	rec := httptest.NewRecorder()

	s := server{
		logger: log.New(os.Stdout, "", log.Lshortfile),
		config: config{
			userID: "testuser",
		},
		cache: &memStorage{
			items: make(map[string]item),
			mu:    &sync.RWMutex{},
		},
	}
	s.cache.Set("qiita", []structs.Entry{
		{Title: "a", URL: "https://example.com/a", CreatedAt: now},
		{Title: "b", URL: "https://example.com/b", CreatedAt: now.Add(1 * time.Hour)},
		{Title: "c", URL: "https://example.com/c", CreatedAt: now.Add(2 * time.Hour)},
	}, 60*time.Second)
	s.cache.Set("hatenablog", []structs.Entry{
		{Title: "d", URL: "https://example.com/d", CreatedAt: now.Add(3 * time.Hour)},
		{Title: "e", URL: "https://example.com/e", CreatedAt: now.Add(4 * time.Hour)},
		{Title: "f", URL: "https://example.com/f", CreatedAt: now.Add(5 * time.Hour)},
	}, 60*time.Second)

	handler := s.handleEntries()

	handler.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
