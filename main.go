package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"whatdash/route"
	"whatdash/wa"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatalln("Please define a $PORT to running on")
	}

	router := mux.NewRouter().StrictSlash(false)

	apiRouter := router.PathPrefix("/api").Subrouter()

	dbSess, _ := wa.ConnectionOpen()

	storage := wa.BucketSession{
		Items:      make(map[string]wa.ConnWrapper),
		MgoSession: dbSess,
	}
	storage.Sync()
	routes := route.InitRoutes(&storage)
	for _, route := range routes {
		apiRouter.
			Methods(route.Method).
			Name(route.Name).
			Path(route.Path).
			Handler(route.Handler)
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	// set timeout request to 3 mins.
	withTimeout := http.TimeoutHandler(router, time.Second*180, "Timeout....")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	http.Handle("/", c.Handler(withTimeout))
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
