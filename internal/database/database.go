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

// represents the contents of DB as map of Chirps
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// represents a Chirp as part of the map of Chirps inside DB
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

/*
	DB STRUCT AS SAVED FILE ON DISK:
	FILE PATH: "path/to/file/database.json"

	INSIDE JSON:
	{	<-- DBStructure struct
		"chirps": {	<-- DBStructure.Chirps map[int]string
			"1": { id: 1, body: "blabla" },		<-- Chirp struct
			"2": { id: 2, body: "blublub" },	<-- Chirp struct
		}
	}
*/

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
	}
	return db.writeDB(dbStructure)
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

// CHIRP IN DB

// Creates a Chirp by loading the whole JSON-DB in-memory,
// determine and setting Chirp.ID by incrementing the length of all
// saved Chirps by 1, setting Chirp.Body with provided string,
// add new Chirp to in-memory DBStructure.Chirps and
// write the updated in-memory JSON-DB back to disk via DB.writeDB()
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// Reads all Chirps in JSON-DB by loading the whole DB in-memory
// and append all Chirps in a []Chirp Slice
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

// NEW DB FOR IN-MEMORY ON SERVER START
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}
