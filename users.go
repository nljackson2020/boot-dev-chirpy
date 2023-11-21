package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type UserParams struct {
	Email string `json:"email"`
}

func handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := UserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	db, err := NewDB("./database.json")
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side creating database")
	}

	dbErr := db.ensureDB()
	if dbErr != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Database not created")
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side accessing database")
		return
	}

	maxID := 0
	for id := range dbStructure.User {
		if id > maxID {
			maxID = id
		}
	}

	newUser := User{
		Email: params.Email,
		ID:    maxID + 1,
	}

	if dbStructure.User == nil {
		dbStructure.User = make(map[int]User)
	}

	dbStructure.User[newUser.ID] = newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error writing data to database")
		return
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}
