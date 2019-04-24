package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Env struct {
	Store Store
}

func jsonFprint(w io.Writer, a interface{}) error {
	json, err := json.Marshal(a)

	if err != nil {
		return err
	}

	fmt.Fprint(w, json)
	return nil
}

func NewRouter(env *Env) *httprouter.Router {
	router := httprouter.New()
	store := env.Store

	router.GET("/show", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		shows, err := store.GetAllShows()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error while reading item")
			return
		}

		jsonFprint(w, shows)
	})

	router.GET("/show/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, _ := strconv.ParseUint(ps.ByName("id"), 10, 64)
		show, err := store.GetShow(uint(id))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error while reading item")
			return
		}

		jsonFprint(w, show)
	})

	return router
}
