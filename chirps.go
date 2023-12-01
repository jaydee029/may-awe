package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Chirp struct {
	Id        int    `json:"id"`
	Body      string `json:"body"`
	Author_id int    `json:"author_id"`
}

func (cfg *apiconfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirps couldn't be fetched")
		return
	}
	author_id := -1
	s := r.URL.Query().Get("author_id")
	if s != "" {
		author_id, err = strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "author id couldnt be parsed")
			return
		}
	}

	chirps_ := []Chirp{}

	for _, chirp := range chirps {
		if author_id != -1 && chirp.Author_id != author_id {
			continue
		}

		chirps_ = append(chirps_, Chirp{
			Id:        chirp.Id,
			Body:      chirp.Body,
			Author_id: chirp.Author_id,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJson(w, http.StatusOK, chirps_)
}

func (cfg *apiconfig) ChirpsbyId(w http.ResponseWriter, r *http.Request) {
	chirpidstr := chi.URLParam(r, "chirpId")
	chirpid, err := strconv.Atoi(chirpidstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse id")
		return
	}

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirps couldn't be fetched")
		return
	}

	for _, chirp := range chirps {
		if chirp.Id == chirpid {
			respondWithJson(w, http.StatusOK, chirp)
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "couldn't get chirp")
	return
}
