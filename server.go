package main

import (
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	fileStore := NewFileStore(itemsFile)
	env := &Env{Store: fileStore}
	router := newRouter(env)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	log.Printf("> %s\n", startMsg)
	log.Fatal(http.ListenAndServe(":8080", n))
}
