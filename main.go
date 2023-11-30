package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jaydee029/Bark/internal/database"
	"github.com/joho/godotenv"
)

type apiconfig struct {
	fileservercounts int
	DB               *database.DB
	jwtsecret        string
}

func main() {
	godotenv.Load(".env")
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT secret key not set")
	}
	apicfg := apiconfig{
		fileservercounts: 0,
		DB:               db,
		jwtsecret:        jwt_secret,
	}

	port := os.Getenv("PORT")

	fmt.Println(port)
	fmt.Println(apicfg.jwtsecret)

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
	s.Post("/users", apicfg.createUser)
	s.Post("/login", apicfg.userLogin)
	s.Post("/refresh", apicfg.verifyRefresh)
	s.Post("/revoke", apicfg.revokeToken)
	s.Put("/users", apicfg.updateUser)
	s.Delete("/chirps/{chirpId}", apicfg.DeleteChirps)
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
