package server

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello groups")
}
