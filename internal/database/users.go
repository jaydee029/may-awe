package database

import (
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password      []byte `json:"password"`
	Email         string `json:"email"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
}

type Res struct {
	Id            int    `json:"id"`
	Email         string `json:"email"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}
type res struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email string, passwd string) (Res, error) {
	database, err := db.loadDB()
	if err != nil {
		return Res{}, err
	}

	id := len(database.Users) + 1

	encrypted, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)

	user := User{
		Password: encrypted,
		Email:    email,
	}
	database.Users[id] = user
	err = db.writeDB(database)
	if err != nil {
		return Res{}, err
	}
	return Res{
		Id:    id,
		Email: email,
	}, nil
}

func (db *DB) GetUser(email string, passwd string) (res, error) {
	database, err := db.loadDB()
	if err != nil {
		return res{}, errors.New("Couldn't load database")
	}

	//user, ok := database.Users[id]

	for id, user := range database.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(passwd))
			if err != nil {
				return res{}, errors.New("Wrong password entered")
			}
			return res{
				ID:            id,
				Email:         user.Email,
				Is_chirpy_red: user.Is_chirpy_red,
			}, nil
		}
	}

	return res{}, os.ErrNotExist

}

func (db *DB) Hashpassword(passwd string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Couldn't Hash the password")
	}

	return string(encrypted), nil
}

func (db *DB) UpdateUser(userid int, userInput User) (res, error) {

	users, err := db.loadDB()

	if err != nil {
		return res{}, errors.New("Couldn't load the database")
	}
	user, ok := users.Users[userid]

	if !ok {
		return res{}, os.ErrNotExist
	}

	user.Email = userInput.Email
	user.Password = userInput.Password
	users.Users[userid] = user

	err = db.writeDB(users)
	if err != nil {
		return res{}, errors.New("Couldn't write into the database")
	}

	response := res{
		ID:    userid,
		Email: userInput.Email,
	}
	return response, nil

}

func (db *DB) Is_red(userid int) (User, error) {
	dBstructure, err := db.loadDB()

	if err != nil {
		return User{}, errors.New("Couldn't load the database")
	}

	user, ok := dBstructure.Users[userid]

	if !ok {
		return User{}, errors.New("User not found")
	}

	user.Is_chirpy_red = true

	dBstructure.Users[userid] = user

	err = db.writeDB(dBstructure)
	if err != nil {
		return User{}, errors.New("Couldn't write into the database")
	}

	return user, nil
}
