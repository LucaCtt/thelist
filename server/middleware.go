package server

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) useStripSlashes(h http.HandlerFunc) http.HandlerFunc {
	return middleware.StripSlashes(h).ServeHTTP
}

func (s *Server) useRecoverer(h http.HandlerFunc) http.HandlerFunc {
	return middleware.Recoverer(h).ServeHTTP
}

func (s *Server) useAllowContent(t string, h http.HandlerFunc) http.HandlerFunc {
	return middleware.AllowContentType(t)(h).ServeHTTP
}

func (s *Server) useContentType(t string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", t)
		h.ServeHTTP(w, r)
	}
}
