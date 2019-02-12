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
	// wsRouter := router.PathPrefix("/ws").Subrouter()

	for _, route := range route.ApiRoutes {
		apiRouter.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Path).
			Handler(route.Handler)
	}

	hub := newHub()
	go hub.run()

	// Upgrade handler websocket connection
	// wsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })

	// wsRouter.Methods("GET").
	// 	Path("/").
	// 	Handler(wsHandler)
	// router.SkipClean(true)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", router)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
