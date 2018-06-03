package main

import (
	"net/http"
	"fmt"
	"log"
)

type server struct {
	router *http.ServeMux
	port int
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleRoot())
}

func (s *server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	}
}

func main() {
	app := server{router: http.NewServeMux()}
	app.routes()
	log.Fatal(http.ListenAndServe(":8080", app.router))
}

