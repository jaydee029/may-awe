package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiconfig) postChirps(w http.ResponseWriter, r *http.Request) {
	type res struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}

	type body struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parameters couldn't be decoded")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	content := profane(params.Body)

	chirp, err := cfg.DB.Createchirp(content)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondWithJson(w, http.StatusCreated, res{
		ID:   chirp.Id,
		Body: chirp.Body,
	})

}

func profane(content string) string {
	contentslice := strings.Split(content, " ")

	for i, word := range contentslice {
		wordl := strings.ToLower(word)
		if wordl == "kerfuffle" || wordl == "sharbert" || wordl == "fornax" {
			contentslice[i] = "****"
		}
	}

	return strings.Join(contentslice, " ")
}

func respondWithError(w http.ResponseWriter, code int, res string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", res)
	}
	type errresponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errresponse{
		Error: res,
	})
}

func respondWithJson(w http.ResponseWriter, code int, res interface{}) {
	w.Header().Set("content-type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
