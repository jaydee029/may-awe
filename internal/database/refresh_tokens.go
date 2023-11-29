package database

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Revoke struct {
	Token      string    `json:"token"`
	Revoked_at time.Time `json:"revoked_at"`
}
type resid struct {
	ID int `json:"id"`
}

func (db *DB) Verifyrevocation(refresh_token string) (bool, error) {

	dBstructure, err := db.loadDB()

	if err != nil {
		return false, err
	}

	_, ok := dBstructure.Revocation[refresh_token]

	if !ok {
		return false, nil
	}
	/*
		if is_revoked.Revoked_at.IsZero() {
			return true, nil
		}*/

	return true, nil
}

func (db *DB) VerifyRefresh(tokenstring, tokenSecret string) (bool, error) {
	type customClaims struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(tokenstring, &customClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return false, errors.New(err.Error()) //"jwt couldn't be parsed"
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return false, errors.New("issuer couldn't be extracted")
	}

	if issuer == "chirpy-refresh" {
		return true, nil
	}
	return false, nil
}

func (db *DB) VerifyRefreshSignature(tokenstring, tokenSecret string) (string, error) {
	type customClaims struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(tokenstring, &customClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })

	if err != nil {
		return "", errors.New(err.Error()) //"jwt couldn't be parsed"
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return "", errors.New("issuer couldn't be extracted")
	}

	if issuer != "chirpy-refresh" {
		return "", errors.New("Not a refresh Token")
	}

	userId, err := token.Claims.GetSubject()

	if err != nil {
		return "", errors.New("User id couldn't be extracted")
	}

	return userId, nil
}

func (db *DB) RevokeToken(tokenstring string) error {
	dbstructure, err := db.loadDB()

	if err != nil {
		return err
	}

	/*User, ok := dbstructure.Revocation[tokenstring]

	if !ok {
		return errors.New("Refresh Token not found")
	}

	User.Revoked_at = time.Now().UTC()
	*/
	dbstructure.Revocation[tokenstring] = Revoke{
		Token:      tokenstring,
		Revoked_at: time.Now().UTC(),
	}

	err = db.writeDB(dbstructure)

	if err != nil {
		return err
	}

	return nil

}
