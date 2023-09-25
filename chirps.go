package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiconfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirps couldn't be fetched")
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJson(w, http.StatusOK, chirps)
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
}
