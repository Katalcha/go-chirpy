package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"syscall"
)

var ErrNotExist = errors.New("resource does not exist")

// represents the DB itself
type DB struct {
	path string
	mu   *sync.RWMutex
}

// represents the contents of DB as map of Chirps and map of Users
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

/*
	DB STRUCT AS SAVED FILE ON DISK:
	FILE PATH: "path/to/file/database.json"

	INSIDE JSON:
	{	<-- DBStructure struct
		"chirps": {	<-- DBStructure.Chirps map[int]string
			"1": { id: 1, body: "blabla" },		<-- Chirp struct
			"2": { id: 2, body: "blublub" },	<-- Chirp struct
		},
		"users" : { <-- DBStructure.Users map[int]string
			"1": { id: 1, email: "blabla@blub.com", password: <hash> },
			"2": { id: 2, email: "blubblub@bla.com", password: <hash> },
		}
	}
*/

// NEW DB FOR IN-MEMORY ON SERVER START
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// Ensures if a JSON-DB is present or not by reading DB.path.
// If err is nil, a JSON-DB was found, otherwise calls DB.createDB()
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

// Creates a new JSON-DB by creating the DBStructure struct
// and handing this struct to the writeDB() method
func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}

// Writes JSON-DB content to provided DBStructure by
// handling mutual exclusions, marshalling content to JSON and
// write a JSON File to Disk by calling os.WriteFile()
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	perm := os.FileMode(syscall.S_IRUSR | syscall.S_IWUSR)
	err = os.WriteFile(db.path, data, perm)
	if err != nil {
		return err
	}

	return nil
}

// Loads a JSON-DB from Disk by
// handling mutual exclusions, creating an empty DBStructure struct
// to fill with data from Disk, Unmarshalling JSON data from Disk to
// to in-memory DBStructure and returning said DBStructure
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}
