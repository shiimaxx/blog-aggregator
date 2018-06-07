package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const defaultListenPort = "8080"

type server struct {
	router *http.ServeMux
	port   string
	logger *log.Logger
	config config
}

type config struct {
	userID string
}

type entriesResponse struct {
	Entries []entry `json:"entries"`
}

func (s *server) routes() {
	s.router.HandleFunc("/api/v1/entries", s.handleEntries())
	s.router.HandleFunc("/", s.handleRoot())
}

func (s *server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}
}

func (s *server) handleEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := fetchQiitaEntries(s.config.userID)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var res entriesResponse
		res.Entries = entries
		if err := json.NewEncoder(w).Encode(res); err != nil {
			s.logger.Print(err.Error())
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
