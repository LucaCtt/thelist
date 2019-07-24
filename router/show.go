package router

import (
	"encoding/json"
	"net/http"

	"github.com/LucaCtt/thelist/data"
	"github.com/go-chi/chi"
)

func showRouter(store data.Store) http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		items, err := store.GetAllItems()
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(items)
	})

	router.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdParam(r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}

		show, err := store.GetItem(id)
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
		var show data.Item
		err := json.NewDecoder(r.Body).Decode(&show)
		if err != nil || !show.IsValid() {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		err = store.CreateItem(&show)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdParam(r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = store.DeleteItem(id)
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
