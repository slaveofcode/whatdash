package main

import (
	"log"
	"net/http"
	"os"
	"whatdash/route"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatalln("Please define a $PORT to running on")
	}

	router := mux.NewRouter().StrictSlash(false)

	apiRouter := router.PathPrefix("/api").Subrouter()

	for _, route := range route.ApiRoutes {
		apiRouter.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Path).
			Handler(route.Handler)
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", router)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
