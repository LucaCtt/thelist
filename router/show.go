package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LucaCtt/thelist/data"
	"github.com/go-chi/chi"
)

func showRouter(store data.Store) http.Handler {
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

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var show data.Show
		err := json.NewDecoder(r.Body).Decode(&show)

		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		err = store.CreateShow(&show)

		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		err = store.DeleteShow(uint(id))
		if err != nil {
			if store.IsRecordNotFoundError(err) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	return router
}
