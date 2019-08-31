package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lucactt/thelist/data"
)

func showRouter(store data.Store) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		items, err := store.All()
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(items)
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDParam(r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		show, err := store.Get(id)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(show)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var show data.Item
		err := json.NewDecoder(r.Body).Decode(&show)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		err = store.Create(&show)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDParam(r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = store.Delete(id)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	return router
}
