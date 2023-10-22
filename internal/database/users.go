package database

import (
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}

type Res struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
type res struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
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

func (db *DB) GetUser(email string, passwd string) (Res, error) {
	database, err := db.loadDB()
	if err != nil {
		return Res{}, err
	}

	//user, ok := database.Users[id]

	for id, user := range database.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(passwd))
			if err != nil {
				return Res{}, errors.New("Wrong password entered")
			}
			return Res{
				Id:    id,
				Email: user.Email,
				//Token: ss,
			}, nil
		}
	}

	return Res{}, os.ErrNotExist

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
