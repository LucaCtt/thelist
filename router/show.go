package router

import (
	"encoding/json"
	"net/http"

	"github.com/LucaCtt/thelist/common/store"
	"github.com/go-chi/chi"
)

func showRouter(s store.Store) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		items, err := s.All()
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

		show, err := s.Get(id)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(show)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var show store.Item
		err := json.NewDecoder(r.Body).Decode(&show)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		err = s.Create(&show)
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

		err = s.Delete(id)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	return router
}
