package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Env struct {
	Store Store
}

func newRouter(env *Env) *httprouter.Router {
	router := httprouter.New()
	store := env.Store

	router.GET("/item", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		items, err := store.GetAllItems()
		itemsJSON, err := json.Marshal(items)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error while reading item")
			return
		}

		fmt.Fprint(w, itemsJSON)
	})

	return router
}
