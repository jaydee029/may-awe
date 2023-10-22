package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	auth "github.com/jaydee029/Bark/internal"
	"github.com/jaydee029/Bark/internal/database"
)

type Input struct {
	Password           string `json:"password"`
	Email              string `json:"email"`
	Expires_in_seconds int    `json:"expires_in_seconds"`
}

type create struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
type User struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}
type res struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiconfig) createUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := create{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondWithJson(w, http.StatusCreated, res{
		Email: user.Email,
		ID:    user.Id,
	})
}

func (cfg *apiconfig) userLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Input{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUser(params.Email, params.Password)

	if params.Expires_in_seconds == 0 {
		params.Expires_in_seconds = 60 * 60 * 24 //default expiration
	}

	token, err := auth.Tokenize(user.Id, cfg.jwtsecret, time.Duration(params.Expires_in_seconds)*time.Second)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, database.Res{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
	})

}

func (cfg *apiconfig) updateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Idstr, err := auth.ValidateToken(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := strconv.Atoi(Idstr)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user Id couldn't be parsed")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := User{}
	err = decoder.Decode(&params)

	hashedPasswd, err := cfg.DB.Hashpassword(string(params.Password))
	params.Password = []byte(hashedPasswd)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	updateduser, err := cfg.DB.UpdateUser(userId, database.User(params))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, res{
		ID:    updateduser.ID,
		Email: updateduser.Email,
	})
}
