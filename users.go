package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReturnUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type UserParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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
		return
	}

	dbErr := db.ensureDB()
	if dbErr != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Database not created")
		return
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error on the server side accessing database")
		return
	}

	for _, user := range dbStructure.User {
		if params.Email == user.Email {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
	}

	maxID := 0
	for id := range dbStructure.User {
		if id > maxID {
			maxID = id
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
	}

	newUser := User{
		Email:    params.Email,
		ID:       maxID + 1,
		Password: string(hashedPassword),
	}

	if dbStructure.User == nil {
		dbStructure.User = make(map[int]User)
	}

	for _, user := range dbStructure.User {
		if newUser.Email == user.Email {
			respondWithError(w, http.StatusConflict, "Email already in use")
			return
		}
	}

	dbStructure.User[newUser.ID] = newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Error writing data to database")
		return
	}

	returnUser := ReturnUser{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, returnUser)
}
