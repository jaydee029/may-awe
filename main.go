package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jaydee029/Bark/internal/database"
)

type apiconfig struct {
	fileservercounts int
	DB               *database.DB
}

func main() {

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	apicfg := apiconfig{
		fileservercounts: 0,
		DB:               db,
	}

	port := "8100"

	r := chi.NewRouter()
	s := chi.NewRouter()
	t := chi.NewRouter()

	fileconfig := apicfg.reqcounts(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fileconfig)
	r.Handle("/app/*", fileconfig)

	s.Get("/healthz", apireadiness)
	s.Post("/chirps", apicfg.postChirps)
	s.Get("/chirps", apicfg.getChirps)
	s.Get("/chirps/{chirpId}", apicfg.ChirpsbyId)
	t.Get("/metrics", apicfg.metrics)

	r.Mount("/api", s)
	r.Mount("/admin", t)
	sermux := corsmiddleware(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
