package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	server "pkg/go_server"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

var (
	allowedOrigins     = handlers.AllowedOrigins([]string{"*"})
	allowedHeaders     = handlers.AllowedHeaders([]string{"Accept", "DNT", "Content-Type", "Referer", "User Agent", "Sec-Fetch-Dest"})
	allowedCredentials = handlers.AllowCredentials()
	allowedMethods     = handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	router := mux.NewRouter()
	db, err := bolt.Open("/Users/Spyro/Developer/go/src/db/wfl.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	rh := &server.RouteHandler{Db: server.Database{Db: db}}
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	router.HandleFunc("/v1/group/query/{groupId}", func(w http.ResponseWriter, r *http.Request) {
		b, err := rh.HandlerQueryDb(w, r, "groupId", server.BktGroup)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(b)
	})

	router.HandleFunc("/v1/group/create", func(w http.ResponseWriter, r *http.Request) {
		ID, err := rh.HandlerCreateGroup(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(ID)
	}).Methods("POST")

	router.HandleFunc("/v1/group/delete", func(w http.ResponseWriter, r *http.Request) {
		if err := rh.HandlerDeleteDb(w, r, "groupId", server.BktGroup); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	spa := spaHandler{staticPath: "/Users/Spyro/Developer/go/src/web/wfl_web/public", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods, allowedCredentials)(router),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Serving on 127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
