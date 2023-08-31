package main

import (
	"log"
	"net/http"
)

func corsmiddleware(app http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}
		app.ServeHTTP(w, r)
	})
}

func main() {
	port := "8100"
	newmux := http.NewServeMux()
	sermux := corsmiddleware(newmux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
