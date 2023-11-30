package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	auth "github.com/jaydee029/Bark/internal"
)

type Res struct {
	Author_id int    `json:"author_id"`
	Body      string `json:"body"`
	ID        int    `json:"id"`
}

type body struct {
	Body string `json:"body"`
}

func (cfg *apiconfig) postChirps(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	authorid, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	authorid_num, err := strconv.Atoi(authorid)

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parameters couldn't be decoded")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	content := profane(params.Body)

	chirp, err := cfg.DB.Createchirp(content, authorid_num)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	respondWithJson(w, http.StatusCreated, Res{
		Author_id: chirp.Author_id,
		Body:      chirp.Body,
		ID:        chirp.Id,
	})

}

func (cfg *apiconfig) DeleteChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	authorid, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	authorid_num, err := strconv.Atoi(authorid)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse author id")
		return
	}

	chirpidstr := chi.URLParam(r, "chirpId")
	chirpid, err := strconv.Atoi(chirpidstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse chirp id")
		return
	}

	err = cfg.DB.Deletechirp(chirpid, authorid_num)

	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, "Chirp deleted")
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
