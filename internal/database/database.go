package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type dBstructure struct {
	Chirps     map[int]chirp     `json:"chirps"`
	Users      map[int]User      `json:"users"`
	Revocation map[string]Revoke `json:"revocation"`
}

type chirp struct {
	Id        int    `json:"id"`
	Body      string `json:"body"`
	Author_id int    `json:"author_id"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err

}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbstructure := dBstructure{
		Chirps:     map[int]chirp{},
		Users:      map[int]User{},
		Revocation: map[string]Revoke{},
	}
	return db.writeDB(dbstructure)
}

func (db *DB) writeDB(dbstructure dBstructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(dbstructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (dBstructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbstructure := dBstructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbstructure, err
	}
	/*if err != nil {
		return dbstructure, err
	}*/
	err = json.Unmarshal(data, &dbstructure)
	if err != nil {
		return dbstructure, err
	}

	return dbstructure, nil
}
func (db *DB) Createchirp(body string, userid int) (chirp, error) {
	dBstructure, err := db.loadDB()
	if err != nil {
		return chirp{}, err
	}
	id := len(dBstructure.Chirps) + 1

	chirp := chirp{
		Id:        id,
		Body:      body,
		Author_id: userid,
	}

	dBstructure.Chirps[id] = chirp

	err = db.writeDB(dBstructure)
	if err != nil {
		return chirp, err
	}
	return chirp, err

}

func (db *DB) Deletechirp(chirpID, userID int) error {
	dBstructure, err := db.loadDB()
	if err != nil {
		return err
	}

	Chirp, ok := dBstructure.Chirps[chirpID]

	if !ok {
		return errors.New("Chirp doesn't exist")
	}
	if Chirp.Author_id != userID {
		return errors.New("chirp doesnt belong to the author")
	}

	chirps := dBstructure.Chirps

	delete(chirps, chirpID)
	dBstructure.Chirps = chirps
	err = db.writeDB(dBstructure)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetChirps() ([]chirp, error) {
	dbstructure, err := db.loadDB()
	if err != nil {
		return []chirp{}, err
	}
	chirps := make([]chirp, 0, len(dbstructure.Chirps))

	for _, chirpe := range dbstructure.Chirps {
		chirps = append(chirps, chirpe)
	}

	return chirps, nil

}
