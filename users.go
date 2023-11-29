package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	auth "github.com/jaydee029/Bark/internal"
	"github.com/jaydee029/Bark/internal/database"
)

type Input struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Token struct {
	Token string `json:"token"`
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
	params := Input{}
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

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := auth.Tokenize(user.Id, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refresh_token, err := auth.RefreshToken(user.Id, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, database.Res{
		Id:            user.Id,
		Email:         user.Email,
		Token:         token,
		Refresh_token: refresh_token,
	})

}

func (cfg *apiconfig) updateUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	is_refresh, err := cfg.DB.VerifyRefresh(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if is_refresh == true {
		respondWithError(w, http.StatusUnauthorized, "Header contains refresh token")
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

func (cfg *apiconfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)

	if err != io.EOF {
		respondWithError(w, http.StatusUnauthorized, "Body is provided")
		return
	}

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.DB.RevokeToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	respondWithJson(w, http.StatusOK, res{})
}

func (cfg *apiconfig) verifyRefresh(w http.ResponseWriter, r *http.Request) {

	/*decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)

	if err != io.EOF {
		respondWithError(w, http.StatusUnauthorized, "Body is provided")
		return
	}*/

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Idstr, err := cfg.DB.VerifyRefreshSignature(token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userid, err := strconv.Atoi(Idstr)

	is_revoked, err := cfg.DB.Verifyrevocation(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if is_revoked == true {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token revoked")
		return
	}

	auth_token, err := auth.Tokenize(userid, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, Token{
		Token: auth_token,
	})
}
