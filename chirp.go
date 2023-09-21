package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validate_chirp(w http.ResponseWriter, r *http.Request) {
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

	content := params.Body
	contentslice := strings.Split(content, " ")

	for i, word := range contentslice {
		wordl := strings.ToLower(word)
		if wordl == "kerfuffle" || wordl == "sharbert" || wordl == "fornax" {
			contentslice[i] = "****"
		}
	}

	

	
}
