package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/common/client"
	"github.com/LucaCtt/thelist/common/store"
	"github.com/LucaCtt/thelist/errors"
	"github.com/go-chi/chi"
)

// Server is an http server which serves a REST API for shows.
type Server struct {
	store  store.Store
	client client.Client
	router chi.Router
}

// New creates a new Server with the given dependencies.
func New(s store.Store, c client.Client) *Server {
	server := &Server{
		store:  s,
		client: c,
		router: chi.NewRouter(),
	}
	server.routes()
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}
	}
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, out interface{}) error {
	return json.NewDecoder(r.Body).Decode(out)
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, e error) {
	code := errors.Code(e)
	http.Error(w, http.StatusText(int(code)), int(code))
	log.Print(e)
}
