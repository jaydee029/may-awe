package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func chirpslength(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}

	type response struct {
		Cleaned_body string `json:"cleaned_body"`
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

	content := params.Body
	contentslice := strings.Split(content, " ")

	for i, word := range contentslice {
		wordl := strings.ToLower(word)
		if wordl == "kerfuffle" || wordl == "sharbert" || wordl == "fornax" {
			contentslice[i] = "****"
		}
	}

	contentreturn := strings.Join(contentslice, " ")

	respondWithJson(w, http.StatusOK, response{
		Cleaned_body: contentreturn,
	})
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
