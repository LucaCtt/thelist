package main

import (
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/store"
	"github.com/urfave/negroni"
)

func main() {
	dbStore, err := store.NewDbStore(&store.DbOptions{})
	defer dbStore.Close()

	if err != nil {
		log.Print(err)
		return
	}

	router := NewRouter(dbStore)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	log.Printf("> %s\n", startMsg)
	log.Print(http.ListenAndServe(":8080", n))
}
