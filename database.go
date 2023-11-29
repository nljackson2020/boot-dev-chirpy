package main

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	User   map[int]User  `json:"users"`
}

func NewDB(path string) (*DB, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	file.Close()

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	return db, nil
}

func (db *DB) ensureDB() error {
	_, err := db.loadDB()
	if err == nil {
		return nil
	}
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
	}
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	file, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	defer file.Close()

	var dbStructure DBStructure
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	file, err := os.OpenFile(db.path, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
