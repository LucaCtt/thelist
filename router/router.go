package router 

import (
	"net/http"

	"github.com/LucaCtt/thelist/data"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(store data.Store) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(contentTypeMiddleware)

	router.Mount("/show", showRouter(store))

	return router
}

func contentTypeMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
    next.ServeHTTP(w, r)
  })
}