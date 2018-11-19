package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/shiimaxx/blog-aggregator/blogservice"
	"github.com/shiimaxx/blog-aggregator/blogservice/hatenablog"
	"github.com/shiimaxx/blog-aggregator/blogservice/qiita"
	"github.com/shiimaxx/blog-aggregator/structs"
)

const defaultListenPort = "8080"
const defaultCacheExpiration = 60 * time.Second

type server struct {
	router      *http.ServeMux
	port        string
	logger      *log.Logger
	config      config
	cache       storage
	blogService *blogservice.BlogService
}

type config struct {
	userID       string
	hatenaID     string
	hatenaBlogID string
	hatenaAPIKey string
}

type entriesResponse struct {
	Entries []structs.Entry `json:"entries"`
}

func (s *server) blogservices() {
	s.blogService.Add(func() ([]structs.Entry, error) {
		return qiita.FetchEntries(context.TODO(), s.config.userID)
	})
	if s.config.hatenaID != "" {
		s.blogService.Add(func() ([]structs.Entry, error) {
			return hatenablog.FetchEntries(context.TODO(), s.config.hatenaID, s.config.hatenaBlogID, s.config.hatenaAPIKey)
		})
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/api/v1/entries", s.handleEntries())
	s.router.HandleFunc("/", s.handleRoot())
}

func (s *server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/entries", http.StatusMovedPermanently)
	}
}

func (s *server) handleEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey := GenerateCacheKey(r.URL.String(), "")
		var entries []structs.Entry
		if cache := s.cache.Get(cacheKey); cache != nil {
			entries = cache
		} else {
			s.logger.Printf("[INFO] %s %s %s %s", r.Method, r.URL.Host, r.URL.Path, "cache miss")
			e, err := s.blogService.Fetch()
			if err != nil {
				return
			}
			s.cache.Set(cacheKey, e, defaultCacheExpiration)
			entries = e
		}

		sort.Slice(entries, func(j, i int) bool {
			return entries[i].CreatedAt.Before(entries[j].CreatedAt)
		})

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var res entriesResponse
		res.Entries = entries
		if err := json.NewEncoder(w).Encode(res); err != nil {
			s.logger.Printf("[ERROR] %s %s %s %s", r.Method, r.URL.Host, r.URL.Path, err.Error())
		}
	}
}

func main() {
	var (
		port   string
		userID string
	)
	if port = os.Getenv("LISTEN_PORT"); port == "" {
		port = defaultListenPort
	}
	if userID = os.Getenv("USER_ID"); userID == "" {
		log.Fatal("USER_ID is required but missing")
	}
	hatenaID := os.Getenv("HATENA_ID")
	hatenaBlogID := os.Getenv("HATENA_BLOG_ID")
	hatenaBlogAPIKey := os.Getenv("HATENA_BLOG_API_KEY")

	app := server{
		router: http.NewServeMux(),
		port:   port,
		logger: log.New(os.Stdout, "", log.Lshortfile),
		config: config{
			userID:       userID,
			hatenaID:     hatenaID,
			hatenaBlogID: hatenaBlogID,
			hatenaAPIKey: hatenaBlogAPIKey,
		},
		cache: &memStorage{
			items: make(map[string]item),
			mu:    &sync.RWMutex{},
		},
		blogService: &blogservice.BlogService{
			FetchFunc: []func() ([]structs.Entry, error){},
		},
	}
	app.blogservices()
	app.routes()
	log.Fatal(http.ListenAndServe(":"+app.port, app.router))
}
