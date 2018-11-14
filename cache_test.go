package main

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/shiimaxx/blog-aggregator/structs"
)

var cache = &memStorage{
	items: make(map[string]item),
	mu:    &sync.RWMutex{},
}
var now = time.Now()

func TestMemStorage_Get(t *testing.T) {
	cases := []struct {
		name string
		key  string
		want []structs.Entry
	}{
		{
			name: "key-1",
			key:  "1",
			want: []structs.Entry{
				{Title: "a", URL: "https://example.com/a", CreatedAt: now},
				{Title: "b", URL: "https://example.com/b", CreatedAt: now.Add(1 * time.Hour)},
				{Title: "c", URL: "https://example.com/c", CreatedAt: now.Add(2 * time.Hour)},
			},
		},
		{
			name: "key-2",
			key:  "2",
			want: []structs.Entry{
				{Title: "d", URL: "https://example.com/d", CreatedAt: now.Add(3 * time.Hour)},
				{Title: "e", URL: "https://example.com/e", CreatedAt: now.Add(4 * time.Hour)},
				{Title: "f", URL: "https://example.com/f", CreatedAt: now.Add(5 * time.Hour)},
			},
		},
		{
			name: "key-3",
			key:  "3",
			want: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := cache.Get(tc.key)
			if !reflect.DeepEqual(c, tc.want) {
				t.Fatalf("got %v; want %v", c, tc.want)
			}
		})
	}
}

func TestMemStorage_Set(t *testing.T) {
	c := cache.Get("4")
	if c != nil {
		t.Fatalf("got %v; want nil", c)
	}

	cache.Set("4", []structs.Entry{
		{Title: "g", URL: "https://example.com/g", CreatedAt: now.Add(9 * time.Hour)},
		{Title: "h", URL: "https://example.com/h", CreatedAt: now.Add(10 * time.Hour)},
		{Title: "e", URL: "https://example.com/e", CreatedAt: now.Add(11 * time.Hour)},
	}, 60*time.Second)

	cc := cache.Get("4")
	if cc == nil {
		t.Fatal("got nil; want not nil")
	}

}

func TestMain(m *testing.M) {
	cache.Set("1", []structs.Entry{
		{Title: "a", URL: "https://example.com/a", CreatedAt: now},
		{Title: "b", URL: "https://example.com/b", CreatedAt: now.Add(1 * time.Hour)},
		{Title: "c", URL: "https://example.com/c", CreatedAt: now.Add(2 * time.Hour)},
	}, 60*time.Second)
	cache.Set("2", []structs.Entry{
		{Title: "d", URL: "https://example.com/d", CreatedAt: now.Add(3 * time.Hour)},
		{Title: "e", URL: "https://example.com/e", CreatedAt: now.Add(4 * time.Hour)},
		{Title: "f", URL: "https://example.com/f", CreatedAt: now.Add(5 * time.Hour)},
	}, 60*time.Second)
	cache.Set("3", []structs.Entry{
		{Title: "d", URL: "https://example.com/d", CreatedAt: now.Add(6 * time.Hour)},
		{Title: "e", URL: "https://example.com/e", CreatedAt: now.Add(7 * time.Hour)},
		{Title: "f", URL: "https://example.com/f", CreatedAt: now.Add(8 * time.Hour)},
	}, 0*time.Second)

	os.Exit(m.Run())
}
