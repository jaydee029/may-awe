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
	Password      []byte `json:"password"`
	Email         string `json:"email"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
}
type res struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
}
type res_login struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
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

	Token, err := auth.Tokenize(user.ID, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Refresh_token, err := auth.RefreshToken(user.ID, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, res_login{
		ID:            user.ID,
		Email:         user.Email,
		Is_chirpy_red: user.Is_chirpy_red,
		Token:         Token,
		Refresh_token: Refresh_token,
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
		ID:            updateduser.ID,
		Email:         updateduser.Email,
		Is_chirpy_red: updateduser.Is_chirpy_red,
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

func (cfg *apiconfig) is_red(w http.ResponseWriter, r *http.Request) {
	type user_struct struct {
		User_id int `json:"user_id"`
	}
	type body struct {
		Event string      `json:"event"`
		Data  user_struct `json:"data"`
	}

	key, err := auth.VerifyAPIkey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if key != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API Key")
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	if params.Event == "user.upgraded" {
		user_res, err := cfg.DB.Is_red(params.Data.User_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJson(w, http.StatusOK, res{
			Email:         user_res.Email,
			Is_chirpy_red: user_res.Is_chirpy_red,
			ID:            params.Data.User_id,
		})
	}

	respondWithJson(w, http.StatusOK, "http request accepted in the webhook")
}
