package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
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

	maxID := 0
	for id := range dbStructure.User {
		if id > maxID {
			maxID = id
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't has password")
		return
	}

	newUser := User{
		Password: string(hashedPassword),
		Email:    params.Email,
		ID:       maxID + 1,
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

	newUserResponse := UserResponse{
		Email: newUser.Email,
		ID:    newUser.ID,
	}

	respondWithJSON(w, http.StatusCreated, newUserResponse)
}
