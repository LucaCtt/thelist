package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LucaCtt/thelist/store"
	"github.com/go-chi/chi"
)

func showRouter(store store.Store) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		shows, err := store.GetAllShows()
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(shows)
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		show, err := store.GetShow(uint(id))
		if err != nil {
			if store.IsRecordNotFoundError(err) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			panic(err)
		}

		json.NewEncoder(w).Encode(show)
	})

	return router
}
