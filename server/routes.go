package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Mount("/api", s.useStripSlashes(s.useRecoverer(s.useContentType("application/json", s.handleAPI()))))
}

func (s *Server) handleAPI() http.HandlerFunc {
	router := chi.NewRouter()

	router.Get("/", s.handleAPIGetAll())
	router.Get("/{name}", s.handleAPISearch())
	router.Post("/", s.handleAPIPost())
	router.Delete("/{id}", s.handleAPIDelete())

	return router.ServeHTTP
}
