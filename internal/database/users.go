package database

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}

type res struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	//Token string `json:"token"`
}

func (db *DB) CreateUser(email string, passwd string) (res, error) {
	database, err := db.loadDB()
	if err != nil {
		return res{}, err
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
		return res{}, err
	}
	return res{
		Id:    id,
		Email: email,
	}, nil
}

func (db *DB) GetUser(email string, passwd string) (res, error) {
	godotenv.Load()
	database, err := db.loadDB()
	if err != nil {
		return res{}, err
	}

	//user, ok := database.Users[id]

	for id, user := range database.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(passwd))
			if err != nil {
				return res{}, errors.New("Wrong password entered")
			}
			/*claims := &jwt.RegisteredClaims{
				Issuer:    "Bark",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now()),
				Subject:   strconv.Itoa(id),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			ss, err := token.SignedString(os.Getenv("JWT_SECRET"))
			if err != nil {
				return res{}, err
			}*/

			return res{
				Id:    id,
				Email: user.Email,
				//Token: ss,
			}, nil
		}
	}

	return res{}, os.ErrNotExist

}
