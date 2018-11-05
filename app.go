package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/shiimaxx/blog-aggregator/hatenablog"
	"github.com/shiimaxx/blog-aggregator/qiita"
	"github.com/shiimaxx/blog-aggregator/structs"
)

const defaultListenPort = "8080"

type server struct {
	router *http.ServeMux
	port   string
	logger *log.Logger
	config config
}

type config struct {
	userID       string
	hatenaBlogID string
	hatenaAPIKey string
}

type entriesResponse struct {
	Entries []structs.Entry `json:"entries"`
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
		q, err := qiita.FetchEntries(s.config.userID)
		if err != nil {
			s.logger.Printf("[ERROR] %s %s %s %s", r.Method, r.URL.Host, r.URL.Path, err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		h, err := hatenablog.FetchEntries(s.config.userID, s.config.hatenaBlogID, s.config.hatenaAPIKey)
		if err != nil {
			s.logger.Printf("[ERROR] %s %s %s %s", r.Method, r.URL.Host, r.URL.Path, err.Error())
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		entries := append(q, h...)

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

	app := server{
		router: http.NewServeMux(),
		port:   port,
		logger: log.New(os.Stdout, "", log.Lshortfile),
		config: config{userID: userID},
	}
	app.routes()
	log.Fatal(http.ListenAndServe(":"+app.port, app.router))
}
