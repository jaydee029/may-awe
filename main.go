package main

import (
	"fmt"
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

func (cfg *apiconfig) reqcounts(app http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileservercounts++
		app.ServeHTTP(w, r)
	})

}
func (cfg *apiconfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileservercounts)))
}

type apiconfig struct {
	fileservercounts int
}

func main() {
	apicfg := apiconfig{
		fileservercounts: 0,
	}

	port := "8100"
	newmux := http.NewServeMux()
	newmux.Handle("/app/", apicfg.reqcounts(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	newmux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})
	newmux.HandleFunc("/metrics", apicfg.metrics)
	sermux := corsmiddleware(newmux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: sermux,
	}

	log.Printf("The server is live on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
