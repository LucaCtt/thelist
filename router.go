package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LucaCtt/thelist/store"
	"github.com/julienschmidt/httprouter"
)

type Env struct {
	Store store.Store
}

func NewRouter(store store.Store) http.Handler {
	router := httprouter.New()

	router.GET("/show", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		shows, err := store.GetAllShows()
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(shows)
	})

	router.GET("/show/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
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
