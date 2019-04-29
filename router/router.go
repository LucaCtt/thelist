package router 

import (
	"net/http"

	"github.com/LucaCtt/thelist/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(store store.Store) http.Handler {
	router := chi.NewRouter()

	router.Mount("/show", showRouter(store))

	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	return router
}
