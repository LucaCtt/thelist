package main

import (
	"log"
	"net/http"

	"github.com/LucaCtt/thelist/router"
	"github.com/LucaCtt/thelist/store"
)

func main() {
	dbStore, err := store.NewDbStore(&store.DbOptions{})
	defer dbStore.Close()

	if err != nil {
		log.Print(err)
		return
	}

	router := router.New(dbStore)

	log.Printf("> %s\n", startMsg)
	log.Print(http.ListenAndServe(":8080", router))
}
