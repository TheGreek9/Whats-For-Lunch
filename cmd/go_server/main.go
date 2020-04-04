package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"pkg/go_server"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", server.HomeHandler)
	r.HandleFunc("/groups", server.GroupsHandler)
	log.Println("serving localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
