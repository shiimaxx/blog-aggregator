package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
)

type server struct {
	db     *leveldb.DB
	router *http.ServeMux
	port   int
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
	db, err := leveldb.OpenFile("./db", nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	app := server{db: db, router: http.NewServeMux()}
	app.routes()
	log.Fatal(http.ListenAndServe(":8080", app.router))
}
