package router

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/LucaCtt/thelist/common"
)

// New returns a new router that handles the /show route
func New(store common.Store) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))
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

func getIDParam(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
