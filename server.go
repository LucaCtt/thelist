package main

import (
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	dbStore, err := NewDbStore(&DbOptions{})
	defer dbStore.Close()

	if err != nil {
		log.Print(err)
		return
	}

	env := &Env{Store: dbStore}
	router := NewRouter(env)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	log.Printf("> %s\n", startMsg)
	log.Print(http.ListenAndServe(":8080", n))
}
